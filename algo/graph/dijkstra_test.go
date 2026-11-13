package graph

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Vertices of CLRS Figures 24.4 and 24.6 (s t x y z).
const vs, vt, vx, vy, vz = 0, 1, 2, 3, 4

// Cormen et al., Introduction to Algorithms, 3rd ed., §24.3, Figure 24.6: edges s→t 10, s→y 5, t→x 1,
// t→y 2, y→t 3, y→x 9, y→z 2, x→z 4, z→x 6, z→s 7; from s: δ = s0 t8 x9 y5 z7, tree t←y x←t y←s z←y.
func TestDijkstra_clrsFigure24_6(t *testing.T) {
	g, _ := NewDirected(5, []Edge{
		{vs, vt, 10}, {vs, vy, 5}, {vt, vx, 1}, {vt, vy, 2}, {vy, vt, 3},
		{vy, vx, 9}, {vy, vz, 2}, {vx, vz, 4}, {vz, vx, 6}, {vz, vs, 7},
	})
	dist, parent, err := Dijkstra(g, vs)
	if err != nil {
		t.Fatal(err)
	}
	if want := []int64{0, 8, 9, 5, 7}; !slices.Equal(dist, want) {
		t.Errorf("dist = %v, want %v", dist, want)
	}
	if want := []int{-1, vy, vt, vs, vy}; !slices.Equal(parent, want) {
		t.Errorf("parent = %v, want %v", parent, want)
	}
}

func TestDijkstra_cases(t *testing.T) {
	g, _ := NewUndirected(5, []Edge{{0, 1, 0}, {1, 2, 0}, {0, 2, 5}, {0, 2, 1}, {3, 3, 4}})
	dist, parent, err := Dijkstra(g, 0)
	if err != nil {
		t.Fatal(err)
	}
	if want := []int64{0, 0, 0, Unreachable, Unreachable}; !slices.Equal(dist, want) {
		t.Errorf("zero weights, parallel edges, unreachable: dist = %v, want %v", dist, want)
	}
	if parent[3] != -1 || parent[4] != -1 || parent[0] != -1 {
		t.Errorf("parent = %v, want -1 for source and unreachable", parent)
	}
	one, _ := NewDirected(1, nil)
	if d, p, err := Dijkstra(one, 0); err != nil || !slices.Equal(d, []int64{0}) || !slices.Equal(p, []int{-1}) {
		t.Errorf("one vertex: (%v, %v, %v)", d, p, err)
	}
	big, _ := NewDirected(2, []Edge{{0, 1, math.MaxInt64 - 1}})
	if d, _, err := Dijkstra(big, 0); err != nil || d[1] != math.MaxInt64-1 {
		t.Errorf("largest representable distance: (%v, %v)", d, err)
	}
}

// Review M1: a sum that overflows upward is "no improvement", not an error, when a shorter path exists.
func TestDijkstra_overflowIsNoImprovement(t *testing.T) {
	g, _ := NewDirected(2, []Edge{{0, 1, 1}, {1, 0, math.MaxInt64}})
	if d, _, err := Dijkstra(g, 0); err != nil || !slices.Equal(d, []int64{0, 1}) {
		t.Errorf("back edge MaxInt64: (%v, %v), want ([0 1], nil)", d, err)
	}
	// The overflowing edge 0→2 is relaxed first; the representable path 0→1→2 is found later.
	h, _ := NewDirected(3, []Edge{{0, 2, math.MaxInt64}, {0, 1, 1}, {1, 2, 1}})
	if d, p, err := Dijkstra(h, 0); err != nil || !slices.Equal(d, []int64{0, 1, 2}) || p[2] != 1 {
		t.Errorf("later shorter path: (%v, %v, %v), want ([0 1 2], parent[2] 1, nil)", d, p, err)
	}
}

func TestDijkstra_errors(t *testing.T) {
	if _, _, err := Dijkstra(nil, 0); !errors.Is(err, ErrNilGraph) {
		t.Errorf("nil: %v", err)
	}
	g, _ := NewDirected(3, []Edge{{0, 1, 1}, {2, 1, -1}})
	if _, _, err := Dijkstra(g, 3); !errors.Is(err, ErrVertexOutOfRange) {
		t.Errorf("bad source: %v", err)
	}
	if d, _, err := Dijkstra(g, 0); !errors.Is(err, ErrNegativeWeight) || d != nil {
		t.Errorf("negative weight (unreachable edge too): (%v, %v)", d, err)
	}
	for _, edges := range [][]Edge{
		{{0, 1, math.MaxInt64}}, // a distance equal to Unreachable
		{{0, 1, math.MaxInt64/2 + 1}, {1, 2, math.MaxInt64/2 + 1}}, // a sum past MaxInt64
	} {
		g, _ := NewDirected(3, edges)
		if _, _, err := Dijkstra(g, 0); !errors.Is(err, ErrOverflow) {
			t.Errorf("%v: err = %v, want ErrOverflow", edges, err)
		}
	}
}

// Oracle: on non-negative weights Dijkstra and BellmanFord agree on every distance.
func TestDijkstra_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		n := 1 + rng.IntN(20)
		edges := randomEdges(rng, n, rng.IntN(4*n), 0, 50)
		g, _ := NewDirected(n, edges)
		if run%2 == 1 {
			g, _ = NewUndirected(n, edges)
		}
		s := rng.IntN(n)
		got, parent, err := Dijkstra(g, s)
		want, _, err2 := BellmanFord(g, s)
		if err != nil || err2 != nil || !slices.Equal(got, want) {
			t.Fatalf("run %d: Dijkstra %v (%v), BellmanFord %v (%v)", run, got, err, want, err2)
		}
		for v, p := range parent { // each tree edge is tight
			if p < 0 {
				continue
			}
			tight := false
			for w, wt := range g.Neighbors(p) {
				tight = tight || (w == v && got[p]+wt == got[v])
			}
			if !tight {
				t.Fatalf("run %d: parent %d of %d is not on a shortest path", run, p, v)
			}
		}
	}
}

func BenchmarkDijkstra(b *testing.B) {
	for _, size := range benchSizes {
		g := benchGraph(size.n, true)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = Dijkstra(g, 0)
			}
		})
	}
}
