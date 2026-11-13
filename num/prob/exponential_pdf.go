package prob

import "math"

// ExponentialPDF returns the exponential density f(x) = λ·e^(−λx) for x ≥ 0, and 0 for x < 0, with λ the
// RATE. NIST writes the scale β = 1/λ; a caller holding β passes 1/β.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when rate ≤ 0.
//
// ID:      NUM-PROB-054
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.7 "Exponential Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3667.htm (read 2026-11-12)
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.5 (exponential random variables)
// Tags:    probability, distribution, exponential, pdf, continuous, waiting-time, decay
// Units:   x: a time or length unit; rate: per that unit; result: probability per that unit
// Example: λ = 2, x = 1 → 2e⁻² ≈ 0.2706705664732254
func ExponentialPDF(x, rate float64) (float64, error) {
	if err := checkPositive("rate", rate); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x < 0 {
		return 0, nil
	}
	return rate * math.Exp(-rate*x), nil // rate ≤ MaxFloat64 and exp ≤ 1, so no overflow
}
