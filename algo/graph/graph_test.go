package graph

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// randomEdges returns m edges on 0..n-1 with weights in [minW, maxW], from rng.
func randomEdges(rng *rand.Rand, n, m int, minW, maxW int64) []Edge {
	edges := make([]Edge, m)
	for i := range edges {
		edges[i] = Edge{From: rng.IntN(n), To: rng.IntN(n), Weight: minW + rng.Int64N(maxW-minW+1)}
	}
	return edges
}

// benchGraph returns a graph with V = n vertices and E = 4n random edges, weights 0..1000.
func benchGraph(n int, directed bool) *Graph {
	rng := rand.New(rand.NewPCG(1, 2))
	edges := randomEdges(rng, n, 4*n, 0, 1000)
	if directed {
		g, _ := NewDirected(n, edges)
		return g
	}
	g, _ := NewUndirected(n, edges)
	return g
}

// neighbors collects the (to, weight) pairs of v.
func neighbors(g *Graph, v int) []arc {
	var out []arc
	for to, w := range g.Neighbors(v) {
		out = append(out, arc{to, w})
	}
	return out
}

// Cormen et al., Introduction to Algorithms, 3rd ed., §22.1, Figure 22.1(a,b): the undirected graph on
// 1..5 (here 0..4) with edges 1–2, 1–5, 2–5, 2–3, 2–4, 3–4, 4–5; list lengths 2 4 2 3 3.
func TestNewUndirected_clrsFigure22_1(t *testing.T) {
	g, err := NewUndirected(5, []Edge{{0, 1, 0}, {0, 4, 0}, {1, 4, 0}, {1, 2, 0}, {1, 3, 0}, {2, 3, 0}, {3, 4, 0}})
	if err != nil {
		t.Fatal(err)
	}
	if g.VertexCount() != 5 || g.EdgeCount() != 7 || g.Directed() {
		t.Errorf("V %d E %d directed %v, want 5 7 false", g.VertexCount(), g.EdgeCount(), g.Directed())
	}
	for v, want := range []int{2, 4, 2, 3, 3} {
		if got := len(neighbors(g, v)); got != want {
			t.Errorf("degree(%d) = %d, want %d", v, got, want)
		}
	}
	if got := neighbors(g, 1); !slices.Equal(got, []arc{{0, 0}, {4, 0}, {2, 0}, {3, 0}}) {
		t.Errorf("Neighbors(1) = %v, want insertion order 0 4 2 3", got)
	}
}

func TestNewDirected_cases(t *testing.T) {
	t.Run("n 0", func(t *testing.T) {
		g, err := NewDirected(0, nil)
		if err != nil || g.VertexCount() != 0 || g.EdgeCount() != 0 || !g.Directed() || len(g.Edges()) != 0 {
			t.Errorf("NewDirected(0, nil) = (%v, %v)", g, err)
		}
	})
	t.Run("arcs one way, weights kept", func(t *testing.T) {
		g, _ := NewDirected(3, []Edge{{0, 1, -5}, {1, 2, math.MaxInt64}, {0, 2, math.MinInt64}})
		if got := neighbors(g, 0); !slices.Equal(got, []arc{{1, -5}, {2, math.MinInt64}}) {
			t.Errorf("Neighbors(0) = %v", got)
		}
		if got := neighbors(g, 2); len(got) != 0 {
			t.Errorf("Neighbors(2) = %v, want none", got)
		}
	})
	t.Run("self-loop and parallel edges", func(t *testing.T) {
		g, _ := NewUndirected(2, []Edge{{0, 0, 1}, {0, 1, 2}, {1, 0, 3}})
		if got := neighbors(g, 0); !slices.Equal(got, []arc{{0, 1}, {1, 2}, {1, 3}}) {
			t.Errorf("Neighbors(0) = %v, want self-loop once and both parallel edges", got)
		}
		if g.EdgeCount() != 3 {
			t.Errorf("EdgeCount = %d, want 3", g.EdgeCount())
		}
		want := []Edge{{0, 0, 1}, {0, 1, 2}, {0, 1, 3}}
		if got := g.Edges(); !slices.Equal(got, want) {
			t.Errorf("Edges = %v, want %v", got, want)
		}
	})
	t.Run("input not modified", func(t *testing.T) {
		edges := []Edge{{1, 0, 4}, {0, 1, 2}}
		before := slices.Clone(edges)
		_, _ = NewUndirected(2, edges)
		if !slices.Equal(edges, before) {
			t.Error("NewUndirected modified edges")
		}
	})
}

