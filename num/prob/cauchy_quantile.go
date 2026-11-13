package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// CauchyQuantile returns the u-quantile x₀ + γ·tan(π(u − ½)) of the Cauchy distribution with location x₀
// and scale γ, computed as x₀ − γ/tan(πu) for u ≤ ½ and x₀ + γ/tan(π(1 − u)) above, so a small tail
// probability keeps its accuracy.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when u is 0 or 1 (an infinite quantile) or scale ≤ 0; ErrOverflow when the quantile
// exceeds the float64 range.
//
// ID:      NUM-PROB-093
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 16 "Cauchy Distribution", §1 (eq. 16.1), inverted in closed form
// Tags:    probability, distribution, cauchy, lorentz, quantile, inverse-cdf, ppf
// Units:   u a probability in (0, 1); location, scale and result: one unit
// Example: x₀ = 0, γ = 1, u = 0.75 → 1; u = 0.5 → 0
func CauchyQuantile(u, location, scale float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkCauchy(location, scale); err != nil {
		return 0, err
	}
	if u == 0 || u == 1 {
		return 0, fmt.Errorf("%w: u = %v, quantile infinite", ErrOutOfRange, u)
	}
	if u == 0.5 {
		return location, nil
	}
	var x float64
	if u < 0.5 {
		x = location - scale/math.Tan(math.Pi*u)
	} else {
		x = location + scale/math.Tan(math.Pi*(1-u))
	}
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: quantile for u = %v, scale = %v", ErrOverflow, u, scale)
	}
	return x, nil
}
