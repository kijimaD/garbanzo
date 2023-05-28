package garbanzo

// Eventはフロント側で1つ1つの通知表示に必要な項目
type Event struct {
	UserName  string
	AvatarURL string
	Title     string
	Body      string
	LatestURL string
}

type Events []*Event

func newEvent(userName string, avatarURL string, title string, body string, latestURL string) *Event {
	return &Event{
		UserName:  userName,
		AvatarURL: avatarURL,
		Title:     title,
		Body:      body,
		LatestURL: latestURL,
	}
}
