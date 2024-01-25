package mockers

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gokul656/go-sockets/internal"
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
	wg := sync.WaitGroup{}
	go ProduceMockData("market", marketch, &wg)
	go ProduceMockData("ticker", tickerch, &wg)

	go PublishMockData("market", marketch, &wg)
	go PublishMockData("ticker", tickerch, &wg)

	wg.Wait()
}

func ProduceMockData(name string, channel chan<- *types.Message, wg *sync.WaitGroup) {
	for {
		message := &types.Message{
			Ch:      name,
			Payload: fmt.Sprintf("Ping from %s", name),
			Ts:      time.Now().UnixMilli(),
		}

		time.Sleep(1 * time.Second)
		channel <- message
	}
}

func PublishMockData(topic string, channel <-chan *types.Message, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for {
		message, ok := <-channel
		if !ok {
			break
		}
		marshalled, _ := json.Marshal(message)
		internal.ConnectionHub.Broadcast(topic, marshalled)
	}
}
