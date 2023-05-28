package garbanzo

import (
	"html/template"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewRouter(templDir string) *echo.Echo {
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob(templDir)),
	}

	e := echo.New()
	e.Renderer = renderer
	e.GET("/", rootHandler)

	return e
}

func rootHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "root.html", nil)
}

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
