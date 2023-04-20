package mockers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gokul656/go-sockets/handlers"
	"github.com/gokul656/go-sockets/types"
)

const (
	ONE_SEC = time.Second * 1
)

var (
	marketch = make(chan *types.Message, 1024)
	tickerch = make(chan *types.Message, 1024)
)

func StartMockers() {
	go ProduceMockData("market", marketch)
	go ProduceMockData("ticket", tickerch)

	go PublishMockData("market", marketch)
	go PublishMockData("ticker", tickerch)
}

func ProduceMockData(name string, channel chan<- *types.Message) {
	for {
		time.Sleep(ONE_SEC)
		message := &types.Message{
			Ch:      name,
			Payload: fmt.Sprintf("Ping from %s", name),
			Ts:      time.Now().UnixMilli(),
		}

		channel <- message
	}
}

func PublishMockData(topic string, channel <-chan *types.Message) {
	for {
		message := <-channel
		marshalled, _ := json.Marshal(message)
		handlers.ConnectionHub.Broadcast(topic, marshalled)
	}
}
