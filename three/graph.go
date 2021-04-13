package main

import (
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
	g.ptr = make([]uint64, size)
	g.dst = make([]uint64, 0)
	g.w = make([]uint64, 0)

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
