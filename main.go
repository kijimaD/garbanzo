package main

import (
	"runtime"

	garbanzo "github.com/kijimaD/garbanzo/pkg"
)

func main() {
	go func() {
		router := garbanzo.NewRouter("pkg/templates/*.html")
		router.Start(":8080")
	}()

	go func() {
		proxyrouter := garbanzo.NewProxyRouter()
		proxyrouter.Start(":8081")
	}()

	runtime.Goexit()
}
