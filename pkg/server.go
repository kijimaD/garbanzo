package garbanzo

import (
	"html/template"
	"io"
	"net/http"
	"os"

	trace "github.com/kijimaD/garbanzo/trace"
	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	data = map[string]interface{}{
		"Host": c.Request().Host,
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewRouter(templDir string) *echo.Echo {
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob(templDir)),
	}

	room := newRoom()
	room.tracer = trace.New(os.Stdout)
	room.initEvent()
	go room.run()

	e := echo.New()
	e.Renderer = renderer
	e.GET("/", rootHandler)
	e.GET("/ws", room.handleWebSocket)
	return e
}

func rootHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "root.html", nil)
}
