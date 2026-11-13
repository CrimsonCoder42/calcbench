package logs

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrOutOfRange: an input is outside the values the function accepts.
	ErrOutOfRange = errors.New("logs: out of range")
	// ErrOverflow: finite inputs give a result outside the float64 range.
	ErrOverflow = errors.New("logs: result outside the float64 range")
)

func finite(xs ...float64) bool {
	for _, x := range xs {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return false
		}
	}
	return true
}

// ExponentialDecayRemaining returns the quantity left after elapsed time from initial, decaying with the
// given half-life: initial·2^(−elapsed/halfLife), the model A(t) = A₀e^{kt} with k = −ln 2/halfLife.
// A negative elapsed gives the quantity at that earlier time.
// Errors: ErrNotFinite for NaN or ±Inf input; ErrOutOfRange when halfLife ≤ 0; ErrOverflow when the
// result is outside the float64 range.
//
// ID:      NUM-LOGS-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Abramson, College Algebra 2e, OpenStax, 2021, §6.7 "Exponential and Logarithmic Models", "Half-Life" (A = A₀e^{kt}, half-life t = −ln 2/k), Example 2
// Tags:    exponential-decay, half-life, decay, radioactive-decay, exponential, scalar
// Units:   initial and the result in one quantity unit; elapsed and halfLife in one time unit
// Example: 100 units, half-life 50, after 100 → 25 (two half-lives, hand-checked)
func ExponentialDecayRemaining(initial, elapsed, halfLife float64) (float64, error) {
	if !finite(initial, elapsed, halfLife) {
		return 0, fmt.Errorf("%w: initial = %v, elapsed = %v, halfLife = %v", ErrNotFinite, initial, elapsed, halfLife)
	}
	if halfLife <= 0 {
		return 0, fmt.Errorf("%w: halfLife = %v, want > 0", ErrOutOfRange, halfLife)
	}
	r := initial * math.Exp2(-elapsed/halfLife)
	if math.IsInf(r, 0) || math.IsNaN(r) {
		return 0, fmt.Errorf("%w: remaining after %v half-lives", ErrOverflow, elapsed/halfLife)
	}
	return r, nil
}

// ExponentialDecayInitial returns the starting quantity that decays to remaining after elapsed time with
// the given half-life: remaining·2^(elapsed/halfLife).
// Errors: ErrNotFinite for NaN or ±Inf input; ErrOutOfRange when halfLife ≤ 0; ErrOverflow when the
// result is outside the float64 range.
//
// ID:      NUM-LOGS-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Abramson, College Algebra 2e, OpenStax, 2021, §6.7 "Exponential and Logarithmic Models", "Half-Life" (A = A₀e^{kt} solved for A₀)
// Tags:    exponential-decay, half-life, decay, initial-quantity, exponential, scalar
// Units:   remaining and the result in one quantity unit; elapsed and halfLife in one time unit
// Example: 25 left after 100 with half-life 50 → 100 at the start (hand-checked)
func ExponentialDecayInitial(remaining, elapsed, halfLife float64) (float64, error) {
	if !finite(remaining, elapsed, halfLife) {
		return 0, fmt.Errorf("%w: remaining = %v, elapsed = %v, halfLife = %v", ErrNotFinite, remaining, elapsed, halfLife)
	}
	if halfLife <= 0 {
		return 0, fmt.Errorf("%w: halfLife = %v, want > 0", ErrOutOfRange, halfLife)
	}
	r := remaining * math.Exp2(elapsed/halfLife)
	if math.IsInf(r, 0) || math.IsNaN(r) {
		return 0, fmt.Errorf("%w: initial before %v half-lives", ErrOverflow, elapsed/halfLife)
	}
	return r, nil
}

// ExponentialDecayHalfLife returns the half-life of an exponential decay from initial to remaining over
// elapsed time: elapsed·ln 2 / ln(initial/remaining).
// Errors: ErrNotFinite for NaN or ±Inf input; ErrOutOfRange unless initial > remaining > 0 and elapsed > 0
// (no decay, growth or an empty quantity has no finite positive half-life); ErrOverflow when the result
// is outside the float64 range.
//
// ID:      NUM-LOGS-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Abramson, College Algebra 2e, OpenStax, 2021, §6.7 "Exponential and Logarithmic Models", "Half-Life" (k = ln(A/A₀)/t, half-life = −ln 2/k)
// Tags:    exponential-decay, half-life, decay, exponential, logarithm, scalar
// Units:   initial and remaining in one quantity unit; elapsed and the result in one time unit
// Example: 100 → 10 over 50 → 50·ln 2/ln 10 = 15.0514997… (hand-checked)
func ExponentialDecayHalfLife(initial, remaining, elapsed float64) (float64, error) {
	if !finite(initial, remaining, elapsed) {
		return 0, fmt.Errorf("%w: initial = %v, remaining = %v, elapsed = %v", ErrNotFinite, initial, remaining, elapsed)
	}
	if !(initial > remaining && remaining > 0) || elapsed <= 0 {
		return 0, fmt.Errorf("%w: initial = %v, remaining = %v, elapsed = %v, want initial > remaining > 0 and elapsed > 0",
			ErrOutOfRange, initial, remaining, elapsed)
	}
	r := elapsed * (math.Ln2 / decayLog(initial, remaining))
	if math.IsInf(r, 0) || math.IsNaN(r) {
		return 0, fmt.Errorf("%w: half-life", ErrOverflow)
	}
	return r, nil
}

