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

const smallGraphEdgeFile = "../datasets/example"
const smallGraphNodeFile = "../datasets/examplePageNum"

const largeGraphEdgeFile = "../datasets/dirLinks.txt"
const largeGraphNodeFile = "../datasets/pageNum.txt"

// test helper
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
			fmt.Println("less than 2 in edge file")
			continue
		}
		src, _ := strconv.Atoi(res[0])
		dest, _ := strconv.Atoi(res[1])
		graph.Link(src, dest)
	}

	observed := make(map[int]float64)
	observedSerial := make(map[int]float64)
	expected := make(map[int]float64)

	graph.Rank(alpha, eps, func(identifier int, rank float64) {
		expected[identifier] = rank
	})

	edges, pages, node_to_index := readGraph(edgeFileName, nodeFileName)
	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))

	pageRankSerial, itersSerial := pageRankSerial(vertexArray, edgeArray, outDegrees, alpha, eps)
	for node, index := range node_to_index {
		observedSerial[node] = pageRankSerial[index]
	}

	pageRank, itersPar := pageRank(vertexArray, edgeArray, outDegrees, alpha, eps)
	for node, index := range node_to_index {
		observed[node] = pageRank[index]
	}

	threshold := 10e-7
	thresholdSerialParallel := 10e-7

	for node := range node_to_index {
		diff := math.Abs(observed[node] - expected[node])
		if diff > threshold {
			t.Errorf("Parallel vs expected: Page rank not matching for node %d; expected: %e, parallel: %e", node, expected[node], observed[node])
		}

		diff = math.Abs(observedSerial[node] - expected[node])
		if diff > threshold {
			t.Errorf("Serial vs expected: Page rank not matching for node %d; expected: %e, serial: %e", node, expected[node], observedSerial[node])
		}

		diff = math.Abs(observedSerial[node] - observed[node])
		if diff > thresholdSerialParallel {
			t.Errorf("Serial vs Parallel: Page rank not matching for node %d; parallel: %e, serial: %e", node, observed[node], observedSerial[node])
		}
	}

	if itersSerial != itersPar {
		t.Errorf("Number of iterations in serial and parallel do not match. Serial: %v, Parallel: %v", itersSerial, itersPar)
	}
}

// benchmark helpers

func benchmarkHelperSerial(b *testing.B, edgeFileName string, nodeFileName string, eps float64) {
	enableLog = false

	edges, pages, _ := readGraph(edgeFileName, nodeFileName)
	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))

	alpha := 0.85

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pageRankSerial(vertexArray, edgeArray, outDegrees, alpha, eps)
	}
}

func benchmarkHelperParallel(b *testing.B, edgeFileName string, nodeFileName string, eps float64) {
	enableLog = false

	edges, pages, _ := readGraph(edgeFileName, nodeFileName)
	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))

	alpha := 0.85

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pageRank(vertexArray, edgeArray, outDegrees, alpha, eps)
	}
}

// tests - small and large

func TestSmallGraph(t *testing.T) {
	testHelperTopoPageRank(t, smallGraphEdgeFile, smallGraphNodeFile, 0.85, 0.000001)
}

func TestLargeGraph(t *testing.T) {
	testHelperTopoPageRank(t, largeGraphEdgeFile, largeGraphNodeFile, 0.85, 0.000001)
}

// Large graph - benchmarks

func BenchmarkLargeGraphSerialE6(b *testing.B) {
	benchmarkHelperSerial(b, largeGraphEdgeFile, largeGraphNodeFile, 10e-6)
}

func BenchmarkLargeGraphE6(b *testing.B) {
	benchmarkHelperParallel(b, largeGraphEdgeFile, largeGraphNodeFile, 10e-6)
}

func BenchmarkLargeGraphSerialE9(b *testing.B) {
	benchmarkHelperSerial(b, largeGraphEdgeFile, largeGraphNodeFile, 10e-9)
}

func BenchmarkLargeGraphE9(b *testing.B) {
	benchmarkHelperParallel(b, largeGraphEdgeFile, largeGraphNodeFile, 10e-9)
}

func BenchmarkLargeGraphSerialE11(b *testing.B) {
	benchmarkHelperSerial(b, largeGraphEdgeFile, largeGraphNodeFile, 10e-11)
}

func BenchmarkLargeGraphE11(b *testing.B) {
	benchmarkHelperParallel(b, largeGraphEdgeFile, largeGraphNodeFile, 10e-11)
}

// Small graph - benchmarks

func BenchmarkSmallGraphE6(b *testing.B) {
	benchmarkHelperParallel(b, smallGraphEdgeFile, smallGraphNodeFile, 10e-6)
}

func BenchmarkSmallGraphE9(b *testing.B) {
	benchmarkHelperParallel(b, smallGraphEdgeFile, smallGraphNodeFile, 10e-9)
}

func BenchmarkSmallGraphE11(b *testing.B) {
	benchmarkHelperParallel(b, smallGraphEdgeFile, smallGraphNodeFile, 10e-11)
}

func BenchmarkSmallGraphSerialE6(b *testing.B) {
	benchmarkHelperSerial(b, smallGraphEdgeFile, smallGraphNodeFile, 10e-6)
}

func BenchmarkSmallGraphSerialE9(b *testing.B) {
	benchmarkHelperSerial(b, smallGraphEdgeFile, smallGraphNodeFile, 10e-9)
}

func BenchmarkSmallGraphSerialE11(b *testing.B) {
	benchmarkHelperSerial(b, smallGraphEdgeFile, smallGraphNodeFile, 10e-11)
}
