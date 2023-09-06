package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

var (
	TimeLayout = "2006-01-02 13:14:05"
)

func YuanToFenNoErr(m string) int64 {
	d, _ := strconv.ParseFloat(m, 64)
	return int64(d * 100)
}

func YuanToFen(m string) (int64, error) {
	d, err := strconv.ParseFloat(m, 64)
	if err != nil {
		return 0, err
	}
	return int64(d * 100), nil
}

func FenInt64ToYuanStr(m int64) string {
	return fmt.Sprintf("%.2f", float64(m)/100.0)
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

func ToJsonString(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return ByteToString(bytes)
}
