package pool

type GoroutinePool struct {
	Queue  chan func() error
	Number int
	Total  int

	result         chan error
	finishCallback func()
}

func (self *GoroutinePool) Init(number int, total int) {
	self.Queue = make(chan func() error, total)
	self.Number = number
	self.Total = total
	self.result = make(chan error, total)
}

func (self *GoroutinePool) Start() {
	//start several goroutines
	for i := 0; i < self.Number; i++ {
		go func() {
			for {
				task, ok := <-self.Queue
				if !ok {
					break
				}

				err := task()
				self.result <- err
			}
		}()
	}

	// get every task's result
	for j := 0; j < self.Total; j++ {
		_, ok := <-self.result
		if !ok {
			break
		}
	}

	// callback when every task are completed
	if self.finishCallback != nil {
		self.finishCallback()
	}
}

func (self *GoroutinePool) Stop() {
	close(self.Queue)
	close(self.result)
}

func (self *GoroutinePool) AddTask(task func() error) {
	self.Queue <- task
}

func (self *GoroutinePool) SetFinishCallback(callback func()) {
	self.finishCallback = callback
}
