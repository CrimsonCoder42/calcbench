// Package graph holds graph algorithms over an adjacency-list Graph on the integer vertices 0..n-1
// (directed or undirected, int64 weights): breadth- and depth-first traversal, shortest paths (BFS,
// Dijkstra, Bellman-Ford, Floyd-Warshall), topological order, connected components, bipartiteness,
// minimum spanning trees (Kruskal, Prim), deep copy, word ladders, and searches on [][]T grids.
//
// ID short form: GRAPH (IDs are ALG-GRAPH-NNN).
//
// Conventions: weights are exact int64 values in the caller's unit; every path or tree sum is checked
// and a sum that does not fit is ErrOverflow. A vertex with no path has distance Unreachable
// (math.MaxInt64) and parent -1. Parallel edges and self-loops are allowed. Traversals return iter.Seq
// values: each range runs a fresh search and may stop early; argument errors come back before iteration.
// ConnectedComponents and Kruskal use algo/unionfind; Dijkstra and Prim use algo/heap with lazy deletion.
// Grids are rectangular [][]T with Cell positions (Cell is unionfind.Cell) and FourNeighbors or
// EightNeighbors adjacency. Only the functions named …InPlace modify their input (the grid); no other
// function modifies an argument. A Graph is not safe for concurrent use while AddEdge runs.
//
// Functions:
//
//	ALG-GRAPH-001  NewDirected  a directed graph from an edge list
//	ALG-GRAPH-002  NewUndirected  an undirected graph from an edge list
//	ALG-GRAPH-003  Graph.AddEdge  add one edge
//	ALG-GRAPH-004  Graph.VertexCount  the number of vertices
//	ALG-GRAPH-005  Graph.EdgeCount  the number of edges
//	ALG-GRAPH-006  Graph.Directed  whether the graph is directed
//	ALG-GRAPH-007  Graph.Neighbors  a vertex's (neighbour, weight) pairs as an iter.Seq2
//	ALG-GRAPH-008  Graph.Edges  every edge once, as a slice
//	ALG-GRAPH-009  BFS  vertices reachable from a source, breadth-first, as an iter.Seq
//	ALG-GRAPH-010  DFS  vertices reachable from a source, depth-first preorder, as an iter.Seq
//	ALG-GRAPH-011  ShortestPathUnweighted  a fewest-edges path between two vertices (BFS)
//	ALG-GRAPH-012  Dijkstra  single-source shortest paths, non-negative weights
//	ALG-GRAPH-013  BellmanFord  single-source shortest paths, any weights, negative-cycle check
//	ALG-GRAPH-014  FloydWarshall  all-pairs shortest-path distances
//	ALG-GRAPH-015  TopologicalSort  a dependency order of a DAG (Kahn), cycle detection
//	ALG-GRAPH-016  ConnectedComponents  a component label per vertex (union-find)
//	ALG-GRAPH-017  IsBipartite  a two-side split with every edge across, if one exists
//	ALG-GRAPH-018  Kruskal  a minimum spanning forest (sorted edges, union-find)
//	ALG-GRAPH-019  Prim  a minimum spanning forest (grown with a heap)
//	ALG-GRAPH-020  CloneGraph  a deep copy of a graph
//	ALG-GRAPH-021  WordLadder  a shortest one-letter-at-a-time word ladder
//	ALG-GRAPH-022  GridBFS  grid cells reachable from a start, breadth-first, as an iter.Seq
//	ALG-GRAPH-023  GridDFS  grid cells reachable from a start, depth-first preorder, as an iter.Seq
//	ALG-GRAPH-024  NumberOfIslands  the number of connected land groups on a grid
//	ALG-GRAPH-025  FloodFillInPlace  repaint the region of one value around a cell (modifies the grid)
//	ALG-GRAPH-026  SurroundedRegionsInPlace  fill regions of one value that cannot reach the border (modifies the grid)
package graph
