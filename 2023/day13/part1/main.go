package main

import (
	"log"
	"path/filepath"
	"yanyu/aoc/2023/util"
)

type Pattern struct {
	matrix []string
}

func transpose(matrix []string) []string {
	var newMatrix = make([]string, len(matrix[0]))

	for col := 0; col < len(matrix[0]); col++ {
		var colVals []byte
		for row := 0; row < len(matrix); row++ {
			colVals = append(colVals, matrix[row][col])
		}
		newMatrix[col] = string(colVals)
	}

	return newMatrix
}

func processLine(lineEmitter <-chan *string) []Pattern {
	var (
		matrix   []string
		patterns []Pattern
	)

	for line := range lineEmitter {
		if *line != "" {
			matrix = append(matrix, *line)
		} else {
			patterns = append(patterns, Pattern{matrix: matrix})
			matrix = nil
		}
	}

	patterns = append(patterns, Pattern{matrix: matrix})

	return patterns
}

func findRowNumber(matrix []string) []int {
	//for _, line := range matrix {
	//	log.Println(line)
	//}
	//
	//log.Println("----------------")

	var rowNums []int
	for i, j := 0, 1; j < len(matrix); {
		if matrix[i] == matrix[j] {
			rowNums = append(rowNums, i)
		}

		i++
		j++
	}

	return rowNums
}

func isValidReflection(matrix []string, rowNum int) bool {
	var isValid = true
	for i, j := rowNum, rowNum+1; i >= 0 && j < len(matrix); {
		if matrix[i] != matrix[j] {
			isValid = false
			break
		}

		i--
		j++
	}

	return isValid
}

func processPattern(pattern *Pattern) int {
	result := 0

	for _, rowNum := range findRowNumber(pattern.matrix) {
		if isValidReflection(pattern.matrix, rowNum) {
			result += 100 * (rowNum + 1)
			break
		}
	}

	if result == 0 {
		pattern.matrix = transpose(pattern.matrix)
		for _, colNum := range findRowNumber(pattern.matrix) {
			if isValidReflection(pattern.matrix, colNum) {
				result += colNum + 1
				break
			}
		}
	}

	return result
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day13", "input.txt"))
	patterns := processLine(lineEmitter)

	result := 0
	for _, pattern := range patterns {
		lineResult := processPattern(&pattern)
		result += lineResult
	}
	log.Println("result is", result)
}
