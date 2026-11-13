package stats

import "github.com/CrimsonCoder42/calcbench/algo/selection"

// Median returns the median of values: the middle value of the sorted values when n is odd, the mean of
// the two middle values when n is even. It selects in linear time without sorting (ALG-SELECT-002) and
// averages without overflow. values is read, never modified: a copy is reordered.
// Errors: ErrEmpty when values is empty; ErrNotFinite when a value is NaN or ±Inf.
//
// ID:      NUM-STATS-012
// Time:    O(n) worst case   n = len(values)
// Space:   O(n) extra (one allocation, the copy) plus O(log n) stack
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm (read 2026-11-12)
// Tags:    median, location, central-tendency, robust, order-statistic, descriptive
// Units:   values in any one unit; result in that unit
// Example: the 12 resistivity values of e-Handbook §7.2.6.2 → (95.1567 + 95.1591)/2 = 95.1579 (hand-checked)
func Median(values []float64) (float64, error) {
	if err := checkSample(values, 1); err != nil {
		return 0, err
	}
	return medianOf(values), nil
}

// medianOf returns the median of finite, non-empty values.
func medianOf(values []float64) float64 {
	lower, upper, _ := selection.Median(values) // cannot fail: values is non-empty
	return midpoint(lower, upper)
}
