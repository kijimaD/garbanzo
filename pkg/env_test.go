package garbanzo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppBase(t *testing.T) {
	e := Env{
		AppHost:     "http://localhost",
		AppPort:     8080,
		ProxyHost:   "http://localhost",
		ProxyPort:   8081,
		GitHubToken: "",
	}
	assert.Equal(t, "http://localhost:8080", e.appBase())
	assert.Equal(t, "http://localhost:8081", e.proxyBase())
}
