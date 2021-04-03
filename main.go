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

	/*
		    //number of nodes
		    n, err := strconv.Atoi(fileContents[0])
			if err != nil {
				panic(err)
			}
	*/

	for _, fc := range fileContents[1:] {
		res := strings.SplitN(fc, ",", -1)
		src, _ := strconv.Atoi(res[0])
		dest, _ := strconv.Atoi(res[1])
		wt, _ := strconv.Atoi(res[2])
		graph.addEdge(int64(src), int64(dest), int64(wt))
	}

	fmt.Println("Adjacency List: ")
	graph.printGraph()

	fmt.Println("Dijkstra's result: ")
	fmt.Println(graph.getPath(5, 6))
}
