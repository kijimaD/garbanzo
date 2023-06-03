package garbanzo

import (
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

func NewProxyRouter() *echo.Echo {
	e := echo.New()
	e.GET("/*", ghHandler)

	return e
}

var proxyCache = make(map[string]string)
var proxyMutex = &sync.RWMutex{}

func ghHandler(c echo.Context) error {
	path := c.Request().URL.String()
	url := "https://github.com" + path

	// load cache
	proxyMutex.RLock()
	val, ok := proxyCache[url]
	proxyMutex.RUnlock()
	if ok {
		return c.HTML(http.StatusOK, val)
	}
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	// save cache
	proxyMutex.Lock()
	proxyCache[url] = string(byteArray)
	proxyMutex.Unlock()
	return c.HTML(http.StatusOK, string(byteArray))
}
