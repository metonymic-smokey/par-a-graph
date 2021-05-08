package main

import "testing"

func TestQuoraGraph(t *testing.T) {
	testHelperTopoPageRank(t, "./quora_edges.txt", "./quora_nodes.txt", 0.85, 0.000001)
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

func BenchmarkQuoraGraphSerialE6(b *testing.B) {
	edgeFileName := "./quora_edges.txt"
	nodeFileName := "./quora_nodes.txt"
	enableLog = false

	edges, pages, _ := readGraph(edgeFileName, nodeFileName)
	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))

	alpha := 0.85
	eps := 10e-6

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRankSerial(vertexArray, edgeArray, outDegrees, alpha, eps)
	}
}

func BenchmarkQuoraGraphSerialE9(b *testing.B) {
	edgeFileName := "./quora_edges.txt"
	nodeFileName := "./quora_nodes.txt"
	enableLog = false

	edges, pages, _ := readGraph(edgeFileName, nodeFileName)
	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))

	alpha := 0.85
	eps := 10e-9

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRankSerial(vertexArray, edgeArray, outDegrees, alpha, eps)
	}
}

func BenchmarkQuoraGraphSerialE11(b *testing.B) {
	edgeFileName := "./quora_edges.txt"
	nodeFileName := "./quora_nodes.txt"
	enableLog = false

	edges, pages, _ := readGraph(edgeFileName, nodeFileName)
	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))

	alpha := 0.85
	eps := 10e-11

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		topoPageRankSerial(vertexArray, edgeArray, outDegrees, alpha, eps)
	}
}
