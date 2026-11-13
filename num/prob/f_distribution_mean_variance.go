package prob

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// FDistributionMeanVariance returns the mean d₂/(d₂ − 2) and the variance
// 2d₂²(d₁ + d₂ − 2)/(d₁(d₂ − 2)²(d₂ − 4)) of Snedecor's F distribution. The variance is finite only for
// d₂ > 4 (the mean exists for d₂ > 2 but the pair is returned together, so d₂ ≤ 4 is an error).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when df1 ≤ 0 or df2 ≤ 4; ErrOverflow
// when the variance exceeds the float64 range.
//
// ID:      NUM-PROB-084
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 27 "F-Distributions", §2 (moments)
// Source:  Abramowitz, Stegun, Handbook of Mathematical Functions, NBS, 1964, §26.6.1 (mean and variance)
// Tags:    probability, distribution, f-distribution, snedecor, mean, variance
// Units:   dimensionless
// Example: d₁ = 5, d₂ = 10 → (1.25, 1.3541666666666667) (2·100·13/(5·64·6) = 2600/1920)
func FDistributionMeanVariance(df1, df2 float64) (mean, variance float64, err error) {
	if err := checkFDistribution(df1, df2); err != nil {
		return 0, 0, err
	}
	if df2 <= 4 {
		return 0, 0, fmt.Errorf("%w: df2 = %v, variance infinite or undefined for df2 ≤ 4", ErrOutOfRange, df2)
	}
	mean = df2 / (df2 - 2)
	variance = 2 * mean * mean * (df1 + df2 - 2) / (df1 * (df2 - 4))
	if !reduce.AllFinite(variance) {
		return 0, 0, fmt.Errorf("%w: df1 = %v, df2 = %v", ErrOverflow, df1, df2)
	}
	return mean, variance, nil
}
