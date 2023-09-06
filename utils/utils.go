package utils

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"time"
)

var (
	TimeLayout = "2006-01-02 13:14:05"
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
	return t.Format(TimeLayout)
}

func StringToTimeNoErr(str string) time.Time {
	t, _ := time.Parse(TimeLayout, str)
	return t
}

func StringToChineseTime(str string) (time.Time, error) {
	t, err := time.Parse(TimeLayout, str)
	return t, err
}

func ToJson(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return ByteToString(bytes)
}
