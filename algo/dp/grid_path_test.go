package dp

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// gridPathsSlow walks every right/down path recursively (grids up to 6×6).
func gridPathsSlow(blocked [][]bool, r, c int) uint64 {
	if r >= len(blocked) || c >= len(blocked[0]) || blocked[r][c] {
		return 0
	}
	if r == len(blocked)-1 && c == len(blocked[0])-1 {
		return 1
	}
	return gridPathsSlow(blocked, r+1, c) + gridPathsSlow(blocked, r, c+1)
}

// minPathSlow walks every right/down path recursively.
func minPathSlow(cost [][]int64, r, c int) int64 {
	rows, cols := len(cost), len(cost[0])
	if r == rows-1 && c == cols-1 {
		return cost[r][c]
	}
	best := int64(math.MaxInt64)
	if r+1 < rows {
		best = min(best, minPathSlow(cost, r+1, c))
	}
	if c+1 < cols {
		best = min(best, minPathSlow(cost, r, c+1))
	}
	return cost[r][c] + best
}

func openGrid(rows, cols int) [][]bool {
	g := make([][]bool, rows)
	for i := range g {
		g[i] = make([]bool, cols)
	}
	return g
}

// binomial computes C(n, k) exactly for small n (test oracle for the unobstructed grid).
func binomial(n, k int) uint64 {
	var c uint64 = 1
	for i := 1; i <= k; i++ {
		c = c * uint64(n-k+i) / uint64(i)
	}
	return c
}

// CLRS Exercise 22.4-2 counts paths in a DAG; on the unobstructed r×c grid the count is the binomial
// C(r+c−2, r−1) (Concrete Mathematics §5.1): 3×7 gives C(8, 2) = 28.
func TestGridPathCount_openGrid(t *testing.T) {
	if got, err := GridPathCount(openGrid(3, 7)); err != nil || got != 28 {
		t.Errorf("GridPathCount(3×7) = (%d, %v), want (28, nil)", got, err)
	}
	for r := 1; r <= 12; r++ {
		for c := 1; c <= 12; c++ {
			if got, err := GridPathCount(openGrid(r, c)); err != nil || got != binomial(r+c-2, r-1) {
				t.Errorf("GridPathCount(%d×%d) = (%d, %v), want %d", r, c, got, err, binomial(r+c-2, r-1))
			}
		}
	}
}

func TestGridPathCount_cases(t *testing.T) {
	centre := openGrid(3, 3)
	centre[1][1] = true
	startBlocked := openGrid(2, 2)
	startBlocked[0][0] = true
	endBlocked := openGrid(2, 2)
	endBlocked[1][1] = true
	wall := openGrid(3, 3)
	wall[0][1], wall[1][1], wall[2][1] = true, true, true
	tests := []struct {
		name string
		grid [][]bool
		want uint64
	}{
		{"no rows", nil, 0},
		{"no columns", [][]bool{{}, {}}, 0},
		{"one open cell", [][]bool{{false}}, 1},
		{"one blocked cell", [][]bool{{true}}, 0},
		{"one row", openGrid(1, 9), 1},
		{"centre blocked", centre, 2},
		{"start blocked", startBlocked, 0},
		{"end blocked", endBlocked, 0},
		{"wall", wall, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GridPathCount(tt.grid); err != nil || got != tt.want {
				t.Errorf("GridPathCount = (%d, %v), want %d", got, err, tt.want)
			}
		})
	}
}

func TestGridPathCount_errors(t *testing.T) {
	if _, err := GridPathCount([][]bool{{false, false}, {false}}); !errors.Is(err, ErrNotRectangular) {
		t.Errorf("ragged: err = %v", err)
	}
	// C(78, 39) ≈ 2.7·10²² exceeds 2⁶⁴.
	if _, err := GridPathCount(openGrid(40, 40)); !errors.Is(err, ErrOverflow) {
		t.Errorf("40×40: err = %v, want ErrOverflow", err)
	}
}

