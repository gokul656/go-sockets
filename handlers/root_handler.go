package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gokul656/go-sockets/types"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[new] connection request incoming from: ", r.RemoteAddr)
	conn, _ := upgrader.Upgrade(w, r, nil)
	ConnectionHub.AddConnection(r.RemoteAddr, conn)
	socketReader(conn)
}

func socketReader(conn *websocket.Conn) {
	for {
		_, p, _ := conn.ReadMessage()
		message := &types.Message{}
		json.Unmarshal([]byte(p), message)

		switch message.Method {
		case types.SUB:
			ConnectionHub.Subscribe(conn.RemoteAddr().String(), message.Ch)
		case types.UNSUB:
			ConnectionHub.UnSubscribe(conn.RemoteAddr().String(), message.Ch)
		default:
			ConnectionHub.Mu.Lock()
			keys := make([]string, 0, len(ConnectionHub.UpgradedSubs))
			for k := range ConnectionHub.UpgradedSubs {
				keys = append(keys, k)
			}

			message = &types.Message{
				Ch:      "error",
				Payload: keys,
				Ts:      time.Hour.Milliseconds(),
			}
			marshalled, _ := json.Marshal(message)
			conn.WriteMessage(websocket.TextMessage, marshalled)
			ConnectionHub.Mu.Unlock()
			continue
		}
	}
}
