package graph

import (
	"errors"
	"fmt"
	"iter"
	"math"
	"strconv"
)

// Errors returned by this package. Each is wrapped with the offending parameter.
var (
	// ErrNilGraph: a *Graph argument or receiver is nil.
	ErrNilGraph = errors.New("graph: graph is nil")
	// ErrNegativeVertexCount: the vertex count n is less than 0.
	ErrNegativeVertexCount = errors.New("graph: vertex count must not be negative")
	// ErrVertexOutOfRange: a vertex is outside 0..n-1.
	ErrVertexOutOfRange = errors.New("graph: vertex out of range")
	// ErrNegativeWeight: Dijkstra found an edge of negative weight.
	ErrNegativeWeight = errors.New("graph: negative edge weight")
	// ErrNegativeCycle: a cycle of negative total weight makes shortest paths undefined.
	ErrNegativeCycle = errors.New("graph: negative-weight cycle")
	// ErrOverflow: a path or tree weight does not fit in int64 (or equals Unreachable).
	ErrOverflow = errors.New("graph: weight sum overflows int64")
	// ErrCycle: a directed graph has a cycle, so it has no topological order.
	ErrCycle = errors.New("graph: directed graph has a cycle")
	// ErrDirected: the function needs an undirected graph.
	ErrDirected = errors.New("graph: graph is directed")
	// ErrUndirected: the function needs a directed graph.
	ErrUndirected = errors.New("graph: graph is undirected")
	// ErrTooLarge: an allocation from a size (V list headers, V² distances, rows×cols cells) would exceed
	// maxAllocBytes.
	ErrTooLarge = errors.New("graph: size too large to allocate")
)

// maxAllocBytes is the largest single allocation this package makes from a size argument: 2⁴⁷ bytes on
// 64-bit, 2³¹ on 32-bit, half the Go runtime's maxAlloc, above which makeslice panics (runtime/malloc.go;
// lead decision L41).
const maxAllocBytes = 1 << (31 + 16*(^uint(0)>>63))

// Bytes per element of the arrays sized from a vertex or edge count: a slice header (three words) per
// adjacency list, an int degree per vertex, and an arc (an int and an int64, rounded up to 16 bytes).
const (
	wordBytes   = strconv.IntSize / 8
	headerBytes = 3 * wordBytes
	arcBytes    = 16
)

// Unreachable is the distance reported for a vertex with no path from the source (math.MaxInt64).
// No real distance equals it: a sum that reaches it is ErrOverflow.
const Unreachable int64 = math.MaxInt64

// Edge is an edge from From to To with an integer Weight (0 for an unweighted graph). Weights are
// exact int64 values in whatever unit the caller chooses (metres, cents, seconds).
type Edge struct {
	From, To int
	Weight   int64
}

// arc is one adjacency-list entry: the head vertex and the edge weight.
type arc struct {
	to     int
	weight int64
}

// Graph is a directed or undirected graph over the vertices 0..n-1 stored as adjacency lists.
// Parallel edges and self-loops are kept; an undirected edge is stored in both endpoints' lists (a
// self-loop once). Neighbours are visited in the order their edges were added.
//
// A Graph comes from NewDirected, NewUndirected or CloneGraph. Methods on a nil *Graph return zero
// values (AddEdge returns ErrNilGraph). A Graph is not safe for concurrent use while AddEdge runs.
type Graph struct {
	adj      [][]arc
	edges    int
	directed bool
}

// NewDirected returns the directed graph on vertices 0..n-1 with the given edges (From → To).
// Every edge is checked before anything is stored. edges is read, never modified.
// Errors: ErrNegativeVertexCount when n < 0; ErrVertexOutOfRange when an endpoint is outside 0..n-1;
// ErrTooLarge when the n list headers (24 bytes each on 64-bit) or the arcs (16 bytes per stored arc)
// would exceed 2⁴⁷ bytes (2³¹ on 32-bit).
//
// ID:      ALG-GRAPH-001
// Time:    O(V + E)   V = n, E = len(edges)
// Space:   O(V + E) for the graph (two allocations: the list headers and one shared arc array)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1 (adjacency-list representation)
// Tags:    graph, adjacency-list, directed, weighted
func NewDirected(n int, edges []Edge) (*Graph, error) {
	return build(n, edges, true)
}

