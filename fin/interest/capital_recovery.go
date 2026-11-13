package interest

// CapitalRecovery returns the factor (A/P,i,n) = i(1+i)ⁿ/((1+i)ⁿ − 1), the end-of-period payment over
// n periods that repays one unit borrowed now. For i ≥ 0 it is evaluated as i + (A/F,i,n) (the
// handbook identity; two positive terms), for i < 0 as 1/(P/A,i,n); both are free of cancellation and
// i = 0 gives exactly 1/n. The factor is at most 1 + i ((F/A) ≥ 1 and (P/A) ≥ 1 when i < 0), so it
// never overflows, even for i near math.MaxFloat64.
// Errors: ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if periods < 1.
//
// ID:      FIN-INT-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, factor table p. 230, Capital Recovery (A/P, i%, n) = i(1 + i)^n/[(1 + i)^n − 1]
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §7.3 (PV ÷ present value annuity factor)
// Tags:    interest, engineering-economy, factor, capital-recovery, loan-payment, uniform-series, annuity
// Units:   ratePerPeriod decimal per period; periods count; result dimensionless (per unit of P)
// Example: i = 1 %, n = 10 → 0.1056 (FE Reference Handbook 10.0.1, factor table p. 233)
func CapitalRecovery(ratePerPeriod float64, periods int) (float64, error) {
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 1); err != nil {
		return 0, err
	}
	if ratePerPeriod >= 0 {
		return ratePerPeriod + 1/uniformSeriesCompoundAmount(ratePerPeriod, periods), nil
	}
	return 1 / uniformSeriesPresentWorth(ratePerPeriod, periods), nil
}
