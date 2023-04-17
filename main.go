package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"net/http"
	"time"

	"github.com/gokul656/go-sockets/handlers"
	"github.com/gokul656/go-sockets/types"
	"github.com/gorilla/websocket"
)

const (
	ONE_SEC = time.Second * 2
)

var (
	addr = flag.String("addr", "localhost:8080", "ws service address")

	MARKET_CH = make(chan types.Message, 1024)
	TICKER_CH = make(chan types.Message, 1024)
)

func main() {
	go setupChannels(MARKET_CH, TICKER_CH)
	go setupPulbishers("market", MARKET_CH)
	go setupPulbishers("ticker", TICKER_CH)

	http.HandleFunc("/ws", handlers.RootHandler)
	http.ListenAndServe(*addr, nil)
}

func setupChannels(market chan types.Message, ticker chan types.Message) {
	go produceMockData("market", market)
	go produceMockData("ticker", ticker)
}

func setupPulbishers(topic string, channel chan types.Message) {
	hub := handlers.GetHub()
	for {
		for k := range hub.UpgradedSubs[topic] {
			conn := hub.Connections[k]
			message := <-channel
			marshalled, _ := json.Marshal(message)
			conn.WriteMessage(websocket.TextMessage, marshalled)
		
		}
	}
}

func produceMockData(name string, channel chan<- types.Message) {
	for {
		time.Sleep(ONE_SEC)
		message := types.Message{
			Ch:      name,
			Method:  types.SUB,
			Payload: fmt.Sprintf("Ping from %s", name),
		}

		channel <- message
	}
}
