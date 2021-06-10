# `groph-parallel`

Parallel implementations of graph algorithms using Go. Parallelization is performed both on CPU (using goroutines) and GPU (using CUDA).

## Algorithms

### Pagerank

Parallel implementation of Pagerank. This was done as part of the project for the course "Heterogeneous Parallelism" (UE18CS342) at PES University.

#### CPU

Details can be found in the [report slides](docs/report-presentation.pdf) and the [README](./cpu).

Summary:
 - Fixed **partitioning** based parallelization using goroutines
 - Tested and **benchmarked on 5 datasets** of varying sizes and characteristics (number of nodes, edges, average degree, density, etc.)
 - **Up to 4.2x speedup** on the Stanford Web graph (on a system with 4 cores, 8 threads)


#### CUDA

TODO: rough implementation can be found in the `cuda_pr` branch.

## Directory Structure

 - [`cpu`](./cpu): CPU-based parallel implementation of Pagerank
 - [`datasets`](./datasets)
 - [`images`](./images): images for README, etc.
 - [`docs`](./docs): documents related to the project - presentation slides, etc.

## License

TODO
