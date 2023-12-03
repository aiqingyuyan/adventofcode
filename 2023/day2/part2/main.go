package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"yanyu/aoc/2023/executor"
	"yanyu/aoc/2023/util"
)

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func updateColorCubeNumber(buffer []byte, numRed *int, numGreen *int, numBlue *int, numCubes int) {
	color := string(buffer)

	switch color {
	case "red":
		*numRed = max(*numRed, numCubes)
	case "green":
		*numGreen = max(*numGreen, numCubes)
	case "blue":
		*numBlue = max(*numBlue, numCubes)
	default:
		panic(fmt.Errorf("invalid color str: %s", color))
	}
}

func processLine(line *string) int {
	var passGame bool
	var numCubes int
	var buffer []byte
	var numRed int
	var numGreen int
	var numBlue int

	for _, c := range []byte(*line) {
		if c != ':' && c != ' ' && c != ';' && c != ',' {
			buffer = append(buffer, c)
		}

		if c == ':' {
			passGame = true
			buffer = nil
		}

		// process cube number
		if passGame && c == ' ' && len(buffer) > 0 {
			number, err := strconv.Atoi(string(buffer))
			if err != nil {
				panic(err)
			}
			numCubes = number
			buffer = nil
		}

		if c == ',' || c == ';' {
			// process color
			updateColorCubeNumber(buffer, &numRed, &numGreen, &numBlue, numCubes)
			buffer = nil
		}
	}

	updateColorCubeNumber(buffer, &numRed, &numGreen, &numBlue, numCubes)

	return numRed * numGreen * numBlue
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day2", "part2", "input.txt"))

	result := executor.Run(6, lineEmitter, processLine)

	log.Printf("result is: %d", result)
}
