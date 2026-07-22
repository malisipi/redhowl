package internal

type ReqAgentRegister struct {
	UUID    string         `json:"uuid"`
	User    MetricsUser    `json:"user"`
	OS      MetricsOS      `json:"os"`
	Machine MetricsMachine `json:"machine"`
}

type WSType int

const (
	WSTypeUndef WSType = iota
	WSTypeMetricSend
)

type WSTypeHeader struct {
	Type WSType `json:"type"`
}

type WSMetricSend struct {
	WSTypeHeader
	CPU     float64        `json:"cpu"`
	Memory  MetricsMemory  `json:"memory"`
	Disk    MetricsDisk    `json:"disk"`
	Network MetricsNetwork `json:"network"`
}
