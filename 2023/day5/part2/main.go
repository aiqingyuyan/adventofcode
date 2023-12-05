package main

import (
	"log"
	"math"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"yanyu/aoc/2023/executor"
	"yanyu/aoc/2023/util"
)

type seedPayload struct {
	start int
	end   int
}

type sourceCategoryToDestination struct {
	sourceStart int
	sourceEnd   int
	destStart   int
}

type computationParameters struct {
	seeds                 []seedPayload
	seedToSoil            []sourceCategoryToDestination
	soilToFertilizer      []sourceCategoryToDestination
	fertilizerToWater     []sourceCategoryToDestination
	waterToLight          []sourceCategoryToDestination
	lightToTemperature    []sourceCategoryToDestination
	temperatureToHumidity []sourceCategoryToDestination
	humidityToLocation    []sourceCategoryToDestination
}

func processLine(lineEmitter <-chan *string) *computationParameters {
	var (
		seeds                 []seedPayload
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

	return &computationParameters{
		seeds:                 seeds,
		seedToSoil:            seedToSoil,
		soilToFertilizer:      soilToFertilizer,
		fertilizerToWater:     fertilizerToWater,
		waterToLight:          waterToLight,
		lightToTemperature:    lightToTemperature,
		temperatureToHumidity: temperatureToHumidity,
		humidityToLocation:    humidityToLocation,
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

func constructSeeds(linesBuffer []string) []seedPayload {
	var seeds []seedPayload

	subStrs := strings.Split(linesBuffer[0], " ")
	for i := 1; i < len(subStrs); {
		start := strToNum(subStrs[i])
		length := strToNum(subStrs[i+1])

		seeds = append(seeds, seedPayload{
			start: start,
			end:   start + length - 1,
		})

		i += 2
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

func lookUpLocation(seedNum int, params *computationParameters) int {
	soil := getDestination(seedNum, params.seedToSoil)
	fertilizer := getDestination(soil, params.soilToFertilizer)
	water := getDestination(fertilizer, params.fertilizerToWater)
	light := getDestination(water, params.waterToLight)
	temperature := getDestination(light, params.lightToTemperature)
	humidity := getDestination(temperature, params.temperatureToHumidity)
	return getDestination(humidity, params.humidityToLocation)
}

func generateTaskFunc(lock *sync.Mutex, currentMin *int, seed seedPayload, params *computationParameters) executor.TaskFunc {
	return func() int {
		for seedNum := seed.start; seedNum <= seed.end; seedNum++ {
			lock.Lock()
			*currentMin = min(*currentMin, lookUpLocation(seedNum, params))
			lock.Unlock()
		}

		return 0
	}
}

func generateTaskEmitter(lock *sync.Mutex, currentMin *int, params *computationParameters) <-chan executor.TaskFunc {
	taskEmitter := make(chan executor.TaskFunc)

	go func() {
		for _, seed := range params.seeds {
			log.Printf("emitting seed %d, %d", seed.start, seed.end)
			taskEmitter <- generateTaskFunc(lock, currentMin, seed, params)
		}

		close(taskEmitter)
	}()

	return taskEmitter
}

func main() {
	e := executor.New(runtime.NumCPU())

	lineEmitter := util.ReadFile(filepath.Join("2023", "day5", "input.txt"))

	params := processLine(lineEmitter)

	var (
		currentMin = math.MaxInt
		lock       sync.Mutex
	)

	taskEmitter := generateTaskEmitter(&lock, &currentMin, params)

	e.Run(taskEmitter)

	log.Println("result is", currentMin)
}
