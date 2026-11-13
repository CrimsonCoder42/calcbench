package interest

// AnnuityPaymentFromPresentValue returns the level end-of-period payment that repays
// presentValueCents over periods periods at ratePerPeriod, PV·(A/P,i,n) = PV·i/(1 − (1+i)⁻ⁿ), rounded
// half-even to the cent. This is the level payment of a fully amortizing loan. At i = 0 it is PV/n.
// The payment is rounded once; a schedule built on it must absorb the rounding in its last payment.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |presentValueCents| > 2⁵³; ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1;
// ErrPeriods if periods < 1; ErrOverflow if the payment exceeds ±2⁵³ cents.
//
// ID:      FIN-INT-018
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §7.3 "Solving for the Size of Annuity Payments", A = PV/present value annuity factor (Eq. 11)
// Source:  J. Dahlquist and R. Knight, Principles of Finance, OpenStax, 2022, §8.2 Eq. 8.30 and §8.3 "Loan Amortization" Eq. 8.31, https://openstax.org/books/principles-finance/pages/8-3-loan-amortization
// Tags:    interest, time-value, annuity, payment, loan-payment, amortization, ordinary-annuity
// Units:   presentValueCents and the result in cents; ratePerPeriod decimal per period; periods count
// Example: $100,000, 0.08/12, 360 → $733.76 (QIA ch. 1, Example 20); $32,000, 0.06/12, 36 → $973.50 (Principles of Finance §8.3)
func AnnuityPaymentFromPresentValue(presentValueCents int64, ratePerPeriod float64, periods int) (int64, error) {
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
	return roundCents(float64(presentValueCents) / uniformSeriesPresentWorth(ratePerPeriod, periods))
}
