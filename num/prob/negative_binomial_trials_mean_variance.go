package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// NegativeBinomialTrialsMeanVariance returns the mean r/p and the variance r(1 − p)/p² of X, the number
// of TRIALS up to and including the r-th success. The failures count X − r has mean r(1 − p)/p and the
// same variance.
// Errors: as NegativeBinomialTrialsPMF; ErrOverflow when the mean or the variance exceeds the float64
// range.
//
// ID:      NUM-PROB-034
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the negative binomial random variable)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (negative binomial distribution)
// Tags:    probability, distribution, negative-binomial, pascal, mean, variance, expected-value, trials
// Units:   r: a count of successes; p: a probability; mean: a count of trials; variance: that count squared
// Example: r = 2, p = 0.5 → (4, 4)
func NegativeBinomialTrialsMeanVariance(r int, p float64) (mean, variance float64, err error) {
	if err := checkNegativeBinomial(r, p); err != nil {
		return 0, 0, err
	}
	mean = float64(r) / p
	variance = mean * (1 - p) / p
	if !reduce.AllFinite(mean) || !reduce.AllFinite(variance) {
		return 0, 0, fmt.Errorf("%w: r = %d, p = %v", ErrOverflow, r, p)
	}
	return mean, variance, nil
}
