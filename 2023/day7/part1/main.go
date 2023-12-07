package main

import (
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"log"
	"path/filepath"
	"strings"
	"sync"
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

// A, K, Q, J, T, 9, 8, 7, 6, 5, 4, 3, or 2
var labelStrength = map[byte]int{
	'2': 0,
	'3': 1,
	'4': 2,
	'5': 3,
	'6': 4,
	'7': 5,
	'8': 6,
	'9': 7,
	'T': 8,
	'J': 9,
	'Q': 10,
	'K': 11,
	'A': 12,
}

func determineHandType(hand string) handType {
	labelMap := make(map[byte]int)
	for _, c := range []byte(hand) {
		if _, ok := labelMap[c]; ok {
			labelMap[c]++
		} else {
			labelMap[c] = 1
		}
	}

	switch len(labelMap) {
	case 1:
		return fiveOfKind
	case 2:
		for _, v := range labelMap {
			if v == 4 {
				return fourOfKind
			}
		}
		return fullHouse
	case 3:
		for _, v := range labelMap {
			if v == 3 {
				return threeOfKind
			}
		}
		return twoPair
	case 4:
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

func processHand(line *string, tree *rbt.Tree, lock *sync.Mutex) {
	subStrs := strings.Split(*line, " ")
	handT := determineHandType(subStrs[0])

	lock.Lock()
	defer lock.Unlock()

	tree.Put(handKey{
		hand: subStrs[0],
		t:    handT,
	}, util.StrToNum(subStrs[1]))
}

func generateTaskFunc(line *string, tree *rbt.Tree, lock *sync.Mutex) executor.TaskFunc {
	return func() int {
		processHand(line, tree, lock)
		return 0
	}
}

func generateTaskEmitter(lineEmitter <-chan *string, tree *rbt.Tree, lock *sync.Mutex) <-chan executor.TaskFunc {
	taskEmitter := make(chan executor.TaskFunc)
	go func() {
		for line := range lineEmitter {
			taskEmitter <- generateTaskFunc(line, tree, lock)
		}
		close(taskEmitter)
	}()
	return taskEmitter
}

func main() {
	e := executor.New(6)

	lineEmitter := util.ReadFile(filepath.Join("2023", "day7", "input.txt"))

	var (
		tree = rbt.NewWith(keyComparator)
		lock sync.Mutex
	)
	taskEmitter := generateTaskEmitter(lineEmitter, tree, &lock)

	e.Run(taskEmitter)

	values := tree.Values()

	result := 0
	for i, v := range values {
		result += v.(int) * (i + 1)
	}

	log.Println("result is", result)
}
