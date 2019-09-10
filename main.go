package main

import (
	"github.com/sergivb01/paste/server"
)

func main() {
	srv, err := server.New("./config.yml")
	if err != nil {
		panic(err)
	}
	srv.Listen()
}
