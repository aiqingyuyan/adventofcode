package main

import (
	"fmt"
	"gonum.org/v1/gonum/stat/combin"
	"log"
	"path/filepath"
	"strings"
	"yanyu/aoc/2023/util"
)

func toNums(groupsOfBroken []string) []int {
	var result []int
	for _, numStr := range groupsOfBroken {
		result = append(result, util.StrToNum(numStr))
	}
	return result
}

func findNumOfExistingBrokenRecordAndIndexesOfUnknown(records []byte) (numOfExistingBrokenRecord int, indexesOfUnknown []int) {
	for idx, r := range records {
		if r == '#' {
			numOfExistingBrokenRecord++
		}

		if r == '?' {
			indexesOfUnknown = append(indexesOfUnknown, idx)
		}
	}

	return
}

func getTotalNumOfBrokenRecords(groupsOfBroken []int) int {
	total := 0
	for _, num := range groupsOfBroken {
		total += num
	}

	return total
}

func findFillsCombination(indexesOfUnknown []int, remainingNumOfBrokenRecords int) [][]int {
	combinations := combin.Combinations(len(indexesOfUnknown), remainingNumOfBrokenRecords)
	for _, c := range combinations {
		for i := 0; i < remainingNumOfBrokenRecords; i++ {
			c[i] = indexesOfUnknown[c[i]]
		}
	}

	return combinations
}

func resetRecords(records []byte, combination []int) {
	for _, recordIdx := range combination {
		records[recordIdx] = '?'
	}
}

func checkCombination(combination []int, records []byte, groupsOfBroken []int) int {
	// fill '#'
	for _, recordIdx := range combination {
		records[recordIdx] = '#'
	}

	//log.Printf("after fill: %s", records)

	defer resetRecords(records, combination)

	// check
	var (
		currentGroupIdx = 0
		numOfDamaged    = 0
	)
	for i, j := 0, 0; j < len(records); {
		if records[j] == '#' {
			i = j

			for j < len(records) && records[j] == '#' {
				j++
			}

			numOfDamaged = j - i

			if numOfDamaged != groupsOfBroken[currentGroupIdx] {
				return 0
			}

			currentGroupIdx++

			if currentGroupIdx >= len(groupsOfBroken) {
				return 1
			}
		} else {
			j++
		}
	}

	return 1
}

func findCombination(records []byte, groupsOfBroken []int) int {
	numOfExistingBrokenRecord, indexesOfUnknown := findNumOfExistingBrokenRecordAndIndexesOfUnknown(records)
	totalNumOfBrokenRecords := getTotalNumOfBrokenRecords(groupsOfBroken)
	remainingNumOfBrokenRecords := totalNumOfBrokenRecords - numOfExistingBrokenRecord
	combinations := findFillsCombination(indexesOfUnknown, remainingNumOfBrokenRecords)

	//log.Printf("totalNumOfBrokenRecords: %d, remainingNumOfBrokenRecords: %d", totalNumOfBrokenRecords, remainingNumOfBrokenRecords)

	result := 0
	for _, c := range combinations {
		//log.Printf("%+v", c)
		result += checkCombination(c, records, groupsOfBroken)

		//log.Printf("after check: %s", records)
	}

	//log.Printf("records: %s, combinations: %d", records, result)

	return result
}

func processLine(lineEmitter <-chan *string) {
	result := 0
	for line := range lineEmitter {
		subStrs := strings.Split(*line, " ")
		groupsOfBroken := toNums(strings.Split(subStrs[1], ","))
		lineResult := findCombination([]byte(subStrs[0]), groupsOfBroken)
		fmt.Printf("line: %s, result: %d\n", *line, lineResult)
		result += lineResult
	}
	log.Println("result is", result)
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day12", "input.txt"))
	processLine(lineEmitter)
}
