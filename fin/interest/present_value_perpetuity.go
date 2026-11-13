package interest

import "fmt"

// PresentValuePerpetuity returns the value now of paymentCents paid at the end of every period
// forever, discounted at ratePerPeriod: C/r, rounded half-even to the cent. The first payment is one
// period from now. The sum converges only for r > 0.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |paymentCents| > 2⁵³; ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1;
// ErrDomain if ratePerPeriod ≤ 0; ErrOverflow if the result exceeds ±2⁵³ cents.
//
// ID:      FIN-INT-026
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  J. Dahlquist and R. Knight, Principles of Finance, OpenStax, 2022, §8.1 "Perpetuities", Eq. 8.1 PV = C/Rs, https://openstax.org/books/principles-finance/pages/8-1-perpetuities
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §6.2, Eq. 13 PV = A/r
// Tags:    interest, time-value, perpetuity, present-value, consol, preferred-stock, capitalized-cost
// Units:   paymentCents and the result in cents; ratePerPeriod decimal per period
// Example: $2.00 at 7 % → $28.57 (Principles of Finance §8.1, Eq. 8.2); £100 at 5 % → £2,000 (QIA ch. 1, Example 14)
func PresentValuePerpetuity(paymentCents int64, ratePerPeriod float64) (int64, error) {
	if err := checkCents("paymentCents", paymentCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if ratePerPeriod <= 0 {
		return 0, fmt.Errorf("%w: ratePerPeriod = %v, a perpetuity needs a rate > 0", ErrDomain, ratePerPeriod)
	}
	return roundCents(float64(paymentCents) / ratePerPeriod)
}
