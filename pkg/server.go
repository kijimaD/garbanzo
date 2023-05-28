package garbanzo

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	e.GET("/hello", helloHandler)

	return e
}

func helloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "hello")
}
