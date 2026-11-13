package interest

// PresentValue returns the value today of futureValueCents received after periods compounding
// periods at ratePerPeriod, FV/(1+i)ⁿ = FV·(1+i)⁻ⁿ, rounded half-even to the cent. periods = 0
// returns the amount unchanged; a very long horizon can round to 0.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |futureValueCents| > 2⁵³; ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1;
// ErrPeriods if periods < 0; ErrOverflow if the result exceeds ±2⁵³ cents (a negative rate).
//
// ID:      FIN-INT-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  J. Dahlquist and R. Knight, Principles of Finance, OpenStax, 2022, §7.4 "Applications of TVM in Finance", Eq. 7.32 PV = FV × 1/(1+r)^n, https://openstax.org/books/principles-finance/pages/7-4-applications-of-tvm-in-finance
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §5, Eq. 8
// Tags:    interest, time-value, present-value, discount, single-payment
// Units:   futureValueCents and the result in cents; ratePerPeriod decimal per period; periods count
// Example: FV $50,000, 5 %, 5 periods → $39,176.31 (Principles of Finance §7.3, Example 2); $1,000, 5 %, 30 → $231.38 (§7.4)
func PresentValue(futureValueCents int64, ratePerPeriod float64, periods int) (int64, error) {
	if err := checkCents("futureValueCents", futureValueCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 0); err != nil {
		return 0, err
	}
	if futureValueCents == 0 || periods == 0 {
		return futureValueCents, nil
	}
	return roundCents(float64(futureValueCents) * compound(ratePerPeriod, -periods))
}
