package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LognormalQuantile returns x = exp(μ + σ·Φ⁻¹(u)), the u-quantile of the lognormal distribution with
// log-mean μ and log-standard-deviation σ; Φ⁻¹ is NormalQuantile's PPND16. A quantile below the smallest
// float64 returns 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when u is outside (0, 1) or
// sigma ≤ 0; ErrOverflow when the quantile exceeds the float64 range.
//
// ID:      NUM-PROB-062
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.9 "Lognormal Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3669.htm (read 2026-11-12) (PPF exp(σΦ⁻¹(p)), scaled by e^μ)
// Source:  Wichura, "Algorithm AS 241: The Percentage Points of the Normal Distribution", Journal of the Royal Statistical Society, Series C (Applied Statistics), 1988, 37(3):477–484, doi:10.2307/2347330 (PPND16)
// Tags:    probability, distribution, lognormal, quantile, inverse-cdf, ppf
// Units:   u: a probability in (0, 1); mu: natural log of the caller's unit; sigma: dimensionless; result: that unit
// Example: μ = 0, σ = 1, u = 0.5 → 1 (the median e^μ)
func LognormalQuantile(u, mu, sigma float64) (float64, error) {
	if err := checkOpenUnit(u); err != nil {
		return 0, err
	}
	if err := checkLognormal(mu, sigma); err != nil {
		return 0, err
	}
	x := math.Exp(mu + sigma*standardNormalQuantile(u))
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: quantile for mu = %v, sigma = %v", ErrOverflow, mu, sigma)
	}
	return x, nil
}
