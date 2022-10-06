package cronjob

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	c:=New()
	c.Start()
	defer c.Stop()
	c.AddTask("test-job",GetEverySecondSpec(30), func() {
		fmt.Println("test job task...at ",time.Now())
	})

	defaultWorker:=new(DefaultWorker)
	c.AddTaskWithData("test-job-data",GetEveryTimeSpec(18,32,20),defaultWorker)
	time.Sleep(2*time.Minute)
}
