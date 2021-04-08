package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
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

	adjacencyMap := make(map[int]map[int]mapItem)

	for _, fc := range fileContents {
		res := strings.SplitN(fc, ",", -1)
		src, _ := strconv.Atoi(res[0])
		dest, _ := strconv.Atoi(res[1])
		wt, _ := strconv.Atoi(res[2])
		graph.addEdge(src, dest, wt)
		if _, ok := adjacencyMap[src]; !ok {
			adjacencyMap[src] = make(map[int]mapItem)
		}
		adjacencyMap[src][dest] = mapItem{float64(wt)}
	}

	fmt.Println("Adjacency List: ")
	graph.printGraph()

	srcNode := 3
	destNode := 6

	dgraphResult := dgraphShortest(adjacencyMap, srcNode, destNode)

	res, path := (graph.getPath(srcNode, destNode))

	if !reflect.DeepEqual(dgraphResult, path) {
		fmt.Println("Incorrect path")
		return
	}

	fmt.Println("Cost: ", res)
	fmt.Println("Path: ", path)
}
