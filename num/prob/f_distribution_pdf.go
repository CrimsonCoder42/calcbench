package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/special"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkFDistribution validates the two degrees of freedom d₁ and d₂: both finite and > 0.
func checkFDistribution(df1, df2 float64) error {
	if err := checkPositive("df1", df1); err != nil {
		return err
	}
	return checkPositive("df2", df2)
}

// FDistributionPDF returns the density of Snedecor's F distribution with d₁ = df1 and d₂ = df2 degrees of
// freedom: with r = d₁x/d₂, f(x) = r^(d₁/2)·(1 + r)^(−(d₁+d₂)/2)/(x·B(d₁/2, d₂/2)) for x > 0, evaluated in log
// space. It is 0 for x < 0. At x = 0 it is 0 for d₁ > 2, 1 for d₁ = 2, and unbounded (ErrOutOfRange) for d₁ < 2.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when df1 ≤ 0 or df2 ≤ 0, or at x = 0
// with df1 < 2; ErrOverflow when the density exceeds the float64 range.
//
// ID:      NUM-PROB-081
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 27 "F-Distributions", §1 (eq. 27.4)
// Source:  Abramowitz, Stegun, Handbook of Mathematical Functions, NBS, 1964, §26.6.1
// Tags:    probability, distribution, f-distribution, snedecor, fisher, pdf, anova
// Units:   dimensionless
// Example: d₁ = 2, d₂ = 2, x = 1 → (1 + 1)⁻² = 0.25
func FDistributionPDF(x, df1, df2 float64) (float64, error) {
	if err := checkFDistribution(df1, df2); err != nil {
		return 0, err
	}
	if err := checkFinite("x", x); err != nil {
		return 0, err
	}
	if x < 0 {
		return 0, nil
	}
	if x == 0 {
		switch {
		case df1 > 2:
			return 0, nil
		case df1 == 2:
			return 1, nil
		default:
			return 0, fmt.Errorf("%w: density unbounded at x = 0 for df1 = %v < 2", ErrOutOfRange, df1)
		}
	}
	lb, err := special.LogBeta(df1/2, df2/2)
	if err != nil {
		return 0, specialErr("log beta", err)
	}
	r := df1 * x / df2
	logR := math.Log(df1) + math.Log(x) - math.Log(df2) // r itself may overflow or underflow
	var log1pR float64
	if math.IsInf(r, 1) || r > 1e300 {
		log1pR = logR
	} else {
		log1pR = math.Log1p(r)
	}
	d := math.Exp(df1/2*logR - (df1+df2)/2*log1pR - math.Log(x) - lb)
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density at x = %v, df1 = %v, df2 = %v", ErrOverflow, x, df1, df2)
	}
	return d, nil
}
