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

func ghHandler(c echo.Context) error {
	path := c.Request().URL.String()
	url := "https://github.com" + path
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	return c.HTML(http.StatusOK, string(byteArray))
}
