package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LaplaceMeanVariance returns the mean μ and the variance 2b² of the Laplace distribution with location μ
// and scale b.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0; ErrOverflow when the
// variance exceeds the float64 range.
//
// ID:      NUM-PROB-110
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 24 "Laplace (Double Exponential) Distributions", §2 (moments)
// Tags:    probability, distribution, laplace, double-exponential, mean, variance
// Units:   location, scale and mean: one unit; variance: that unit squared
// Example: μ = 3, b = 2 → (3, 8)
func LaplaceMeanVariance(location, scale float64) (mean, variance float64, err error) {
	if err := checkLocationScale(location, scale); err != nil {
		return 0, 0, err
	}
	variance = 2 * scale * scale
	if !reduce.AllFinite(variance) {
		return 0, 0, fmt.Errorf("%w: scale = %v", ErrOverflow, scale)
	}
	return location, variance, nil
}
