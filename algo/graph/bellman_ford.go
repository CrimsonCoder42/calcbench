package graph

import "fmt"

// BellmanFord returns the shortest-path distance from source to every vertex of a graph whose weights
// may be negative, and each vertex's parent on one shortest path. A vertex with no path has dist
// Unreachable and parent -1; parent[source] is -1. In an undirected graph a negative edge is itself a
// negative cycle (walk it back and forth).
// Errors: ErrNilGraph; ErrVertexOutOfRange when source is not a vertex; ErrNegativeCycle when a cycle of
// negative total weight is reachable from source; ErrOverflow when a tentative distance falls below
// math.MinInt64 (this can be reported before a negative cycle that causes it is detected), or when some
// vertex reachable from source has a shortest distance of math.MaxInt64 (Unreachable) or more. A path whose
// sum overflows upward while a shorter path to the same vertex exists is not an improvement, never an error.
//
// ID:      ALG-GRAPH-013
// Time:    O(V·(V + E))   V = VertexCount(), E = EdgeCount(); at most V passes over all edges, stopping early after a pass that changes nothing
// Space:   O(1) extra (O(V) once a path sum overflows upward), O(V) for the results
// Source:  Bellman, "On a routing problem", Quarterly of Applied Mathematics, 1958, 16(1):87–90, doi:10.1090/qam/102435
// Source:  Ford, "Network flow theory", RAND Corporation, paper P-923, 1956
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §24.1 (BELLMAN-FORD with the negative-cycle check)
// Tags:    graph, bellman-ford, shortest-path, negative-weights, negative-cycle, single-source
func BellmanFord(g *Graph, source int) (dist []int64, parent []int, err error) {
	if g == nil {
		return nil, nil, ErrNilGraph
	}
	if err := g.checkVertex("source", source); err != nil {
		return nil, nil, err
	}
	n := len(g.adj)
	dist, parent = make([]int64, n), make([]int, n)
	for v := range dist {
		dist[v], parent[v] = Unreachable, -1
	}
	var tooFar []bool // allocated only when some path sum does not fit
	dist[source] = 0
	// Passes 1..n-1 relax; pass n only checks. A change in pass n means a reachable negative cycle.
	for pass := 1; pass <= n; pass++ {
		changed := false
		for u, list := range g.adj {
			if dist[u] == Unreachable {
				continue
			}
			for _, a := range list {
				nd, ok := addWeight(dist[u], a.weight)
				if !ok {
					if a.weight < 0 {
						// Below math.MinInt64: a real improvement that int64 cannot hold.
						return nil, nil, fmt.Errorf("%w: distance to %d via %d is below math.MinInt64", ErrOverflow, a.to, u)
					}
					// At or above math.MaxInt64: never shorter than a distance already found. Remember the
					// vertex: if no representable path reaches it, that is ErrOverflow.
					if tooFar == nil {
						tooFar = make([]bool, n)
					}
					tooFar[a.to] = true
					continue
				}
				if nd < dist[a.to] {
					if pass == n {
						return nil, nil, fmt.Errorf("%w: reachable from source %d (edge %d→%d still relaxes)", ErrNegativeCycle, source, u, a.to)
					}
					dist[a.to], parent[a.to] = nd, u
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}
	for v, far := range tooFar {
		if far && dist[v] == Unreachable {
			return nil, nil, fmt.Errorf("%w: shortest distance to %d is math.MaxInt64 or more", ErrOverflow, v)
		}
	}
	return dist, parent, nil
}
