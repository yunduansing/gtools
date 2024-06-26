package workqueue

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type workerQueues struct {
	sync.Mutex
	tasks    []func()
	max      int32
	running  int32
	blocking int32
	state    int32 //0-closed  1-open
	delay time.Duration
}

func newWorkerQueue(max, blocking int32,delay time.Duration) (w *workerQueues) {
	w=new(workerQueues)
	w.max = max
	w.blocking = blocking
	w.delay=delay
	go w.Run()
	return
}

func (w *workerQueues) Running() int32 {
	return atomic.LoadInt32(&w.running)
}

func (w *workerQueues) Submit(t func()) {
	w.Lock()
	defer func() {
		w.Unlock()
	}()
	if int32(len(w.tasks))>w.max{
		return
	}
	w.tasks = append(w.tasks, t)
}

func (w *workerQueues) Run() {
	if atomic.LoadInt32(&w.state) == 0 {
		atomic.AddInt32(&w.state, 1)
		timer := time.NewTimer(time.Second)
		for {
			select {
			case <-timer.C:
			taskListLoop:
				for {
					if len(w.tasks) > 0 && atomic.LoadInt32(&w.running) < w.blocking {
						t := w.PopTask()
						atomic.AddInt32(&w.running, 1)
						go func() {
							defer func() {
								atomic.AddInt32(&w.running, -1)
							}()
							fmt.Println("running", w.running)

							t()
							time.Sleep(w.delay)
						}()
					} else {
						break taskListLoop
					}
				}
				timer.Reset(w.delay)
			}
		}
	}
	return
}

func (w *workerQueues) Release() {
	atomic.StoreInt32(&w.state, 0)
}

func (w *workerQueues) PopTask() (t func()) {
	w.Lock()
	t = w.tasks[0]
	w.tasks = w.tasks[1:]
	w.Unlock()
	return
}

func (w *workerQueues) Cap() (l int32) {

	return atomic.LoadInt32(&w.max)
}

func (w *workerQueues) Free() (l int32) {

	return atomic.LoadInt32(&w.max) - atomic.LoadInt32(&w.running)
}
