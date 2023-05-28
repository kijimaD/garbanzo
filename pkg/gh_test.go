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
	assert.Equal(t, true, len(Evs) > 0)
}
