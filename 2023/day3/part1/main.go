package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"yanyu/aoc/2023/executor"
	"yanyu/aoc/2023/util"
)

func checkUpperLeft(matrix []string, rowIdx, columnIdx int) bool {
	if rowIdx < 0 || columnIdx < 0 {
		return false
	}

	return matrix[rowIdx][columnIdx] != '.' && (matrix[rowIdx][columnIdx] < '0' || matrix[rowIdx][columnIdx] > '9')
}

func checkUpper(matrix []string, rowIdx, columnIdx int) bool {
	if rowIdx < 0 {
		return false
	}

	return matrix[rowIdx][columnIdx] != '.' && (matrix[rowIdx][columnIdx] < '0' || matrix[rowIdx][columnIdx] > '9')
}

func checkUpperRight(matrix []string, rowIdx, columnIdx int) bool {
	if rowIdx < 0 || columnIdx > len(matrix[rowIdx])-1 {
		return false
	}

	return matrix[rowIdx][columnIdx] != '.' && (matrix[rowIdx][columnIdx] < '0' || matrix[rowIdx][columnIdx] > '9')
}

func checkLeft(matrix []string, rowIdx, columnIdx int) bool {
	if columnIdx < 0 {
		return false
	}

	return matrix[rowIdx][columnIdx] != '.' && (matrix[rowIdx][columnIdx] < '0' || matrix[rowIdx][columnIdx] > '9')
}

func checkRight(matrix []string, rowIdx, columnIdx int) bool {
	if columnIdx > len(matrix[rowIdx])-1 {
		return false
	}

	return matrix[rowIdx][columnIdx] != '.' && (matrix[rowIdx][columnIdx] < '0' || matrix[rowIdx][columnIdx] > '9')
}

func checkLowerLeft(matrix []string, rowIdx, columnIdx int) bool {
	if rowIdx > len(matrix)-1 || columnIdx < 0 {
		return false
	}

	return matrix[rowIdx][columnIdx] != '.' && (matrix[rowIdx][columnIdx] < '0' || matrix[rowIdx][columnIdx] > '9')
}

func checkLower(matrix []string, rowIdx, columnIdx int) bool {
	if rowIdx > len(matrix)-1 {
		return false
	}

	return matrix[rowIdx][columnIdx] != '.' && (matrix[rowIdx][columnIdx] < '0' || matrix[rowIdx][columnIdx] > '9')
}

func checkLowerRight(matrix []string, rowIdx, columnIdx int) bool {
	if rowIdx > len(matrix)-1 || columnIdx > len(matrix[rowIdx])-1 {
		return false
	}

	return matrix[rowIdx][columnIdx] != '.' && (matrix[rowIdx][columnIdx] < '0' || matrix[rowIdx][columnIdx] > '9')
}

func isNumberAdjacentToSymbol(matrix []string, rowIdx, columnIdx int) bool {
	return checkUpperLeft(matrix, rowIdx-1, columnIdx-1) || checkUpper(matrix, rowIdx-1, columnIdx) ||
		checkUpperRight(matrix, rowIdx-1, columnIdx+1) || checkLeft(matrix, rowIdx, columnIdx-1) ||
		checkRight(matrix, rowIdx, columnIdx+1) || checkLowerLeft(matrix, rowIdx+1, columnIdx-1) ||
		checkLower(matrix, rowIdx+1, columnIdx) || checkLowerRight(matrix, rowIdx+1, columnIdx+1)
}

func isValidNumber(matrix []string, rowIdx, numStartIdx, numEndIdx int) bool {
	return isNumberAdjacentToSymbol(matrix, rowIdx, numStartIdx) || isNumberAdjacentToSymbol(matrix, rowIdx, numEndIdx)
}

func processLine(matrix []string, currentLineIdx int) int {
	currentLine := matrix[currentLineIdx]

	var (
		numBuffer   []byte
		numStartIdx = -1
		numEndIdx   = -1
		result      = 0
	)

	for columnIdx, c := range []byte(currentLine) {
		if c >= '0' && c <= '9' {
			numBuffer = append(numBuffer, c)
			if numStartIdx < 0 {
				numStartIdx = columnIdx
			}
		}

		if (c < '0' || c > '9') && numStartIdx >= 0 {
			numEndIdx = columnIdx - 1
			if isValidNumber(matrix, currentLineIdx, numStartIdx, numEndIdx) {
				num, err := strconv.Atoi(string(numBuffer))
				if err != nil {
					fmt.Printf("%+v\n", err)
					panic(err)
				}
				result += num
			}
			numBuffer = nil
			numStartIdx = -1
			numEndIdx = -1
		}
	}

	if len(numBuffer) > 0 {
		numEndIdx = len(currentLine) - 1
		if isValidNumber(matrix, currentLineIdx, numStartIdx, numEndIdx) {
			num, err := strconv.Atoi(string(numBuffer))
			if err != nil {
				fmt.Printf("%+v\n", err)
				panic(err)
			}
			result += num
		}
	}

	return result
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
	lineEmitter := util.ReadFile(filepath.Join("2023", "day3", "part1", "input.txt"))

	var matrix []string
	for line := range lineEmitter {
		matrix = append(matrix, *line)
	}

	e := executor.New(6)

	taskEmitter := createTaskEmitter(matrix)
	result := e.Run(taskEmitter)

	log.Println("result is ", result)
}
