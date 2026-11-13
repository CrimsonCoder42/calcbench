package prob

import "math"

// ParetoCDF returns P(X ≤ x) = 1 − (x_m/x)ᵅ for x ≥ x_m (computed as −expm1(α·ln(x_m/x)), so x just above
// x_m keeps its relative accuracy) and 0 below, x_m the scale and α the shape.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0 or shape ≤ 0.
//
// ID:      NUM-PROB-116
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 20 "Pareto Distributions", §3 (Pareto distribution of the first kind)
// Tags:    probability, distribution, pareto, power-law, cdf, continuous, heavy-tail
// Units:   x and scale: one unit; shape dimensionless; result: a probability in [0, 1]
// Example: x_m = 1, α = 2, x = 2 → 0.75
func ParetoCDF(x, scale, shape float64) (float64, error) {
	if err := checkPareto(scale, shape); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x <= scale {
		return 0, nil
	}
	return -math.Expm1(shape * math.Log(scale/x)), nil
}
