package interest

// GrowingAnnuityPaymentFromPresentValue returns the first end-of-period payment of a series of
// periods payments growing by growthPerPeriod each period whose present value at ratePerPeriod is
// presentValueCents: PV·(r − g)/(1 − ((1+g)/(1+r))ⁿ), rounded half-even to the cent. r = g gives
// PV·(1+r)/n.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |presentValueCents| > 2⁵³; ErrRate if ratePerPeriod or growthPerPeriod is
// NaN, ±Inf or ≤ −1; ErrPeriods if periods < 1; ErrOverflow if the result exceeds ±2⁵³ cents.
//
// ID:      FIN-INT-024
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  A. Damodaran, Applied Corporate Finance, 3rd ed., Wiley, 2010, Appendix 3 "Time Value of Money", "Growing Annuities", PV formula solved for the first payment A(1+g), https://pages.stern.nyu.edu/~adamodar/pdfiles/acf3E/book/app3.pdf
// Tags:    interest, time-value, annuity, growing-annuity, payment, graduated-payment
// Units:   presentValueCents and the result in cents; ratePerPeriod and growthPerPeriod decimal per period; periods count
// Example: PV $16,145,979.98, r = 10 %, g = 3 %, 20 → $1,545,000.00 (inverse of the gold mine; hand-checked)
func GrowingAnnuityPaymentFromPresentValue(presentValueCents int64, ratePerPeriod, growthPerPeriod float64, periods int) (int64, error) {
	if err := checkGrowingAnnuity("presentValueCents", presentValueCents, ratePerPeriod, growthPerPeriod, periods, 1); err != nil {
		return 0, err
	}
	if presentValueCents == 0 {
		return 0, nil
	}
	return roundCents(float64(presentValueCents) / growingAnnuityPresentWorth(ratePerPeriod, growthPerPeriod, periods))
}
