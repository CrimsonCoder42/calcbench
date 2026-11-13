package prob

import "math"

// standardNormalSurvival returns Q(z) = 1 − Φ(z) = ½·erfc(z/√2), accurate in the upper tail.
func standardNormalSurvival(z float64) float64 {
	return 0.5 * math.Erfc(z/math.Sqrt2)
}

// NormalSurvival returns P(X > x) = Q(z) = ½·erfc(z/√2), z = (x − μ)/σ, the upper tail of the normal
// distribution with full relative accuracy (1 − Φ(z) gives 0 from z ≈ 8.3 on).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when sd ≤ 0.
//
// ID:      NUM-PROB-049
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.1 "Normal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3661.htm (read 2026-11-12)
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §7.2 eqs. 7.2.1–7.2.2 (erf, erfc) and §7.1 (Φ(z) = ½erfc(−z/√2), Q(z) = ½erfc(z/√2))
// Tags:    probability, distribution, normal, gaussian, survival, upper-tail, continuous, p-value
// Units:   x, mean, sd: one unit; result: a probability in [0, 1]
// Example: z = 1.96 → 0.02499789514822043; z = 10 → 7.62e-24 (1 − Φ(10) gives 0)
func NormalSurvival(x, mean, sd float64) (float64, error) {
	if err := checkNormal(mean, sd); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	return standardNormalSurvival(standardize(x, mean, sd)), nil
}