// decayLog returns ln(initial/remaining) for initial > remaining > 0, without forming a ratio that
// overflows or underflows.
func decayLog(initial, remaining float64) float64 {
	q := initial / remaining
	if math.IsInf(q, 0) || q == 0 {
		return math.Log(initial) - math.Log(remaining)
	}
	return math.Log(q)
}

// ExponentialDecayElapsed returns the time for initial to decay to remaining with the given half-life:
// halfLife·log₂(initial/remaining). A remaining above initial gives a negative time (the earlier instant).
// Errors: ErrNotFinite for NaN or ±Inf input; ErrOutOfRange when initial ≤ 0, remaining ≤ 0 or
// halfLife ≤ 0; ErrOverflow when the result is outside the float64 range.
//
// ID:      NUM-LOGS-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Abramson, College Algebra 2e, OpenStax, 2021, §6.7 "Exponential and Logarithmic Models", Example 3 "Finding the Age of a Bone" (t = ln(A/A₀)/k)
// Tags:    exponential-decay, half-life, decay, elapsed-time, dating, logarithm, scalar
// Units:   initial and remaining in one quantity unit; halfLife and the result in one time unit
// Example: 20 % left, half-life 5 730 years → 5 730·log₂ 5 = 13 304.6 years (Abramson Example 3 prints ≈ 13 301 from a rounded k)
func ExponentialDecayElapsed(initial, remaining, halfLife float64) (float64, error) {
	if !finite(initial, remaining, halfLife) {
		return 0, fmt.Errorf("%w: initial = %v, remaining = %v, halfLife = %v", ErrNotFinite, initial, remaining, halfLife)
	}
	if initial <= 0 || remaining <= 0 || halfLife <= 0 {
		return 0, fmt.Errorf("%w: initial = %v, remaining = %v, halfLife = %v, want all > 0", ErrOutOfRange, initial, remaining, halfLife)
	}
	var lg float64
	if initial >= remaining {
		lg = decayLog(initial, remaining)
	} else {
		lg = -decayLog(remaining, initial)
	}
	r := halfLife * (lg / math.Ln2)
	if math.IsInf(r, 0) || math.IsNaN(r) {
		return 0, fmt.Errorf("%w: elapsed time", ErrOverflow)
	}
	return r, nil
}

// DecayConstantFromHalfLife returns the decay constant λ = ln 2/halfLife (the k of A₀e^{−λt} with the sign
// made positive).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when halfLife ≤ 0; ErrOverflow when halfLife is so
// small that λ is outside the float64 range.
//
// ID:      NUM-LOGS-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Abramson, College Algebra 2e, OpenStax, 2021, §6.7 "Exponential and Logarithmic Models", "Half-Life" (half-life = −ln 2/k, so |k| = ln 2/half-life)
// Source:  Ling, Sanny and Moebs, University Physics Volume 3, OpenStax, 2016, §10.3 "Radioactive Decay", eq. 10.15 (λ = 0.693/T½)
// Tags:    exponential-decay, half-life, decay-constant, rate-constant, scalar
// Units:   halfLife in a time unit; the result per that time unit
// Example: half-life 5 730 years → ln(0.5)/5 730 in magnitude = 1.20968…×10⁻⁴ per year (Abramson Example 2)
func DecayConstantFromHalfLife(halfLife float64) (float64, error) {
	if !finite(halfLife) {
		return 0, fmt.Errorf("%w: halfLife = %v", ErrNotFinite, halfLife)
	}
	if halfLife <= 0 {
		return 0, fmt.Errorf("%w: halfLife = %v, want > 0", ErrOutOfRange, halfLife)
	}
	r := math.Ln2 / halfLife
	if math.IsInf(r, 0) {
		return 0, fmt.Errorf("%w: halfLife = %v", ErrOverflow, halfLife)
	}
	return r, nil
}

// MeanLifetimeFromHalfLife returns the mean lifetime τ = 1/λ = halfLife/ln 2 of an exponential decay: the
// time for the quantity to fall by a factor of e.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when halfLife ≤ 0; ErrOverflow when the result is
// outside the float64 range.
//
// ID:      NUM-LOGS-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ling, Sanny and Moebs, University Physics Volume 3, OpenStax, 2016, §10.3 "Radioactive Decay", eq. 10.15 (λ = 0.693/T½ = ln 2/T½) and eq. 10.16 (lifetime = 1/λ)
// Tags:    exponential-decay, half-life, mean-lifetime, time-constant, scalar
// Units:   halfLife and the result in one time unit
// Example: half-life ln 2 → τ = 1; half-life 50 → 72.1347520… (hand-checked)
func MeanLifetimeFromHalfLife(halfLife float64) (float64, error) {
	if !finite(halfLife) {
		return 0, fmt.Errorf("%w: halfLife = %v", ErrNotFinite, halfLife)
	}
	if halfLife <= 0 {
		return 0, fmt.Errorf("%w: halfLife = %v, want > 0", ErrOutOfRange, halfLife)
	}
	r := halfLife / math.Ln2
	if math.IsInf(r, 0) {
		return 0, fmt.Errorf("%w: halfLife = %v", ErrOverflow, halfLife)
	}
	return r, nil
}
