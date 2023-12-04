package main

import (
	"log"
	"path/filepath"
	"sync"
	"yanyu/aoc/2023/executor"
	"yanyu/aoc/2023/util"
)

type cardPayload struct {
	card    *string
	cardIdx int
}

func generateCardPayloadEmitter(cardDeck []*string) <-chan *cardPayload {
	cardEmitter := make(chan *cardPayload)
	go func() {
		for id, card := range cardDeck {
			cardEmitter <- &cardPayload{
				card:    card,
				cardIdx: id,
			}
		}
		close(cardEmitter)
	}()
	return cardEmitter
}

func processCard(payload *cardPayload, cardCount []int, winningCardCountMap []int, lock *sync.Mutex) {
	var (
		seenColon                 bool
		seenBar                   bool
		countOfFoundWinningNumber int
		numBuffer                 []byte
		winningNums               = make(map[string]bool)
	)

	for _, c := range []byte(*payload.card) {
		if c == ':' && !seenColon {
			seenColon = true
			continue
		}

		if c == '|' && !seenBar {
			seenBar = true
			continue
		}

		if seenColon {
			if util.IsByteANumber(c) {
				numBuffer = append(numBuffer, c)
			}

			if c == ' ' && len(numBuffer) > 0 {
				if !seenBar {
					winningNums[string(numBuffer)] = true
				} else if winningNums[string(numBuffer)] {
					countOfFoundWinningNumber++
				}
				numBuffer = nil
			}
		}
	}

	if winningNums[string(numBuffer)] {
		countOfFoundWinningNumber++
	}

	lock.Lock()
	for i := 1; i <= countOfFoundWinningNumber; i++ {
		cardCount[payload.cardIdx+i]++
	}
	lock.Unlock()

	winningCardCountMap[payload.cardIdx] = countOfFoundWinningNumber
}

func generateTaskFunc(payload *cardPayload, cardCount []int, winningCardCountMap []int, lock *sync.Mutex) executor.TaskFunc {
	return func() int {
		processCard(payload, cardCount, winningCardCountMap, lock)
		return 0
	}
}

func generateTaskEmitter(cardPayloadEmitter <-chan *cardPayload, cardCount []int, winningCardCountMap []int, lock *sync.Mutex) chan executor.TaskFunc {
	taskEmitter := make(chan executor.TaskFunc)
	go func() {
		for payload := range cardPayloadEmitter {
			taskEmitter <- generateTaskFunc(payload, cardCount, winningCardCountMap, lock)
		}
		close(taskEmitter)
	}()
	return taskEmitter
}

func updateCardCount(cardCount []int, winningCardCountMap []int) {
	for cardId, winningCardCount := range winningCardCountMap {
		if cardCount[cardId] > 1 {
			for i := 1; i <= winningCardCount; i++ {
				cardCount[cardId+i] += (cardCount[cardId] - 1) * 1
			}
		}
	}
}

func main() {
	e := executor.New(6)

	var (
		cardDeck  []*string
		cardCount []int
		lock      sync.Mutex
	)

	lineEmitter := util.ReadFile(filepath.Join("2023", "day4", "input.txt"))

	for line := range lineEmitter {
		cardCount = append(cardCount, 1)
		cardDeck = append(cardDeck, line)
	}

	winningCardCountMap := make([]int, len(cardDeck))

	cardPayloadEmitter := generateCardPayloadEmitter(cardDeck)

	taskEmitter := generateTaskEmitter(cardPayloadEmitter, cardCount, winningCardCountMap, &lock)

	e.Run(taskEmitter)

	updateCardCount(cardCount, winningCardCountMap)

	result := 0
	for i, count := range cardCount {
		log.Printf("card %d, count: %d", i, count)
		result += count
	}
	log.Println("result is ", result)
}
