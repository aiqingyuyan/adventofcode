package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"yanyu/aoc/2023/executor"
	"yanyu/aoc/2023/util"
)

func isColorCubeNumberValid(buffer []byte, numCubes int) bool {
	color := string(buffer)

	switch color {
	case "red":
		if numCubes > 12 {
			return false
		}
	case "green":
		if numCubes > 13 {
			return false
		}
	case "blue":
		if numCubes > 14 {
			return false
		}
	default:
		panic(fmt.Errorf("invalid color str: %s", color))
	}

	return true
}

func processLine(line *string) int {
	var seenGame bool
	var gameId int
	var numCubes int
	var buffer []byte

	for _, c := range []byte(*line) {
		if c != ':' && c != ' ' && c != ';' && c != ',' {
			buffer = append(buffer, c)
		}

		if !seenGame && len(buffer) == 4 && string(buffer) == "Game" {
			seenGame = true
			buffer = nil
		}

		// process game id
		if c == ':' {
			number, err := strconv.Atoi(string(buffer))
			if err != nil {
				panic(err)
			}
			gameId = number
			buffer = nil
		}

		// process cube number
		if c == ' ' && len(buffer) > 0 {
			number, err := strconv.Atoi(string(buffer))
			if err != nil {
				panic(err)
			}
			numCubes = number
			buffer = nil
		}

		// process color
		if c == ',' || c == ';' {
			if !isColorCubeNumberValid(buffer, numCubes) {
				return 0
			}

			buffer = nil
		}
	}

	if !isColorCubeNumberValid(buffer, numCubes) {
		return 0
	}

	return gameId
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day2", "part1", "input.txt"))

	result := executor.Run(6, lineEmitter, processLine)

	log.Printf("result is: %d", result)
}
