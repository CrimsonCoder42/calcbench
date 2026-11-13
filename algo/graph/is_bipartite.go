package graph

// IsBipartite reports whether the vertices of an undirected graph split into two sides with every edge
// joining the two sides, and returns the split: side[v] is false for one side and true for the other,
// with the smallest vertex of each connected component on the false side. When the graph is not
// bipartite (an odd cycle, a self-loop included) side is nil and ok is false. It 2-colours each component
// by breadth-first search.
// Errors: ErrNilGraph; ErrDirected when g is directed.
//
// ID:      ALG-GRAPH-017
// Time:    O(V + E)   V = VertexCount(), E = EdgeCount()
// Space:   O(V) extra (colour and queue arrays), O(V) for the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2, Exercise 22.2-7 (split into two kinds with every rivalry across, by BFS, O(n + r))
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §4.1 (two-colorability: is a graph bipartite?)
// Tags:    graph, bipartite, two-coloring, bfs, odd-cycle
func IsBipartite(g *Graph) (side []bool, ok bool, err error) {
	if g == nil {
		return nil, false, ErrNilGraph
	}
	if g.directed {
		return nil, false, ErrDirected
	}
	n := len(g.adj)
	const unseen, left, right = 0, 1, 2
	colour := make([]uint8, n)
	queue := make([]int, 0, n)
	for s := range n {
		if colour[s] != unseen {
			continue
		}
		colour[s] = left
		queue = append(queue[:0], s)
		for head := 0; head < len(queue); head++ {
			v := queue[head]
			other := uint8(left + right - int(colour[v]))
			for _, a := range g.adj[v] {
				switch colour[a.to] {
				case unseen:
					colour[a.to] = other
					queue = append(queue, a.to)
				case colour[v]:
					return nil, false, nil
				}
			}
		}
	}
	side = make([]bool, n)
	for v, c := range colour {
		side[v] = c == right
	}
	return side, true, nil
}
