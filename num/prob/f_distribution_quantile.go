package prob

import "fmt"

// FDistributionQuantile returns the u-quantile of Snedecor's F distribution with d₁ = df1 and d₂ = df2
// degrees of freedom: the x with P(X ≤ x) = u, so u = 0.95 gives the upper critical value of an F test at
// α = 0.05. It is found by bisection on the CDF (u ≤ ½) or the survival (u > ½) on a bracket that starts at
// [0, 1] and doubles, until the bracket is within 4 ulps. u = 0 returns 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrProbability when u is outside [0, 1];
// ErrOutOfRange when u = 1 (an infinite quantile), df1 ≤ 0 or df2 ≤ 0, or the incomplete beta does not
// converge; ErrOverflow when the quantile exceeds the float64 range.
//
// ID:      NUM-PROB-083
// Time:    O(B·(√max(d₁, d₂) + log(1/ε)))   B ≤ 1,300 bisection steps (about 70 in practice)
// Space:   O(1)
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 27 "F-Distributions", §1 (cdf), inverted
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.6.7.3 "Upper Critical Values of the F Distribution", https://www.itl.nist.gov/div898/handbook/eda/section3/eda3673.htm (the table the tests reproduce)
// Tags:    probability, distribution, f-distribution, snedecor, quantile, inverse-cdf, critical-value, anova
// Units:   dimensionless
// Example: d₁ = 2, d₂ = 2, u = 0.5 → 1 (F(x) = x/(1 + x)); d₁ = 5, d₂ = 10, u = 0.95 → 3.326 (NIST table)
func FDistributionQuantile(u, df1, df2 float64) (float64, error) {
	if err := checkProbability("u", u); err != nil {
		return 0, err
	}
	if err := checkFDistribution(df1, df2); err != nil {
		return 0, err
	}
	if u == 1 {
		return 0, fmt.Errorf("%w: u = 1, quantile infinite", ErrOutOfRange)
	}
	if u == 0 {
		return 0, nil
	}
	return quantileByBisection(u, 0, 1, func(x float64) (float64, float64, error) {
		if x <= 0 {
			return 0, 1, nil
		}
		return fCDFSurvival(x, df1, df2)
	})
}
