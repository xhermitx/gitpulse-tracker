package main

import (
	"github.com/xhermitx/gitpulse-tracker/github-service/servers"
)

// SERVICE TO FETCH AND HANDLE GITHUB DATA
func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Panic("Error loading the environment variables")
	// }
	servers.HttpServer()
}
