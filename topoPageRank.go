package main

// /*
//    void pageRankCuda(int *vertexArray, int vertexArray_size, int *edgeArray,
//                      int edgeArray_size, int *outDegrees, int outDegree_size,
//                      double alpha, double eps);
//    #cgo LDFLAGS: -L. -lpagerank
// */
// import "C"
// import (
// 	"bufio"
// 	"fmt"
// 	"math"
// 	"os"
// 	"sort"
// 	"strconv"
// 	"strings"
// 	"sync"
// 	// "log"
// 	// "runtime"
// )

// func pageRankGPU(
// 	vertexArray []int,
// 	edgeArray []int,
// 	outDegrees []int,
// 	alpha float64,
// 	eps float64) {

// 	var vertexArray_size C.int
// 	var edgeArray_size C.int
// 	var outDegrees_size C.int

// 	vertexArray_size = C.int(len(vertexArray))
// 	edgeArray_size = C.int(len(edgeArray))
// 	outDegrees_size = C.int(len(outDegrees))

// 	_vertexArray := make([]C.int, vertexArray_size)
// 	_edgeArray := make([]C.int, edgeArray_size)

// 	_outDegrees := make([]C.int, outDegrees_size)
// 	fmt.Println(vertexArray[3])
// 	for i, v := range vertexArray {
// 		_vertexArray[i] = C.int(v)
// 	}
// 	for i, v := range edgeArray {
// 		_edgeArray[i] = C.int(v)
// 	}
// 	for i, v := range outDegrees {
// 		_outDegrees[i] = C.int(v)
// 	}

// 	C.pageRankCuda(&_vertexArray[0], vertexArray_size, &_edgeArray[0], edgeArray_size, &_outDegrees[0], outDegrees_size, C.double(alpha), C.double(eps))
// }

// func pageRankSerial(
// 	vertexArray []int,
// 	edgeArray []int,
// 	outDegrees []int,
// 	alpha float64,
// 	eps float64) ([]float64, int) {

// 	n := len(vertexArray) - 1

// 	// pagerank vector
// 	x := make([]float64, n)
// 	new_x := make([]float64, n)

// 	for i := 0; i < n; i++ {
// 		x[i] = 1 / float64(n)
// 	}

// 	//all the nodes in 1 slice
// 	var nodes []int
// 	for i := 0; i < n; i++ {
// 		nodes = append(nodes, i)
// 	}

// 	delta := make([]float64, n)

// 	iters := 0

// 	for {
// 		iters++

// 		deltaSum := 0.0
// 		var leak float64

// 		for _, v := range nodes {
// 			if outDegrees[v] == 0 { //dangling nodes
// 				leak += x[v]
// 				// log.Println("found leak", x[v], outDegrees[v], v)
// 			}
// 		}

// 		leak *= alpha

// 		for _, v := range nodes {
// 			tmp := x[v]
// 			sum_value := 0.0
// 			for w := vertexArray[v]; w < vertexArray[v+1]; w++ {
// 				sum_value += x[edgeArray[w]] / float64(outDegrees[edgeArray[w]])
// 			}
// 			new_x[v] = (1-alpha)/float64(n) + alpha*sum_value + leak/float64(n)
// 			delta[v] = math.Abs(new_x[v] - tmp)
// 			deltaSum += delta[v]
// 		}

// 		for i, new_val := range new_x {
// 			x[i] = new_val
// 		}
// 		// log.Println(deltaSum)

// 		if deltaSum < eps {
// 			break
// 		}
// 	}

// 	// log.Println("serial iters: ", iters)

// 	norm := 0.0
// 	for _, v := range x {
// 		norm += v
// 	}

// 	for i := range x {
// 		x[i] /= norm
// 	}

// 	return x, iters

// }

// func pageRank(
// 	vertexArray []int,
// 	edgeArray []int,
// 	outDegrees []int,
// 	alpha float64,
// 	eps float64) ([]float64, int) {

// 	n := len(vertexArray) - 1

// 	// pagerank vector
// 	x := make([]float64, n)
// 	new_x := make([]float64, n)

// 	numParallel := 16
// 	if n < numParallel {
// 		numParallel = n
// 	}
// 	var wg sync.WaitGroup
// 	var leak float64
// 	leaks := make([]float64, numParallel)

// 	numNodes := float64(n)
// 	alphaTerm := (1 - alpha) / (numNodes)

// 	// new idea
// 	// partition total nodes
// 	// give set number of nodes to each goroutine
// 	blockSize := n / numParallel

