package graph

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/CrimsonCoder42/calcbench/algo/unionfind"
)

// Kruskal returns a minimum spanning forest of an undirected graph (a minimum spanning tree when the graph
// is connected) and its total weight. Edges are taken in increasing weight (ties in the order of g.Edges())
// and kept when they join two different trees of an algo/unionfind DisjointSet. Each returned edge has
// From <= To. Negative weights are allowed. g is not modified.
// Errors: ErrNilGraph; ErrDirected when g is directed; ErrOverflow when the total weight does not fit in int64.
//
// ID:      ALG-GRAPH-018
// Time:    O(V + E log E)   V = VertexCount(), E = EdgeCount() (the sort dominates; the unions are O(E α(V)))
// Space:   O(V + E) extra (a sorted copy of the edges and the disjoint-set forest), O(V) for the result
// Source:  Kruskal, "On the shortest spanning subtree of a graph and the traveling salesman problem", Proceedings of the American Mathematical Society, 1956, 7(1):48–50, doi:10.1090/S0002-9939-1956-0078686-7
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §23.2 (MST-KRUSKAL)
// Tags:    graph, mst, minimum-spanning-tree, kruskal, greedy, union-find
func Kruskal(g *Graph) (tree []Edge, totalWeight int64, err error) {
	if g == nil {
		return nil, 0, ErrNilGraph
	}
	if g.directed {
		return nil, 0, ErrDirected
	}
	n := len(g.adj)
	edges := g.Edges()
	slices.SortStableFunc(edges, func(a, b Edge) int { return cmp.Compare(a.Weight, b.Weight) })
	sets, err := unionfind.NewDisjointSet(n)
	if err != nil {
		return nil, 0, err // unreachable: n >= 0
	}
	tree = make([]Edge, 0, max(n-1, 0))
	for _, e := range edges {
		if len(tree) == n-1 {
			break
		}
		merged, err := sets.Union(e.From, e.To)
		if err != nil {
			return nil, 0, err // unreachable: edges come from g
		}
		if !merged {
			continue
		}
		sum, ok := checkedSum(totalWeight, e.Weight)
		if !ok {
			return nil, 0, fmt.Errorf("%w: tree weight after edge %d–%d", ErrOverflow, e.From, e.To)
		}
		totalWeight = sum
		tree = append(tree, e)
	}
	return tree, totalWeight, nil
}
