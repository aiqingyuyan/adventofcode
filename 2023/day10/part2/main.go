package main

import (
	"log"
	"path/filepath"
	"yanyu/aoc/2023/util"
)

type Coordinate struct {
	row int
	col int
}

func generateGrid(lineEmitter <-chan *string) [][]byte {
	var grid [][]byte
	for line := range lineEmitter {
		grid = append(grid, []byte(*line))
	}

	return grid
}

func findStartPoint(grid [][]byte) Coordinate {
	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[row]); col++ {
			if grid[row][col] == 'S' {
				return Coordinate{
					row: row,
					col: col,
				}
			}
		}
	}

	return Coordinate{
		-1, -1,
	}
}

func getNorthNeighbour(point Coordinate) Coordinate {
	return Coordinate{
		row: point.row - 1,
		col: point.col,
	}
}

func getSouthNeighbour(point Coordinate) Coordinate {
	return Coordinate{
		row: point.row + 1,
		col: point.col,
	}
}

func getWestNeighbour(point Coordinate) Coordinate {
	return Coordinate{
		row: point.row,
		col: point.col - 1,
	}
}

func getEastNeighbour(point Coordinate) Coordinate {
	return Coordinate{
		row: point.row,
		col: point.col + 1,
	}
}

func getNextConnectedNeighbour(grid [][]byte, point Coordinate, previousConnectedNeighbour Coordinate) Coordinate {
	tile := grid[point.row][point.col]

	switch tile {
	case '|':
		north := getNorthNeighbour(point)
		if north != previousConnectedNeighbour {
			return north
		}
		return getSouthNeighbour(point)
	case '-':
		west := getWestNeighbour(point)
		if west != previousConnectedNeighbour {
			return west
		}
		return getEastNeighbour(point)
	case 'L':
		north := getNorthNeighbour(point)
		if north != previousConnectedNeighbour {
			return north
		}
		return getEastNeighbour(point)
	case 'J':
		north := getNorthNeighbour(point)
		if north != previousConnectedNeighbour {
			return north
		}
		return getWestNeighbour(point)
	case '7':
		west := getWestNeighbour(point)
		if west != previousConnectedNeighbour {
			return west
		}
		return getSouthNeighbour(point)
	case 'F':
		east := getEastNeighbour(point)
		if east != previousConnectedNeighbour {
			return east
		}
		return getSouthNeighbour(point)
	case 'S':
		south := getSouthNeighbour(point)
		if grid[south.row][south.col] == '|' || grid[south.row][south.col] == 'L' || grid[south.row][south.col] == 'J' {
			return south
		}

		return getEastNeighbour(point)
	default:
		return Coordinate{-1, -1}
	}
}

func findLoop(grid [][]byte, startPoint Coordinate) map[Coordinate]bool {
	neighbour := getNextConnectedNeighbour(grid, startPoint, Coordinate{})

	previousConnectedNeighbour := startPoint

	loopPoints := map[Coordinate]bool{
		previousConnectedNeighbour: true,
	}

	for neighbour != startPoint {
		nextNeighbour := getNextConnectedNeighbour(grid, neighbour, previousConnectedNeighbour)

		previousConnectedNeighbour = neighbour
		neighbour = nextNeighbour

		loopPoints[previousConnectedNeighbour] = true
	}

	grid[startPoint.row][startPoint.col] = 'F'

	return loopPoints
}

func pointBelongToLoopPipe(point Coordinate, loopPoints map[Coordinate]bool) bool {
	_, ok := loopPoints[point]

	return ok
}

func isPointInsideLoop(loopPoints map[Coordinate]bool, row int, leftPart []byte) bool {
	countOfVerticalBarLJ := 0
	for col, c := range leftPart {
		if pointBelongToLoopPipe(Coordinate{row, col}, loopPoints) {
			if c == '|' || c == 'L' || c == 'J' {
				countOfVerticalBarLJ++
			}
		}
	}

	if countOfVerticalBarLJ%2 == 0 {
		return false
	}

	return true
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day10", "input.txt"))
	grid := generateGrid(lineEmitter)
	startPoint := findStartPoint(grid)

	log.Printf("start: %+v", startPoint)

	loopPoints := findLoop(grid, startPoint)

	result := 0
	for row := 0; row < len(grid); row++ {
		for col := 0; col < len(grid[row]); col++ {
			if !pointBelongToLoopPipe(Coordinate{row, col}, loopPoints) {
				// ray casting - (check if inside polygon)
				if isPointInsideLoop(loopPoints, row, grid[row][0:col]) {
					result++
				}
			}
		}
	}

	log.Println("result is", result)
}
