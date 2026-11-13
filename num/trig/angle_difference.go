package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// AngleDifference returns the shortest signed turn from fromRad to toRad, in (−π, π]: positive is
// counter-clockwise (toRad ahead of fromRad), and a half turn gives +π.
// It reduces toRad − fromRad as NormalizeAngleSigned does. When both angles are huge (accumulated
// rotations), toRad − fromRad loses digits; normalise the inputs first.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOverflow when toRad − fromRad is beyond float64.
//
// ID:      NUM-TRIG-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  IEEE Std 1003.1-2024 (The Open Group Base Specifications Issue 8), remainder(), DESCRIPTION ("r=x−ny … n is the integral value nearest the exact value x/y"), https://pubs.opengroup.org/onlinepubs/9799919799/functions/remainder.html (read 2026-11-12)
// Tags:    trig, angle, difference, heading
// Units:   fromRad, toRad in radians; result in radians, in (−π, π]
// Example: from 350° to 10° → +20°; from 10° to 350° → −20° (hand-derived, angles passed in radians)
func AngleDifference(fromRad, toRad float64) (float64, error) {
	if !reduce.AllFinite(fromRad) {
		return 0, fmt.Errorf("%w: fromRad", ErrNotFinite)
	}
	if !reduce.AllFinite(toRad) {
		return 0, fmt.Errorf("%w: toRad", ErrNotFinite)
	}
	d := toRad - fromRad
	if math.IsInf(d, 0) {
		return 0, fmt.Errorf("%w: toRad - fromRad", ErrOverflow)
	}
	return remainderTwoPi(d), nil
}
