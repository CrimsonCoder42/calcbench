package window

// RollingMean appends to dst the mean of each window of windowLen consecutive values (the single
// moving average) and returns the extended slice: n-windowLen+1 means for n values, none when
// windowLen > n. Each mean is the window sum of RollingSum (compensated, restarted every windowLen
// windows) divided by windowLen. values is read, never modified.
// Errors: ErrWindowLen when windowLen < 1; ErrNotFinite when a value is NaN or ±Inf; ErrOverflow when
// a window's sum, or a partial sum inside it, exceeds the float64 range, even if the mean would not
// (doc/ASSUMPTIONS.md A-007). On error dst is returned with its length unchanged.
//
// ID:      ALG-WINDOW-002
// Time:    O(n)   n = len(values); at most 3n additions
// Space:   O(1) extra; O(n-w+1) for the result (w = windowLen), one allocation only if dst lacks capacity
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §6.4.2.1 "Single Moving Average", https://www.itl.nist.gov/div898/handbook/pmc/section4/pmc421.htm (read 2026-11-12)
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §4.3 (compensated summation)
// Tags:    window, sliding-window, rolling, moving-average, mean, smoothing, time-series
func RollingMean(dst, values []float64, windowLen int) ([]float64, error) {
	return appendWindowSums(dst, values, windowLen, float64(windowLen))
}
