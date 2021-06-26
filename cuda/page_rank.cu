#include <cstring>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <cuda.h>
__global__
void prGPU(double *pagerank_vector,
                              double *new_pagerank, long long *vertexArray,
	                      long long *edgeArray, long long *outDegrees, double alpha,
                              double *deltaSum, long long n, double *leak) {

  long long index = blockIdx.x * blockDim.x + threadIdx.x;
  long long numthreads = blockDim.x * gridDim.x;
  
  for (long long i = index; i < n; i += numthreads) {
    double temp = pagerank_vector[i];
    double sumVal = 0;
    for (long long neighbour = vertexArray[i]; neighbour < vertexArray[i + 1]; neighbour++) {
      double num = pagerank_vector[edgeArray[neighbour]];
      long long denom = outDegrees[edgeArray[neighbour]];
      sumVal += num / denom;
    }
    new_pagerank[i] =
      (1 - alpha) / (double) n + alpha * sumVal + *leak / (double) n;
    double delta = abs(new_pagerank[i] - temp);

    atomicAdd(deltaSum, delta);
  }
}

__global__
void calcLeak(long long n, long long* outDegrees, double* pageRank, double* leak, double alpha) {
  long long index = blockIdx.x * blockDim.x + threadIdx.x;
  long long numthreads = blockDim.x * gridDim.x;

  double localLeak = 0.0;
  for (long long i = index; i < n; i += numthreads) {
    if (outDegrees[i] == 0) {
      localLeak += pageRank[i] * alpha;
    }
  }

  atomicAdd(leak, localLeak);
}
