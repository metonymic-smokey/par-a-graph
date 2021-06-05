package main

import "testing"

const stanfordEdgeFileName = "../datasets/stanford-edges.txt"
const stanfordNodeFileName = "../datasets/stanford-nodes.txt"

func TestStanfordVoteGraph(t *testing.T) {
	testHelperTopoPageRank(t, stanfordEdgeFileName, stanfordNodeFileName, 0.85, 10e-6)
}

func BenchmarkStanfordVoteGraphE6(b *testing.B) {
	benchmarkHelperParallel(b, stanfordEdgeFileName, stanfordNodeFileName, 10e-6)
}

func BenchmarkStanfordVoteGraphE9(b *testing.B) {
	benchmarkHelperParallel(b, stanfordEdgeFileName, stanfordNodeFileName, 10e-9)
}

func BenchmarkStanfordVoteGraphE11(b *testing.B) {
	benchmarkHelperParallel(b, stanfordEdgeFileName, stanfordNodeFileName, 10e-11)
}

func BenchmarkStanfordVoteGraphSerialE6(b *testing.B) {
	benchmarkHelperSerial(b, stanfordEdgeFileName, stanfordNodeFileName, 10e-6)
}

func BenchmarkStanfordVoteGraphSerialE9(b *testing.B) {
	benchmarkHelperSerial(b, stanfordEdgeFileName, stanfordNodeFileName, 10e-9)
}

func BenchmarkStanfordVoteGraphSerialE11(b *testing.B) {
	benchmarkHelperSerial(b, stanfordEdgeFileName, stanfordNodeFileName, 10e-11)
}
