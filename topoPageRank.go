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
			x[v] = (1-alpha)/float64(len(nodes)) + alpha*sum_value + leak/float64(len(nodes))
			delta[v] = math.Abs(x[v] - tmp)
			deltaSum += delta[v]
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
	//e := len(edges)

	// pagerank vector
	x := make([]float64, n)
	new_x := make([]float64, n)

	var wgAssign sync.WaitGroup
	wgAssign.Add(n)

	for i := 0; i < n; i++ {
		go func(i int) {
			defer wgAssign.Done()
			x[i] = 1 / float64(n)
		}(i)
	}

	wgAssign.Wait()

	//all the nodes in 1 slice
	var nodes []int

	for _, v := range node_to_index {
		nodes = append(nodes, v)
	}

	// out degree of each node
	degree_out := make([]float64, n)

	wgAssign.Add(n)

	for i, nodes := range adj_array {
		go func(i int, nodes []int) {
			defer wgAssign.Done()
			degree_out[i] = float64(len(nodes))
		}(i, nodes)
	}

	wgAssign.Wait()

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

	numParallel := 8
	ch := make(chan int, 1000)
	var wg sync.WaitGroup
	var leak float64

	for i := 0; i < numParallel; i++ {
		go func() {
			for {
				v := <-ch
				tmp := x[v]
				sumValue := 0.0
				if _, ok := s[v]; ok {
					for _, w := range s[v] {
						sumValue += x[w] / degree_out[w]
					}
				}
				new_x[v] = (1-alpha)/float64(len(nodes)) + alpha*sumValue + leak/float64(len(nodes))
				delta[v] = math.Abs(new_x[v] - tmp)
				wg.Done()
			}
		}()
	}

	for {

		deltaSum := 0.0
		leak = 0.0

		for _, v := range nodes {
			if len(s[v]) == 0 { //dangling nodes
				leak += x[v]
			}
		}

		leak *= alpha

		wg.Add(len(nodes))

		for _, v := range nodes {
			ch <- v
		}
		wg.Wait()

		for _, val := range delta {
			deltaSum += val
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
