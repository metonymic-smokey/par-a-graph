package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"time"
	"unsafe"

	"gorgonia.org/cu"
)

func setupGPU() *cu.CUContext {
	devices, _ := cu.NumDevices()
	if devices == 0 {
		log.Fatalln("No device!")
	}

	// fmt.Println("Getting device")
	dev := cu.Device(0)

	// fmt.Println("Creating new Context")
	context, err := dev.MakeContext(cu.SchedAuto)
	if err != nil {
		log.Fatalf("Err while trying to MakeContext: %v", err)
	}
	if err = cu.SetCurrentContext(context); err != nil {
		log.Fatalln("Error while setting current context!", err)
	}
	return &context
}

func pageRank(
	vertexArray []int,
	edgeArray []int,
	outDegrees []int,
	alpha float64,
	eps float64) ([]float64, int) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ctx := setupGPU()
	defer ctx.Destroy()

	n := len(vertexArray) - 1

	// pagerank vector
	x := make([]float64, n)
	new_x := make([]float64, n)

	for i := 0; i < n; i++ {
		x[i] = 1 / float64(n)
	}

	//all the nodes in 1 slice
	/* 	var nodes []int
	   	for i := 0; i < n; i++ {
	   		nodes = append(nodes, i)
	   	} */

	// delta := make([]float64, n)

	iters := 0
	// fmt.Println("Allocating and copying memory to GPU")
	// memory being initialized
	var memVertex, memOut, memEdge, memPRVec, memPRNewVec, memDeltaSum, memLeak cu.DevicePtr
	// fmt.Println("VertexArray")
	// fmt.Println("size of vertexarry = ", len(vertexArray))
	var err error
	memVertex, err = cu.AllocAndCopy(unsafe.Pointer(&vertexArray[0]), int64(len(vertexArray)*8))
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("OutDegrees")
	memOut, err = cu.AllocAndCopy(unsafe.Pointer(&outDegrees[0]), int64(len(outDegrees)*8))
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("EdgeArray")
	memEdge, err = cu.AllocAndCopy(unsafe.Pointer(&edgeArray[0]), int64(len(edgeArray)*8))
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("PageRank")
	memPRVec, err = cu.AllocAndCopy(unsafe.Pointer(&x[0]), int64(len(x)*8))
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("NewPageRank")
	memPRNewVec, err = cu.AllocAndCopy(unsafe.Pointer(&new_x[0]), int64(len(new_x)*8))
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("DeltaSum")
	memDeltaSum, err = cu.MemAlloc(8)
	if err != nil {
		log.Fatalln(err)
	}

	memLeak, err = cu.MemAlloc(8)
	if err != nil {
		log.Fatalln(err)
	}
	// copying memory from host to device
	// if err = cu.Al(memPRVec, unsafe.Pointer(&x[0]), int64(len(x)*8)); err != nil {
	// 	log.Fatalf("Error while copying pagerank vec to device: %v", err)
	// }
	// if err = cu.MemcpyHtoD(memPRNewVec, unsafe.Pointer(&new_x[0]), int64(len(new_x)*8)); err != nil {
	// 	log.Fatalf("Error while copying new pagerank vec to device: %v", err)
	// }

	// fmt.Println("Loading ptx")
	mod, err := cu.Load("pagerank.ptx")
	if err != nil {
		log.Fatalln("Error while loading ptx: ", err)
	}
	// fmt.Println("Loading function")
	fn, err := mod.Function("_Z5prGPUPdS_PlS0_S0_dS_iS_") //
	if err != nil {
		log.Fatalln("Cannot load function!: ", err)
	}

	fnLeak, err := mod.Function("_Z8calcLeakiPlPdS0_d")
	if err != nil {
		log.Fatalln("Cannot load leak function:", err)
	}

	// fmt.Printf("Printing out degrees: %v\n", outDegrees)
	// fmt.Printf("Printing edgeArray: %v\n", edgeArray)
	// fmt.Printf("Printing vertexArray: %v\n", vertexArray)

	// fmt.Println("Going into loop")
	// return x, 0
	start := time.Now()
	for {

		iters++
		deltaSum := make([]float64, 1)
		// leak := make([]float64, 1)
		/*   		for _, v := range nodes {
			if outDegrees[v] == 0 {
				leak += x[v]
			}
		} */
		// if err = cu.MemcpyHtoD(memLeak, unsafe.Pointer(&leak[0]), 8); err != nil {
		// 	log.Fatalf("Error while copying leak to device!: %v", err)
		// }
		if err = cu.MemsetD32(memLeak, 0, 2); err != nil {
			log.Fatalf("Error while setting memLeak: %v", err)
		}
		if err = cu.MemsetD32(memDeltaSum, 0, 2); err != nil {
			log.Fatalf("Error while setting memDeltaSum: %v", err)
		}
		// leak *= alpha

		argsLeak := []unsafe.Pointer{
			unsafe.Pointer(&n),
			unsafe.Pointer(&memOut),
			unsafe.Pointer(&memPRVec),
			unsafe.Pointer(&memLeak),
			unsafe.Pointer(&alpha),
		}

		if err = fnLeak.LaunchAndSync(32, 1, 1, 256, 1, 1, 0, cu.Stream{}, argsLeak); err != nil {
			log.Fatalf("Failed to launch leak kernel! %v", err)
		}

		// if err = cu.MemcpyHtoD(memDeltaSum, unsafe.Pointer(&deltaSum[0]), 8); err != nil {
		// 	log.Fatalf("Error while copying deltasum to device: %v", err)
		// }

		// function arguments
		args := []unsafe.Pointer{
			unsafe.Pointer(&memPRVec),
			unsafe.Pointer(&memPRNewVec),
			unsafe.Pointer(&memVertex),
			unsafe.Pointer(&memEdge),
			unsafe.Pointer(&memOut),
			unsafe.Pointer(&alpha),
			unsafe.Pointer(&memDeltaSum),
			unsafe.Pointer(&n),
			unsafe.Pointer(&memLeak),
		}

		// fmt.Println("Launching kernel")
		// Launching and syncing to wait for kernel to finish
		if err = fn.LaunchAndSync(32, 1, 1, 256, 1, 1, 0, cu.Stream{}, args); err != nil {
			log.Fatalf("Failed to launch page rank kernel! %v", err)
		}

		// copying memory from device to host
		// if err = cu.MemcpyDtoH(unsafe.Pointer(&x[0]), memPRVec, int64(len(x)*8)); err != nil {
		// 	log.Fatalf("Error while copying pagerank from device to host! %v", err)
		// }
		/*
			if err = cu.MemcpyDtoH(unsafe.Pointer(&new_x[0]), memPRNewVec, int64(len(new_x)*8)); err != nil {
				log.Fatalf("Error while copying new_pagerank from device to host! %v", err)
			} */

		// if err = cu.MemcpyDtoD(memPRVec, memPRNewVec, int64(len(x)*8)); err != nil {
		// 	log.Fatalf("Error while updating pagerank with new pr values!: %v", err)
		// }

		temp := memPRVec
		memPRVec = memPRNewVec
		memPRNewVec = temp

		if err = cu.MemcpyDtoH(unsafe.Pointer(&deltaSum[0]), memDeltaSum, 8); err != nil {
			log.Fatalf("Error while copying deltasum to host: %v", err)
		}

		// updating page rank
		/* 		copy(x, new_x) */
		// log.Println(deltaSum)
		// checking for convergence
		if deltaSum[0] < 10e-6 {
			break
		} else if iters > 100 {
			break
		}
	}
	fmt.Println("Time taken", time.Since(start))

	// copy pagerank back to host
	if err = cu.MemcpyDtoH(unsafe.Pointer(&x[0]), memPRVec, int64(len(x)*8)); err != nil {
		log.Fatalf("Error while copying pagerank from device to host! %v", err)
	}

	cu.MemFree(memEdge)
	cu.MemFree(memOut)
	cu.MemFree(memPRNewVec)
	cu.MemFree(memPRVec)
	cu.MemFree(memVertex)

	// normalizing the pagerank
	norm := 0.0
	for _, v := range x {
		norm += v
	}

	for i := range x {
		x[i] /= norm
	}

	return x, iters

}

