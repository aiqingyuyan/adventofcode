package main

import (
	"log"
	"path/filepath"
	"strconv"
	"yanyu/aoc/2023/executor"
	"yanyu/aoc/2023/util"
)

var digitWordMap = map[string]byte{
	"one":   '1',
	"two":   '2',
	"three": '3',
	"four":  '4',
	"five":  '5',
	"six":   '6',
	"seven": '7',
	"eight": '8',
	"nine":  '9',
}

func processLine(line *string) int {
	numSlice := make([]byte, 2)
	var digit3LetterSlice []byte
	var digit4LetterSlice []byte
	var digit5LetterSlice []byte
	for _, c := range *line {
		if c >= '0' && c <= '9' {
			if numSlice[0] == 0 {
				numSlice[0] = byte(c)
			}

			numSlice[1] = byte(c)
		} else if c >= 'a' && c <= 'z' {
			digit3LetterSlice = append(digit3LetterSlice, byte(c))
			digit4LetterSlice = append(digit4LetterSlice, byte(c))
			digit5LetterSlice = append(digit5LetterSlice, byte(c))

			if len(digit3LetterSlice) == 3 {
				digitWord := string(digit3LetterSlice)
				if v, ok := digitWordMap[digitWord]; ok {
					if numSlice[0] == 0 {
						numSlice[0] = v
					}

					numSlice[1] = v
				}
				digit3LetterSlice = digit3LetterSlice[1:]
			}

			if len(digit4LetterSlice) == 4 {
				digitWord := string(digit4LetterSlice)
				if v, ok := digitWordMap[digitWord]; ok {
					if numSlice[0] == 0 {
						numSlice[0] = v
					}

					numSlice[1] = v
				}
				digit4LetterSlice = digit4LetterSlice[1:]
			}

			if len(digit5LetterSlice) == 5 {
				digitWord := string(digit5LetterSlice)
				if v, ok := digitWordMap[digitWord]; ok {
					if numSlice[0] == 0 {
						numSlice[0] = v
					}

					numSlice[1] = v
				}
				digit5LetterSlice = digit5LetterSlice[1:]
			}
		}
	}

	// line doesn't include number char or word
	if numSlice[0] == 0 {
		return 0
	}

	num, err := strconv.Atoi(string(numSlice))
	if err != nil {
		panic(err)
	}

	return num
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

	lineEmitter := util.ReadFile(filepath.Join("2023", "day1", "part2", "input.txt"))
	taskEmitter := createTaskEmitter(lineEmitter)
	result := e.Run(taskEmitter)

	log.Printf("result is %d", result)
}
