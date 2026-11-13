package matrix

// GameOfLifeStep returns the next generation of Conway's Game of Life on grid (true = live): a live cell
// with 2 or 3 live neighbours survives, a dead cell with exactly 3 is born, every other cell is dead next;
// neighbours are the 8 surrounding cells. Cells outside the grid are dead (a finite board, no wrap).
// dst is reused when it already has grid's shape, else a new grid is allocated; dst must not share
// memory with grid (pass the previous generation's buffer to alternate two buffers). grid is read, never
// modified. Errors: ErrRagged.
//
// ID:      ALG-MATRIX-005
// Time:    O(mn)   m = rows, n = columns; 8 reads per cell
// Space:   O(1) extra; O(mn) for the result, two allocations unless dst is reused
// Source:  Gardner, "Mathematical Games: The fantastic combinations of John Conway's new solitaire game 'life'", Scientific American 223(4):120–123, 1970, doi:10.1038/scientificamerican1070-120 (the rules of survival, death and birth)
// Tags:    game-of-life, cellular-automaton, grid, simulation
func GameOfLifeStep(dst, grid [][]bool) ([][]bool, error) {
	m, n, err := gridShape(grid)
	if err != nil {
		return dst, err
	}
	out := newGrid(dst, m, n)
	for r := range m {
		for c := range n {
			live := 0
			for dr := -1; dr <= 1; dr++ {
				rr := r + dr
				if rr < 0 || rr >= m {
					continue
				}
				for dc := -1; dc <= 1; dc++ {
					cc := c + dc
					if (dr != 0 || dc != 0) && cc >= 0 && cc < n && grid[rr][cc] {
						live++
					}
				}
			}
			out[r][c] = live == 3 || (live == 2 && grid[r][c])
		}
	}
	return out, nil
}