func TestGraphEdges_roundTrip(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		n := 1 + rng.IntN(12)
		edges := randomEdges(rng, n, rng.IntN(30), -9, 9)
		for _, directed := range []bool{true, false} {
			build := NewUndirected
			if directed {
				build = NewDirected
			}
			g, _ := build(n, edges)
			h, err := build(n, g.Edges())
			if err != nil || len(g.Edges()) != len(edges) {
				t.Fatalf("run %d: Edges has %d edges, want %d (%v)", run, len(g.Edges()), len(edges), err)
			}
			for v := range n {
				a, b := neighbors(g, v), neighbors(h, v)
				slices.SortFunc(a, cmpArc)
				slices.SortFunc(b, cmpArc)
				if !slices.Equal(a, b) {
					t.Fatalf("run %d directed %v: vertex %d lists differ: %v vs %v", run, directed, v, a, b)
				}
			}
		}
	}
}

func cmpArc(a, b arc) int {
	if a.to != b.to {
		return a.to - b.to
	}
	switch {
	case a.weight < b.weight:
		return -1
	case a.weight > b.weight:
		return 1
	}
	return 0
}

func TestGraphAddEdge_cases(t *testing.T) {
	g, _ := NewUndirected(3, []Edge{{0, 1, 1}, {1, 2, 1}})
	if err := g.AddEdge(0, 2, 7); err != nil {
		t.Fatal(err)
	}
	// Row 0 was full (capacity = degree), so the append must not overwrite row 1 in the shared array.
	if got := neighbors(g, 1); !slices.Equal(got, []arc{{0, 1}, {2, 1}}) {
		t.Errorf("Neighbors(1) after AddEdge(0, 2) = %v, want [{0 1} {2 1}]", got)
	}
	if got := neighbors(g, 2); !slices.Equal(got, []arc{{1, 1}, {0, 7}}) {
		t.Errorf("Neighbors(2) = %v", got)
	}
	if g.EdgeCount() != 3 {
		t.Errorf("EdgeCount = %d, want 3", g.EdgeCount())
	}
	if err := g.AddEdge(2, 2, 0); err != nil || len(neighbors(g, 2)) != 3 {
		t.Errorf("undirected self-loop: err %v, degree %d, want nil 3", err, len(neighbors(g, 2)))
	}
	d, _ := NewDirected(2, nil)
	_ = d.AddEdge(1, 0, 4)
	if len(neighbors(d, 0)) != 0 || !slices.Equal(neighbors(d, 1), []arc{{0, 4}}) {
		t.Error("directed AddEdge stored the wrong arcs")
	}
}

func TestGraphNeighbors_edges(t *testing.T) {
	g, _ := NewDirected(2, []Edge{{0, 1, 1}, {0, 0, 2}})
	for _, v := range []int{-1, 2, math.MaxInt} {
		if got := neighbors(g, v); len(got) != 0 {
			t.Errorf("Neighbors(%d) = %v, want empty", v, got)
		}
	}
	count := 0
	for range g.Neighbors(0) {
		count++
		break
	}
	if count != 1 {
		t.Errorf("break after one: %d", count)
	}
}

func TestGraph_nil(t *testing.T) {
	var g *Graph
	if g.VertexCount() != 0 || g.EdgeCount() != 0 || g.Directed() || g.Edges() != nil || len(neighbors(g, 0)) != 0 {
		t.Error("nil Graph methods: want zero values")
	}
	if err := g.AddEdge(0, 0, 0); !errors.Is(err, ErrNilGraph) {
		t.Errorf("nil AddEdge err = %v, want ErrNilGraph", err)
	}
}

