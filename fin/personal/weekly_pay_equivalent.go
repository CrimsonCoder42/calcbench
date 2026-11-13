package personal

import (
	"errors"
	"fmt"
)

// Sentinel errors of package personal. Every error a function returns wraps one of these with the name
// of the offending parameter.
var (
	// ErrOutOfRange: a money input is outside 0…2⁵³ cents, or a count is out of range.
	ErrOutOfRange = errors.New("personal: input out of range")
	// ErrNegative: an input that must not be negative is negative.
	ErrNegative = errors.New("personal: negative input")
	// ErrNotFinite: an input is NaN or ±Inf.
	ErrNotFinite = errors.New("personal: value is NaN or infinite")
	// ErrDomain: inputs outside the formula's domain.
	ErrDomain = errors.New("personal: inputs outside the formula's domain")
)

// maxExactCents is the largest money magnitude, in cents, that float64 holds exactly (2⁵³).
const maxExactCents = 1 << 53

// WeeklyPayEquivalent returns the weekly wage equivalent to payCents paid periodsPerYear times a year:
// payCents × periodsPerYear ÷ 52, rounded half-even to the cent in exact integer arithmetic. The rule
// is the Fair Labor Standards Act's: a monthly salary × 12 ÷ 52, a semimonthly salary × 24 ÷ 52; the same
// rule gives 26 for biweekly, 4 for quarterly and 1 for annual pay. Multiply the result by 52 and divide
// by another period count to convert between any two pay periods.
// Errors: ErrNegative if payCents < 0; ErrOutOfRange if payCents > 2⁵³ or periodsPerYear is outside 1…366.
//
// ID:      FIN-PERSON-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  29 CFR 778.113(b) "Salary for periods other than workweek" ("A monthly salary is subject to translation to its equivalent weekly wage by multiplying by 12 (the number of months) and dividing by 52 (the number of weeks). A semimonthly salary is translated into its equivalent weekly wage by multiplying by 24 and dividing by 52."), https://www.law.cornell.edu/cfr/text/29/778.113 (read 2026-11-13)
// Tags:    personal, salary, pay-period, weekly-wage, payroll, flsa
// Units:   payCents in cents per pay period; periodsPerYear pay periods per year; result cents per week
// Example: hand-checked from the rule: $8,666.67 a month × 12 ÷ 52 = $2,000.00 a week
func WeeklyPayEquivalent(payCents int64, periodsPerYear int) (int64, error) {
	if payCents < 0 {
		return 0, fmt.Errorf("%w: payCents = %d", ErrNegative, payCents)
	}
	if payCents > maxExactCents {
		return 0, fmt.Errorf("%w: payCents = %d is above 2^53", ErrOutOfRange, payCents)
	}
	if periodsPerYear < 1 || periodsPerYear > 366 {
		return 0, fmt.Errorf("%w: periodsPerYear = %d, want 1…366", ErrOutOfRange, periodsPerYear)
	}
	n := payCents * int64(periodsPerYear) // ≤ 2⁵³ × 366 < 2⁶³
	q, r := n/52, n%52
	if 2*r > 52 || (2*r == 52 && q%2 == 1) {
		q++
	}
	return q, nil
}
