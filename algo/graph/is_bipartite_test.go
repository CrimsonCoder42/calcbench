package graph

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand examples (catalog/graph.md, ALG-GRAPH-017), in the setting of CLRS Exercise 22.2-7: rivalries
// that can all cross between two teams.
func TestIsBipartite_cases(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		edges []Edge
		want  []bool // nil: not bipartite
	}{
		{"empty graph", 0, nil, []bool{}},
		{"isolated vertices", 3, nil, []bool{false, false, false}},
		{"one edge", 2, []Edge{{0, 1, 0}}, []bool{false, true}},
		{"4-cycle", 4, []Edge{{0, 1, 0}, {1, 2, 0}, {2, 3, 0}, {3, 0, 0}}, []bool{false, true, false, true}},
		{"triangle", 3, []Edge{{0, 1, 0}, {1, 2, 0}, {2, 0, 0}}, nil},
		{"self-loop", 2, []Edge{{0, 1, 0}, {1, 1, 0}}, nil},
		{"parallel edges", 2, []Edge{{0, 1, 0}, {1, 0, 0}}, []bool{false, true}},
		{"two components, second starts at 2", 4, []Edge{{0, 1, 0}, {3, 2, 0}}, []bool{false, true, false, true}},
		{"odd cycle in the second component", 6, []Edge{{0, 1, 0}, {2, 3, 0}, {3, 4, 0}, {4, 5, 0}, {5, 2, 0}, {2, 4, 0}}, nil},
	}
	for _, tt := range tests {
		g, _ := NewUndirected(tt.n, tt.edges)
		side, ok, err := IsBipartite(g)
		if err != nil || ok != (tt.want != nil) || !slices.Equal(side, tt.want) || (side == nil) != (tt.want == nil) {
			t.Errorf("%s: (%v, %v, %v), want %v", tt.name, side, ok, err, tt.want)
		}
	}
}

func TestIsBipartite_errors(t *testing.T) {
	if _, _, err := IsBipartite(nil); !errors.Is(err, ErrNilGraph) {
		t.Errorf("nil: %v", err)
	}
	d, _ := NewDirected(2, []Edge{{0, 1, 0}})
	if _, _, err := IsBipartite(d); !errors.Is(err, ErrDirected) {
		t.Errorf("directed: %v", err)
	}
}

// Oracle: brute force over all 2^n colourings for n ≤ 10.
func TestIsBipartite_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		n := 1 + rng.IntN(10)
		edges := randomEdges(rng, n, rng.IntN(2*n), 0, 0)
		g, _ := NewUndirected(n, edges)
		want := false
		for mask := 0; mask < 1<<n && !want; mask++ {
			good := true
			for _, e := range edges {
				good = good && (mask>>e.From&1) != (mask>>e.To&1)
			}
			want = good
		}
		side, ok, err := IsBipartite(g)
		if err != nil || ok != want {
			t.Fatalf("run %d: ok = %v (%v), brute force %v", run, ok, err, want)
		}
		for _, e := range edges {
			if ok && side[e.From] == side[e.To] {
				t.Fatalf("run %d: edge %v inside one side", run, e)
			}
		}
	}
}

func BenchmarkIsBipartite(b *testing.B) {
	for _, size := range benchSizes {
		// A bipartite graph (edges between even and odd vertices) so the whole graph is searched.
		rng := rand.New(rand.NewPCG(1, 2))
		edges := make([]Edge, 4*size.n)
		for i := range edges {
			edges[i] = Edge{From: 2 * rng.IntN(size.n/2), To: 2*rng.IntN(size.n/2) + 1}
		}
		g, _ := NewUndirected(size.n, edges)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = IsBipartite(g)
			}
		})
	}
}
