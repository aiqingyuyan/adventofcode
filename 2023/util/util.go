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

		//log.Println("done reading all lines")
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

func Min(a int, rest ...int) int {
	min := a

	for _, next := range rest {
		if next < min {
			min = next
		}
	}

	return min
}

func Max(a int, rest ...int) int {
	max := a

	for _, next := range rest {
		if next > max {
			max = next
		}
	}

	return max
}

func IdxOfMaxElement(elements []int) int {
	currentMaxIdx := 0
	currentMax := elements[0]

	for id, e := range elements {
		if e > currentMax {
			currentMaxIdx = id
			currentMax = e
		}
	}

	return currentMaxIdx
}

func Transpose(matrix []string) []string {
	var newMatrix = make([]string, len(matrix[0]))

	for col := 0; col < len(matrix[0]); col++ {
		var colVals []byte
		for row := 0; row < len(matrix); row++ {
			colVals = append(colVals, matrix[row][col])
		}
		newMatrix[col] = string(colVals)
	}

	return newMatrix
}
