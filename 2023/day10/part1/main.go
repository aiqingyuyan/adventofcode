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

type Boundary struct {
	right  int
	bottom int
}

func generateGrid(lineEmitter <-chan *string) []string {
	var grid []string
	for line := range lineEmitter {
		grid = append(grid, *line)
	}

	return grid
}

func findStartPoint(grid []string) Coordinate {
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

func getNextConnectedNeighbour(grid []string, point Coordinate, previousConnectedNeighbour Coordinate) Coordinate {
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

func findLoop(grid []string, startPoint Coordinate) int {
	neighbour := getNextConnectedNeighbour(grid, startPoint, Coordinate{})
	previousConnectedNeighbour := startPoint
	stepsToFindLoop := 1

	for neighbour != startPoint {
		nextNeighbour := getNextConnectedNeighbour(grid, neighbour, previousConnectedNeighbour)
		previousConnectedNeighbour = neighbour
		neighbour = nextNeighbour
		stepsToFindLoop++
	}

	log.Println("stepToFindLoop", stepsToFindLoop)

	return stepsToFindLoop
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day10", "input.txt"))
	grid := generateGrid(lineEmitter)
	startPoint := findStartPoint(grid)

	log.Printf("start: %+v", startPoint)

	stepToFindLoop := findLoop(grid, startPoint)

	log.Println("step", stepToFindLoop/2)
}
