package prob

import (
	"fmt"
	"math"
)

// NegativeBinomialFailuresQuantile returns the smallest k ≥ 0 with P(Y ≤ k) ≥ u, Y the number of FAILURES
// before the r-th success. It brackets k by doubling from the mean r(1 − p)/p and bisects on
// NegativeBinomialFailuresCDF (NUM-PROB-100), so the result agrees with the CDF as this package computes
// it. u = 0 or p = 1 returns 0.
// Errors: ErrNotFinite when u or p is NaN or ±Inf; ErrProbability when u or p is outside [0, 1];
// ErrOutOfRange when r < 1, p = 0, u = 1 with p < 1 (no finite quantile), or as the CDF; ErrOverflow when
// k exceeds the int range.
//
// ID:      NUM-PROB-101
// Time:    O(log(k + 2) · C)   k = the result, C = the cost of one NegativeBinomialFailuresCDF
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 5 "Negative Binomial Distribution", §5.1 (definition)
// Tags:    probability, distribution, negative-binomial, quantile, inverse-cdf, discrete, failures
// Units:   u, p: probabilities; r: a count of successes; result: a count of failures
// Example: r = 2, p = 0.5: u = 0.6 → 2 (F(1) = 0.5 < 0.6 ≤ F(2) = 0.6875)
func NegativeBinomialFailuresQuantile(u float64, r int, p float64) (int, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkNegativeBinomial(r, p); err != nil {
		return 0, err
	}
	if u == 0 || p == 1 {
		return 0, nil
	}
	if u == 1 {
		return 0, fmt.Errorf("%w: u = 1 with p < 1, quantile infinite", ErrOutOfRange)
	}
	cdf := func(k int) (float64, error) { return NegativeBinomialFailuresCDF(k, r, p) }
	lo, hi := -1, max(1, int(math.Min(math.Ceil(float64(r)*(1-p)/p), 1<<60)))
	for {
		c, err := cdf(hi)
		if err != nil {
			return 0, err
		}
		if c >= u {
			break
		}
		lo = hi
		if hi > (math.MaxInt-r)/2 {
			return 0, fmt.Errorf("%w: quantile for u = %v, r = %d, p = %v", ErrOverflow, u, r, p)
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
