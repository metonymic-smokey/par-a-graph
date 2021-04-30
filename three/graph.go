package main

import (
	"math/rand"
	"fmt"
	"sort"
)

type graph struct {
	dst []uint64
	w []uint64
	ptr []uint64
}

func EmptyGraph(size int) graph {
	g := graph{}
	g.ptr = make([]uint64, size + 1)
	g.dst = make([]uint64, 0)
	g.w = make([]uint64, 0)

	return g
}

func RandomGraph1(nodes int, p float64) graph {
	ag := newGraph()

	for i := 0; i < nodes; i++ {
		ag.nodes[i] = make([]edge, 0)
	}

	for i := 0; i < nodes; i++ {
		for j := 0; j < nodes; j++ {
			rnum := rand.Float64()

			if rnum < p {
				// graph[i] = append(graph[i], {j, })
				// graph[i][j] = mapItem{rand.Float64()}
				ag.addEdge(i, j, rand.Intn(1000))
			}
		}
	}

	g := GraphFromAdjList(*ag)

	return g
}

// returns a sorted slice of keys of a map
func sortedKeys(m map[int][]edge) []int {
	keys := make([]int, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	sort.Ints(keys)

	return keys
}

// TODO: there's an issue here. if a destination node
// doesn't have any outgoing links, then it fails to show
// up in g.ptr and everthing dies
func GraphFromAdjList(g adjGraph) graph {
	newG := EmptyGraph(len(g.nodes))

	lastDegree := 0

	for _, src := range sortedKeys(g.nodes) {
		newG.ptr[src] = uint64(lastDegree)

		for _, dest := range g.nodes[src] {
			newG.dst = append(newG.dst, uint64(dest.node))

			newG.w = append(newG.w, uint64(dest.weight))

			lastDegree += 1
		}
	}

	newG.ptr[len(g.nodes)] = uint64(lastDegree)

	return newG
}

func (g *graph) PrintGraph() {
	fmt.Println("Destination array: ")
	for _, v := range g.dst {
		fmt.Printf(" %v", v)
	}
	fmt.Println()

	fmt.Println("Weight array: ")
	for _, v := range g.w {
		fmt.Printf(" %v", v)
	}
	fmt.Println()

	fmt.Println("Pointer array: ")
	for _, v := range g.ptr {
		fmt.Printf(" %v", v)
	}
	fmt.Println()
}

// func (g *graph) EdgeExists(u uint64, v uint64) bool {
// 	first := g.ptr[u]
// 	last := g.ptr[u+1]

// 	for i := first; i < last; i++ {
// 		if (g.dst[i] == v) {
// 			return true
// 		}
// 	}

// 	return false
// }

func (g *graph) GetEdgeRange(node uint64) (uint64, uint64) {
	return g.ptr[node], g.ptr[node+1]
}
