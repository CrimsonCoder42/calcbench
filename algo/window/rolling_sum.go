package window

import (
	"errors"
	"fmt"
	"math"
	"slices"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrWindowLen: windowLen is less than 1.
	ErrWindowLen = errors.New("window: window length must be at least 1")
	// ErrNotFinite: a value is NaN or ±Inf, which a sum cannot slide past.
	ErrNotFinite = errors.New("window: value is NaN or infinite")
	// ErrOverflow: a window's sum is beyond the float64 range.
	ErrOverflow = errors.New("window: window sum exceeds the float64 range")
)

// RollingSum appends to dst the sum of each window of windowLen consecutive values and returns the
// extended slice: n-windowLen+1 sums for n values, none when windowLen > n.
// Sums are compensated (Neumaier) and restarted from scratch every windowLen windows, so a window's
// error depends only on the at most 2·windowLen-1 values since the restart, never on earlier history.
// values is read, never modified.
// Errors: ErrWindowLen when windowLen < 1; ErrNotFinite when a value is NaN or ±Inf; ErrOverflow when
// a window's sum, or a partial sum inside it, exceeds the float64 range. On error dst is returned with
// its length unchanged.
//
// ID:      ALG-WINDOW-001
// Time:    O(n)   n = len(values); at most 3n additions
// Space:   O(1) extra; O(n-w+1) for the result (w = windowLen), one allocation only if dst lacks capacity
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §6.4.2.1 "Single Moving Average", https://www.itl.nist.gov/div898/handbook/pmc/section4/pmc421.htm (read 2026-11-12)
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §4.3 (compensated summation)
// Tags:    window, sliding-window, rolling, moving-sum, compensated-summation, time-series
func RollingSum(dst, values []float64, windowLen int) ([]float64, error) {
	return appendWindowSums(dst, values, windowLen, 1)
}

// appendWindowSums appends each window's compensated sum divided by divisor. RollingSum passes 1 and
// RollingMean passes windowLen.
func appendWindowSums(dst, values []float64, windowLen int, divisor float64) ([]float64, error) {
	if windowLen < 1 {
		return dst, fmt.Errorf("%w: windowLen = %d", ErrWindowLen, windowLen)
	}
	for i, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return dst, fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, v)
		}
	}
	n := len(values)
	if windowLen > n {
		return dst, nil
	}
	out := slices.Grow(dst, n-windowLen+1)
	var sum, comp float64
	sinceRestart := windowLen // forces a restart at the first window
	for end := windowLen - 1; end < n; end++ {
		if sinceRestart == windowLen {
			// Restart from scratch every windowLen windows: rounding left in comp by values that
			// have since left the window is dropped, so a window's error depends only on the at
			// most 2·windowLen-1 values since the restart. Costs n extra additions in total.
			sum, comp = 0, 0
			for _, v := range values[end-windowLen+1 : end+1] {
				sum, comp = reduce.NeumaierAdd(sum, comp, v)
			}
			sinceRestart = 0
		} else {
			// Remove the leaving value before adding the entering one, so the running sum never
			// holds more than one window.
			sum, comp = reduce.NeumaierAdd(sum, comp, -values[end-windowLen])
			sum, comp = reduce.NeumaierAdd(sum, comp, values[end])
		}
		sinceRestart++
		r := (sum + comp) / divisor
		if math.IsInf(r, 0) || math.IsNaN(r) {
			return dst, fmt.Errorf("%w: window ending at values[%d]", ErrOverflow, end)
		}
		out = append(out, r)
	}
	return out, nil
}
