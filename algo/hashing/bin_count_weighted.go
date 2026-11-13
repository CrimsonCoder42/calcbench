package hashing

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Errors for the weighted count.
var (
	// ErrLength: two slices that must have the same length do not.
	ErrLength = errors.New("hashing: slices differ in length")
	// ErrNotFinite: a weight is NaN or ±Inf.
	ErrNotFinite = errors.New("hashing: weight is NaN or Inf")
	// ErrOverflow: every weight is finite but a bin's total exceeds ±MaxFloat64.
	ErrOverflow = errors.New("hashing: bin total overflows float64")
)

// BinCountWeighted returns, for each integer 0…bins-1, the sum of the weights of the values equal to it:
// sums[v] = Σ weights[i] over i with values[i] == v. Each bin is a Neumaier-compensated sum (ALG-RED-002),
// so its error does not grow with the number of values in it. values and weights are read, never modified.
// Errors: ErrOutOfRange when bins < 0, or when a value is negative or ≥ bins; ErrLength when len(values) !=
// len(weights); ErrNotFinite when a weight is NaN or ±Inf; ErrOverflow when a bin total overflows. All
// are wrapped with the offending parameter or index.
//
// ID:      ALG-HASH-015
// Time:    O(n + bins)   n = len(values)
// Space:   O(bins) extra (the compensation terms, one allocation); O(bins) for the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §11.1 (direct-address tables)
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 (compensated sum per bin)
// Tags:    count, histogram, bincount, weighted, direct-address, compensated, sum
func BinCountWeighted(values []int, weights []float64, bins int) ([]float64, error) {
	if bins < 0 {
		return nil, fmt.Errorf("%w: bins = %d, want ≥ 0", ErrOutOfRange, bins)
	}
	if len(values) != len(weights) {
		return nil, fmt.Errorf("%w: len(values) = %d, len(weights) = %d", ErrLength, len(values), len(weights))
	}
	sums := make([]float64, bins)
	comps := make([]float64, bins)
	for i, v := range values {
		if v < 0 || v >= bins {
			return nil, fmt.Errorf("%w: values[%d] = %d, want 0..%d", ErrOutOfRange, i, v, bins-1)
		}
		w := weights[i]
		if !reduce.AllFinite(w) {
			return nil, fmt.Errorf("%w: weights[%d] = %v", ErrNotFinite, i, w)
		}
		sums[v], comps[v] = reduce.NeumaierAdd(sums[v], comps[v], w)
		if !reduce.AllFinite(sums[v]) {
			return nil, fmt.Errorf("%w: bin %d", ErrOverflow, v)
		}
	}
	for v := range sums {
		sums[v] += comps[v]
		if !reduce.AllFinite(sums[v]) {
			return nil, fmt.Errorf("%w: bin %d", ErrOverflow, v)
		}
	}
	return sums, nil
}
