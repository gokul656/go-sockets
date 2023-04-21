package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/url"
	"time"

	"github.com/gokul656/go-sockets/types"
	"github.com/gorilla/websocket"
)

var (
	address = flag.String("address", "localhost:8080", "")
	serverURL     = url.URL{Scheme: "ws", Host: *address, Path: "/ws"}

	pongWait = 5 * time.Second
)

type Pong struct {
	Pong int64 `json:"pong"`
}

func main() {
	conn, err := connect()
	handleSocketErr(conn, err)
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPingHandler(func(appData string) error {
		pong:= &Pong{Pong: time.Now().UnixMilli()}
		marshaled, _ := json.Marshal(pong)
		conn.WriteMessage(websocket.PongMessage, []byte(marshaled))

		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	subscribe(conn, "ticker", "market", "unknown")
	
	for {
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			log.Println("[socket] closed", conn.RemoteAddr())
			break
		}

		handleMessage(conn, messageType, payload)
	}
}

func connect() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(serverURL.String(), nil)
	return conn, err
}

func subscribe(conn *websocket.Conn, topics ...string) {
	for _, topic := range topics {
		body := &types.Message{
			Ch: topic,
			Method: "sub",
			Ts: time.Now().UnixMilli(),
		}

		marshalled, _ := json.Marshal(body)
		err := conn.WriteMessage(websocket.TextMessage, marshalled)
		if err == nil {
			log.Println("[socket] subscribed to", topic)
		}

	}
}

func handleMessage(conn *websocket.Conn, messageType int, message []byte) {
	log.Println("[message]", string(message))
}

func handleSocketErr(conn *websocket.Conn, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered & Retrying")
		}
	}()

	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
			log.Println("[socket]", "Connection abnormally closed by client")
			conn.Close()
		} else {
			panic(err)
		}
	}
}
