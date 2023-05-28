package main

import garbanzo "github.com/kijimaD/garbanzo/pkg"

func main() {
	router := garbanzo.NewRouter("pkg/templates/*.html")
	router.Start(":8080")
}
