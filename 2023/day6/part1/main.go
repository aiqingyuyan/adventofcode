package main

import (
	"log"
	"math"
	"path/filepath"
	"regexp"
	"sync"
	"yanyu/aoc/2023/executor"
	"yanyu/aoc/2023/util"
)

type computationParameter struct {
	time     int
	distance int
}

func bufferLines(lineEmitter <-chan *string) []string {
	var linesBuffer []string
	for line := range lineEmitter {
		linesBuffer = append(linesBuffer, *line)
	}

	return linesBuffer
}

func transformLinesToComputationParam(lineBuffer []string) <-chan computationParameter {
	parameterEmitter := make(chan computationParameter)

	go func() {
		regex := regexp.MustCompile(`\s+`)
		timeLineSubStrs := regex.Split(lineBuffer[0], -1)
		distanceSubStrs := regex.Split(lineBuffer[1], -1)

		for i := 1; i < len(timeLineSubStrs); i++ {
			parameterEmitter <- computationParameter{
				time:     util.StrToNum(timeLineSubStrs[i]),
				distance: util.StrToNum(distanceSubStrs[i]),
			}
		}

		close(parameterEmitter)
	}()

	return parameterEmitter
}

func generateTaskFunc(param computationParameter, result *int, lock *sync.Mutex) executor.TaskFunc {
	return func() int {
		delta := param.time*param.time - 4*-1*-param.distance
		sqrtOfDelta := math.Sqrt(float64(delta))
		x1 := int(math.Floor((float64(-param.time)+sqrtOfDelta)/-2 + 1))
		x2 := int(math.Ceil((float64(-param.time)-sqrtOfDelta)/-2 - 1))

		log.Printf("[%d, %d]", x1, x2)

		lock.Lock()
		defer lock.Unlock()

		*result *= x2 - x1 + 1

		return 0
	}
}

func generateTaskEmitter(paramsEmitter <-chan computationParameter, result *int, lock *sync.Mutex) <-chan executor.TaskFunc {
	taskEmitter := make(chan executor.TaskFunc)

	go func() {
		for param := range paramsEmitter {
			taskEmitter <- generateTaskFunc(param, result, lock)
		}

		close(taskEmitter)
	}()

	return taskEmitter
}

func main() {
	e := executor.New(6)

	lineEmitter := util.ReadFile(filepath.Join("2023", "day6", "input.txt"))
	lineBuffer := bufferLines(lineEmitter)
	paramsEmitter := transformLinesToComputationParam(lineBuffer)

	var (
		result = 1
		lock   sync.Mutex
	)
	taskEmitter := generateTaskEmitter(paramsEmitter, &result, &lock)

	e.Run(taskEmitter)

	log.Println("result is", result)
}
