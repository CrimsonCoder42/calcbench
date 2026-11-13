package graph

// NumberOfIslands returns the number of islands in grid: groups of cells for which isLand is true, joined
// through the 4 or 8 neighbours that nb names. An empty grid has 0 islands. grid is never modified.
// Errors: ErrNeighborhood; ErrNotRectangular; ErrTooLarge; ErrNilPredicate when isLand is nil.
//
// ID:      ALG-GRAPH-024
// Time:    O(R·C)   R = rows, C = columns (isLand is called at most 9 times per cell)
// Space:   O(R·C) extra (a visited bitmap and one reused BFS queue)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (BFS from each undiscovered vertex counts the connected components of the implicit grid graph)
// Source:  Smith, "Tint fill", Proceedings of SIGGRAPH '79, 1979, pp. 276–283, doi:10.1145/800249.807456 (4- and 8-connected regions of a raster)
// Tags:    grid, islands, connected-components, bfs, flood-fill
func NumberOfIslands[T any](grid [][]T, isLand func(T) bool, nb Neighborhood) (int, error) {
	steps, err := nb.offsets()
	if err != nil {
		return 0, err
	}
	rows, cols, err := gridShape(grid)
	if err != nil {
		return 0, err
	}
	if isLand == nil {
		return 0, ErrNilPredicate
	}
	seen := make([]bool, rows*cols)
	var queue []Cell
	islands := 0
	for r0 := range rows {
		for c0 := range cols {
			if seen[r0*cols+c0] || !isLand(grid[r0][c0]) {
				continue
			}
			islands++
			seen[r0*cols+c0] = true
			queue = append(queue[:0], Cell{Row: r0, Col: c0})
			for head := 0; head < len(queue); head++ {
				c := queue[head]
				for _, d := range steps {
					r, k := c.Row+d.Row, c.Col+d.Col
					if r < 0 || r >= rows || k < 0 || k >= cols || seen[r*cols+k] || !isLand(grid[r][k]) {
						continue
					}
					seen[r*cols+k] = true
					queue = append(queue, Cell{Row: r, Col: k})
				}
			}
		}
	}
	return islands, nil
}
