package stats

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/search"
)

// FrequencyDistribution returns the number of values in each class defined by edges: class j (0-based)
// holds the values v with edges[j] ≤ v < edges[j+1], and the last class is closed on the right, so it
// also holds v = edges[len(edges)-1]. The result has len(edges) − 1 counts. A value outside
// [edges[0], edges[len(edges)-1]] is an error, never silently dropped.
// values and edges are read, never modified.
// Errors: ErrEmpty when values is empty; ErrTooFew when edges has fewer than 2 elements; ErrNotFinite
// when a value or an edge is NaN or ±Inf; ErrNotAscending when edges are not strictly ascending;
// ErrOutOfRange when a value is outside the edges.
//
// ID:      NUM-STATS-008
// Time:    O(n log b)   n = len(values), b = len(edges) − 1 classes; one binary search per value (ALG-SEARCH-002)
// Space:   O(b) for the result (one allocation); O(1) extra
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.3.14 "Histogram", https://www.itl.nist.gov/div898/handbook/eda/section3/histogra.htm (read 2026-11-12)
// Tags:    frequency-distribution, histogram, counts, classes, bins, descriptive
// Units:   values and edges in one unit; counts dimensionless
// Example: values 1, 2, 2, 3, 5 with edges 0, 2, 4, 6 → 1, 3, 1 (hand-checked)
func FrequencyDistribution(values, edges []float64) ([]int, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("%w: values", ErrEmpty)
	}
	if err := checkFinite("values", values); err != nil {
		return nil, err
	}
	if err := checkEdges(edges); err != nil {
		return nil, err
	}
	b := len(edges) - 1
	first, last := edges[0], edges[b]
	counts := make([]int, b)
	for i, v := range values {
		if v < first || v > last {
			return nil, fmt.Errorf("%w: values[%d] = %v, want %v..%v", ErrOutOfRange, i, v, first, last)
		}
		// UpperBound gives the first edge > v; the class is the one before it. v = last edge would
		// fall past the end, so it joins the last class.
		counts[min(search.UpperBound(edges, v)-1, b-1)]++
	}
	return counts, nil
}

// checkEdges returns ErrTooFew for fewer than 2 edges, ErrNotFinite for a NaN or ±Inf edge and
// ErrNotAscending when the edges are not strictly ascending.
func checkEdges(edges []float64) error {
	if len(edges) < 2 {
		return fmt.Errorf("%w: len(edges) = %d, want ≥ 2", ErrTooFew, len(edges))
	}
	if err := checkFinite("edges", edges); err != nil {
		return err
	}
	for i := 1; i < len(edges); i++ {
		if !(edges[i] > edges[i-1]) {
			return fmt.Errorf("%w: edges[%d] = %v ≤ edges[%d] = %v", ErrNotAscending, i, edges[i], i-1, edges[i-1])
		}
	}
	return nil
}
