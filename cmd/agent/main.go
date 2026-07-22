package main

import (
	"log"
	"net/url"
	"redhowl/cmd/internal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	time.Sleep(2 * time.Second) // Wait to start server, test-only

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:4000", Path: "/api/agent-com/ws"}
	log.Printf("Agent trying to connect to the server: %s\n", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("Can't connect server: %v\n", err)
		return
	}
	defer conn.Close()

	stats := internal.WSMetricSend{
		WSTypeHeader: internal.WSTypeHeader{Type: internal.WSTypeMetricSend},
		CPU:          74.1,
		Memory:       internal.MetricsMemory{Used: 1.2, Total: 4.7},
		Disk:         internal.MetricsDisk{Used: 24, Total: 68, MountPoint: "/"},
		Network:      internal.MetricsNetwork{IPv4: "64.12.15.47", IPv6: "::0", LocalIPv4: "127.0.0.1", LocalIPv6: "::0", MAC: "00:11:22:33:44:55"},
	}

	err = conn.WriteJSON(stats)
	if err != nil {
		log.Printf("Failed to send data due to: %v\n", err)
		return
	}

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Connection lost due to: %v\n", err)
			return // TODO: try reconnect later
		}

		log.Printf("Got data from server: %v\n", msg)
	}
}
