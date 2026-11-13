package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// GeometricTrialsMeanVariance returns the mean 1/p and the variance (1 − p)/p² of X, the number of
// TRIALS up to and including the first success (D-013). The failures count X − 1 has mean (1 − p)/p and
// the same variance.
// Errors: as GeometricTrialsPMF; ErrOverflow when the variance exceeds the float64 range (p below about
// 1.5e-154).
//
// ID:      NUM-PROB-031
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §4.8 (other discrete distributions: the geometric random variable)
// Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §3.2 "Discrete Distributions" (geometric distribution)
// Tags:    probability, distribution, geometric, mean, variance, expected-value, trials
// Units:   p: a probability; mean: a count of trials; variance: that count squared
// Example: p = 0.2 → (5, 20)
func GeometricTrialsMeanVariance(p float64) (mean, variance float64, err error) {
	if err := checkGeometric(p); err != nil {
		return 0, 0, err
	}
	mean = 1 / p
	variance = (1 - p) * mean / p
	if !reduce.AllFinite(mean) || !reduce.AllFinite(variance) {
		return 0, 0, fmt.Errorf("%w: variance for p = %v", ErrOverflow, p)
	}
	return mean, variance, nil
}
