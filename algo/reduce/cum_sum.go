package reduce

import (
	"fmt"
	"math"
)

// CumSum returns the running sums of values: out[i] = values[0] + … + values[i], each one the
// Neumaier-compensated sum of its prefix (the same value Sum returns for values[:i+1]). The result is
// written into dst (see the package comment). An empty input gives an empty result.
// Errors: ErrNotFinite (with the index) when a value is NaN or ±Inf, checked before anything is written;
// ErrOverflow when a partial sum overflows, in which case dst may have been partly written. On error the
// returned slice is nil.
//
// ID:      ALG-RED-005
// Time:    O(n)   n = len(values)
// Space:   O(1) extra (no allocation when cap(dst) >= n, else one); O(n) for the result
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §4.3
// Tags:    cumsum, prefix-sum, running-sum, scan, compensated, neumaier, dst
func CumSum(dst, values []float64) ([]float64, error) {
	if err := checkFinite(values); err != nil {
		return nil, err
	}
	out := resize(dst, len(values))
	var sum, comp float64
	for i, x := range values {
		sum, comp = NeumaierAdd(sum, comp, x)
		total := sum + comp
		if math.IsInf(sum, 0) || math.IsInf(total, 0) {
			return nil, fmt.Errorf("%w: partial sum through values[%d]", ErrOverflow, i)
		}
		out[i] = total
	}
	return out, nil
}

// checkFinite returns ErrNotFinite wrapped with the index of the first NaN or ±Inf in values.
func checkFinite(values []float64) error {
	for i, x := range values {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, x)
		}
	}
	return nil
}

// resize returns dst[:n] when dst has the capacity, otherwise a new slice of length n.
func resize[T any](dst []T, n int) []T {
	if cap(dst) >= n {
		return dst[:n]
	}
	return make([]T, n)
}
