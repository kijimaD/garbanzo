package garbanzo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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
	md := string(data)

	homedir, _ := os.UserHomeDir()
	conf := NewConfig(homedir)
	b, _ := os.ReadFile(conf.feedFilePath())
	fss := conf.loadFeedSources(b)
	md = md + fss.dumpFeedsTable()

	html := string(mdToHTML([]byte(md)))
	return c.HTML(http.StatusOK, html)
}

var proxyCache = make(map[string]string)
var proxyMutex = &sync.RWMutex{}

func ghHandler(c echo.Context) error {
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

	byteArray, _ := ioutil.ReadAll(resp.Body)
	// save cache
	proxyMutex.Lock()
	proxyCache[u] = string(byteArray)
	proxyMutex.Unlock()
	return c.HTML(http.StatusOK, string(byteArray))
}
