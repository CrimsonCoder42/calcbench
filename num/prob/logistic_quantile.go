package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LogisticQuantile returns x = μ + s·ln(u/(1 − u)), the inverse of LogisticCDF, formed as
// μ + s·(ln u − log1p(−u)) so neither tail loses accuracy, for 0 < u < 1.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when u = 0 or u = 1 (infinite quantile) or scale ≤ 0; ErrOverflow when the result
// exceeds the float64 range.
//
// ID:      NUM-PROB-113
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 23 "Logistic Distribution", §1 (definition, inverted: the logit)
// Tags:    probability, distribution, logistic, quantile, inverse-cdf, logit
// Units:   u: a probability; location, scale and result: one unit
// Example: μ = 0, s = 1, u = 0.75 → ln 3 ≈ 1.0986122886681098
func LogisticQuantile(u, location, scale float64) (float64, error) {
	if err := checkQuantileOpen(u); err != nil {
		return 0, err
	}
	if err := checkLocationScale(location, scale); err != nil {
		return 0, err
	}
	x := location + scale*(math.Log(u)-math.Log1p(-u))
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: quantile for u = %v, location = %v, scale = %v", ErrOverflow, u, location, scale)
	}
	return x, nil
}
