package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// GumbelQuantile returns x = μ − β·ln(−ln u), the inverse of GumbelCDF, for 0 < u < 1. Above u = ½ the
// inner −ln u is formed as −log1p(u − 1) (u − 1 exact), so the upper tail keeps its accuracy.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when u = 0 or u = 1 (infinite quantile) or scale ≤ 0; ErrOverflow when the result
// exceeds the float64 range.
//
// ID:      NUM-PROB-105
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 22 "Extreme Value Distributions", §1 (type 1, eq. 22.1, inverted)
// Tags:    probability, distribution, gumbel, extreme-value, quantile, inverse-cdf, return-level
// Units:   u: a probability; location, scale and result: one unit
// Example: μ = 0, β = 1, u = e⁻¹ → 0
func GumbelQuantile(u, location, scale float64) (float64, error) {
	if err := checkQuantileOpen(u); err != nil {
		return 0, err
	}
	if err := checkLocationScale(location, scale); err != nil {
		return 0, err
	}
	var t float64
	if u > 0.5 {
		t = -math.Log1p(u - 1)
	} else {
		t = -math.Log(u)
	}
	x := location - scale*math.Log(t)
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: quantile for u = %v, location = %v, scale = %v", ErrOverflow, u, location, scale)
	}
	return x, nil
}
