package graph

import "fmt"

// TopologicalSort returns the vertices of a directed acyclic graph in an order where every edge u → v has
// u before v (Kahn's method: repeatedly output a vertex with no remaining incoming edges). The vertices
// with no incoming edges start the queue in increasing order, then the queue is first-in first-out, so the
// order is deterministic for the same graph.
// Errors: ErrNilGraph; ErrUndirected when g is undirected; ErrCycle when g has a directed cycle
// (a self-loop included).
//
// ID:      ALG-GRAPH-015
// Time:    O(V + E)   V = VertexCount(), E = EdgeCount()
// Space:   O(V) extra (the in-degree array; the result doubles as the queue), O(V) for the result
// Source:  Kahn, "Topological sorting of large networks", Communications of the ACM, 1962, 5(11):558–562, doi:10.1145/368996.369025
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.4 (topological sort; Exercise 22.4-5 gives this in-degree method)
// Tags:    graph, topological-sort, kahn, dag, cycle-detection, scheduling, dependency-order
func TopologicalSort(g *Graph) ([]int, error) {
	if g == nil {
		return nil, ErrNilGraph
	}
	if !g.directed {
		return nil, ErrUndirected
	}
	n := len(g.adj)
	indegree := make([]int, n)
	for _, list := range g.adj {
		for _, a := range list {
			indegree[a.to]++
		}
	}
	order := make([]int, 0, n)
	for v, d := range indegree {
		if d == 0 {
			order = append(order, v)
		}
	}
	for head := 0; head < len(order); head++ {
		for _, a := range g.adj[order[head]] {
			indegree[a.to]--
			if indegree[a.to] == 0 {
				order = append(order, a.to)
			}
		}
	}
	if len(order) < n {
		return nil, fmt.Errorf("%w: %d of %d vertices lie on or behind a cycle", ErrCycle, n-len(order), n)
	}
	return order, nil
}
