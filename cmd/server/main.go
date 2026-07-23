package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var redhowlTmpDir string

func main() {
	redhowlTmpDir = filepath.Join(os.TempDir(), "redhowl")

	err := os.MkdirAll(redhowlTmpDir, os.ModePerm)
	if err != nil {
		log.Fatalln("Unable to create tmpDir")
		return
	}

	go func() {
		agentMux := http.NewServeMux()

		agentMux.HandleFunc("GET /api/version", handlerGetVersion)
		agentMux.HandleFunc("GET /api/agent-com/ws", handlerAgentComWS)
		agentMux.HandleFunc("POST /api/agent-com/register", handlerAgentComRegister)

		log.Println("Agent backend will be started on http://0.0.0.0:4000")
		err := http.ListenAndServe("0.0.0.0:4000", agentMux)

		if err != nil {
			log.Fatalf("Agent backend is failed due to: %v", err)
		}
	}()

	mux := http.NewServeMux()

	// version
	mux.HandleFunc("GET /api/version", handlerGetVersion)

	// agents
	mux.HandleFunc("GET /api/agents", handlerGetAgents)
	mux.HandleFunc("POST /api/agents/authorize", handlerPostAgentsAuthorize)
	mux.HandleFunc("POST /api/agents/unauthorize", handlerPostAgentsUnauthorize)

	// tasks
	mux.HandleFunc("POST /api/agents/tasks", handlerPostAgentsTasks)
	mux.HandleFunc("GET /api/agents/tasks", handlerGetAgentsTasks)
	mux.HandleFunc("GET /api/agents/tasks/{taskId}", handlerGetAgentsTasksById)
	mux.HandleFunc("DELETE /api/agents/tasks/{taskId}", handlerDeleteAgentsTasksById)

	fs := http.FileServer(http.Dir("www"))
	mux.Handle("GET /", fs)

	log.Println("Admin dashboard will be started on http://127.0.0.1:3000")

	err = http.ListenAndServe("127.0.0.1:3000", mux)
	if err != nil {
		log.Fatalf("Admin dashboard is failed due to: %v", err)
	}
}
