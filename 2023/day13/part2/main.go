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

func compareDiff(str1, str2 string) int {
	countOfDiff := 0

	for i, c := range []byte(str1) {
		if c != str2[i] {
			countOfDiff++
		}
	}

	return countOfDiff
}

func isValidReflection(matrix []string, rowNum int) bool {
	countOfDiff := 0

	for i, j := rowNum, rowNum+1; i >= 0 && j < len(matrix); {
		diff := compareDiff(matrix[i], matrix[j])
		if diff == 1 {
			countOfDiff++
		}

		if diff > 2 {
			return false
		}

		i--
		j++
	}

	if countOfDiff == 1 {
		return true
	} else {
		return false
	}
}

func findRowNumber(matrix []string, transposed bool) int {
	for i := 0; i < len(matrix); i++ {
		if isValidReflection(matrix, i) {
			if !transposed {
				return 100 * (i + 1)
			} else {
				return i + 1
			}
		}
	}

	return 0
}

func processPattern(pattern *Pattern) int {
	result := findRowNumber(pattern.matrix, false)

	if result == 0 {
		pattern.matrix = transpose(pattern.matrix)
		result = findRowNumber(pattern.matrix, true)
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
