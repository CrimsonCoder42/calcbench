package stream

import (
	"errors"
	"fmt"
	"math"
)

// ErrAlpha: the smoothing constant is not in (0, 1].
var ErrAlpha = errors.New("stream: smoothing constant alpha must be in (0, 1]")

// EWMAStream is the single exponential smoothing (exponentially weighted moving average) of the values
// added so far, kept one value at a time: the first value is its own average, and each later value x
// gives S ← α·x + (1 − α)·S. It is the streaming form of ALG-WINDOW-014 ExponentialMovingAverage, with
// the same recurrence, the same initialisation and the same floating-point expression, so after the
// i-th Add Value equals result i of ExponentialMovingAverage on the same values. α is a parameter (no
// default): α = 1 follows the values; small α smooths more. The zero value is not usable: its Add
// returns ErrUninitialized. Not safe for concurrent use. No Merge: an EWMA of a later part of a stream
// cannot be combined without knowing that part's length and first value.
//
// Methods: Add(x) adds a value and returns the new average. Value returns the average, or (0, false)
// before the first value. Alpha returns α. Reset forgets every value and keeps α.
//
// ID:      ALG-STREAM-009
// Time:    O(1) per Add and per reader
// Space:   O(1), no allocation after NewEWMAStream
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §6.4.3.1 "Single Exponential Smoothing", https://www.itl.nist.gov/div898/handbook/pmc/section4/pmc431.htm (read 2026-11-12) (Sₜ = α·yₜ₋₁ + (1 − α)·Sₜ₋₁, S₂ = y₁)
// Tags:    stream, exponential-moving-average, ema, ewma, exponential-smoothing, smoothing, time-series
// Units:   values x in any unit U; Value in U; alpha unitless, in (0, 1]
// Example: alpha 0.5, Add 2, 4, 8 → 2, 3, 5.5
type EWMAStream struct {
	alpha, keep float64 // keep = 1 − alpha, formed once as ExponentialMovingAverage does
	s           float64
	started     bool
}

// NewEWMAStream returns an empty EWMAStream with smoothing constant alpha.
// Errors: ErrAlpha when alpha is not in (0, 1] (NaN included).
func NewEWMAStream(alpha float64) (*EWMAStream, error) {
	if !(alpha > 0 && alpha <= 1) {
		return nil, fmt.Errorf("%w: alpha = %v", ErrAlpha, alpha)
	}
	return &EWMAStream{alpha: alpha, keep: 1 - alpha}, nil
}

// Add adds x and returns the new average. O(1).
// Errors: ErrNotFinite when x is NaN or ±Inf; ErrOverflow when the average rounds beyond the float64
// range (possible only for values within an ulp of it); ErrUninitialized on the zero value. On error
// nothing changes.
func (e *EWMAStream) Add(x float64) (float64, error) {
	if e.alpha == 0 {
		return 0, ErrUninitialized
	}
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	if !e.started {
		e.s, e.started = x, true
		return x, nil
	}
	// The expression of ExponentialMovingAverage, term for term, so the two agree value for value.
	s := e.alpha*x + e.keep*e.s
	if math.IsInf(s, 0) {
		return 0, fmt.Errorf("%w: x = %v", ErrOverflow, x)
	}
	e.s = s
	return s, nil
}

// Value returns the current average, or (0, false) before the first value. O(1).
func (e *EWMAStream) Value() (float64, bool) {
	if !e.started {
		return 0, false
	}
	return e.s, true
}

// Alpha returns the smoothing constant (0 for the zero value). O(1).
func (e *EWMAStream) Alpha() float64 { return e.alpha }

// Reset forgets every value and keeps alpha. O(1).
func (e *EWMAStream) Reset() { e.s, e.started = 0, false }