func TestNewDirected_errors(t *testing.T) {
	for _, build := range []func(int, []Edge) (*Graph, error){NewDirected, NewUndirected} {
		if g, err := build(-1, nil); !errors.Is(err, ErrNegativeVertexCount) || g != nil {
			t.Errorf("n -1: (%v, %v), want (nil, ErrNegativeVertexCount)", g, err)
		}
		for _, e := range []Edge{{-1, 0, 0}, {0, -1, 0}, {3, 0, 0}, {0, 3, 0}, {math.MinInt, math.MaxInt, 0}} {
			if g, err := build(3, []Edge{{0, 1, 0}, e}); !errors.Is(err, ErrVertexOutOfRange) || g != nil {
				t.Errorf("edge %v: (%v, %v), want (nil, ErrVertexOutOfRange)", e, g, err)
			}
		}
		if _, err := build(0, []Edge{{0, 0, 0}}); !errors.Is(err, ErrVertexOutOfRange) {
			t.Errorf("edge on empty graph: err = %v, want ErrVertexOutOfRange", err)
		}
		// L41: a vertex count whose list headers exceed 2⁴⁷ bytes (2³¹ on 32-bit) is an error, not a panic.
		for _, n := range []int{math.MaxInt, maxAllocBytes/headerBytes + 1} {
			if g, err := build(n, nil); !errors.Is(err, ErrTooLarge) || g != nil {
				t.Errorf("n %d: (%v, %v), want (nil, ErrTooLarge)", n, g, err)
			}
		}
	}
	g, _ := NewUndirected(2, nil)
	for _, e := range [][2]int{{-1, 0}, {0, 2}} {
		if err := g.AddEdge(e[0], e[1], 0); !errors.Is(err, ErrVertexOutOfRange) {
			t.Errorf("AddEdge%v err = %v, want ErrVertexOutOfRange", e, err)
		}
	}
	if g.EdgeCount() != 0 {
		t.Error("failed AddEdge changed EdgeCount")
	}
}

func TestAddWeight_overflow(t *testing.T) {
	tests := []struct {
		a, b int64
		ok   bool
	}{
		{1, 2, true}, {math.MaxInt64 - 2, 1, true}, {math.MaxInt64 - 1, 1, false}, {math.MaxInt64, 1, false},
		{math.MinInt64, -1, false}, {math.MinInt64 + 1, -1, true}, {-5, 3, true},
	}
	for _, tt := range tests {
		if _, ok := addWeight(tt.a, tt.b); ok != tt.ok {
			t.Errorf("addWeight(%d, %d) ok = %v, want %v", tt.a, tt.b, ok, tt.ok)
		}
	}
	if _, ok := checkedSum(math.MaxInt64-1, 1); !ok {
		t.Error("checkedSum(MaxInt64-1, 1) = overflow, want ok (MaxInt64 is a valid tree weight)")
	}
}

func BenchmarkNewDirected(b *testing.B) {
	for _, size := range benchSizes {
		edges := randomEdges(rand.New(rand.NewPCG(1, 2)), size.n, 4*size.n, 0, 1000)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = NewDirected(size.n, edges)
			}
		})
	}
}

func BenchmarkNewUndirected(b *testing.B) {
	for _, size := range benchSizes {
		edges := randomEdges(rand.New(rand.NewPCG(1, 2)), size.n, 4*size.n, 0, 1000)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = NewUndirected(size.n, edges)
			}
		})
	}
}

func BenchmarkGraphAddEdge(b *testing.B) {
	for _, size := range benchSizes {
		edges := randomEdges(rand.New(rand.NewPCG(1, 2)), size.n, size.n, 0, 1000)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				g, _ := NewDirected(size.n, nil)
				for _, e := range edges {
					_ = g.AddEdge(e.From, e.To, e.Weight)
				}
			}
		})
	}
}

func BenchmarkGraphNeighbors(b *testing.B) {
	for _, size := range benchSizes {
		g, _ := NewDirected(2, nil)
		for range size.n {
			_ = g.AddEdge(0, 1, 1)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for range g.Neighbors(0) {
				}
			}
		})
	}
}

func BenchmarkGraphEdges(b *testing.B) {
	for _, size := range benchSizes {
		g := benchGraph(size.n, false)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = g.Edges()
			}
		})
	}
}

func BenchmarkGraphVertexCount(b *testing.B) {
	g := benchGraph(1000, true)
	b.ReportAllocs()
	for b.Loop() {
		_ = g.VertexCount()
	}
}

func BenchmarkGraphEdgeCount(b *testing.B) {
	g := benchGraph(1000, true)
	b.ReportAllocs()
	for b.Loop() {
		_ = g.EdgeCount()
	}
}

func BenchmarkGraphDirected(b *testing.B) {
	g := benchGraph(1000, true)
	b.ReportAllocs()
	for b.Loop() {
		_ = g.Directed()
	}
}
