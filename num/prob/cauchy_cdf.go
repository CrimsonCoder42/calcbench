package prob

import "math"

// CauchyCDF returns P(X ≤ x) = ½ + arctan(z)/π, z = (x − x₀)/γ, for the Cauchy distribution with location x₀
// and scale γ. It is computed as atan2(1, −z)/π, which equals ½ + arctan(z)/π and keeps the far lower
// tail's relative accuracy (no cancellation against ½).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0.
//
// ID:      NUM-PROB-092
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 16 "Cauchy Distribution", §1 (eq. 16.1)
// Tags:    probability, distribution, cauchy, lorentz, cdf, continuous
// Units:   x, location and scale: one unit; result: a probability in [0, 1]
// Example: x₀ = 0, γ = 1, x = 1 → 0.75; x = −1 → 0.25
func CauchyCDF(x, location, scale float64) (float64, error) {
	if err := checkCauchy(location, scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	z := x/scale - location/scale
	return math.Atan2(1, -z) / math.Pi, nil
}
