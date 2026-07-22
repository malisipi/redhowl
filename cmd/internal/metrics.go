package internal

import "time"

type MetricsNetwork struct {
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
	MAC  string `json:"mac"`
}

type MetricsMachine struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Vendor    string `json:"vendor"`
	ModelName string `json:"modelName"`
}

type MetricsOS struct {
	Name             string    `json:"name"`
	Kernel           string    `json:"kernel"`
	Generic          string    `json:"generic"`
	Arch             string    `json:"arch"`
	Shell            string    `json:"shell"`
	StartupTimestamp time.Time `json:"startupTimestamp"`
}

type MetricsUser struct {
	Name    string `json:"name"`
	UID     int    `json:"uid"`
	IsAdmin bool   `json:"isAdmin"`
}

type MetricsDisk struct {
	Used       float64 `json:"used"`
	Total      float64 `json:"total"`
	MountPoint string  `json:"mountPoint"`
}

type MetricsMemory struct {
	Used  float64 `json:"used"`
	Total float64 `json:"total"`
}
