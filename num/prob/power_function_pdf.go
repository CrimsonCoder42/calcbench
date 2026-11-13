package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PowerFunctionPDF returns the power-function density f(x) = a·x^(a−1) on [0, 1] (the beta distribution
// with parameters a and 1; numpy's "power") and 0 outside. At x = 0 it is 0 for a > 1, 1 for a = 1, and
// unbounded (ErrOutOfRange) for a < 1.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when shape ≤ 0, or at x = 0 with
// shape < 1; ErrOverflow when the density exceeds the float64 range.
//
// ID:      NUM-PROB-127
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 25 "Beta Distributions", §1 (the power-function distribution, beta(a, 1))
// Tags:    probability, distribution, power-function, beta, pdf, continuous
// Units:   x and result dimensionless (x a fraction of the range)
// Example: a = 3, x = 0.5 → 3·0.25 = 0.75
func PowerFunctionPDF(x, shape float64) (float64, error) {
	if err := checkPositive("shape", shape); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	switch {
	case x < 0 || x > 1:
		return 0, nil
	case x == 0:
		switch {
		case shape > 1:
			return 0, nil
		case shape == 1:
			return 1, nil
		default:
			return 0, fmt.Errorf("%w: density unbounded at x = 0 for shape = %v < 1", ErrOutOfRange, shape)
		}
	}
	d := shape * math.Pow(x, shape-1)
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density at x = %v, shape = %v", ErrOverflow, x, shape)
	}
	return d, nil
}
