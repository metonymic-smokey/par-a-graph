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

	// chans := make([]chan float64, n)
	// for i := range chans {
	// 	chans[i] = make(chan float64, 100)
	// }
	type ch chan float64
	p_channel := make([]ch, n)
	for i := range p_channel {
		p_channel[i] = make(chan float64, 10)
	}

	for looping > 0 {
		for i := 0; i < n; i++ {
			p_new[i] = 0.0
			// chans[i] <- 0.0
			p_channel[i] <- 0.0
		}

		rowel := 0
		curcol := 0
		var wg sync.WaitGroup
		wg.Add(n/granularity + 1)

		for i := 0; i < n; i = i + granularity {
			rowel = row_ptr[i+1] - row_ptr[i]
			go func(i, rowel, curcol int) {
				defer wg.Done()
				for j := 0; j < rowel; j++ {
					temp := <-p_channel[col_ind[curcol]]
					value := temp + val[curcol]*p[i]
					p_channel[col_ind[curcol]] <- value
					curcol++
				}
			}(i, rowel, curcol)
			// p_new[i] = <-chans[i]
			curcol += rowel
		}
		wg.Wait()
		for i := range p_channel {
			p_new[i] = <-p_channel[i]
		}
		// for i := 0; i < n; i++ {
		// 	p_new[i] = <-chans[i]
		// }

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

	// fmt.Println("Parallel", k)

	return p
}
