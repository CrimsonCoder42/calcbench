package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ExponentialQuantile returns x = −ln(1 − u)/λ, the u-quantile of the exponential distribution with
// rate λ, computed as −log1p(−u)/λ, which is accurate for a small u.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when u = 1 (an infinite quantile) or rate ≤ 0; ErrOverflow when the quantile exceeds the
// float64 range.
//
// ID:      NUM-PROB-060
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.6.7 "Exponential Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3667.htm (read 2026-11-12) (PPF G(p) = −β ln(1 − p), β = 1/λ)
// Tags:    probability, distribution, exponential, quantile, inverse-cdf, ppf
// Units:   u: a probability in [0, 1); rate: per a time or length unit; result: that unit
// Example: λ = 2, u = 0.5 → ln 2/2 ≈ 0.34657359027997264
func ExponentialQuantile(u, rate float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if u == 1 {
		return 0, fmt.Errorf("%w: u = 1, quantile infinite", ErrOutOfRange)
	}
	if err := checkPositive("rate", rate); err != nil {
		return 0, err
	}
	x := -math.Log1p(-u) / rate
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: quantile for rate = %v", ErrOverflow, rate)
	}
	return x, nil
}
