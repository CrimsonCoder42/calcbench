package series

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// AverageRateOfChangePerDay returns the average rate of change of a quantity between two time-stamped
// readings, per day of elapsed time: (to.Value − from.Value)/Δt, Δt = to.At − from.At in days of
// 86 400 s (the slope of the secant through the two points). Δt is computed from Unix seconds and
// nanoseconds, so spans beyond time.Duration's ±292 years are exact. A minimum span, a calendar month
// length or rounding are the caller's choices.
// Errors: ErrNotFinite when a value is NaN or ±Inf; ErrUnsorted when to.At is before from.At;
// ErrOutOfRange when to.At equals from.At (no elapsed time); ErrOverflow when the change or the rate is
// beyond the float64 range.
//
// ID:      NUM-SERIES-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Strang G, Herman E, et al., Calculus Volume 1, OpenStax, 2016, §3.4 "Derivatives as Rates of Change", https://openstax.org/books/calculus-volume-1/pages/3-4-derivatives-as-rates-of-change (read 2026-11-12) ("the average rate of change of the function f over that same interval is the ratio of the amount of change over that interval to the corresponding change in the x values", (f(a + h) − f(a))/h)
// Source:  BIPM, The International System of Units (SI Brochure), 9th ed., 2019, Table 8 (day: 1 d = 24 h = 86 400 s)
// Tags:    rate-of-change, velocity, secant-slope, trend, per-day, time-series
// Units:   values in any one unit; result in that unit per day
// Example: 50 at 2026-01-01T00:00Z, 53 at 2026-01-31T00:00Z → 3/30 = 0.1 per day (hand-checked)
func AverageRateOfChangePerDay(from, to Point) (float64, error) {
	if !reduce.AllFinite(from.Value) {
		return 0, fmt.Errorf("%w: from.Value = %v", ErrNotFinite, from.Value)
	}
	if !reduce.AllFinite(to.Value) {
		return 0, fmt.Errorf("%w: to.Value = %v", ErrNotFinite, to.Value)
	}
	if to.At.Before(from.At) {
		return 0, fmt.Errorf("%w: to.At = %v is before from.At = %v", ErrUnsorted, to.At, from.At)
	}
	if to.At.Equal(from.At) {
		return 0, fmt.Errorf("%w: to.At equals from.At, no elapsed time", ErrOutOfRange)
	}
	change := to.Value - from.Value
	if !reduce.AllFinite(change) {
		return 0, fmt.Errorf("%w: change from %v to %v", ErrOverflow, from.Value, to.Value)
	}
	rate := change / elapsedDays(from.At, to.At)
	if !reduce.AllFinite(rate) {
		return 0, fmt.Errorf("%w: rate of change", ErrOverflow)
	}
	return rate, nil
}
