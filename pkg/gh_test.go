//go:build gh

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
	assert.Equal(t, true, len(gh.notifications) > 0)
}

func TestGetNotificationsDup(t *testing.T) {
	gh := newGitHub()
	err := gh.getNotifications()
	if err != nil {
		t.Error(err)
	}

	// 同じ通知は追加しない
	count := len(gh.notifications)
	gh.getNotifications()
	gh.getNotifications()
	assert.Equal(t, count, len(gh.notifications))
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
	assert.Equal(t, true, len(gh.events) > 0)
	assert.Equal(t, len(gh.events), len(gh.notifications))
}
