package main

import (
	"fmt"
	"sync"
)

type edge struct {
	node   int64
	weight int64
}

type graph struct {
	nodes map[int64][]edge
}

func newGraph() *graph {
	return &graph{nodes: make(map[int64][]edge)}
}

func (g *graph) addEdge(origin int64, destination int64, weight int64) {
	g.nodes[origin] = append(g.nodes[origin], edge{node: destination, weight: weight})
}

//prints graph as adjacency list
func (g *graph) printGraph() {

	for k, v := range g.nodes {
		fmt.Print(k, ": ")
		for _, dest := range v {
			fmt.Print(dest, ",")
		}
		fmt.Println()
	}

}

func (g *graph) getPath(origin int64, destination int64) (int64, []int64) {
	h := newHeap()
	h.push(path{value: 0, nodes: []int64{origin}})
	visited := make(map[int64]bool)

	for len(*h.values) > 0 {

		// Find the nearest yet to visit node
		p := h.pop()
		node := p.nodes[len(p.nodes)-1]

		if visited[node] {
			continue
		}

		if node == destination {
			return p.value, p.nodes
		}

		g_edges := g.nodes[node]

		var wg sync.WaitGroup
		wg.Add(len(g_edges))

		for _, e := range g_edges {
			go func(e edge) {
				defer wg.Done()
				if !visited[e.node] {
					h.push(path{value: p.value + e.weight, nodes: append([]int64{}, append(p.nodes, e.node)...)})
				}
			}(e)
		}
		wg.Wait()

		visited[node] = true
	}

	return 0, nil
}
