package concurrent

type Pool2 struct {
	sem     chan struct{}
	argChan chan any
}

func NewPool2(size int) *Pool2 {
	inst := &Pool2{
		sem:     make(chan struct{}, size),
		argChan: make(chan any, size),
	}
	return inst
}

func (pool *Pool2) NewTask(task func(any), val interface{}) {
	select {
	case pool.sem <- struct{}{}:
		pool.argChan <- val
		go pool.Worker(task)
	}
}

func (pool *Pool2) Worker(task func(any)) {
	defer func() {
		<-pool.sem
	}()

	val := <-pool.argChan
	task(val)
}
