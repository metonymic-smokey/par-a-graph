# Parallel Pagerank using goroutines

Parallel implementation of Pagerank in Go benchmarked on various datasets, achieving upto 4.2x speedup (on 4 core, 8 threads system).

## Implementation Details

### Parallelization

Various techniques were tried, with iterative improvements. More details regarding the other techniques can be found in the [report slides](../docs/report-presentation.pdf) and the commit history. We describe the final implementation here.

We partition the input graph into a fixed number of blocks ([16 was chosen for our experiments](topoPageRank.go#L105)). Partitioning is based on the nodes (and not edges). Each block has the following associated with it:
 - a local [`deltaSum`](topoPageRank.go#L122) to store the error (delta) for each partition. This will be summed later at the end of every iteration.
 - a local [`leak`](topoPageRank.go#L111) term to store the leak from nodes with no out-neighbours. This will be summed at the end of every iteration.
 - a [signal channel](topoPageRank.go#L125) to signal the start of an iteration. A signal on this channel makes the goroutines start computation.
 - a goroutine that performs the computation for that particular block.

Each goroutine does the following:
 - Initialize values of the pagerank vector with $1/n$ as soon as they are launched. These are seen as yellow marks at the top of each goroutine in the visualization below.
 - When a signal is given over the signal channel, it performs the pagerank computation. The wait group is used to synchronise all of the goroutines.

Below is a visualization of the goroutines using [gotrace](https://github.com/divan/gotrace). At the center is the main goroutine and the rest are each of the goroutines for each partition. The yellow parts are CPU intensive operations - initializing pagerank vector and pagerank computation. The blue lines are channel send/receive - over the signal channels. The short breaks are between each iteration of pagerank.

![](https://drive.google.com/uc?export=view&id=1iL0QNMGY4xqN-R0NqCh_rVHAZ-wc-E_J)

Note: some small changes were made for the sake of visualization. A few explicit pauses had to be added to show all of the phases.

## Usage

### Testing

Correctness tests compare the outputs from [Daniel Cadenas' implementation in Go](https://github.com/dcadenas/pagerank), serial Pagerank and parallel Pagerank.

To run all correctness tests on all datasets:
```
go test
```

To run tests for a specific dataset. For example, the Wiki vote dataset (note the first letter should be in uppercase):
```
go test -run Wiki
```
Refer [Datasets](#datasets) for a list of available datasets.

### Benchmarks

Benchmarks compare the execution time of Serial vs Parallel on three different values of `eps` (error threshold): `10e-6`, `10e-9` and `10e-11`.

To run all benchmarks on all datasets:
```
go test -run None -bench .
```

To run benchmarks for a specific dataset. For example, the Stanford web graph (note: the first letter is in uppercase):
```
go test -run None -bench Stanford
```
Refer [Datasets](#datasets) for a list of available datasets.

## Code Structure

### Implementation

Serial: [`topoPageRank.go`](./topoPageRank.go#L10)

Parallel: [`topoPageRank.go`](./topoPageRank.go#L89)

### Utilities

 - [`main.go`](./main.go):
    - utilities for reading datasets
    - utilities for converting graphs to adjacency array and CSR format
    - main function to be used for `go tool pprof` profiling`
       - Note: binaries built from this will not include any tests or benchmarks

### Tests and Benchmarks

 - [`topoPageRank_test.go`](./topoPageRank_test.go): Test helpers and tests & benchmarks for small & large graph
 - [`wikiVote_test.go`](./wikiVote_test.go): Tests and benchmarks for Wiki Vote Dataset
 - [`stanford_test.go`](./stanford_test.go): Tests and benchmarks for Stanford Web graph
 - [`quora_test.go`](./quora_test.go): Tests and benchmarks for Quora Question pairs graph

## Datasets

List of datasets which are tested and benchmarked on (these values can be passed to `-run` or `-bench`):
 - `Small`
 - `Wiki`
 - `Quora`
 - `Stanford`
 - `Large`

Refer [Datasets](../datasets/) for more details.

## Results

Benchmark results as tested on a 4-core, 8-threads system. Refer the slides for information on cache characteristics of different graphs which explains the results we observed.

### Large Graph - 1.58x speedup

![Results on the Large Graph showing upto 1.58x speedup from serial to parallel](../images/cpu_1_large_graph.png)

### Quora Question pairs Graph - 2.64x speedup

![Results on the Quora Graph showing upto 2.64x speedup from serial to parallel](../images/cpu_2_quora.png)

### Stanford Web Graph - 4.23x speedup

![Results on the Stanford Web Graph showing upto 4.23x speedup from serial to parallel](../images/cpu_3_stanford.png)

### Wiki Vote Graph - 2.02x speedup

![Results on the Wiki Vote Graph showing upto 2.02x speedup from serial to parallel](../images/cpu_4_wiki.png)

### Small Graph - 42x _slowdown_

![Results on the Small Graph showing a huge 2.02x slowdown from serial to parallel](../images/cpu_5_small.png)
