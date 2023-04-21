package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gokul656/go-sockets/types"
	"github.com/gorilla/websocket"
)

const pongWait = time.Second * 5

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[new] connection request incoming from: ", r.RemoteAddr)
	conn, _ := upgrader.Upgrade(w, r, nil)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	go handlePing(conn)
	
	conn.SetPongHandler(func(data string) error {
		log.Println("[pong]", data)
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	ConnectionHub.AddConnection(r.RemoteAddr, conn)
	socketReader(conn)
}

func handlePing(conn *websocket.Conn) {
	pingTicker := time.NewTicker(time.Second * 3)
	defer pingTicker.Stop()
	defer conn.Close()

	for {
		<-pingTicker.C
		ping := &types.Ping{Ping: time.Now().UnixMilli(), Message: "Hi"}
		marshal, _ := json.Marshal(&ping)

		err := conn.WriteMessage(websocket.PingMessage, marshal)
		if err != nil {
			break
		}
	}
}

func socketReader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("[socket] closed", conn.RemoteAddr())
			break
		}

		if messageType == websocket.PingMessage {
			conn.SetReadDeadline(time.Now().Add(pongWait))
			return
		}

		message := &types.Message{}
		json.Unmarshal([]byte(p), message)
		ConnectionHub.Subscribe(conn.RemoteAddr().String(), message.Ch)

		switch message.Method {
		case types.SUB:
			err = ConnectionHub.Subscribe(conn.RemoteAddr().String(), message.Ch)
			if err != nil {
				message = &types.Message{
					Ch: message.Ch,
					Payload: err.Error(),
					Ts: time.Now().UnixMilli(),
				}
				socketWriter(conn, *message)
			}
		case types.UNSUB:
			ConnectionHub.UnSubscribe(conn.RemoteAddr().String(), message.Ch)
		default:
			message = &types.Message{
				Ch:      "error",
				Payload: "invalid request",
				Ts:      time.Hour.Milliseconds(),
			}
			socketWriter(conn, *message)
			continue
		}
	}
}

func socketWriter(conn *websocket.Conn, message types.Message) error {
	ConnectionHub.Mu.Lock()
	marshalled, _ := json.Marshal(message)
	ConnectionHub.Mu.Unlock()
	return conn.WriteMessage(websocket.TextMessage, marshalled)
}
