package prob

import "math"

// GumbelCDF returns P(X ≤ x) = exp(−e^(−z)), z = (x − μ)/β, for the Gumbel (type I extreme value,
// maximum) distribution with location μ and scale β.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0.
//
// ID:      NUM-PROB-104
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 22 "Extreme Value Distributions", §1 (type 1, eq. 22.1)
// Tags:    probability, distribution, gumbel, extreme-value, cdf, continuous, return-period
// Units:   x, location and scale: one unit; result: a probability in [0, 1]
// Example: μ = 0, β = 1, x = 0 → e⁻¹ ≈ 0.36787944117144233
func GumbelCDF(x, location, scale float64) (float64, error) {
	if err := checkLocationScale(location, scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	return math.Exp(-math.Exp(-(x - location) / scale)), nil
}
