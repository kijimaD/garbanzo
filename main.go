package main

import (
	"log"
	"os"
	"runtime"
	"strconv"

	garbanzo "github.com/kijimaD/garbanzo/pkg"
)

func main() {
	// 設定ファイルがない場合は作成する
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
	}
	c := garbanzo.NewConfig(homedir)
	c.PutConfDir()

	go func() {
		router := garbanzo.NewRouter(c, "templates/*.html", "static/*")
		router.Start(":" + strconv.FormatUint(uint64(garbanzo.Envar.AppPort), 10))
	}()

	go func() {
		proxyrouter := garbanzo.NewProxyRouter(c)
		proxyrouter.Start(":" + strconv.FormatUint(uint64(garbanzo.Envar.ProxyPort), 10))
	}()

	runtime.Goexit()
}
