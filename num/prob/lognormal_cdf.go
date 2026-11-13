package prob

import "math"

// LognormalCDF returns P(X ≤ x) = Φ((ln x − μ)/σ) for x > 0, and 0 for x ≤ 0, X lognormal with log-mean μ
// and log-standard-deviation σ.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when sigma ≤ 0.
//
// ID:      NUM-PROB-052
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.9 "Lognormal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3669.htm (read 2026-11-12)
// Tags:    probability, distribution, lognormal, cdf, continuous
// Units:   x: the caller's unit; mu: natural log of that unit; sigma: dimensionless; result: a probability in [0, 1]
// Example: μ = 0, σ = 1, x = 1 → 0.5
func LognormalCDF(x, mu, sigma float64) (float64, error) {
	if err := checkLognormal(mu, sigma); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	return standardNormalCDF(standardize(math.Log(x), mu, sigma)), nil
}
