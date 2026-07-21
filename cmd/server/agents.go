package main

import "time"

func agentGetList() []Agent {
	return []Agent{
		Agent{
			UUID:               "agent-1",
			Status:             "authorized",
			ConnectedTimestamp: time.Now(),
			Metrics: AgentMetrics{
				CPU: .452,
				Memory: MetricsMemory{
					Used:  3.4,
					Total: 3.9,
				},
				Disk: MetricsDisk{
					Used:       248,
					Total:      480,
					MountPoint: "/",
				},
				User: MetricsUser{
					Name:    "redwolf",
					UID:     1000,
					IsAdmin: false,
				},
				OS: MetricsOS{
					Name:             "CachyOS",
					Kernel:           "linux",
					Generic:          "linux",
					Arch:             "amd64",
					Shell:            "/bin/bash",
					StartupTimestamp: time.Now(),
				},
				Machine: MetricsMachine{
					ID:        "machine-id",
					Name:      "eye-of-the-wolf",
					Vendor:    "MONSTER",
					ModelName: "Tulpar",
				},
				Network: MetricsNetwork{
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
				Memory: MetricsMemory{
					Used:  5.4,
					Total: 7.9,
				},
				Disk: MetricsDisk{
					Used:       278,
					Total:      880,
					MountPoint: "/",
				},
				User: MetricsUser{
					Name:    "redwolf",
					UID:     1000,
					IsAdmin: false,
				},
				OS: MetricsOS{
					Name:             "CachyOS",
					Kernel:           "linux",
					Generic:          "macos",
					Arch:             "arm64",
					Shell:            "/bin/bash",
					StartupTimestamp: time.Now(),
				},
				Machine: MetricsMachine{
					ID:        "machine-id",
					Name:      "eye-of-the-wolf",
					Vendor:    "MONSTER",
					ModelName: "Tulpar",
				},
				Network: MetricsNetwork{
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
				Memory: MetricsMemory{
					Used:  15.4,
					Total: 15.7,
				},
				Disk: MetricsDisk{
					Used:       948,
					Total:      956,
					MountPoint: "C:\\",
				},
				User: MetricsUser{
					Name:    "redwolf",
					UID:     1000,
					IsAdmin: false,
				},
				OS: MetricsOS{
					Name:             "CachyOS",
					Kernel:           "linux",
					Generic:          "windows",
					Arch:             "amd64",
					Shell:            "/bin/bash",
					StartupTimestamp: time.Now(),
				},
				Machine: MetricsMachine{
					ID:        "machine-id",
					Name:      "eye-of-the-wolf",
					Vendor:    "MONSTER",
					ModelName: "Tulpar",
				},
				Network: MetricsNetwork{
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
