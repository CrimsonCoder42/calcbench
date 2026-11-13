package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// UniformMeanVariance returns the mean (a + b)/2 and the variance (b − a)²/12 of the continuous uniform
// distribution on [a, b]. The mean is a + (b − a)/2, so a + b cannot overflow. The variance overflows
// once b − a exceeds about 1.2e154; that is the only ErrOverflow (the mean alone would be representable,
// but the pair is not, so no partial result is returned).
// Errors: ErrNotFinite when a or b is NaN or ±Inf; ErrOutOfRange when a ≥ b; ErrOverflow when the
// variance exceeds the float64 range.
//
// ID:      NUM-PROB-046
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.2 "Uniform Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3662.htm (read 2026-11-12)
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.3 (uniform random variable)
// Tags:    probability, distribution, uniform, mean, variance, expected-value, continuous
// Units:   a, b: one unit; mean: that unit; variance: that unit squared
// Example: a = 0, b = 10 → (5, 100/12 ≈ 8.333333333333334)
func UniformMeanVariance(a, b float64) (mean, variance float64, err error) {
	w, err := checkUniform(a, b)
	if err != nil {
		return 0, 0, err
	}
	variance = w * w / 12
	if !reduce.AllFinite(variance) { // also when b − a itself overflows
		return 0, 0, fmt.Errorf("%w: variance (b − a)²/12 for a = %v, b = %v", ErrOverflow, a, b)
	}
	return a + w/2, variance, nil
}
