package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gokul656/go-sockets/mockers"
	"github.com/gokul656/go-sockets/rest_handlers"
	"github.com/gokul656/go-sockets/socket_handlers"
)

var addr = flag.String("addr", ":8080", "ws service address")

func main() {
	flag.Parse()
	log.Println("[server] up and running on", *addr)

	go mockers.StartMockers()
	go socket_handlers.PingHandler()

	http.HandleFunc("/api/details", rest_handlers.GetHubOverview)
	http.HandleFunc("/api/topics", rest_handlers.GetAvailableTopics)
	http.HandleFunc("/ws", socket_handlers.RootHandler)
	http.ListenAndServe(*addr, nil)
}
