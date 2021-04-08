package main

import (
	"math/rand"
)

func GenRandom(p float64, N int) adjMap {
	graph := make(adjMap)

	for i := 0; i < N; i++ {
		graph[i] = make(adjList)
	}

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			rnum := rand.Float64()

			if rnum < p {
				// graph[i] = append(graph[i], {j, })
				graph[i][j] = mapItem{rand.Float64()}
			}
		}
	}

	return graph
}
