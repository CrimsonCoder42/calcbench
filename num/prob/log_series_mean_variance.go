package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LogSeriesMeanVariance returns the mean μ = a·p/(1 − p) and variance a·p(1 − a·p)/(1 − p)², a = −1/ln(1 − p),
// of the log-series distribution. ln(1 − p) is math.Log1p(−p).
// Errors: ErrNotFinite when p is NaN or ±Inf; ErrOutOfRange unless 0 < p < 1; ErrOverflow when p is so close
// to 1 that the variance exceeds the float64 range.
//
// ID:      NUM-PROB-137
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kemp, Kotz, Univariate Discrete Distributions, 3rd ed., Wiley, 2005, ch. 7, §7.1 (moments)
// Tags:    probability, distribution, log-series, logarithmic, mean, variance, moments, discrete
// Units:   p: a probability; mean: a count; variance: a count squared
// Example: p = 0.5 → mean 1/ln 2 ≈ 1.4426950408889634, variance (1 − 0.5/ln 2)/(0.5·ln 2) ≈ 0.804021100772319
func LogSeriesMeanVariance(p float64) (mean, variance float64, err error) {
	if err := checkLogSeries(p); err != nil {
		return 0, 0, err
	}
	a := -1 / math.Log1p(-p)
	q := 1 - p
	mean = a * p / q
	variance = a * p * (1 - a*p) / (q * q)
	if !reduce.AllFinite(mean, variance) {
		return 0, 0, fmt.Errorf("%w: log-series moments for p = %v", ErrOverflow, p)
	}
	return mean, variance, nil
}
