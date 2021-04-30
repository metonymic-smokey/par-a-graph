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

	for true {

		deltaSum := 0.0
		var leak float64

		for _, v := range nodes {
			if len(s[v]) == 0 { //dangling nodes
				leak += x[v]
			}
		}

		leak *= alpha

		var wg sync.WaitGroup
		wg.Add(len(nodes))

		for _, v := range nodes {
			go func(v int) {
				defer wg.Done()
				tmp := x[v]
				sum_value := 0.0
				if _, ok := s[v]; ok {

					for _, w := range s[v] {
						sum_value += x[w] / degree_out[w]
					}
				}
				x[v] = (1-alpha)/float64(len(nodes)) + alpha*sum_value + leak/float64(len(nodes))
				delta[v] = math.Abs(x[v] - tmp)
			}(v)
		}
		wg.Wait()

		for _, val := range delta {
			deltaSum += val
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
