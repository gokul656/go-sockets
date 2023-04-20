package handlers

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Mu           sync.Mutex
	Connections  map[string]*websocket.Conn
	UpgradedSubs map[string]map[string]struct{}
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

func (h *Hub) GetConnection(id string) *websocket.Conn {
	return h.Connections[id]
}

func (h *Hub) Subscribe(conn string, topic string) {
	h.Mu.Lock()
	if _, ok := h.UpgradedSubs[topic]; !ok {
		h.UpgradedSubs[topic] = map[string]struct{}{}
	}

	h.UpgradedSubs[topic][conn] = struct{}{}

	defer h.Mu.Unlock()
}

func (h *Hub) GetSubscriptions(conn string) []string {
	keys := make([]string, 0, len(h.UpgradedSubs))
	for k := range h.UpgradedSubs {
		keys = append(keys, k)
	}

	return keys
}

func (h *Hub) UnSubscribe(conn string, topic string) {
	h.Mu.Lock()
	delete(h.UpgradedSubs[topic], conn)
	defer h.Mu.Unlock()
}

func (h *Hub) Close(conn string) error {
	return h.Connections[conn].Close()
}

func (h *Hub) Broadcast(topic string, event []byte) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	for peer := range h.UpgradedSubs[topic] {
		conn := h.GetConnection(peer)
		if conn == nil {
			delete(h.UpgradedSubs, peer)
		} else {
			err := conn.WriteMessage(websocket.TextMessage, event)
			if err != nil {
				log.Println("[error]", err)
			}
		}
	}
}
