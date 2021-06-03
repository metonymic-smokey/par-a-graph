package main

import "testing"

const wikiEdgeFileName = "../datasets/wiki-vote-edges.txt"
const wikiNodeFileName = "../datasets/wiki-vote-nodes.txt"

func TestWikiVoteGraph(t *testing.T) {
	testHelperTopoPageRank(t, wikiEdgeFileName, wikiNodeFileName, 0.85, 10e-6)
}

func BenchmarkWikiVoteGraphE6(b *testing.B) {
	edgeFileName := wikiEdgeFileName
	nodeFileName := wikiNodeFileName
	enableLog = false

	edges, pages, _ := readGraph(edgeFileName, nodeFileName)
	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))

	alpha := 0.85
	eps := 10e-6

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pageRank(vertexArray, edgeArray, outDegrees, alpha, eps)
	}
}

func BenchmarkWikiVoteGraphE9(b *testing.B) {
	edgeFileName := wikiEdgeFileName
	nodeFileName := wikiNodeFileName
	enableLog = false

	edges, pages, _ := readGraph(edgeFileName, nodeFileName)
	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))

	alpha := 0.85
	eps := 10e-9

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pageRank(vertexArray, edgeArray, outDegrees, alpha, eps)
	}
}

func BenchmarkWikiVoteGraphE11(b *testing.B) {
	edgeFileName := wikiEdgeFileName
	nodeFileName := wikiNodeFileName
	enableLog = false

	edges, pages, _ := readGraph(edgeFileName, nodeFileName)
	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))

	alpha := 0.85
	eps := 10e-11

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pageRank(vertexArray, edgeArray, outDegrees, alpha, eps)
	}
}

func BenchmarkWikiVoteGraphSerialE6(b *testing.B) {
	edgeFileName := wikiEdgeFileName
	nodeFileName := wikiNodeFileName
	enableLog = false

	edges, pages, _ := readGraph(edgeFileName, nodeFileName)
	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))

	alpha := 0.85
	eps := 10e-6

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pageRankSerial(vertexArray, edgeArray, outDegrees, alpha, eps)
	}
}

func BenchmarkWikiVoteGraphSerialE9(b *testing.B) {
	edgeFileName := wikiEdgeFileName
	nodeFileName := wikiNodeFileName
	enableLog = false

	edges, pages, _ := readGraph(edgeFileName, nodeFileName)
	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))

	alpha := 0.85
	eps := 10e-9

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pageRankSerial(vertexArray, edgeArray, outDegrees, alpha, eps)
	}
}

func BenchmarkWikiVoteGraphSerialE11(b *testing.B) {
	edgeFileName := wikiEdgeFileName
	nodeFileName := wikiNodeFileName
	enableLog = false

	edges, pages, _ := readGraph(edgeFileName, nodeFileName)
	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))

	alpha := 0.85
	eps := 10e-11

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pageRankSerial(vertexArray, edgeArray, outDegrees, alpha, eps)
	}
}
