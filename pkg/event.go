package garbanzo

import "time"

// Eventはフロント側で1つ1つの通知表示に必要な項目
// 将来的にGitHubだけじゃなくなるかもなので、汎用的にしておく
type Event struct {
	NotificationID string    // 通知ID
	UserName       string    // GitHubユーザの名前
	AvatarURL      string    // GitHubユーザのアイコン画像
	Title          string    // 通知のタイトル
	Body           string    // 通知の本文
	HTMLURL        string    // iframe遷移に使う、ホストがリバースプロキシ先で置き換えられたURL
	RepoName       string    // フルリポジトリ名 golang/go
	When           time.Time // 更新時間(updated_at)
}

type Events map[string]*Event // keyにNotificationIDを使う

func newEvent(notificationID string, userName string, avatarURL string, title string, body string, HTMLURL string, repoName string, when time.Time) *Event {
	return &Event{
		NotificationID: notificationID,
		UserName:       userName,
		AvatarURL:      avatarURL,
		Title:          title,
		Body:           body,
		HTMLURL:        HTMLURL,
		RepoName:       repoName,
		When:           when,
	}
}
