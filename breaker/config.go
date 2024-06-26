package breaker

import (
	"github.com/sony/gobreaker"
	"github.com/yunduansing/gtools/redistool"
)

type Config struct {
	Sensitivity float64 //熔断器灵敏度，数值越大越敏感
	MaxRequest  uint32
	Redis       *redistool.RedisCli
}

const DefaultSensitivity = 0.8

var config *Config

func NewConfig(conf *Config) {
	config = conf
	if config.Sensitivity == 0 {
		config.Sensitivity = DefaultSensitivity
	}
}

type Breaker struct {
	*gobreaker.CircuitBreaker
}

func NewBreaker(conf *Config) *Breaker {
	var setting gobreaker.Settings
	if conf.Sensitivity == 0 {
		conf.Sensitivity = DefaultSensitivity
	}
	if conf.MaxRequest > 0 {
		setting.MaxRequests = conf.MaxRequest
	}

	setting.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return failureRatio > config.Sensitivity
	}
	return &Breaker{CircuitBreaker: gobreaker.NewCircuitBreaker(setting)}
}
