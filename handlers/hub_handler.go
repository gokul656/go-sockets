package handlers

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Mu           sync.Mutex
	Connections  map[string]*websocket.Conn
	UpgradedSubs map[string]map[string]struct{} // using struct, it consumes 0 bytes
}

var (
	ConnectionHub *Hub = setupHub()
)

func setupHub() *Hub {
	hub := &Hub{
		Connections:  make(map[string]*websocket.Conn),
		UpgradedSubs: map[string]map[string]struct{}{},
	}

	return hub
}

func (h *Hub) AddConnection(name string, conn *websocket.Conn) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	h.Connections[name] = conn
}

func (h *Hub) RemoveConnection(name string) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	// removing connections from actual connections list
	_ = h.close(name)
	delete(h.Connections, name)

	// removing connection reference from subscriptions
	for k := range h.UpgradedSubs {
		delete(h.UpgradedSubs[k], name)
	}
}

func (h *Hub) GetConnection(name string) *websocket.Conn {
	return h.Connections[name]
}

func (h *Hub) Subscribe(conn string, topic string) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	// creating topic if not exists
	if _, ok := h.UpgradedSubs[topic]; !ok {
		h.UpgradedSubs[topic] = map[string]struct{}{}
	}

	// appending to topics[] if already exists
	h.UpgradedSubs[topic][conn] = struct{}{}
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
	return h.Connections[conn].Close()
}

func (h *Hub) Broadcast(topic string, event []byte) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	for peer := range h.UpgradedSubs[topic] {
		conn := h.GetConnection(peer)
		if conn == nil {
			h.RemoveConnection(peer)
		} else {
			err := conn.WriteMessage(websocket.TextMessage, event)
			if err != nil {
				log.Println("[error]", err)
			}
		}
	}
}