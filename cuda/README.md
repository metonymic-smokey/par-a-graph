# Datasets

## Data Format

The graph must be in two files:
 - Edges file: each line must be the (comma-separated) value of each (directed) edge. An **edge list**. Ex: `1,2` means an edge from node 1 to node 2 (the numbers can be any valid 64-bit unsigned integer).
 - Nodes file: each line must be comma-separated with two values - node number and a name (or number) for the node. This **must enumerate all possible nodes**. Ex: `1,b` means that the node 1 is given the name `b` (the name does not matter).

Although, note that this can be modified to accept any kind of graph inputs by changing the appropriate loading functions. (For pagerank - `readGraph` in [`main.go`](../cpu/main.go)).

## Pre-processing

Some of the datasets may have to be pre-processed to get the two files (as described above). For Wiki Vote and Stanford graph datasets, we provide a pre-processing script - [`process.py`](./process.py).

## List of Datasets

### Example/Small graph

Files: [`example`](./example) and [`examplePageNum`](./examplePageNum).  
Description: This is only to show the file format and for basic tests.  
Nodes: 17  
Edges: 11  
Name (in code): `Small`

### Large Graph Dataset (English Wikipedia graph)

Files: [`dirLinks.txt`](./dirLinks.txt) and [`pageNum.txt`](./pageNum.txt)  
Link: [Wiki](http://cfinder.org/wiki/?n=Main.Data#toc1)  
Description: Network of pages in the English Wikipedia.  
Nodes: 2080370  
Edges: 46092177  
Name (in code): `Large`

### Quora Question Pairs

Files: [`quora_edges.txt`](./quora_edges.txt) and [`quora_nodes.txt`](./quora_nodes.txt)  
Link: [Quora question pairs](https://www.kaggle.com/c/quora-question-pairs/data)  
Description: Quora question pairs dataset.  
Nodes: 537935  
Edges: 404291  
Name (in code): `Quora`

### Wiki Vote Graph

Files: [`wiki-vote-edges.txt`](./wiki-vote-edges.txt) and [`wiki-vote-nodes.txt`](./wiki-vote-nodes.txt)  
Link: [Wikipedia dataset](http://snap.stanford.edu/data/wiki-Vote.html)  
Description: The network contains all the Wikipedia voting data from the inception of Wikipedia till January 2008.  
Nodes: 7115  
Edges: 103689  
Name (in code): `Wiki`

### Stanford Web Graph

Files: [`stanford-edges.txt`](./stanford-edges.txt) and [`stanford-nodes.txt`](./stanford-nodes.txt)  
Links: [Stanford Web](http://snap.stanford.edu/data/web-Stanford.html)  
Description: Pages from Stanford university (as of 2002) and links between them.  
Nodes: 281903  
Edges: 2312497  
Name (in code): `Stanford`
