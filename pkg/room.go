package garbanzo

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kijimaD/garbanzo/trace"
	"github.com/labstack/echo/v4"
)

type room struct {
	// fetchはGitHubから取得してきたデータを保持するためのチャネル
	fetch chan *Event
	// forwardはクライアントに転送するためのメッセージを保持するためのチャネル
	forward chan *Event
	// joinはクライアントの接続要求を保持するためのチャネル
	join chan *wsClient
	// leaveは切断しようとしているクライアントのためのチャネル
	leave chan *wsClient
	// wsClientsには接続しているすべてのクライアントが保持される
	wsClients map[*wsClient]bool
	// 既読にしようとしている通知IDを保持するためのチャネル
	markRead   chan *mark
	stats      chan *Stats
	statsStore *Stats
	// tracerは操作のログを受け取る
	tracer trace.Tracer
	events Events
	mu     *sync.RWMutex
}

func newRoom() *room {
	return &room{
		fetch:      make(chan *Event),
		forward:    make(chan *Event),
		join:       make(chan *wsClient),
		leave:      make(chan *wsClient),
		wsClients:  make(map[*wsClient]bool),
		markRead:   make(chan *mark),
		stats:      make(chan *Stats, 1), // 同じゴルーチン内で送信と受信をするため、容量ゼロだと止まってしまう
		statsStore: newStats(),
		tracer:     trace.Off(), // デフォルトではログ出力はされない
		events:     make(Events),
		mu:         &sync.RWMutex{},
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

const syncSecond = 4

func (r *room) run() {
	go func() {
		t1 := time.NewTicker(syncSecond * time.Second)
		defer t1.Stop()
		t2 := time.NewTicker(30 * time.Second)
		defer t2.Stop()
		t3 := time.NewTicker(2 * time.Second)
		defer t3.Stop()
		for {
			select {
			case <-t1.C:
				go func() {
					// r.eventsをクライアントと同期する
					// 一旦ソートしたスライスを作成して送信する
					keys := make([]string, 0, len(r.events))
					r.mu.RLock()
					for key := range r.events {
						keys = append(keys, key)
					}
					sort.SliceStable(keys, func(i, j int) bool {
						return r.events[keys[i]].When < r.events[keys[j]].When
					})
					sorted := make([]*Event, 0, len(r.events))
					for _, k := range keys {
						sorted = append(sorted, r.events[k])
					}
					r.mu.RUnlock()
					for _, v := range sorted {
						r.forward <- v
					}
				}()
			case <-t2.C:
				go func() {
					r.fetchEvent()

					err := r.fetchCache()
					if err != nil {
						log.Println(err)
					}
				}()
			case <-t3.C:
				go func() {
					r.mu.RLock()
					r.statsStore.EventCount = len(r.events)
					r.mu.RUnlock()

					proxyMutex.RLock()
					r.statsStore.CacheCount = len(proxyCache)
					proxyMutex.RUnlock()
				}()
			}
		}
	}()

	for {
		select {
		case wsClient := <-r.join:
			// 参加
			r.wsClients[wsClient] = true
			r.tracer.Trace("join client")
		case wsClient := <-r.leave:
			// 退室
			delete(r.wsClients, wsClient)
			close(wsClient.send)
			r.tracer.Trace("leave client")
		case mark := <-r.markRead:
			err := r.markThreadRead(mark.ID)
			if err != nil {
				log.Println("mark thread read err:", err)
			}
			delete(r.events, mark.ID)
			delete(proxyCache, mark.URL)
			r.statsStore.ReadCount++
			r.stats <- r.statsStore
		case stats := <-r.stats:
			for wsClient := range r.wsClients {
				select {
				case wsClient.stats <- stats:
					r.tracer.Trace("send stats to client")
				}
			}
		case fetch := <-r.fetch:
			r.mu.Lock()
			r.events[fetch.NotificationID] = fetch
			r.mu.Unlock()
			r.stats <- r.statsStore
		case forward := <-r.forward:
			for wsClient := range r.wsClients {
				wsClient.mu.RLock()
				_, exists := wsClient.done[forward.NotificationID]
				wsClient.mu.RUnlock()
				if exists {
					continue
				}
				select {
				case wsClient.send <- forward:
					// roomごとのEventsをwsClientごとのEventsと同期する
					r.tracer.Trace("send message to client")
				default:
					// 送信に失敗
					delete(r.wsClients, wsClient)
					close(wsClient.send)
				}
			}
		}
	}
}

func (r *room) eventHandler(c echo.Context) error {
	socket, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return err
	}

	wsc := &wsClient{
		socket: socket,
		send:   make(chan *Event, messageBufferSize),
		stats:  make(chan *Stats, messageBufferSize),
		room:   r,
		done:   make(map[string]bool),
		mu:     &sync.RWMutex{},
	}

	r.join <- wsc
	r.stats <- r.statsStore
	defer func() { r.leave <- wsc }()
	go wsc.write() // websocketに書き込む
	wsc.read()     // このメソッドの中の無限ループによって接続は保持され、終了を指示されるまで他の処理をブロックする
	return nil
}

func (r *room) fetchEvent() error {
	gh := newGitHub()
	err := gh.getNotifications()
	if err != nil {
		return err
	}
	err = gh.processNotification(r)
	if err != nil {
		return err
	}
	return nil
}

func (r *room) markThreadRead(id string) error {
	gh := newGitHub()
	ctx := context.Background()
	_, err := gh.Client.Activity.MarkThreadRead(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

// HTMLページのキャッシュを取得する
func (r *room) fetchCache() error {
	for _, v := range r.events {
		proxyMutex.RLock()
		_, exists := proxyCache[v.ProxyURL]
		proxyMutex.RUnlock()
		if exists {
			continue
		}

		resp, _ := http.Get(v.ProxyURL)
		defer resp.Body.Close()
		byteArray, _ := ioutil.ReadAll(resp.Body)
		proxyMutex.Lock()
		proxyCache[v.ProxyURL] = string(byteArray)
		proxyMutex.Unlock()

		time.Sleep(time.Second * 1)
	}
	return nil
}
