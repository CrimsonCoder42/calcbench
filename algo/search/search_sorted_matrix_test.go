package search

import (
	"math/rand/v2"
	"testing"
)

// knuthMatrix is the 16 sorted keys of Knuth, TAOCP Vol. 3, 2nd ed., §6.2.1 laid out 4 by 4, row by row.
var knuthMatrix = [][]int{
	{61, 87, 154, 170},
	{275, 426, 503, 509},
	{512, 612, 653, 677},
	{703, 765, 897, 908},
}

// Knuth §6.2.1: 653 is the 11th key, so row 2, column 2 in the 4×4 layout; 400 is absent.
func TestSearchSortedMatrix_knuthExample(t *testing.T) {
	if r, c, ok := SearchSortedMatrix(knuthMatrix, 653); !ok || r != 2 || c != 2 {
		t.Errorf("SearchSortedMatrix(653) = (%d, %d, %v), want (2, 2, true)", r, c, ok)
	}
	if r, c, ok := SearchSortedMatrix(knuthMatrix, 400); ok || r != -1 || c != -1 {
		t.Errorf("SearchSortedMatrix(400) = (%d, %d, %v), want (-1, -1, false)", r, c, ok)
	}
}

// Our own data: a 3×4 grid that is ascending read row by row; 19 is at (1, 2), 14 is absent.
func TestSearchSortedMatrix_exercises(t *testing.T) {
	grid := [][]int{{2, 4, 8, 9}, {12, 15, 19, 21}, {26, 28, 33, 40}}
	if r, c, ok := SearchSortedMatrix(grid, 19); !ok || r != 1 || c != 2 {
		t.Errorf("SearchSortedMatrix(19) = (%d, %d, %v), want (1, 2, true)", r, c, ok)
	}
	if _, _, ok := SearchSortedMatrix(grid, 14); ok {
		t.Error("SearchSortedMatrix(14) found an absent value")
	}
}

func TestSearchSortedMatrix_cases(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]int
		target int
		row    int
		col    int
	}{
		{"no rows", nil, 1, -1, -1},
		{"one cell, present", [][]int{{5}}, 5, 0, 0},
		{"one cell, absent", [][]int{{5}}, 4, -1, -1},
		{"one row", [][]int{{1, 2, 3}}, 3, 0, 2},
		{"one column", [][]int{{1}, {2}, {3}}, 2, 1, 0},
		{"below every element", knuthMatrix, 1, -1, -1},
		{"above every element", knuthMatrix, 1000, -1, -1},
		{"first cell", knuthMatrix, 61, 0, 0},
		{"last cell", knuthMatrix, 908, 3, 3},
		{"first of a row", knuthMatrix, 512, 2, 0},
		{"gap between rows", knuthMatrix, 510, -1, -1},
		{"ragged rows", [][]int{{1}, {4, 5, 6}, {9, 10}}, 10, 2, 1},
		{"negative", [][]int{{-9, -7}, {-3, 0}}, -7, 0, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, c, ok := SearchSortedMatrix(tt.matrix, tt.target)
			if r != tt.row || c != tt.col || ok != (tt.row >= 0) {
				t.Errorf("SearchSortedMatrix(%v, %d) = (%d, %d, %v), want (%d, %d)", tt.matrix, tt.target, r, c, ok, tt.row, tt.col)
			}
		})
	}
}

func TestSearchSortedMatrix_emptyRowsDoNotPanic(t *testing.T) {
	matrix := [][]int{{1, 2}, {}, {5, 6}, {}}
	for target := 0; target <= 7; target++ {
		SearchSortedMatrix(matrix, target)
	}
}

func TestSearchSortedMatrix_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { benchSink, _, _ = SearchSortedMatrix(knuthMatrix, 653) }); allocs != 0 {
		t.Errorf("SearchSortedMatrix allocates %v times per call, want 0", allocs)
	}
}

// sortedMatrix lays out m·n strictly ascending fixed-seed values row by row.
func sortedMatrix(rng *rand.Rand, m, n int) [][]int {
	matrix := make([][]int, m)
	v := 0
	for r := range matrix {
		matrix[r] = make([]int, n)
		for c := range matrix[r] {
			v += 1 + rng.IntN(3)
			matrix[r][c] = v
		}
	}
	return matrix
}

func TestSearchSortedMatrix_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		matrix := sortedMatrix(rng, 1+rng.IntN(8), 1+rng.IntN(8))
		for target := 0; target <= 200; target++ {
			wantR, wantC := -1, -1
			for r, row := range matrix {
				for c, v := range row {
					if v == target {
						wantR, wantC = r, c
					}
				}
			}
			if r, c, ok := SearchSortedMatrix(matrix, target); r != wantR || c != wantC || ok != (wantR >= 0) {
				t.Fatalf("trial %d: SearchSortedMatrix(%d) = (%d, %d, %v), want (%d, %d)", trial, target, r, c, ok, wantR, wantC)
			}
		}
	}
}

// Square matrices of about 1e3, 1e4, 1e5 cells.
var matrixBenchSizes = []struct {
	name string
	side int
}{{"1e3", 32}, {"1e4", 100}, {"1e5", 317}}

func BenchmarkSearchSortedMatrix(b *testing.B) {
	for _, size := range matrixBenchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		matrix := sortedMatrix(rng, size.side, size.side)
		top := matrix[size.side-1][size.side-1]
		targets := make([]int, 1024)
		for i := range targets {
			targets[i] = rng.IntN(top + 1)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for i := 0; b.Loop(); i++ {
				r, _, _ := SearchSortedMatrix(matrix, targets[i&1023])
				sink += r
			}
			benchSink = sink
		})
	}
}
