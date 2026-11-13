package graph

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/heap"
)

// Dijkstra returns the shortest-path distance from source to every vertex of a graph whose weights are
// all non-negative, and each vertex's parent on one shortest path (follow parent from v back to source to
// read the path). A vertex with no path has dist Unreachable and parent -1; parent[source] is -1.
// Among equal-length paths the parent is the first one found (deterministic for the same graph).
// Errors: ErrNilGraph; ErrVertexOutOfRange when source is not a vertex; ErrNegativeWeight when any edge
// of g (reachable or not) has a negative weight; ErrOverflow when some vertex reachable from source has a
// shortest distance of math.MaxInt64 (Unreachable) or more. A path whose sum overflows while a shorter
// path to the same vertex exists is simply not an improvement, never an error.
//
// ID:      ALG-GRAPH-012
// Time:    O((V + E) log E)   V = VertexCount(), E = EdgeCount(); a binary heap with lazy deletion (a vertex is pushed once per improvement and stale entries are skipped), and log E ≤ 2 log V
// Space:   O(V + E) extra (the heap holds at most one entry per edge relaxation; O(V) more once a path sum overflows), O(V) for the results
// Source:  Dijkstra, "A note on two problems in connexion with graphs", Numerische Mathematik, 1959, 1(1):269–271, doi:10.1007/BF01386390 (problem 2)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §24.3 (DIJKSTRA; binary-heap bound)
// Tags:    graph, dijkstra, shortest-path, weighted, single-source, heap
func Dijkstra(g *Graph, source int) (dist []int64, parent []int, err error) {
	if g == nil {
		return nil, nil, ErrNilGraph
	}
	if err := g.checkVertex("source", source); err != nil {
		return nil, nil, err
	}
	for v, list := range g.adj {
		for _, a := range list {
			if a.weight < 0 {
				return nil, nil, fmt.Errorf("%w: %d→%d weight %d", ErrNegativeWeight, v, a.to, a.weight)
			}
		}
	}
	n := len(g.adj)
	dist, parent = make([]int64, n), make([]int, n)
	for v := range dist {
		dist[v], parent[v] = Unreachable, -1
	}
	type entry struct {
		v int
		d int64
	}
	pq, _ := heap.NewHeap(func(a, b entry) bool { return a.d < b.d }, n)
	var tooFar []bool // allocated only when some path sum does not fit
	dist[source] = 0
	pq.Push(entry{source, 0})
	for pq.Len() > 0 {
		top, _ := pq.Pop()
		if top.d > dist[top.v] {
			continue // stale: a shorter distance to top.v was already settled
		}
		for _, a := range g.adj[top.v] {
			nd, ok := addWeight(top.d, a.weight)
			if !ok {
				// Weights are ≥ 0, so the sum is ≥ math.MaxInt64: never shorter than a distance already
				// found. Remember the vertex: if no representable path reaches it, that is ErrOverflow.
				if tooFar == nil {
					tooFar = make([]bool, n)
				}
				tooFar[a.to] = true
				continue
			}
			if nd < dist[a.to] {
				dist[a.to], parent[a.to] = nd, top.v
				pq.Push(entry{a.to, nd})
			}
		}
	}
	for v, far := range tooFar {
		if far && dist[v] == Unreachable {
			return nil, nil, fmt.Errorf("%w: shortest distance to %d is math.MaxInt64 or more", ErrOverflow, v)
		}
	}
	return dist, parent, nil
}
