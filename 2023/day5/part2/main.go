package main

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"yanyu/aoc/2023/util"
)

type seedPayload struct {
	start int
	end   int
}

type sourceCategoryToDestination struct {
	sourceStart int
	destStart   int
	rangeLength int
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

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
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
			destStart:   destStart,
			rangeLength: rangeLength,
		})
	}

	return sourceToDestMap
}

func getSource(dest int, sToDs []sourceCategoryToDestination) int {
	source := dest
	for _, sToD := range sToDs {
		if sToD.destStart <= dest {
			if dest-sToD.destStart+1 <= sToD.rangeLength {
				offset := dest - sToD.destStart
				source = sToD.sourceStart + offset

				return source
			}
		}
	}

	return source
}

func reverseLookUp(location int, param *computationParameters) bool {
	humidity := getSource(location, param.humidityToLocation)
	temperature := getSource(humidity, param.temperatureToHumidity)
	light := getSource(temperature, param.lightToTemperature)
	water := getSource(light, param.waterToLight)
	fertilizer := getSource(water, param.fertilizerToWater)
	soil := getSource(fertilizer, param.soilToFertilizer)
	seed := getSource(soil, param.seedToSoil)

	for _, s := range param.seeds {
		if seed >= s.start && seed <= s.end {
			return true
		}
	}

	return false
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day5", "input.txt"))

	params := processLine(lineEmitter)

	limit := params.seeds[0].end
	for i := 1; i < len(params.seeds); i++ {
		limit = max(limit, params.seeds[i].end)
	}

	var location = 0
	for ; location <= limit; location++ {
		if reverseLookUp(location, params) {
			break
		}
	}

	log.Println("result is", location)

}
