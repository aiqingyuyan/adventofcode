package main

import (
	"fmt"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"log"
	"path/filepath"
	"strings"
	"yanyu/aoc/2023/util"
)

func hash(s string) (boxNum int32, label string, isRemoveOp bool, focalLength int) {
	labelBuilder := strings.Builder{}
	for _, c := range s {
		if c == '-' {
			isRemoveOp = true
		} else if c == '=' {
			focalLength = util.StrToNum(s[len(s)-1:])
			break
		} else {
			labelBuilder.WriteRune(c)
			boxNum += c
			boxNum *= 17
			boxNum %= 256
		}
	}

	label = labelBuilder.String()

	return
}

func checkBoxes(boxes []*linkedhashmap.Map) {
	for i, box := range boxes {
		if !box.Empty() {
			contentsBuilder := strings.Builder{}
			box.Each(func(key interface{}, value interface{}) {
				contentsBuilder.WriteString(fmt.Sprintf("[%s %d] ", key, value))
			})

			log.Println("box", i, ":", contentsBuilder.String())
		}
	}
}

func processSequences(sequences []string) []*linkedhashmap.Map {
	boxes := make([]*linkedhashmap.Map, 256)
	for i := range boxes {
		boxes[i] = linkedhashmap.New()
	}

	for _, s := range sequences {
		if boxNum, label, isRemoveOp, focalLength := hash(s); isRemoveOp {
			boxes[boxNum].Remove(label)
		} else {
			boxes[boxNum].Put(label, focalLength)
		}

		//log.Println("after", s)
		//checkBoxes(boxes)
	}

	return boxes
}

func processLines(lineEmitter <-chan *string) []*linkedhashmap.Map {
	line := <-lineEmitter
	sequences := strings.Split(*line, ",")

	return processSequences(sequences)
}

func calculatePower(resultMap []*linkedhashmap.Map) int {
	result := 0
	for i, box := range resultMap {
		if !box.Empty() {
			slotNumber := 1
			box.Each(func(key interface{}, value interface{}) {
				result += (i + 1) * slotNumber * value.(int)
				slotNumber++
			})
		}
	}

	return result
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day15", "input.txt"))
	resultMap := processLines(lineEmitter)
	result := calculatePower(resultMap)

	log.Println("result is", result)
}
