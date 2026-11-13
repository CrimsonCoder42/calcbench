package stats

// Midrange returns (smallest + largest)/2 of values, computed without overflow. NIST notes that it is
// not robust: one extreme value moves it. values is read, never modified.
// Errors: ErrEmpty when values is empty; ErrNotFinite when a value is NaN or ±Inf.
//
// ID:      NUM-STATS-013
// Time:    O(n)   n = len(values); one pass
// Space:   O(1) extra, no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm (read 2026-11-12)
// Tags:    midrange, mid-range, location, central-tendency, descriptive
// Units:   values in any one unit; result in that unit
// Example: the 12 resistivity values of e-Handbook §7.2.6.2 → (95.0610 + 95.1990)/2 = 95.1300 (hand-checked)
func Midrange(values []float64) (float64, error) {
	if err := checkSample(values, 1); err != nil {
		return 0, err
	}
	lo, hi := minMax(values)
	return midpoint(lo, hi), nil
}
