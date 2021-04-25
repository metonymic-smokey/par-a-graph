package main

import (
	// "fmt"
	"math"
	"sync"
	"go.uber.org/atomic"
)

const eps float64 = 0.000001

func pagerankSerial(n int, granularity int, row_ptr []int, col_ind []int, val []float64) []float64 {
	d := 0.85
	p := make([]float64, n)
	for i := 0; i < n; i++ {
		p[i] = (float64)(1.0 / n)
	}

	looping := 1
	k := 0
	//parallel := 0

	p_new := make([]*atomic.Float64, n)
	p_new_l := make([]float64, n)

	for looping > 0 {
		for i := 0; i < n; i++ {
			p_new[i] = atomic.NewFloat64(0.0)
		}

		rowel := 0
		curcol := 0

		for i := 0; i < n; i = i + granularity {
			rowel = row_ptr[i+1] - row_ptr[i]
			for j := 0; j < rowel; j++ {
				// p_new[col_ind[curcol]] += val[curcol]*p[i]
				p_new[col_ind[curcol]].Add(val[curcol] * p[i])
				curcol++
			}
		}

		for i := 0; i < n; i++ {
			p_new_l[i] = d*p_new[i].Load() + ((1.0 - d) / float64(n))
		}

		err := 0.0
		for i := 0; i < n; i++ {
			err = err + (float64)(math.Abs(p_new_l[i]-p[i]))
		}

		if err < eps {
			looping = 0
		}

		for i := 0; i < n; i++ {
			p[i] = p_new_l[i]
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

	looping := true
	k := 0

	p_new := make([]*atomic.Float64, n)
	p_new_l := make([]float64, n)

	var wg sync.WaitGroup
	type UpdateVal struct {
		i int
		curcol int
	}
	numParallel := 10
	pCh := make(chan UpdateVal, n/granularity + 1)
	quit := make(chan struct{})

	for i := 0; i < numParallel; i++ {
		go func() {
			for {
				v := <-pCh
					rowel := row_ptr[v.i+1] - row_ptr[v.i]
					curcol := v.curcol
					for j := 0; j < rowel; j++ {
						// p_new[col_ind[curcol]] += val[curcol] * p[v.i]
						p_new[col_ind[curcol]].Add(val[curcol] * p[v.i])
						curcol++
					}

					wg.Done()
			}
		}()
	}

		for i := 0; i < n; i++ {
			p_new[i] = atomic.NewFloat64(0.0)
		}

	for looping {
		for i := 0; i < n; i++ {
			p_new[i].Store(0.0)
		}

		curcol := 0
		wg.Add(n/granularity + 1)

		for i := 0; i < n; i = i + granularity {
			pCh <- UpdateVal{i, curcol}
			curcol += row_ptr[i+1] - row_ptr[i]
		}

		wg.Wait()

		for i := 0; i < n; i++ {
			p_new_l[i] = d*p_new[i].Load() + ((1.0 - d) / float64(n))
		}

		err := 0.0
		for i := 0; i < n; i++ {
			err = err + (float64)(math.Abs(p_new_l[i]-p[i]))
		}

		if err < eps {
			looping = false
		}

		for i := 0; i < n; i++ {
			p[i] = p_new_l[i]
		}

		k = k + 1
	}

	close(quit)

	return p
}
