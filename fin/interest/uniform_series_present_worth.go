package interest

// UniformSeriesPresentWorth returns the factor (P/A,i,n) = ((1+i)ⁿ − 1)/(i(1+i)ⁿ) = (1 − (1+i)⁻ⁿ)/i,
// the present worth of one unit paid at the end of each of n periods. It is evaluated as
// expm1(−x)/(−x) · n · log1p(i)/i with x = n·log1p(i), with no cancellation for small i; i = 0 gives
// exactly n.
// Errors: ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if periods < 0; ErrOverflow if
// the factor exceeds float64 (a negative rate over many periods).
//
// ID:      FIN-INT-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, factor table p. 230, Uniform Series Present Worth (P/A, i%, n) = [(1 + i)^n − 1]/[i(1 + i)^n]
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §6.1, Eq. 11
// Tags:    interest, engineering-economy, factor, present-worth, uniform-series, annuity
// Units:   ratePerPeriod decimal per period; periods count; result dimensionless (per unit of A)
// Example: i = 1 %, n = 10 → 9.4713 (FE Reference Handbook 10.0.1, factor table p. 233)
func UniformSeriesPresentWorth(ratePerPeriod float64, periods int) (float64, error) {
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 0); err != nil {
		return 0, err
	}
	return factorResult(uniformSeriesPresentWorth(ratePerPeriod, periods))
}
