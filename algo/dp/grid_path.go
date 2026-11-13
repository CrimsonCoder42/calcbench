package dp

import "fmt"

// GridPathCount returns the number of paths from the top-left cell to the bottom-right cell of a grid
// that move one cell right or one cell down at a time and never enter a blocked cell
// (blocked[row][col] = true). With no blocked cell the count is the binomial C(r+c−2, r−1). A grid
// with no cells, or whose first or last cell is blocked, has 0 paths; a 1×1 open grid has 1.
// blocked is read, never modified.
// Errors: ErrNotRectangular when rows differ in length; ErrOverflow when the count for some cell
// exceeds 2⁶⁴−1 (a cell's count is part of the recurrence's state even if the final count fits, as L28).
//
// ID:      ALG-DP-013
// Time:    O(r·c)   r = rows, c = columns
// Space:   O(c) extra (one row of counts, one allocation)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Exercise 22.4-2 (the number of paths between two vertices of a DAG in linear time: a vertex's count is the sum over its predecessors, in topological order; here row-major order)
// Source:  Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §5.1 (the addition formula, Pascal's rule, which the unobstructed counts satisfy)
// Tags:    dp, grid, lattice-path, counting, obstacles, unique-paths
func GridPathCount(blocked [][]bool) (uint64, error) {
	cols, err := gridWidth(blocked)
	if err != nil || cols == 0 {
		return 0, err
	}
	ways := make([]uint64, cols) // ways[j] = paths to (current row, j)
	for r, row := range blocked {
		for j, isBlocked := range row {
			switch {
			case isBlocked:
				ways[j] = 0
			case r == 0 && j == 0:
				ways[j] = 1
			case j > 0:
				s, ok := addUint64(ways[j], ways[j-1])
				if !ok {
					return 0, fmt.Errorf("%w: paths to cell (%d, %d)", ErrOverflow, r, j)
				}
				ways[j] = s
			}
		}
	}
	return ways[cols-1], nil
}

// MinPathSumGrid returns the smallest sum of the cells on a path from the top-left cell to the
// bottom-right cell of cost that moves one cell right or one cell down at a time; both end cells count.
// Costs may be negative (every path has the same number of steps, and the grid is acyclic).
// cost is read, never modified.
// Errors: ErrEmpty when the grid has no cells; ErrNotRectangular when rows differ in length; ErrOverflow
// when a partial path sum the recurrence holds does not fit in int64 (as L28).
//
// ID:      ALG-DP-014
// Time:    O(r·c)   r = rows, c = columns
// Space:   O(c) extra (one row, one allocation)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §24.2 (shortest paths in a directed acyclic graph by relaxing in topological order; row-major order is topological for right and down moves, with each cell's cost on its incoming edges)
// Tags:    dp, grid, shortest-path, dag, min-path-sum
func MinPathSumGrid(cost [][]int64) (int64, error) {
	cols, err := gridWidth(cost)
	if err != nil {
		return 0, err
	}
	if cols == 0 {
		return 0, fmt.Errorf("%w: grid has no cells", ErrEmpty)
	}
	dist := make([]int64, cols) // dist[j] = least sum to (current row, j)
	for r, row := range cost {
		for j, v := range row {
			var from int64
			switch {
			case r == 0 && j == 0:
				from = 0
			case r == 0:
				from = dist[j-1]
			case j == 0:
				from = dist[j]
			default:
				from = min(dist[j], dist[j-1])
			}
			s, ok := addInt64(from, v)
			if !ok {
				return 0, fmt.Errorf("%w: path sum to cell (%d, %d)", ErrOverflow, r, j)
			}
			dist[j] = s
		}
	}
	return dist[cols-1], nil
}

// gridWidth returns the common row length of grid, or ErrNotRectangular. A grid with no rows has
// width 0.
func gridWidth[T any](grid [][]T) (int, error) {
	if len(grid) == 0 {
		return 0, nil
	}
	cols := len(grid[0])
	for r, row := range grid {
		if len(row) != cols {
			return 0, fmt.Errorf("%w: row 0 has %d cells, row %d has %d", ErrNotRectangular, cols, r, len(row))
		}
	}
	return cols, nil
}
