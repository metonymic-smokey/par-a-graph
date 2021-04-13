package main

import "fmt"

type edge struct {
	node   int
	weight int
}

type adjGraph struct {
	nodes map[int][]edge
}

func newGraph() *adjGraph {
	return &adjGraph{nodes: make(map[int][]edge)}
}

func (g *adjGraph) addEdge(origin int, destination int, weight int) {
	g.nodes[origin] = append(g.nodes[origin], edge{node: destination, weight: weight})
}

//prints adjGraph as adjacency list
func (g *adjGraph) printGraph() {

	for k, v := range g.nodes {
		fmt.Print(k, ": ")
		for _, dest := range v {
			fmt.Print(dest, ",")
		}
		fmt.Println()
	}

}
