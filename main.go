package main

import garbanzo "github.com/kijimaD/garbanzo/pkg"

func main() {
	router := garbanzo.NewRouter()
	router.Start(":8080")
}
