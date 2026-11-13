package prob

import "math"

// standardNormalCDF returns Φ(z) = ½·erfc(−z/√2). erfc keeps full relative accuracy in the lower tail,
// where ½(1 + erf) would cancel.
func standardNormalCDF(z float64) float64 {
	return 0.5 * math.Erfc(-z/math.Sqrt2)
}

// NormalCDF returns P(X ≤ x) = Φ((x − μ)/σ) for X normal with mean μ and standard deviation σ, computed
// as ½·erfc(−z/√2), which is accurate in the lower tail. Upper tails come from NormalSurvival, never
// from 1 − NormalCDF.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when sd ≤ 0.
//
// ID:      NUM-PROB-048
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.1 "Normal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3661.htm (read 2026-11-12)
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §7.2 eqs. 7.2.1–7.2.2 (erf, erfc) and §7.1 (Φ(z) = ½erfc(−z/√2), Q(z) = ½erfc(z/√2))
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.4 (normal random variables)
// Tags:    probability, distribution, normal, gaussian, cdf, continuous, z-score, phi
// Units:   x, mean, sd: one unit; result: a probability in [0, 1]
// Example: NIST §1.3.6.7.1: the area from 0 to 1.96 is 0.47500, so Φ(1.96) = 0.97500 (float64 0.9750021048517795)
func NormalCDF(x, mean, sd float64) (float64, error) {
	if err := checkNormal(mean, sd); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	return standardNormalCDF(standardize(x, mean, sd)), nil
}
