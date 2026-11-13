package stats

import (
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// MeanAbsoluteDeviation returns the average absolute deviation of values about their mean,
// Σ|x − x̄|/n (about the MEAN; for the median absolute deviation about the median use
// MedianAbsoluteDeviation). The sum is compensated. values is read, never modified.
// Errors: ErrEmpty when values is empty; ErrNotFinite when a value is NaN or ±Inf.
//
// ID:      NUM-STATS-011
// Time:    O(n)   n = len(values); two passes (mean, then deviations)
// Space:   O(1) extra, no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.6 "Measures of Scale", https://www.itl.nist.gov/div898/handbook/eda/section3/eda356.htm (read 2026-11-12)
// Tags:    mean-absolute-deviation, average-absolute-deviation, aad, dispersion, scale, descriptive
// Units:   values in any one unit; result in that unit
// Example: NumAcc1 (10000001, 10000003, 10000002) → deviations 1, 1, 0 → 2/3 (hand-checked)
func MeanAbsoluteDeviation(values []float64) (float64, error) {
	if err := checkSample(values, 1); err != nil {
		return 0, err
	}
	mean := meanOf(values)
	n := float64(len(values))
	if s := absoluteDeviationSum(values, mean, 1); reduce.AllFinite(s) {
		return s / n, nil
	}
	// Some |x − x̄| or their sum overflowed: scale by 2^-k with 2^k > 2n, so the sum (at most n times a
	// deviation below 2·MaxFloat64) stays finite. The result is at most the half-range, so it is finite.
	k := bitLen(len(values)) + 1
	s := absoluteDeviationSum(values, mean, math.Ldexp(1, -k))
	return math.Ldexp(s/n, k), nil
}

// absoluteDeviationSum returns Σ|x·scale − mean·scale|, compensated.
func absoluteDeviationSum(values []float64, mean, scale float64) float64 {
	m := mean * scale
	var s compensatedSum
	for _, v := range values {
		s.add(math.Abs(v*scale - m))
	}
	return s.total()
}
