package matrix

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
)

// The reduction of Bentley, Programming Pearls, Column 8 (pairs of rows) with the ≤ k bound; our own
// example, checked by listing the rectangles by hand: in [[1 0 1] [0 −2 3]] the best sum ≤ 2 is 2
// (e.g. the right 2×2 block, 0+1−2+3), and the best ≤ 3 is 3 (the single cell 3).
func TestMaxRectangleSumNoLargerThanK_workedExample(t *testing.T) {
	g := [][]int{{1, 0, 1}, {0, -2, 3}}
	for k, want := range map[int]int{2: 2, 3: 3, 4: 4, 100: 4, -2: -2, 0: 0} {
		if got, err := MaxRectangleSumNoLargerThanK(g, k); err != nil || got != want {
			t.Errorf("k=%d: got %d, %v; want %d", k, got, err, want)
		}
	}
}

func TestMaxRectangleSumNoLargerThanK_cases(t *testing.T) {
	tests := []struct {
		name string
		grid [][]int
		k    int
		want int
	}{
		{"one cell equal to k", [][]int{{5}}, 5, 5},
		{"one cell below k", [][]int{{-5}}, 0, -5},
		{"row", [][]int{{2, -1, 2}}, 2, 2},
		{"column (tall grid uses the other side)", [][]int{{2}, {-1}, {2}, {4}, {-3}}, 3, 3},
		{"all negative", [][]int{{-3, -1}, {-2, -4}}, -2, -2},
		{"all equal", [][]int{{1, 1}, {1, 1}}, 3, 2},
		{"whole grid is best", [][]int{{1, 2}, {3, 4}}, 10, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MaxRectangleSumNoLargerThanK(tt.grid, tt.k); err != nil || got != tt.want {
				t.Errorf("got %d, %v; want %d", got, err, tt.want)
			}
		})
	}
}

func TestMaxRectangleSumNoLargerThanK_errors(t *testing.T) {
	tests := []struct {
		name string
		grid [][]int64
		k    int64
		want error
	}{
		{"no rows", nil, 0, ErrEmpty},
		{"no columns", [][]int64{{}, {}}, 0, ErrEmpty},
		{"ragged", [][]int64{{1}, {1, 2}}, 0, ErrRagged},
		{"nothing at most k", [][]int64{{5, 6}}, 4, ErrNotFound},
		{"line sum overflows", [][]int64{{math.MaxInt64, 0}, {1, 0}}, -1, ErrOverflow},
		{"prefix sum overflows", [][]int64{{math.MaxInt64, 1}}, -1, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := MaxRectangleSumNoLargerThanK(tt.grid, tt.k); !errors.Is(err, tt.want) {
				t.Errorf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

// Small element types: sums are carried in int64; an answer outside int8 is ErrOverflow, one inside is fine.
func TestMaxRectangleSumNoLargerThanK_int8(t *testing.T) {
	if got, err := MaxRectangleSumNoLargerThanK([][]int8{{100, 100}}, 127); err != nil || got != 100 {
		t.Errorf("got %d, %v; want 100", got, err)
	}
	if _, err := MaxRectangleSumNoLargerThanK([][]int8{{-100, -100}}, -128); !errors.Is(err, ErrOverflow) {
		t.Errorf("answer −200: err = %v, want ErrOverflow", err)
	}
}

// Extremes: the differences of prefix sums would overflow int64 but are compared exactly.
func TestMaxRectangleSumNoLargerThanK_extremes(t *testing.T) {
	g := [][]int64{{math.MaxInt64, math.MinInt64, math.MaxInt64}}
	// Prefix sums 0, Max, −1, Max−1; best ≤ Max is Max (the first cell).
	if got, err := MaxRectangleSumNoLargerThanK(g, math.MaxInt64); err != nil || got != math.MaxInt64 {
		t.Errorf("got %d, %v; want MaxInt64", got, err)
	}
	// Best ≤ MinInt64 is the middle cell.
	if got, err := MaxRectangleSumNoLargerThanK(g, math.MinInt64); err != nil || got != math.MinInt64 {
		t.Errorf("got %d, %v; want MinInt64", got, err)
	}
}

// bruteForce is the slow obvious oracle: every rectangle, summed from scratch.
func bruteForce(g [][]int, k int) (int, bool) {
	best, found := 0, false
	for r1 := range g {
		for r2 := r1; r2 < len(g); r2++ {
			for c1 := range g[0] {
				for c2 := c1; c2 < len(g[0]); c2++ {
					s := 0
					for r := r1; r <= r2; r++ {
						for c := c1; c <= c2; c++ {
							s += g[r][c]
						}
					}
					if s <= k && (!found || s > best) {
						best, found = s, true
					}
				}
			}
		}
	}
	return best, found
}

func TestMaxRectangleSumNoLargerThanK_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		m, n := 1+rng.IntN(6), 1+rng.IntN(6)
		g := randomGrid(rng, m, n, -10, 11)
		before := cloneGrid(g)
		k := rng.IntN(41) - 20
		want, found := bruteForce(g, k)
		got, err := MaxRectangleSumNoLargerThanK(g, k)
		switch {
		case !found && !errors.Is(err, ErrNotFound):
			t.Fatalf("trial %d: %v k=%d: got %d, %v; want ErrNotFound", trial, g, k, got, err)
		case found && (err != nil || got != want):
			t.Fatalf("trial %d: %v k=%d: got %d, %v; want %d", trial, g, k, got, err, want)
		}
		if !gridEqual(g, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

// BenchmarkMaxRectangleSumNoLargerThanK: square n×n grids, O(n³ log n): n = 10, 32, 100 (1e2…1e4 cells).
// Every cell is even and k is odd, so no sum equals k and the early exit never fires: the full work.
func BenchmarkMaxRectangleSumNoLargerThanK(b *testing.B) {
	for _, n := range []int{10, 32, 100} {
		g := randomGrid(rand.New(rand.NewPCG(1, 2)), n, n, -50, 51)
		for _, row := range g {
			for j := range row {
				row[j] *= 2
			}
		}
		b.Run(fmt.Sprint(n*n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				v, _ := MaxRectangleSumNoLargerThanK(g, 17)
				benchSink = v
			}
		})
	}
}
