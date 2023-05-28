//go:build gh

package garbanzo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNotifications(t *testing.T) {
	s := newStore()
	gh := newGitHub()
	err := gh.getNotifications(s)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, true, len(s.notifications) > 0)
}

func TestGetNotificationsDup(t *testing.T) {
	s := newStore()
	gh := newGitHub()
	err := gh.getNotifications(s)
	if err != nil {
		t.Error(err)
	}

	// 同じ通知は追加しない
	count := len(s.notifications)
	gh.getNotifications(s)
	gh.getNotifications(s)
	assert.Equal(t, count, len(s.notifications))
}

func TestProcess(t *testing.T) {
	s := newStore()
	gh := newGitHub()
	err := gh.getNotifications(s)
	if err != nil {
		t.Error(err)
	}
	err = gh.processNotification(s)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, true, len(s.events) > 0)
	assert.Equal(t, len(s.events), len(s.notifications))
}
