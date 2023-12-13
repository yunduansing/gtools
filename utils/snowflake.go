package utils

import (
	"errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sony/sonyflake"
	"net"
	"time"
)

// UUID 生成UUID
func UUID() string {
	return uuid.NewV4().String()
}

var (
	snowflake = newSnowflake()
)

func newSnowflake() *sonyflake.Sonyflake {
	//startTime, _ := time.ParseInLocation("2006-01-02", "2021-12-01", time.Local)

	snowflake := sonyflake.NewSonyflake(sonyflake.Settings{StartTime: time.Now()})
	if snowflake == nil {
		panic("创建snowflake失败, snowflake实例为nil")
	}

	return snowflake
}

// SnowflakeID 生成雪花ID
func SnowflakeID() (uint64, error) {
	return snowflake.NextID()
}

func privateIPv4() (net.IP, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if isPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ip address")
}

func isPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func lower16BitPrivateIP() (uint16, error) {
	ip, err := privateIPv4()
	if err != nil {
		return 0, err
	}

	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

var (
	startTime, _                      = time.ParseInLocation("2006-01-02", "2021-12-01", time.Local)
	core         *sonyflake.Sonyflake = sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: func() (uint16, error) {
			return 0, nil
		},
		CheckMachineID: func(u uint16) bool {
			return u == 0
		},
	})
)

func Init(serviceName string, serviceId uint64) {
	var v int64 = 1
	for k, c := range serviceName {
		v = v * int64(c) * int64(k+1)
	}
	id := uint16(v + int64(serviceId))
	core = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Now(),
		MachineID: func() (uint16, error) {
			return id, nil
		},
		CheckMachineID: func(u uint16) bool {
			return u == id
		},
	})
}

func InitV2(serviceName string) {
	var v int64 = 1
	for k, c := range serviceName {
		v = v * int64(c) * int64(k+1)
	}

	ip, _ := lower16BitPrivateIP()
	id := uint16(v) + ip
	//logger.Info("snow flake lower ip=", ip)
	core = sonyflake.NewSonyflake(sonyflake.Settings{
		//StartTime: time.Now(),
		MachineID: func() (uint16, error) {
			return id, nil
		},
		CheckMachineID: func(u uint16) bool {
			return u == id
		},
	})
}

func Int64() int64 {
	id, _ := core.NextID()
	return int64(id)
}

func Uint64() uint64 {
	id, _ := core.NextID()
	return id
}
