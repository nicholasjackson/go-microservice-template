package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nicholasjackson/microservice-template/global"
	"github.com/nicholasjackson/microservice-template/handlers"
)

func main() {
	config := os.Args[1]
	rootfolder := os.Args[2]

	fmt.Println("Loading Config:" + config)
	global.LoadConfig(config, rootfolder)

	setupHandlers()
}

func setupHandlers() {
	http.Handle("/", handlers.GetRouter())

	fmt.Println("Listening for connections on port", 8001)
	http.ListenAndServe(fmt.Sprintf(":%v", 8001), nil)
}
