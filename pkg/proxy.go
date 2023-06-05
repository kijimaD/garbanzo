package garbanzo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

func NewProxyRouter() *echo.Echo {
	e := echo.New()
	e.GET("/", homeHandler)
	e.GET("/*", ghHandler)

	return e
}

func homeHandler(c echo.Context) error {
	data, err := fss.ReadFile("static/home.md")
	if err != nil {
		fmt.Println(err)
		return err
	}
	html := string(mdToHTML([]byte(data)))
	return c.HTML(http.StatusOK, html)
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
