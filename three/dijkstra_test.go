package main

import (
	"fmt"
	"testing"
	"time"
)

// stolen from stackoverflow: https://stackoverflow.com/a/15312097/11199009
func testEq(a, b []uint64) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func compareTestSerialParallel(g *graph, sourceId uint64, t *testing.T) {
	start := time.Now()
	shortestDistancesSerial := Dijkstra(g, sourceId)
	fmt.Printf("Serial took:\t %v\n", time.Since(start))

	// for i, v := range shortestDistancesSerial {
	// 	fmt.Printf("From vertex %v to %v = %v\n", sourceId, i, v)
	// }

	start2 := time.Now()
	shortestDistancesParallel := DijkstraParallel(g, sourceId)
	fmt.Printf("Parallel took:\t %v\n", time.Since(start2))

	// for i, v := range shortestDistancesParallel {
	// 	fmt.Printf("From vertex %v to %v = %v\n", sourceId, i, v)
	// }

	if !testEq(shortestDistancesSerial, shortestDistancesParallel) {
		t.Errorf("Distances from serial and parallel do not match. Serial: %v, Parallel: %v", shortestDistancesSerial, shortestDistancesParallel)
	}
}

func TestOurFirstGraph(t *testing.T) {
	agraph := newGraph()
	agraph.addEdge(4, 0, 4)
	agraph.addEdge(4, 1, 2)
	agraph.addEdge(0, 1, 1)
	agraph.addEdge(0, 2, 5)
	agraph.addEdge(1, 2, 8)
	agraph.addEdge(1, 3, 10)
	agraph.addEdge(2, 3, 2)
	agraph.addEdge(2, 5, 6)
	agraph.addEdge(3, 5, 2)
	agraph.printGraph()

	graph := GraphFromAdjList(*agraph)
	// hack for that issue in GraphFromAdjList
	graph.ptr = append(graph.ptr, graph.ptr[len(graph.ptr) - 1])
	graph.PrintGraph()

	sourceId := uint64(4)
	
	compareTestSerialParallel(&graph, sourceId, t)
}

func TestRandomGraph(t *testing.T) {
	g := RandomGraph1(100, 0.5)

	// g.PrintGraph()

	compareTestSerialParallel(&g, 0, t)
}
