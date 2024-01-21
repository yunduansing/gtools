package workqueue

import (
	"sync"
	"sync/atomic"
)

type Pool struct {
	cap    int32
	closed int32
	sync.Pool
	sync.Cond
	sync.Locker
	running     int32
	workers     []worker
	maxBlocking int32
}

func NewPool(size int) *Pool {
	p := &Pool{
		cap:    int32(size),
		closed: 0,
		Pool:   sync.Pool{},
		Cond:   sync.Cond{},
		Locker: &sync.Mutex{},
	}

	p.Pool.New = func() any {
		return defaultWorker{
			tasks: make(chan func(), 0),
			p:     p,
		}
	}

	return p
}

func (p *Pool) IsClosed() bool {
	return atomic.LoadInt32(&p.closed) == 1
}

func (p *Pool) Submit(t func()) {
	w := p.getWorker()
	if w != nil {
		w.insertTask(t)
	}
}

func (p *Pool) getWorker() worker {
	for {
		p.Lock()
		if p.Running() < p.Cap() && len(p.workers) > 0 {

			w := p.workers[0] //取最先进入队列的worker
			p.workers[0] = nil
			p.workers = p.workers[1:]
			return w
		}
		if atomic.LoadInt32(&p.running) < int32(p.cap) {
			return p.Pool.Get().(worker)
		}
		//等待
		p.Wait()
	}
}

func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.cap))
}

func (p *Pool) Running() int {
	return int(atomic.LoadInt32(&p.running))
}

func (p *Pool) revertWorker(w worker) bool {
	if (p.Running() > p.Cap()) || p.IsClosed() {
		p.Broadcast()
		return false
	}
	p.Lock()
	if p.closed == 1 {
		p.Unlock()
		return false
	}
	p.workers = append(p.workers, w)
	p.Signal()
	p.Unlock()
	return true
}

func (p *Pool) addRunning(delta int32) {
	atomic.AddInt32(&p.running, delta)
}
