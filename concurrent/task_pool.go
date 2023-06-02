package concurrent

import (
	"errors"
	"sync"
)

type Task struct {
	WorkerId int
	Handler  func(val any)
	Args     any
}

type Worker struct {
	id       int
	taskChan chan Task
	wg       *sync.WaitGroup
}

func (w *Worker) Start() {
	go func() {
		defer w.wg.Done()
		for task := range w.taskChan {
			task.Handler(task.Args)
		}
	}()
}

type TaskPool struct {
	workerNum int
	taskChan  chan Task
	workers   []*Worker
	wg        sync.WaitGroup
}

func NewTaskPool(workerNum int, taskNum int) *TaskPool {
	taskChan := make(chan Task, taskNum)
	workers := make([]*Worker, workerNum)
	for i := 0; i < workerNum; i++ {
		workers[i] = &Worker{
			id:       i,
			taskChan: make(chan Task),
			wg:       &sync.WaitGroup{},
		}
		workers[i].wg.Add(1)
		workers[i].Start()
	}
	inst := &TaskPool{
		workerNum: workerNum,
		taskChan:  taskChan,
		workers:   workers,
	}
	inst.start()
	return inst
}

func (p *TaskPool) start() {
	go func() {
		for task := range p.taskChan {
			worker := p.workers[task.WorkerId]
			worker.taskChan <- task
			p.wg.Done()
		}
	}()
}

func (p *TaskPool) AddTask(task Task) error {
	if task.WorkerId > len(p.workers) {
		return errors.New("WorkerIdInvalid")
	}

	p.wg.Add(1)
	p.taskChan <- task

	return nil
}

func (p *TaskPool) Wait() {
	close(p.taskChan)
	p.wg.Wait()
	for _, worker := range p.workers {
		close(worker.taskChan)
		worker.wg.Wait()
	}
}
