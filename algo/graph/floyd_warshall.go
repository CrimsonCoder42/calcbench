package graph

import "fmt"

// FloydWarshall returns the matrix of shortest-path distances between every pair of vertices: d[i][j]
// is the least total weight of a path from i to j, Unreachable when there is none, and d[i][i] = 0.
// Weights may be negative. The rows share one backing array of V² values; the caller bounds V (V = 10⁴
// is 800 MB).
// Errors: ErrNilGraph; ErrNegativeCycle when any cycle has negative total weight (for an undirected graph:
// any negative edge); ErrOverflow when a distance does not fit in int64; ErrTooLarge when the V² table
// (8 bytes per entry) would exceed 2⁴⁷ bytes (V > 2²², about 4.2 million vertices; 2³¹ bytes and V > 2¹⁴
// on 32-bit).
//
// ID:      ALG-GRAPH-014
// Time:    Θ(V³ + E)   V = VertexCount(), E = EdgeCount()
// Space:   O(1) extra, Θ(V²) for the result (two allocations)
// Source:  Floyd, "Algorithm 97: Shortest path", Communications of the ACM, 1962, 5(6):345, doi:10.1145/367766.368168
// Source:  Warshall, "A theorem on Boolean matrices", Journal of the ACM, 1962, 9(1):11–12, doi:10.1145/321105.321107
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §25.2 (FLOYD-WARSHALL)
// Tags:    graph, floyd-warshall, all-pairs, shortest-path, dense, dynamic-programming
func FloydWarshall(g *Graph) ([][]int64, error) {
	if g == nil {
		return nil, ErrNilGraph
	}
	n := len(g.adj)
	if !distanceTableFits(n) {
		return nil, fmt.Errorf("%w: V = %d, V² distances", ErrTooLarge, n)
	}
	backing := make([]int64, n*n)
	d := make([][]int64, n)
	for i := range d {
		row := backing[i*n : (i+1)*n : (i+1)*n]
		for j := range row {
			row[j] = Unreachable
		}
		row[i] = 0
		d[i] = row
	}
	for u, list := range g.adj {
		for _, a := range list {
			if a.weight < d[u][a.to] {
				d[u][a.to] = a.weight
			}
		}
	}
	for i := range d {
		if d[i][i] < 0 {
			return nil, fmt.Errorf("%w: self-loop at %d", ErrNegativeCycle, i)
		}
	}
	for k := range n {
		dk := d[k]
		for i := range n {
			di := d[i]
			dik := di[k]
			if dik == Unreachable {
				continue
			}
			for j, dkj := range dk {
				if dkj == Unreachable {
					continue
				}
				s, ok := addWeight(dik, dkj)
				if !ok {
					return nil, fmt.Errorf("%w: path %d→%d through %d", ErrOverflow, i, j, k)
				}
				if s < di[j] {
					if i == j {
						return nil, fmt.Errorf("%w: through vertices %d and %d", ErrNegativeCycle, i, k)
					}
					di[j] = s
				}
			}
		}
	}
	return d, nil
}

// distanceTableFits reports whether an n×n table of int64 fits in maxAllocBytes (L41).
func distanceTableFits(n int) bool {
	return n == 0 || uint64(n) <= maxAllocBytes/8/uint64(n)
}
