package garbanzo

import "time"

// Eventはフロント側で1つ1つの通知表示に必要な項目
// 将来的にGitHubだけじゃなくなるかもなので、汎用的にしておく
type Event struct {
	NotificationID  string    // 通知ID
	UserName        string    // GitHubユーザの名前
	AvatarURL       string    // GitHubユーザのアイコン画像
	Title           string    // 通知のタイトル(プレーンテキスト)
	TitleHTML       string    // 通知のタイトル(HTML)
	Body            string    // 通知の本文
	BodyHTML        string    // 通知の本文(HTML)
	HTMLURL         string    // 通常の、ホストがgithub.comのURL
	ProxyURL        string    // iframe遷移に使う、ホストがリバースプロキシ先で置き換えられたURL
	RepoName        string    // フルリポジトリ名 golang/go
	When            string    // 更新時間(updated_at) // TODO: リネームする
	Kind            string    // 種別
	UpdatedAt       time.Time // 更新時刻
	IsNotifyBrowser bool      // ブラウザ通知するかどうか
}

type Events map[string]*Event // keyにNotificationIDを使う

func newEvent(notificationID string, userName string, avatarURL string, title string, titleHTML string, body string, bodyHTML, HTMLURL string, ProxyURL string, repoName string, when string, kind string, updatedAt time.Time) *Event {
	return &Event{
		NotificationID:  notificationID,
		UserName:        userName,
		AvatarURL:       avatarURL,
		Title:           title,
		TitleHTML:       titleHTML,
		Body:            body,
		BodyHTML:        bodyHTML,
		HTMLURL:         HTMLURL,
		ProxyURL:        ProxyURL,
		RepoName:        repoName,
		When:            when,
		Kind:            kind,
		UpdatedAt:       updatedAt,
		IsNotifyBrowser: false,
	}
}

type Stats struct {
	ReadCount  int
	EventCount int
	CacheCount int
}

func newStats() *Stats {
	return &Stats{
		ReadCount:  0,
		EventCount: 0,
		CacheCount: 0,
	}
}
