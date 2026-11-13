package interest

// UniformSeriesSinkingFund returns the factor (A/F,i,n) = i/((1+i)ⁿ − 1), the end-of-period payment
// that grows to one unit after n periods. It is 1/(F/A,i,n) with (F/A) evaluated without
// cancellation; i = 0 gives exactly 1/n, and a factor too small for float64 returns 0.
// (F/A,i,n) = 1 + (1+i) + … + (1+i)ⁿ⁻¹ ≥ 1 for n ≥ 1, so the factor is at most 1 and never overflows.
// Errors: ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if periods < 1.
//
// ID:      FIN-INT-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, factor table p. 230, Uniform Series Sinking Fund (A/F, i%, n) = i/[(1 + i)^n − 1]
// Tags:    interest, engineering-economy, factor, sinking-fund, uniform-series, annuity
// Units:   ratePerPeriod decimal per period; periods count; result dimensionless (per unit of F)
// Example: i = 1 %, n = 10 → 0.0956 (FE Reference Handbook 10.0.1, factor table p. 233)
func UniformSeriesSinkingFund(ratePerPeriod float64, periods int) (float64, error) {
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 1); err != nil {
		return 0, err
	}
	return 1 / uniformSeriesCompoundAmount(ratePerPeriod, periods), nil
}
