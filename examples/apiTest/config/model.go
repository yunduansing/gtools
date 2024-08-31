package config

type ServiceConfig struct {
	ServiceName          string `json:"serviceName"`
	IsMetricsOpen        bool   `json:"isMetricsOpen"`
	IsTracingOpen        bool   `json:"isTracingOpen"`
	IsRequestLimiterOpen bool   `json:"isRequestLimiterOpen"`
	Env                  string `json:"env"`
}

type UptraceConfig struct {
	Version string `json:"version"`
	Dsn     string `json:"dsn"`
}
