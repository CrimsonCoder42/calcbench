package prob

import "math"

// PowerFunctionQuantile returns x = u^(1/a), the inverse of PowerFunctionCDF; u = 0 returns 0 and u = 1
// returns 1.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when shape ≤ 0.
//
// ID:      NUM-PROB-129
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 25 "Beta Distributions", §1 (the power-function distribution, inverted)
// Tags:    probability, distribution, power-function, beta, quantile, inverse-cdf
// Units:   u and result dimensionless
// Example: a = 3, u = 0.125 → 0.5
func PowerFunctionQuantile(u, shape float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkPositive("shape", shape); err != nil {
		return 0, err
	}
	return math.Pow(u, 1/shape), nil
}
