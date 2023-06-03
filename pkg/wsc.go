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
	mu   *sync.RWMutex
}

// 無限ループでwebsocketを受信し続ける
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
	for send := range wsc.send {
		// doneに存在しないときだけ書き込み
		wsc.mu.RLock()
		exists := wsc.done[send.NotificationID]
		wsc.mu.RUnlock()
		if exists {
			continue
		}
		err := wsc.socket.WriteJSON(send)
		if err != nil {
			break
		}
		wsc.mu.Lock()
		wsc.done[send.NotificationID] = true
		wsc.mu.Unlock()
	}
	wsc.socket.Close()
}
