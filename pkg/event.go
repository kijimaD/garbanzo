package garbanzo

import "time"

type notificationID string

// Eventはフロント側で1つ1つの通知表示に必要な項目
type Event struct {
	NotificationID string
	UserName       string
	AvatarURL      string
	Title          string
	Body           string
	LatestURL      string
	When           time.Time
}

type Events map[notificationID]*Event

func newEvent(notificationID string, userName string, avatarURL string, title string, body string, latestURL string, when time.Time) *Event {
	return &Event{
		NotificationID: notificationID,
		UserName:       userName,
		AvatarURL:      avatarURL,
		Title:          title,
		Body:           body,
		LatestURL:      latestURL,
		When:           when,
	}
}
