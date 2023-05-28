package garbanzo

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	router := NewRouter()
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	req, _ := http.NewRequest("GET", testServer.URL+"/hello", nil)

	client := new(http.Client)
	resp, _ := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "hello", string(respBody))
}
