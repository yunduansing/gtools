package gen

import (
	"fmt"
	"testing"
	"time"
)

func TestSnowflakeID(t *testing.T) {
	for i := 0; i < 100; i++ {
		id1 := Uint64()
		fmt.Print(int(id1), " ")
		time.Sleep(1000 * time.Millisecond)
	}
}
