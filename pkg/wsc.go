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

// 無限ループで待機
func (wsc *wsClient) read() {
	for {
	}
	wsc.socket.Close()
}

// c.sendの内容をwebsocketに書き込む
func (wsc *wsClient) write() {
	for send := range wsc.send {
		if err := wsc.socket.WriteJSON(send); err != nil {
			break
		}
	}
	wsc.socket.Close()
}
