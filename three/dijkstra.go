package main

import (
	"math"
)

// Finds the vertex with minimum distance value, from the set of vertices not yet included in shortest path tree
func minVertex(
	shortestDistances []uint64,
	finalizedVertices []bool,
	sourceId uint64) uint64 {
	minId := sourceId

	var min uint64 = math.MaxUint64

	for i, v := range finalizedVertices {
		if !v && shortestDistances[i] <= min {
			min = shortestDistances[i]
			minId = uint64(i)
		}
	}

	return minId
}

func Dijkstra(g *graph, sourceId uint64) []uint64 {
	numVertices := len(g.ptr) - 1

	shortestDistances := make([]uint64, numVertices)

	for i := range shortestDistances {
		shortestDistances[i] = math.MaxUint64
	}

	shortestDistances[sourceId] = 0

	finalizedVertices := make([]bool, numVertices)

	// question: why only N-1 iterations?
	// iterations of dijkstra
	for i := 0; i < numVertices-1; i++ {
		currentVertex := minVertex(shortestDistances, finalizedVertices, sourceId)
		// fmt.Println(shortestDistances, finalizedVertices, currentVertex)

		finalizedVertices[currentVertex] = true

		// for j := 0; j < numVertices; j++ {
		// 	if !finalizedVertices[uint64(j)] && g.EdgeExists(uint64(currentVertex), uint64(j)) && shortestDistances[currentVertex] != math.MaxUint64 && shortestDistances[currentVertex] > 0 {
		// 		// ok issues here in the condition
		// 		// how to get the weight corresponding to two given nodes?
		// 		// ok new idea, iterate through the edges like adjacency list
		// 	}
		// }

		first, last := g.GetEdgeRange(uint64(currentVertex))

		for j := first; j < last; j++ {
			if v := g.dst[j]; !finalizedVertices[v] &&
				shortestDistances[currentVertex] != math.MaxUint64 &&
				shortestDistances[currentVertex]+g.w[j] < shortestDistances[v] {
				shortestDistances[v] = shortestDistances[currentVertex] + g.w[j]
			}
		}
	}

	return shortestDistances
}
