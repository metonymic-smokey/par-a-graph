package main

import (
	"fmt"
	"testing"
	"time"
)

// stolen from stackoverflow: https://stackoverflow.com/a/15312097/11199009
func testEq(a, b []int) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func testPathEqual(adjacencyMap map[int]map[int]mapItem, srcNode int, dstNode int, want []int, t *testing.T) {
	start := time.Now()
	got := dgraphShortest(adjacencyMap, srcNode, dstNode)
  fmt.Printf("took %v\n", time.Since(start))

	if !testEq(got, want) {
		t.Errorf("Got wrong path from dgraphShortest. Got: %v, Want: %v", got, want)
	}
}

func TestOurFirstGraph(t *testing.T) {
	adjacencyMap := make(map[int]map[int]mapItem)
	adjacencyMap[5] = map[int]mapItem{1: {4}, 2: {2}}
	adjacencyMap[1] = map[int]mapItem{2: {1}, 3: {5}}
	adjacencyMap[2] = map[int]mapItem{3: {8}, 4: {10}}
	adjacencyMap[3] = map[int]mapItem{4: {2}, 6: {6}}
	adjacencyMap[4] = map[int]mapItem{6: {2}}

	testPathEqual(adjacencyMap, 5, 6, []int{5, 1, 3, 4, 6}, t)
	testPathEqual(adjacencyMap, 2, 6, []int{2, 4, 6}, t)
	testPathEqual(adjacencyMap, 1, 2, []int{1, 2}, t)
	testPathEqual(adjacencyMap, 3, 6, []int{3, 4, 6}, t)
}
