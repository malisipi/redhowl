# RedHowl Agent Backend API

|Method|URL                              |Description|
|------|---------------------------------|-----------|
||||
||**Version Information**||
|GET   |/version                         |Get API Version|
||||
||**Agent Communications**||
|GET   |/agents-com/ws                   |Get the auth status and connect WebSocker server|
|POST  |/agents-com/register             |Register itself as a agent to server|
|GET   |/agents-com/file-transfer        |Get file from server|
|POST  |/agents-com/file-transfer        |Upload file to server|

## GET /version

Gets server version, same with WebUI backend api

## POST /agents-com/register

To register agent itself for authentication

Request:
```json
{
    "uuid": "AGENT-LONG-UUID",
    "os": {
        "name": "CachyOS", // os name itself "Windows 11", "Ubuntu 22.04", "Arch", "CachyOS" etc.
        "kernel": "Linux 7.1.3-2-cachyos",
        "generic": "linux", // generic names like "linux", "windows", "macos"
        "arch": "amd64", // amd64, riscv64, arm64 etc
        "shell": "/bin/bash",
        "startupTimestamp": "2011-10-05T14:48:00.000Z"
    },
    "machine": {
        "id": "",
        "name": "",
        "vendor": "MONSTER",
        "modelName": "TULPAR T5"
    }
}
```

Response Code: 202

## GET /agents-com/file-transfer

For mostly downloading WASM executables from server

Query Parameters:

```json
{
    "uuid": "AGENT-LONG-UUID",
    "fileId": "FILE-ID"
}
```

Response: as `application/octet-stream`

Response Code: 200, 403, 404

## POST /agents-com/file-transfer

For mostly uploading task artifacts to server

Request: as `multipart/form-data`

```json
{
    "uuid": "AGENT-LONG-UUID",
    "fileId": "FILE-ID",
    "file": "" // raw binary
}
```

Response Code: 200, 403, 404

## GET /agents-com/ws

Query Parameters:

```json
{
    "uuid": "AGENT-LONG-UUID"
}
```

Response: It will be upgraded to WS when authenticated, if not will return 403.

### WebSocker Transfer Information

