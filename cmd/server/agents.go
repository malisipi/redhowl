package main

import (
	"redhowl/cmd/internal"
	"sync"
	"time"
)

var agentsList []Agent
var agentsListLock sync.RWMutex

func agentRegister(agentReq internal.ReqAgentRegister) {
	agentsListLock.Lock()
	defer agentsListLock.Unlock()

	for i := range agentsList { // return if exist
		if agentsList[i].UUID == agentReq.UUID {
			return
		}
	}

	agentsList = append(agentsList, Agent{
		UUID:               agentReq.UUID,
		Status:             "unauthorized",
		ConnectedTimestamp: time.Now(),
		Metrics: AgentMetrics{
			User:    agentReq.User,
			OS:      agentReq.OS,
			Machine: agentReq.Machine,
		},
	})
}

func agentGetList() []Agent {
	agentsListLock.RLock()
	defer agentsListLock.RUnlock()

	copiedList := make([]Agent, len(agentsList))
	copy(copiedList, agentsList)
	return copiedList
}

func agentExist(agentUUID string) bool {
	agentsListLock.RLock()
	defer agentsListLock.RUnlock()

	for i := range agentsList {
		if agentsList[i].UUID == agentUUID {
			return true
		}
	}
	return false
}

func agentAuthorized(agentUUID string) bool {
	agentsListLock.RLock()
	defer agentsListLock.RUnlock()

	for i := range agentsList {
		if agentsList[i].UUID == agentUUID {
			return agentsList[i].Status == "authorized" // MUST be AUTHORIZED as EXPLICITLY
		}
	}
	return false
}

// must also handle "*" value
func agentAuthorize(agentUUID string) {
	agentsListLock.Lock()
	defer agentsListLock.Unlock()

	is_any_agent := agentUUID == "*"
	for i := range agentsList {
		if agentsList[i].UUID == agentUUID || is_any_agent {
			agentsList[i].Status = "authorized"
			if !is_any_agent {
				return
			}
		}
	}
}

func agentUnauthorize(agentUUID string) {
	agentsListLock.Lock()
	defer agentsListLock.Unlock()

	is_any_agent := agentUUID == "*"
	for i := range agentsList {
		if agentsList[i].UUID == agentUUID || is_any_agent {
			agentsList[i].Status = "unauthorized"
			agentsList[i].WSConn.Close() // Can't even wait to get or send any transfer or data with "unauthorized" client
			if !is_any_agent {
				return
			}
		}
	}
}
