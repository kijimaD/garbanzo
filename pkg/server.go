package garbanzo

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	e.GET("/", rootHandler)

	return e
}

func rootHandler(c echo.Context) error {
	return c.String(http.StatusOK, "this is root")
}
