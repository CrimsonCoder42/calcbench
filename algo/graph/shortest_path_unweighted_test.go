package graph

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// Cormen et al., Introduction to Algorithms, 3rd ed., §22.2, Figure 22.3: δ(s, y) = 3, and the only
// shortest path is s w x y.
func TestShortestPathUnweighted_clrsFigure22_3(t *testing.T) {
	g := figure22_3(t)
	got, err := ShortestPathUnweighted(g, fs, fy)
	if want := []int{fs, fw, fx, fy}; err != nil || !slices.Equal(got, want) {
		t.Errorf("s→y = (%v, %v), want (%v, nil)", got, err, want)
	}
	got, _ = ShortestPathUnweighted(g, fv, fu)
	if len(got) != 6 || got[0] != fv || got[5] != fu { // v r s w t|x u: δ(v, u) = δ(v, s) + δ(s, u) = 2 + 3 edges
		t.Errorf("v→u = %v, want 6 vertices", got)
	}
}

func TestShortestPathUnweighted_cases(t *testing.T) {
	g, _ := NewDirected(4, []Edge{{0, 1, 99}, {1, 2, 99}, {0, 2, 1000}, {3, 0, 0}})
	tests := []struct {
		name           string
		source, target int
		want           []int
	}{
		{"same vertex", 2, 2, []int{2}},
		{"direct edge beats two light edges (weights ignored)", 0, 2, []int{0, 2}},
		{"one edge", 0, 1, []int{0, 1}},
		{"against the arrow: unreachable", 2, 0, nil},
		{"two hops", 3, 2, []int{3, 0, 2}},
	}
	for _, tt := range tests {
		got, err := ShortestPathUnweighted(g, tt.source, tt.target)
		if err != nil || !slices.Equal(got, tt.want) || (tt.want == nil) != (got == nil) {
			t.Errorf("%s: (%v, %v), want (%v, nil)", tt.name, got, err, tt.want)
		}
	}
}

func TestShortestPathUnweighted_errors(t *testing.T) {
	if _, err := ShortestPathUnweighted(nil, 0, 0); !errors.Is(err, ErrNilGraph) {
		t.Errorf("nil: %v", err)
	}
	g, _ := NewUndirected(2, nil)
	if _, err := ShortestPathUnweighted(g, 2, 0); !errors.Is(err, ErrVertexOutOfRange) {
		t.Errorf("bad source: %v", err)
	}
	if _, err := ShortestPathUnweighted(g, 0, -1); !errors.Is(err, ErrVertexOutOfRange) {
		t.Errorf("bad target: %v", err)
	}
}

func TestShortestPathUnweighted_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		n := 1 + rng.IntN(15)
		g, _ := NewDirected(n, randomEdges(rng, n, rng.IntN(3*n), 0, 0))
		s, d := rng.IntN(n), rng.IntN(n)
		path, _ := ShortestPathUnweighted(g, s, d)
		want := levels(g, s)[d]
		if want < 0 {
			if path != nil {
				t.Fatalf("run %d: path %v to an unreachable vertex", run, path)
			}
			continue
		}
		if len(path) != want+1 || path[0] != s || path[len(path)-1] != d {
			t.Fatalf("run %d: path %v, want %d edges from %d to %d", run, path, want, s, d)
		}
		for i := 1; i < len(path); i++ {
			found := false
			for w := range g.Neighbors(path[i-1]) {
				found = found || w == path[i]
			}
			if !found {
				t.Fatalf("run %d: %d→%d is not an edge", run, path[i-1], path[i])
			}
		}
	}
}

func BenchmarkShortestPathUnweighted(b *testing.B) {
	for _, size := range benchSizes {
		g := benchGraph(size.n, true)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = ShortestPathUnweighted(g, 0, size.n-1)
			}
		})
	}
}
