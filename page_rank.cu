#include <cstring>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <cuda.h>
__global__
void prGPU(double *pagerank_vector,
                              double *new_pagerank, long *vertexArray,
	                      long *edgeArray, long *outDegrees, double alpha,
                              double *deltaSum, int n, double *leak) {

  int index = blockIdx.x * blockDim.x + threadIdx.x;
  int numthreads = blockDim.x * gridDim.x;
  
  for (int i = index; i < n; i += numthreads) {
    double temp = pagerank_vector[i];
    double sumVal = 0;
    for (int neighbour = vertexArray[i]; neighbour < vertexArray[i + 1]; neighbour++) {
      double num = pagerank_vector[edgeArray[neighbour]];
      double denom = double(outDegrees[edgeArray[neighbour]]);
      /* printf("num = %lf, denom = %lf\n", num, denom); */
      sumVal += num / denom;
    }
    new_pagerank[i] =
        (1 - alpha) / double(n) + alpha * sumVal + *leak / double(n);
    double delta = abs(new_pagerank[i] - temp);

    atomicAdd(deltaSum, delta);
  }
}

__global__
void calcLeak(int n, long* outDegrees, double* pageRank, double* leak, double alpha) {
  int index = blockIdx.x * blockDim.x + threadIdx.x;
  int numthreads = blockDim.x * gridDim.x;

  double localLeak = 0.0;
  for (int i = index; i < n; i += numthreads) {
    if (outDegrees[i] == 0) {
      localLeak += pageRank[i] * alpha;
    }
  }

  atomicAdd(leak, localLeak);
}

__global__ void 