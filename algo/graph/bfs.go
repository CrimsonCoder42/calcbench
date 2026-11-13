package graph

import "iter"

// BFS returns the vertices reachable from source in breadth-first order: source, then every vertex at
// distance 1, then distance 2, and so on; within a level, in the order their edges were added. Each range
// over the sequence runs a fresh search; breaking out of the loop stops it. The graph is read while the
// sequence runs. Edge weights are ignored.
// Errors (returned before any iteration): ErrNilGraph; ErrVertexOutOfRange when source is not a vertex.
//
// ID:      ALG-GRAPH-009
// Time:    O(V + E) for a full iteration   V = VertexCount(), E = EdgeCount()
// Space:   O(V) extra per iteration (a visited array and the queue; two allocations)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (BFS)
// Source:  The Go Programming Language Specification, "For statements with range clause" (range over function iterators); package iter (iter.Seq)
// Tags:    graph, bfs, breadth-first, traversal, iterator
func BFS(g *Graph, source int) (iter.Seq[int], error) {
	if g == nil {
		return nil, ErrNilGraph
	}
	if err := g.checkVertex("source", source); err != nil {
		return nil, err
	}
	return func(yield func(int) bool) {
		n := len(g.adj)
		seen := make([]bool, n)
		queue := make([]int, 1, n)
		queue[0] = source
		seen[source] = true
		for head := 0; head < len(queue); head++ {
			v := queue[head]
			if !yield(v) {
				return
			}
			for _, a := range g.adj[v] {
				if !seen[a.to] {
					seen[a.to] = true
					queue = append(queue, a.to)
				}
			}
		}
	}, nil
}
