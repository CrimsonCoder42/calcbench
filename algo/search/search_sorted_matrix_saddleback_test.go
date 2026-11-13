package search

import (
	"math/rand/v2"
	"testing"
)

// Gries (1981) and Bird (MPC 2006) search a matrix sorted along rows and columns. The row-major Knuth
// matrix is also column-sorted, so 653 is at (2, 2) here too. Our own column-and-row-sorted 4×4 grid
// adds a case where rows overlap: 5 is at (1, 1), 15 is absent.
func TestSearchSortedMatrixSaddleback_example(t *testing.T) {
	if r, c, ok := SearchSortedMatrixSaddleback(knuthMatrix, 653); !ok || r != 2 || c != 2 {
		t.Errorf("Saddleback(knuth, 653) = (%d, %d, %v), want (2, 2, true)", r, c, ok)
	}
	grid := [][]int{{1, 4, 7, 11}, {2, 5, 8, 12}, {3, 6, 9, 16}, {10, 13, 14, 17}}
	if r, c, ok := SearchSortedMatrixSaddleback(grid, 5); !ok || r != 1 || c != 1 {
		t.Errorf("Saddleback(grid, 5) = (%d, %d, %v), want (1, 1, true)", r, c, ok)
	}
	if r, c, ok := SearchSortedMatrixSaddleback(grid, 15); ok || r != -1 || c != -1 {
		t.Errorf("Saddleback(grid, 15) = (%d, %d, %v), want (-1, -1, false)", r, c, ok)
	}
}

func TestSearchSortedMatrixSaddleback_cases(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]int
		target int
		found  bool
	}{
		{"no rows", nil, 1, false},
		{"one empty row", [][]int{{}}, 1, false},
		{"one cell", [][]int{{5}}, 5, true},
		{"one row", [][]int{{1, 2, 3}}, 1, true},
		{"one column", [][]int{{1}, {2}, {3}}, 3, true},
		{"duplicates", [][]int{{1, 2, 2}, {2, 2, 3}}, 2, true},
		{"below every element", [][]int{{1, 2}, {3, 4}}, 0, false},
		{"above every element", [][]int{{1, 2}, {3, 4}}, 5, false},
		{"negative", [][]int{{-9, -4}, {-8, -1}}, -8, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, c, ok := SearchSortedMatrixSaddleback(tt.matrix, tt.target)
			if ok != tt.found || (ok && tt.matrix[r][c] != tt.target) || (!ok && (r != -1 || c != -1)) {
				t.Errorf("Saddleback(%v, %d) = (%d, %d, %v), want found %v", tt.matrix, tt.target, r, c, ok, tt.found)
			}
		})
	}
}

func TestSearchSortedMatrixSaddleback_shortRowsDoNotPanic(t *testing.T) {
	matrix := [][]int{{1, 2, 3}, {4}, {}, {5, 6, 7}}
	for target := 0; target <= 8; target++ {
		SearchSortedMatrixSaddleback(matrix, target)
	}
}

func TestSearchSortedMatrixSaddleback_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { benchSink, _, _ = SearchSortedMatrixSaddleback(knuthMatrix, 653) }); allocs != 0 {
		t.Errorf("SearchSortedMatrixSaddleback allocates %v times per call, want 0", allocs)
	}
}

// youngTableau returns an m×n matrix with ascending rows and columns: cell (r, c) is the larger of the
// cells above and to the left plus a fixed-seed step in [0, 3), so duplicates occur.
func youngTableau(rng *rand.Rand, m, n int) [][]int {
	matrix := make([][]int, m)
	for r := range matrix {
		matrix[r] = make([]int, n)
		for c := range matrix[r] {
			base := 0
			if r > 0 {
				base = matrix[r-1][c]
			}
			if c > 0 {
				base = max(base, matrix[r][c-1])
			}
			matrix[r][c] = base + rng.IntN(3)
		}
	}
	return matrix
}

func TestSearchSortedMatrixSaddleback_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		matrix := youngTableau(rng, 1+rng.IntN(8), 1+rng.IntN(8))
		for target := -1; target <= 40; target++ {
			want := false
			for _, row := range matrix {
				for _, v := range row {
					want = want || v == target
				}
			}
			r, c, ok := SearchSortedMatrixSaddleback(matrix, target)
			if ok != want || (ok && matrix[r][c] != target) {
				t.Fatalf("trial %d: Saddleback(%d) = (%d, %d, %v), want found %v", trial, target, r, c, ok, want)
			}
		}
	}
}

// Square matrices of about 1e3, 1e4, 1e5 cells: O(m+n) grows with the side, ~√cells.
func BenchmarkSearchSortedMatrixSaddleback(b *testing.B) {
	for _, size := range matrixBenchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		matrix := youngTableau(rng, size.side, size.side)
		top := matrix[size.side-1][size.side-1]
		targets := make([]int, 1024)
		for i := range targets {
			targets[i] = rng.IntN(top + 1)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for i := 0; b.Loop(); i++ {
				r, _, _ := SearchSortedMatrixSaddleback(matrix, targets[i&1023])
				sink += r
			}
			benchSink = sink
		})
	}
}
