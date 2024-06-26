package workqueue

type worker interface {
	run()
	insertTask(func())
	finish()
}

type defaultWorker struct {
	tasks chan func()
	p     *Pool
}

func (w *defaultWorker) run() {
	w.p.addRunning(1)
	go func() {
		defer func() {
			w.p.addRunning(-1)
			w.p.Pool.Put(w)
			w.p.Cond.Signal() //notice other goroutines call
		}()
		for t := range w.tasks {

			t()
			//执行完毕后存入队
			if ok := w.p.revertWorker(w); !ok {
				return
			}
		}
	}()

}

func (w *defaultWorker) insertTask(t func()) {
	w.tasks <- t
}

func (w *defaultWorker) finish() {
	w.tasks = nil
}
