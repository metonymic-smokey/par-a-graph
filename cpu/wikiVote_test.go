package main

import "testing"

const wikiEdgeFileName = "../datasets/wiki-vote-edges.txt"
const wikiNodeFileName = "../datasets/wiki-vote-nodes.txt"

func TestWikiVoteGraph(t *testing.T) {
	testHelperTopoPageRank(t, wikiEdgeFileName, wikiNodeFileName, 0.85, 10e-6)
}

func BenchmarkWikiVoteGraphE6(b *testing.B) {
	benchmarkHelperParallel(b, wikiEdgeFileName, wikiNodeFileName, 10e-6)
}

func BenchmarkWikiVoteGraphE9(b *testing.B) {
	benchmarkHelperParallel(b, wikiEdgeFileName, wikiNodeFileName, 10e-9)
}

func BenchmarkWikiVoteGraphE11(b *testing.B) {
	benchmarkHelperParallel(b, wikiEdgeFileName, wikiNodeFileName, 10e-11)
}

func BenchmarkWikiVoteGraphSerialE6(b *testing.B) {
	benchmarkHelperSerial(b, wikiEdgeFileName, wikiNodeFileName, 10e-6)
}

func BenchmarkWikiVoteGraphSerialE9(b *testing.B) {
	benchmarkHelperSerial(b, wikiEdgeFileName, wikiNodeFileName, 10e-9)
}

func BenchmarkWikiVoteGraphSerialE11(b *testing.B) {
	benchmarkHelperSerial(b, wikiEdgeFileName, wikiNodeFileName, 10e-11)
}
