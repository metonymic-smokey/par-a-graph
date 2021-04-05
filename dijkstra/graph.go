package main

import (
	"fmt"
	"sync"
)

type edge struct {
	node   int
	weight int
}

type graph struct {
	nodes map[int][]edge
}

func newGraph() *graph {
	return &graph{nodes: make(map[int][]edge)}
}

func (g *graph) addEdge(origin int, destination int, weight int) {
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

func (g *graph) getPath(origin int, destination int) (int, []int) {
	h := newHeap()
	h.push(path{value: 0, nodes: []int{origin}})
	visited := make(map[int]bool)

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
					h.push(path{value: p.value + e.weight, nodes: append([]int{}, append(p.nodes, e.node)...)})
				}
			}(e)
		}
		wg.Wait()

		visited[node] = true
	}

	return 0, nil
}
