package garbanzo

import (
	"testing"
)

func TestNewGitHub(t *testing.T) {
	gh := newGitHub()
	err := gh.getNotifications()

	if err != nil {
		t.Error(err)
	}
}
