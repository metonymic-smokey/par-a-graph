package main

import "testing"

const quoraEdgeFileName = "../datasets/quora_edges.txt"
const quoraNodeFileName = "../datasets/quora_nodes.txt"

func TestQuoraGraph(t *testing.T) {
	testHelperTopoPageRank(t, quoraEdgeFileName, quoraNodeFileName, 0.85, 0.000001)
}

func BenchmarkQuoraGraphE6(b *testing.B) {
	edgeFileName := quoraEdgeFileName
	nodeFileName := quoraNodeFileName
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

func BenchmarkQuoraGraphE9(b *testing.B) {
	edgeFileName := quoraEdgeFileName
	nodeFileName := quoraNodeFileName
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

func BenchmarkQuoraGraphE11(b *testing.B) {
	edgeFileName := quoraEdgeFileName
	nodeFileName := quoraNodeFileName
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

func BenchmarkQuoraGraphSerialE6(b *testing.B) {
	edgeFileName := quoraEdgeFileName
	nodeFileName := quoraNodeFileName
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

func BenchmarkQuoraGraphSerialE9(b *testing.B) {
	edgeFileName := quoraEdgeFileName
	nodeFileName := quoraNodeFileName
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

func BenchmarkQuoraGraphSerialE11(b *testing.B) {
	edgeFileName := quoraEdgeFileName
	nodeFileName := quoraNodeFileName
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
