package nacos

import (
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"sync"
)

var nClient naming_client.INamingClient
var once sync.Once

// 注册nacos服务
func Register(config Config) (err error) {
	once.Do(func() {
		nClient, err = clients.NewNamingClient(
			vo.NacosClientParam{
				ClientConfig: &constant.ClientConfig{
					Username:            config.Username,
					Password:            config.Password,
					NamespaceId:         config.NamespaceId, //namespace id
					TimeoutMs:           5000,
					NotLoadCacheAtStart: true,
					LogDir:              "./log",
					CacheDir:            "./nacos/cache",
					//RotateTime:          "1h",
					//MaxAge:              3,
					LogLevel: "info",
				},
				ServerConfigs: []constant.ServerConfig{
					{
						IpAddr: config.ServerIp,
						Port:   config.ServerPort,
					},
				},
			},
		)
		if err != nil {
			return
		}
		var b bool
		b, err = nClient.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          config.ClientIp,
			Port:        config.ClientPort,
			Enable:      true,
			Healthy:     true,
			Weight:      10,
			Metadata:    config.Metadata,
			ServiceName: config.ServiceName,
			GroupName:   config.GroupName,
			//ClusterName: "default",
			Ephemeral: true,
		})
		if err != nil {
			return
		}
		if b != true {
			err = errors.New("nacos注册失败")
			return
		}
		return
	})
	return
}

// 注册nacos服务
func RegisterWithMultiServer(config ConfigMultiServer) (err error) {
	if len(config.Servers) == 0 {
		return errors.New("servers not be empty")
	}
	once.Do(func() {
		cfg := vo.NacosClientParam{
			ClientConfig: &constant.ClientConfig{
				Username:            config.Username,
				Password:            config.Password,
				NamespaceId:         config.NamespaceId, //namespace id
				TimeoutMs:           5000,
				NotLoadCacheAtStart: true,
				LogDir:              "./log",
				CacheDir:            "./nacos/cache",
				//RotateTime:          "1h",
				//MaxAge:              3,
				LogLevel: "debug",
			},
		}
		for _, v := range config.Servers {
			cfg.ServerConfigs = append(cfg.ServerConfigs, constant.ServerConfig{
				IpAddr: v.ServerIp,
				Port:   v.ServerPort,
			})
		}
		nClient, err = clients.NewNamingClient(cfg)
		if err != nil {
			return
		}
		var b bool
		b, err = nClient.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          config.ClientIp,
			Port:        config.ClientPort,
			Enable:      true,
			Healthy:     true,
			Weight:      10,
			Metadata:    config.Metadata,
			ServiceName: config.ServiceName,
			GroupName:   config.GroupName,
			//ClusterName: "default",
			Ephemeral: true,
		})
		if err != nil {
			return
		}
		if b != true {
			err = errors.New("nacos注册失败")
			return
		}
		return
	})
	return
}

// 服务解析
// 返回http://ip:port/
func Resolve(serviceName string) (string, error) {
	if nClient == nil {
		return "", errors.New("nacos未注册")
	}
	instance, err := nClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{ServiceName: serviceName})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://%s:%d/", instance.Ip, instance.Port), nil
}
