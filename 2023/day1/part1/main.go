package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func readFile() <-chan *string {
	lineEmitter := make(chan *string, 15)

	go func() {
		file, err := os.Open(filepath.Join("2023", "day1", "part1", "input.txt"))
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
	numSlice := make([]byte, 2)
	for _, c := range *line {
		if c >= '0' && c <= '9' {
			if numSlice[0] == 0 {
				numSlice[0] = byte(c)
			}

			numSlice[1] = byte(c)
		}
	}

	// line doesn't include number char
	if numSlice[0] == 0 {
		return 0
	}

	num, err := strconv.Atoi(string(numSlice))
	if err != nil {
		panic(err)
	}

	return num
}

func main() {
	result := 0
	lineEmitter := readFile()

	for line := range lineEmitter {
		result += processLine(line)
	}

	log.Printf("result is %d", result)
}
