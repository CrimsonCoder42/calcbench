package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// WeibullQuantile returns the u-quantile λ·(−ln(1 − u))^(1/k) of the Weibull distribution with shape k and
// scale λ, computed with log1p so a small u keeps its accuracy.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when u = 1 (an infinite quantile), shape ≤ 0 or scale ≤ 0; ErrOverflow when the quantile
// exceeds the float64 range.
//
// ID:      NUM-PROB-089
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 21 "Weibull Distributions", §1 (eq. 21.3), inverted in closed form
// Tags:    probability, distribution, weibull, quantile, inverse-cdf, ppf, reliability, b10-life
// Units:   u a probability in [0, 1); scale and result: one unit; shape dimensionless
// Example: k = 2, λ = 1, u = 1 − e⁻¹ → 1
func WeibullQuantile(u, shape, scale float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkWeibull(shape, scale); err != nil {
		return 0, err
	}
	if u == 1 {
		return 0, fmt.Errorf("%w: u = 1, quantile infinite", ErrOutOfRange)
	}
	x := scale * math.Pow(-math.Log1p(-u), 1/shape)
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: quantile for u = %v, shape = %v, scale = %v", ErrOverflow, u, shape, scale)
	}
	return x, nil
}
