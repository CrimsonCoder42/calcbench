package graph

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/heap"
)

// Prim returns a minimum spanning forest of an undirected graph (a minimum spanning tree when the graph is
// connected) and its total weight. It grows a tree from vertex 0, always adding the lightest edge that
// leaves the tree, then restarts from the smallest vertex not yet reached. Each returned edge has From in
// the tree and To the vertex it added, in the order added. Negative weights are allowed. With equal
// weights the edges may differ from Kruskal's; the total weight is the same.
// Errors: ErrNilGraph; ErrDirected when g is directed; ErrOverflow when the total weight does not fit in int64.
//
// ID:      ALG-GRAPH-019
// Time:    O(V + E log E)   V = VertexCount(), E = EdgeCount(); a binary heap with lazy deletion (stale edges into the tree are skipped when popped), and log E ≤ 2 log V
// Space:   O(V + E) extra (the heap holds at most 2E edges; a visited array), O(V) for the result
// Source:  Prim, "Shortest connection networks and some generalizations", Bell System Technical Journal, 1957, 36(6):1389–1401, doi:10.1002/j.1538-7305.1957.tb01515.x
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §23.2 (MST-PRIM)
// Tags:    graph, mst, minimum-spanning-tree, prim, greedy, heap
func Prim(g *Graph) (tree []Edge, totalWeight int64, err error) {
	if g == nil {
		return nil, 0, ErrNilGraph
	}
	if g.directed {
		return nil, 0, ErrDirected
	}
	n := len(g.adj)
	inTree := make([]bool, n)
	pq, _ := heap.NewHeap(func(a, b Edge) bool { return a.Weight < b.Weight }, n)
	tree = make([]Edge, 0, max(n-1, 0))
	grow := func(v int) {
		inTree[v] = true
		for _, a := range g.adj[v] {
			if !inTree[a.to] {
				pq.Push(Edge{v, a.to, a.weight})
			}
		}
	}
	for root := range n {
		if inTree[root] {
			continue
		}
		grow(root)
		for pq.Len() > 0 {
			e, _ := pq.Pop()
			if inTree[e.To] {
				continue // stale: e.To joined the tree through a lighter edge
			}
			sum, ok := checkedSum(totalWeight, e.Weight)
			if !ok {
				return nil, 0, fmt.Errorf("%w: tree weight after edge %d–%d", ErrOverflow, e.From, e.To)
			}
			totalWeight = sum
			tree = append(tree, e)
			grow(e.To)
		}
	}
	return tree, totalWeight, nil
}
