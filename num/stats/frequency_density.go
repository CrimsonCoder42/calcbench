package stats

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// FrequencyDensity returns the histogram density estimate in each class defined by edges: the count of
// values in class j divided by n·(edges[j+1] − edges[j]), n = len(values), so that the bars' areas sum
// to 1. Classes follow FrequencyDistribution (left-closed, the last one closed on both sides), and a value
// outside the edges is an error, never silently dropped. values and edges are read, never modified.
// Errors: ErrEmpty when values is empty; ErrTooFew when edges has fewer than 2 elements; ErrNotFinite
// when a value or an edge is NaN or ±Inf; ErrNotAscending when edges are not strictly ascending;
// ErrOutOfRange when a value is outside the edges; ErrOverflow when a density is beyond float64.
//
// ID:      NUM-STATS-044
// Time:    O(n log b)   n = len(values), b = len(edges) − 1 classes
// Space:   O(b) for the result; O(b) extra for the counts (two allocations)
// Source:  Scott, Multivariate Density Estimation: Theory, Practice, and Visualization, 2nd ed., Wiley, 2015, §3.1, eq. 3.1 (histogram estimator f̂(x) = ν_k/(n·h_k) for x in bin k)
// Tags:    histogram, density, density-estimation, frequency-distribution, bins, descriptive
// Units:   values and edges in one unit; result in 1/unit
// Example: values 1, 2, 2, 3, 5 with edges 0, 2, 4, 6 → counts 1, 3, 1 → 0.1, 0.3, 0.1 (hand-checked)
func FrequencyDensity(values, edges []float64) ([]float64, error) {
	counts, err := FrequencyDistribution(values, edges)
	if err != nil {
		return nil, err
	}
	n := float64(len(values))
	out := make([]float64, len(counts))
	for j, c := range counts {
		halfWidth := edges[j+1]/2 - edges[j]/2 // cannot overflow
		d := float64(c) / n / 2 / halfWidth
		if !reduce.AllFinite(d) {
			return nil, fmt.Errorf("%w: density of class %d", ErrOverflow, j)
		}
		out[j] = d
	}
	return out, nil
}
