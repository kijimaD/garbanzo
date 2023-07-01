package garbanzo

import "strconv"

var Envar Env

type Env struct {
	AppPort     uint16 `envconfig:"APP_PORT" default:"8080"`
	ProxyHost   string `envconfig:"PROXY_BASE" default:"http://localhost"`
	ProxyPort   uint16 `envconfig:"PROXY_PORT" default:"8081"`
	GitHubToken string `envconfig:"GH_TOKEN"`
}

// 環境変数からベースURLを組み立てる
// 例: http://localhost:8080
func (e *Env) proxyBase() string {
	proxyBase := e.ProxyHost + ":" + strconv.FormatUint(uint64(e.ProxyPort), 10)
	return proxyBase
}
