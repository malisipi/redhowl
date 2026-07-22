package main

import (
	"redhowl/cmd/internal"
	"time"
)

type ReqAuthorize struct {
	UUID string `json:"uuid"`
}

type ResErr struct {
	Error string `json:"error"`
}

type ResVersion struct {
	Version []int `json:"version"`
}

type ResAgents struct {
	Agents []Agent `json:"agents"`
}

type Agent struct {
	UUID               string       `json:"uuid"`
	Status             string       `json:"status"`
	ConnectedTimestamp time.Time    `json:"connectedTimestamp"`
	Metrics            AgentMetrics `json:"metrics"`
}

type AgentMetrics struct {
	CPU     float64                 `json:"cpu"`
	Memory  internal.MetricsMemory  `json:"memory"`
	Disk    internal.MetricsDisk    `json:"disk"`
	User    internal.MetricsUser    `json:"user"`
	OS      internal.MetricsOS      `json:"os"`
	Machine internal.MetricsMachine `json:"machine"`
	Network internal.MetricsNetwork `json:"network"`
}

type TaskGeneric struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Timestamp   time.Time  `json:"timestamp"`
	Status      TaskStatus `json:"status"`
}

type TaskStatus struct {
	Running  int `json:"running"`
	Finished int `json:"finished"`
	Failed   int `json:"failed"`
}

type TaskDetailed struct {
	TaskGeneric
	Permission TaskPermission             `json:"permission"`
	Reports    map[string]TaskAgentReport `json:"reports"`
}

type TaskAgentReport struct {
	Logs      []string       `json:"logs"`
	WasmLogs  []string       `json:"wasmLogs"`
	Artifacts []TaskArtifact `json:"artifacts"`
	Status    string         `json:"status"`
}

type TaskArtifact struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

type ReqTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions TaskPermission
	Agents      []string `json:"agents"`
	// Wasm part must be parsed as raw binary
}

type TaskPermission struct {
	Filesystem bool     `json:"filesystem"`
	Network    bool     `json:"network"`
	Execution  []string `json:"execution"`
	Artifacts  bool     `json:"artifacts"`
}

type ResTaskId struct {
	Task struct {
		ID string `json:"id"`
	} `json:"task"`
}

type ResTasks struct {
	Tasks []TaskGeneric `json:"tasks"`
}
