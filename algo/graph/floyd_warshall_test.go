package graph

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// Cormen et al., Introduction to Algorithms, 3rd ed., §25.2, Figure 25.4: vertices 1..5 (here 0..4),
// edges 1→2 3, 1→3 8, 1→5 −4, 2→4 1, 2→5 7, 3→2 4, 4→1 2, 4→3 −5, 5→4 6; the final matrix D⁽⁵⁾.
func TestFloydWarshall_clrsFigure25_4(t *testing.T) {
	g, _ := NewDirected(5, []Edge{
		{0, 1, 3}, {0, 2, 8}, {0, 4, -4}, {1, 3, 1}, {1, 4, 7}, {2, 1, 4}, {3, 0, 2}, {3, 2, -5}, {4, 3, 6},
	})
	d, err := FloydWarshall(g)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]int64{
		{0, 1, -3, 2, -4},
		{3, 0, -4, 1, -1},
		{7, 4, 0, 5, 3},
		{2, -1, -5, 0, -2},
		{8, 5, 1, 6, 0},
	}
	for i := range want {
		if !slices.Equal(d[i], want[i]) {
			t.Errorf("row %d = %v, want %v", i+1, d[i], want[i])
		}
	}
}

func TestFloydWarshall_cases(t *testing.T) {
	empty, _ := NewDirected(0, nil)
	if d, err := FloydWarshall(empty); err != nil || len(d) != 0 {
		t.Errorf("empty: (%v, %v)", d, err)
	}
	g, _ := NewUndirected(3, []Edge{{0, 1, 7}, {0, 1, 2}, {1, 1, 0}})
	d, err := FloydWarshall(g)
	want := [][]int64{{0, 2, Unreachable}, {2, 0, Unreachable}, {Unreachable, Unreachable, 0}}
	if err != nil || !slices.EqualFunc(d, want, slices.Equal) {
		t.Errorf("parallel edges, zero self-loop, isolated vertex: (%v, %v), want %v", d, err, want)
	}
	d[0][0] = 99
	if d[1][0] != 2 || d[0][1] != 2 {
		t.Error("rows alias each other")
	}
}

func TestFloydWarshall_errors(t *testing.T) {
	if _, err := FloydWarshall(nil); !errors.Is(err, ErrNilGraph) {
		t.Errorf("nil: %v", err)
	}
	cycle, _ := NewDirected(3, []Edge{{0, 1, 1}, {1, 2, 1}, {2, 0, -3}})
	if d, err := FloydWarshall(cycle); !errors.Is(err, ErrNegativeCycle) || d != nil {
		t.Errorf("negative cycle: (%v, %v)", d, err)
	}
	loop, _ := NewDirected(2, []Edge{{1, 1, -1}})
	if _, err := FloydWarshall(loop); !errors.Is(err, ErrNegativeCycle) {
		t.Errorf("negative self-loop: %v", err)
	}
	und, _ := NewUndirected(2, []Edge{{0, 1, -1}})
	if _, err := FloydWarshall(und); !errors.Is(err, ErrNegativeCycle) {
		t.Errorf("undirected negative edge: %v", err)
	}
	over, _ := NewDirected(3, []Edge{{0, 1, math.MaxInt64 - 1}, {1, 2, 1}})
	if _, err := FloydWarshall(over); !errors.Is(err, ErrOverflow) {
		t.Errorf("sum reaching MaxInt64: %v", err)
	}
}

// L41: the V² table must fit in 2⁴⁷ bytes (2³¹ on 32-bit); checked against exact 128-bit-free arithmetic.
func TestFloydWarshall_tooLarge(t *testing.T) {
	limit := uint64(maxAllocBytes) / 8 // int64 entries that fit
	for _, n := range []int{0, 1, 1 << 14, 1<<14 + 1, 1 << 20, 1 << 22, 1<<22 + 1, 1 << 30, math.MaxInt} {
		want := n == 0 || uint64(n) <= limit/uint64(n)
		if got := distanceTableFits(n); got != want {
			t.Errorf("distanceTableFits(%d) = %v, want %v", n, got, want)
		}
	}
	if testing.Short() || strconv.IntSize != 64 {
		t.Skip("the end-to-end case builds a graph of 2²²+1 vertices (about 130 MB)")
	}
	g, _ := NewDirected(1<<22+1, nil)
	if d, err := FloydWarshall(g); !errors.Is(err, ErrTooLarge) || d != nil {
		t.Errorf("V = 2²²+1: (%d rows, %v), want (nil, ErrTooLarge)", len(d), err)
	}
}

// Oracle: each row of FloydWarshall equals BellmanFord from that vertex (random graphs, weights ≥ 0
// directed or undirected, and negative-weight DAGs).
func TestFloydWarshall_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		n := 1 + rng.IntN(12)
		var g *Graph
		switch run % 3 {
		case 0:
			g, _ = NewDirected(n, randomEdges(rng, n, rng.IntN(3*n), 0, 30))
		case 1:
			g, _ = NewUndirected(n, randomEdges(rng, n, rng.IntN(3*n), 0, 30))
		default:
			var edges []Edge
			for range rng.IntN(3 * n) {
				a, b := rng.IntN(n), rng.IntN(n)
				if a != b {
					edges = append(edges, Edge{From: min(a, b), To: max(a, b), Weight: rng.Int64N(41) - 20})
				}
			}
			g, _ = NewDirected(n, edges)
		}
		d, err := FloydWarshall(g)
		if err != nil {
			t.Fatalf("run %d: %v", run, err)
		}
		for s := range n {
			want, _, _ := BellmanFord(g, s)
			if !slices.Equal(d[s], want) {
				t.Fatalf("run %d: row %d = %v, BellmanFord %v", run, s, d[s], want)
			}
		}
	}
}

func BenchmarkFloydWarshall(b *testing.B) {
	// Θ(V³): V = 100, 200, 400 (10⁶, 8·10⁶, 6.4·10⁷ inner steps).
	for _, size := range []struct {
		name string
		n    int
	}{{"V100", 100}, {"V200", 200}, {"V400", 400}} {
		g := benchGraph(size.n, true)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = FloydWarshall(g)
			}
		})
	}
}
