package garbanzo

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootHandler(t *testing.T) {
	router := NewRouter("templates/*.html", "static/*")
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	req, _ := http.NewRequest("GET", testServer.URL, nil)

	client := new(http.Client)
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(respBody), "Garbanzo")
}

func TestFaviconHandler(t *testing.T) {
	router := NewRouter("templates/*.html", "static/*")
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	req, _ := http.NewRequest("GET", testServer.URL+"/favicon.ico", nil)

	client := new(http.Client)
	resp, _ := client.Do(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Proxy Server.
func TestHomeHandler(t *testing.T) {
	router := NewProxyRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	req, _ := http.NewRequest("GET", testServer.URL, nil)

	client := new(http.Client)
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(respBody), "Garbanzo is fast GitHub viewer!")
}

func TestGHHandler(t *testing.T) {
	router := NewProxyRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	req, _ := http.NewRequest("GET", testServer.URL+"/kijimaD", nil)

	client := new(http.Client)
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(respBody), "GitHub")
}
