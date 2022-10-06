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
}

func newWorkerQueue(max, blocking int32) (w workerQueues) {
	w.max = max
	w.blocking = blocking
	go w.Run()
	return
}

func (w workerQueues) Running() int32 {
	return atomic.LoadInt32(&w.running)
}

func (w workerQueues) Submit(t func()) {
	w.Lock()
	if atomic.LoadInt32(&w.running) >= w.blocking {
		w.tasks = append(w.tasks, t)
	} else {
		go func() {
			fmt.Println("running", w.running)
			atomic.AddInt32(&w.running, 1)
			t()
			time.Sleep(300 * time.Millisecond)
			atomic.AddInt32(&w.running, -1)
		}()
	}

	w.Unlock()
}

func (w workerQueues) Run() {
	w.Lock()
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
						go func() {
							fmt.Println("running", w.running)
							atomic.AddInt32(&w.running, 1)
							t()
							time.Sleep(300 * time.Millisecond)
							atomic.AddInt32(&w.running, -1)
						}()
					} else {
						break taskListLoop
					}
				}
				timer.Reset(time.Second)
			}
		}
	}

	w.Unlock()
	return
}

func (w workerQueues) Release() {
	atomic.StoreInt32(&w.state, 0)
}

func (w workerQueues) PopTask() (t func()) {
	w.Lock()
	t = w.tasks[0]
	w.tasks = w.tasks[1:]
	w.Unlock()
	return
}

func (w workerQueues) Cap() (l int32) {

	return atomic.LoadInt32(&w.max)
}

func (w workerQueues) Free() (l int32) {

	return atomic.LoadInt32(&w.max) - atomic.LoadInt32(&w.running)
}
