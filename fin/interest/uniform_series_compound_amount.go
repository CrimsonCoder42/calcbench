package interest

import "math"

// expm1OverX returns (eˣ − 1)/x, with its limit 1 at x = 0, without cancellation for small x.
func expm1OverX(x float64) float64 {
	if x == 0 {
		return 1
	}
	return math.Expm1(x) / x
}

// log1pOverX returns ln(1+x)/x, with its limit 1 at x = 0, without cancellation for small x.
func log1pOverX(x float64) float64 {
	if x == 0 {
		return 1
	}
	return math.Log1p(x) / x
}

// uniformSeriesCompoundAmount returns ((1+i)ⁿ − 1)/i for checked inputs, as
// expm1(x)/x · n · log1p(i)/i with x = n·log1p(i). Every factor is free of cancellation, and i = 0
// gives exactly n.
func uniformSeriesCompoundAmount(i float64, n int) float64 {
	x := float64(n) * math.Log1p(i)
	return expm1OverX(x) * float64(n) * log1pOverX(i)
}

// uniformSeriesPresentWorth returns (1 − (1+i)⁻ⁿ)/i for checked inputs, as
// expm1(−x)/(−x) · n · log1p(i)/i with x = n·log1p(i); i = 0 gives exactly n.
func uniformSeriesPresentWorth(i float64, n int) float64 {
	x := float64(n) * math.Log1p(i)
	return expm1OverX(-x) * float64(n) * log1pOverX(i)
}

// UniformSeriesCompoundAmount returns the factor (F/A,i,n) = ((1+i)ⁿ − 1)/i, the future value at the
// end of period n of one unit paid at the end of each of n periods. It is evaluated as
// expm1(x)/x · n · log1p(i)/i with x = n·log1p(i), which has no cancellation for small i and gives
// exactly n at i = 0.
// Errors: ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if periods < 0; ErrOverflow if
// the factor exceeds float64.
//
// ID:      FIN-INT-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, factor table p. 230, Uniform Series Compound Amount (F/A, i%, n) = [(1 + i)^n − 1]/i
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §4.1, Eq. 7
// Tags:    interest, engineering-economy, factor, compound-amount, uniform-series, annuity
// Units:   ratePerPeriod decimal per period; periods count; result dimensionless (per unit of A)
// Example: i = 1 %, n = 10 → 10.4622 (FE Reference Handbook 10.0.1, factor table p. 233)
func UniformSeriesCompoundAmount(ratePerPeriod float64, periods int) (float64, error) {
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 0); err != nil {
		return 0, err
	}
	return factorResult(uniformSeriesCompoundAmount(ratePerPeriod, periods))
}
