package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ChiSquareMeanVariance returns the mean k and the variance 2k of a chi-square variable with k = df
// degrees of freedom.
// Errors: ErrNotFinite when df is NaN or ±Inf; ErrOutOfRange when df ≤ 0; ErrOverflow when 2k exceeds the
// float64 range.
//
// ID:      NUM-PROB-097
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 1, 2nd ed., Wiley, 1994, ch. 18 "Chi-Square Distributions Including Chi and Rayleigh", §3 (moments)
// Tags:    probability, distribution, chi-square, mean, variance
// Units:   dimensionless
// Example: df = 3 → (3, 6)
func ChiSquareMeanVariance(df float64) (mean, variance float64, err error) {
	if err := checkPositive("df", df); err != nil {
		return 0, 0, err
	}
	if !reduce.AllFinite(2 * df) {
		return 0, 0, fmt.Errorf("%w: variance 2·df, df = %v", ErrOverflow, df)
	}
	return df, 2 * df, nil
}
