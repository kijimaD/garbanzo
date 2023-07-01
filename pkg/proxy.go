package garbanzo

import (
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/labstack/echo/v4"
)

var proxyCache = make(map[string]string)
var proxyMutex = &sync.RWMutex{}

type proxyServ struct {
	config *Config
}

func NewProxyRouter(c *Config) *echo.Echo {
	ps := proxyServ{config: c}
	e := echo.New()
	e.GET("/", ps.homeHandler)
	e.GET("/*", ps.ghHandler)

	return e
}

func (p *proxyServ) homeHandler(c echo.Context) error {
	md, err := buildHomeMD(p.config)
	if err != nil {
		return err
	}
	html := string(mdToHTML([]byte(md)))
	return c.HTML(http.StatusOK, html)
}

func (p *proxyServ) ghHandler(c echo.Context) error {
	var u string
	reqpath := c.Request().URL.String()
	h, err := url.Parse(reqpath)
	if err != nil {
		return err
	}
	originHost := h.Query()["origin"]
	if len(originHost) == 1 {
		u = "https://" + originHost[0] + reqpath
	} else {
		// FIXME: iframe内で開くのが相対リンクの場合、暗黙的にホストがlocalshostになる。ホストがわからないから、元のページを開けない。とりあえずgithub.comにしておく...
		u = "https://github.com" + reqpath
	}

	// load cache
	proxyMutex.RLock()
	val, ok := proxyCache[u]
	proxyMutex.RUnlock()
	if ok {
		return c.HTML(http.StatusOK, val)
	}
	resp, _ := http.Get(u)
	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)
	// save cache
	proxyMutex.Lock()
	proxyCache[u] = string(byteArray)
	proxyMutex.Unlock()
	return c.HTML(http.StatusOK, string(byteArray))
}
