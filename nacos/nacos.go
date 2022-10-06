package nacos

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yunduansing/gocommon/httputils"
	"github.com/yunduansing/gocommon/logger"
	"go.uber.org/zap"
	"net/http"
)

type Config struct {
	ServerIp   string `json:"server_ip"`
	ServerPort int    `json:"server_port"`
}

func (c *Config) Register(serviceIp string, servicePort int, serviceName string) error {
	path := "/nacos/v1/ns/instance"
	url := fmt.Sprintf("http://%s:%d%s?port=%d&healthy=true&ip=%s&weight=1.0&serviceName=%s&encoding=GBK&namespaceId=&ephemeral=false", c.ServerIp, c.ServerPort, path, servicePort, serviceIp, serviceName)
	_, code, err := httputils.HttpPost(url, nil, nil)
	if err != nil || code != http.StatusOK {
		if err != nil {
			logger.Logger.Error("service register to nacos", zap.Error(err))
		}
		return errors.New("error")
	}
	return nil
}

type ClientHost struct {
	Valid      bool        `json:"valid"`
	Marked     bool        `json:"marked"`
	InstanceId string      `json:"instanceId"`
	Port       int         `json:"port"`
	Ip         string      `json:"ip"`
	Weight     float64     `json:"weight"`
	Metadata   interface{} `json:"metadata"`
}

type Result struct {
	Dom             string       `json:"dom"`
	CacheMillis     int          `json:"cacheMillis"`
	UseSpecifiedURL bool         `json:"useSpecifiedURL"`
	Checksum        string       `json:"checksum"`
	LastRefTime     int          `json:"lastRefTime"`
	Env             string       `json:"env"`
	Clusters        string       `json:"clusters"`
	Hosts           []ClientHost `json:"hosts"`
}

// Resolve 返回http://ip:port
func (c *Config) Resolve(serviceName string) string {
	url := fmt.Sprintf("http://%s:%d/nacos/v1/ns/instance/list?serviceName=%s", c.ServerIp, c.ServerPort, serviceName)
	data, code, err := httputils.HttpGet(url, nil)
	if err != nil || code != http.StatusOK {
		return ""
	}
	var result Result
	err = json.Unmarshal(data, &result)
	if err != nil || len(result.Hosts) == 0 {
		return ""
	}

	return fmt.Sprintf("http://%s:%d", result.Hosts[0].Ip, result.Hosts[0].Port)
}
