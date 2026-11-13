package interest

// AnnuityDuePaymentFromPresentValue returns the level start-of-period payment over periods periods whose present value at ratePerPeriod is presentValueCents, PV·i/((1 − (1+i)⁻ⁿ)(1+i)), rounded half-even to the cent.
// It is AnnuityPaymentFromPresentValue's payment divided by (1+i) (payments in advance, as in a lease); i = 0 gives PV/n.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |presentValueCents| > 2⁵³; ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1;
// ErrPeriods if periods < 1. |payment| ≤ |presentValueCents|, so it cannot overflow.
//
// ID:      FIN-INT-020
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §6.1, footnote 7 (annuity due PV) solved for A
// Source:  J. Dahlquist and R. Knight, Principles of Finance, OpenStax, 2022, §8.2, Eq. 8.17, https://openstax.org/books/principles-finance/pages/8-2-annuities
// Tags:    interest, time-value, annuity, payment, annuity-due, lease-payment
// Units:   presentValueCents and the result in cents; ratePerPeriod decimal per period; periods count
// Example: PV $340.18, 12 %, 4 → $100.00 (inverse of QIA ch. 1 §6.1 annuity due; hand-checked)
func AnnuityDuePaymentFromPresentValue(presentValueCents int64, ratePerPeriod float64, periods int) (int64, error) {
	if err := checkCents("presentValueCents", presentValueCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 1); err != nil {
		return 0, err
	}
	if presentValueCents == 0 {
		return 0, nil
	}
	return roundCents(float64(presentValueCents) / (uniformSeriesPresentWorth(ratePerPeriod, periods) * (1 + ratePerPeriod)))
}
