package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

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

		// process color
		if c == ',' || c == ';' {
			color := string(buffer)

			switch color {
			case "red":
				numRed = max(numRed, numCubes)
			case "green":
				numGreen = max(numGreen, numCubes)
			case "blue":
				numBlue = max(numBlue, numCubes)
			default:
				panic(fmt.Errorf("invalid color str: %s", color))
			}

			buffer = nil
		}
	}

	color := string(buffer)
	switch color {
	case "red":
		numRed = max(numRed, numCubes)
	case "green":
		numGreen = max(numGreen, numCubes)
	case "blue":
		numBlue = max(numBlue, numCubes)
	default:
		panic(fmt.Errorf("invalid color str: %s", color))
	}

	return numRed * numGreen * numBlue
}

func main() {
	lineEmitter := readFile()

	result := 0
	for line := range lineEmitter {
		result += processLine(line)
	}

	log.Printf("result is: %d", result)
}
