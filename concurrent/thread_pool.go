package concurrent

import "sync"

type ThreadPool struct {
	sem     chan struct{}
	argChan chan any
	wg      sync.WaitGroup
	handler func(key any)
}

func NewThreadPool(size int, callback func(key any)) *ThreadPool {
	inst := &ThreadPool{
		sem:     make(chan struct{}, size),
		argChan: make(chan any, size),
		wg:      sync.WaitGroup{},
		handler: callback,
	}
	return inst
}

func (this *ThreadPool) NewTask(task func(any), val interface{}) {
	select {
	case this.sem <- struct{}{}:
		this.wg.Add(1)
		this.argChan <- val
		go this.Worker(task)
	}
}

func (this *ThreadPool) Worker(task func(any)) {
	defer func() {
		<-this.sem
	}()

	val := <-this.argChan
	task(val)
	this.handler(val)
	this.wg.Done()
}

func (this *ThreadPool) Wait() {
	this.wg.Wait()
	close(this.argChan)
}
