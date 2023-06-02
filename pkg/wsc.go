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
		var event *Event
		if err := wsc.socket.ReadJSON(&event); err == nil {
		} else {
			// 読み込めないと終了
			// このループを抜けるとハンドラの実行が終了する。deferによってleaveチャンネルに送られ、送信対象から外される
			break
		}
	}
	wsc.socket.Close()
}

// c.sendの内容をwebsocketに書き込む
func (wsc *wsClient) write() {
	mu := &sync.RWMutex{}
	for send := range wsc.send {
		// doneに存在しないときだけ書き込み
		mu.RLock()
		exists := wsc.done[send.NotificationID]
		mu.RUnlock()
		if exists {
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
