package garbanzo

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type wsClient struct {
	// socketはこのクライアントのためのWebSocket
	socket *websocket.Conn
	// sendはイベントが送られるチャネル。WebSocketを通じてユーザのブラウザに送られるのを待機する
	send chan *Event
	// 統計
	stats chan *Stats
	// roomはこのクライアントが参加している接続
	room *room
	// doneはクライアントに送信済みの通知IDを保持する
	done map[string]bool
	mu   *sync.RWMutex
}

type mark struct {
	ID  string
	URL string // proxy URL
}

// 無限ループでwebsocketを受信し続ける
func (wsc *wsClient) read() {
	for {
		var m *mark
		if err := wsc.socket.ReadJSON(&m); err == nil {
			wsc.room.markRead <- m
		} else {
			// 読み込めないと終了
			// このループを抜けるとハンドラの実行が終了する。deferによってleaveチャンネルに送られ、送信対象から外される
			break
		}
	}
	wsc.socket.Close()
}

// 直近〜分だけブラウザ通知する
const notifyMinutesAgo = 60

// clientへのメッセージをwebsocketに書き込む
func (wsc *wsClient) write() {
	go func() {
		for stats := range wsc.stats {
			wsc.mu.Lock()
			err := wsc.socket.WriteJSON(stats)
			wsc.mu.Unlock()
			if err != nil {
				break
			}
		}
		wsc.socket.Close()
	}()

	go func() {
		for send := range wsc.send {
			// doneに存在しないときだけ書き込み
			wsc.mu.RLock()
			_, exists := wsc.done[send.NotificationID]
			wsc.mu.RUnlock()
			if exists {
				continue
			}

			// 直近のイベントだけブラウザ通知する
			now := time.Now()
			minutesAgo := now.Add(-notifyMinutesAgo * time.Minute)
			// 「更新時間」が、「更新時刻よりN分前」より未来にあるか?
			// (過去) ---> 今-N分前 ---> |-> 通知有効期間 <-| ---> 今 ---> (未来)
			if send.UpdatedAt.After(minutesAgo) {
				send.IsNotifyBrowser = true
			}

			wsc.mu.Lock()
			err := wsc.socket.WriteJSON(send)
			wsc.mu.Unlock()
			if err != nil {
				break
			}
			wsc.mu.Lock()
			wsc.done[send.NotificationID] = true
			wsc.mu.Unlock()
		}
		wsc.socket.Close()
	}()
}
