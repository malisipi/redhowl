package main

import (
	"encoding/json"
	"log"
	"net/http"
	"redhowl/cmd/internal"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func handlerAgentComWS(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("WS Connection failed due to, %v\n", err)
		return
	}
	defer conn.Close()

	log.Printf("New agent connected: %s\n", conn.RemoteAddr())

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Agent is disconnected from host: %s, %v\n", conn.RemoteAddr(), err)
		}

		var header internal.WSTypeHeader
		if err := json.Unmarshal(msg, &header); err != nil {
			log.Println("Borked package was recieved, skipping")
			continue
		}

		switch header.Type {
		case internal.WSTypeMetricSend:
			{
				var metrics internal.WSMetricSend
				if err := json.Unmarshal(msg, &metrics); err == nil {
					log.Printf("Got metrics data\nCPU: %v\nMemory| Used:%v Total:%v|\n"+
						"Disk| Used:%v Total:%v Mount:%v\nNetwork| IPv4:%v, IPv6:%v, MAC:%v\n",
						metrics.CPU, metrics.Memory.Used, metrics.Memory.Total,
						metrics.Disk.Used, metrics.Disk.Total, metrics.Disk.MountPoint, metrics.Network.IPv4, metrics.Network.IPv6, metrics.Network.MAC)
				} else {
					log.Println("Can't parsed JSON, Type and Struct is mismatched")
				}
			}
		default:
			{
				log.Println("Unknown WSType data. Is client and server running same version or not?")
			}
		}
	}
}