// func pageRankGPU(
// 	vertexArray []int,
// 	edgeArray []int,
// 	outDegrees []int,
// 	alpha float64,
// 	eps float64) ([]float64, int) {

// 	var vertexArray_size C.int
// 	var edgeArray_size C.int
// 	var outDegrees_size C.int

// 	vertexArray_size = C.int(len(vertexArray))
// 	edgeArray_size = C.int(len(edgeArray))
// 	outDegrees_size = C.int(len(outDegrees))

// 	_vertexArray := make([]C.int, vertexArray_size)
// 	_edgeArray := make([]C.int, edgeArray_size)

// 	_outDegrees := make([]C.int, outDegrees_size)
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

// 	f, err := os.Open("pr_cuda_res.txt")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	scanner := bufio.NewScanner(f)

// 	var fileContents []string
// 	for scanner.Scan() {
// 		fileContents = append(fileContents, scanner.Text())
// 	}

// 	pagerank := fileContents[1 : len(fileContents)-3]
// 	_pagerank := make([]float64, len(pagerank))
// 	for i := 0; i < len(pagerank); i++ {
// 		_pagerank[i], _ = strconv.ParseFloat(pagerank[i], 64)
// 	}
// 	iters, _ := strconv.Atoi(fileContents[len(fileContents)-1])
// 	return _pagerank, iters
//}

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

	for {
		iters++

		deltaSum := 0.0
		var leak float64

		for _, v := range nodes {
			if outDegrees[v] == 0 { //dangling nodes
				leak += x[v]
				// log.Println("found leak", x[v], outDegrees[v], v)
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
		// log.Println(deltaSum)

		if deltaSum < eps {
			break
		}
	}

	// log.Println("serial iters: ", iters)

	norm := 0.0
	for _, v := range x {
		norm += v
	}

	for i := range x {
		x[i] /= norm
	}

	return x, iters

}

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
