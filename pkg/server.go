package garbanzo

import (
	"embed"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	trace "github.com/kijimaD/garbanzo/trace"
	"github.com/labstack/echo/v4"
)

//go:embed templates
var fst embed.FS

//go:embed static
var fss embed.FS

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	data = map[string]interface{}{
		"Host": c.Request().Host,
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewRouter(c *Config, templDir string, publicDir string) *echo.Echo {
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseFS(fst, templDir)),
	}

	room := newRoom()
	room.tracer = trace.New(os.Stdout)
	room.config = c
	go room.fetchEvent()                                            // 初回実行
	go room.fetchFeeds()                                            // 初回実行
	go func() { time.Sleep(10 * time.Second); room.fetchCache() }() // 初回実行
	go room.run()

	e := echo.New()
	e.Renderer = renderer
	e.GET("/", rootHandler)
	e.GET("/events", room.eventHandler)
	// TODO: static系を1ハンドラにまとめる
	e.GET("/favicon.ico", faviconHandler)
	e.GET("/rssicon", rssIconHandler)
	return e
}

func rootHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "root.html", nil)
}

func faviconHandler(c echo.Context) error {
	data, err := fss.ReadFile("static/favicon.ico")
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "image/x-icon", data)
}

func rssIconHandler(c echo.Context) error {
	data, err := fss.ReadFile("static/rss.svg")
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "image/svg+xml", data)
}
