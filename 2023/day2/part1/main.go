package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func readFile() <-chan *string {
	lineEmitter := make(chan *string, 15)

	go func() {
		file, err := os.Open(filepath.Join("2023", "day2", "part1", "input.txt"))
		if err != nil {
			panic(err)
		}

		fileScanner := bufio.NewScanner(file)

		fileScanner.Split(bufio.ScanLines)

		for fileScanner.Scan() {
			line := fileScanner.Text()
			lineEmitter <- &line
		}

		close(lineEmitter)

		log.Println("done reading all lines")
	}()

	return lineEmitter
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
			color := string(buffer)

			switch color {
			case "red":
				if numCubes > 12 {
					return 0
				}
			case "green":
				if numCubes > 13 {
					return 0
				}
			case "blue":
				if numCubes > 14 {
					return 0
				}
			default:
				panic(fmt.Errorf("invalid color str: %s", color))
			}

			buffer = nil
		}
	}

	color := string(buffer)
	switch color {
	case "red":
		if numCubes > 12 {
			return 0
		}
	case "green":
		if numCubes > 13 {
			return 0
		}
	case "blue":
		if numCubes > 14 {
			return 0
		}
	default:
		panic(fmt.Errorf("invalid color str: %s", color))
	}

	return gameId
}

func main() {
	lineEmitter := readFile()

	result := 0
	for line := range lineEmitter {
		result += processLine(line)
	}

	log.Printf("result is: %d", result)
}
