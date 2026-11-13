package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ExponentialMeanVariance returns the mean 1/λ and the variance 1/λ² of an exponential variable with
// rate λ.
// Errors: ErrNotFinite when rate is NaN or ±Inf; ErrOutOfRange when rate ≤ 0; ErrOverflow when the mean
// or the variance exceeds the float64 range (rate below about 1.5e-154).
//
// ID:      NUM-PROB-056
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.7 "Exponential Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3667.htm (read 2026-11-12)
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.5 (exponential random variables)
// Tags:    probability, distribution, exponential, mean, variance, expected-value
// Units:   rate: per a time or length unit; mean: that unit; variance: that unit squared
// Example: λ = 2 → (0.5, 0.25)
func ExponentialMeanVariance(rate float64) (mean, variance float64, err error) {
	if err := checkPositive("rate", rate); err != nil {
		return 0, 0, err
	}
	mean = 1 / rate
	variance = mean / rate
	if !reduce.AllFinite(mean) || !reduce.AllFinite(variance) {
		return 0, 0, fmt.Errorf("%w: rate = %v", ErrOverflow, rate)
	}
	return mean, variance, nil
}
