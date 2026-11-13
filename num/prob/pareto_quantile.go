package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ParetoQuantile returns x = x_m·(1 − u)^(−1/α), the inverse of ParetoCDF, formed as x_m·e^(−log1p(−u)/α);
// u = 0 returns x_m.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when u = 1 (infinite quantile), scale ≤ 0 or shape ≤ 0; ErrOverflow when the result
// exceeds the float64 range.
//
// ID:      NUM-PROB-117
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 20 "Pareto Distributions", §3 (Pareto distribution of the first kind, inverted)
// Tags:    probability, distribution, pareto, power-law, quantile, inverse-cdf
// Units:   u: a probability; scale and result: one unit; shape dimensionless
// Example: x_m = 1, α = 2, u = 0.75 → 2
func ParetoQuantile(u, scale, shape float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkPareto(scale, shape); err != nil {
		return 0, err
	}
	if u == 1 {
		return 0, fmt.Errorf("%w: u = 1, quantile infinite", ErrOutOfRange)
	}
	x := scale * math.Exp(-math.Log1p(-u)/shape)
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: quantile for u = %v, scale = %v, shape = %v", ErrOverflow, u, scale, shape)
	}
	return x, nil
}
