package garbanzo

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootHandler(t *testing.T) {
	c := NewConfig(".")
	c.PutConfDir()
	defer os.RemoveAll(".garbanzo")

	router := NewRouter(c, "templates/*.html", "static/*")
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
	c := NewConfig(".")
	c.PutConfDir()
	defer os.RemoveAll(".garbanzo")

	router := NewRouter(c, "templates/*.html", "static/*")
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	req, _ := http.NewRequest("GET", testServer.URL+"/favicon.ico", nil)

	client := new(http.Client)
	resp, _ := client.Do(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
