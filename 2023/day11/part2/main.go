package main

import (
	"log"
	"path/filepath"
	"yanyu/aoc/2023/util"
)

type Galaxy struct {
	assignedNumber int
	row            int
	col            int
}

type GalaxyPair struct {
	top, bottom Galaxy
}

func generateMatrix(lineEmitter <-chan *string) [][]byte {
	var matrix [][]byte
	for line := range lineEmitter {
		matrix = append(matrix, []byte(*line))
	}

	return matrix
}

func isValidRow(row []byte) (int, bool) {
	for colNum, c := range row {
		if c == '#' {
			return colNum, true
		}
	}
	return -1, false
}

func isValidColumn(matrix [][]byte, col int) bool {
	for _, row := range matrix {
		if row[col] == '#' {
			return true
		}
	}

	return false
}

func expandMatrix(matrix [][]byte) (rowsNeedsToBeExpanded, colsNeedsToBeExpanded map[int]bool) {
	rowsNeedsToBeExpanded = make(map[int]bool)

	for row := 0; row < len(matrix); row++ {
		if _, ok := isValidRow(matrix[row]); !ok {
			rowsNeedsToBeExpanded[row] = true
		}
	}

	colsNeedsToBeExpanded = make(map[int]bool)
	for col := 0; col < len(matrix[0]); col++ {
		if !isValidColumn(matrix, col) {
			colsNeedsToBeExpanded[col] = true
		}
	}

	return
}

func getAllGalaxies(matrix [][]byte, galaxyNumber *int, row, col int) []Galaxy {
	var galaxies []Galaxy
	for ; col < len(matrix[row]); col++ {
		if matrix[row][col] == '#' {
			galaxies = append(galaxies, Galaxy{
				assignedNumber: *galaxyNumber,
				row:            row,
				col:            col,
			})
			*galaxyNumber++
		}
	}

	return galaxies
}

func findAllGalaxies(matrix [][]byte) []Galaxy {
	var (
		galaxies       []Galaxy
		assignedNumber int
	)
	for rowNum, row := range matrix {
		if col, ok := isValidRow(row); ok {
			foundGalaxies := getAllGalaxies(matrix, &assignedNumber, rowNum, col)

			galaxies = append(galaxies, foundGalaxies...)
		}
	}

	return galaxies
}

func generateGalaxiesPair(galaxies []Galaxy) []GalaxyPair {
	var galaxyPairs []GalaxyPair
	for i := 0; i < len(galaxies)-1; i++ {
		for j := i + 1; j < len(galaxies); j++ {
			galaxyPairs = append(galaxyPairs, GalaxyPair{
				top:    galaxies[i],
				bottom: galaxies[j],
			})
		}
	}

	return galaxyPairs
}

const multiplier = 1000000

func calculateDistanceOfGalaxiesPair(pair GalaxyPair, rowsNeedsToBeExpanded, colsNeedsToBeExpanded map[int]bool) int {
	minRow := util.Min(pair.top.row, pair.bottom.row)
	maxRow := util.Max(pair.top.row, pair.bottom.row)
	rowDistance := 0
	for i, j := minRow, minRow+1; j <= maxRow; {
		if rowsNeedsToBeExpanded[j] {
			rowDistance += (j - i) * multiplier
		} else {
			rowDistance++
		}

		i++
		j++
	}

	minCol := util.Min(pair.top.col, pair.bottom.col)
	maxCol := util.Max(pair.top.col, pair.bottom.col)
	colDistance := 0
	for i, j := minCol, minCol+1; j <= maxCol; {
		if colsNeedsToBeExpanded[j] {
			colDistance += (j - i) * multiplier
		} else {
			colDistance++
		}

		i++
		j++
	}

	return rowDistance + colDistance
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day11", "input.txt"))
	matrix := generateMatrix(lineEmitter)
	rowsNeedsToBeExpanded, colsNeedsToBeExpanded := expandMatrix(matrix)
	galaxies := findAllGalaxies(matrix)
	pairsOfGalaxies := generateGalaxiesPair(galaxies)

	sum := 0
	for _, pair := range pairsOfGalaxies {
		pairDistance := calculateDistanceOfGalaxiesPair(pair, rowsNeedsToBeExpanded, colsNeedsToBeExpanded)

		minDistance := pairDistance
		sum += minDistance
	}

	log.Println("result is", sum)
}
