package internal

import (
	"github.com/gokul656/go-sockets/types"
)

func GetHubSummary() types.HubDetails {
	length := len(ConnectionHub.Connections)
	conns := make([]types.ConnectionDetails, 0)

	for _, conn := range ConnectionHub.Connections {
		sub := make([]string, 0)
		sub = append(sub, ConnectionHub.GetSubscriptions(conn.Conn.RemoteAddr().String())...)

		conns = append(conns, types.ConnectionDetails{
			ConnectionId:    conn.ConnectionId,
			Status:          conn.Status,
			SubscribedTopic: sub,
			ConnectedAt:     conn.ConnectedAt,
			DisconnectedAt:  conn.DisconnectedAt,
		})
	}

	return types.HubDetails{
		ActiveConnections: uint64(length),
		Connections:       conns,
	}
}

func GetAvailableTopics() []types.TopicData {
	topics := make([]types.TopicData, 0)
	for key := range ConnectionHub.AvailableTopics {
		topics = append(topics, types.TopicData{
			Topic:  key,
			Status: ConnectionHub.AvailableTopics[key],
		})
	}
	return topics
}
