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
	c.PutConfDir()
	c.PutConfDir()
	defer os.RemoveAll(".garbanzo")

	sb, err := os.ReadFile(c.saveFilePath())
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "# marked list\n", string(sb))

	fb, err := os.ReadFile(c.feedFilePath())
	if err != nil {
		t.Error(err)
	}

	expect := "# feed list\n- desc: Zenn Go\n  url: https://zenn.dev/topics/go/feed\n- desc: oreilly ebook soon\n  url: https://www.oreilly.co.jp/ebook/new_release.atom\n- desc: Russ Cox blog\n  url: https://research.swtch.com/feed.atom\n"
	assert.Equal(t, expect, string(fb))
}

func TestLoadFeedSources(t *testing.T) {
	c := NewConfig(".")
	b := []byte(`
- desc: RFC1
  url: https://www.rfc-editor.org/rfcrss.xml
- desc: RFC2
  url: https://www.rfc-editor.org/rfcrss.xml
`)
	ss := c.loadFeedSources(b)
	assert.Equal(t, "RFC1", ss[0].Desc)
	assert.Equal(t, "RFC2", ss[1].Desc)
	assert.Equal(t, "https://www.rfc-editor.org/rfcrss.xml", ss[0].URL)
	assert.Equal(t, "https://www.rfc-editor.org/rfcrss.xml", ss[1].URL)
	assert.Equal(t, 2, len(ss))
}
