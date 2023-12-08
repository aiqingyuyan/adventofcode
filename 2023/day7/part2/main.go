package main

import (
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"log"
	"path/filepath"
	"strings"
	"yanyu/aoc/2023/executor"
	"yanyu/aoc/2023/util"
)

type handType int

type handKey struct {
	hand string
	t    handType
}

const (
	highCard = 1 << iota
	onePair
	twoPair
	threeOfKind
	fullHouse
	fourOfKind
	fiveOfKind
)

var labelStrength = map[byte]int{
	'J': 0,
	'2': 1,
	'3': 2,
	'4': 3,
	'5': 4,
	'6': 5,
	'7': 6,
	'8': 7,
	'9': 8,
	'T': 9,
	'Q': 10,
	'K': 11,
	'A': 12,
}

func determineHandType(hand string) handType {
	numOfJ := 0
	labelMap := make(map[byte]int)
	for _, c := range []byte(hand) {
		if c == 'J' {
			numOfJ++
			continue
		}

		if _, ok := labelMap[c]; ok {
			labelMap[c]++
		} else {
			labelMap[c] = 1
		}
	}

	switch len(labelMap) {
	case 0: // 5 J
		return fiveOfKind
	case 1:
		return fiveOfKind
	case 2: // j can appear max of 3
		if numOfJ == 0 {
			for _, v := range labelMap {
				if v == 4 {
					return fourOfKind
				}
			}
			return fullHouse
		} else if numOfJ == 1 {
			currentMaxNum := 2
			for _, v := range labelMap {
				if v > currentMaxNum {
					currentMaxNum = v
				}
			}

			if currentMaxNum == 2 {
				return fullHouse
			} else {
				return fourOfKind
			}
		} else {
			return fourOfKind
		}
	case 3: // J can appear max of 2
		if numOfJ == 0 {
			for _, v := range labelMap {
				if v == 3 {
					return threeOfKind
				}
			}
			return twoPair
		} else {
			return threeOfKind
		}
	case 4: // J can appear max of 1
		return onePair
	default:
		return highCard
	}
}

func keyComparator(a interface{}, b interface{}) int {
	key1 := a.(handKey)
	key2 := b.(handKey)

	if key1.t > key2.t {
		return 1
	} else if key1.t < key2.t {
		return -1
	} else {
		for i, v := range []byte(key1.hand) {
			if labelStrength[v] > labelStrength[key2.hand[i]] {
				return 1
			} else if labelStrength[v] < labelStrength[key2.hand[i]] {
				return -1
			}
		}

		return 0
	}
}

type handPayload struct {
	handKey handKey
	bid     int
}

func processHand(line *string) handPayload {
	subStrs := strings.Split(*line, " ")
	handT := determineHandType(subStrs[0])

	return handPayload{
		handKey: handKey{
			hand: subStrs[0],
			t:    handT,
		},
		bid: util.StrToNum(subStrs[1]),
	}
}

func generateTaskFunc(line *string) executor.TaskFunc {
	return func() any {
		return processHand(line)
	}
}

func generateTaskEmitter(lineEmitter <-chan *string) <-chan executor.TaskFunc {
	taskEmitter := make(chan executor.TaskFunc)
	go func() {
		for line := range lineEmitter {
			taskEmitter <- generateTaskFunc(line)
		}
		close(taskEmitter)
	}()
	return taskEmitter
}

func main() {
	e := executor.New(6)

	lineEmitter := util.ReadFile(filepath.Join("2023", "day7", "input.txt"))

	taskEmitter := generateTaskEmitter(lineEmitter)

	tree := rbt.NewWith(keyComparator)
	resultHandleFunc := func(taskFuncResult any) {
		payload := taskFuncResult.(handPayload)
		tree.Put(payload.handKey, payload.bid)
	}

	e.Run(taskEmitter, resultHandleFunc)

	values := tree.Values()

	result := 0
	for i, v := range values {
		result += v.(int) * (i + 1)
	}

	log.Println("result is", result)
}
