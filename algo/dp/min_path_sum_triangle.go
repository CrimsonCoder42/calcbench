package dp

import (
	"fmt"
	"slices"
)

// MinPathSumTriangle returns the smallest sum of a path from the apex of a number triangle to its base
// that steps from entry j of a row to entry j or j+1 of the next; every entry on the path counts.
// Row i must have i+1 entries. Entries may be negative. triangle is read, never modified.
// Errors: ErrEmpty when the triangle has no rows; ErrNotTriangle when row i does not have i+1 entries;
// ErrOverflow when a partial path sum the recurrence holds does not fit in int64 (as L28).
//
// ID:      ALG-DP-024
// Time:    O(N)   N = r(r+1)/2 entries, r = rows
// Space:   O(r) extra (a copy of the base row, one allocation)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §24.2 (shortest paths in a directed acyclic graph in reverse topological order: each entry's best path to the base is its value plus the better of its two children's)
// Tags:    dp, triangle, shortest-path, dag, min-path-sum
func MinPathSumTriangle(triangle [][]int64) (int64, error) {
	r := len(triangle)
	if r == 0 {
		return 0, fmt.Errorf("%w: triangle has no rows", ErrEmpty)
	}
	for i, row := range triangle {
		if len(row) != i+1 {
			return 0, fmt.Errorf("%w: row %d has %d entries", ErrNotTriangle, i, len(row))
		}
	}
	best := slices.Clone(triangle[r-1]) // best[j] = least sum from (row i, j) down to the base
	for i := r - 2; i >= 0; i-- {
		for j, v := range triangle[i] {
			s, ok := addInt64(v, min(best[j], best[j+1]))
			if !ok {
				return 0, fmt.Errorf("%w: path sum from row %d, entry %d", ErrOverflow, i, j)
			}
			best[j] = s
		}
	}
	return best[0], nil
}
