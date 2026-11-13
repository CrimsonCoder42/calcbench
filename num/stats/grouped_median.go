package stats

// GroupedMedian returns the estimated median of grouped data, L + ((N/2 − F)/f)·c, where the median
// class is the first whose cumulative count reaches N/2, L is its lower edge, c its width, f its count,
// F the cumulative count below it and N = Σ counts. Class j (0-based) runs from edges[j] to edges[j+1].
// It assumes the observations are spread uniformly within each class; it is an estimate, not the data's
// own median. edges and counts are read, never modified.
// Errors: ErrEmpty when counts is empty; ErrLengthMismatch when len(edges) ≠ len(counts) + 1;
// ErrNotFinite when an edge is NaN or ±Inf; ErrNotAscending when edges are not strictly ascending;
// ErrNegative when a count is negative; ErrZeroTotal when every count is 0; ErrOverflow when the count
// total exceeds the int range.
//
// ID:      NUM-STATS-006
// Time:    O(b)   b = len(counts); two passes (total, then cumulative)
// Space:   O(1) extra, no allocation
// Source:  NIST, Dataplot Reference Manual, "GROUPED PERCENTILE", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/groupper.htm (read 2026-11-12)
// Tags:    grouped-data, median, interpolation, frequency, location, central-tendency, descriptive
// Units:   edges in the unit of the data; counts dimensionless; result in the unit of the data
// Example: classes [0,10) [10,20) [20,30) with counts 2, 5, 3 → 10 + (5 − 2)/5·10 = 16 (hand-checked)
func GroupedMedian(edges []float64, counts []int) (float64, error) {
	total, err := checkGrouped(edges, counts)
	if err != nil {
		return 0, err
	}
	below := 0 // F: the cumulative count before class j
	for j, f := range counts {
		// The median class is the first with below + f ≥ N/2, written below + f ≥ N − (below + f) so
		// that nothing overflows. Then below < N/2, so f > 0.
		if cum := below + f; cum >= total-cum {
			t := (float64(total)/2 - float64(below)) / float64(f)
			return lerp(edges[j], edges[j+1], min(max(t, 0), 1)), nil
		}
		below += f
	}
	return edges[len(edges)-1], nil // unreachable: the last cumulative count is N ≥ N/2
}
