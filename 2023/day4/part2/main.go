package main

import (
	"log"
	"path/filepath"
	"yanyu/aoc/2023/executor"
	"yanyu/aoc/2023/util"
)

type cardPayload struct {
	card    *string
	cardIdx int
}

func processLine(lineEmitter <-chan *string) (cardCount []int, cardDeck []*string) {
	for line := range lineEmitter {
		cardCount = append(cardCount, 1)
		cardDeck = append(cardDeck, line)
	}
	return cardCount, cardDeck
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

type processedCardResult struct {
	cardIdx                              int
	countOfFoundWinningNumberForThisCard int
}

func processCard(payload *cardPayload) processedCardResult {
	var (
		seenColon                            bool
		seenBar                              bool
		countOfFoundWinningNumberForThisCard int
		numBuffer                            []byte
		winningNumsToCheck                   = make(map[string]bool)
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
					winningNumsToCheck[string(numBuffer)] = true
				} else if winningNumsToCheck[string(numBuffer)] {
					countOfFoundWinningNumberForThisCard++
				}
				numBuffer = nil
			}
		}
	}

	if winningNumsToCheck[string(numBuffer)] {
		countOfFoundWinningNumberForThisCard++
	}

	return processedCardResult{
		cardIdx:                              payload.cardIdx,
		countOfFoundWinningNumberForThisCard: countOfFoundWinningNumberForThisCard,
	}
}

func generateTaskFunc(payload *cardPayload) executor.TaskFunc {
	return func() any {
		return processCard(payload)
	}
}

func generateTaskEmitter(cardPayloadEmitter <-chan *cardPayload) chan executor.TaskFunc {
	taskEmitter := make(chan executor.TaskFunc)
	go func() {
		for payload := range cardPayloadEmitter {
			taskEmitter <- generateTaskFunc(payload)
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

	lineEmitter := util.ReadFile(filepath.Join("2023", "day4", "input.txt"))

	cardCount, cardDeck := processLine(lineEmitter)

	cardPayloadEmitter := generateCardPayloadEmitter(cardDeck)

	taskEmitter := generateTaskEmitter(cardPayloadEmitter)

	winningCardCountMap := make([]int, len(cardDeck))

	resultHandleFunc := func(taskFuncResult any) {
		result := taskFuncResult.(processedCardResult)

		for i := 1; i <= result.countOfFoundWinningNumberForThisCard; i++ {
			cardCount[result.cardIdx+i]++
		}

		winningCardCountMap[result.cardIdx] = result.countOfFoundWinningNumberForThisCard
	}

	e.Run(taskEmitter, resultHandleFunc)

	updateCardCount(cardCount, winningCardCountMap)

	result := 0
	for i, count := range cardCount {
		log.Printf("card %d, count: %d", i, count)
		result += count
	}
	log.Println("result is ", result)
}
