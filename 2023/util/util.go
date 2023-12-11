package util

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func ReadFile(path string) <-chan *string {
	lineEmitter := make(chan *string, 15)

	go func() {
		file, err := os.Open(path)
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

func IsByteANumber(c byte) bool {
	if c >= '0' && c <= '9' {
		return true
	}

	return false
}

func StrToNum(str string) int {
	if num, err := strconv.Atoi(str); err == nil {
		return num
	} else {
		log.Printf("err: %+v", err)
		panic(err)
	}
}

func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
