package gen

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sony/sonyflake"
	"reflect"
	"regexp"
	"time"
	"unsafe"
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

func Int64() int64 {
	id, _ := core.NextID()
	return int64(id)
}

func Uint64() uint64 {
	id, _ := core.NextID()
	return id
}

// ByteToString String and []byte buffers may converted without memory allocations
//This is an unsafe way, the result string and []byte buffer share the same bytes.
//Please make sure not to modify the bytes in the []byte buffer if the string still survives!
func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToByte String and []byte buffers may converted without memory allocations
//This is an unsafe way, the result string and []byte buffer share the same bytes.
//Please make sure not to modify the bytes in the []byte buffer if the string still survives!
func StringToByte(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

// ValidPhoneNumber 正则验证手机号
func ValidPhoneNumber(phone string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)
}
