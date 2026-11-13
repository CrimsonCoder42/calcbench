package interest

import "math"

// PresentValueContinuous returns the value now of futureValueCents due after years years, discounted
// continuously at the nominal ratePerYear, FV·e^(−r·t), rounded half-even to the cent. years may be
// fractional; years = 0 returns the amount unchanged.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |futureValueCents| > 2⁵³; ErrRate if ratePerYear is NaN, ±Inf or ≤ −1;
// ErrNotFinite if years is NaN or ±Inf; ErrNegative if years < 0; ErrOverflow if the result exceeds
// ±2⁵³ cents (a negative rate).
//
// ID:      FIN-INT-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  G. Strang and E. Herman, Calculus Volume 1, OpenStax, 2016, §6.8 "Exponential Growth and Decay", https://openstax.org/books/calculus-volume-1/pages/6-8-exponential-growth-and-decay (Balance = Pe^(rt), solved for P)
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §3.2, Eq. 4 solved for PV
// Tags:    interest, time-value, present-value, discount, continuous-compounding, single-payment
// Units:   futureValueCents and the result in cents; ratePerYear decimal per year; years in years
// Example: $1,000,000 in 40 years at 5 % → $135,335.28; at 6 % → $90,717.95 (Calculus Volume 1, Example 6.43)
func PresentValueContinuous(futureValueCents int64, ratePerYear, years float64) (int64, error) {
	if err := checkCents("futureValueCents", futureValueCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerYear", ratePerYear); err != nil {
		return 0, err
	}
	if err := checkYears("years", years); err != nil {
		return 0, err
	}
	if futureValueCents == 0 || years == 0 {
		return futureValueCents, nil
	}
	return roundCents(float64(futureValueCents) * math.Exp(-ratePerYear*years))
}
