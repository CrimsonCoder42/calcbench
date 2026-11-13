package graph

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// Cormen et al., Introduction to Algorithms, 3rd ed., §21.1, Figure 21.1: vertices a..j, edges (b,d) (e,g)
// (a,c) (h,i) (a,b) (e,f) (b,c); components {a,b,c,d} {e,f,g} {h,i} {j}.
func TestConnectedComponents_clrsFigure21_1(t *testing.T) {
	const a, b, c, d, e, f, g, h, i, j = 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
	gr, _ := NewUndirected(10, []Edge{{b, d, 0}, {e, g, 0}, {a, c, 0}, {h, i, 0}, {a, b, 0}, {e, f, 0}, {b, c, 0}})
	comp, count, err := ConnectedComponents(gr)
	if want := []int{0, 0, 0, 0, 1, 1, 1, 2, 2, 3}; err != nil || count != 4 || !slices.Equal(comp, want) {
		t.Errorf("= (%v, %d, %v), want (%v, 4, nil)", comp, count, err, want)
	}
}

func TestConnectedComponents_cases(t *testing.T) {
	empty, _ := NewUndirected(0, nil)
	if comp, count, err := ConnectedComponents(empty); err != nil || count != 0 || len(comp) != 0 {
		t.Errorf("empty: (%v, %d, %v)", comp, count, err)
	}
	// Directed: weak components; labels follow each component's smallest vertex.
	dg, _ := NewDirected(5, []Edge{{4, 0, 0}, {3, 1, 0}, {1, 3, 0}, {2, 2, 0}})
	comp, count, err := ConnectedComponents(dg)
	if want := []int{0, 1, 2, 1, 0}; err != nil || count != 3 || !slices.Equal(comp, want) {
		t.Errorf("directed: (%v, %d, %v), want (%v, 3, nil)", comp, count, err, want)
	}
	if _, _, err := ConnectedComponents(nil); !errors.Is(err, ErrNilGraph) {
		t.Errorf("nil: %v", err)
	}
}

// Oracle: two vertices share a label exactly when BFS on the undirected version joins them.
func TestConnectedComponents_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		n := 1 + rng.IntN(20)
		edges := randomEdges(rng, n, rng.IntN(2*n), 0, 0)
		und, _ := NewUndirected(n, edges)
		g := und
		if run%2 == 1 {
			g, _ = NewDirected(n, edges)
		}
		comp, count, err := ConnectedComponents(g)
		if err != nil {
			t.Fatal(err)
		}
		next := 0
		for s := range n {
			reach := levels(und, s)
			for v := range n {
				if (reach[v] >= 0) != (comp[v] == comp[s]) {
					t.Fatalf("run %d: %d and %d labels %d %d, reachable %v", run, s, v, comp[s], comp[v], reach[v] >= 0)
				}
			}
			if comp[s] == next {
				next++
			} else if comp[s] > next {
				t.Fatalf("run %d: labels not in order of smallest vertex: %v", run, comp)
			}
		}
		if next != count {
			t.Fatalf("run %d: count %d, labels used %d", run, count, next)
		}
	}
}

func BenchmarkConnectedComponents(b *testing.B) {
	for _, size := range benchSizes {
		g := benchGraph(size.n, false)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = ConnectedComponents(g)
			}
		})
	}
}
