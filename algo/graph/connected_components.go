package graph

import "github.com/CrimsonCoder42/calcbench/algo/unionfind"

// ConnectedComponents returns, for each vertex, the number of its connected component, and the number of
// components. Components are numbered 0..count-1 in order of their smallest vertex. For a directed graph
// the components are the weakly connected ones (edge direction is ignored). It unions the endpoints of
// every edge in an algo/unionfind DisjointSet.
// Errors: ErrNilGraph.
//
// ID:      ALG-GRAPH-016
// Time:    O(V + E α(V))   V = VertexCount(), E = EdgeCount(), α = the inverse Ackermann function
// Space:   O(V) extra (the disjoint-set forest and a label per root), O(V) for the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.1 (CONNECTED-COMPONENTS with disjoint sets)
// Source:  Tarjan, "Efficiency of a good but not linear set union algorithm", Journal of the ACM, 1975, 22(2):215–225, doi:10.1145/321879.321884
// Tags:    graph, connected-components, weakly-connected, union-find
func ConnectedComponents(g *Graph) (component []int, count int, err error) {
	if g == nil {
		return nil, 0, ErrNilGraph
	}
	n := len(g.adj)
	sets, err := unionfind.NewDisjointSet(n)
	if err != nil {
		return nil, 0, err // unreachable: n = len(g.adj) >= 0
	}
	for u, list := range g.adj {
		for _, a := range list {
			if _, err := sets.Union(u, a.to); err != nil {
				return nil, 0, err // unreachable: every arc's endpoints are vertices of g
			}
		}
	}
	// rootLabel[r] is the label of the component whose root is r, -1 until its smallest vertex is seen.
	rootLabel := make([]int, n)
	for v := range rootLabel {
		rootLabel[v] = -1
	}
	component = make([]int, n)
	for v := range n {
		root, _ := sets.Find(v)
		if rootLabel[root] == -1 {
			rootLabel[root] = count
			count++
		}
		component[v] = rootLabel[root]
	}
	return component, count, nil
}
