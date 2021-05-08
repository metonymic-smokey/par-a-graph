package main

import "testing"

func TestWikiVoteGraph(t *testing.T) {
	testHelperTopoPageRank(t, "./wiki-vote-edges.txt", "./wiki-vote-nodes.txt", 0.85, 0.000001)
}

func BenchmarkWikiVoteGraphE6(b *testing.B) {
	edgeFileName := "./wiki-vote-edges.txt"
	nodeFileName := "./wiki-vote-nodes.txt"
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

func BenchmarkWikiVoteGraphE9(b *testing.B) {
	edgeFileName := "./wiki-vote-edges.txt"
	nodeFileName := "./wiki-vote-nodes.txt"
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

func BenchmarkWikiVoteGraphE11(b *testing.B) {
	edgeFileName := "./wiki-vote-edges.txt"
	nodeFileName := "./wiki-vote-nodes.txt"
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

func BenchmarkWikiVoteGraphSerialE6(b *testing.B) {
	edgeFileName := "./wiki-vote-edges.txt"
	nodeFileName := "./wiki-vote-nodes.txt"
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

func BenchmarkWikiVoteGraphSerialE9(b *testing.B) {
	edgeFileName := "./wiki-vote-edges.txt"
	nodeFileName := "./wiki-vote-nodes.txt"
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

func BenchmarkWikiVoteGraphSerialE11(b *testing.B) {
	edgeFileName := "./wiki-vote-edges.txt"
	nodeFileName := "./wiki-vote-nodes.txt"
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
