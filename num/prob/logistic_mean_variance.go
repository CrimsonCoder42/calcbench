package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LogisticMeanVariance returns the mean μ and the variance π²s²/3 of the logistic distribution with
// location μ and scale s.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0; ErrOverflow when the
// variance exceeds the float64 range.
//
// ID:      NUM-PROB-114
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 23 "Logistic Distribution", §2 (moments)
// Tags:    probability, distribution, logistic, mean, variance
// Units:   location, scale and mean: one unit; variance: that unit squared
// Example: μ = 0, s = 1 → (0, π²/3 ≈ 3.289868133696453)
func LogisticMeanVariance(location, scale float64) (mean, variance float64, err error) {
	if err := checkLocationScale(location, scale); err != nil {
		return 0, 0, err
	}
	variance = math.Pi * math.Pi / 3 * scale * scale
	if !reduce.AllFinite(variance) {
		return 0, 0, fmt.Errorf("%w: scale = %v", ErrOverflow, scale)
	}
	return location, variance, nil
}
