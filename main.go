package main

import (
	"github.com/sergivb01/paste/server"
)

func main() {
	srv, err := server.New()
	if err != nil {
		panic(err)
	}
	srv.Listen()
}
