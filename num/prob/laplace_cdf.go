package prob

import "math"

// LaplaceCDF returns P(X ≤ x) for the Laplace distribution with location μ and scale b: ½e^z for
// z = (x − μ)/b < 0 and 1 − ½e^(−z) for z ≥ 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0.
//
// ID:      NUM-PROB-108
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 24 "Laplace (Double Exponential) Distributions", §1 (definition)
// Tags:    probability, distribution, laplace, double-exponential, cdf, continuous
// Units:   x, location and scale: one unit; result: a probability in [0, 1]
// Example: μ = 0, b = 1, x = 1 → 1 − e⁻¹/2 ≈ 0.8160602794142788
func LaplaceCDF(x, location, scale float64) (float64, error) {
	if err := checkLocationScale(location, scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	z := (x - location) / scale
	if z < 0 {
		return 0.5 * math.Exp(z), nil
	}
	return 1 - 0.5*math.Exp(-z), nil
}
