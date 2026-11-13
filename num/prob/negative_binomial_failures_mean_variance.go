package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// NegativeBinomialFailuresMeanVariance returns the mean r(1 − p)/p and the variance r(1 − p)/p² of Y, the
// number of FAILURES before the r-th success (the trials count Y + r has mean r/p and the same variance).
// Errors: as NegativeBinomialFailuresPMF; ErrOverflow when the mean or the variance exceeds the float64
// range.
//
// ID:      NUM-PROB-102
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 5 "Negative Binomial Distribution", §5.4 (moments)
// Tags:    probability, distribution, negative-binomial, mean, variance, failures, overdispersion
// Units:   r: a count of successes; p: a probability; mean: a count of failures; variance: that count squared
// Example: r = 2, p = 0.5 → (2, 4)
func NegativeBinomialFailuresMeanVariance(r int, p float64) (mean, variance float64, err error) {
	if err := checkNegativeBinomial(r, p); err != nil {
		return 0, 0, err
	}
	mean = float64(r) * (1 - p) / p
	variance = mean / p
	if !reduce.AllFinite(mean) || !reduce.AllFinite(variance) {
		return 0, 0, fmt.Errorf("%w: r = %d, p = %v", ErrOverflow, r, p)
	}
	return mean, variance, nil
}
