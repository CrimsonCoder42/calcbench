package stats

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/search"
)

// Digitize returns, for each value, the index of the class it falls in among the half-open intervals cut
// by edges: result[i] = the number of edges ≤ values[i], so 0 means below edges[0], j means
// edges[j−1] ≤ v < edges[j], and len(edges) means v ≥ the last edge. Unlike FrequencyDistribution,
// values outside the edges are allowed (that is what 0 and len(edges) mark). values and edges are read,
// never modified.
// Errors: ErrTooFew when edges is empty; ErrNotFinite when a value or an edge is NaN or ±Inf;
// ErrNotAscending when edges are not strictly ascending. Empty values give an empty result.
//
// ID:      NUM-STATS-046
// Time:    O(n log b)   n = len(values), b = len(edges); one binary search per value (ALG-SEARCH upper bound)
// Space:   O(n) for the result (one allocation); O(1) extra
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.1, Algorithm B (binary search for the insertion point)
// Tags:    histogram, bins, digitize, bin-index, classify, binary-search, descriptive
// Units:   values and edges in one unit; indices dimensionless
// Example: values −1, 0, 0.5, 1, 3 with edges 0, 1, 2 → 0, 1, 1, 2, 3 (hand-checked)
func Digitize(values, edges []float64) ([]int, error) {
	if len(edges) == 0 {
		return nil, fmt.Errorf("%w: edges is empty, want ≥ 1", ErrTooFew)
	}
	if err := checkFinite("values", values); err != nil {
		return nil, err
	}
	if err := checkFinite("edges", edges); err != nil {
		return nil, err
	}
	for i := 1; i < len(edges); i++ {
		if !(edges[i] > edges[i-1]) {
			return nil, fmt.Errorf("%w: edges[%d] = %v ≤ edges[%d] = %v", ErrNotAscending, i, edges[i], i-1, edges[i-1])
		}
	}
	out := make([]int, len(values))
	for i, v := range values {
		out[i] = search.UpperBound(edges, v)
	}
	return out, nil
}
