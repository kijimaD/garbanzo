package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/kelseyhightower/envconfig"
	garbanzo "github.com/kijimaD/garbanzo/pkg"
)

// TODO: 別パッケージにもある。一箇所にまとめたい
type Env struct {
	AppPort   uint16 `envconfig:"APP_PORT" default:"8080"`
	ProxyPort uint16 `envconfig:"PROXY_PORT" default:"8081"`
}

func main() {
	var env Env
	err := envconfig.Process("", &env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't parse environment variables: %s\n", err.Error())
		os.Exit(1)
	}

	go func() {
		router := garbanzo.NewRouter("templates/*.html", "static/*")
		router.Start(":" + strconv.FormatUint(uint64(env.AppPort), 10))
	}()

	go func() {
		proxyrouter := garbanzo.NewProxyRouter()
		proxyrouter.Start(":" + strconv.FormatUint(uint64(env.ProxyPort), 10))
	}()

	runtime.Goexit()
}
