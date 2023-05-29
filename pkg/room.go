package garbanzo

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type room struct {
	// forwardはクライアントに転送するためのメッセージを保持するためのチャネル
	forward chan *Event
	// clientsには在室しているすべてのクライアントが保持される
	clients map[*wsClient]bool
}

func newRoom() *room {
	return &room{
		forward: make(chan *Event),
		clients: make(map[*wsClient]bool),
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
		case msg := <-r.forward:
			for client := range r.clients {
				select {
				case client.send <- msg:
					// メッセージを送信
				default:
					// 送信に失敗
					delete(r.clients, client)
					close(client.send)
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

	go wsc.write() // c.sendの内容をwebsocketに書き込む
	wsc.read()     // 接続は保持され、終了を指示されるまで他の処理をブロックする

	return nil
}
