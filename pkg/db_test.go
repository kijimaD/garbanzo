package garbanzo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMark(t *testing.T) {
	c := NewConfig(".")
	c.PutConfDir()
	c.markToFile("http://ismark-example.com")
	assert.Equal(t, false, c.isMarked("NOT_EXISTS"))
	assert.Equal(t, true, c.isMarked("http://ismark-example.com"))
	defer os.RemoveAll(".garbanzo")
}
