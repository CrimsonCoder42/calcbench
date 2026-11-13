package matrix

import "fmt"

// gridShape returns the number of rows and the common row length of grid, or ErrRagged.
func gridShape[T any](grid [][]T) (rows, cols int, err error) {
	if len(grid) == 0 {
		return 0, 0, nil
	}
	cols = len(grid[0])
	for r, row := range grid {
		if len(row) != cols {
			return 0, 0, fmt.Errorf("%w: row 0 has %d cells, row %d has %d", ErrRagged, cols, r, len(row))
		}
	}
	return len(grid), cols, nil
}

// newGrid returns dst when it already has rows×cols shape, else a fresh rows×cols grid on one backing
// array. rows·cols never exceeds the cell count of an input grid that exists, so it cannot overflow.
func newGrid[T any](dst [][]T, rows, cols int) [][]T {
	if len(dst) == rows {
		ok := true
		for _, row := range dst {
			if len(row) != cols {
				ok = false
				break
			}
		}
		if ok {
			return dst
		}
	}
	cells := make([]T, rows*cols)
	out := make([][]T, rows)
	for r := range out {
		out[r] = cells[r*cols : (r+1)*cols : (r+1)*cols]
	}
	return out
}

// Transpose returns the transpose of the m×n grid, an n×m grid with out[c][r] = grid[r][c]. dst is
// reused when it is already n×m, else a new grid is allocated; dst must not share memory with grid.
// grid is read, never modified. An m×0 grid's transpose has no rows.
// Errors: ErrRagged.
//
// ID:      ALG-MATRIX-002
// Time:    O(mn)   m = rows, n = columns
// Space:   O(1) extra; O(mn) for the result, two allocations unless dst is reused
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Appendix D.1 (the transpose, (Aᵀ)ᵢⱼ = aⱼᵢ)
// Tags:    matrix, transpose, grid
func Transpose[T any](dst, grid [][]T) ([][]T, error) {
	m, n, err := gridShape(grid)
	if err != nil {
		return dst, err
	}
	out := newGrid(dst, n, m)
	for r, row := range grid {
		for c, v := range row {
			out[c][r] = v
		}
	}
	return out, nil
}
