package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/dcadenas/pagerank"
)

func TestTopoPageRank(t *testing.T) {

	f, err := os.Open("dirLinks.txt")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	graph := pagerank.New()

	var fileContents []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		fileContents = append(fileContents, scanner.Text())
	}

	for _, fc := range fileContents {
		res := strings.Split(fc, ",")
		if len(res) < 2 {
			fmt.Println("less than 2 in dirLinks")
			continue
		}
		src, _ := strconv.Atoi(res[0])
		dest, _ := strconv.Atoi(res[1])
		graph.Link(src, dest)
	}

	alpha := 0.85
	eps := 0.000001

	observed := make(map[int]float64)
	expected := make(map[int]float64)

	graph.Rank(alpha, eps, func(identifier int, rank float64) {
		expected[identifier] = rank
	})

	edges, pages, node_to_index := readGraph()
	adj_array := makeAdjArray(edges, len(pages))
	pageRank := topoPageRank(edges, pages, alpha, adj_array, node_to_index)
	for node, index := range node_to_index {
		observed[node] = pageRank[index]
	}

	for node, _ := range node_to_index {
        diff := math.Abs(observed[node] - expected[node])
        fmt.Println("DIFF: ", diff)
		if diff > 10e-6{
			t.Errorf("Page rank not matching for node %d; expected: %e, observed: %e", node, expected[node], observed[node])
		}
	}
}
