package executor

import "fmt"

func startWorker(schedulerChan chan chan *string, resultChan chan int, taskFunc TaskFunc) {
	workerTaskChan := make(chan *string)

	go func() {
		for {
			schedulerChan <- workerTaskChan

			line := <-workerTaskChan

			resultChan <- taskFunc(line)
		}
	}()
}

type TaskFunc func(*string) int

func Run(numOfWorker int, lineEmitter <-chan *string, processLineFunc TaskFunc) int {
	if numOfWorker == 0 {
		panic(fmt.Errorf("invalid worker number: %d, must be > 0", numOfWorker))
	}

	result := 0

	resultChan := make(chan int, numOfWorker)
	schedulerChan := make(chan chan *string)
	for i := 0; i < numOfWorker; i++ {
		startWorker(schedulerChan, resultChan, processLineFunc)
	}

	var doneEmittingLine bool
	var workerQueue []chan *string
	var taskQueue []*string

	for {
		if doneEmittingLine && len(taskQueue) == 0 && len(workerQueue) == numOfWorker && len(resultChan) == 0 {
			break
		}

		var activeTask *string
		var activeWorker chan *string

		if len(taskQueue) > 0 && len(workerQueue) > 0 {
			activeTask = taskQueue[0]
			activeWorker = workerQueue[0]
		}

		select {
		case line, ok := <-lineEmitter:
			if ok {
				taskQueue = append(taskQueue, line)
			} else {
				doneEmittingLine = true
			}
		case workerChan := <-schedulerChan:
			workerQueue = append(workerQueue, workerChan)
		case activeWorker <- activeTask:
			taskQueue = taskQueue[1:]
			workerQueue = workerQueue[1:]
		case number := <-resultChan:
			result += number
		}
	}

	return result
}
