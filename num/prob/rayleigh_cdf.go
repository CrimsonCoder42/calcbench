package prob

import "math"

// RayleighCDF returns P(X ≤ x) = 1 − e^(−x²/(2σ²)) for x ≥ 0 (as −expm1, so a small x keeps its relative
// accuracy) and 0 below, σ the scale.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0.
//
// ID:      NUM-PROB-120
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 18 "Chi-Square Distributions Including Chi and Rayleigh", §10 (Rayleigh distribution)
// Tags:    probability, distribution, rayleigh, cdf, continuous, wind-speed
// Units:   x and scale: one unit; result: a probability in [0, 1]
// Example: σ = 1, x = 1 → 1 − e^(−½) ≈ 0.3934693402873666
func RayleighCDF(x, scale float64) (float64, error) {
	if err := checkPositive("scale", scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	z := x / scale
	return -math.Expm1(-0.5 * z * z), nil
}
