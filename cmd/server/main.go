package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerGetVersion(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := ResVersion{
		Version: []int{0, 0, 1},
	}
	json.NewEncoder(w).Encode(res)
}

func handlerGetAgents(w http.ResponseWriter, req *http.Request)             {}
func handlerPostAgentsAuthorize(w http.ResponseWriter, req *http.Request)   {}
func handlerPostAgentsUnauthorize(w http.ResponseWriter, req *http.Request) {}
func handlerPostAgentsTasks(w http.ResponseWriter, req *http.Request)       {}
func handlerGetAgentsTasks(w http.ResponseWriter, req *http.Request)        {}

func handlerGetAgentsTasksById(w http.ResponseWriter, req *http.Request) {
	taskId := req.PathValue("taskId")
	log.Println(taskId)
	w.Header().Set("Content-Type", "application/json")

	if taskId == "" || !agentTaskExist(taskId) {
		w.WriteHeader(http.StatusOK)
		errRes := ResErr{
			Error: "taskId is not valid",
		}
		json.NewEncoder(w).Encode(errRes)
		return
	}
	agentTaskGet(taskId)
	w.WriteHeader(http.StatusNoContent)
}

func handlerDeleteAgentsTasksById(w http.ResponseWriter, req *http.Request) {
	taskId := req.PathValue("taskId")
	log.Println(taskId)
	w.Header().Set("Content-Type", "application/json")

	if taskId == "" || !agentTaskExist(taskId) {
		w.WriteHeader(http.StatusNotFound)
		errRes := ResErr{
			Error: "taskId is not valid",
		}
		json.NewEncoder(w).Encode(errRes)
		return
	}
	agentTaskDelete(taskId)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
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

	log.Println("Server is started")

	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatalf("Server is failed due to: %v", err)
	}
}
