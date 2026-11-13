package graph

// ShortestPathUnweighted returns a path with the fewest edges from source to target, as the list of
// vertices source … target, found by breadth-first search (edge weights are ignored). Among equally short
// paths it returns the one BFS discovers first (deterministic for the same graph). source == target gives
// [source]; no path gives nil with a nil error.
// Errors: ErrNilGraph; ErrVertexOutOfRange when source or target is not a vertex.
//
// ID:      ALG-GRAPH-011
// Time:    O(V + E)   V = VertexCount(), E = EdgeCount(); the search stops when target is reached
// Space:   O(V) extra (parent array and queue), O(p) for the result, p = the path's vertex count
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (Theorem 22.5: BFS finds shortest-path distances; PRINT-PATH)
// Tags:    graph, bfs, shortest-path, unweighted, path-reconstruction
func ShortestPathUnweighted(g *Graph, source, target int) ([]int, error) {
	if g == nil {
		return nil, ErrNilGraph
	}
	if err := g.checkVertex("source", source); err != nil {
		return nil, err
	}
	if err := g.checkVertex("target", target); err != nil {
		return nil, err
	}
	if source == target {
		return []int{source}, nil
	}
	n := len(g.adj)
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -1
	}
	parent[source] = source
	queue := make([]int, 1, n)
	queue[0] = source
	found := false
	for head := 0; head < len(queue) && !found; head++ {
		v := queue[head]
		for _, a := range g.adj[v] {
			if parent[a.to] == -1 {
				parent[a.to] = v
				if a.to == target {
					found = true
					break
				}
				queue = append(queue, a.to)
			}
		}
	}
	if !found {
		return nil, nil
	}
	length := 1
	for v := target; v != source; v = parent[v] {
		length++
	}
	path := make([]int, length)
	for i, v := length-1, target; i >= 0; i-- {
		path[i] = v
		v = parent[v]
	}
	return path, nil
}
