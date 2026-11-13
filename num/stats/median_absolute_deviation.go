package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// MedianAbsoluteDeviation returns the median absolute deviation of values about their median,
// median|x − x̃| with x̃ = Median(values), unscaled (Dataplot's scaled MADN = MAD/0.6745 estimates a normal
// standard deviation; it is left to the caller). For the deviation about the mean use MeanAbsoluteDeviation.
// values is read, never modified: one buffer (a copy, then the deviations) is reordered.
// Errors: ErrEmpty when values is empty; ErrNotFinite when a value is NaN or ±Inf; ErrOverflow when the
// result exceeds the float64 range.
//
// ID:      NUM-STATS-031
// Time:    O(n) worst case   n = len(values); two linear-time selections
// Space:   O(n) extra (one allocation: a buffer holding a copy of values, then the deviations, each selected in place) plus O(log n) stack
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.6 "Measures of Scale", https://www.itl.nist.gov/div898/handbook/eda/section3/eda356.htm (read 2026-11-12)
// Source:  NIST, Dataplot Reference Manual, "MEDIAN ABSOLUTE DEVIATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/mad.htm (read 2026-11-12)
// Tags:    median-absolute-deviation, mad, robust, dispersion, scale, descriptive
// Units:   values in any one unit; result in that unit
// Example: 1, 1, 2, 2, 4, 6, 9 → median 2, absolute deviations 1, 1, 0, 0, 2, 4, 7 → 1 (hand-checked)
func MedianAbsoluteDeviation(values []float64) (float64, error) {
	if err := checkSample(values, 1); err != nil {
		return 0, err
	}
	// One buffer serves both selections: a copy of values for the median, then the deviations.
	dev := make([]float64, len(values))
	copy(dev, values)
	med := medianInPlace(dev)
	scale := 1.0
	for i, v := range values {
		dev[i] = math.Abs(v - med)
		if !reduce.AllFinite(dev[i]) {
			scale = 0.5
		}
	}
	if scale != 1 {
		// A deviation overflowed; halving both terms keeps every one finite.
		for i, v := range values {
			dev[i] = math.Abs(v*scale - med*scale)
		}
	}
	r := medianInPlace(dev) / scale
	if !reduce.AllFinite(r) {
		return 0, fmt.Errorf("%w: median absolute deviation of values", ErrOverflow)
	}
	return r, nil
}
