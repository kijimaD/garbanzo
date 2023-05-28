package garbanzo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNotifications(t *testing.T) {
	gh := newGitHub()
	err := gh.getNotifications()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, true, len(notifications) > 0)
}

func TestProcess(t *testing.T) {
	gh := newGitHub()
	err := gh.getNotifications()
	if err != nil {
		t.Error(err)
	}
	err = gh.processNotification()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, true, len(events) > 0)
}
