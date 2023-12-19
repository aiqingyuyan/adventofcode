package main

import (
	"log"
	"path/filepath"
	"strings"
	"yanyu/aoc/2023/util"
)

func processLines(lineEmitter <-chan *string) []string {
	var matrix []string
	for line := range lineEmitter {
		matrix = append(matrix, *line)
	}
	return matrix
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

func fillRoundReverse(row []byte, prevStopIdx, currStopIdx, countOfRound int) {
	for i := prevStopIdx; i > currStopIdx; i-- {
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

func moveRoundRocksReverse(row []byte) {
	var (
		prevStopIdx  = len(row) - 1
		currStopIdx  = len(row) - 1
		countOfRound int
	)

	for i := len(row) - 1; i >= 0; i-- {
		if row[i] == '#' {
			currStopIdx = i
			fillRoundReverse(row, prevStopIdx, currStopIdx, countOfRound)
			prevStopIdx = currStopIdx
			countOfRound = 0
		}

		if row[i] == 'O' {
			countOfRound++
		}
	}

	if countOfRound > 0 {
		fillRoundReverse(row, prevStopIdx, -1, countOfRound)
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

func tiltReverse(matrix []string) {
	tempMatrix := make([][]byte, len(matrix))

	for i := 0; i < len(matrix); i++ {
		tempMatrix[i] = []byte(matrix[i])
	}

	for _, row := range tempMatrix {
		moveRoundRocksReverse(row)
	}

	for i := 0; i < len(matrix); i++ {
		matrix[i] = string(tempMatrix[i])
	}
}

func spin(matrix []string) []string {
	// north
	matrix = util.Transpose(matrix)
	tilt(matrix)
	matrix = util.Transpose(matrix)

	// west
	tilt(matrix)

	// south
	matrix = util.Transpose(matrix)
	tiltReverse(matrix)
	matrix = util.Transpose(matrix)

	// east
	tiltReverse(matrix)

	return matrix
}

func toStr(matrix []string) string {
	strBuilder := strings.Builder{}
	for _, row := range matrix {
		strBuilder.Write([]byte(row))
		strBuilder.WriteByte('\n')
	}

	return strBuilder.String()
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day14", "input.txt"))
	matrix := processLines(lineEmitter)

	const numOfCycles = 1000000000
	cache := make(map[string]int)
	for i := 0; i < numOfCycles; {
		matrix = spin(matrix)

		if v, ok := cache[toStr(matrix)]; ok {
			log.Printf("i %d, v %d", i, v)

			offset := i - v
			numOfOffsetToSkip := (numOfCycles - i) / offset
			i = i + offset*numOfOffsetToSkip
		}

		cache[toStr(matrix)] = i
		i++
	}

	matrix = util.Transpose(matrix)

	load := calculateLoad(matrix)

	log.Println("result is", load)
}
