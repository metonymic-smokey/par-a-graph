package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

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

	n := 281903 + 4
	e := 2312497
	granularity := 8

	val := make([]float64, e)
	col_ind := make([]int, e)
	row_ptr := make([]int, n+1)

	row_ptr[0] = 0
	cur_row := 0
	i := 0
	j := 0
	elrow := 0
	curel := 0

	for _, fc := range fileContents {
		res := strings.SplitN(fc, ",", -1)
		if len(res) < 2 {
			continue
		}

		src, _ := strconv.Atoi(res[0])
		dest, _ := strconv.Atoi(res[1])
		//wt, _ := strconv.Atoi(res[2])

		if src > cur_row {
			curel = curel + elrow
			for k := cur_row + 1; k <= src; k++ {
				row_ptr[k] = curel
			}
			elrow = 0
			cur_row = src
		}

		val[i] = 1.0
		col_ind[i] = dest
		elrow++
		i++
	}
	row_ptr[cur_row+1] = curel + elrow - 1

	out_link := make([]int, n)
	for a := 0; a < n; a++ {
		out_link[a] = 0
	}

	rowel := 0
	for i = 0; i < n; i++ {
		if row_ptr[i+1] != 0 {
			rowel = row_ptr[i+1] - row_ptr[i]
			out_link[i] = rowel
		}
	}

	curcol := 0
	for i = 0; i < n; i++ {
		rowel = row_ptr[i+1] - row_ptr[i]
		for j := 0; j < rowel; j++ {
			val[curcol] = val[curcol] / float64(out_link[i])
			curcol++
		}
	}

	d := 0.85
	p := make([]float64, n)
	for i := 0; i < n; i++ {
		p[i] = (float64)(1.0 / n)
	}

	looping := 1
	k := 0
	//parallel := 0

	p_new := make([]float64, n)

	start := time.Now()

	for looping > 0 {
		for i = 0; i < n; i++ {
			p_new[i] = 0.0
		}

		rowel := 0
		curcol := 0

		for i = 0; i < n; i = i + granularity {
			rowel = row_ptr[i+1] - row_ptr[i]
			for j = 0; j < rowel; j++ {
				p_new[col_ind[curcol]] = p_new[col_ind[curcol]] + val[curcol]*p[i]
				curcol++
			}
		}

		for i = 0; i < n; i++ {
			p_new[i] = d*p_new[i] + ((1.0 - d) / float64(n))
		}

		err := 0.0
		for i = 0; i < n; i++ {
			err = err + (float64)(math.Abs(p_new[i]-p[i]))
		}

		if err < 0.000001 {
			looping = 0
		}

		for i = 0; i < n; i++ {
			p[i] = p_new[i]
		}

		k = k + 1
	}

	fmt.Printf("took %v\n", time.Since(start))

	fmt.Println("Final vals:")
	for i = 0; i < n; i++ {
		fmt.Println(p[i])
	}

}
