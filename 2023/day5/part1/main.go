package main

import (
	"log"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"yanyu/aoc/2023/util"
)

type sourceCategoryToDestination struct {
	sourceStart int
	sourceEnd   int
	destStart   int
}

func printLinesBuffer(linesBuffer []string) {
	for _, line := range linesBuffer {
		log.Println(line)
	}
	log.Println()
}

func printMap(mapName string, mapToPrint map[int]int) {
	log.Println(mapName)
	for k, v := range mapToPrint {
		log.Printf("key %d, value %d", k, v)
	}
}

func strToNum(str string) int {
	if num, err := strconv.Atoi(str); err == nil {
		return num
	} else {
		log.Printf("err: %+v", err)
		panic(err)
	}
}

func constructSeeds(linesBuffer []string) []int {
	var seeds []int
	subStrs := strings.Split(linesBuffer[0], " ")

	for id, seedNum := range subStrs {
		if id == 0 {
			continue
		}

		if num, err := strconv.Atoi(seedNum); err == nil {
			seeds = append(seeds, num)
		} else {
			panic(err)
		}
	}

	return seeds
}

func constructMap(linesBuffer []string) []sourceCategoryToDestination {
	var sourceToDestMap []sourceCategoryToDestination

	for _, line := range linesBuffer {
		subStrs := strings.Split(line, " ")

		destStart := strToNum(subStrs[0])
		sourceStart := strToNum(subStrs[1])
		rangeLength := strToNum(subStrs[2])

		sourceToDestMap = append(sourceToDestMap, sourceCategoryToDestination{
			sourceStart: sourceStart,
			sourceEnd:   sourceStart + rangeLength - 1,
			destStart:   destStart,
		})
	}

	return sourceToDestMap
}

func getDestination(source int, sToDs []sourceCategoryToDestination) int {
	dest := source
	for _, sToD := range sToDs {
		if source >= sToD.sourceStart && source <= sToD.sourceEnd {
			offset := source - sToD.sourceStart
			dest = sToD.destStart + offset
		}
	}

	return dest
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func getSeedLocation(seed int, seedToSoil, soilToFertilizer, fertilizerToWater, waterToLight, lightToTemperature, temperatureToHumidity, humidityToLocation []sourceCategoryToDestination) int {
	soil := getDestination(seed, seedToSoil)
	fertilizer := getDestination(soil, soilToFertilizer)
	water := getDestination(fertilizer, fertilizerToWater)
	light := getDestination(water, waterToLight)
	temperature := getDestination(light, lightToTemperature)
	humidity := getDestination(temperature, temperatureToHumidity)
	location := getDestination(humidity, humidityToLocation)

	return location
}

func processLine(lineEmitter <-chan *string) {
	var (
		seeds                 []int
		seedToSoil            []sourceCategoryToDestination
		soilToFertilizer      []sourceCategoryToDestination
		fertilizerToWater     []sourceCategoryToDestination
		waterToLight          []sourceCategoryToDestination
		lightToTemperature    []sourceCategoryToDestination
		temperatureToHumidity []sourceCategoryToDestination
		humidityToLocation    []sourceCategoryToDestination
	)

	var linesBuffer []string
	for line := range lineEmitter {
		switch {
		case strings.HasPrefix(*line, "seed-to-soil"):
			seeds = constructSeeds(linesBuffer)
			linesBuffer = nil
		case strings.HasPrefix(*line, "soil-to-fertilizer"):
			seedToSoil = constructMap(linesBuffer)
			linesBuffer = nil
		case strings.HasPrefix(*line, "fertilizer-to-water"):
			soilToFertilizer = constructMap(linesBuffer)
			linesBuffer = nil
		case strings.HasPrefix(*line, "water-to-light"):
			fertilizerToWater = constructMap(linesBuffer)
			linesBuffer = nil
		case strings.HasPrefix(*line, "light-to-temperature"):
			waterToLight = constructMap(linesBuffer)
			linesBuffer = nil
		case strings.HasPrefix(*line, "temperature-to-humidity"):
			lightToTemperature = constructMap(linesBuffer)
			linesBuffer = nil
		case strings.HasPrefix(*line, "humidity-to-location"):
			temperatureToHumidity = constructMap(linesBuffer)
			linesBuffer = nil
		case strings.TrimSpace(*line) != "":
			linesBuffer = append(linesBuffer, *line)
		}
	}

	humidityToLocation = constructMap(linesBuffer)

	var currentMin = math.MaxInt
	for _, seed := range seeds {
		currentMin = min(currentMin, getSeedLocation(seed, seedToSoil, soilToFertilizer, fertilizerToWater, waterToLight, lightToTemperature, temperatureToHumidity, humidityToLocation))
	}

	log.Println("result is ", currentMin)
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day5", "input.txt"))

	processLine(lineEmitter)
}
