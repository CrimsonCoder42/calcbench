package matrix

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

// benchSink keeps benchmark results observable. Shared by every benchmark in this package.
var benchSink int

// gridEqual reports whether two grids have the same rows.
func gridEqual[T comparable](a, b [][]T) bool {
	return slices.EqualFunc(a, b, func(x, y []T) bool { return slices.Equal(x, y) })
}

// cloneGrid returns a deep copy of g.
func cloneGrid[T any](g [][]T) [][]T {
	out := make([][]T, len(g))
	for i, row := range g {
		out[i] = slices.Clone(row)
	}
	return out
}

// randomGrid returns an m×n grid of values in [lo, hi).
func randomGrid(rng *rand.Rand, m, n, lo, hi int) [][]int {
	g := make([][]int, m)
	for i := range g {
		g[i] = make([]int, n)
		for j := range g[i] {
			g[i][j] = lo + rng.IntN(hi-lo)
		}
	}
	return g
}

// The rotation (r, c) → (c, n−1−r) on a 3×3 grid, hand-derived from the definition (Knuth Vol. 1,
// §1.3.3 supplies the in-place cycle method, not this example).
func TestRotate90InPlace_threeByThree(t *testing.T) {
	g := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	if err := Rotate90InPlace(g); err != nil {
		t.Fatal(err)
	}
	if want := [][]int{{7, 4, 1}, {8, 5, 2}, {9, 6, 3}}; !gridEqual(g, want) {
		t.Errorf("got %v, want %v", g, want)
	}
}

func TestRotate90InPlace_cases(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		want [][]int
	}{
		{"empty", nil, nil},
		{"one", [][]int{{5}}, [][]int{{5}}},
		{"two", [][]int{{1, 2}, {3, 4}}, [][]int{{3, 1}, {4, 2}}},
		{"four", [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}, {13, 14, 15, 16}},
			[][]int{{13, 9, 5, 1}, {14, 10, 6, 2}, {15, 11, 7, 3}, {16, 12, 8, 4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := cloneGrid(tt.in)
			if err := Rotate90InPlace(g); err != nil || !gridEqual(g, tt.want) {
				t.Errorf("got %v, %v; want %v", g, err, tt.want)
			}
		})
	}
}

func TestRotate90InPlace_fourTurnsIsIdentity(t *testing.T) {
	g := [][]string{{"a", "b", "c"}, {"d", "e", "f"}, {"g", "h", "i"}}
	orig := cloneGrid(g)
	for range 4 {
		_ = Rotate90InPlace(g)
	}
	if !gridEqual(g, orig) {
		t.Errorf("four turns: %v", g)
	}
}

func TestRotate90InPlace_errors(t *testing.T) {
	for _, g := range [][][]int{{{1, 2}}, {{1, 2}, {3}}, {{1}, {2}}} {
		before := cloneGrid(g)
		if err := Rotate90InPlace(g); !errors.Is(err, ErrNotSquare) || !gridEqual(g, before) {
			t.Errorf("%v: err = %v, grid %v", before, err, g)
		}
	}
}

func TestRotate90InPlace_noAllocation(t *testing.T) {
	g := randomGrid(rand.New(rand.NewPCG(1, 2)), 8, 8, 0, 9)
	if allocs := testing.AllocsPerRun(100, func() { _ = Rotate90InPlace(g) }); allocs != 0 {
		t.Errorf("Rotate90InPlace allocates %v times per call, want 0", allocs)
	}
}

// TestRotate90InPlace_oracle: the lost alternative (transpose, then reverse each row) as the oracle.
func TestRotate90InPlace_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := rng.IntN(9)
		g := randomGrid(rng, n, n, -50, 50)
		want, _ := Transpose(nil, g)
		for _, row := range want {
			slices.Reverse(row)
		}
		if err := Rotate90InPlace(g); err != nil || !gridEqual(g, want) && n > 0 {
			t.Fatalf("trial %d: got %v, want %v", trial, g, want)
		}
	}
}

func BenchmarkRotate90InPlace(b *testing.B) {
	for _, n := range []int{32, 100, 316} { // ~1e3, 1e4, 1e5 cells
		g := randomGrid(rand.New(rand.NewPCG(1, 2)), n, n, 0, 100)
		b.Run(fmt.Sprint(n*n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = Rotate90InPlace(g)
			}
			benchSink = g[0][0]
		})
	}
}
