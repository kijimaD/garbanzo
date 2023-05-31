package garbanzo

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kijimaD/garbanzo/trace"
	"github.com/labstack/echo/v4"
)

type room struct {
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
}

func newRoom() *room {
	return &room{
		forward:   make(chan *Event),
		join:      make(chan *wsClient),
		leave:     make(chan *wsClient),
		wsClients: make(map[*wsClient]bool),
		tracer:    trace.Off(), // デフォルトではログ出力はされない
		events:    make(Events),
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) run() {
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
		case forward := <-r.forward:
			for wsClient := range r.wsClients {
				select {
				case wsClient.send <- forward:
					// メッセージを送信
					r.tracer.Trace("send message to client: " + forward.NotificationID)
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
	}

	r.join <- wsc
	defer func() { r.leave <- wsc }()
	go wsc.write() // c.sendの内容をwebsocketに書き込む

	// 初回実行
	go func() {
		for _, v := range r.events {
			r.forward <- v
		}
	}()

	// キャッシュ保存
	go func() {
		for _, v := range r.events {
			resp, _ := http.Get(v.HTMLURL)
			defer resp.Body.Close()
			byteArray, _ := ioutil.ReadAll(resp.Body)
			proxyCache[v.HTMLURL] = string(byteArray)
			time.Sleep(time.Second * 1)
		}
	}()

	wsc.read() // 接続は保持され、終了を指示されるまで他の処理をブロックする
	return nil
}

func (r *room) initEvent() error {
	s := newStore()
	gh := newGitHub()
	err := gh.getNotifications(s)
	if err != nil {
		panic(err)
	}
	err = gh.processNotification(s)
	if err != nil {
		panic(err)
	}
	for k, v := range s.events {
		r.events[k] = v
	}
	return nil
}
