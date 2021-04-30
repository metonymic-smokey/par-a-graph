package main

import (
	"fmt"
	"math"
	"sync"
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

// same as above but returns minvalue too
func minVertex2(
	shortestDistances []uint64,
	finalizedVertices []bool,
	sourceId uint64) (uint64, uint64) {
	minId := sourceId

	var min uint64 = math.MaxUint64

	for i, v := range finalizedVertices {
		if !v && shortestDistances[i] <= min {
			min = shortestDistances[i]
			minId = uint64(i)
		}
	}

	return min, minId
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

func dijkstraInnerLoop(
	j uint64,
	g *graph,
	currentVertex uint64,
	shortestDistances []uint64,
	finalizedVertices []bool,
	wg *sync.WaitGroup) {
	defer wg.Done()

	v := g.dst[j]

	if !finalizedVertices[v] &&
		shortestDistances[currentVertex] != math.MaxUint64 &&
		shortestDistances[currentVertex]+g.w[j] < shortestDistances[v] {
		shortestDistances[v] = shortestDistances[currentVertex] + g.w[j]
	}
}

func DijkstraParallel(g *graph, sourceId uint64) []uint64 {
	numVertices := len(g.ptr) - 1

	shortestDistances := make([]uint64, numVertices)
	// instatic := make([]uint64, numVertices)

	incoming_mins := make([]uint64, numVertices)
	outgoing_mins := make([]uint64, numVertices)
	for i := range incoming_mins {
		incoming_mins[i] = math.MaxUint64
	}
	incoming_min_ch := make(chan struct {
		ind uint64
		min uint64
	})

	// update incoming_mins using fan in
	// TODO: how to synchronise and close this when done?
	// TODO: can get more parallelization here by using numVertices goroutines
	go func() {
		for i := range incoming_min_ch {
			if i.min < incoming_mins[i.ind] {
				incoming_mins[i.ind] = i.min
			}
		}
	}()
	// outstatic := make([]uint64, numVertices)

	for i := range shortestDistances {
		shortestDistances[i] = math.MaxUint64
	}

	shortestDistances[sourceId] = 0

	finalizedVertices := make([]bool, numVertices)

	for i := 0; i < numVertices; i++ {
		first, last := g.GetEdgeRange(uint64(i))

		var outgoing_min uint64 = math.MaxUint64
		for j := first; j < last; j++ {
			if g.w[j] < outgoing_min {
				outgoing_min = g.w[j]
			}
			incoming_min_ch <- struct {
				ind uint64
				min uint64
			}{g.dst[j], g.w[j]}
		}

		outgoing_mins[i] = outgoing_min
	}

	close(incoming_min_ch)
	fmt.Println(incoming_mins)
	fmt.Println(outgoing_mins)

	// start phases here
	for phase_no := 0; ; phase_no++ {
		// identification
		minValue, _ := minVertex2(shortestDistances, finalizedVertices, sourceId)

		identifiedVertices := make([]uint64, 0)

		for i := 0; i < numVertices; i++ {
			if !finalizedVertices[i] && shortestDistances[i] != math.MaxUint64 && (shortestDistances[i]-incoming_mins[i] <= minValue ||
				shortestDistances[i] <= minValue+outgoing_mins[i]) {
				identifiedVertices = append(identifiedVertices, uint64(i))
			}
		}

        // fmt.Println("identifiedVertices", identifiedVertices)

		// TODO: break condition
		if len(identifiedVertices) == 0 {
			break
		}

		// settling?
		// numIdentVertices := len(identifiedVertices)
		for _, i := range identifiedVertices {
            finalizedVertices[i] = true
			start, end := g.GetEdgeRange(i)

			for j := start; j < end; j++ {
				if v := g.dst[j]; !finalizedVertices[v] &&
					shortestDistances[i] != math.MaxUint64 &&
					shortestDistances[i]+g.w[j] < shortestDistances[v] {
					shortestDistances[v] = shortestDistances[i] + g.w[j]
				}
			}
		}
	}

	// // question: why only N-1 iterations?
	// // iterations of dijkstra
	// for i := 0; i < numVertices-1; i++ {

	// 	finalizedVertices[currentVertex] = true

	// 	first, last := g.GetEdgeRange(uint64(currentVertex))

	// 	for j := first; j < last; j++ {
	// 		v := g.dst[j]

	// 		if !finalizedVertices[v] &&
	// 			shortestDistances[currentVertex] != math.MaxUint64 &&
	// 			shortestDistances[currentVertex]+g.w[j] < shortestDistances[v] {
	// 			shortestDistances[v] = shortestDistances[currentVertex] + g.w[j]
	// 		}
	// 	}

	// }

	return shortestDistances
}
