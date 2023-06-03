package garbanzo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	// tracerは操作のログを受け取る
	tracer trace.Tracer
	events Events
	mu     *sync.RWMutex
}

func newRoom() *room {
	return &room{
		fetch:     make(chan *Event),
		forward:   make(chan *Event),
		join:      make(chan *wsClient),
		leave:     make(chan *wsClient),
		wsClients: make(map[*wsClient]bool),
		tracer:    trace.Off(), // デフォルトではログ出力はされない
		events:    make(Events),
		mu:        &sync.RWMutex{},
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
		for {
			select {
			case <-t1.C:
				go func() {
					// r.eventsをクライアントと同期する
					r.mu.RLock()
					for _, v := range r.events {
						r.forward <- v
					}
					r.mu.RUnlock()
				}()
			case <-t2.C:
				go func() {
					r.fetchEvent()

					err := r.fetchCache()
					if err != nil {
						log.Println(err)
					}
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
		case fetch := <-r.fetch:
			r.mu.Lock()
			r.events[fetch.NotificationID] = fetch
			r.mu.Unlock()
		case forward := <-r.forward:
			for wsClient := range r.wsClients {
				wsClient.mu.RLock()
				exists := wsClient.done[forward.NotificationID]
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

func (r *room) handleWebSocket(c echo.Context) error {
	socket, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return err
	}

	wsc := &wsClient{
		socket: socket,
		send:   make(chan *Event, messageBufferSize),
		done:   make(map[string]bool),
		mu:     &sync.RWMutex{},
	}

	r.join <- wsc
	defer func() { r.leave <- wsc }()
	go wsc.write() // c.sendの内容をwebsocketに書き込む
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
	fmt.Print("e")
	return nil
}

// HTMLページのキャッシュを取得する
func (r *room) fetchCache() error {
	for _, v := range r.events {
		if _, exists := proxyCache[v.ProxyURL]; exists {
			continue
		}
		resp, _ := http.Get(v.ProxyURL)
		defer resp.Body.Close()
		byteArray, _ := ioutil.ReadAll(resp.Body)
		proxyCache[v.ProxyURL] = string(byteArray)
		time.Sleep(time.Second * 1)
		fmt.Print("c")
	}
	return nil
}
