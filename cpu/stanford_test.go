package main

import "testing"

const stanfordEdgeFileName = "../datasets/stanford-edges.txt"
const stanfordNodeFileName = "../datasets/stanford-nodes.txt"

func TestStanfordVoteGraph(t *testing.T) {
	testHelperTopoPageRank(t, stanfordEdgeFileName, stanfordNodeFileName, 0.85, 10e-6)
}

func BenchmarkStanfordVoteGraphE6(b *testing.B) {
	edgeFileName := stanfordEdgeFileName
	nodeFileName := stanfordNodeFileName
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

func BenchmarkStanfordVoteGraphE9(b *testing.B) {
	edgeFileName := stanfordEdgeFileName
	nodeFileName := stanfordNodeFileName
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

func BenchmarkStanfordVoteGraphE11(b *testing.B) {
	edgeFileName := stanfordEdgeFileName
	nodeFileName := stanfordNodeFileName
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

func BenchmarkStanfordVoteGraphSerialE6(b *testing.B) {
	edgeFileName := stanfordEdgeFileName
	nodeFileName := stanfordNodeFileName
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

func BenchmarkStanfordVoteGraphSerialE9(b *testing.B) {
	edgeFileName := stanfordEdgeFileName
	nodeFileName := stanfordNodeFileName
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

func BenchmarkStanfordVoteGraphSerialE11(b *testing.B) {
	edgeFileName := stanfordEdgeFileName
	nodeFileName := stanfordNodeFileName
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
