package main

import (
	// "fmt"
	"math"
	"sync"
)

func pagerankSerial(n int, granularity int, row_ptr []int, col_ind []int, val []float64) []float64 {
	d := 0.85
	p := make([]float64, n)
	for i := 0; i < n; i++ {
		p[i] = (float64)(1.0 / n)
	}

	looping := 1
	k := 0
	//parallel := 0

	p_new := make([]float64, n)

	for looping > 0 {
		for i := 0; i < n; i++ {
			p_new[i] = 0.0
		}

		rowel := 0
		curcol := 0

		for i := 0; i < n; i = i + granularity {
			rowel = row_ptr[i+1] - row_ptr[i]
			for j := 0; j < rowel; j++ {
				p_new[col_ind[curcol]] = p_new[col_ind[curcol]] + val[curcol]*p[i]
				curcol++
			}
		}

		for i := 0; i < n; i++ {
			p_new[i] = d*p_new[i] + ((1.0 - d) / float64(n))
		}

		err := 0.0
		for i := 0; i < n; i++ {
			err = err + (float64)(math.Abs(p_new[i]-p[i]))
		}

		if err < 0.000001 {
			looping = 0
		}

		for i := 0; i < n; i++ {
			p[i] = p_new[i]
		}

		k = k + 1
	}

	// fmt.Println("Serial", k)

	return p
}

func pagerank(n int, granularity int, row_ptr []int, col_ind []int, val []float64) []float64 {
	d := 0.85
	p := make([]float64, n)
	for i := 0; i < n; i++ {
		p[i] = (float64)(1.0 / n)
	}

	looping := 1
	k := 0

	p_new := make([]float64, n)

	var wg sync.WaitGroup
	type UpdateVal struct {
		ind int
		val float64
	}
	pCh := make(chan UpdateVal, n)
	quit := make(chan struct{})
	numParallel := 8
	for i := 0; i < numParallel; i++ {
		go func() {
			for {
				select {

				case v := <-pCh:
					p_new[v.ind] += v.val

					wg.Done()

				case <-quit:
					return
				}
			}
		}()
	}

	for looping > 0 {
		for i := 0; i < n; i++ {
			p_new[i] = 0.0
		}

		rowel := 0
		curcol := 0

		for i := 0; i < n; i = i + granularity {
			rowel = row_ptr[i+1] - row_ptr[i]
			if (rowel > 0) {
				wg.Add(rowel)
			}
			for j := 0; j < rowel; j++ {
				pCh <- UpdateVal{col_ind[curcol], val[curcol] * p[i]}
				curcol++
			}
			curcol += rowel
		}

		wg.Wait()

		for i := 0; i < n; i++ {
			p_new[i] = d*p_new[i] + ((1.0 - d) / float64(n))
		}

		err := 0.0
		for i := 0; i < n; i++ {
			err = err + (float64)(math.Abs(p_new[i]-p[i]))
		}

		if err < 0.000001 {
			looping = 0
		}

		for i := 0; i < n; i++ {
			p[i] = p_new[i]
		}

		k = k + 1
	}

	close(quit)

	return p
}
