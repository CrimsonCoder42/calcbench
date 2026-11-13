package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// finiteResult returns r, or ErrOverflow when finite inputs produced ±Inf.
func finiteResult(r float64, what string) (float64, error) {
	if math.IsInf(r, 0) || math.IsNaN(r) {
		return 0, fmt.Errorf("%w: %s", ErrOverflow, what)
	}
	return r, nil
}

// PercentOf returns the part that is percent % of total: part = (percent/100)·total.
// Errors: ErrNotFinite for NaN or ±Inf input; ErrOverflow when the part is outside the float64 range.
//
// ID:      NUM-ALGEB-029
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kirk, Contemporary Mathematics, OpenStax, 2023, §6.1 "Understanding Percent", "Calculate the Total, Percent, or Part" (part = percent × total), Example 6.4
// Tags:    percent, percentage, part, proportion, scalar, algebra
// Units:   percent in percent (35 means 35 %); total and the result in one unit
// Example: 70 % of 3 500 = 2 450; 156 % of 720 = 1 123.2 (Kirk §6.1, Example 6.4)
func PercentOf(percent, total float64) (float64, error) {
	if !reduce.AllFinite(percent, total) {
		return 0, fmt.Errorf("%w: percent = %v, total = %v", ErrNotFinite, percent, total)
	}
	return finiteResult(percent/100*total, "percent of total")
}

// PercentFromPart returns the percent that part is of total: 100·part/total.
// Errors: ErrNotFinite for NaN or ±Inf input; ErrDegenerate when total = 0; ErrOverflow when the result
// is outside the float64 range.
//
// ID:      NUM-ALGEB-030
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kirk, Contemporary Mathematics, OpenStax, 2023, §6.1 "Understanding Percent", "Calculate the Total, Percent, or Part" (percent = part/total), Example 6.6
// Tags:    percent, percentage, what-percent, proportion, scalar, algebra
// Units:   part and total in one unit; the result in percent
// Example: 175 of 500 = 35 % (Kirk §6.1, Example 6.6)
func PercentFromPart(part, total float64) (float64, error) {
	if !reduce.AllFinite(part, total) {
		return 0, fmt.Errorf("%w: part = %v, total = %v", ErrNotFinite, part, total)
	}
	if total == 0 {
		return 0, fmt.Errorf("%w: total = 0", ErrDegenerate)
	}
	return finiteResult(100*(part/total), "part / total")
}

// TotalFromPercent returns the total of which part is percent %: total = 100·part/percent.
// Errors: ErrNotFinite for NaN or ±Inf input; ErrDegenerate when percent = 0; ErrOverflow when the result
// is outside the float64 range.
//
// ID:      NUM-ALGEB-031
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kirk, Contemporary Mathematics, OpenStax, 2023, §6.1 "Understanding Percent", "Calculate the Total, Percent, or Part" (total = part/percent), Example 6.5
// Tags:    percent, percentage, total, base, proportion, scalar, algebra
// Units:   part and the result in one unit; percent in percent
// Example: 70 is 35 % of 200 (Kirk §6.1, Example 6.5)
func TotalFromPercent(part, percent float64) (float64, error) {
	if !reduce.AllFinite(part, percent) {
		return 0, fmt.Errorf("%w: part = %v, percent = %v", ErrNotFinite, part, percent)
	}
	if percent == 0 {
		return 0, fmt.Errorf("%w: percent = 0", ErrDegenerate)
	}
	return finiteResult(part/percent*100, "part / percent")
}

// PercentChange returns the percent change from original to updated, 100·(updated − original)/|original|:
// positive for an increase, negative for a decrease.
// Errors: ErrNotFinite for NaN or ±Inf input; ErrDegenerate when original = 0; ErrOverflow when the result
// is outside the float64 range.
//
// ID:      NUM-ALGEB-032
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Marecek, Anthony-Smith and Honeycutt Mathis, Prealgebra 2e, OpenStax, 2020, §6.2 "Solve General Applications of Percent", "Find Percent Increase and Percent Decrease" (increase = new − original; percent increase = increase / original)
// Tags:    percent, percent-change, percent-increase, percent-decrease, relative-change, scalar, algebra
// Units:   original and updated in one unit; the result in percent
// Example: 40 → 50 is +25 %; 50 → 40 is −20 % (hand-checked)
func PercentChange(original, updated float64) (float64, error) {
	if !reduce.AllFinite(original, updated) {
		return 0, fmt.Errorf("%w: original = %v, updated = %v", ErrNotFinite, original, updated)
	}
	if original == 0 {
		return 0, fmt.Errorf("%w: original = 0", ErrDegenerate)
	}
	return finiteResult(100*((updated-original)/math.Abs(original)), "percent change")
}

// ApplyPercentChange returns original changed by percent %: original·(1 + percent/100). A decrease is a
// negative percent (−20 for a 20 % decrease).
// Errors: ErrNotFinite for NaN or ±Inf input; ErrOverflow when the result is outside the float64 range.
//
// ID:      NUM-ALGEB-033
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Marecek, Anthony-Smith and Honeycutt Mathis, Prealgebra 2e, OpenStax, 2020, §6.2 "Solve General Applications of Percent", "Find Percent Increase and Percent Decrease" (new = original + percent × original)
// Tags:    percent, percent-increase, percent-decrease, markup, scalar, algebra
// Units:   original and the result in one unit; percent in percent
// Example: 500 increased by 10 % = 550; 500 decreased by 10 % (percent −10) = 450 (hand-checked)
func ApplyPercentChange(original, percent float64) (float64, error) {
	if !reduce.AllFinite(original, percent) {
		return 0, fmt.Errorf("%w: original = %v, percent = %v", ErrNotFinite, original, percent)
	}
	return finiteResult(original+original*(percent/100), "changed value")
}
