package main

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"github.com/emirpasic/gods/trees/binaryheap"
	"log"
	"path/filepath"
	"strings"
	"yanyu/aoc/2023/util"
)

func processLines(lineEmitter <-chan *string) (graph Graph) {
	graph.vertexNeighbourWeightMap = make(map[string]map[string]int)

	for line := range lineEmitter {
		subStrs := strings.Split(*line, ": ")
		src := subStrs[0]
		neighbours := strings.Split(subStrs[1], " ")

		for _, neighbour := range neighbours {
			graph.addEdge(src, neighbour)

			graph.addEdge(src, neighbour)
		}
	}

	return
}

type Graph struct {
	vertexNeighbourWeightMap map[string]map[string]int
}

func (g Graph) numberOfVertex() int {
	return len(g.vertexNeighbourWeightMap)
}

func (g Graph) numberOfEdges() int {
	sum := 0
	for v := range g.vertexNeighbourWeightMap {
		sum += len(g.vertexNeighbourWeightMap[v])
	}
	return sum / 2
}

func (g Graph) addEdge(u, v string) {
	if _, ok := g.vertexNeighbourWeightMap[u]; !ok {
		g.vertexNeighbourWeightMap[u] = make(map[string]int)
	}

	g.vertexNeighbourWeightMap[u][v] = 1

	if _, ok := g.vertexNeighbourWeightMap[v]; !ok {
		g.vertexNeighbourWeightMap[v] = make(map[string]int)
	}

	g.vertexNeighbourWeightMap[v][u] = 1
}

func (g Graph) getVertices() (vertices *hashset.Set) {
	vertices = hashset.New()
	for v := range g.vertexNeighbourWeightMap {
		vertices.Add(v)
	}
	return
}

func (g Graph) getVertexWeight(v string) int {
	sum := 0
	for _, w := range g.vertexNeighbourWeightMap[v] {
		sum += w
	}
	return sum
}

func (g Graph) getEdgeWeight(u, v string) (weight int, ok bool) {
	weight, ok = g.vertexNeighbourWeightMap[u][v]
	return
}

func (g Graph) mergeVertices(s, t string) {
	tNeighbours := g.vertexNeighbourWeightMap[t]

	for tNeighbour, edgeWeight := range tNeighbours {
		if _, ok := g.vertexNeighbourWeightMap[s][tNeighbour]; ok {
			g.vertexNeighbourWeightMap[s][tNeighbour] += edgeWeight
			g.vertexNeighbourWeightMap[tNeighbour][s] += edgeWeight
		} else if tNeighbour != s {
			g.vertexNeighbourWeightMap[s][tNeighbour] = edgeWeight
			g.vertexNeighbourWeightMap[tNeighbour][s] = edgeWeight
		}
	}

	delete(g.vertexNeighbourWeightMap, t)
	for v := range g.vertexNeighbourWeightMap {
		delete(g.vertexNeighbourWeightMap[v], t)
	}
}

type CutOfPhase struct {
	s, t      string
	cutWeight int
}

type VertexWithWeightSum struct {
	v         string
	weightSum int
}

func vertexWithWeightSumComparator(a, b VertexWithWeightSum) int {
	if a.weightSum > b.weightSum {
		return 1
	} else if a.weightSum < b.weightSum {
		return -1
	} else {
		return 0
	}
}

func (g Graph) minCutPhase() CutOfPhase {
	vertices := g.getVertices()
	setA := linkedhashset.New(vertices.Values()[0])
	numberOfVertices := vertices.Size()

	// max heap
	maxHeap := binaryheap.NewWith(func(a, b interface{}) int {
		return -vertexWithWeightSumComparator(a.(VertexWithWeightSum), b.(VertexWithWeightSum))
	})

	var vertexTWithWeightSum VertexWithWeightSum
	for setA.Size() != numberOfVertices {
		for _, vertex := range vertices.Values() {
			if !setA.Contains(vertex) {
				vertexWithWeightSum := VertexWithWeightSum{v: vertex.(string)}

				for _, nextVertexInA := range setA.Values() {
					if w, ok := g.getEdgeWeight(vertex.(string), nextVertexInA.(string)); ok {
						vertexWithWeightSum.weightSum += w
					}
				}

				maxHeap.Push(vertexWithWeightSum)

				if setA.Size() == numberOfVertices-1 {
					vertexTWithWeightSum = vertexWithWeightSum
				}
			}
		}

		mostTightlyConnectedVertex, _ := maxHeap.Pop()
		setA.Add(mostTightlyConnectedVertex.(VertexWithWeightSum).v)
		vertices.Remove(mostTightlyConnectedVertex)

		maxHeap.Clear()
	}

	setAValues := setA.Values()
	s, t := setAValues[len(setAValues)-2].(string), setAValues[len(setAValues)-1].(string)

	// merge s, t (t into s)
	g.mergeVertices(s, t)

	return CutOfPhase{s, t, vertexTWithWeightSum.weightSum}
}

func (g Graph) minCut() (CutOfPhase, []string) {
	var (
		firstPartition       []string
		bestCut              CutOfPhase
		partitionOfEachPhase = make(map[string][]string)
	)

	for v := range g.vertexNeighbourWeightMap {
		partitionOfEachPhase[v] = []string{v}
	}

	for len(g.vertexNeighbourWeightMap) > 1 {
		cutOfPhase := g.minCutPhase()

		if bestCut.cutWeight == 0 || cutOfPhase.cutWeight < bestCut.cutWeight {
			bestCut = cutOfPhase
			firstPartition = partitionOfEachPhase[bestCut.t]
		}

		partitionOfEachPhase[cutOfPhase.s] = append(partitionOfEachPhase[cutOfPhase.s], partitionOfEachPhase[cutOfPhase.t]...)
	}

	return bestCut, firstPartition
}

func main() {
	lineEmitter := util.ReadFile(filepath.Join("2023", "day25", "input2.txt"))
	graph := processLines(lineEmitter)

	numOfVertices := graph.numberOfVertex()

	bestCut, firstPartition := graph.minCut()

	log.Println("first partition", firstPartition)
	log.Printf("minCut: %d, result %d", bestCut.cutWeight, (numOfVertices-len(firstPartition))*len(firstPartition))
}
