package stats

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Range returns the largest value of values minus the smallest. values is read, never modified.
// Errors: ErrEmpty when values is empty; ErrNotFinite when a value is NaN or ±Inf; ErrOverflow when the
// difference exceeds the float64 range (for example −MaxFloat64 and +MaxFloat64).
//
// ID:      NUM-STATS-018
// Time:    O(n)   n = len(values); one pass
// Space:   O(1) extra, no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.6 "Measures of Scale", https://www.itl.nist.gov/div898/handbook/eda/section3/eda356.htm (read 2026-11-12)
// Tags:    range, spread, dispersion, scale, descriptive, min, max
// Units:   values in any one unit; result in that unit
// Example: the 12 resistivity values of e-Handbook §7.2.6.2 → 95.1990 − 95.0610 = 0.1380 (hand-checked)
func Range(values []float64) (float64, error) {
	if err := checkSample(values, 1); err != nil {
		return 0, err
	}
	lo, hi := minMax(values)
	r := hi - lo
	if !reduce.AllFinite(r) {
		return 0, fmt.Errorf("%w: range of values", ErrOverflow)
	}
	return r, nil
}

// minMax returns the least and greatest of finite, non-empty values.
func minMax(values []float64) (lo, hi float64) {
	lo, hi = values[0], values[0]
	for _, v := range values[1:] {
		if v < lo {
			lo = v
		} else if v > hi {
			hi = v
		}
	}
	return lo, hi
}
