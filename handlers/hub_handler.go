package handlers

import (
	"errors"
	"log"
	"sync"

	"github.com/gokul656/go-sockets/types"
	"github.com/gorilla/websocket"
)

type Hub struct {
	Mu           sync.Mutex
	Connections  map[string]*types.Connection
	UpgradedSubs map[string]map[string]struct{} // using struct, it consumes 0 bytes
}

var (
	ConnectionHub   *Hub            = setupHub()
	availableTopics map[string]bool = make(map[string]bool)
)

func setupHub() *Hub {
	hub := &Hub{
		Connections:  make(map[string]*types.Connection),
		UpgradedSubs: map[string]map[string]struct{}{},
	}

	availableTopics["ticker"] = true
	availableTopics["market"] = true

	return hub
}

func (h *Hub) AddConnection(name string, soc *types.Connection) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	h.Connections[name] = soc
}

func (h *Hub) RemoveConnection(name string) {
	// TO-DO: Fix lock issue
	// h.Mu.TryLock()
	// defer h.Mu.Unlock()

	// removing connections from actual connections list
	h.close(name)
	delete(h.Connections, name)

	// removing connection reference from subscriptions
	for k := range h.UpgradedSubs {
		delete(h.UpgradedSubs[k], name)
	}
}

func (h *Hub) GetConnection(name string) *types.Connection {
	return h.Connections[name]
}

func (h *Hub) Subscribe(conn string, topic string) error {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	if _, ok := availableTopics[topic]; !ok {
		return errors.New("topic is unavailable or invalid")
	}

	// creating topic if not exists
	if _, ok := h.UpgradedSubs[topic]; !ok {
		h.UpgradedSubs[topic] = map[string]struct{}{}
	}

	// appending to topics[] if already exists
	h.UpgradedSubs[topic][conn] = struct{}{}

	return nil
}

func (h *Hub) GetSubscriptions(conn string) []string {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	keys := make([]string, 0, len(h.UpgradedSubs))
	for k := range h.UpgradedSubs {
		keys = append(keys, k)
	}

	return keys
}

func (h *Hub) UnSubscribe(conn string, topic string) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	delete(h.UpgradedSubs[topic], conn)
}

func (h *Hub) close(conn string) error {
	return h.Connections[conn].Conn.Close()
}

func (h *Hub) Broadcast(topic string, event []byte) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	for peer := range h.UpgradedSubs[topic] {
		soc := h.GetConnection(peer)
		conn := soc.Conn
		if conn == nil {
			h.RemoveConnection(peer)
		} else {
			soc.ConnMu.Lock()
			err := conn.WriteMessage(websocket.TextMessage, event)
			if err != nil {
				log.Println("[error] socket abruptly closed")
				h.RemoveConnection(peer)
			}

			soc.ConnMu.Unlock()
		}
	}
}
