package garbanzo

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewProxyRouter() *echo.Echo {
	e := echo.New()
	e.GET("/*", ghHandler)

	return e
}

var proxyCache map[string]string

func ghHandler(c echo.Context) error {
	if proxyCache == nil {
		proxyCache = make(map[string]string)
	}

	path := c.Request().URL.String()
	url := "https://github.com" + path

	// load cache
	if val, ok := proxyCache[url]; ok {
		return c.HTML(http.StatusOK, val)
	}
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	// save cache
	proxyCache[url] = string(byteArray)
	return c.HTML(http.StatusOK, string(byteArray))
}
