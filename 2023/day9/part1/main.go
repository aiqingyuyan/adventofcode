package main

import (
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"yanyu/aoc/2023/executor"
	"yanyu/aoc/2023/util"
)

func transformToNums(numStrs []string) []int {
	var nums []int
	for _, numStr := range numStrs {
		nums = append(nums, util.StrToNum(numStr))
	}

	return nums
}

func getNextLine(nums []int) ([]int, bool) {
	var (
		length  = len(nums)
		allZero = true
	)

	for i, j := 0, 1; j < length; {
		nums[i] = nums[j] - nums[i]

		if nums[i] != 0 {
			allZero = false
		}

		i++
		j++
	}

	return nums[0 : length-1], allZero
}

func predictNextNum(nums []int) int {
	var (
		lastValueOfEachLine = []int{nums[len(nums)-1]}
		allZeros            bool
	)

	for !allZeros {
		nums, allZeros = getNextLine(nums)
		lastValueOfEachLine = append(lastValueOfEachLine, nums[len(nums)-1])
	}

	result := 0
	for _, val := range lastValueOfEachLine {
		result += val
	}

	return result
}

func processLines(line *string) int {
	nums := transformToNums(strings.Split(*line, " "))
	return predictNextNum(nums)
}

func generateTaskFunc(line *string) executor.TaskFunc {
	return func() any {
		return processLines(line)
	}
}

func generateTaskEmitter(lineEmitter <-chan *string) <-chan executor.TaskFunc {
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
	e := executor.New(runtime.NumCPU())

	lineEmitter := util.ReadFile(filepath.Join("2023", "day9", "input.txt"))
	taskEmitter := generateTaskEmitter(lineEmitter)

	var (
		result           = 0
		resultHandleFunc = func(taskResult any) {
			result += taskResult.(int)
		}
	)

	e.Run(taskEmitter, resultHandleFunc)

	log.Println("result is", result)
}
