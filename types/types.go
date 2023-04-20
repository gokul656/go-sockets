package types

type Method string

type Topic string

const (
	SUB   Method = "sub"
	UNSUB Method = "unsub"
)

type Message struct {
	Ch      string `json:"ch,omitempty"`
	Method  Method `json:"method,omitempty"`
	Payload any `json:"payload,omitempty"`
}

type SocketConnection interface {
	Subscribe(conn string, topic string)
	GetSubscriptions(conn string) []string
	UnSubscribe(conn string, topic string)
	Close(conn string) error
}
