package graph

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Cormen et al., Introduction to Algorithms, 3rd ed., §24.1, Figure 24.4: edges s→t 6, s→y 7, t→x 5,
// t→y 8, t→z −4, x→t −2, y→x −3, y→z 9, z→s 2, z→x 7; from s: δ = s0 t2 x4 y7 z−2, tree t←x x←y y←s z←t.
func TestBellmanFord_clrsFigure24_4(t *testing.T) {
	g, _ := NewDirected(5, []Edge{
		{vs, vt, 6}, {vs, vy, 7}, {vt, vx, 5}, {vt, vy, 8}, {vt, vz, -4},
		{vx, vt, -2}, {vy, vx, -3}, {vy, vz, 9}, {vz, vs, 2}, {vz, vx, 7},
	})
	dist, parent, err := BellmanFord(g, vs)
	if err != nil {
		t.Fatal(err)
	}
	if want := []int64{0, 2, 4, 7, -2}; !slices.Equal(dist, want) {
		t.Errorf("dist = %v, want %v", dist, want)
	}
	if want := []int{-1, vx, vy, vs, vt}; !slices.Equal(parent, want) {
		t.Errorf("parent = %v, want %v", parent, want)
	}
}

func TestBellmanFord_cases(t *testing.T) {
	// A negative cycle 2→3→2 that the source cannot reach does not matter.
	g, _ := NewDirected(4, []Edge{{0, 1, -5}, {2, 3, -1}, {3, 2, -1}})
	dist, _, err := BellmanFord(g, 0)
	if want := []int64{0, -5, Unreachable, Unreachable}; err != nil || !slices.Equal(dist, want) {
		t.Errorf("unreachable negative cycle: (%v, %v), want (%v, nil)", dist, err, want)
	}
	one, _ := NewDirected(1, nil)
	if d, p, err := BellmanFord(one, 0); err != nil || !slices.Equal(d, []int64{0}) || !slices.Equal(p, []int{-1}) {
		t.Errorf("one vertex: (%v, %v, %v)", d, p, err)
	}
	low, _ := NewDirected(2, []Edge{{0, 1, math.MinInt64}})
	if d, _, err := BellmanFord(low, 0); err != nil || d[1] != math.MinInt64 {
		t.Errorf("most negative distance: (%v, %v)", d, err)
	}
	zeroCycle, _ := NewDirected(2, []Edge{{0, 1, 3}, {1, 0, -3}})
	if d, _, err := BellmanFord(zeroCycle, 0); err != nil || !slices.Equal(d, []int64{0, 3}) {
		t.Errorf("zero-weight cycle is allowed: (%v, %v)", d, err)
	}
}

// Review M1: a sum that overflows upward is "no improvement", not an error, when a shorter path exists.
func TestBellmanFord_overflowIsNoImprovement(t *testing.T) {
	g, _ := NewDirected(2, []Edge{{0, 1, 1}, {1, 0, math.MaxInt64}})
	if d, _, err := BellmanFord(g, 0); err != nil || !slices.Equal(d, []int64{0, 1}) {
		t.Errorf("back edge MaxInt64: (%v, %v), want ([0 1], nil)", d, err)
	}
	h, _ := NewDirected(3, []Edge{{0, 2, math.MaxInt64}, {0, 1, 1}, {1, 2, -5}})
	if d, _, err := BellmanFord(h, 0); err != nil || !slices.Equal(d, []int64{0, 1, -4}) {
		t.Errorf("later shorter path: (%v, %v), want ([0 1 -4], nil)", d, err)
	}
}

func TestBellmanFord_errors(t *testing.T) {
	if _, _, err := BellmanFord(nil, 0); !errors.Is(err, ErrNilGraph) {
		t.Errorf("nil: %v", err)
	}
	g, _ := NewDirected(2, []Edge{{0, 1, 1}, {1, 0, -2}})
	if _, _, err := BellmanFord(g, -1); !errors.Is(err, ErrVertexOutOfRange) {
		t.Errorf("bad source: %v", err)
	}
	if d, _, err := BellmanFord(g, 0); !errors.Is(err, ErrNegativeCycle) || d != nil {
		t.Errorf("negative cycle: (%v, %v)", d, err)
	}
	loop, _ := NewDirected(1, []Edge{{0, 0, -1}})
	if _, _, err := BellmanFord(loop, 0); !errors.Is(err, ErrNegativeCycle) {
		t.Errorf("negative self-loop: %v", err)
	}
	und, _ := NewUndirected(2, []Edge{{0, 1, -1}})
	if _, _, err := BellmanFord(und, 0); !errors.Is(err, ErrNegativeCycle) {
		t.Errorf("undirected negative edge: %v", err)
	}
	over, _ := NewDirected(3, []Edge{{0, 1, math.MinInt64}, {1, 2, -1}})
	if _, _, err := BellmanFord(over, 0); !errors.Is(err, ErrOverflow) {
		t.Errorf("sum below MinInt64: %v", err)
	}
	far, _ := NewDirected(3, []Edge{{0, 1, math.MaxInt64 - 1}, {1, 2, 5}})
	if _, _, err := BellmanFord(far, 0); !errors.Is(err, ErrOverflow) {
		t.Errorf("only path weighs more than MaxInt64: %v", err)
	}
	exact, _ := NewDirected(2, []Edge{{0, 1, math.MaxInt64}})
	if _, _, err := BellmanFord(exact, 0); !errors.Is(err, ErrOverflow) {
		t.Errorf("only path weighs exactly MaxInt64 (= Unreachable): %v", err)
	}
}

// Oracle: on random DAGs (edges from a lower to a higher vertex) with negative weights, BellmanFord
// matches a dynamic program over the vertices in increasing order.
func TestBellmanFord_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		n := 1 + rng.IntN(15)
		var edges []Edge
		for range rng.IntN(3 * n) {
			a, b := rng.IntN(n), rng.IntN(n)
			if a == b {
				continue
			}
			edges = append(edges, Edge{From: min(a, b), To: max(a, b), Weight: rng.Int64N(41) - 20})
		}
		g, _ := NewDirected(n, edges)
		s := rng.IntN(n)
		want := make([]int64, n)
		for i := range want {
			want[i] = Unreachable
		}
		want[s] = 0
		for u := s; u < n; u++ {
			if want[u] == Unreachable {
				continue
			}
			for v, w := range g.Neighbors(u) {
				want[v] = min(want[v], want[u]+w)
			}
		}
		got, _, err := BellmanFord(g, s)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("run %d: BellmanFord = (%v, %v), DP = %v", run, got, err, want)
		}
	}
}

func BenchmarkBellmanFord(b *testing.B) {
	// O(V·(V + E)) worst case; random graphs settle in a few passes, so sizes stay at 1e3/1e4/1e5.
	for _, size := range benchSizes {
		g := benchGraph(size.n, true)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = BellmanFord(g, 0)
			}
		})
	}
}
