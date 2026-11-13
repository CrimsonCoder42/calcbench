package graph

// FloodFillInPlace MODIFIES grid: it sets to replacement the start cell and every cell joined to it (through
// the 4 or 8 neighbours that nb names) by cells holding the start cell's value, and returns how many
// cells it changed. When the start cell already holds replacement nothing changes and it returns 0.
// Values are compared with ==, so a floating-point NaN start cell matches no cell, and only the start
// cell changes.
// Errors (grid unchanged): ErrNeighborhood; ErrNotRectangular; ErrTooLarge; ErrCellOutOfRange when start
// is outside the grid.
//
// ID:      ALG-GRAPH-025
// Time:    O(R·C)   R = rows, C = columns (only the filled region and its border are visited)
// Space:   O(R·C) extra worst case (the BFS queue; no visited array: a changed cell no longer matches)
// Source:  Smith, "Tint fill", Proceedings of SIGGRAPH '79, 1979, pp. 276–283, doi:10.1145/800249.807456 (filling a 4- or 8-connected region of one value)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (BFS on the implicit grid graph)
// Tags:    grid, flood-fill, paint-bucket, bfs, in-place, image
func FloodFillInPlace[T comparable](grid [][]T, start Cell, replacement T, nb Neighborhood) (int, error) {
	rows, cols, steps, err := gridArgs(grid, start, nb)
	if err != nil {
		return 0, err
	}
	original := grid[start.Row][start.Col]
	if original == replacement {
		return 0, nil
	}
	grid[start.Row][start.Col] = replacement
	queue := []Cell{start}
	for head := 0; head < len(queue); head++ {
		c := queue[head]
		for _, d := range steps {
			r, k := c.Row+d.Row, c.Col+d.Col
			if r < 0 || r >= rows || k < 0 || k >= cols || grid[r][k] != original {
				continue
			}
			grid[r][k] = replacement
			queue = append(queue, Cell{Row: r, Col: k})
		}
	}
	return len(queue), nil
}
