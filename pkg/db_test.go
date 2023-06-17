package garbanzo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutDir(t *testing.T) {
	putConfDir("./")
	putConfDir("./")
	defer os.RemoveAll("./.garbanzo/")
}

func TestMarkToFile(t *testing.T) {
	putConfDir("./")
	markToFile("http://example.com")
	defer os.RemoveAll("./.garbanzo/")
}

func TestIsMark(t *testing.T) {
	putConfDir("./")
	markToFile("http://example.com")
	assert.Equal(t, false, isMarked("NOT_EXISTS"))
	assert.Equal(t, true, isMarked("http://example.com"))
	defer os.RemoveAll("./.garbanzo/")
}
