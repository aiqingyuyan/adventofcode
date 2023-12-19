package main

import (
	"log"
	"path/filepath"
	"yanyu/aoc/2023/util"
)

func processLines(lineEmitter <-chan *string) []string {
	var matrix []string
	for line := range lineEmitter {
		matrix = append(matrix, *line)
	}
	return util.Transpose(matrix)
}

func fillRound(row []byte, prevStopIdx, currStopIdx, countOfRound int) {
	for i := prevStopIdx; i < currStopIdx; i++ {
		if countOfRound > 0 && row[i] != '#' { // in case prevStopIdx is '#'
			row[i] = 'O'
			countOfRound--
		} else if row[i] != '#' { // in case prevStopIdx is '#'
			row[i] = '.'
		}
	}
}

func moveRoundRocks(row []byte) {
	var (
		prevStopIdx  int
		currStopIdx  int
		countOfRound int
	)

	for i := 0; i < len(row); i++ {
		if row[i] == '#' {
			currStopIdx = i
			fillRound(row, prevStopIdx, currStopIdx, countOfRound)
			prevStopIdx = currStopIdx
			countOfRound = 0
		}

		if row[i] == 'O' {
			countOfRound++
		}
	}

	if countOfRound > 0 {
		fillRound(row, prevStopIdx, len(row), countOfRound)
	}
}

func calculateRowLoad(row string) int {
	result := 0
	for i, c := range row {
		if c == 'O' {
			result += len(row) - i
		}
	}
	return result
}

func calculateLoad(matrix []string) int {
	result := 0
	for _, row := range matrix {
		result += calculateRowLoad(row)
	}
	return result
}

func tilt(matrix []string) {
	tempMatrix := make([][]byte, len(matrix))

	for i := 0; i < len(matrix); i++ {
		tempMatrix[i] = []byte(matrix[i])
	}

	for _, row := range tempMatrix {
		moveRoundRocks(row)
	}

	for i := 0; i < len(matrix); i++ {
		matrix[i] = string(tempMatrix[i])
	}
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day14", "input.txt"))
	matrix := processLines(lineEmitter)
	tilt(matrix)
	load := calculateLoad(matrix)

	log.Println("result is", load)
}
