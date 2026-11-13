package prob

import "math"

// WeibullCDF returns P(X ≤ x) = 1 − e^(−(x/λ)^k) for the Weibull distribution with shape k and scale λ,
// computed as −expm1(−(x/λ)^k) so a small x keeps its relative accuracy; 0 for x ≤ 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when shape ≤ 0 or scale ≤ 0.
//
// ID:      NUM-PROB-088
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 21 "Weibull Distributions", §1 (eq. 21.3)
// Tags:    probability, distribution, weibull, cdf, continuous, reliability, failure-probability
// Units:   x and scale: one unit; shape dimensionless; result: a probability in [0, 1]
// Example: k = 2, λ = 1, x = 1 → 1 − e⁻¹ ≈ 0.6321205588285577
func WeibullCDF(x, shape, scale float64) (float64, error) {
	if err := checkWeibull(shape, scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	return -math.Expm1(-math.Pow(x/scale, shape)), nil
}
