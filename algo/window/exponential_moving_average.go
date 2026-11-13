package window

import (
	"errors"
	"fmt"
	"math"
	"slices"
)

// ErrAlpha: the smoothing constant is not in (0, 1].
var ErrAlpha = errors.New("window: smoothing constant alpha must be in (0, 1]")

// ExponentialMovingAverage appends to dst the single exponential smoothing of values and returns the
// extended slice: n smoothed values for n values. Result 0 is values[0] and result i is
// alpha·values[i] + (1 − alpha)·result[i−1]. In the handbook's numbering (observations y₁…yₙ, smoothed
// Sₜ with S₂ = y₁ and Sₜ = α·yₜ₋₁ + (1 − α)·Sₜ₋₁), result i is S_{i+2}: the smoothed value after
// observation i+1, which the handbook uses as the forecast for the next period. alpha = 1 returns
// values unchanged; small alpha smooths more. Results are in the unit of values.
// values is read, never modified. dst must not overlap values.
// Errors: ErrAlpha when alpha is not in (0, 1] (NaN included); ErrNotFinite when a value is NaN or ±Inf;
// ErrOverflow when a result rounds beyond the float64 range (possible only for values within an ulp of
// it). On error dst is returned with its length unchanged.
//
// ID:      ALG-WINDOW-014
// Time:    O(n)   n = len(values); two multiplications and one addition per value
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §6.4.3.1 "Single Exponential Smoothing", https://www.itl.nist.gov/div898/handbook/pmc/section4/pmc431.htm (read 2026-11-12)
// Tags:    window, exponential-moving-average, ema, exponential-smoothing, smoothing, time-series
func ExponentialMovingAverage(dst, values []float64, alpha float64) ([]float64, error) {
	if !(alpha > 0 && alpha <= 1) {
		return dst, fmt.Errorf("%w: alpha = %v", ErrAlpha, alpha)
	}
	for i, v := range values {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return dst, fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, v)
		}
	}
	if len(values) == 0 {
		return dst, nil
	}
	out := slices.Grow(dst, len(values))
	keep := 1 - alpha
	s := values[0]
	out = append(out, s)
	for i, v := range values[1:] {
		// A convex combination of two finite values; it can exceed the range only by rounding.
		s = alpha*v + keep*s
		if math.IsInf(s, 0) {
			return dst, fmt.Errorf("%w: smoothed value at values[%d]", ErrOverflow, i+1)
		}
		out = append(out, s)
	}
	return out, nil
}
