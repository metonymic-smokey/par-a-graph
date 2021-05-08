package main

import (
	"math"
	"sync"
)

func topoPageRankSerial(edges [][2]int, pages [][2]string, alpha float64, eps float64, adj_array map[int][]int, node_to_index map[int]int) []float64 {

	n := len(pages)
	//e := len(edges)

	// pagerank vector
	x := make([]float64, n)
	new_x := make([]float64, n)

	for i := 0; i < n; i++ {
		x[i] = 1 / float64(n)
	}

	//all the nodes in 1 slice
	var nodes []int

	for _, v := range node_to_index {
		nodes = append(nodes, v)
	}

	// out degree of each node
	degree_out := make([]float64, n)

	for i, nodes := range adj_array {
		degree_out[i] = float64(len(nodes))
	}

	//t := adj_array
	// node -> list of nodes connecting it
	s := make(map[int][]int)

	for node := range adj_array {
		out_neighbours := adj_array[node]
		for _, out_node := range out_neighbours {
			if _, ok := s[out_node]; !ok {
				s[out_node] = make([]int, 0)
			}
			s[out_node] = append(s[out_node], node)
		}
	}

	delta := make([]float64, n)

	for true {

		deltaSum := 0.0
		var leak float64

		for _, v := range nodes {
			if len(s[v]) == 0 { //dangling nodes
				leak += x[v]
			}
		}

		leak *= alpha

		for _, v := range nodes {
			tmp := x[v]
			sum_value := 0.0
			if _, ok := s[v]; ok {

				for _, w := range s[v] {
					sum_value += x[w] / degree_out[w]
				}
			}
			new_x[v] = (1-alpha)/float64(len(nodes)) + alpha*sum_value + leak/float64(len(nodes))
			delta[v] = math.Abs(new_x[v] - tmp)
			deltaSum += delta[v]
		}

		for i, new_val := range new_x {
			x[i] = new_val
		}

		if deltaSum < eps {
			break
		}
	}

	norm := 0.0
	for _, v := range x {
		norm += v
	}

	for i := range x {
		x[i] /= norm
	}

	return x

}

func topoPageRank(edges [][2]int, pages [][2]string, alpha float64, eps float64, adj_array map[int][]int, node_to_index map[int]int) []float64 {

	n := len(pages)
	// e := len(edges)

	// pagerank vector
	x := make([]float64, n)
	new_x := make([]float64, n)

	numParallel := 16
	if n < numParallel {
		numParallel = n
	}
	var wg sync.WaitGroup
	var leak float64
	leaks := make([]float64, numParallel)

	numNodes := float64(n)
	alphaTerm := (1 - alpha) / (numNodes)

	// out degree of each node
	degree_out := make([]float64, n)

	//t := adj_array
	// node -> list of nodes connecting it
	s := make(map[int][]int)

	for node := range adj_array {
		out_neighbours := adj_array[node]
		for _, out_node := range out_neighbours {
			if _, ok := s[out_node]; !ok {
				s[out_node] = make([]int, 0)
			}
			s[out_node] = append(s[out_node], node)
		}
	}

	// new idea
	// partition total nodes
	// give set number of nodes to each goroutine
	blockSize := n / numParallel

	// store deltaSum of each partition separately
	deltaSums := make([]float64, numParallel)

	// channels go brr
	signallers := make([]chan struct{}, numParallel)

	// to wait for initialization
	wg.Add(numParallel);
	for i := 0; i < numParallel; i++ {
		signallers[i] = make(chan struct{})

		var sliceStart int = blockSize * i
		var sliceEnd int

		if i == numParallel-1 {
			sliceEnd = n
		} else {
			sliceEnd = blockSize * (i + 1)
		}

		go func(parIndex int, sliceStart int, sliceEnd int) {
			// initialize pageranks
			// initialize out degree
			for v := sliceStart; v < sliceEnd; v++ {
				x[v] = 1 / float64(n)
				degree_out[v] = float64(len(adj_array[v]))
			}
			wg.Done()

			for {
				<-signallers[parIndex]
				for v := sliceStart; v < sliceEnd; v++ {
					sumValue := 0.0
					if len(s[v]) == 0 {
						leaks[parIndex] += x[v]
					} else {
						for _, w := range s[v] {
							// could improve cache locality here using GAS, PCPM
							sumValue += x[w] / degree_out[w]
						}
					}
					new_x[v] = alphaTerm + alpha*sumValue + leak/numNodes
					deltaSums[parIndex] += math.Abs(new_x[v] - x[v])
				}
				wg.Done()
			}
		}(i, sliceStart, sliceEnd)
	}

	// wait for initialization to complete
	wg.Wait()

	leak = 0.0
	for v := 0; v < n; v++ {
		if len(s[v]) == 0 { //dangling nodes
			leak += x[v]
		}
	}
	leak *= alpha

	for {
		deltaSum := 0.0

		wg.Add(numParallel)
		for i := 0; i < numParallel; i++ {
			signallers[i] <- struct{}{}
		}
		wg.Wait()

		leak = 0
		for i := 0; i < numParallel; i++ {
			leak += leaks[i]
			leaks[i] = 0.0

			deltaSum += deltaSums[i]
			deltaSums[i] = 0.0
		}
		leak *= alpha

		// swap x and new_x instead of replacing values
		temp := x
		x = new_x
		new_x = temp

		if deltaSum < eps {
			break
		}
	}

	norm := 0.0
	for _, v := range x {
		norm += v
	}

	for i := range x {
		x[i] /= norm
	}

	return x

}
