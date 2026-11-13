package prob

import "fmt"

// StudentTMeanVariance returns the mean 0 and the variance ν/(ν − 2) of Student's t distribution with
// ν = df degrees of freedom. The variance is finite only for ν > 2 (the mean exists for ν > 1 but the pair
// is returned together, so ν ≤ 2 is an error).
// Errors: ErrNotFinite when df is NaN or ±Inf; ErrOutOfRange when df ≤ 2.
//
// ID:      NUM-PROB-080
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 28 "t-Distributions", §2 (moments)
// Tags:    probability, distribution, student-t, t-distribution, mean, variance
// Units:   dimensionless
// Example: ν = 4 → (0, 2)
func StudentTMeanVariance(df float64) (mean, variance float64, err error) {
	if err := checkPositive("df", df); err != nil {
		return 0, 0, err
	}
	if df <= 2 {
		return 0, 0, fmt.Errorf("%w: df = %v, variance infinite or undefined for df ≤ 2", ErrOutOfRange, df)
	}
	return 0, df / (df - 2), nil
}
