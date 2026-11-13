package stack

import (
	"errors"
	"math/rand/v2"
	"testing"
)

// parseGrid turns rows of '0'/'1' into a bool grid.
func parseGrid(rows ...string) [][]bool {
	grid := make([][]bool, len(rows))
	for r, row := range rows {
		grid[r] = make([]bool, len(row))
		for c := range row {
			grid[r][c] = row[c] == '1'
		}
	}
	return grid
}

// rectangleOracle tries every rectangle: O(R²C²) rectangles, each checked cell by cell.
func rectangleOracle(grid [][]bool) int {
	best := 0
	for r0 := range grid {
		for c0 := range grid[r0] {
			for r1 := r0; r1 < len(grid); r1++ {
				for c1 := c0; c1 < len(grid[r0]); c1++ {
					all := true
					for r := r0; r <= r1 && all; r++ {
						for c := c0; c <= c1; c++ {
							if !grid[r][c] {
								all = false
								break
							}
						}
					}
					if all {
						best = max(best, (r1-r0+1)*(c1-c0+1))
					}
				}
			}
		}
	}
	return best
}

// Hand-derived (catalog), our own grid: the 4×2 block in columns 1–2 (beats the 2×3 block in rows 1–2).
func TestMaximalRectangle_example(t *testing.T) {
	grid := parseGrid("01101", "11101", "11111", "01110")
	if got, err := MaximalRectangle(grid); err != nil || got != 8 {
		t.Errorf("MaximalRectangle = (%d, %v), want (8, nil)", got, err)
	}
}

func TestMaximalRectangle_cases(t *testing.T) {
	tests := []struct {
		name string
		grid [][]bool
		want int
	}{
		{"empty", nil, 0},
		{"empty rows", [][]bool{{}, {}}, 0},
		{"one false", parseGrid("0"), 0},
		{"one true", parseGrid("1"), 1},
		{"all true", parseGrid("111", "111"), 6},
		{"all false", parseGrid("000", "000"), 0},
		{"one column", parseGrid("1", "1", "0", "1"), 2},
		{"one row", parseGrid("0111011"), 3},
		{"tall beats wide", parseGrid("1100", "1100", "1111"), 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MaximalRectangle(tt.grid); err != nil || got != tt.want {
				t.Errorf("MaximalRectangle(%v) = (%d, %v), want (%d, nil)", tt.grid, got, err, tt.want)
			}
		})
	}
}

func TestMaximalRectangle_errors(t *testing.T) {
	if got, err := MaximalRectangle(parseGrid("11", "1")); !errors.Is(err, ErrRagged) || got != 0 {
		t.Errorf("MaximalRectangle(ragged) = (%d, %v), want (0, ErrRagged)", got, err)
	}
}

func TestMaximalRectangle_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		rows, cols := rng.IntN(7), rng.IntN(7)
		grid := make([][]bool, rows)
		for r := range grid {
			grid[r] = make([]bool, cols)
			for c := range grid[r] {
				grid[r][c] = rng.IntN(4) != 0
			}
		}
		want := rectangleOracle(grid)
		if got, err := MaximalRectangle(grid); err != nil || got != want {
			t.Fatalf("trial %d: MaximalRectangle(%v) = (%d, %v), want %d", trial, grid, got, err, want)
		}
	}
}

// Sizes are cell counts: square grids of about 1e3, 1e4, 1e5 cells.
func BenchmarkMaximalRectangle(b *testing.B) {
	for _, size := range []struct {
		name string
		side int
	}{{"1e3", 32}, {"1e4", 100}, {"1e5", 316}} {
		rng := rand.New(rand.NewPCG(1, 2))
		grid := make([][]bool, size.side)
		for r := range grid {
			grid[r] = make([]bool, size.side)
			for c := range grid[r] {
				grid[r][c] = rng.IntN(5) != 0
			}
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = MaximalRectangle(grid)
			}
		})
	}
}
