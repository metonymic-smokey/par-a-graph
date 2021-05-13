#include <cstring>
#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <cuda.h>

__global__
void prGPU(int *nodes, double *pagerank_vector,
                              double *new_pagerank, int *vertexArray,
                              int *edgeArray, int *outDegrees, double alpha,
                              double *deltaSum, int n, double leak) {
  int index = threadIdx.x;
  int numthreads = blockDim.x;
  for (int i = index; i < n; i += numthreads) {
    double temp = pagerank_vector[i];
    double sumVal = 0;
    for (int neighbour = vertexArray[i]; neighbour < vertexArray[i + 1];
         neighbour++) {
      sumVal += pagerank_vector[edgeArray[neighbour]] /
                double(outDegrees[edgeArray[neighbour]]);
    }
    new_pagerank[i] =
        (1 - alpha) / double(n) + alpha * sumVal + leak / double(n);
    /* printf("New pagerank = %lf\n",new_pagerank[i]); */
    double delta = abs(new_pagerank[i] - temp);
    atomicAdd(deltaSum, delta);
  }
}

extern "C" {
void pageRankCuda(int *vertexArray, int vertexArray_size, int *edgeArray,
                  int edgeArray_size, int *outDegrees, int outDegree_size,
                  double alpha, double eps) {

  FILE* fptr;
  fptr = fopen("pr_cuda_res.txt", "w");
  int n = vertexArray_size - 1;
  double *res = (double *)malloc((n + 1) * sizeof(double));
  double *pagerank_vector = (double *)malloc(n * sizeof(double));

  for (int i = 0; i < n; i++) {
    pagerank_vector[i] = 1 / double(n);
  }

  double *new_pagerank = (double *)malloc(n * sizeof(double));

  int *nodes = (int *)malloc(n * sizeof(int));
  for (int i = 0; i < n; i++) {
    nodes[i] = i;
  }

  double *delta = (double *)malloc(n * sizeof(double));
  int iters = 0;

  int *_nodes, *_vertexArray, *_edgeArray, *_outDegrees;
  cudaMalloc((void **)&_nodes, n * sizeof(int));
  cudaMalloc((void **)&_vertexArray, vertexArray_size * sizeof(int));
  cudaMalloc((void **)&_edgeArray, edgeArray_size * sizeof(int));
  cudaMalloc((void **)&_outDegrees, outDegree_size * sizeof(int));

  double *_pagerank_vector, *_new_pagerank;
  cudaMalloc((void **)&_pagerank_vector, n * sizeof(double));
  cudaMalloc((void **)&_new_pagerank, n * sizeof(double));

  double *_deltaSum, _leak, _alpha;
  cudaMalloc((void **)&_deltaSum, sizeof(double));
  cudaMalloc((void **)&_leak, sizeof(double));
  cudaMalloc((void **)&_alpha, sizeof(double));

  int _n;
  cudaMalloc((void **)&_n, sizeof(n));

  while (1) {
    iters++;
    double deltaSum = 0;
    double leak = 0;

    for (int i = 0; i < n; i++) {
      if (outDegrees[i] == 0) {
        leak += pagerank_vector[i];
      }
    }

    leak *= alpha;

    // memcpy go brrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr
    cudaMemcpy(_nodes, nodes, n * sizeof(int), cudaMemcpyHostToDevice);
    cudaMemcpy(_vertexArray, vertexArray, vertexArray_size * sizeof(int),
               cudaMemcpyHostToDevice);
    cudaMemcpy(_edgeArray, edgeArray, edgeArray_size * sizeof(int),
               cudaMemcpyHostToDevice);
    cudaMemcpy(_outDegrees, outDegrees, outDegree_size * sizeof(int),
               cudaMemcpyHostToDevice);
    cudaMemcpy(_pagerank_vector, pagerank_vector, n * sizeof(double),
               cudaMemcpyHostToDevice);
    cudaMemcpy(_new_pagerank, new_pagerank, n * sizeof(double),
               cudaMemcpyHostToDevice);
    cudaMemcpy(_deltaSum, &deltaSum, sizeof(double), cudaMemcpyHostToDevice);
    
    prGPU<<<1, 256>>>(_nodes, _pagerank_vector, _new_pagerank,
                              _vertexArray, _edgeArray, _outDegrees, alpha,
                              _deltaSum, n, leak);
    cudaDeviceSynchronize();

    // memcpy go brrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr
    cudaMemcpy(nodes, _nodes, n * sizeof(int), cudaMemcpyDeviceToHost);
    cudaMemcpy(vertexArray, _vertexArray, vertexArray_size * sizeof(int),
               cudaMemcpyDeviceToHost);
    cudaMemcpy(edgeArray, _edgeArray, edgeArray_size * sizeof(int),
               cudaMemcpyDeviceToHost);
    cudaMemcpy(outDegrees, _outDegrees, outDegree_size * sizeof(int),
               cudaMemcpyDeviceToHost);
    cudaMemcpy(pagerank_vector, _pagerank_vector, n * sizeof(double),
               cudaMemcpyDeviceToHost);
    cudaMemcpy(new_pagerank, _new_pagerank, n * sizeof(double),
               cudaMemcpyDeviceToHost);
    cudaMemcpy(&deltaSum, _deltaSum, sizeof(double), cudaMemcpyDeviceToHost);
    memcpy(pagerank_vector, new_pagerank, n * sizeof(double));

    if (deltaSum < eps) {
      break;
    }
  }

  cudaFree(_nodes);
  cudaFree(_vertexArray);
  cudaFree(_edgeArray);
  cudaFree(_outDegrees);
  cudaFree(_pagerank_vector);
  cudaFree(_new_pagerank);
  cudaFree(_deltaSum);

  double norm = 0;
  for (int i = 0; i < n; i++) {
    norm += i;
  }

  for (int i = 0; i < n; i++) {
    pagerank_vector[i] /= norm;
  }

  fprintf(fptr, "Latest pagerank is:\n");
  for (int i = 0; i < n; i++) {
    fprintf(fptr, "%lf\n", new_pagerank[i]);
  }
  fprintf(fptr, "\nIterations = \n%d\n", iters);
  fclose(fptr);
  return;
 }
}
