package interest

// GrowingAnnuityPaymentFromFutureValue returns the first end-of-period payment of a series of
// periods payments growing by growthPerPeriod each period that accumulates to futureValueCents at
// the end of period n at ratePerPeriod: FV·(r − g)/((1+r)ⁿ − (1+g)ⁿ), rounded half-even to the cent.
// r = g gives FV/(n·(1+r)ⁿ⁻¹).
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |futureValueCents| > 2⁵³; ErrRate if ratePerPeriod or growthPerPeriod is
// NaN, ±Inf or ≤ −1; ErrPeriods if periods < 1; ErrOverflow if the result exceeds ±2⁵³ cents.
//
// ID:      FIN-INT-025
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  A. Damodaran, Applied Corporate Finance, 3rd ed., Wiley, 2010, Appendix 3 "Time Value of Money", "Growing Annuities", PV formula times (1 + r)^n, solved for the first payment, https://pages.stern.nyu.edu/~adamodar/pdfiles/acf3E/book/app3.pdf
// Tags:    interest, time-value, annuity, growing-annuity, payment, savings-goal
// Units:   futureValueCents and the result in cents; ratePerPeriod and growthPerPeriod decimal per period; periods count
// Example: FV $108,622,079.49, r = 10 %, g = 3 %, 20 → $1,545,000.00 (hand-checked)
func GrowingAnnuityPaymentFromFutureValue(futureValueCents int64, ratePerPeriod, growthPerPeriod float64, periods int) (int64, error) {
	if err := checkGrowingAnnuity("futureValueCents", futureValueCents, ratePerPeriod, growthPerPeriod, periods, 1); err != nil {
		return 0, err
	}
	if futureValueCents == 0 {
		return 0, nil
	}
	f := growingAnnuityPresentWorth(ratePerPeriod, growthPerPeriod, periods) * compound(ratePerPeriod, periods)
	return roundCents(float64(futureValueCents) / f)
}
