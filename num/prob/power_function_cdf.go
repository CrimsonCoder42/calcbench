package prob

import "math"

// PowerFunctionCDF returns P(X ≤ x) = xᵃ on [0, 1], 0 below and 1 above, for the power-function
// distribution with shape a.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when shape ≤ 0.
//
// ID:      NUM-PROB-128
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 25 "Beta Distributions", §1 (the power-function distribution, beta(a, 1))
// Tags:    probability, distribution, power-function, beta, cdf, continuous
// Units:   x and result dimensionless
// Example: a = 3, x = 0.5 → 0.125
func PowerFunctionCDF(x, shape float64) (float64, error) {
	if err := checkPositive("shape", shape); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	switch {
	case x <= 0:
		return 0, nil
	case x >= 1:
		return 1, nil
	}
	return math.Pow(x, shape), nil
}
