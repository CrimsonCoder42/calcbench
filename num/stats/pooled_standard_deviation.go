package stats

import (
	"fmt"
	"math"
)

// PooledStandardDeviation returns the pooled standard deviation of two samples from their sizes and
// sample standard deviations, s_p = √(((n1 − 1)·s1² + (n2 − 1)·s2²)/(n1 + n2 − 2)), computed with
// math.Hypot so that no square overflows. It never exceeds max(sd1, sd2), so it is always finite.
// Errors: ErrTooFew when n1 or n2 is below 1 or n1 + n2 < 3; ErrNotFinite when sd1 or sd2 is NaN or ±Inf;
// ErrNegative when sd1 or sd2 is negative.
//
// ID:      NUM-STATS-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.3 "Two-Sample t-Test for Equal Means", https://www.itl.nist.gov/div898/handbook/eda/section3/eda353.htm (read 2026-11-12)
// Tags:    pooled-standard-deviation, standard-deviation, two-sample, dispersion, descriptive
// Units:   n1, n2 dimensionless counts; sd1, sd2 in the data's unit; result in that unit
// Example: e-Handbook §1.3.5.3, AUTO83B.DAT: n1 = 249, s1 = 6.41470; n2 = 79, s2 = 6.10771 → s_p = 6.34260 (as printed)
func PooledStandardDeviation(n1 int, sd1 float64, n2 int, sd2 float64) (float64, error) {
	if n1 < 1 {
		return 0, fmt.Errorf("%w: n1 = %d, want ≥ 1", ErrTooFew, n1)
	}
	if n2 < 1 {
		return 0, fmt.Errorf("%w: n2 = %d, want ≥ 1", ErrTooFew, n2)
	}
	if err := checkScalar("sd1", sd1); err != nil {
		return 0, err
	}
	if err := checkScalar("sd2", sd2); err != nil {
		return 0, err
	}
	if sd1 < 0 {
		return 0, fmt.Errorf("%w: sd1 = %v", ErrNegative, sd1)
	}
	if sd2 < 0 {
		return 0, fmt.Errorf("%w: sd2 = %v", ErrNegative, sd2)
	}
	// float64 counts: n1 + n2 cannot overflow, and are exact below 2^53.
	df := float64(n1) + float64(n2) - 2
	if df < 1 {
		return 0, fmt.Errorf("%w: n1 + n2 = %v, want ≥ 3", ErrTooFew, df+2)
	}
	r := math.Hypot(math.Sqrt((float64(n1)-1)/df)*sd1, math.Sqrt((float64(n2)-1)/df)*sd2)
	// s_p is a weighted root mean square of sd1 and sd2, so it never exceeds the larger; the clamp
	// only undoes rounding, and keeps s_p finite.
	return min(r, max(sd1, sd2)), nil
}
