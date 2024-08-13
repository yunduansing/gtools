package utils

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"net/http"
	"reflect"
	"time"
	"unsafe"
)

var (
	ChineseTimeLayout = "2006-01-02 15:04:05"
)

func YuanToFenNoErr(m string) int64 {
	d, _ := decimal.NewFromString(m)
	return d.Mul(decimal.NewFromInt(100)).IntPart()
}

func YuanToFen(m string) (int64, error) {
	d, err := decimal.NewFromString(m)
	if err != nil {
		return 0, err
	}
	return d.Mul(decimal.NewFromInt(100)).IntPart(), nil
}

func FenInt64ToYuanStr(m int64) string {
	return decimal.NewFromInt(m).Div(decimal.NewFromInt(100)).StringFixed(2)
}

// TimeToChineseStr 2006-01-02 13:14:05
func TimeToChineseStr(t time.Time) string {
	return t.Format(ChineseTimeLayout)
}

func StringToTimeNoErr(str string) time.Time {
	t, _ := time.Parse(ChineseTimeLayout, str)
	return t
}

func StringToChineseTime(str string) (time.Time, error) {
	t, err := time.Parse(ChineseTimeLayout, str)
	return t, err
}

func ToJsonString(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return ByteToString(bytes)
}

// ByteToString String and []byte buffers may converted without memory allocations
// This is an unsafe way, the result string and []byte buffer share the same bytes.
// Please make sure not to modify the bytes in the []byte buffer if the string still survives!
func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToByte String and []byte buffers may converted without memory allocations
// This is an unsafe way, the result string and []byte buffer share the same bytes.
// Please make sure not to modify the bytes in the []byte buffer if the string still survives!
func StringToByte(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

func GetClientIp(req http.Request) string {
	forwarded := req.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return req.RemoteAddr

}
