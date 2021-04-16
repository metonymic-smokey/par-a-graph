package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	graph := newGraph()

	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var fileContents []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		fileContents = append(fileContents, scanner.Text())
	}

	var adjacencyMap adjMap
	adjacencyMap.edges = make(map[int]map[int]mapItem)
	adjacencyMap.nodes = make(map[int]*node)

	for _, fc := range fileContents {
		res := strings.SplitN(fc, ",", -1)
		src, _ := strconv.Atoi(res[0])
		dest, _ := strconv.Atoi(res[1])
		wt, _ := strconv.Atoi(res[2])
		graph.addEdge(src, dest, wt)
		if _, ok := adjacencyMap.nodes[src]; !ok {
			adjacencyMap.nodes[src] = &node{
				weight:   0,
				outbound: 0,
			}
		}
		if _, ok := adjacencyMap.nodes[dest]; !ok {
			adjacencyMap.nodes[dest] = &node{
				weight:   0,
				outbound: 0,
			}
		}
		if _, ok := adjacencyMap.edges[src]; !ok {
			adjacencyMap.edges[src] = make(map[int]mapItem)
		}

		adjacencyMap.nodes[src].outbound += float64(wt)
		adjacencyMap.edges[src][dest] = mapItem{float64(wt)}
	}

	fmt.Println("Adjacency List: ")
	graph.printGraph()

	Rank(adjacencyMap, 0.85, 0.000001, func(node int, rank float64) {
		fmt.Println("Node", node, "has a rank of", rank)
	})

}
