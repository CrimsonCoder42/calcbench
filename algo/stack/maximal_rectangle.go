package stack

import (
	"errors"
	"fmt"
)

// ErrRagged: the rows of a grid have different lengths.
var ErrRagged = errors.New("stack: grid rows must all have the same length")

// MaximalRectangle returns the area, in cells, of the largest axis-aligned rectangle of grid whose cells
// are all true. grid[r][c] is row r, column c. An empty grid, or one with no true cell, gives 0.
// Each row is the base of a histogram whose bar in column c is the run of trues ending there; the
// answer is the largest of those histograms' rectangles. grid is read, never modified.
// Errors: ErrRagged when rows differ in length.
//
// ID:      ALG-STACK-006
// Time:    O(R·C)   R = len(grid), C = len(grid[0])
// Space:   O(C) extra for the bar heights and the stack (two allocations)
// Source:  Berkman, Schieber, Vishkin, "Optimal doubly logarithmic parallel algorithms based on finding all nearest smaller values", Journal of Algorithms, 1993, 14(3):344–370, doi:10.1006/jagm.1993.1018 (the histogram step)
// Source:  Naamad, Lee, Hsu, "On the maximum empty rectangle problem", Discrete Applied Mathematics, 1984, 8(3):267–277, doi:10.1016/0166-218X(84)90124-0 (the problem family; row reduction proved in the catalog)
// Tags:    stack, monotonic-stack, matrix, grid, binary-matrix, largest-rectangle, area
func MaximalRectangle(grid [][]bool) (int, error) {
	if len(grid) == 0 {
		return 0, nil
	}
	cols := len(grid[0])
	for r, row := range grid {
		if len(row) != cols {
			return 0, fmt.Errorf("%w: row %d has %d cells, row 0 has %d", ErrRagged, r, len(row), cols)
		}
	}
	heights := make([]int64, cols)
	stack := make([]int, 0, cols+1)
	var best int64
	for _, row := range grid {
		for c, cell := range row {
			if cell {
				heights[c]++
			} else {
				heights[c] = 0
			}
		}
		// Heights are at most R and widths at most C; R·C cells already exist in memory, so no overflow.
		area, s, _ := largestRectangle(heights, stack)
		stack = s
		best = max(best, area)
	}
	return int(best), nil
}
