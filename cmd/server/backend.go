package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func handlerGetVersion(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := ResVersion{
		Version: []int{0, 0, 1},
	}
	json.NewEncoder(w).Encode(res)
}

func handlerGetAgents(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res := ResAgents{
		Agents: agentGetList(),
	}
	json.NewEncoder(w).Encode(res)
}

func handlerPostAgentsAuthorize(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var reqBody ReqAuthorize
	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode((ResErr{Error: "Invalid JSON request"}))
		return
	}

	if !agentExist(reqBody.UUID) && reqBody.UUID != "*" {
		w.WriteHeader(http.StatusBadRequest)
		errRes := ResErr{
			Error: "UUID is not valid",
		}
		json.NewEncoder(w).Encode(errRes)
		return
	}
	agentAuthorize(reqBody.UUID)
	w.WriteHeader(http.StatusNoContent)
}

func handlerPostAgentsUnauthorize(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var reqBody ReqAuthorize
	err := json.NewDecoder(req.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode((ResErr{Error: "Invalid JSON request"}))
		return
	}

	if !agentExist(reqBody.UUID) && reqBody.UUID != "*" {
		w.WriteHeader(http.StatusBadRequest)
		errRes := ResErr{
			Error: "UUID is not valid",
		}
		json.NewEncoder(w).Encode(errRes)
		return
	}
	agentUnauthorize(reqBody.UUID)
	w.WriteHeader(http.StatusNoContent)
}

func handlerPostAgentsTasks(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req.ParseMultipartForm(32 << 20) // max allowed 32 MiB

	formJSON := req.FormValue("json")

	var reqBody ReqTask
	err := json.Unmarshal([]byte(formJSON), &reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode((ResErr{Error: "Invalid JSON part of form"}))
		return
	}

	file, _, err := req.FormFile("wasm")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode((ResErr{Error: "Missing/Corrupted WASM file"}))
		return
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp(redhowlTmpDir, "task_*.wasm")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode((ResErr{Error: "Unable to copy wasm file"}))
		return
	}
	defer tmpFile.Close()

	{
		_, err := io.Copy(tmpFile, file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode((ResErr{Error: "Unable to copy wasm file"}))
			return
		}
	}

	fileId := tmpFile.Name()
	log.Printf("WASM file is saved to: %s\n", fileId)

	res := ResTaskId{}
	res.Task.ID = agentTaskRun(reqBody)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(res)
}

func handlerGetAgentsTasks(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader((http.StatusOK))
	res := ResTasks{
		Tasks: agentTaskList(),
	}
	json.NewEncoder(w).Encode(res)
}

func handlerGetAgentsTasksById(w http.ResponseWriter, req *http.Request) {
	taskId := req.PathValue("taskId")
	w.Header().Set("Content-Type", "application/json")

	if taskId == "" || !agentTaskExist(taskId) {
		w.WriteHeader(http.StatusNotFound)
		errRes := ResErr{
			Error: "taskId is not valid",
		}
		json.NewEncoder(w).Encode(errRes)
		return
	}
	agentTaskGet(taskId) // TODO: return the data
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
