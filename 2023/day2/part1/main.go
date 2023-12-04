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

func generateTaskFunc(linePtr *string) executor.TaskFunc {
	return func() int {
		return processLine(linePtr)
	}
}

func createTaskEmitter(lineEmitter <-chan *string) <-chan executor.TaskFunc {
	taskEmitter := make(chan executor.TaskFunc)
	go func() {
		for line := range lineEmitter {
			task := generateTaskFunc(line)
			taskEmitter <- task
		}
		close(taskEmitter)
	}()
	return taskEmitter
}

func main() {
	e := executor.New(6)

	lineEmitter := util.ReadFile(filepath.Join("2023", "day2", "part1", "input.txt"))
	taskEmitter := createTaskEmitter(lineEmitter)
	result := e.Run(taskEmitter)

	log.Printf("result is: %d", result)
}
