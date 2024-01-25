package types

import (
	"time"

	"github.com/gokul656/go-sockets/pkg"
)

type Method string

type Topic string

const (
	SUB   Method = "sub"
	UNSUB Method = "unsub"
)

type Ping struct {
	Ping    int64 `json:"ping"`
	Message string
}

type Message struct {
	Ch      string `json:"ch,omitempty"`
	Method  Method `json:"method,omitempty"`
	Payload any    `json:"payload,omitempty"`
	Ts      int64  `json:"ts"`
}

type SocketConnection interface {
	Subscribe(conn string, topic string) error
	GetSubscriptions(conn string) []string
	UnSubscribe(conn string, topic string)
	Close(conn string) error
}

type ConnectionDetails struct {
	ConnectionId    string               `json:"connection_id"`
	Status          pkg.ConnectionStatus `json:"status"`
	SubscribedTopic []string             `json:"subscribed_topics"`
	ConnectedAt     time.Time            `json:"connected_at"`
	DisconnectedAt  time.Time            `json:"disconnected_at"`
}

type HubDetails struct {
	ActiveConnections uint64              `json:"active_connections"`
	Connections       []ConnectionDetails `json:"connections"`
}

type TopicData struct {
	Topic  string `json:"topic"`
	Status bool   `json:"status"`
}
