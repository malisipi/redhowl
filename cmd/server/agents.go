package main

import (
	"redhowl/cmd/internal"
	"time"
)

func agentGetList() []Agent {
	return []Agent{
		Agent{
			UUID:               "agent-1",
			Status:             "authorized",
			ConnectedTimestamp: time.Now(),
			Metrics: AgentMetrics{
				CPU: .452,
				Memory: internal.MetricsMemory{
					Used:  3.4,
					Total: 3.9,
				},
				Disk: internal.MetricsDisk{
					Used:       248,
					Total:      480,
					MountPoint: "/",
				},
				User: internal.MetricsUser{
					Name:    "redwolf",
					UID:     1000,
					IsAdmin: false,
				},
				OS: internal.MetricsOS{
					Name:             "CachyOS",
					Kernel:           "linux",
					Generic:          "linux",
					Arch:             "amd64",
					Shell:            "/bin/bash",
					StartupTimestamp: time.Now(),
				},
				Machine: internal.MetricsMachine{
					ID:        "machine-id",
					Name:      "eye-of-the-wolf",
					Vendor:    "MONSTER",
					ModelName: "Tulpar",
				},
				Network: internal.MetricsNetwork{
					IPv4: "127.0.0.1",
					IPv6: "::0",
					MAC:  "00:00:00:00:00:00",
				},
			},
		},
		Agent{
			UUID:               "agent-2",
			Status:             "authorized",
			ConnectedTimestamp: time.Now(),
			Metrics: AgentMetrics{
				CPU: .252,
				Memory: internal.MetricsMemory{
					Used:  5.4,
					Total: 7.9,
				},
				Disk: internal.MetricsDisk{
					Used:       278,
					Total:      880,
					MountPoint: "/",
				},
				User: internal.MetricsUser{
					Name:    "redwolf",
					UID:     1000,
					IsAdmin: false,
				},
				OS: internal.MetricsOS{
					Name:             "CachyOS",
					Kernel:           "linux",
					Generic:          "macos",
					Arch:             "arm64",
					Shell:            "/bin/bash",
					StartupTimestamp: time.Now(),
				},
				Machine: internal.MetricsMachine{
					ID:        "machine-id",
					Name:      "eye-of-the-wolf",
					Vendor:    "MONSTER",
					ModelName: "Tulpar",
				},
				Network: internal.MetricsNetwork{
					IPv4: "127.0.0.1",
					IPv6: "::0",
					MAC:  "00:00:00:00:00:00",
				},
			},
		},
		Agent{
			UUID:               "agent-3",
			Status:             "authorized",
			ConnectedTimestamp: time.Now(),
			Metrics: AgentMetrics{
				CPU: .85,
				Memory: internal.MetricsMemory{
					Used:  15.4,
					Total: 15.7,
				},
				Disk: internal.MetricsDisk{
					Used:       948,
					Total:      956,
					MountPoint: "C:\\",
				},
				User: internal.MetricsUser{
					Name:    "redwolf",
					UID:     1000,
					IsAdmin: false,
				},
				OS: internal.MetricsOS{
					Name:             "CachyOS",
					Kernel:           "linux",
					Generic:          "windows",
					Arch:             "amd64",
					Shell:            "/bin/bash",
					StartupTimestamp: time.Now(),
				},
				Machine: internal.MetricsMachine{
					ID:        "machine-id",
					Name:      "eye-of-the-wolf",
					Vendor:    "MONSTER",
					ModelName: "Tulpar",
				},
				Network: internal.MetricsNetwork{
					IPv4: "127.0.0.1",
					IPv6: "::0",
					MAC:  "00:00:00:00:00:00",
				},
			},
		},
	}
}

func agentExist(agentUUID string) bool {
	return agentUUID != ""
}

// must also handle "any" value
func agentAuthorize(agentUUID string) {
	return
}

func agentUnauthorize(agentUUID string) {
	return
}
