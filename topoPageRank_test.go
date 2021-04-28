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

func testHelperTopoPageRank(t *testing.T, edgeFileName string, nodeFileName string, alpha float64, eps float64) {

	f, err := os.Open(edgeFileName)
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

	observed := make(map[int]float64)
	expected := make(map[int]float64)

	graph.Rank(alpha, eps, func(identifier int, rank float64) {
		expected[identifier] = rank
	})

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	adj_array := makeAdjArray(edges, len(pages))
	pageRank := topoPageRank(edges, pages, alpha, eps, adj_array, node_to_index)
	for node, index := range node_to_index {
		observed[node] = pageRank[index]
	}

	for node := range node_to_index {
		diff := math.Abs(observed[node] - expected[node])
		if diff > 10e-7 {
			t.Errorf("Page rank not matching for node %d; expected: %e, observed: %e", node, expected[node], observed[node])
		}
	}
}

func TestSmallGraph(t *testing.T) {
	testHelperTopoPageRank(t, "./example", "./examplePageNum", 0.85, 0.000001)
}

func TestQuoraGraph(t *testing.T) {
	testHelperTopoPageRank(t, "./quora_edges.txt", "./quora_nodes.txt", 0.85, 0.000001)
}

func TestLargeGraph(t *testing.T) {
	testHelperTopoPageRank(t, "./dirLinks.txt", "./pageNum.txt", 0.85, 0.000001)
}

func BenchmarkLargeGraphSerialE6(b *testing.B) {
	edgeFileName := "./dirLinks.txt"
	nodeFileName := "./pageNum.txt"
	enableLog = false

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	adj_array := makeAdjArray(edges, len(pages))

	alpha := 0.85
	eps := 10e-6

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRankSerial(edges, pages, alpha, eps, adj_array, node_to_index)
	}
}

func BenchmarkLargeGraphE6(b *testing.B) {
	edgeFileName := "./dirLinks.txt"
	nodeFileName := "./pageNum.txt"
	enableLog = false

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	adj_array := makeAdjArray(edges, len(pages))

	alpha := 0.85
	eps := 10e-6

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRank(edges, pages, alpha, eps, adj_array, node_to_index)
	}
}

func BenchmarkSmallGraphE6(b *testing.B) {
	edgeFileName := "./example"
	nodeFileName := "./examplePageNum"
	enableLog = false

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	adj_array := makeAdjArray(edges, len(pages))

	alpha := 0.85
	eps := 10e-6

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRank(edges, pages, alpha, eps, adj_array, node_to_index)
	}
}

func BenchmarkLargeGraphSerialE9(b *testing.B) {
	edgeFileName := "./dirLinks.txt"
	nodeFileName := "./pageNum.txt"
	enableLog = false

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	adj_array := makeAdjArray(edges, len(pages))

	alpha := 0.85
	eps := 10e-9

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRankSerial(edges, pages, alpha, eps, adj_array, node_to_index)
	}
}

func BenchmarkLargeGraphE9(b *testing.B) {
	edgeFileName := "./dirLinks.txt"
	nodeFileName := "./pageNum.txt"

	enableLog = false

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	adj_array := makeAdjArray(edges, len(pages))

	alpha := 0.85
	eps := 10e-9

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRank(edges, pages, alpha, eps, adj_array, node_to_index)
	}
}

func BenchmarkSmallGraphE9(b *testing.B) {
	edgeFileName := "./example"
	nodeFileName := "./examplePageNum"
	enableLog = false

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	adj_array := makeAdjArray(edges, len(pages))

	alpha := 0.85
	eps := 10e-9

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRank(edges, pages, alpha, eps, adj_array, node_to_index)
	}
}

func BenchmarkLargeGraphSerialE11(b *testing.B) {
	edgeFileName := "./dirLinks.txt"
	nodeFileName := "./pageNum.txt"
	enableLog = false

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	adj_array := makeAdjArray(edges, len(pages))

	alpha := 0.85
	eps := 10e-11

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRankSerial(edges, pages, alpha, eps, adj_array, node_to_index)
	}
}

func BenchmarkLargeGraphE11(b *testing.B) {
	edgeFileName := "./dirLinks.txt"
	nodeFileName := "./pageNum.txt"
	enableLog = false

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	adj_array := makeAdjArray(edges, len(pages))

	alpha := 0.85
	eps := 10e-11

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRank(edges, pages, alpha, eps, adj_array, node_to_index)
	}
}

func BenchmarkSmallGraphE11(b *testing.B) {
	edgeFileName := "./example"
	nodeFileName := "./examplePageNum"
	enableLog = false

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	adj_array := makeAdjArray(edges, len(pages))

	alpha := 0.85
	eps := 10e-11

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRank(edges, pages, alpha, eps, adj_array, node_to_index)
	}
}

func BenchmarkQuoraGraphE6(b *testing.B) {
	edgeFileName := "./quora_edges.txt"
	nodeFileName := "./quora_nodes.txt"
	enableLog = false

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	adj_array := makeAdjArray(edges, len(pages))

	alpha := 0.85
	eps := 10e-6

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRank(edges, pages, alpha, eps, adj_array, node_to_index)
	}
}

func BenchmarkQuoraGraphE9(b *testing.B) {
	edgeFileName := "./quora_edges.txt"
	nodeFileName := "./quora_nodes.txt"
	enableLog = false

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	adj_array := makeAdjArray(edges, len(pages))

	alpha := 0.85
	eps := 10e-9

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRank(edges, pages, alpha, eps, adj_array, node_to_index)
	}
}

func BenchmarkQuoraGraphE11(b *testing.B) {
	edgeFileName := "./quora_edges.txt"
	nodeFileName := "./quora_nodes.txt"
	enableLog = false

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	adj_array := makeAdjArray(edges, len(pages))

	alpha := 0.85
	eps := 10e-11

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRank(edges, pages, alpha, eps, adj_array, node_to_index)
	}
}