// 	// store deltaSum of each partition separately
// 	deltaSums := make([]float64, numParallel)

// 	// channels go brr
// 	signallers := make([]chan struct{}, numParallel)

// 	// to wait for initialization
// 	wg.Add(numParallel)

// 	for i := 0; i < numParallel; i++ {
// 		signallers[i] = make(chan struct{})

// 		var sliceStart int = blockSize * i
// 		var sliceEnd int

// 		if i == numParallel-1 {
// 			sliceEnd = n
// 		} else {
// 			sliceEnd = blockSize * (i + 1)
// 		}
// 		// log.Println("allocating block ", sliceStart, " ", sliceEnd)

// 		go func(parIndex int, sliceStart int, sliceEnd int) {
// 			// initialize pageranks
// 			for v := sliceStart; v < sliceEnd; v++ {
// 				x[v] = 1 / numNodes
// 			}
// 			wg.Done()

// 			// runtime.LockOSThread()
// 			for {
// 				<-signallers[parIndex]
// 				// _, ok := <-signallers[parIndex]
// 				// if !ok {
// 				// 	// runtime.UnlockOSThread()
// 				// 	return
// 				// }
// 				for v := sliceStart; v < sliceEnd; v++ {
// 					sumValue := 0.0
// 					for w := vertexArray[v]; w < vertexArray[v+1]; w++ {
// 						// could improve cache locality here using GAS, PCPM
// 						sumValue += x[edgeArray[w]] / float64(outDegrees[edgeArray[w]])
// 					}
// 					new_x[v] = alphaTerm + alpha*sumValue + leak/numNodes
// 					deltaSums[parIndex] += math.Abs(new_x[v] - x[v])
// 					if outDegrees[v] == 0 {
// 						leaks[parIndex] += new_x[v]
// 					}
// 				}
// 				wg.Done()
// 			}
// 		}(i, sliceStart, sliceEnd)
// 	}

// 	// wait for initialization to complete
// 	wg.Wait()

// 	leak = 0.0
// 	for v := 0; v < n; v++ {
// 		if outDegrees[v] == 0 { //dangling nodes
// 			leak += x[v]
// 		}
// 	}
// 	leak *= alpha

// 	iters := 0

// 	for {
// 		iters++

// 		deltaSum := 0.0
// 		// log.Println("leak", leak)

// 		wg.Add(numParallel)
// 		for i := 0; i < numParallel; i++ {
// 			signallers[i] <- struct{}{}
// 		}
// 		wg.Wait()

// 		leak = 0.0
// 		for i := 0; i < numParallel; i++ {
// 			leak += leaks[i]
// 			leaks[i] = 0.0

// 			deltaSum += deltaSums[i]
// 			deltaSums[i] = 0.0
// 		}
// 		leak *= alpha

// 		// swap x and new_x instead of replacing values
// 		temp := x
// 		x = new_x
// 		new_x = temp
// 		// log.Println(deltaSum)

// 		if deltaSum < eps {
// 			break
// 		}
// 	}

// 	// log.Println("parallel iters: ", iters)

// 	// for i := 0; i < numParallel; i++ {
// 	// 	close(signallers[i])
// 	// }

// 	norm := 0.0
// 	for _, v := range x {
// 		norm += v
// 	}

// 	for i := range x {
// 		x[i] /= norm
// 	}

// 	return x, iters

// }

// var enableLog = true

// func readGraph(edgeFileName string, nodeFileName string) ([][2]int, [][2]string, map[int]int) {
// 	f, err := os.Open(edgeFileName)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer f.Close()

// 	var dirLinks [][2]int
// 	var fileContents []string
// 	scanner := bufio.NewScanner(f)

// 	for scanner.Scan() {
// 		fileContents = append(fileContents, scanner.Text())
// 	}

// 	for _, fc := range fileContents {
// 		res := strings.Split(fc, ",")
// 		if len(res) < 2 {
// 			fmt.Println("less than 2 in dirLinks")
// 			continue
// 		}
// 		src, _ := strconv.Atoi(res[0])
// 		dest, _ := strconv.Atoi(res[1])
// 		dirLinks = append(dirLinks, [2]int{src, dest})
// 	}

// 	f, err = os.Open(nodeFileName)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer f.Close()

// 	var pageNum [][2]string
// 	var fc2 []string
// 	scanner = bufio.NewScanner(f)

// 	for scanner.Scan() {
// 		fc2 = append(fc2, scanner.Text())
// 	}

