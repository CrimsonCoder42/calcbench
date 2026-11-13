package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// twoPi is float64(2π) = 2·math.Pi, exact as a doubling. Reductions are modulo this value, not true 2π.
const twoPi = 2 * math.Pi

// NormalizeAngle returns angleRad reduced to [0, 2π): angleRad − 2π·⌊angleRad/(2π)⌋.
// The reduction is exact relative to float64(2π) (math.Mod is exact); because 2π itself is not
// representable, for |angleRad| far above 2π the result differs from a reduction by true 2π by about
// |angleRad|·1.2e-16. A negative angle that reduces to within half an ulp of 2π gives 0, and −0 gives +0.
// Errors: ErrNotFinite when angleRad is NaN or ±Inf.
//
// ID:      NUM-TRIG-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  IEEE Std 1003.1-2024 (The Open Group Base Specifications Issue 8), fmod(), DESCRIPTION ("x−i*y … the result has the same sign as x and magnitude less than the magnitude of y"), https://pubs.opengroup.org/onlinepubs/9799919799/functions/fmod.html (read 2026-11-12)
// Tags:    trig, angle, normalize, wrap
// Units:   angleRad in radians; result in radians, in [0, 2π)
// Example: 7π/2 → 3π/2; −π/2 → 3π/2; 2π → 0 (hand-derived from the fmod definition)
func NormalizeAngle(angleRad float64) (float64, error) {
	if !reduce.AllFinite(angleRad) {
		return 0, fmt.Errorf("%w: angleRad", ErrNotFinite)
	}
	r := math.Mod(angleRad, twoPi)
	if r < 0 {
		r += twoPi
		if r >= twoPi {
			r = 0
		}
	}
	if r == 0 {
		return 0, nil // also turns −0 into +0
	}
	return r, nil
}
