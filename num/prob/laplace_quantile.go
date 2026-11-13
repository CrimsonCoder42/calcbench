package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LaplaceQuantile returns the inverse of LaplaceCDF: μ + b·ln(2u) for u ≤ ½ and μ − b·ln(2(1 − u)) for
// u > ½ (1 − u exact there), for 0 < u < 1.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when u = 0 or u = 1 (infinite quantile) or scale ≤ 0; ErrOverflow when the result
// exceeds the float64 range.
//
// ID:      NUM-PROB-109
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 24 "Laplace (Double Exponential) Distributions", §1 (definition, inverted)
// Tags:    probability, distribution, laplace, double-exponential, quantile, inverse-cdf
// Units:   u: a probability; location, scale and result: one unit
// Example: μ = 0, b = 1, u = 0.75 → ln 2 ≈ 0.6931471805599453
func LaplaceQuantile(u, location, scale float64) (float64, error) {
	if err := checkQuantileOpen(u); err != nil {
		return 0, err
	}
	if err := checkLocationScale(location, scale); err != nil {
		return 0, err
	}
	var x float64
	if u <= 0.5 {
		x = location + scale*math.Log(2*u)
	} else {
		x = location - scale*math.Log(2*(1-u))
	}
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: quantile for u = %v, location = %v, scale = %v", ErrOverflow, u, location, scale)
	}
	return x, nil
}
