package interest

import (
	"fmt"
	"math"
)

// checkYears returns ErrNotFinite for NaN or ±Inf and ErrNegative for a negative span.
func checkYears(name string, years float64) error {
	if err := checkFinite(name, years); err != nil {
		return err
	}
	if years < 0 {
		return fmt.Errorf("%w: %s = %v", ErrNegative, name, years)
	}
	return nil
}

// FutureValueContinuous returns the value of presentValueCents after years years compounded
// continuously at the nominal ratePerYear, PV·e^(r·t), rounded half-even to the cent. years may be
// fractional; years = 0 returns the amount unchanged.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |presentValueCents| > 2⁵³; ErrRate if ratePerYear is NaN, ±Inf or ≤ −1;
// ErrNotFinite if years is NaN or ±Inf; ErrNegative if years < 0; ErrOverflow if the result exceeds
// ±2⁵³ cents.
//
// ID:      FIN-INT-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §3.2 "Continuous Compounding", Eq. 4 FV_N = PV·e^(r_s·N)
// Source:  G. Strang and E. Herman, Calculus Volume 1, OpenStax, 2016, §6.8 "Exponential Growth and Decay", https://openstax.org/books/calculus-volume-1/pages/6-8-exponential-growth-and-decay
// Tags:    interest, time-value, future-value, continuous-compounding, single-payment
// Units:   presentValueCents and the result in cents; ratePerYear decimal per year; years in years
// Example: $10,000 at 8 % for 2 years → $11,735.11 (Quantitative Investment Analysis ch. 1, Example 6)
func FutureValueContinuous(presentValueCents int64, ratePerYear, years float64) (int64, error) {
	if err := checkCents("presentValueCents", presentValueCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerYear", ratePerYear); err != nil {
		return 0, err
	}
	if err := checkYears("years", years); err != nil {
		return 0, err
	}
	if presentValueCents == 0 || years == 0 {
		return presentValueCents, nil
	}
	return roundCents(float64(presentValueCents) * math.Exp(ratePerYear*years))
}