// NewUndirected returns the undirected graph on vertices 0..n-1 with the given edges; each edge is
// stored in both endpoints' lists (a self-loop once). Every edge is checked before anything is stored.
// edges is read, never modified.
// Errors: ErrNegativeVertexCount when n < 0; ErrVertexOutOfRange when an endpoint is outside 0..n-1;
// ErrTooLarge when the n list headers (24 bytes each on 64-bit) or the arcs (16 bytes per stored arc)
// would exceed 2⁴⁷ bytes (2³¹ on 32-bit).
//
// ID:      ALG-GRAPH-002
// Time:    O(V + E)   V = n, E = len(edges)
// Space:   O(V + E) for the graph (two allocations: the list headers and one shared arc array)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1 (adjacency-list representation; an undirected edge appears in both lists)
// Tags:    graph, adjacency-list, undirected, weighted
func NewUndirected(n int, edges []Edge) (*Graph, error) {
	return build(n, edges, false)
}

// build sizes every list exactly (a degree count pass), then fills them from one backing array.
// Each list is clipped to its degree, so a later append reallocates that list alone.
func build(n int, edges []Edge, directed bool) (*Graph, error) {
	if n < 0 {
		return nil, fmt.Errorf("%w: n = %d", ErrNegativeVertexCount, n)
	}
	for i, e := range edges {
		if e.From < 0 || e.From >= n || e.To < 0 || e.To >= n {
			return nil, fmt.Errorf("%w: edges[%d] = %d→%d, n = %d", ErrVertexOutOfRange, i, e.From, e.To, n)
		}
	}
	if n > maxAllocBytes/headerBytes {
		return nil, fmt.Errorf("%w: n = %d vertices", ErrTooLarge, n)
	}
	// At most 2 arcs per edge; len(edges) ≤ maxAlloc/24 (an Edge is 24 bytes), so 2·len(edges) cannot overflow.
	if 2*len(edges) > maxAllocBytes/arcBytes {
		return nil, fmt.Errorf("%w: %d edges", ErrTooLarge, len(edges))
	}
	deg := make([]int, n)
	total := 0
	for _, e := range edges {
		deg[e.From]++
		total++
		if !directed && e.From != e.To {
			deg[e.To]++
			total++
		}
	}
	g := &Graph{adj: make([][]arc, n), edges: len(edges), directed: directed}
	backing := make([]arc, total)
	start := 0
	for v := range n {
		end := start + deg[v]
		g.adj[v] = backing[start:start:end]
		start = end
	}
	for _, e := range edges {
		g.adj[e.From] = append(g.adj[e.From], arc{e.To, e.Weight})
		if !directed && e.From != e.To {
			g.adj[e.To] = append(g.adj[e.To], arc{e.From, e.Weight})
		}
	}
	return g, nil
}

// AddEdge adds an edge from → to of the given weight (in both lists when g is undirected, once for a
// self-loop).
// Errors: ErrNilGraph when g is nil; ErrVertexOutOfRange when from or to is outside 0..VertexCount()-1.
//
// ID:      ALG-GRAPH-003
// Time:    O(1) amortized (the endpoint's list doubles when full)
// Space:   O(1) extra amortized
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1
// Tags:    graph, adjacency-list, add-edge
func (g *Graph) AddEdge(from, to int, weight int64) error {
	if g == nil {
		return ErrNilGraph
	}
	if err := g.checkVertex("from", from); err != nil {
		return err
	}
	if err := g.checkVertex("to", to); err != nil {
		return err
	}
	g.adj[from] = append(g.adj[from], arc{to, weight})
	if !g.directed && from != to {
		g.adj[to] = append(g.adj[to], arc{from, weight})
	}
	g.edges++
	return nil
}

