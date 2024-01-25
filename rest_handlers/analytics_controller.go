package rest_handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gokul656/go-sockets/internal"
)

func GetHubOverview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	marshalled, _ := json.Marshal(internal.GetHubSummary())
	w.Write([]byte(marshalled))
}

func GetAvailableTopics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	marshalled, _ := json.Marshal(internal.GetAvailableTopics())
	w.Write(marshalled)
}
