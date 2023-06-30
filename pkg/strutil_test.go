package garbanzo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBulidHomeMD(t *testing.T) {
	homedir, _ := os.UserHomeDir()
	c := NewConfig(homedir)
	s, _ := buildTokenStatus(c)
	assert.Equal(t, true, len(s) > 0)
}

func TestBulidFeedStatus(t *testing.T) {
	homedir, _ := os.UserHomeDir()
	c := NewConfig(homedir)
	s, _ := buildFeedStatus(c)
	assert.Equal(t, true, len(s) > 0)
}
