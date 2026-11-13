package prob

import "math"

// ExponentialCDF returns P(X ≤ x) = 1 − e^(−λx) for x ≥ 0, and 0 for x < 0, computed as −expm1(−λx),
// which is accurate for a small λx.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when rate ≤ 0.
//
// ID:      NUM-PROB-055
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.7 "Exponential Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3667.htm (read 2026-11-12)
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.5 (exponential random variables)
// Tags:    probability, distribution, exponential, cdf, continuous, waiting-time, decay
// Units:   x: a time or length unit; rate: per that unit; result: a probability in [0, 1]
// Example: λ = 2, x = 1 → 1 − e⁻² ≈ 0.8646647167633873
func ExponentialCDF(x, rate float64) (float64, error) {
	if err := checkPositive("rate", rate); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	return -math.Expm1(-rate * x), nil
}
