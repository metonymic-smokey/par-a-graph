package main

import (
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

const nThreads = 50;

func DijkstraParallel(g *graph, sourceId uint64) []uint64 {
	numVertices := len(g.ptr) - 1

	shortestDistances := make([]uint64, numVertices)

	for i := range shortestDistances {
		shortestDistances[i] = math.MaxUint64
	}

	shortestDistances[sourceId] = 0

	finalizedVertices := make([]bool, numVertices)

	ch := make(chan struct{j uint64; currentVertex uint64})

	var wg sync.WaitGroup
	wg.Add(nThreads)
	for i:=0; i < nThreads; i++ {
		go func() {
			for {
				if el, ok := <-ch; !ok {
					wg.Done()
					return
				} else {
					v := g.dst[el.j]

					if !finalizedVertices[v] &&
						shortestDistances[el.currentVertex] != math.MaxUint64 &&
						shortestDistances[el.currentVertex]+g.w[el.j] < shortestDistances[v] {
						shortestDistances[v] = shortestDistances[el.currentVertex] + g.w[el.j]
					}
				}
			}
		}()
	}

	// question: why only N-1 iterations?
	// iterations of dijkstra
	for i := 0; i < numVertices-1; i++ {
		currentVertex := minVertex(shortestDistances, finalizedVertices, sourceId)

		finalizedVertices[currentVertex] = true

		first, last := g.GetEdgeRange(uint64(currentVertex))

		for j := first; j < last; j++ {
			// wg.Add(1)
			// go dijkstraInnerLoop(j, g, currentVertex, shortestDistances, finalizedVertices, &wg)
			ch <- struct{j uint64; currentVertex uint64}{j, currentVertex}
		}


	}

	close(ch)
		wg.Wait()

	return shortestDistances
}
