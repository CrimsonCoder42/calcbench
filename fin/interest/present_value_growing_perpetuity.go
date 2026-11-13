package interest

import "fmt"

// PresentValueGrowingPerpetuity returns the value now of a payment stream that starts with
// firstPaymentCents one period from now and grows by growthPerPeriod every period forever,
// discounted at ratePerPeriod: C/(r − g), rounded half-even to the cent. This is also the terminal
// value of a constant-growth cash flow and the Gordon growth price. The sum converges only for r > g.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |firstPaymentCents| > 2⁵³; ErrRate if ratePerPeriod or growthPerPeriod is
// NaN, ±Inf or ≤ −1; ErrDomain if ratePerPeriod ≤ growthPerPeriod; ErrOverflow if the result exceeds
// ±2⁵³ cents.
//
// ID:      FIN-INT-027
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  J. Dahlquist and R. Knight, Principles of Finance, OpenStax, 2022, §8.1 "Perpetuities", Eq. 8.3 PV = C/(Rs − G), https://openstax.org/books/principles-finance/pages/8-1-perpetuities
// Source:  A. Damodaran, Applied Corporate Finance, 3rd ed., Wiley, 2010, Appendix 3 "Time Value of Money", "Growing Perpetuities", PV = CF1/(r − g), https://pages.stern.nyu.edu/~adamodar/pdfiles/acf3E/book/app3.pdf
// Tags:    interest, time-value, perpetuity, growing-perpetuity, present-value, terminal-value, gordon-growth
// Units:   firstPaymentCents and the result in cents; ratePerPeriod and growthPerPeriod decimal per period
// Example: $1.75, 5.8 %, 1.2 % → $38.04; $4.00, 3 %, 0.2 % → $142.86 (Principles of Finance §8.1)
func PresentValueGrowingPerpetuity(firstPaymentCents int64, ratePerPeriod, growthPerPeriod float64) (int64, error) {
	if err := checkCents("firstPaymentCents", firstPaymentCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkRate("growthPerPeriod", growthPerPeriod); err != nil {
		return 0, err
	}
	if ratePerPeriod <= growthPerPeriod {
		return 0, fmt.Errorf("%w: ratePerPeriod = %v must exceed growthPerPeriod = %v", ErrDomain, ratePerPeriod, growthPerPeriod)
	}
	return roundCents(float64(firstPaymentCents) / (ratePerPeriod - growthPerPeriod))
}
