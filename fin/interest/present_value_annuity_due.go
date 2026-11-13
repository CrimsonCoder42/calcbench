package interest

// PresentValueAnnuityDue returns the value now of paymentCents paid at the start of each of periods periods (an annuity due), PMT·(1 − (1+i)⁻ⁿ)/i·(1+i), rounded half-even to the cent.
// The first payment is not discounted; i = 0 gives PMT·n and periods = 0 gives 0.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |paymentCents| > 2⁵³; ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1;
// ErrPeriods if periods < 0; ErrOverflow if the result exceeds ±2⁵³ cents.
//
// ID:      FIN-INT-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §6.1, footnote 7: PV(annuity due) = A{[1 − (1 + r)^−N]/r}(1 + r)
// Source:  J. Dahlquist and R. Knight, Principles of Finance, OpenStax, 2022, §8.2, Eq. 8.17, https://openstax.org/books/principles-finance/pages/8-2-annuities
// Tags:    interest, time-value, annuity, present-value, annuity-due, lease
// Units:   paymentCents and the result in cents; ratePerPeriod decimal per period; periods count
// Example: four payments of $100 at 12 %, first today → $340.18 (Quantitative Investment Analysis ch. 1, §6.1)
func PresentValueAnnuityDue(paymentCents int64, ratePerPeriod float64, periods int) (int64, error) {
	if err := checkCents("paymentCents", paymentCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 0); err != nil {
		return 0, err
	}
	if paymentCents == 0 {
		return 0, nil
	}
	return roundCents(float64(paymentCents) * uniformSeriesPresentWorth(ratePerPeriod, periods) * (1 + ratePerPeriod))
}
