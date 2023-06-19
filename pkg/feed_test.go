package garbanzo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDumpFeedsTable(t *testing.T) {
	c := NewConfig(".")
	b := []byte(`
- desc: RFC1
  url: https://www.rfc-editor.org/rfcrss.xml
- desc: RFC2
  url: https://www.rfc-editor.org/rfcrss.xml
`)
	ss := c.loadFeedSources(b)
	result := ss.dumpFeedsTable()
	expect := `| Description |                 Feed                  |
|-------------|---------------------------------------|
| RFC1        | https://www.rfc-editor.org/rfcrss.xml |
| RFC2        | https://www.rfc-editor.org/rfcrss.xml |
`
	assert.Equal(t, expect, result)
}
