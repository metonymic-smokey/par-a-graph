package main

import (
	"fmt"
	"math"
)

type mapItem struct {
	//attr int
	cost float64
}

type node struct {
	weight   float64
	outbound float64
}

type adjMap struct {
	edges map[int]map[int]mapItem
	nodes map[int]*node
}

type adjList map[int]mapItem

func PrintAdjMap(adjacencyMap adjMap) {
	for n1, list := range adjacencyMap.edges {
		fmt.Printf("%v ->", n1)
		for n2, item := range list {
			fmt.Printf(" (%v, %v)", n2, item.cost)
		}
		fmt.Println()
	}
}

func Rank(graph adjMap, alpha, epsilon float64, callback func(id int, rank float64)) {

	delta := float64(1.0)

	inverse := 1 / float64(len(graph.nodes))

	for source := range graph.edges {
		if graph.nodes[source].outbound > 0 {
			for target := range graph.edges[source] {
				//fmt.Printf("%T\n",graph.edges[source][target])
				var temp mapItem
				temp.cost = graph.edges[source][target].cost / graph.nodes[source].outbound
				graph.edges[source][target] = temp
			}
		}
	}

	for key := range graph.nodes {
		graph.nodes[key].weight = inverse
	}

	for delta > epsilon {
		leak := float64(0)
		nodes := map[int]float64{}

		for key, value := range graph.nodes {
			nodes[key] = value.weight

			if value.outbound == 0 {
				leak += value.weight
			}

			graph.nodes[key].weight = 0
		}

		leak *= alpha

		for source := range graph.nodes {
			for target, _ := range graph.edges[source] {
				var weight mapItem
				weight = graph.edges[source][target]
				//fmt.Printf("%T %T %T\n",nodes[source],alpha,weight)
				//graph.nodes[target].weight += alpha * nodes[source] * weight.cost
				_, _ = target, weight
			}

			graph.nodes[source].weight += (1-alpha)*inverse + leak*inverse
		}

		delta = 0

		for key, value := range graph.nodes {
			delta += math.Abs(value.weight - nodes[key])
		}

		for key, value := range graph.nodes {
			callback(key, value.weight)
		}
	}

}
