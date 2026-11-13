package graph

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"

	"github.com/CrimsonCoder42/calcbench/algo/unionfind"
)

// figure23_1 is the undirected graph of Cormen et al., Introduction to Algorithms, 3rd ed., §23.1,
// Figure 23.1 (a..i = 0..8): a–b 4, a–h 8, b–c 8, b–h 11, c–d 7, c–f 4, c–i 2, d–e 9, d–f 14, e–f 10,
// f–g 2, g–h 1, g–i 6, h–i 7. Its minimum spanning trees weigh 37 (Figures 23.4 and 23.5).
func figure23_1(t testing.TB) *Graph {
	const a, b, c, d, e, f, g, h, i = 0, 1, 2, 3, 4, 5, 6, 7, 8
	gr, err := NewUndirected(9, []Edge{
		{a, b, 4}, {a, h, 8}, {b, c, 8}, {b, h, 11}, {c, d, 7}, {c, f, 4}, {c, i, 2},
		{d, e, 9}, {d, f, 14}, {e, f, 10}, {f, g, 2}, {g, h, 1}, {g, i, 6}, {h, i, 7},
	})
	if err != nil {
		t.Fatal(err)
	}
	return gr
}

// checkSpanningForest fails unless tree is an acyclic set of edges of g with one edge fewer than the
// vertices in each component, weighing total.
func checkSpanningForest(t *testing.T, g *Graph, tree []Edge, total int64) {
	t.Helper()
	_, components, _ := ConnectedComponents(g)
	if len(tree) != g.VertexCount()-components {
		t.Fatalf("%d tree edges, want V - components = %d", len(tree), g.VertexCount()-components)
	}
	sets, _ := unionfind.NewDisjointSet(g.VertexCount())
	var sum int64
	for _, e := range tree {
		found := false
		for w, wt := range g.Neighbors(e.From) {
			found = found || (w == e.To && wt == e.Weight)
		}
		if merged, _ := sets.Union(e.From, e.To); !merged || !found {
			t.Fatalf("tree edge %v is a cycle edge or not in g", e)
		}
		sum += e.Weight
	}
	if sum != total {
		t.Fatalf("edges sum to %d, reported %d", sum, total)
	}
}

// bruteForceMSF is the oracle: the least weight over all acyclic edge subsets of size V − components.
func bruteForceMSF(g *Graph) int64 {
	edges := g.Edges()
	_, components, _ := ConnectedComponents(g)
	need := g.VertexCount() - components
	best := int64(math.MaxInt64)
	for mask := 0; mask < 1<<len(edges); mask++ {
		sets, _ := unionfind.NewDisjointSet(g.VertexCount())
		var sum int64
		size, ok := 0, true
		for k, e := range edges {
			if mask>>k&1 == 1 {
				merged, _ := sets.Union(e.From, e.To)
				ok = ok && merged
				sum += e.Weight
				size++
			}
		}
		if ok && size == need {
			best = min(best, sum)
		}
	}
	return best
}

// CLRS Figure 23.4: MST-KRUSKAL on Figure 23.1 takes 8 edges of total weight 37.
func TestKruskal_clrsFigure23_4(t *testing.T) {
	g := figure23_1(t)
	tree, total, err := Kruskal(g)
	if err != nil || total != 37 || len(tree) != 8 {
		t.Fatalf("Kruskal = (%v, %d, %v), want 8 edges weighing 37", tree, total, err)
	}
	checkSpanningForest(t, g, tree, total)
	if tree[0] != (Edge{From: 6, To: 7, Weight: 1}) { // the lightest edge g–h goes first
		t.Errorf("first edge %v, want g–h 1", tree[0])
	}
}

func TestKruskal_cases(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		edges []Edge
		total int64
	}{
		{"empty", 0, nil, 0},
		{"one vertex, self-loop", 1, []Edge{{0, 0, -5}}, 0},
		{"parallel edges take the lighter", 2, []Edge{{0, 1, 9}, {1, 0, 3}}, 3},
		{"negative weights", 3, []Edge{{0, 1, -4}, {1, 2, -1}, {0, 2, -2}}, -6},
		{"forest of two trees", 5, []Edge{{0, 1, 1}, {1, 2, 2}, {0, 2, 3}, {3, 4, 4}}, 7},
		{"total exactly MaxInt64", 3, []Edge{{0, 1, math.MaxInt64 - 1}, {1, 2, 1}}, math.MaxInt64},
	}
	for _, tt := range tests {
		g, _ := NewUndirected(tt.n, tt.edges)
		tree, total, err := Kruskal(g)
		if err != nil || total != tt.total {
			t.Errorf("%s: (%v, %d, %v), want total %d", tt.name, tree, total, err, tt.total)
			continue
		}
		checkSpanningForest(t, g, tree, total)
	}
}

func TestKruskal_errors(t *testing.T) {
	if _, _, err := Kruskal(nil); !errors.Is(err, ErrNilGraph) {
		t.Errorf("nil: %v", err)
	}
	d, _ := NewDirected(2, []Edge{{0, 1, 1}})
	if _, _, err := Kruskal(d); !errors.Is(err, ErrDirected) {
		t.Errorf("directed: %v", err)
	}
	over, _ := NewUndirected(3, []Edge{{0, 1, math.MaxInt64}, {1, 2, math.MaxInt64}})
	if _, _, err := Kruskal(over); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: %v", err)
	}
	under, _ := NewUndirected(3, []Edge{{0, 1, math.MinInt64}, {1, 2, -1}})
	if _, _, err := Kruskal(under); !errors.Is(err, ErrOverflow) {
		t.Errorf("underflow: %v", err)
	}
}

// Oracle: Kruskal, Prim and brute force agree on the forest weight (n ≤ 7, at most 12 edges).
func TestKruskal_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		n := 1 + rng.IntN(7)
		g, _ := NewUndirected(n, randomEdges(rng, n, rng.IntN(13), -10, 10))
		want := bruteForceMSF(g)
		kt, kw, err1 := Kruskal(g)
		pt, pw, err2 := Prim(g)
		if err1 != nil || err2 != nil || kw != want || pw != want {
			t.Fatalf("run %d: Kruskal %d (%v), Prim %d (%v), brute force %d", run, kw, err1, pw, err2, want)
		}
		checkSpanningForest(t, g, kt, kw)
		checkSpanningForest(t, g, pt, pw)
	}
}

func BenchmarkKruskal(b *testing.B) {
	for _, size := range benchSizes {
		g := benchGraph(size.n, false)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = Kruskal(g)
			}
		})
	}
}
