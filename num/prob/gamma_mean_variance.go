package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// GammaMeanVariance returns the mean kθ and the variance kθ² of a gamma variable with shape k and scale θ.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when shape ≤ 0 or scale ≤ 0;
// ErrOverflow when the mean or the variance exceeds the float64 range.
//
// ID:      NUM-PROB-072
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 17 "Gamma Distributions", §3 (moments)
// Tags:    probability, distribution, gamma, mean, variance, expected-value
// Units:   scale and mean: one unit; variance: that unit squared; shape dimensionless
// Example: k = 3, θ = 2 → (6, 12)
func GammaMeanVariance(shape, scale float64) (mean, variance float64, err error) {
	if err := checkGamma(shape, scale); err != nil {
		return 0, 0, err
	}
	mean = shape * scale
	variance = mean * scale
	if !reduce.AllFinite(mean) || !reduce.AllFinite(variance) {
		return 0, 0, fmt.Errorf("%w: shape = %v, scale = %v", ErrOverflow, shape, scale)
	}
	return mean, variance, nil
}
