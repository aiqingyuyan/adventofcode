package main

import (
	"log"
	"path/filepath"
	"strings"
	"yanyu/aoc/2023/util"
)

func hash(s string) int32 {
	var current int32
	for _, c := range s {
		current += c
		current *= 17
		current %= 256
	}

	log.Println(s, "becomes", current)

	return current
}

func processSequences(sequences []string) int32 {
	var result int32
	for _, s := range sequences {
		result += hash(s)
	}

	return result
}

func processLines(lineEmitter <-chan *string) int32 {
	line := <-lineEmitter
	sequences := strings.Split(*line, ",")

	return processSequences(sequences)
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day15", "input.txt"))
	result := processLines(lineEmitter)

	log.Println("result is", result)
}
