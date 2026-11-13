package graph

import "iter"

// DFS returns the vertices reachable from source in depth-first preorder (each vertex when it is first
// discovered), following edges in the order they were added: the same order as the recursive DFS-VISIT,
// computed with an explicit stack, so deep graphs cannot overflow the call stack. Each range over the
// sequence runs a fresh search; breaking out of the loop stops it. The graph is read while the sequence
// runs. Edge weights are ignored.
// Errors (returned before any iteration): ErrNilGraph; ErrVertexOutOfRange when source is not a vertex.
//
// ID:      ALG-GRAPH-010
// Time:    O(V + E) for a full iteration   V = VertexCount(), E = EdgeCount()
// Space:   O(V) extra per iteration (a visited array and a stack of at most V frames)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.3 (DFS-VISIT)
// Source:  Tarjan, "Depth-first search and linear graph algorithms", SIAM Journal on Computing, 1972, 1(2):146–160, doi:10.1137/0201010
// Source:  The Go Programming Language Specification, "For statements with range clause" (range over function iterators); package iter (iter.Seq)
// Tags:    graph, dfs, depth-first, preorder, traversal, iterator
func DFS(g *Graph, source int) (iter.Seq[int], error) {
	if g == nil {
		return nil, ErrNilGraph
	}
	if err := g.checkVertex("source", source); err != nil {
		return nil, err
	}
	return func(yield func(int) bool) {
		// frame is a vertex on the current path and the index of its next arc to try.
		type frame struct{ v, next int }
		seen := make([]bool, len(g.adj))
		seen[source] = true
		if !yield(source) {
			return
		}
		stack := []frame{{source, 0}}
		for len(stack) > 0 {
			top := len(stack) - 1
			arcs := g.adj[stack[top].v]
			if stack[top].next >= len(arcs) {
				stack = stack[:top]
				continue
			}
			w := arcs[stack[top].next].to
			stack[top].next++
			if seen[w] {
				continue
			}
			seen[w] = true
			if !yield(w) {
				return
			}
			stack = append(stack, frame{w, 0})
		}
	}, nil
}
