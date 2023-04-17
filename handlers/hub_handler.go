package handlers

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu           sync.Mutex
	Connections  map[string]*websocket.Conn
	UpgradedSubs map[string]map[string]struct{}
}

var (
	hub *Hub = setupHub()
)

func GetHub() *Hub {
	return hub
}

func setupHub() *Hub {
	hub := &Hub{
		Connections:  make(map[string]*websocket.Conn),
		UpgradedSubs: map[string]map[string]struct{}{},
	}

	return hub
}

func (h *Hub) Subscribe(conn string, topic string) {
	h.mu.Lock()
	if _, ok := h.UpgradedSubs[topic]; !ok {
		h.UpgradedSubs[topic] = map[string]struct{}{}
	}

	h.UpgradedSubs[topic][conn] = struct{}{}

	defer h.mu.Lock()
}

func (h *Hub) GetSubscriptions(conn string) []string {
	keys := make([]string, 0, len(h.UpgradedSubs))
	for k := range h.UpgradedSubs {
		keys = append(keys, k)
	}

	return keys
}

func (h *Hub) UnSubscribe(conn string, topic string) {
	h.mu.Lock()
	delete(h.UpgradedSubs[topic], conn)
	defer h.mu.Lock()
}

func (h *Hub) Close(conn string) error {
	return h.Connections[conn].Close()
}
