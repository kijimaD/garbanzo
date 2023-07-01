package garbanzo

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Proxy Server.
func TestHomeHandler(t *testing.T) {
	c := NewConfig(".")
	c.PutConfDir()
	defer os.RemoveAll(".garbanzo")

	router := NewProxyRouter(c)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	req, _ := http.NewRequest("GET", testServer.URL, nil)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(respBody), "Garbanzo is fast notification viewer!")
}

func TestGHHandler(t *testing.T) {
	c := NewConfig(".")
	c.PutConfDir()
	defer os.RemoveAll(".garbanzo")

	router := NewProxyRouter(c)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	req, _ := http.NewRequest("GET", testServer.URL+"/kijimaD?origin=github.com", nil)

	client := new(http.Client)
	resp, _ := client.Do(req)
	respBody, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(respBody), "<title>kijimaD (Kijima Daigo) · GitHub</title>")
}
