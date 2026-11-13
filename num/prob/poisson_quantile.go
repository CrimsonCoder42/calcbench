package prob

import (
	"fmt"
	"math"
)

// PoissonQuantile returns the smallest k ≥ 0 with P(X ≤ k) ≥ u for X ~ Poisson(λ). It brackets k by
// doubling from ⌈λ⌉ and then bisects on PoissonCDF (NUM-PROB-042), so the result agrees with the CDF as this
// package computes it. u = 0 returns 0.
// Errors: ErrNotFinite when u or λ is NaN or ±Inf; ErrProbability when u is outside [0, 1]; ErrOutOfRange
// when λ < 0, when u = 1 with λ > 0 (no finite quantile), or as PoissonCDF; ErrOverflow when k exceeds the
// int range.
//
// ID:      NUM-PROB-086
// Time:    O(log(k + 2) · C)   k = the result, C = the cost of one PoissonCDF, O(√λ + 1)
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 4 "Poisson Distribution", §4.2 (cdf), inverted as the least k with F(k) ≥ u
// Source:  Devroye, Non-Uniform Random Variate Generation, Springer, 1986, §III.2.1 (inversion by search for discrete distributions)
// Tags:    probability, distribution, poisson, quantile, inverse-cdf, discrete, count
// Units:   u: a probability; lambda: the mean count per interval; result: a count
// Example: λ = 1, u = 0.9 → 2 (F(1) = 2/e ≈ 0.7358 < 0.9 ≤ F(2) = 2.5/e ≈ 0.9197)
func PoissonQuantile(u, lambda float64) (int, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkPoisson(lambda); err != nil {
		return 0, err
	}
	if u == 0 || lambda == 0 {
		return 0, nil
	}
	if u == 1 {
		return 0, fmt.Errorf("%w: u = 1 with lambda > 0, quantile infinite", ErrOutOfRange)
	}
	cdf := func(k int) (float64, error) { return PoissonCDF(k, lambda) }
	lo, hi := -1, max(1, int(math.Min(math.Ceil(lambda), 1<<60)))
	for {
		c, err := cdf(hi)
		if err != nil {
			return 0, err
		}
		if c >= u {
			break
		}
		lo = hi
		if hi > math.MaxInt/2 {
			return 0, fmt.Errorf("%w: quantile for u = %v, lambda = %v", ErrOverflow, u, lambda)
		}
		hi *= 2
	}
	for hi-lo > 1 {
		mid := lo + (hi-lo)/2
		c, err := cdf(mid)
		if err != nil {
			return 0, err
		}
		if c >= u {
			hi = mid
		} else {
			lo = mid
		}
	}
	return hi, nil
}
