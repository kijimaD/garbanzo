package garbanzo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBulidFeedStatus(t *testing.T) {
	homedir, _ := os.UserHomeDir()
	c := NewConfig(homedir)
	s, _ := buildFeedStatus(c)
	assert.Equal(t, true, len(s) > 0)
}

func TestBulidTokenStatus(t *testing.T) {
	c := NewConfig(".")
	c.PutConfDir()
	defer os.RemoveAll(".garbanzo")
	f, err := os.Create(c.tokenFilePath())
	defer f.Close()
	_, err = f.Write([]byte("THIS IS TOKEN"))
	if err != nil {
		t.Error(err)
	}

	s, _ := buildTokenStatus(c)
	expect := "## GitHub Token\n" + "`~/.garbanzo/token`\n\n" + "🟢 ok"
	assert.Equal(t, expect, s)
}
