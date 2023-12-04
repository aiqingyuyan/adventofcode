package main

import (
	"log"
	"path/filepath"
	"yanyu/aoc/2023/executor"
	"yanyu/aoc/2023/util"
)

func updateResult(result, countOfFoundWinningNumber int) int {
	if countOfFoundWinningNumber == 1 {
		return 1
	} else {
		return result * 2
	}
}

func processLine(line *string) int {
	var (
		seenColon                 bool
		seenBar                   bool
		countOfFoundWinningNumber int
		result                    int
		numBuffer                 []byte
		winningNumsToCheck        = make(map[string]bool)
	)

	for _, c := range []byte(*line) {
		if c == ':' && !seenColon {
			seenColon = true
			continue
		}

		if c == '|' && !seenBar {
			seenBar = true
			continue
		}

		if seenColon {
			if util.IsByteANumber(c) {
				numBuffer = append(numBuffer, c)
			}

			if c == ' ' && len(numBuffer) > 0 {
				if !seenBar {
					winningNumsToCheck[string(numBuffer)] = true
				} else if winningNumsToCheck[string(numBuffer)] {
					countOfFoundWinningNumber++
					result = updateResult(result, countOfFoundWinningNumber)
				}
				numBuffer = nil
			}
		}
	}

	if winningNumsToCheck[string(numBuffer)] {
		countOfFoundWinningNumber++
		result = updateResult(result, countOfFoundWinningNumber)
	}

	return result
}

func generateTaskFunc(linePtr *string) executor.TaskFunc {
	return func() int {
		return processLine(linePtr)
	}
}

func createTaskEmitter(lineEmitter <-chan *string) <-chan executor.TaskFunc {
	taskEmitter := make(chan executor.TaskFunc)
	go func() {
		for line := range lineEmitter {
			taskEmitter <- generateTaskFunc(line)
		}
		close(taskEmitter)
	}()
	return taskEmitter
}

func main() {
	e := executor.New(6)

	lineEmitter := util.ReadFile(filepath.Join("2023", "day4", "input.txt"))
	taskEmitter := createTaskEmitter(lineEmitter)

	result := e.Run(taskEmitter)

	log.Println("result is ", result)
}
