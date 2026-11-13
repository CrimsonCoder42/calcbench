package prob

import "math"

// LogisticCDF returns P(X ≤ x) = 1/(1 + e^(−z)), z = (x − μ)/s, for the logistic distribution with
// location μ and scale s; for z < 0 it is e^z/(1 + e^z), so the lower tail keeps its relative accuracy.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0.
//
// ID:      NUM-PROB-112
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 23 "Logistic Distribution", §1 (definition)
// Tags:    probability, distribution, logistic, sigmoid, cdf, continuous, elo, win-probability
// Units:   x, location and scale: one unit; result: a probability in [0, 1]
// Example: μ = 0, s = 1, x = ln 3 → 0.75
func LogisticCDF(x, location, scale float64) (float64, error) {
	if err := checkLocationScale(location, scale); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	z := (x - location) / scale
	if z < 0 {
		e := math.Exp(z)
		return e / (1 + e), nil
	}
	return 1 / (1 + math.Exp(-z)), nil
}
