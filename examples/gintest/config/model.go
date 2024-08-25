package config

type ServiceConfig struct {
	ServiceName   string `json:"serviceName"`
	IsMetricsOpen bool   `json:"isMetricsOpen"`
	IsTracingOpen bool   `json:"isTracingOpen"`
}
