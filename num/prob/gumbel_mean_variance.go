package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// eulerGamma is the Euler–Mascheroni constant γ (DLMF 5.2.3).
const eulerGamma = 0.57721566490153286060651209008240243

// GumbelMeanVariance returns the mean μ + γβ (γ the Euler–Mascheroni constant) and the variance π²β²/6 of
// the Gumbel (type I extreme value, maximum) distribution with location μ and scale β.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when scale ≤ 0; ErrOverflow when the
// mean or the variance exceeds the float64 range.
//
// ID:      NUM-PROB-106
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 22 "Extreme Value Distributions", §2 (moments of the type 1 distribution)
// Source:  NIST DLMF, §5.2.3 (Euler's constant γ)
// Tags:    probability, distribution, gumbel, extreme-value, mean, variance
// Units:   location, scale and mean: one unit; variance: that unit squared
// Example: μ = 0, β = 1 → (γ ≈ 0.5772156649015329, π²/6 ≈ 1.6449340668482264)
func GumbelMeanVariance(location, scale float64) (mean, variance float64, err error) {
	if err := checkLocationScale(location, scale); err != nil {
		return 0, 0, err
	}
	mean = location + eulerGamma*scale
	variance = math.Pi * math.Pi / 6 * scale * scale
	if !reduce.AllFinite(mean) || !reduce.AllFinite(variance) {
		return 0, 0, fmt.Errorf("%w: location = %v, scale = %v", ErrOverflow, location, scale)
	}
	return mean, variance, nil
}