func TestGridPathCount_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		g := openGrid(1+rng.IntN(6), 1+rng.IntN(6))
		for _, row := range g {
			for j := range row {
				row[j] = rng.IntN(5) == 0
			}
		}
		if got, err := GridPathCount(g); err != nil || got != gridPathsSlow(g, 0, 0) {
			t.Fatalf("trial %d: GridPathCount(%v) = (%d, %v), want %d", trial, g, got, err, gridPathsSlow(g, 0, 0))
		}
	}
}

// Worked by hand: in [[2 9 1] [4 1 8] [7 3 2]] the cheapest path is 2→4→1→3→2 = 12.
func TestMinPathSumGrid_example(t *testing.T) {
	cost := [][]int64{{2, 9, 1}, {4, 1, 8}, {7, 3, 2}}
	if got, err := MinPathSumGrid(cost); err != nil || got != 12 {
		t.Errorf("MinPathSumGrid = (%d, %v), want (12, nil)", got, err)
	}
}

func TestMinPathSumGrid_cases(t *testing.T) {
	tests := []struct {
		name string
		cost [][]int64
		want int64
	}{
		{"one cell", [][]int64{{-4}}, -4},
		{"one row", [][]int64{{1, 2, 3}}, 6},
		{"one column", [][]int64{{1}, {2}, {3}}, 6},
		{"negative costs", [][]int64{{0, -5}, {-1, 0}}, -5},
		{"zeros", [][]int64{{0, 0}, {0, 0}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MinPathSumGrid(tt.cost); err != nil || got != tt.want {
				t.Errorf("MinPathSumGrid(%v) = (%d, %v), want %d", tt.cost, got, err, tt.want)
			}
		})
	}
}

func TestMinPathSumGrid_errors(t *testing.T) {
	tests := []struct {
		name string
		cost [][]int64
		want error
	}{
		{"no rows", nil, ErrEmpty},
		{"no columns", [][]int64{{}}, ErrEmpty},
		{"ragged", [][]int64{{1, 2}, {3}}, ErrNotRectangular},
		{"overflow", [][]int64{{math.MaxInt64, 1}}, ErrOverflow},
		{"negative overflow", [][]int64{{math.MinInt64}, {-1}}, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := MinPathSumGrid(tt.cost); !errors.Is(err, tt.want) {
				t.Errorf("MinPathSumGrid error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestMinPathSumGrid_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		rows, cols := 1+rng.IntN(6), 1+rng.IntN(6)
		cost := make([][]int64, rows)
		for i := range cost {
			cost[i] = make([]int64, cols)
			for j := range cost[i] {
				cost[i][j] = rng.Int64N(41) - 20
			}
		}
		if got, err := MinPathSumGrid(cost); err != nil || got != minPathSlow(cost, 0, 0) {
			t.Fatalf("trial %d: MinPathSumGrid(%v) = (%d, %v), want %d", trial, cost, got, err, minPathSlow(cost, 0, 0))
		}
	}
}

// Square grids of about 10³, 10⁴ and 10⁵ cells (sides 32, 100, 316).
var gridSides = []struct {
	name string
	side int
}{{"1e3", 32}, {"1e4", 100}, {"1e5", 316}}

func BenchmarkGridPathCount(b *testing.B) {
	for _, size := range gridSides {
		g := openGrid(size.side, size.side)
		// Every path crosses the anti-diagonal i+j = 3, which is blocked: counts stay 0 (no overflow)
		// while every cell is still visited.
		for i := range g {
			for j := range g[i] {
				g[i][j] = (i+j)%7 == 3 && i != j
			}
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = GridPathCount(g)
			}
		})
	}
}

func BenchmarkMinPathSumGrid(b *testing.B) {
	for _, size := range gridSides {
		cost := make([][]int64, size.side)
		for i := range cost {
			cost[i] = benchInt64s(size.side, 100)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MinPathSumGrid(cost)
			}
		})
	}
}
