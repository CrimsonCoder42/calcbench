package interest

// FutureValueAnnuityDue returns the value at the end of period n of paymentCents paid at the start of each of periods periods (an annuity due), PMT·((1+i)ⁿ − 1)/i·(1+i), rounded half-even to the cent.
// Each payment earns one more period of interest than in FutureValueAnnuity; i = 0 gives PMT·n and periods = 0 gives 0.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |paymentCents| > 2⁵³; ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1;
// ErrPeriods if periods < 0; ErrOverflow if the result exceeds ±2⁵³ cents.
//
// ID:      FIN-INT-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  J. Dahlquist and R. Knight, Principles of Finance, OpenStax, 2022, §8.2 "Annuities", Eq. 8.25 (ordinary annuity FV times one more period of interest, (1 + i)), https://openstax.org/books/principles-finance/pages/8-2-annuities
// Tags:    interest, time-value, annuity, future-value, annuity-due, savings
// Units:   paymentCents and the result in cents; ratePerPeriod decimal per period; periods count
// Example: $1,000, 5 %, 5 periods → $5,801.91 = $5,525.63 × 1.05 (hand-checked from QIA ch. 1 Figure 3)
func FutureValueAnnuityDue(paymentCents int64, ratePerPeriod float64, periods int) (int64, error) {
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
	return roundCents(float64(paymentCents) * uniformSeriesCompoundAmount(ratePerPeriod, periods) * (1 + ratePerPeriod))
}
