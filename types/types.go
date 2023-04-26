package types

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
