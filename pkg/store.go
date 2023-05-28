package garbanzo

import "github.com/google/go-github/v48/github"

type store struct {
	notifications map[notificationID]*github.Notification
	events        Events
}

func newStore() *store {
	s := store{
		notifications: make(map[notificationID]*github.Notification),
		events:        make(Events),
	}
	return &s
}
