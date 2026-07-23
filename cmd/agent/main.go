package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"redhowl/cmd/internal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	conn      *websocket.Conn
	connAlive bool
	connMutex sync.Mutex
)

func main() {
	agentUUID := os.Getenv("REDHOWL-AGENT-UUID")
	if agentUUID == "" {
		agentUUID = "wolf-agent-1" // TODO: Replace that static string with real id gen and load
	}
	agentHomeServer := os.Getenv("REDHOWL-AGENT-HOME")
	if agentHomeServer == "" {
		agentHomeServer = "127.0.0.1:4000" // it defaults to localhost if not set
	}

	registerUri := url.URL{Scheme: "http", Host: agentHomeServer, Path: "/api/agent-com/register"}
	registerBody := internal.ReqAgentRegister{
		UUID:    agentUUID,
		User:    getUserInfo(),
		OS:      getOSInfo(),
		Machine: getMachineInfo(),
	}

	go func() {
		for {
			connMutex.Lock()
			if connAlive {
				stats := internal.WSMetricSend{
					WSTypeHeader: internal.WSTypeHeader{Type: internal.WSTypeMetricSend},
					CPU:          getCpuUsage(),
					Memory:       getMemoryStats(),
					Disk:         getDiskStats(),
					Network:      getNetworkStats(),
				}

				err := conn.WriteJSON(stats)
				if err != nil {
					log.Printf("Failed to send data due to: %v\n", err)
				}
			}
			connMutex.Unlock()
			time.Sleep(1 * time.Second)
		}
	}()

registerItselfBefore:
	for { // try again till send a register event
		var registerBodyBuffer bytes.Buffer
		err := json.NewEncoder(&registerBodyBuffer).Encode(registerBody)
		if err != nil {
			log.Fatalln("Failed to create JSON for registration")
		}

		req, err := http.NewRequest(http.MethodPost, registerUri.String(), &registerBodyBuffer)
		if err != nil {
			log.Fatalln("Failed to create new request")
		}

		_, err = (&http.Client{}).Do(req)
		if err != nil {
			log.Println("Failed to send register")
		} else {
			break
		}
		time.Sleep(2 * time.Second)
	}

	wsUri := url.URL{Scheme: "ws", Host: agentHomeServer, Path: "/api/agent-com/ws", RawQuery: "uuid=" + agentUUID}
	log.Printf("Agent trying to connect to the server: %s\n", wsUri.String())

	var err error
	conn, _, err = websocket.DefaultDialer.Dial(wsUri.String(), nil)
	if err != nil {
		log.Printf("Can't connect server: %v\n", err)
	} else {
		connMutex.Lock()
		connAlive = true
		connMutex.Unlock()
		for {
			var msg map[string]interface{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				connMutex.Lock()
				connAlive = false
				log.Printf("Connection lost due to: %v\n", err)
				connMutex.Unlock()
				goto registerItselfBefore
			}

			log.Printf("Got data from server: %v\n", msg)
		}
	}
	time.Sleep(1 * time.Second)
	goto registerItselfBefore
}
