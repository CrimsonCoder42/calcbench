package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LognormalMeanVariance returns the mean exp(μ + σ²/2) and the variance expm1(σ²)·exp(2μ + σ²) of X
// lognormal with log-mean μ and log-standard-deviation σ. expm1 keeps the variance accurate for a small σ.
// When a factor overflows or underflows on its own, the variance is formed from its logarithm instead.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when sigma ≤ 0; ErrOverflow when the
// mean or the variance exceeds the float64 range.
//
// ID:      NUM-PROB-053
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.9 "Lognormal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3669.htm (read 2026-11-12) (the m = 1 form, scaled by e^μ)
// Tags:    probability, distribution, lognormal, mean, variance, expected-value
// Units:   mu: natural log of the caller's unit; sigma: dimensionless; mean: that unit; variance: that unit squared
// Example: μ = 0, σ = 1 → (√e ≈ 1.6487212707001282, (e − 1)·e ≈ 4.670774270471604)
func LognormalMeanVariance(mu, sigma float64) (mean, variance float64, err error) {
	if err := checkLognormal(mu, sigma); err != nil {
		return 0, 0, err
	}
	s2 := sigma * sigma
	mean = math.Exp(mu + 0.5*s2)
	if !reduce.AllFinite(mean) {
		return 0, 0, fmt.Errorf("%w: mean for mu = %v, sigma = %v", ErrOverflow, mu, sigma)
	}
	em1, scale := math.Expm1(s2), math.Exp(2*mu+s2)
	if reduce.AllFinite(em1) && reduce.AllFinite(scale) && scale > 0 {
		variance = em1 * scale
	} else {
		// ln var = 2μ + σ² + ln(expm1(σ²)), with ln(expm1(s)) = s + log1p(−e^(−s)) for a large s.
		lv := 2*mu + s2
		if reduce.AllFinite(em1) {
			lv += math.Log(em1)
		} else {
			lv += s2 + math.Log1p(-math.Exp(-s2))
		}
		variance = math.Exp(lv)
	}
	if !reduce.AllFinite(variance) {
		return 0, 0, fmt.Errorf("%w: variance for mu = %v, sigma = %v", ErrOverflow, mu, sigma)
	}
	return mean, variance, nil
}
