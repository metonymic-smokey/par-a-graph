package main

import (
    "math"
)

func topoPageRank(edges [][2]int, pages [][2]string, alpha float64, adj_array map[int][]int, node_to_index map[int]int) []float64 {

	n := len(pages)
	//e := len(edges)

	// pagerank vector
	var x []float64

	for i := 0; i < n; i++ {
		x = append(x, 1/float64(n))
	}

	// error between iterations
	eps := 0.000001

	//
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

	for node, _ := range adj_array {
		out_neighbours := adj_array[node]
		for _, out_node := range out_neighbours {
			if _, ok := s[out_node]; !ok {
				s[out_node] = make([]int, 0)
			}
			s[out_node] = append(s[out_node], node)
		}
	}

	//max_delta := 0.0
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
