package reduce

import (
	"errors"
	"fmt"
	"math"
)

// Errors returned by this package. They are wrapped with the offending parameter or index; test with
// errors.Is.
var (
	// ErrEmpty: the input has no elements and the function has no value to return for it.
	ErrEmpty = errors.New("reduce: empty input")
	// ErrNotFinite: an input is NaN or ±Inf.
	ErrNotFinite = errors.New("reduce: value is NaN or Inf")
	// ErrOverflow: every input is finite but the result, or a partial result, exceeds ±MaxFloat64.
	ErrOverflow = errors.New("reduce: result overflows float64")
	// ErrOutOfRange: an integer parameter is outside the values the function accepts.
	ErrOutOfRange = errors.New("reduce: parameter out of range")
)

// Sum returns the sum of values by Neumaier's compensated summation: a running correction term collects
// the low-order bits each addition loses, and is added back once at the end. The error bound is
// 2u·Σ|xᵢ| + O(nu²)·Σ|xᵢ| (u the unit roundoff), independent of n to first order, against (n-1)u·Σ|xᵢ|
// for a plain loop. An empty slice sums to 0.
// Errors: ErrNotFinite (with the index) when a value is NaN or ±Inf; ErrOverflow when a partial sum
// overflows, even if later values would bring it back.
//
// ID:      ALG-RED-001
// Time:    O(n)   n = len(values)
// Space:   O(1) (no allocation)
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §4.3 (compensated summation)
// Tags:    sum, summation, compensated, neumaier, kahan, reduction, numerically-stable, no-alloc
func Sum(values []float64) (float64, error) {
	var sum, comp float64
	for i, x := range values {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return 0, fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, x)
		}
		sum, comp = NeumaierAdd(sum, comp, x)
		if math.IsInf(sum, 0) {
			return 0, fmt.Errorf("%w: partial sum through values[%d]", ErrOverflow, i)
		}
	}
	total := sum + comp
	if math.IsInf(total, 0) {
		return 0, fmt.Errorf("%w: sum", ErrOverflow)
	}
	return total, nil
}
