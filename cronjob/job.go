package cronjob

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
	"sync/atomic"
	"time"
)

type Job struct {
	cron *cron.Cron
	names sync.Map
	isRunning int32
	sync.Once
}

func GetEveryHourSpec(hour int) string {
	return fmt.Sprintf("0 0 0/%d * * ?",hour)
}

func GetEveryMinuteSpec(minute int) string {
	return fmt.Sprintf("0 0/%d * * * ?",minute)
}

func GetEverySecondSpec(second int) string {
	return fmt.Sprintf("%d * * * * ?",second)
}

func GetEveryTimeSpec(hour,minute,second int) string {
	return fmt.Sprintf("%d %d %d * * ?",second,minute,hour)
}

//You must implements all method
type Worker interface {
	Task(ctx context.Context,data chan interface{}) //task
	OnBefore()
	OnData(data interface{})  //operate data
	OnComplete()
	Lock() //lock must implement
	UnLock() //unlock must implement
}

type DefaultWorker struct {
	locker sync.Mutex
	Data chan interface{}
	startTime time.Time
}

func (d *DefaultWorker) Task(ctx context.Context, data chan interface{}) {
	ts:=time.Now().Unix()
	time.Sleep(5*time.Second)
	data<-ts
}

func (d *DefaultWorker) OnBefore()  {
	fmt.Println("on before")
	d.startTime=time.Now()
}
func (d *DefaultWorker) OnComplete()  {
	fmt.Println("on complete,time cost ",time.Since(d.startTime))
}
func (d *DefaultWorker) Lock()  {
	d.locker.Lock()
}
func (d *DefaultWorker) UnLock()  {
	d.locker.Unlock()
}
func (d *DefaultWorker) OnData(data interface{})  {
	fmt.Println("on data ts:",data)
}

func New() *Job {
	return &Job{
		cron: cron.New(cron.WithParser(cron.NewParser(
			cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		))),
	}
}
//add task with data
func (c *Job) AddTask(name,spec string,f func()) (cron.EntryID,error) {
	eid,ok:=c.names.Load(name)
	if ok{
		return eid.(cron.EntryID),nil
	}
	entryId,err:=c.cron.AddFunc(spec,f)
	if err!=nil{
		return 0, err
	}
	c.names.Store(name,entryId)
	return  entryId,err
}
//add task with data operation
func (c *Job) AddTaskWithData(name,spec string,w Worker) (cron.EntryID,error) {
	eid,ok:=c.names.Load(name)
	if ok{
		return eid.(cron.EntryID),nil
	}
	entryId,err:=c.cron.AddFunc(spec, func() {
		w.Lock()
		w.OnBefore()
		data:=make(chan interface{})
		defer func() {
			close(data)
			w.OnComplete()
			w.UnLock()
		}()
		go func() {
			for d:=range data{
				w.OnData(d)
			}
		}()
		w.Task(context.Background(),data)
	})
	if err!=nil{
		return 0, err
	}
	c.names.Store(name,entryId)

	return  entryId,err
}
func (c *Job) RemoveTask(id cron.EntryID)  {
	c.cron.Remove(id)
}
func (c *Job) Start()  {
	c.Do(func() {
		c.isRunning=1
		c.cron.Start()
	})

}
func (c *Job) Stop()  {
	c.cron.Stop()
	atomic.StoreInt32(&c.isRunning,0)
}