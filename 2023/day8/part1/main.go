package main

import (
	"log"
	"path/filepath"
	"regexp"
	"yanyu/aoc/2023/util"
)

type nodeNeighbours struct {
	left  string
	right string
}

func processLines(lineEmitter <-chan *string) (instructions []byte, nodeMaps map[string]nodeNeighbours) {
	nodeValuesRegex := regexp.MustCompile(`([A-Z]{3})[^A-Z]+([A-Z]{3})[^A-Z]+([A-Z]{3})\)$`)
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
		}
	}

	return
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day8", "part1", "input.txt"))
	instructions, nodeMaps := processLines(lineEmitter)

	var (
		currentNode       = "AAA"
		currentNeighbours nodeNeighbours
		step              int
		found             bool
	)

	for !found {
		for _, instruction := range instructions {
			currentNeighbours = nodeMaps[currentNode]

			if instruction == 'L' {
				currentNode = currentNeighbours.left
			} else {
				currentNode = currentNeighbours.right
			}

			step++

			if currentNode == "ZZZ" {
				found = true
				break
			}
		}
	}

	log.Println("steps", step)
}
