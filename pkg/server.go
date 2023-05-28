package garbanzo

import (
	"html/template"
	"io"
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
