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

func handlerAgentComRegister(w http.ResponseWriter, req *http.Request) {
	var reqBody internal.ReqAgentRegister
	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode((ResErr{Error: "Invalid JSON request"}))
		return
	}

	log.Println("New agent registered")

	agentRegister(reqBody)
	w.WriteHeader(http.StatusNoContent)
}

func handlerAgentComWS(w http.ResponseWriter, req *http.Request) {
	agentUUID := req.FormValue("uuid")

	if !agentAuthorized(agentUUID) {
		log.Println("New agent request was revoked since it's unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Printf("WS Connection failed due to, %v\n", err)
		return
	}
	defer conn.Close()

	log.Printf("New agent connected: %s\n", conn.RemoteAddr())

	agentsListLock.Lock()
	for i := range agentsList {
		if agentsList[i].UUID == agentUUID {
			agentsList[i].WSConn = conn
			break
		}
	}
	agentsListLock.Unlock()

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
					agentsListLock.Lock()
					for i := range agentsList {
						if agentsList[i].UUID == agentUUID {
							agentsList[i].Metrics.CPU = metrics.CPU
							agentsList[i].Metrics.Memory = metrics.Memory
							agentsList[i].Metrics.Disk = metrics.Disk
							agentsList[i].Metrics.Network = metrics.Network
							break
						}
					}
					agentsListLock.Unlock()
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
