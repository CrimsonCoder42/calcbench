package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// WeibullMeanVariance returns the mean λ·Γ(1 + 1/k) and the variance λ²·(Γ(1 + 2/k) − Γ(1 + 1/k)²) of a
// Weibull variable with shape k and scale λ.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when shape ≤ 0 or scale ≤ 0;
// ErrOverflow when the mean or the variance exceeds the float64 range (a small shape: Γ(1 + 2/k) overflows
// below k ≈ 0.0118).
//
// ID:      NUM-PROB-090
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 21 "Weibull Distributions", §2 (moments, eq. 21.9)
// Tags:    probability, distribution, weibull, mean, variance, mttf, reliability
// Units:   scale and mean: one unit; variance: that unit squared; shape dimensionless
// Example: k = 1, λ = 2 → (2, 4) (the exponential case); k = 2, λ = 1 → (√π/2, 1 − π/4)
func WeibullMeanVariance(shape, scale float64) (mean, variance float64, err error) {
	if err := checkWeibull(shape, scale); err != nil {
		return 0, 0, err
	}
	g1 := math.Gamma(1 + 1/shape)
	g2 := math.Gamma(1 + 2/shape)
	mean = scale * g1
	variance = scale * scale * (g2 - g1*g1)
	if !reduce.AllFinite(mean) || !reduce.AllFinite(variance) || !reduce.AllFinite(g2) {
		return 0, 0, fmt.Errorf("%w: shape = %v, scale = %v", ErrOverflow, shape, scale)
	}
	return mean, variance, nil
}
