package utils

import (
	"fmt"
	"math"
	"testing"
)

func TestSnowflakeID(t *testing.T) {
	var i int64 = math.MaxInt64 - 1
	fmt.Println(i, int(i))
	for i := 0; i < 100; i++ {
		id1 := Uint64()
		fmt.Println(int(id1), " ")
	}
}
