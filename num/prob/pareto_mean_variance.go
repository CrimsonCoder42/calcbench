package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ParetoMeanVariance returns the mean α·x_m/(α − 1) and the variance α·x_m²/((α − 1)²·(α − 2)) of the
// Pareto (first kind) distribution with scale x_m and shape α. The variance is finite only for α > 2 (the
// mean for α > 1); the pair is returned together, so α ≤ 2 is an error.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0 or shape ≤ 2;
// ErrOverflow when the mean or the variance exceeds the float64 range.
//
// ID:      NUM-PROB-118
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 20 "Pareto Distributions", §3 (moments of the first kind)
// Tags:    probability, distribution, pareto, power-law, mean, variance
// Units:   scale and mean: one unit; variance: that unit squared; shape dimensionless
// Example: x_m = 1, α = 3 → (1.5, 0.75)
func ParetoMeanVariance(scale, shape float64) (mean, variance float64, err error) {
	if err := checkPareto(scale, shape); err != nil {
		return 0, 0, err
	}
	if shape <= 2 {
		return 0, 0, fmt.Errorf("%w: shape = %v, variance infinite for shape ≤ 2", ErrOutOfRange, shape)
	}
	mean = shape * scale / (shape - 1)
	variance = mean / (shape - 1) * scale / (shape - 2)
	if !reduce.AllFinite(mean) || !reduce.AllFinite(variance) {
		return 0, 0, fmt.Errorf("%w: scale = %v, shape = %v", ErrOverflow, scale, shape)
	}
	return mean, variance, nil
}
