package synctask

import (
	"fmt"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	t1:= func() {
		fmt.Println("t1")
		time.Sleep(2*time.Second)
	}
	t2:= func() {
		fmt.Println("t2")
		time.Sleep(4*time.Second)
	}
	Exec(t1,t2)

}

func TestExecWithReturn(t *testing.T) {
	t1:= func() interface{} {
		fmt.Println("t1")
		time.Sleep(2*time.Second)
		return 1
	}
	t2:= func() interface{}{
		fmt.Println("t2")
		time.Sleep(4*time.Second)
		return 2
	}
	res:=ExecWithReturn(t1,t2)
	fmt.Println(res)
}
