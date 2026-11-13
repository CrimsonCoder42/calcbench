package stats

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/search"
)

// FrequencyDistribution2D returns the bivariate histogram counts of the pairs (x[t], y[t]): entry
// (i, j), at index i·by + j (row-major; bx = len(xEdges) − 1 rows, by = len(yEdges) − 1 columns), is the
// number of pairs with x in x-class i and y in y-class j. Classes in each direction follow
// FrequencyDistribution (left-closed, the last one closed on both sides). A pair outside the edges is
// an error, never silently dropped. All inputs are read, never modified.
// Errors: ErrEmpty when x is empty; ErrLengthMismatch when len(y) ≠ len(x); ErrTooFew when an edge
// slice has fewer than 2 elements; ErrNotFinite when a value or an edge is NaN or ±Inf; ErrNotAscending
// when an edge slice is not strictly ascending; ErrOutOfRange when a value is outside its edges.
//
// ID:      NUM-STATS-045
// Time:    O(n·(log bx + log by))   n = len(x); two binary searches per pair
// Space:   O(bx·by) for the result (one allocation); O(1) extra
// Source:  Scott, Multivariate Density Estimation: Theory, Practice, and Visualization, 2nd ed., Wiley, 2015, §3.4 (the multivariate histogram on a rectangular mesh; bivariate bin counts ν_kl)
// Tags:    histogram, histogram-2d, bivariate, frequency-distribution, bins, counts, descriptive
// Units:   x and xEdges in one unit, y and yEdges in one unit; counts dimensionless
// Example: pairs (0.5, 0.5), (1.5, 0.5), (1.5, 1.5), (2, 2) with both edges 0, 1, 2 → [1 0; 1 2] (hand-checked)
func FrequencyDistribution2D(x, y, xEdges, yEdges []float64) ([]int, error) {
	if len(x) == 0 {
		return nil, fmt.Errorf("%w: x", ErrEmpty)
	}
	if len(y) != len(x) {
		return nil, fmt.Errorf("%w: len(y) = %d, len(x) = %d", ErrLengthMismatch, len(y), len(x))
	}
	if err := checkFinite("x", x); err != nil {
		return nil, err
	}
	if err := checkFinite("y", y); err != nil {
		return nil, err
	}
	if err := checkEdges(xEdges); err != nil {
		return nil, fmt.Errorf("xEdges: %w", err)
	}
	if err := checkEdges(yEdges); err != nil {
		return nil, fmt.Errorf("yEdges: %w", err)
	}
	bx, by := len(xEdges)-1, len(yEdges)-1
	counts := make([]int, bx*by)
	for t := range x {
		if x[t] < xEdges[0] || x[t] > xEdges[bx] {
			return nil, fmt.Errorf("%w: x[%d] = %v, want %v..%v", ErrOutOfRange, t, x[t], xEdges[0], xEdges[bx])
		}
		if y[t] < yEdges[0] || y[t] > yEdges[by] {
			return nil, fmt.Errorf("%w: y[%d] = %v, want %v..%v", ErrOutOfRange, t, y[t], yEdges[0], yEdges[by])
		}
		i := min(search.UpperBound(xEdges, x[t])-1, bx-1)
		j := min(search.UpperBound(yEdges, y[t])-1, by-1)
		counts[i*by+j]++
	}
	return counts, nil
}
