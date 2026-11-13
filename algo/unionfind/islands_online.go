package unionfind

import "fmt"

// Cell is a grid position: Row and Col count from 0 at the top-left corner.
// algo/graph uses the same type (graph.Cell is an alias of it).
type Cell struct {
	Row, Col int
}

// IslandsOnline returns, after each cell of additions turns from water to land on a rows×cols grid that
// starts all water, the number of islands: groups of land cells joined through their 4 neighbours (up,
// down, left, right). Adding a cell that is already land leaves the count unchanged.
// The counts are appended to dst[:0] (dst == nil allocates once, len(additions)); additions is read,
// never modified. Only the added cells are stored, so memory does not depend on rows×cols.
// Errors: ErrNegativeSize when rows < 0 or cols < 0; ErrOutOfRange when an addition is outside the grid.
//
// ID:      ALG-UNION-008
// Time:    O(k α(k)) expected   k = len(additions), α = the inverse Ackermann function (expected: one hash lookup per neighbour)
// Space:   O(k) extra (a map and a DisjointSet over the land cells), O(k) for the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.1 (connected components of a graph whose edges arrive dynamically, maintained with disjoint sets)
// Source:  Tarjan, "Efficiency of a good but not linear set union algorithm", Journal of the ACM, 1975, 22(2):215–225, doi:10.1145/321879.321884
// Tags:    union-find, grid, islands, online, dynamic-connectivity
func IslandsOnline(dst []int, rows, cols int, additions []Cell) ([]int, error) {
	if rows < 0 {
		return nil, fmt.Errorf("%w: rows = %d", ErrNegativeSize, rows)
	}
	if cols < 0 {
		return nil, fmt.Errorf("%w: cols = %d", ErrNegativeSize, cols)
	}
	for i, c := range additions {
		if c.Row < 0 || c.Row >= rows || c.Col < 0 || c.Col >= cols {
			return nil, fmt.Errorf("%w: additions[%d] = (%d, %d), grid %d×%d", ErrOutOfRange, i, c.Row, c.Col, rows, cols)
		}
	}
	if dst == nil {
		dst = make([]int, 0, len(additions))
	}
	dst = dst[:0]
	id := make(map[Cell]int, len(additions))
	set := &DisjointSet{
		parent: make([]int, 0, len(additions)),
		rank:   make([]uint8, 0, len(additions)),
	}
	offsets := [4]Cell{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	for _, c := range additions {
		if _, land := id[c]; !land {
			x := set.MakeSet()
			id[c] = x
			for _, d := range offsets {
				// Off-grid neighbours are never keys of id, so no bounds check is needed.
				if y, ok := id[Cell{c.Row + d.Row, c.Col + d.Col}]; ok {
					set.union(x, y)
				}
			}
		}
		dst = append(dst, set.count)
	}
	return dst, nil
}