// VertexCount returns the number of vertices n (0 for a nil *Graph).
//
// ID:      ALG-GRAPH-004
// Time:    O(1)
// Space:   O(1)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1 (|V|)
// Tags:    graph, order, vertex-count
func (g *Graph) VertexCount() int {
	if g == nil {
		return 0
	}
	return len(g.adj)
}

// EdgeCount returns the number of edges added (an undirected edge counts once; 0 for a nil *Graph).
//
// ID:      ALG-GRAPH-005
// Time:    O(1)
// Space:   O(1)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1 (|E|)
// Tags:    graph, size, edge-count
func (g *Graph) EdgeCount() int {
	if g == nil {
		return 0
	}
	return g.edges
}

// Directed reports whether g is directed (false for a nil *Graph).
//
// ID:      ALG-GRAPH-006
// Time:    O(1)
// Space:   O(1)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Appendix B.4 (directed and undirected graphs)
// Tags:    graph, directed
func (g *Graph) Directed() bool {
	return g != nil && g.directed
}

// Neighbors returns the sequence of (neighbour, weight) pairs of v's edges, in the order they were
// added. A vertex outside 0..VertexCount()-1 (or a nil *Graph) gives an empty sequence. The list is
// read when the sequence runs.
//
// ID:      ALG-GRAPH-007
// Time:    O(deg(v)) for a full iteration   deg(v) = the number of edges in v's list
// Space:   O(1) extra (no allocation beyond the closure)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1 (Adj[u])
// Source:  The Go Programming Language Specification, "For statements with range clause" (range over function iterators); package iter (iter.Seq2)
// Tags:    graph, adjacency-list, neighbors, iterator
func (g *Graph) Neighbors(v int) iter.Seq2[int, int64] {
	return func(yield func(int, int64) bool) {
		if g == nil || v < 0 || v >= len(g.adj) {
			return
		}
		for _, a := range g.adj[v] {
			if !yield(a.to, a.weight) {
				return
			}
		}
	}
}

// Edges returns every edge of g once: for a directed graph each arc as added, grouped by From; for an
// undirected graph each edge with From <= To (endpoints may be swapped from how it was added), grouped
// by From. A nil *Graph gives nil.
//
// ID:      ALG-GRAPH-008
// Time:    O(V + E)   V = VertexCount(), E = EdgeCount()
// Space:   O(E) for the result (one allocation)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1
// Tags:    graph, edge-list
func (g *Graph) Edges() []Edge {
	if g == nil {
		return nil
	}
	out := make([]Edge, 0, g.edges)
	for from, list := range g.adj {
		for _, a := range list {
			if g.directed || from <= a.to {
				out = append(out, Edge{from, a.to, a.weight})
			}
		}
	}
	return out
}

// checkVertex returns ErrVertexOutOfRange, naming the parameter, when v is not a vertex of g.
func (g *Graph) checkVertex(name string, v int) error {
	if v < 0 || v >= len(g.adj) {
		return fmt.Errorf("%w: %s = %d, n = %d", ErrVertexOutOfRange, name, v, len(g.adj))
	}
	return nil
}

// addWeight returns a + b and true, or false when the sum overflows int64 or equals Unreachable.
func addWeight(a, b int64) (int64, bool) {
	s := a + b
	if (b > 0 && s < a) || (b < 0 && s > a) || s == Unreachable {
		return 0, false
	}
	return s, true
}

// checkedSum returns a + b and true, or false when the sum overflows int64. Unlike addWeight it allows
// math.MaxInt64: a tree weight is not a distance.
func checkedSum(a, b int64) (int64, bool) {
	s := a + b
	if (b > 0 && s < a) || (b < 0 && s > a) {
		return 0, false
	}
	return s, true
}
