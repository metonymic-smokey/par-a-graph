package main

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"fmt"
)

func main() {
	agraph := newGraph()

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

	for _, fc := range fileContents {
		res := strings.SplitN(fc, ",", -1)
		src, _ := strconv.Atoi(res[0])
		dest, _ := strconv.Atoi(res[1])
		wt, _ := strconv.Atoi(res[2])
		agraph.addEdge(src, dest, wt)
	}

	agraph.printGraph()

	graph := GraphFromAdjList(*agraph)
	// hack for that issue in GraphFromAdjList
	graph.ptr = append(graph.ptr, graph.ptr[len(graph.ptr) - 1])

	graph.PrintGraph()

	sourceId := uint64(4)
	shortestDistances := Dijkstra(&graph, sourceId)

	for i, v := range shortestDistances {
		fmt.Printf("From vertex %v to %v = %v\n", sourceId, i, v)
	}
}
