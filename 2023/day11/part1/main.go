package main

import (
	"log"
	"math"
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

func expandColumn(matrix [][]byte, col int) {
	for rowNum := range matrix {
		matrix[rowNum] = append(matrix[rowNum][0:col+1], matrix[rowNum][col:]...)
	}
}

func expandMatrix(matrix [][]byte) [][]byte {
	for row := 0; row < len(matrix); row++ {
		if _, ok := isValidRow(matrix[row]); !ok {
			matrix = append(matrix[0:row+1], matrix[row:]...)
			row++
		}
	}

	for col := 0; col < len(matrix[0]); col++ {
		if !isValidColumn(matrix, col) {
			expandColumn(matrix, col)
			col++
		}
	}

	return matrix
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

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day11", "input.txt"))
	matrix := generateMatrix(lineEmitter)
	matrix = expandMatrix(matrix)
	galaxies := findAllGalaxies(matrix)
	pairsOfGalaxies := generateGalaxiesPair(galaxies)

	sum := 0
	for _, pair := range pairsOfGalaxies {
		minDistance := math.Abs(float64(pair.top.row-pair.bottom.row)) + math.Abs(float64(pair.top.col-pair.bottom.col))
		sum += int(minDistance)
	}

	log.Println("result is", sum)
}
