package main

import (
	"bufio"
	hp "container/heap"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type route struct {
	route       *[]pathInfo
	totalWeight float64
}

type pathInfo struct {
	uid int
	//attr int
}

type queueItem struct {
	uid  int     // uid of the node.
	cost float64 // cost of taking the path till this uid.
	// number of hops taken to reach this node. This is useful in finding out if we need to
	// expandOut after poping an element from the heap. We only expandOut if item.hop > numHops
	// otherwise expanding would be useless.
	hop   int
	index int
	path  route // used in k shortest path.
}

type priorityQueue []*queueItem

type mapItem struct {
	//attr int
	cost float64
}

type nodeInfo struct {
	mapItem
	parent int
	// Pointer to the item in heap. Used to update priority
	node *queueItem
}

func (h priorityQueue) Len() int { return len(h) }

func (h priorityQueue) Less(i, j int) bool { return h[i].cost < h[j].cost }

func (h priorityQueue) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *priorityQueue) Push(val interface{}) {
	n := len(*h)
	item := val.(*queueItem)
	item.index = int(n)
	*h = append(*h, item)
}

func (h *priorityQueue) Pop() interface{} {
	old := *h
	n := len(old)
	val := old[n-1]
	*h = old[0 : n-1]
	val.index = -1
	return val
}

func dgraphShortest(adjacencyMap map[int]map[int]mapItem, src int, dest int) []int {
	pq := make(priorityQueue, 0)
	srcNode := &queueItem{
		uid:  src,
		cost: 0,
		hop:  0,
	}
	hp.Push(&pq, srcNode)

	var maxHops int
	maxHops = math.MaxInt32

	if maxHops == 0 {
		return []int{}
	}

	dist := make(map[int]nodeInfo)
	dist[srcNode.uid] = nodeInfo{
		parent: 0,
		node:   srcNode,
		mapItem: mapItem{
			cost: 0,
		},
	}

	for pq.Len() > 0 {
		item := hp.Pop(&pq).(*queueItem)
		if item.uid == dest { //reached destination
			break
		}

		neighbours := adjacencyMap[item.uid]

		for toUID, neighbour := range neighbours {
			d, ok := dist[toUID]
			nodeCost := item.cost + float64(neighbour.cost)

			if ok && d.cost <= nodeCost {
				continue
			}

			var node *queueItem
			if !ok {
				// This is the first time we're seeing this node. So
				// create a new node and add it to the heap and map.
				node = &queueItem{
					uid:  int(toUID),
					cost: nodeCost,
					hop:  item.hop + 1,
				}
				hp.Push(&pq, node)
			} else {
				// We've already seen this node. So, just update the cost
				// and fix the priority in the heap and map.
				node = dist[toUID].node
				node.cost = nodeCost
				node.hop = item.hop + 1
				hp.Fix(&pq, node.index)
			}
			dist[toUID] = nodeInfo{
				parent: item.uid,
				node:   node,
				mapItem: mapItem{
					cost: nodeCost,
				},
			}

		}

	}

	var result []int
	cur := dest

	for i := 0; i < len(dist); i++ {
		result = append(result, cur)
		if cur == src {
			break
		}
		cur = dist[cur].parent
	}

	l := len(result)
	// Reverse the list.
	for i := 0; i < l/2; i++ {
		result[i], result[l-i-1] = result[l-i-1], result[i]
	}

	return result
}

func main() {
	graph := newGraph()

	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	var fileContents []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		fileContents = append(fileContents, scanner.Text())
	}

	adjacencyMap := make(map[int]map[int]mapItem)

	for _, fc := range fileContents {
		res := strings.SplitN(fc, ",", -1)
		src, _ := strconv.Atoi(res[0])
		dest, _ := strconv.Atoi(res[1])
		wt, _ := strconv.Atoi(res[2])
		graph.addEdge(src, dest, wt)
		if _, ok := adjacencyMap[src]; !ok {
			adjacencyMap[src] = make(map[int]mapItem)
		}
		adjacencyMap[src][dest] = mapItem{float64(wt)}
	}

	fmt.Println("Adjacency List: ")
	graph.printGraph()

	srcNode := 5
	destNode := 6

	dgraphResult := dgraphShortest(adjacencyMap, srcNode, destNode)

	fmt.Println("Path from trial code: ", dgraphResult)

	fmt.Println("Dijkstra's result: ")
	res, path := (graph.getPath(srcNode, destNode))
	fmt.Println(res)
	fmt.Println("Correct path: ", path)
}