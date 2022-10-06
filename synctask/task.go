package synctask

import "sync"

func Exec(tasks ...func())  {
	taskSize:=len(tasks)
	var wg sync.WaitGroup
	wg.Add(taskSize)
	for i,_:=range tasks{

		go func(index int) {
			defer wg.Done()
			tasks[index]()
		}(i)
	}
	wg.Wait()
}

func ExecWithReturn(tasks ...func() interface{}) []interface{} {
	taskSize:=len(tasks)
	resList:=make([]interface{},taskSize)
	var wg sync.WaitGroup
	var lock sync.Mutex
	wg.Add(taskSize)
	for i,_:=range tasks{

		go func(index int) {
			defer wg.Done()
			res:=tasks[index]()
			lock.Lock()
			resList[index]=res
			lock.Unlock()
		}(i)
	}
	wg.Wait()
	return resList
}
