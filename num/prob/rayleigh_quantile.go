package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// RayleighQuantile returns x = σ·√(−2·ln(1 − u)), the inverse of RayleighCDF, with ln(1 − u) as log1p(−u);
// u = 0 returns 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when u = 1 (infinite quantile) or scale ≤ 0; ErrOverflow when the result exceeds the
// float64 range.
//
// ID:      NUM-PROB-121
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 18 "Chi-Square Distributions Including Chi and Rayleigh", §10 (Rayleigh distribution, inverted)
// Tags:    probability, distribution, rayleigh, quantile, inverse-cdf
// Units:   u: a probability; scale and result: one unit
// Example: σ = 1, u = 1 − e⁻² → 2
func RayleighQuantile(u, scale float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkPositive("scale", scale); err != nil {
		return 0, err
	}
	if u == 1 {
		return 0, fmt.Errorf("%w: u = 1, quantile infinite", ErrOutOfRange)
	}
	x := scale * math.Sqrt(-2*math.Log1p(-u))
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: quantile for u = %v, scale = %v", ErrOverflow, u, scale)
	}
	return x, nil
}
