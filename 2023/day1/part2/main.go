package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

func readFile() <-chan *string {
	lineEmitter := make(chan *string, 15)

	go func() {
		file, err := os.Open(filepath.Join("2023", "day1", "part2", "input.txt"))
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

func startWorker(schedulerChan chan chan *string, resultChan chan int) {
	workerTaskChan := make(chan *string)

	go func() {
		for {
			schedulerChan <- workerTaskChan

			line := <-workerTaskChan

			resultChan <- processLine(line)
		}
	}()
}

func start(numOfWorker int) int {
	if numOfWorker == 0 {
		panic(fmt.Errorf("invalid worker number: %d, must be > 0", numOfWorker))
	}

	result := 0

	lineEmitter := readFile()

	resultChan := make(chan int, numOfWorker)
	schedulerChan := make(chan chan *string)
	for i := 0; i < numOfWorker; i++ {
		startWorker(schedulerChan, resultChan)
	}

	var doneReadingFile bool
	var workerQueue []chan *string
	var taskQueue []*string
	for {
		if doneReadingFile && len(taskQueue) == 0 && len(workerQueue) == numOfWorker && len(resultChan) == 0 {
			break
		}

		var activeTask *string
		var activeWorker chan *string

		if len(taskQueue) > 0 && len(workerQueue) > 0 {
			activeTask = taskQueue[0]
			activeWorker = workerQueue[0]
		}

		select {
		case line, ok := <-lineEmitter:
			if ok {
				taskQueue = append(taskQueue, line)
			} else {
				doneReadingFile = true
			}
		case workerChan := <-schedulerChan:
			workerQueue = append(workerQueue, workerChan)
		case activeWorker <- activeTask:
			taskQueue = taskQueue[1:]
			workerQueue = workerQueue[1:]
		case number := <-resultChan:
			result += number
		}
	}

	return result
}

func main() {
	result := start(6)

	log.Printf("result is %d", result)
}
