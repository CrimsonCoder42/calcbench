package graph

// SurroundedRegionsInPlace MODIFIES grid: every cell holding region that cannot reach the grid's border
// through cells holding region (moving to the 4 or 8 neighbours that nb names) is set to fill; region
// cells joined to a border cell stay. It returns how many cells it changed (0 when region == fill).
// Values are compared with ==.
// Errors (grid unchanged): ErrNeighborhood; ErrNotRectangular; ErrTooLarge.
//
// ID:      ALG-GRAPH-026
// Time:    O(R·C)   R = rows, C = columns
// Space:   O(R·C) extra (a visited bitmap and the BFS queue)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (BFS; seeding the queue with every border cell is a BFS from one virtual source joined to them)
// Source:  Smith, "Tint fill", Proceedings of SIGGRAPH '79, 1979, pp. 276–283, doi:10.1145/800249.807456 (connected regions of one value)
// Tags:    grid, surrounded-regions, multi-source-bfs, flood-fill, in-place, enclosed
func SurroundedRegionsInPlace[T comparable](grid [][]T, region, fill T, nb Neighborhood) (int, error) {
	steps, err := nb.offsets()
	if err != nil {
		return 0, err
	}
	rows, cols, err := gridShape(grid)
	if err != nil {
		return 0, err
	}
	if region == fill || rows == 0 || cols == 0 {
		return 0, nil
	}
	seen := make([]bool, rows*cols)
	var queue []Cell
	seed := func(r, k int) {
		if !seen[r*cols+k] && grid[r][k] == region {
			seen[r*cols+k] = true
			queue = append(queue, Cell{Row: r, Col: k})
		}
	}
	for k := range cols {
		seed(0, k)
		seed(rows-1, k)
	}
	for r := range rows {
		seed(r, 0)
		seed(r, cols-1)
	}
	for head := 0; head < len(queue); head++ {
		c := queue[head]
		for _, d := range steps {
			r, k := c.Row+d.Row, c.Col+d.Col
			if r < 0 || r >= rows || k < 0 || k >= cols {
				continue
			}
			seed(r, k)
		}
	}
	changed := 0
	for r, row := range grid {
		for k, v := range row {
			if v == region && !seen[r*cols+k] {
				row[k] = fill
				changed++
			}
		}
	}
	return changed, nil
}
