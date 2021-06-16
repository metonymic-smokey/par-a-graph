package main

import (
	"math"
	"sync"
)

func pageRankSerial(
	vertexArray []int,
	edgeArray []int,
	outDegrees []int,
	alpha float64,
	eps float64) ([]float64, int) {

	n := len(vertexArray) - 1

	// pagerank vector
	x := make([]float64, n)
	new_x := make([]float64, n)

	for i := 0; i < n; i++ {
		x[i] = 1 / float64(n)
	}

	//all the nodes in 1 slice
	var nodes []int
	for i := 0; i < n; i++ {
		nodes = append(nodes, i)
	}

	delta := make([]float64, n)

	iters := 0

	for true {
		iters++

		deltaSum := 0.0
		var leak float64

		for _, v := range nodes {
			if outDegrees[v] == 0 { //dangling nodes have 0 outdegrees
				leak += x[v]
			}
		}

		leak *= alpha

		for _, v := range nodes {
			tmp := x[v]
			sum_value := 0.0
			for w := vertexArray[v]; w < vertexArray[v+1]; w++ {
				sum_value += x[edgeArray[w]] / float64(outDegrees[edgeArray[w]])
			}
			new_x[v] = (1-alpha)/float64(n) + alpha*sum_value + leak/float64(n)
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

	//first degree normalization
	for i := range x {
		x[i] /= norm
	}

	return x, iters

}

func pageRank(
	vertexArray []int,
	edgeArray []int,
	outDegrees []int,
	alpha float64,
	eps float64) ([]float64, int) {

	n := len(vertexArray) - 1

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

	// new idea
	// partition total nodes
	// give set number of nodes to each goroutine
	blockSize := n / numParallel

	// store deltaSum of each partition separately
	deltaSums := make([]float64, numParallel)

	// channels go brr
	signallers := make([]chan struct{}, numParallel)

	// to wait for initialization
	wg.Add(numParallel)

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
			for v := sliceStart; v < sliceEnd; v++ {
				x[v] = 1 / numNodes
			}
			wg.Done()

			for {
				<-signallers[parIndex]
				for v := sliceStart; v < sliceEnd; v++ {
					sumValue := 0.0
					for w := vertexArray[v]; w < vertexArray[v+1]; w++ {
						// could improve cache locality here using GAS, PCPM
						// lots of random access happening here
						sumValue += x[edgeArray[w]] / float64(outDegrees[edgeArray[w]])
					}
					new_x[v] = alphaTerm + alpha*sumValue + leak/numNodes
					deltaSums[parIndex] += math.Abs(new_x[v] - x[v])
					if outDegrees[v] == 0 {
						leaks[parIndex] += new_x[v]
					}
				}
				wg.Done()
			}
		}(i, sliceStart, sliceEnd)
	}

	// wait for initialization to complete
	wg.Wait()

	leak = 0.0
	for v := 0; v < n; v++ {
		if outDegrees[v] == 0 { //dangling nodes have 0 outdegree
			leak += x[v]
		}
	}
	leak *= alpha

	//variable for number of iterations
	iters := 0

	for {
		iters++

		deltaSum := 0.0

		wg.Add(numParallel)
		for i := 0; i < numParallel; i++ {
			signallers[i] <- struct{}{}
		}
		wg.Wait()

		leak = 0.0
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

	//first degree normalization
	for i := range x {
		x[i] /= norm
	}

	return x, iters

}
