package graph

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand example (catalog/graph.md), a 4- and 8-connected fill in the sense of Smith, "Tint fill" (1979).
func TestFloodFillInPlace_handExample(t *testing.T) {
	grid := toGrid(
		"aab",
		"abb",
		"bba",
	)
	n, err := FloodFillInPlace(grid, Cell{Row: 1, Col: 1}, 'c', FourNeighbors)
	if want := toGrid("aac", "acc", "cca"); n != 5 || err != nil || !slices.EqualFunc(grid, want, slices.Equal) {
		t.Errorf("fill b from (1,1): %d %v %q, want 5 %q", n, err, grid, want)
	}
	grid = toGrid("aab", "abb", "bba")
	n, _ = FloodFillInPlace(grid, Cell{}, 'z', EightNeighbors)
	if want := toGrid("zzb", "zbb", "bba"); n != 3 || !slices.EqualFunc(grid, want, slices.Equal) {
		t.Errorf("fill a from (0,0), 8 neighbours: %d %q, want 3 %q", n, grid, want)
	}
}

func TestFloodFillInPlace_cases(t *testing.T) {
	grid := toGrid("xx", "xx")
	if n, err := FloodFillInPlace(grid, Cell{Row: 1, Col: 1}, 'x', FourNeighbors); n != 0 || err != nil {
		t.Errorf("same colour: (%d, %v), want (0, nil)", n, err)
	}
	if n, _ := FloodFillInPlace(grid, Cell{}, 'y', FourNeighbors); n != 4 {
		t.Errorf("whole grid: %d, want 4", n)
	}
	one := [][]int{{7}}
	if n, _ := FloodFillInPlace(one, Cell{}, -1, EightNeighbors); n != 1 || one[0][0] != -1 {
		t.Errorf("1×1: %d %v", n, one)
	}
	nan := [][]float64{{math.NaN(), math.NaN()}}
	if n, _ := FloodFillInPlace(nan, Cell{}, 0, FourNeighbors); n != 1 || nan[0][0] != 0 || !math.IsNaN(nan[0][1]) {
		t.Errorf("NaN start matches nothing, so only the start changes: %d %v", n, nan)
	}
}

func TestFloodFillInPlace_errors(t *testing.T) {
	grid := toGrid("ab", "cd")
	before := toGrid("ab", "cd")
	if _, err := FloodFillInPlace(grid, Cell{}, 'z', Neighborhood(5)); !errors.Is(err, ErrNeighborhood) {
		t.Errorf("neighborhood: %v", err)
	}
	if _, err := FloodFillInPlace(toGrid("ab", "c"), Cell{}, 'z', FourNeighbors); !errors.Is(err, ErrNotRectangular) {
		t.Errorf("ragged: %v", err)
	}
	if _, err := FloodFillInPlace(grid, Cell{Row: 2}, 'z', FourNeighbors); !errors.Is(err, ErrCellOutOfRange) {
		t.Errorf("start out of range: %v", err)
	}
	if !slices.EqualFunc(grid, before, slices.Equal) {
		t.Error("grid changed on error")
	}
}

// Oracle: the changed cells are exactly GridBFS's region of the start value, now holding replacement.
func TestFloodFillInPlace_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		rows, cols := 1+rng.IntN(8), 1+rng.IntN(8)
		grid := make([][]byte, rows)
		for r := range grid {
			grid[r] = make([]byte, cols)
			for c := range grid[r] {
				grid[r][c] = "abc"[rng.IntN(3)]
			}
		}
		nb := FourNeighbors
		if run%2 == 1 {
			nb = EightNeighbors
		}
		start := Cell{Row: rng.IntN(rows), Col: rng.IntN(cols)}
		orig := grid[start.Row][start.Col]
		region := slices.Collect(must(GridBFS(grid, start, func(b byte) bool { return b == orig }, nb)))
		before := slices.Clone(grid)
		for r := range before {
			before[r] = slices.Clone(grid[r])
		}
		n, err := FloodFillInPlace(grid, start, 'z', nb)
		if err != nil || n != len(region) {
			t.Fatalf("run %d: changed %d (%v), region %d", run, n, err, len(region))
		}
		for r := range rows {
			for c := range cols {
				in := slices.Contains(region, Cell{Row: r, Col: c})
				if (in && grid[r][c] != 'z') || (!in && grid[r][c] != before[r][c]) {
					t.Fatalf("run %d: cell (%d,%d) wrong", run, r, c)
				}
			}
		}
	}
}

func BenchmarkFloodFillInPlace(b *testing.B) {
	for _, size := range gridBenchSides {
		grid := openGrid(size.side)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				// Alternate the colour so every iteration repaints the same region.
				_, _ = FloodFillInPlace(grid, Cell{}, ',', FourNeighbors)
				_, _ = FloodFillInPlace(grid, Cell{}, '.', FourNeighbors)
			}
		})
	}
}
