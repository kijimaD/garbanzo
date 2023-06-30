package garbanzo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMark(t *testing.T) {
	c := NewConfig(".")
	c.PutConfDir()
	fdb := newfdb(c)
	fdb.markToFile("http://ismark-example.com")
	assert.Equal(t, false, fdb.isMarked("NOT_EXISTS"))
	assert.Equal(t, true, fdb.isMarked("http://ismark-example.com"))
	defer os.RemoveAll(".garbanzo")
}
