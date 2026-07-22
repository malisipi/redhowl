# RedHowl WebUI Backend API

|Method|URL                              |Description|
|------|---------------------------------|-----------|
||||
||**Version Information**||
|GET   |/version                         |Get API Version|
||||
||**Agent Information**||
|GET   |/agents                          |Get agents list and generic data|
|POST  |/agents/authorize                |Authorize the agent|
|POST  |/agents/unauthorize              |Unauthorize the agent|
||||
||**Task Management**||
|POST  |/agents/tasks/                   |Upload an WASM file and run on agents|
|GET   |/agents/tasks/                   |Get task list|
|GET   |/agents/tasks/{taskId}           |Get task information|
|DELETE|/agents/tasks/{taskId}           |Delete task information|

## GET /version

Gets version of current API.

Response:
```json
{
    "version": [1,0,0]
}
```

Response Code: 200

## GET /agents

Gets the list of agents with filtering options.

Query Parameters:
```json
{
    "status": "authorized", // "*" (default), "authorized", "unauthorized"
    "uuid": "*", // "*" (default) | "AGENT-UUID"
}
```
Response:
```json
{
    "agents": [
        {
            "uuid": "AGENT-LONG-UUID",
            "status": "authorized",
            "connectedTimestamp": "2011-10-05T14:48:00.000Z", // it changes also status changes
            "metrics": {
                "cpu": 65, // percentage, via WSMetricSend
                "memory": { // via WSMetricSend
                    "used": 5.2,
                    "total": 8 // as GiB
                },
                "disk": { // via WSMetricSend
                    "used": 127.8, // as also GiB
                    "total": 937.4,
                    "mountPoint": "C:\\"
                },
                "user": { // via register
                    "name": "redwolf",
                    "uid": 255,
                    "isAdmin": false
                },
                "os": { // via Register
                    "name": "CachyOS", // os name itself "Windows 11", "Ubuntu 22.04", "Arch", "CachyOS" etc.
                    "kernel": "Linux 7.1.3-2-cachyos",
                    "generic": "linux", // generic names like "linux", "windows", "macos"
                    "arch": "amd64", // amd64, riscv64, arm64 etc
                    "shell": "/bin/bash",
                    "startupTimestamp": "2011-10-05T14:48:00.000Z"
                },
                "machine": { // via Register
                    "id": "",
                    "name": "",
                    "vendor": "MONSTER",
                    "modelName": "TULPAR T5"
                },
                "network": { // via WSMetricSend
                    "ipv4": "192.168.1.1",
                    "ipv6": "::0",
                    "localIpv4": "192.168.1.1",
                    "localIpv6": "::0",
                    "mac": "00:00:00:00:00:00"
                }
            }
        }
    ]
}
```

Response Code: 200

## POST /agents/authorize

Authorizes a agent.

Request:
```json
{
    "uuid": "*", // "*" (default) | "AGENT-UUID"
}
```

Response Code: 204

## POST /agents/unauthorize

Unauthorizes a agent.

Request:
```json
{
    "uuid": "*", // "*" (default) | "AGENT-UUID"
}
```

Response Code: 204

## POST /agents/tasks/

Uploads a WASM module and run a task with the module.

Request: as `multipart/form-data`
```json
// "json" part
{
    "name": "Task name",
    "description": "Task Description",
    "permissions": {
        "filesystem": false,
        "network": false,
        "execution": [], // binaries, or directly "*"
        "artifacts": true
    },
    "agents": [], // uuids or directly "*"
}
// "wasm" as binary
```
Response:
```json
{
    "task": {
        "id": ""
    }
}
```

Response Code: 202

## GET /agents/tasks/

Gets tasks list and basic state information of tasks.

Response:
```json
{
    "tasks": [
        {
            "id": "",
            "name": "...",
            "description": "...",
            "timestamp": "2011-10-05T14:48:00.000Z",
            "status": {
                "running": 16,
                "finished": 5,
                "failed": 7
            }
        }
    ]
}
```

Response Code: 200

## GET /agents/tasks/{taskId}

Gets detailed information of task.

Response:
```json
{
    "task": {
        "id": "",
        "name": "...",
        "description": "...",
        "permissions": {...},
        "agents": [...],
        "timestamp": "2011-10-05T14:48:00.000Z",
        "status": {
            "running": 16,
            "finished": 5,
            "failed": 7
        },
        "reports": {
            [agentUuid]: {
                "logs": [],
                "wasmLogs": [],
                "artifacts": [
                    {
                        "id": "",
                        "name": "",
                        "timestamp": "2011-10-05T14:48:00.000Z"
                    }
                ],
                "status": "running" // "running", "finished", "failed" (due offline or not exist)
            }
        }
    }
}
```

Response Code: 200, 404

## DELETE /agents/tasks/{taskId}

Deletes the task and related artifacts.

Response Code: 204, 404