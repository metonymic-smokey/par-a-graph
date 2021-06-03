package main

import "testing"

const quoraEdgeFileName = "../datasets/quora_edges.txt"
const quoraNodeFileName = "../datasets/quora_nodes.txt"

func TestQuoraGraph(t *testing.T) {
	testHelperTopoPageRank(t, quoraEdgeFileName, quoraNodeFileName, 0.85, 0.000001)
}

func BenchmarkQuoraGraphE6(b *testing.B) {
	benchmarkHelperParallel(b, quoraEdgeFileName, quoraNodeFileName, 10e-6)
}

func BenchmarkQuoraGraphE9(b *testing.B) {
	benchmarkHelperParallel(b, quoraEdgeFileName, quoraNodeFileName, 10e-9)
}

func BenchmarkQuoraGraphE11(b *testing.B) {
	benchmarkHelperParallel(b, quoraEdgeFileName, quoraNodeFileName, 10e-11)
}

func BenchmarkQuoraGraphSerialE6(b *testing.B) {
	benchmarkHelperSerial(b, quoraEdgeFileName, quoraNodeFileName, 10e-6)
}

func BenchmarkQuoraGraphSerialE9(b *testing.B) {
	benchmarkHelperSerial(b, quoraEdgeFileName, quoraNodeFileName, 10e-9)
}

func BenchmarkQuoraGraphSerialE11(b *testing.B) {
	benchmarkHelperSerial(b, quoraEdgeFileName, quoraNodeFileName, 10e-11)
}
