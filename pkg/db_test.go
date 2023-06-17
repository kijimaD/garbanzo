package garbanzo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppDirPath(t *testing.T) {
	c := NewConfig(".")
	assert.Equal(t, ".garbanzo", c.appDirPath())

	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{
			name:   "dot",
			input:  ".",
			expect: ".garbanzo",
		},
		{
			name:   "dot slash",
			input:  "./",
			expect: ".garbanzo",
		},
		{
			name:   "test",
			input:  "test",
			expect: "test/.garbanzo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfig(tt.input)
			assert.Equal(t, tt.expect, c.appDirPath())
		})
	}
}

func TestSaveFilePath(t *testing.T) {
	c := NewConfig(".")
	assert.Equal(t, ".garbanzo/mark.csv", c.saveFilePath())
}

func TestPutDir(t *testing.T) {
	c := NewConfig(".")
	c.putConfDir()
	c.putConfDir()
	defer os.RemoveAll(".garbanzo")
}

func TestMarkToFile(t *testing.T) {
	c := NewConfig(".")
	c.putConfDir()
	c.markToFile("http://example.com")
	defer os.RemoveAll(".garbanzo")
}

func TestIsMark(t *testing.T) {
	c := NewConfig(".")
	c.putConfDir()
	c.markToFile("http://example.com")
	assert.Equal(t, false, c.isMarked("NOT_EXISTS"))
	assert.Equal(t, true, c.isMarked("http://example.com"))
	defer os.RemoveAll(".garbanzo")
}
