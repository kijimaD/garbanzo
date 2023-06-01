package garbanzo

import (
	"sync"

	"github.com/gorilla/websocket"
)

type wsClient struct {
	// socketはこのクライアントのためのWebSocket
	socket *websocket.Conn
	// sendはメッセージが送られるチャネル。WebSocketを通じてユーザのブラウザに送られるのを待機する
	send chan *Event
	done map[string]bool
}

// 無限ループで待機
func (wsc *wsClient) read() {
	for {
	}
	wsc.socket.Close()
}

// c.sendの内容をwebsocketに書き込む
func (wsc *wsClient) write() {
	mu := &sync.Mutex{}
	for send := range wsc.send {
		// doneに存在しないときだけ書き込み
		if _, exist := wsc.done[send.NotificationID]; exist {
			continue
		}
		err := wsc.socket.WriteJSON(send)
		if err != nil {
			break
		}
		mu.Lock()
		wsc.done[send.NotificationID] = true
		mu.Unlock()
	}
	wsc.socket.Close()
}
