package prob

import "math"

// fCDFSurvival returns P(X ≤ x) and P(X > x) for validated d₁, d₂ and x > 0. With r = d₁x/d₂ the incomplete
// beta argument w = r/(1 + r) and its complement 1/(1 + r) are both formed without cancellation.
func fCDFSurvival(x, df1, df2 float64) (cdf, survival float64, err error) {
	r := df1 * x / df2
	if math.IsInf(r, 1) {
		return 1, 0, nil
	}
	return betaSplit(r/(1+r), 1/(1+r), df1/2, df2/2)
}

// FDistributionCDF returns P(X ≤ x) for Snedecor's F distribution with d₁ = df1 and d₂ = df2 degrees of
// freedom: I_w(d₁/2, d₂/2) with w = d₁x/(d₁x + d₂), from the regularized incomplete beta (num/special
// NUM-SPEC-005), and 0 for x ≤ 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when df1 ≤ 0 or df2 ≤ 0, or when the
// incomplete beta does not converge (wrapping special.ErrNoConvergence).
//
// ID:      NUM-PROB-082
// Time:    O(√max(d₁, d₂) + log(1/ε))   ε = 2⁻⁵²
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 27 "F-Distributions", §1
// Source:  Abramowitz, Stegun, Handbook of Mathematical Functions, NBS, 1964, §26.6.2 (Q(F|ν₁,ν₂) = I_x(ν₂/2, ν₁/2), x = ν₂/(ν₂ + ν₁F))
// Tags:    probability, distribution, f-distribution, snedecor, fisher, cdf, anova, p-value
// Units:   dimensionless; result a probability in [0, 1]
// Example: d₁ = 2, d₂ = 2, x = 1 → x/(1 + x) = 0.5
func FDistributionCDF(x, df1, df2 float64) (float64, error) {
	if err := checkFDistribution(df1, df2); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x <= 0 {
		return 0, nil
	}
	c, _, err := fCDFSurvival(x, df1, df2)
	if err != nil {
		return 0, err
	}
	return c, nil
}
