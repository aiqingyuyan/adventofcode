package main

import (
	"log"
	"math"
	"path/filepath"
	"regexp"
	"strings"
	"yanyu/aoc/2023/util"
)

type computationParameter struct {
	time     int
	distance int
}

func bufferLines(lineEmitter <-chan *string) []string {
	var linesBuffer []string
	for line := range lineEmitter {
		linesBuffer = append(linesBuffer, *line)
	}

	return linesBuffer
}

func transformLinesToComputationParam(lineBuffer []string) computationParameter {
	regex := regexp.MustCompile(`\s+`)
	timeLineSubStrs := regex.Split(lineBuffer[0], -1)
	distanceSubStrs := regex.Split(lineBuffer[1], -1)

	var (
		time     []string
		distance []string
	)
	for i := 1; i < len(timeLineSubStrs); i++ {
		time = append(time, timeLineSubStrs[i])
		distance = append(distance, distanceSubStrs[i])
	}

	return computationParameter{
		time:     util.StrToNum(strings.Join(time, "")),
		distance: util.StrToNum(strings.Join(distance, "")),
	}
}

func computeResult(param computationParameter) int64 {
	delta := int64(param.time*param.time) - int64(4*-1*-param.distance)
	sqrtOfDelta := math.Sqrt(float64(delta))
	x1 := int64(math.Floor((float64(-param.time)+sqrtOfDelta)/-2 + 1))
	x2 := int64(math.Ceil((float64(-param.time)-sqrtOfDelta)/-2 - 1))

	log.Printf("[%d, %d]", x1, x2)

	result := x2 - x1 + 1
	return result
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day6", "input.txt"))
	lineBuffer := bufferLines(lineEmitter)
	param := transformLinesToComputationParam(lineBuffer)
	result := computeResult(param)

	log.Println("result is", result)
}
