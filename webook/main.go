package main

import (
	"github.com/kisara71/WeBook/webook/ioc"
)

func main() {
	server := ioc.InitWebServer()
	if err := server.Run(":8080"); err != nil {
		panic(err)
		return
	}
}
