package personal

import (
	"fmt"
	"math"
)

// HourlyRateFromWeeklyPay returns the regular hourly rate of a weekly salary: weeklyPayCents divided by
// the hours the salary is meant to cover, rounded half-even to the cent. For a salary paid over another
// period, convert it first with WeeklyPayEquivalent.
// Errors: ErrNegative if weeklyPayCents < 0; ErrOutOfRange if weeklyPayCents > 2⁵³; ErrNotFinite if
// hoursPerWeek is NaN or ±Inf; ErrDomain if hoursPerWeek ≤ 0 or above 168 (the hours in a week);
// ErrOutOfRange if the rate would exceed 2⁵³ cents (a tiny hoursPerWeek).
//
// ID:      FIN-PERSON-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  29 CFR 778.113(a) "Weekly salary" ("the regular hourly rate of pay ... is computed by dividing the salary by the number of hours which the salary is intended to compensate"), https://www.law.cornell.edu/cfr/text/29/778.113 (read 2026-11-13)
// Tags:    personal, salary, hourly-rate, regular-rate, payroll, flsa
// Units:   weeklyPayCents in cents per week; hoursPerWeek hours; result cents per hour
// Example: hand-checked from the rule: $2,000.00 a week for 40 hours → $50.00 an hour
func HourlyRateFromWeeklyPay(weeklyPayCents int64, hoursPerWeek float64) (int64, error) {
	if weeklyPayCents < 0 {
		return 0, fmt.Errorf("%w: weeklyPayCents = %d", ErrNegative, weeklyPayCents)
	}
	if weeklyPayCents > maxExactCents {
		return 0, fmt.Errorf("%w: weeklyPayCents = %d is above 2^53", ErrOutOfRange, weeklyPayCents)
	}
	if math.IsNaN(hoursPerWeek) || math.IsInf(hoursPerWeek, 0) {
		return 0, fmt.Errorf("%w: hoursPerWeek = %v", ErrNotFinite, hoursPerWeek)
	}
	if hoursPerWeek <= 0 || hoursPerWeek > 168 {
		return 0, fmt.Errorf("%w: hoursPerWeek = %v, want (0, 168]", ErrDomain, hoursPerWeek)
	}
	r := math.RoundToEven(float64(weeklyPayCents) / hoursPerWeek)
	if r > maxExactCents {
		return 0, fmt.Errorf("%w: hourly rate %v cents is above 2^53 (hoursPerWeek = %v)", ErrOutOfRange, r, hoursPerWeek)
	}
	return int64(r), nil
}
