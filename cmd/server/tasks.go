package main

import "time"

func agentTaskExist(taskId string) bool {
	return true
}

func agentTaskDelete(taskId string) {
	return
}

func agentTaskGet(taskId string) {
	return
}

func agentTaskRun(task ReqTask) string {
	return "0000"
}

func agentTaskList() []TaskGeneric {
	return []TaskGeneric{
		TaskGeneric{
			ID:          "123",
			Name:        "Task Name",
			Description: "Description",
			Timestamp:   time.Now(),
			Status: TaskStatus{
				Running:  12,
				Finished: 4,
				Failed:   2,
			},
		},
	}
}
