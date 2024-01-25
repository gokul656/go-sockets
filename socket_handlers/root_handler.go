package socket_handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gokul656/go-sockets/internal"
	"github.com/gokul656/go-sockets/pkg"
	"github.com/gokul656/go-sockets/types"
	"github.com/google/uuid"
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
	cid, _ := uuid.NewRandom()
	soc := &pkg.Connection{
		Conn:         conn,
		ConnectionId: cid.String(),
	}

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(data string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	internal.ConnectionHub.AddConnection(r.RemoteAddr, soc)
	go socketReader(soc)
}

func PingHandler() {
	pingTicker := time.NewTicker(time.Second * 3)
	defer pingTicker.Stop()

	for {
		<-pingTicker.C

		ping := &types.Ping{Ping: time.Now().UnixMilli(), Message: "Hi"}
		marshal, _ := json.Marshal(&ping)
		internal.ConnectionHub.Broadcast("ping", marshal)
	}
}

func socketReader(soc *pkg.Connection) {
	conn := soc.Conn
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("[socket] closed", conn.RemoteAddr())
			soc.Status = pkg.CLOSED
			soc.DisconnectedAt = time.Now()
			break
		}

		if messageType == websocket.PingMessage {
			conn.SetReadDeadline(time.Now().Add(pongWait))
			return
		}

		message := &types.Message{}
		json.Unmarshal([]byte(p), message)
		internal.ConnectionHub.Subscribe(conn.RemoteAddr().String(), message.Ch)

		switch message.Method {
		case types.SUB:
			err = internal.ConnectionHub.Subscribe(conn.RemoteAddr().String(), message.Ch)
			if err != nil {
				message = &types.Message{
					Ch:      message.Ch,
					Payload: err.Error(),
					Ts:      time.Now().UnixMilli(),
				}

				marshalled, _ := json.Marshal(message)
				soc.SendMessage(websocket.TextMessage, marshalled)
			}
		case types.UNSUB:
			internal.ConnectionHub.UnSubscribe(conn.RemoteAddr().String(), message.Ch)
		default:
			message = &types.Message{
				Ch:      "error",
				Payload: "invalid request. Input should in the format of { ch: string, metod: SUB | UNSUB, payload: string}",
				Ts:      time.Hour.Milliseconds(),
			}
			marshalled, _ := json.Marshal(message)
			soc.SendMessage(websocket.TextMessage, marshalled)
			continue
		}
	}
}
