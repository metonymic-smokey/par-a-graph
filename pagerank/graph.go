package main

import (
	"fmt"
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
