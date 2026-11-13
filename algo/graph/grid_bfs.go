package graph

import (
	"errors"
	"fmt"
	"iter"

	"github.com/CrimsonCoder42/calcbench/algo/unionfind"
)

// Grid errors.
var (
	// ErrNotRectangular: the rows of a grid have different lengths.
	ErrNotRectangular = errors.New("graph: grid rows differ in length")
	// ErrCellOutOfRange: a cell is outside the grid.
	ErrCellOutOfRange = errors.New("graph: cell out of range")
	// ErrNilPredicate: the passable or isLand function is nil.
	ErrNilPredicate = errors.New("graph: predicate is nil")
	// ErrNeighborhood: the Neighborhood is neither FourNeighbors nor EightNeighbors.
	ErrNeighborhood = errors.New("graph: neighborhood must be FourNeighbors or EightNeighbors")
)

// Cell is a grid position, Row and Col counting from 0 at the top-left. It is unionfind.Cell (one type).
type Cell = unionfind.Cell

// Neighborhood says which cells are adjacent on a grid.
type Neighborhood int

const (
	// FourNeighbors: the cells above, right, below and left (in that order).
	FourNeighbors Neighborhood = 4
	// EightNeighbors: the four plus the diagonals, clockwise from above.
	EightNeighbors Neighborhood = 8
)

// fourOffsets and eightOffsets are the (row, col) steps of FourNeighbors and EightNeighbors, clockwise
// from up (constant tables, never written).
var (
	fourOffsets  = [4]Cell{{Row: -1, Col: 0}, {Row: 0, Col: 1}, {Row: 1, Col: 0}, {Row: 0, Col: -1}}
	eightOffsets = [8]Cell{{Row: -1, Col: 0}, {Row: -1, Col: 1}, {Row: 0, Col: 1}, {Row: 1, Col: 1}, {Row: 1, Col: 0}, {Row: 1, Col: -1}, {Row: 0, Col: -1}, {Row: -1, Col: -1}}
)

// offsets returns the neighbour steps of nb, or ErrNeighborhood.
func (nb Neighborhood) offsets() ([]Cell, error) {
	switch nb {
	case FourNeighbors:
		return fourOffsets[:], nil
	case EightNeighbors:
		return eightOffsets[:], nil
	}
	return nil, fmt.Errorf("%w: nb = %d", ErrNeighborhood, int(nb))
}

// gridShape returns the row and column counts of a rectangular grid, or ErrNotRectangular, or
// ErrTooLarge when rows×cols (one visited byte per cell) exceeds maxAllocBytes.
func gridShape[T any](grid [][]T) (rows, cols int, err error) {
	rows = len(grid)
	if rows == 0 {
		return 0, 0, nil
	}
	cols = len(grid[0])
	for r, row := range grid {
		if len(row) != cols {
			return 0, 0, fmt.Errorf("%w: row %d has %d cells, row 0 has %d", ErrNotRectangular, r, len(row), cols)
		}
	}
	if cols > 0 && uint64(rows) > maxAllocBytes/uint64(cols) {
		return 0, 0, fmt.Errorf("%w: %d×%d cells", ErrTooLarge, rows, cols)
	}
	return rows, cols, nil
}

// gridArgs checks the grid, the start cell and the neighbourhood shared by the grid searches.
func gridArgs[T any](grid [][]T, start Cell, nb Neighborhood) (rows, cols int, steps []Cell, err error) {
	if steps, err = nb.offsets(); err != nil {
		return 0, 0, nil, err
	}
	if rows, cols, err = gridShape(grid); err != nil {
		return 0, 0, nil, err
	}
	if start.Row < 0 || start.Row >= rows || start.Col < 0 || start.Col >= cols {
		return 0, 0, nil, fmt.Errorf("%w: start = (%d, %d), grid %d×%d", ErrCellOutOfRange, start.Row, start.Col, rows, cols)
	}
	return rows, cols, steps, nil
}

// GridBFS returns the cells reachable from start through passable cells, in breadth-first order (start,
// then cells 1 step away, then 2, …; within a step in neighbour order), moving to the 4 or 8 neighbours
// that nb names. A start cell that is not passable gives an empty sequence. Each range over the sequence
// runs a fresh search and reads the grid as it goes; breaking out stops it. grid is never modified.
// Errors (returned before any iteration): ErrNeighborhood; ErrNotRectangular; ErrTooLarge;
// ErrCellOutOfRange when start is outside the grid; ErrNilPredicate when passable is nil.
//
// ID:      ALG-GRAPH-022
// Time:    O(R·C) for a full iteration   R = rows, C = columns (each cell is queued once and checks at most 8 neighbours)
// Space:   O(R·C) extra per iteration (a visited bitmap and the queue)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (BFS, here on the implicit grid graph)
// Source:  The Go Programming Language Specification, "For statements with range clause" (range over function iterators); package iter (iter.Seq)
// Tags:    grid, bfs, breadth-first, traversal, iterator, maze, implicit-graph
func GridBFS[T any](grid [][]T, start Cell, passable func(T) bool, nb Neighborhood) (iter.Seq[Cell], error) {
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
		seen := make([]bool, rows*cols)
		seen[start.Row*cols+start.Col] = true
		queue := []Cell{start}
		for head := 0; head < len(queue); head++ {
			c := queue[head]
			if !yield(c) {
				return
			}
			for _, d := range steps {
				r, k := c.Row+d.Row, c.Col+d.Col
				if r < 0 || r >= rows || k < 0 || k >= cols || seen[r*cols+k] || !passable(grid[r][k]) {
					continue
				}
				seen[r*cols+k] = true
				queue = append(queue, Cell{Row: r, Col: k})
			}
		}
	}, nil
}
