package internal

import (
	"github.com/gokul656/go-sockets/pkg"
)

var (
	ConnectionHub   *pkg.Hub        = nil
	availableTopics map[string]bool = make(map[string]bool)
)

func init() {
	ConnectionHub = setupHub()
}

func setupHub() *pkg.Hub {
	availableTopics["ticker"] = true
	availableTopics["market"] = true
	availableTopics["kline"] = false

	hub := &pkg.Hub{
		Connections:     make(map[string]*pkg.Connection),
		UpgradedSubs:    map[string]map[string]struct{}{},
		AvailableTopics: availableTopics,
	}

	return hub
}
