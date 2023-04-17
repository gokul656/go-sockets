package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gokul656/go-sockets/types"
	"github.com/gorilla/websocket"
)

var (
	upgrader      = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[new] connection request incoming from: ", r.RemoteAddr)
	newConn, _ := upgrader.Upgrade(w, r, nil)

	hub.mu.Lock()
	hub.Connections[r.RemoteAddr] = newConn
	hub.mu.Unlock()

	reader(newConn)
}

func reader(conn *websocket.Conn) {
	for {
		_, p, _ := conn.ReadMessage()
		message := &types.Message{}
		json.Unmarshal([]byte(p), message)
		
		switch message.Method {
		case types.SUB:
			hub.Subscribe(conn.RemoteAddr().String(), message.Ch)
		case types.UNSUB:
			hub.UnSubscribe(conn.RemoteAddr().String(), message.Ch)
		default:
			message = &types.Message{
				Ch: "error",
				Payload: "Invalid request",
			}
			marshalled, _ := json.Marshal(message)
			conn.WriteMessage(websocket.TextMessage, marshalled);
			continue
		}
	}
}