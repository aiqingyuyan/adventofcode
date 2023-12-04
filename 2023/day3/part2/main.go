package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"yanyu/aoc/2023/executor"
	"yanyu/aoc/2023/util"
)

func getLeft(matrix []string, rowIdx, columnIdx int) []byte {
	if columnIdx < 0 {
		return nil
	}

	if matrix[rowIdx][columnIdx] >= '0' && matrix[rowIdx][columnIdx] <= '9' {
		startIdx := columnIdx
		countOfNum := 1
		for columnIdx > 0 {
			columnIdx--

			if matrix[rowIdx][columnIdx] >= '0' && matrix[rowIdx][columnIdx] <= '9' {
				startIdx = columnIdx
				countOfNum++
			} else {
				break
			}
		}

		return []byte(matrix[rowIdx][startIdx : startIdx+countOfNum])
	}

	return nil
}

func getRight(matrix []string, rowIdx, columnIdx int) []byte {
	if columnIdx > len(matrix[rowIdx])-1 {
		return nil
	}

	if matrix[rowIdx][columnIdx] >= '0' && matrix[rowIdx][columnIdx] <= '9' {
		startIdx := columnIdx
		countOfNum := 1
		for columnIdx < len(matrix[rowIdx])-1 {
			columnIdx++

			if matrix[rowIdx][columnIdx] >= '0' && matrix[rowIdx][columnIdx] <= '9' {
				countOfNum++
			} else {
				break
			}
		}

		return []byte(matrix[rowIdx][startIdx : startIdx+countOfNum])
	}

	return nil
}

func getNumbers(matrix []string, rowIdx, columnIdx int) []int {
	if rowIdx < 0 {
		return nil
	}

	if rowIdx > len(matrix)-1 {
		return nil
	}

	var result []int
	if matrix[rowIdx][columnIdx] < '0' || matrix[rowIdx][columnIdx] > '9' {
		numBufferLeft := getLeft(matrix, rowIdx, columnIdx-1)
		numBufferRight := getRight(matrix, rowIdx, columnIdx+1)

		if numBufferLeft != nil {
			if num, err := strconv.Atoi(string(numBufferLeft)); err == nil {
				result = append(result, num)
			} else {
				fmt.Printf("%+v\n", err)
				panic(err)
			}
		}

		if numBufferRight != nil {
			if num, err := strconv.Atoi(string(numBufferRight)); err == nil {
				result = append(result, num)
			} else {
				fmt.Printf("%+v\n", err)
				panic(err)
			}
		}
	} else {
		numBufferLeft := getLeft(matrix, rowIdx, columnIdx-1)
		numBufferRight := getRight(matrix, rowIdx, columnIdx+1)

		numBuffer := append(numBufferLeft, matrix[rowIdx][columnIdx])
		numBuffer = append(numBuffer, numBufferRight...)

		if num, err := strconv.Atoi(string(numBuffer)); err == nil {
			result = append(result, num)
		} else {
			fmt.Printf("%+v\n", err)
			panic(err)
		}
	}

	return result
}

func isGear(matrix []string, rowIdx, columnIdx int) (gearRatio int, ok bool) {
	upperResult := getNumbers(matrix, rowIdx-1, columnIdx)
	currentLineResult := getNumbers(matrix, rowIdx, columnIdx)
	lowerResult := getNumbers(matrix, rowIdx+1, columnIdx)

	countOfAdjacentNumber := len(upperResult) + len(currentLineResult) + len(lowerResult)

	switch {
	case countOfAdjacentNumber > 2 || countOfAdjacentNumber < 2:
		return -1, false
	default:
		allNumbers := append(upperResult, currentLineResult...)
		allNumbers = append(allNumbers, lowerResult...)

		log.Printf("at row %d, col: %d, found gear, numbers: %d, %d", rowIdx, columnIdx, allNumbers[0], allNumbers[1])

		return allNumbers[0] * allNumbers[1], true
	}
}

func processLine(matrix []string, currentLineIdx int) int {
	currentLine := matrix[currentLineIdx]

	lineResult := 0
	for columnIdx, c := range []byte(currentLine) {
		if c == '*' {
			if gearRatio, ok := isGear(matrix, currentLineIdx, columnIdx); ok {
				lineResult += gearRatio
			}
		}
	}

	return lineResult
}

func generateTaskFunc(matrix []string, lineIdx int) executor.TaskFunc {
	return func() int {
		return processLine(matrix, lineIdx)
	}
}

func createTaskEmitter(matrix []string) <-chan executor.TaskFunc {
	taskEmitter := make(chan executor.TaskFunc)
	go func() {
		for lineIdx := 0; lineIdx < len(matrix); lineIdx++ {
			taskEmitter <- generateTaskFunc(matrix, lineIdx)
		}
		close(taskEmitter)
	}()
	return taskEmitter
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day3", "part2", "input.txt"))

	var matrix []string
	for line := range lineEmitter {
		matrix = append(matrix, *line)
	}

	e := executor.New(6)

	taskEmitter := createTaskEmitter(matrix)
	result := e.Run(taskEmitter)

	log.Println("result is ", result)
}
