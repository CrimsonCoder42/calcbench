package graph

import "iter"

// GridDFS returns the cells reachable from start through passable cells, in depth-first preorder (each
// cell when first discovered, trying neighbours in nb's order), moving to the 4 or 8 neighbours that nb
// names. It uses an explicit stack, so a long winding path cannot overflow the call stack. A start cell
// that is not passable gives an empty sequence. Each range over the sequence runs a fresh search and reads
// the grid as it goes; breaking out stops it. grid is never modified.
// Errors (returned before any iteration): ErrNeighborhood; ErrNotRectangular; ErrTooLarge;
// ErrCellOutOfRange when start is outside the grid; ErrNilPredicate when passable is nil.
//
// ID:      ALG-GRAPH-023
// Time:    O(R·C) for a full iteration   R = rows, C = columns
// Space:   O(R·C) extra per iteration (a visited bitmap and a stack of at most R·C frames)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.3 (DFS-VISIT, here on the implicit grid graph)
// Source:  The Go Programming Language Specification, "For statements with range clause" (range over function iterators); package iter (iter.Seq)
// Tags:    grid, dfs, depth-first, traversal, iterator, maze, implicit-graph
func GridDFS[T any](grid [][]T, start Cell, passable func(T) bool, nb Neighborhood) (iter.Seq[Cell], error) {
	rows, cols, steps, err := gridArgs(grid, start, nb)
	if err != nil {
		return nil, err
	}
	if passable == nil {
		return nil, ErrNilPredicate
	}
	return func(yield func(Cell) bool) {
		if !passable(grid[start.Row][start.Col]) {
			return
		}
		// frame is a cell on the current path and the index of its next neighbour step to try.
		type frame struct {
			c    Cell
			next int
		}
		seen := make([]bool, rows*cols)
		seen[start.Row*cols+start.Col] = true
		if !yield(start) {
			return
		}
		stack := []frame{{start, 0}}
		for len(stack) > 0 {
			top := len(stack) - 1
			if stack[top].next == len(steps) {
				stack = stack[:top]
				continue
			}
			c, d := stack[top].c, steps[stack[top].next]
			stack[top].next++
			r, k := c.Row+d.Row, c.Col+d.Col
			if r < 0 || r >= rows || k < 0 || k >= cols || seen[r*cols+k] || !passable(grid[r][k]) {
				continue
			}
			seen[r*cols+k] = true
			if !yield(Cell{Row: r, Col: k}) {
				return
			}
			stack = append(stack, frame{Cell{Row: r, Col: k}, 0})
		}
	}, nil
}
