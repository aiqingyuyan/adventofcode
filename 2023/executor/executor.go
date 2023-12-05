package executor

import (
	"fmt"
)

type Executor interface {
	Run(taskEmitter <-chan TaskFunc) int
}

type TaskFunc func() int

type executor struct {
	executorChan   chan chan TaskFunc
	resultChan     chan int
	numberOfWorker int
}

func New(numberOfWorker int) Executor {
	if numberOfWorker == 0 {
		panic(fmt.Errorf("invalid worker number: %d, must be > 0", numberOfWorker))
	}

	executor := executor{
		executorChan:   make(chan chan TaskFunc),
		resultChan:     make(chan int, numberOfWorker),
		numberOfWorker: numberOfWorker,
	}

	for i := 0; i < numberOfWorker; i++ {
		startWorker(executor.executorChan, executor.resultChan /*, processLineFunc*/)
	}

	return &executor
}

func startWorker(executorChan chan chan TaskFunc, resultChan chan int) {
	workerTaskChan := make(chan TaskFunc)

	go func() {
		for {
			executorChan <- workerTaskChan

			taskFunc := <-workerTaskChan

			resultChan <- taskFunc()
		}
	}()
}

func (e *executor) Run(taskEmitter <-chan TaskFunc) int {
	result := 0

	var doneEmittingTask bool
	var workerQueue []chan TaskFunc
	var taskQueue []TaskFunc

	for {
		if doneEmittingTask && len(taskQueue) == 0 && len(workerQueue) == e.numberOfWorker && len(e.resultChan) == 0 {
			break
		}

		var activeTask TaskFunc
		var activeWorker chan TaskFunc

		if len(taskQueue) > 0 && len(workerQueue) > 0 {
			activeTask = taskQueue[0]
			activeWorker = workerQueue[0]
		}

		select {
		case task, ok := <-taskEmitter:
			if ok {
				taskQueue = append(taskQueue, task)
			} else {
				doneEmittingTask = true
			}
		case workerChan := <-e.executorChan:
			workerQueue = append(workerQueue, workerChan)
		case activeWorker <- activeTask:
			taskQueue = taskQueue[1:]
			workerQueue = workerQueue[1:]
		case number := <-e.resultChan:
			result += number
		}
	}

	for i := 0; i < e.numberOfWorker; i++ {
		go func(id int) {
			e.executorChan <- workerQueue[id]
		}(i)
	}

	return result
}
