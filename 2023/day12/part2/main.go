package main

import (
	"log"
	"path/filepath"
	"strings"
	"yanyu/aoc/2023/util"
)

func toNums(groupsOfDamaged []string) []int {
	var result []int
	for _, numStr := range groupsOfDamaged {
		result = append(result, util.StrToNum(numStr))
	}
	return result
}

type DPFuncParameters struct {
	currentFoundDamaged      int
	indexIntoGroupsOfDamaged int
	subRecords               string
}

// subRecords - for dp memorization & debug
func findCombinations(records []byte, groupsOfDamaged []int, currentFoundDamaged int, indexIntoGroupsOfDamaged int, subRecords string, mem map[DPFuncParameters]int) int {
	if v, ok := mem[DPFuncParameters{currentFoundDamaged, indexIntoGroupsOfDamaged, subRecords}]; ok {
		return v
	}

	if len(records) == 0 {
		if indexIntoGroupsOfDamaged == len(groupsOfDamaged) && currentFoundDamaged == 0 {
			return 1
		}

		if indexIntoGroupsOfDamaged == len(groupsOfDamaged)-1 && groupsOfDamaged[indexIntoGroupsOfDamaged] == currentFoundDamaged {
			return 1
		}

		return 0
	}

	numOfCombinations := 0
	switch records[0] {
	case '#':
		numOfCombinations = findCombinations(records[1:], groupsOfDamaged, currentFoundDamaged+1, indexIntoGroupsOfDamaged, subRecords+"#", mem)
	case '.':
		if indexIntoGroupsOfDamaged < len(groupsOfDamaged) && groupsOfDamaged[indexIntoGroupsOfDamaged] == currentFoundDamaged {
			numOfCombinations = findCombinations(records[1:], groupsOfDamaged, 0, indexIntoGroupsOfDamaged+1, subRecords+".", mem)
		} else if currentFoundDamaged == 0 { // if currentDamagedGroupNum > 0: current damaged group != groupsOfDamaged[0]
			numOfCombinations = findCombinations(records[1:], groupsOfDamaged, currentFoundDamaged, indexIntoGroupsOfDamaged, subRecords+".", mem)
		}
	case '?':
		// .
		if indexIntoGroupsOfDamaged < len(groupsOfDamaged) && groupsOfDamaged[indexIntoGroupsOfDamaged] == currentFoundDamaged {
			numOfCombinations = findCombinations(records[1:], groupsOfDamaged, 0, indexIntoGroupsOfDamaged+1, subRecords+"?", mem)
		} else if currentFoundDamaged == 0 { // if currentDamagedGroupNum > 0: current damaged group != groupsOfDamaged[0]
			numOfCombinations = findCombinations(records[1:], groupsOfDamaged, currentFoundDamaged, indexIntoGroupsOfDamaged, subRecords+"?", mem)
		}

		// #
		numOfCombinations += findCombinations(records[1:], groupsOfDamaged, currentFoundDamaged+1, indexIntoGroupsOfDamaged, subRecords+"?", mem)
	}

	mem[DPFuncParameters{currentFoundDamaged, indexIntoGroupsOfDamaged, subRecords}] = numOfCombinations

	return numOfCombinations
}

func processLine(lineEmitter <-chan *string) {
	result := 0
	for line := range lineEmitter {
		subStrs := strings.Split(*line, " ")
		lineResult := findCombinations(
			[]byte(strings.Join([]string{subStrs[0], subStrs[0], subStrs[0], subStrs[0], subStrs[0]}, "?")),
			toNums(strings.Split(strings.Join([]string{subStrs[1], subStrs[1], subStrs[1], subStrs[1], subStrs[1]}, ","), ",")),
			0,
			0,
			"",
			make(map[DPFuncParameters]int))
		log.Printf("line: %s, result: %d\n", *line, lineResult)
		result += lineResult
	}
	log.Println("result is", result)
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day12", "input.txt"))
	processLine(lineEmitter)
}
