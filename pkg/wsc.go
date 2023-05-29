package garbanzo

import (
	"github.com/gorilla/websocket"
)

type wsClient struct {
	// socketはこのクライアントのためのWebSocket
	socket *websocket.Conn
	// sendはメッセージが送られるチャネル。WebSocketを通じてユーザのブラウザに送られるのを待機する
	send chan *Event
}

func (wsc *wsClient) read() {
	for {
	}
	wsc.socket.Close()
}

// c.sendの内容をwebsocketに書き込む
func (wsc *wsClient) write() {
	for msg := range wsc.send {
		if err := wsc.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	wsc.socket.Close()
}
