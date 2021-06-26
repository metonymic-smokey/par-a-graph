package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var enableLog = true

func readGraph(edgeFileName string, nodeFileName string) ([][2]int, [][2]string, map[int]int) {
	f, err := os.Open(edgeFileName)
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
			fmt.Println("less than 2 in dirLinks")
			continue
		}
		src, _ := strconv.Atoi(res[0])
		dest, _ := strconv.Atoi(res[1])
		dirLinks = append(dirLinks, [2]int{src, dest})
	}

	f, err = os.Open(nodeFileName)
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
			fmt.Println("less than 2 in pageNum")
			continue
		}
		pageNum = append(pageNum, [2]string{res[0], res[1]})
	}

	if enableLog {
		fmt.Println("No. of edges: ", len(dirLinks))
		fmt.Println("No. of nodes: ", len(pageNum))
	}
	//fmt.Println("No. of nodes: ", 2080370)

	v := 0
	node_to_index := make(map[int]int)

	for _, node := range pageNum {
		temp, _ := strconv.Atoi(node[0])
		node_to_index[temp] = v
		v += 1
	}

	for i, edge := range dirLinks {
		if _, ok := node_to_index[edge[0]]; !ok {
			fmt.Printf("node %v does not exist in node_to_index", edge[0])
		}
		if _, ok := node_to_index[edge[1]]; !ok {
			fmt.Printf("node %v does not exist in node_to_index", edge[1])
		}
		dirLinks[i][0] = node_to_index[edge[0]]
		dirLinks[i][1] = node_to_index[edge[1]]
	}

	return dirLinks, pageNum, node_to_index
}

func makeAdjArray(edges [][2]int, n int) map[int][]int {

	adj_array := make(map[int][]int)

	for i := 0; i < n; i++ {
		adj_array[i] = make([]int, 0)
	}

	for _, edge := range edges {
		adj_array[edge[0]] = append(adj_array[edge[0]], edge[1])
	}

	return adj_array

}

// converts the edge list to CSR format with
// a vertexArray consisting of cumulative in degrees
// an edgeArray consisting of the in-vertices
func makeCSR(edges [][2]int, n int) (vertexArray, edgeArray, outDegrees []int) {
	// stores indegrees of all vertices
	inDegrees := make([]int, n)

	// stores outdegrees of all vertices
	// needed in pagerank computation
	outDegrees = make([]int, n)

	// vertex array of CSR format
	vertexArray = make([]int, n+1)
	// edge or destination array
	edgeArray = make([]int, len(edges))

	// populate indegrees and outDegrees
	//  - increment inDegree of the destination vertex
	//  - increment out degree of the source vertex
	for i := 0; i < len(edges); i++ {
		inDegrees[edges[i][1]]++
		outDegrees[edges[i][0]]++
	}

	// populate vertex  array
	// with cumulative indegree
	//   previous value + indegree of previous node
	// first value is always 0
	for i := 1; i <= n; i++ {
		vertexArray[i] = vertexArray[i-1] + inDegrees[i-1]
	}

	// we will reuse inDegrees to store number of in-vertices
	//  currently processed for that node
	// this is used to index into edgeArray
	for i := 0; i < n; i++ {
		inDegrees[i] = 0
	}

	// populate edge array
	for i := 0; i < len(edges); i++ {
		// since we need incoming vertices, we store
		//   destination to source mapping
		from := edges[i][1]
		to := edges[i][0]
		edgeArray[vertexArray[from]+inDegrees[from]] = to

		inDegrees[from]++
	}

	// the vertices in edgeArray are currently randomly ordered
	// let's sort them to reduce randomm accesses
	for i := 0; i < n; i++ {
		// if at least 2 edges exist, sort them
		if vertexArray[i+1] > vertexArray[i]+1 {
			sort.Ints(edgeArray[vertexArray[i]:vertexArray[i+1]])
		}
	}

	return
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

func main() {
	edges, pages, _ := readGraph("./wiki-vote-edges.txt", "./wiki-vote-nodes.txt")
	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))
	// pageRankCU(vertexArray, edgeArray, outDegrees, 0.85, 10e-6)
	pageRank(vertexArray, edgeArray, outDegrees, 0.85, 10e-6)
}
