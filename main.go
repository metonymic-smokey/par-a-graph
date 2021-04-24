package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func readGraph() ([][2]int, [][2]string, map[int]int) {
	f, err := os.Open("dirLinks.txt")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var dirLinks [][2]int
	var fileContents []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		fileContents = append(fileContents, scanner.Text())
	}

	for _, fc := range fileContents {
		res := strings.Split(fc, ",")
		if len(res) < 2 {
			continue
		}
		src, _ := strconv.Atoi(res[0])
		dest, _ := strconv.Atoi(res[1])
		dirLinks = append(dirLinks, [2]int{src, dest})
	}

	f, err = os.Open("pageNum.txt")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var pageNum [][2]string
	var fc2 []string
	scanner = bufio.NewScanner(f)

	for scanner.Scan() {
		fc2 = append(fc2, scanner.Text())
	}

	for _, fc := range fc2 {
		res := strings.Split(fc, ",")
		if len(res) < 2 {
			continue
		}
		pageNum = append(pageNum, [2]string{res[0], res[1]})
	}

	fmt.Println("No. of edges: ", len(dirLinks))
	fmt.Println("No. of nodes: ", len(pageNum))
	//fmt.Println("No. of nodes: ", 2080370)

	v := 0
	node_to_index := make(map[int]int)

	for _, node := range pageNum {
		temp, _ := strconv.Atoi(node[0])
		node_to_index[temp] = v
		v += 1
	}

	for _, edge := range dirLinks {
		edge[0] = node_to_index[edge[0]]
		edge[1] = node_to_index[edge[1]]
	}

	return dirLinks, pageNum, node_to_index
}

func makeAdjArray(edges [][2]int) map[int][]int {

	adj_array := make(map[int][]int)
	for _, edge := range edges {
		_, ok := adj_array[edge[0]]
		if !ok {
			adj_array[edge[0]] = make([]int, 0)
		}
		adj_array[edge[0]] = append(adj_array[edge[0]], edge[1])
	}

	return adj_array

}

func max(arr []float64) float64 {

	res := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] > res {
			res = arr[i]
		}
	}

	return res
}

func topoPageRank(edges [][2]int, pages [][2]string, alpha float64, adj_array map[int][]int, node_to_index map[int]int) []float64 {

	n := len(pages)
	//e := len(edges)

	// pagerank vector
	var x []float64
	for i := 0; i < n; i++ {
		x = append(x, 1-alpha)
	}
	// error between iterations
	eps := 0.000001

	// 
	var nodes []int
	for _, v := range node_to_index {
		nodes = append(nodes, v)
	}

	fmt.Println(len(adj_array), n)

	// out degree of each node
	degree_out := make([]float64, n)
	for node, _ := range adj_array {
		degree_out[node] = float64(len(adj_array[node]))
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

	max_delta := 0.0

	for true {
		for _, v := range nodes {
			tmp := x
			sum_value := 0.0
			if _, ok := s[v]; ok {
				for _, w := range s[v] {
					sum_value += x[w] / degree_out[w]
				}
			}
			x[v] = alpha*sum_value + (1 - alpha)
			max_delta = math.Max(max_delta, x[v]-tmp[v])
		}
		if max_delta < eps {
			break
		}
	}

	return x

}

func main() {

	edges, pages, node_to_index := readGraph()
	adj_array := makeAdjArray(edges)
	pageRank := topoPageRank(edges, pages, 0.85, adj_array, node_to_index)
	for i := 0; i < len(pageRank); i++ {
		fmt.Println(pageRank[i])
	}
}