// 	for _, fc := range fc2 {
// 		res := strings.Split(fc, ",")
// 		if len(res) < 2 {
// 			fmt.Println("less than 2 in pageNum")
// 			continue
// 		}
// 		pageNum = append(pageNum, [2]string{res[0], res[1]})
// 	}

// 	if enableLog {
// 		fmt.Println("No. of edges: ", len(dirLinks))
// 		fmt.Println("No. of nodes: ", len(pageNum))
// 	}
// 	//fmt.Println("No. of nodes: ", 2080370)

// 	v := 0
// 	node_to_index := make(map[int]int)

// 	for _, node := range pageNum {
// 		temp, _ := strconv.Atoi(node[0])
// 		node_to_index[temp] = v
// 		v += 1
// 	}

// 	for i, edge := range dirLinks {
// 		if _, ok := node_to_index[edge[0]]; !ok {
// 			fmt.Printf("node %v does not exist in node_to_index", edge[0])
// 		}
// 		if _, ok := node_to_index[edge[1]]; !ok {
// 			fmt.Printf("node %v does not exist in node_to_index", edge[1])
// 		}
// 		dirLinks[i][0] = node_to_index[edge[0]]
// 		dirLinks[i][1] = node_to_index[edge[1]]
// 	}

// 	return dirLinks, pageNum, node_to_index
// }

// func makeAdjArray(edges [][2]int, n int) map[int][]int {

// 	adj_array := make(map[int][]int)

// 	for i := 0; i < n; i++ {
// 		adj_array[i] = make([]int, 0)
// 	}

// 	for _, edge := range edges {
// 		adj_array[edge[0]] = append(adj_array[edge[0]], edge[1])
// 	}

// 	return adj_array

// }

// // converts the edge list to CSR format with
// // a vertexArray consisting of cumulative in degrees
// // an edgeArray consisting of the in-vertices
// func makeCSR(edges [][2]int, n int) (vertexArray, edgeArray, outDegrees []int) {
// 	// stores indegrees of all vertices
// 	inDegrees := make([]int, n)

// 	// stores outdegrees of all vertices
// 	// needed in pagerank computation
// 	outDegrees = make([]int, n)

// 	// vertex array of CSR format
// 	vertexArray = make([]int, n+1)
// 	// edge or destination array
// 	edgeArray = make([]int, len(edges))

// 	// populate indegrees and outDegrees
// 	//  - increment inDegree of the destination vertex
// 	//  - increment out degree of the source vertex
// 	for i := 0; i < len(edges); i++ {
// 		inDegrees[edges[i][1]]++
// 		outDegrees[edges[i][0]]++
// 	}

// 	// populate vertex  array
// 	// with cumulative indegree
// 	//   previous value + indegree of previous node
// 	// first value is always 0
// 	for i := 1; i <= n; i++ {
// 		vertexArray[i] = vertexArray[i-1] + inDegrees[i-1]
// 	}

// 	// we will reuse inDegrees to store number of in-vertices
// 	//  currently processed for that node
// 	// this is used to index into edgeArray
// 	for i := 0; i < n; i++ {
// 		inDegrees[i] = 0
// 	}

// 	// populate edge array
// 	for i := 0; i < len(edges); i++ {
// 		// since we need incoming vertices, we store
// 		//   destination to source mapping
// 		from := edges[i][1]
// 		to := edges[i][0]
// 		edgeArray[vertexArray[from]+inDegrees[from]] = to

// 		inDegrees[from]++
// 	}

// 	// the vertices in edgeArray are currently randomly ordered
// 	// let's sort them to reduce randomm accesses
// 	for i := 0; i < n; i++ {
// 		// if at least 2 edges exist, sort them
// 		if vertexArray[i+1] > vertexArray[i]+1 {
// 			sort.Ints(edgeArray[vertexArray[i]:vertexArray[i+1]])
// 		}
// 	}

// 	return
// }

// func max(arr []float64) float64 {

// 	res := arr[0]
// 	for i := 1; i < len(arr); i++ {
// 		if arr[i] > res {
// 			res = arr[i]
// 		}
// 	}

// 	return res
// }

// func main() {
// 	f, err := os.Open("wiki-vote-edges.txt")
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer f.Close()

// 	edges, pages, _ := readGraph("wiki-vote-edges.txt", "wiki-vote-nodes.txt")
// 	vertexArray, edgeArray, outDegrees := makeCSR(edges, len(pages))
// 	pageRankGPU(vertexArray, edgeArray, outDegrees, 0.85, 10e-6)

// }
