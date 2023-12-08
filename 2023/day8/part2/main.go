package main

import (
	"log"
	"path/filepath"
	"regexp"
	"yanyu/aoc/2023/executor"
	"yanyu/aoc/2023/util"
)

type nodeNeighbours struct {
	left  string
	right string
}

func processLines(lineEmitter <-chan *string) (instructions []byte, startNodes []string, nodeMaps map[string]nodeNeighbours) {
	nodeValuesRegex := regexp.MustCompile(`([A-Z0-9]{3})[^A-Z0-9]+([A-Z0-9]{3})[^A-Z0-9]+([A-Z0-9]{3})\)$`)
	nodeMaps = make(map[string]nodeNeighbours)

	for line := range lineEmitter {
		if len(instructions) == 0 {
			instructions = []byte(*line)
			continue
		}

		if len(*line) != 0 && len(instructions) > 0 {
			matches := nodeValuesRegex.FindAllStringSubmatch(*line, -1)

			nodeValue := matches[0][1]
			left := matches[0][2]
			right := matches[0][3]

			nodeMaps[nodeValue] = nodeNeighbours{
				left:  left,
				right: right,
			}

			if nodeValue[2] == 'A' {
				startNodes = append(startNodes, nodeValue)
			}
		}
	}

	return
}

func getNumOfStepsToReachEndFirstTime(instructions []byte, nodeMaps map[string]nodeNeighbours, currentNode string) int {
	var (
		step  int
		found bool
	)

	for !found {
		for _, instruction := range instructions {
			currentNeighbours := nodeMaps[currentNode]

			if instruction == 'L' {
				currentNode = currentNeighbours.left
			} else {
				currentNode = currentNeighbours.right
			}

			step++

			if currentNode[2] == 'Z' {
				found = true
				break
			}
		}
	}

	return step
}

func generateTaskFunc(instructions []byte, nodeMaps map[string]nodeNeighbours, startNode string) executor.TaskFunc {
	return func() any {
		steps := getNumOfStepsToReachEndFirstTime(instructions, nodeMaps, startNode)

		log.Printf("node %s, need %d step to reach Z for the first time", startNode, steps)

		return steps
	}
}

func generateTaskEmitter(instructions []byte, nodeMaps map[string]nodeNeighbours, startNodes []string) <-chan executor.TaskFunc {
	taskEmitter := make(chan executor.TaskFunc)

	go func() {
		for _, startNode := range startNodes {
			taskEmitter <- generateTaskFunc(instructions, nodeMaps, startNode)
		}

		close(taskEmitter)
	}()

	return taskEmitter
}

// greatest common divisor (GCD) - Euclidean algorithm
func gcd(a, b int) int {
	for b != 0 {
		temp := b
		b = a % b
		a = temp
	}

	return a
}

// Least Common Multiple (LCM)
func lcm(a, b int, rest ...int) int {
	currentLcm := a * b / gcd(a, b)

	for _, num := range rest {
		currentLcm = currentLcm * num / gcd(currentLcm, num)
	}

	return currentLcm
}

func main() {
	e := executor.New(6)

	lineEmitter := util.ReadFile(filepath.Join("2023", "day8", "part2", "input.txt"))
	instructions, startNodes, nodeMaps := processLines(lineEmitter)

	taskEmitter := generateTaskEmitter(instructions, nodeMaps, startNodes)

	var (
		steps            []int
		resultHandleFunc = func(result any) {
			steps = append(steps, result.(int))
		}
	)

	e.Run(taskEmitter, resultHandleFunc)

	result := lcm(steps[0], steps[1], steps[2:]...)

	log.Println("result is", result)
}
