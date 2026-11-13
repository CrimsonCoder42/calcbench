package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// NormalizeAngleSigned returns angleRad reduced to (−π, π]: angleRad − 2π·n, n the integer nearest
// angleRad/(2π). The reduction is exact relative to float64(2π) (math.Remainder is exact); the same
// large-|angleRad| caveat as NormalizeAngle applies. A remainder of exactly −π (a tie) is returned as +π.
// Errors: ErrNotFinite when angleRad is NaN or ±Inf.
//
// ID:      NUM-TRIG-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  IEEE Std 1003.1-2024 (The Open Group Base Specifications Issue 8), remainder(), DESCRIPTION ("r=x−ny … n is the integral value nearest the exact value x/y"), https://pubs.opengroup.org/onlinepubs/9799919799/functions/remainder.html (read 2026-11-12)
// Tags:    trig, angle, normalize, wrap
// Units:   angleRad in radians; result in radians, in (−π, π]
// Example: 7π/2 → −π/2; 3π/2 → −π/2; −π → π (hand-derived from the remainder definition)
func NormalizeAngleSigned(angleRad float64) (float64, error) {
	if !reduce.AllFinite(angleRad) {
		return 0, fmt.Errorf("%w: angleRad", ErrNotFinite)
	}
	return remainderTwoPi(angleRad), nil
}

// remainderTwoPi reduces a finite x to (−π, π] by math.Remainder with the −π tie moved to +π.
func remainderTwoPi(x float64) float64 {
	r := math.Remainder(x, twoPi)
	if r <= -math.Pi {
		r = math.Pi
	}
	return r
}
