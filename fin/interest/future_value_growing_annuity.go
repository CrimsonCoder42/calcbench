package interest

// FutureValueGrowingAnnuity returns the value at the end of period n of periods end-of-period
// payments, the first firstPaymentCents and each later one growthPerPeriod larger, compounded at
// ratePerPeriod: C·((1+r)ⁿ − (1+g)ⁿ)/(r − g) = PresentValueGrowingAnnuity·(1+r)ⁿ, rounded half-even to
// the cent. r = g gives C·n·(1+r)ⁿ⁻¹.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |firstPaymentCents| > 2⁵³; ErrRate if ratePerPeriod or growthPerPeriod is
// NaN, ±Inf or ≤ −1; ErrPeriods if periods < 0; ErrOverflow if the result exceeds ±2⁵³ cents.
//
// ID:      FIN-INT-023
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  A. Damodaran, Applied Corporate Finance, 3rd ed., Wiley, 2010, Appendix 3 "Time Value of Money", "Growing Annuities" (PV), carried forward n periods by FV = PV(1 + r)^n (R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 §3, Eq. 2), https://pages.stern.nyu.edu/~adamodar/pdfiles/acf3E/book/app3.pdf
// Tags:    interest, time-value, annuity, growing-annuity, future-value, savings
// Units:   firstPaymentCents and the result in cents; ratePerPeriod and growthPerPeriod decimal per period; periods count
// Example: Damodaran App. 3 gold mine ($1,545,000 first payment, 10 %, 3 %, 20 years): $16,145,979.98 × 1.1²⁰ = $108,622,079.49 (hand-checked, exact rational arithmetic)
func FutureValueGrowingAnnuity(firstPaymentCents int64, ratePerPeriod, growthPerPeriod float64, periods int) (int64, error) {
	if err := checkGrowingAnnuity("firstPaymentCents", firstPaymentCents, ratePerPeriod, growthPerPeriod, periods, 0); err != nil {
		return 0, err
	}
	if firstPaymentCents == 0 || periods == 0 {
		return 0, nil
	}
	f := growingAnnuityPresentWorth(ratePerPeriod, growthPerPeriod, periods) * compound(ratePerPeriod, periods)
	return roundCents(float64(firstPaymentCents) * f)
}
