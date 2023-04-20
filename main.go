package main

import (
	"flag"
	"log"

	"net/http"

	"github.com/gokul656/go-sockets/handlers"
	"github.com/gokul656/go-sockets/mockers"
)

var (
	addr = flag.String("addr", ":8080", "ws service address")
)

func main() {
	flag.Parse()
	log.Println("[server] up and running on", *addr)

	mockers.StartMockers()

	http.HandleFunc("/ws", handlers.RootHandler)
	http.ListenAndServe(*addr, nil)
}
