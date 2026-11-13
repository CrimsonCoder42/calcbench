package graph

// CloneGraph returns a deep copy of g: the same vertices, edges (in the same list order), weights and
// direction, sharing no memory with g, so AddEdge on either leaves the other unchanged. A nil g gives nil.
//
// ID:      ALG-GRAPH-020
// Time:    O(V + E)   V = VertexCount(), E = EdgeCount()
// Space:   O(V + E) for the copy (two allocations plus the Graph value)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1 (adjacency lists; the copy duplicates every list)
// Tags:    graph, clone, deep-copy, adjacency-list
func CloneGraph(g *Graph) *Graph {
	if g == nil {
		return nil
	}
	total := 0
	for _, list := range g.adj {
		total += len(list)
	}
	backing := make([]arc, total)
	c := &Graph{adj: make([][]arc, len(g.adj)), edges: g.edges, directed: g.directed}
	start := 0
	for v, list := range g.adj {
		end := start + copy(backing[start:], list)
		c.adj[v] = backing[start:end:end]
		start = end
	}
	return c
}
