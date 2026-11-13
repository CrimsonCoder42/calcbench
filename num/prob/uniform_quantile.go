package prob

import (
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// UniformQuantile returns x = a + u·(b − a), the u-quantile of the continuous uniform distribution on
// [a, b], with one rounding (math.FMA), clamped to [a, b]. When b − a exceeds the float64 range it is
// 2·(a/2 + u·(b/2 − a/2)), every term finite, so the result is always returned (a = −1e308, b = 1e308,
// u = ½ → 0).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when a ≥ b.
//
// ID:      NUM-PROB-059
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.2 "Uniform Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3662.htm (read 2026-11-12)
// Source:  Ross, A First Course in Probability, 10th ed., Pearson, 2019, §5.3 (uniform random variable)
// Tags:    probability, distribution, uniform, quantile, inverse-cdf, ppf
// Units:   u: a probability in [0, 1]; a, b and result: one unit
// Example: a = 0, b = 10, u = 0.3 → 3
func UniformQuantile(u, a, b float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	w, err := checkUniform(a, b)
	if err != nil {
		return 0, err
	}
	if !reduce.AllFinite(w) {
		return min(b, max(a, 2*math.FMA(u, b/2-a/2, a/2))), nil
	}
	return min(b, max(a, math.FMA(u, w, a))), nil
}
