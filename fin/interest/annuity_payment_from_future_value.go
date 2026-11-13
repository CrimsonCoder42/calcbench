package interest

// AnnuityPaymentFromFutureValue returns the level end-of-period payment that grows to futureValueCents after periods periods at ratePerPeriod, FV·i/((1+i)ⁿ − 1) = FV·(A/F,i,n), rounded half-even to the cent.
// This is the sinking-fund deposit; i = 0 gives FV/n.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |futureValueCents| > 2⁵³; ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1;
// ErrPeriods if periods < 1. |payment| ≤ |futureValueCents|, so it cannot overflow.
//
// ID:      FIN-INT-019
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, factor table p. 230, Uniform Series Sinking Fund (A/F, i%, n); A = F·(A/F, i%, n)
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §4.1, Eq. 7 solved for A
// Tags:    interest, time-value, annuity, payment, sinking-fund, ordinary-annuity, savings-goal
// Units:   futureValueCents and the result in cents; ratePerPeriod decimal per period; periods count
// Example: FV $5,525.63, 5 %, 5 → $1,000.00 (inverse of QIA ch. 1 Figure 3; hand-checked)
func AnnuityPaymentFromFutureValue(futureValueCents int64, ratePerPeriod float64, periods int) (int64, error) {
	if err := checkCents("futureValueCents", futureValueCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 1); err != nil {
		return 0, err
	}
	if futureValueCents == 0 {
		return 0, nil
	}
	return roundCents(float64(futureValueCents) / uniformSeriesCompoundAmount(ratePerPeriod, periods))
}
