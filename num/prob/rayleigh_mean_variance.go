package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// RayleighMeanVariance returns the mean σ·√(π/2) and the variance (4 − π)/2·σ² of the Rayleigh
// distribution with scale σ.
// Errors: ErrNotFinite when scale is NaN or ±Inf; ErrOutOfRange when scale ≤ 0; ErrOverflow when the
// mean or the variance exceeds the float64 range.
//
// ID:      NUM-PROB-122
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 18 "Chi-Square Distributions Including Chi and Rayleigh", §10 (Rayleigh distribution, moments)
// Tags:    probability, distribution, rayleigh, mean, variance
// Units:   scale and mean: one unit; variance: that unit squared
// Example: σ = 1 → (√(π/2) ≈ 1.2533141373155003, (4 − π)/2 ≈ 0.42920367320510344)
func RayleighMeanVariance(scale float64) (mean, variance float64, err error) {
	if err := checkPositive("scale", scale); err != nil {
		return 0, 0, err
	}
	mean = scale * math.Sqrt(math.Pi/2)
	variance = (4 - math.Pi) / 2 * scale * scale
	if !reduce.AllFinite(mean) || !reduce.AllFinite(variance) {
		return 0, 0, fmt.Errorf("%w: scale = %v", ErrOverflow, scale)
	}
	return mean, variance, nil
}
