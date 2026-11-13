package interest

// SimpleInterest returns the simple interest on principalCents at ratePerYear for yearFraction
// years, I = P·r·t, rounded half-even to the cent. Interest is not compounded. yearFraction is
// typically YearFraction30360 or YearFractionActual365 of the term. The payoff amount is P + I.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |principalCents| > 2⁵³; ErrRate if ratePerYear is NaN, ±Inf or ≤ −1;
// ErrNotFinite if yearFraction is NaN or ±Inf; ErrNegative if yearFraction < 0; ErrOverflow if the
// result exceeds ±2⁵³ cents.
//
// ID:      FIN-INT-035
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  D. Kirk, Contemporary Mathematics, OpenStax, 2023, §6.3 "Simple Interest", https://openstax.org/books/contemporary-mathematics/pages/6-3-simple-interest: I = P × r × t
// Tags:    interest, simple-interest, accrued-interest, short-term-loan
// Units:   principalCents and the result in cents; ratePerYear decimal per year; yearFraction years
// Example: $4,000 at 5.5 % for 4 years → $880.00 (Example 6.27); $11,995 at 14.9 % for 2 years → $3,574.51 (Example 6.28)
func SimpleInterest(principalCents int64, ratePerYear, yearFraction float64) (int64, error) {
	if err := checkCents("principalCents", principalCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerYear", ratePerYear); err != nil {
		return 0, err
	}
	if err := checkYears("yearFraction", yearFraction); err != nil {
		return 0, err
	}
	return roundCents(float64(principalCents) * ratePerYear * yearFraction)
}
