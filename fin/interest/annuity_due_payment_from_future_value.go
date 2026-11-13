package interest

// AnnuityDuePaymentFromFutureValue returns the level start-of-period payment that grows to futureValueCents at the end of period n at ratePerPeriod, FV·i/(((1+i)ⁿ − 1)(1+i)), rounded half-even to the cent.
// It is AnnuityPaymentFromFutureValue's payment divided by (1+i); i = 0 gives FV/n.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |futureValueCents| > 2⁵³; ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1;
// ErrPeriods if periods < 1; ErrOverflow if the result exceeds ±2⁵³ cents.
//
// ID:      FIN-INT-021
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  J. Dahlquist and R. Knight, Principles of Finance, OpenStax, 2022, §8.2, Eq. 8.25 (annuity due FV) solved for PYMT, https://openstax.org/books/principles-finance/pages/8-2-annuities
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §4.1, Eq. 7
// Tags:    interest, time-value, annuity, payment, annuity-due, sinking-fund, savings-goal
// Units:   futureValueCents and the result in cents; ratePerPeriod decimal per period; periods count
// Example: FV $5,801.91, 5 %, 5 → $1,000.00 (inverse of FutureValueAnnuityDue's example; hand-checked)
func AnnuityDuePaymentFromFutureValue(futureValueCents int64, ratePerPeriod float64, periods int) (int64, error) {
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
	return roundCents(float64(futureValueCents) / (uniformSeriesCompoundAmount(ratePerPeriod, periods) * (1 + ratePerPeriod)))
}
