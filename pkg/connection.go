package pkg

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ConnMu         sync.Mutex
	ConnectionId   string
	Status         ConnectionStatus
	Conn           *websocket.Conn
	ConnectedAt    time.Time
	DisconnectedAt time.Time
}

func (s *Connection) SendMessage(messageType int, message []byte) error {
	s.ConnMu.Lock()
	defer s.ConnMu.Unlock()

	return s.Conn.WriteMessage(messageType, message)
}
