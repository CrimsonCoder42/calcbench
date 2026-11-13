package geometry

import (
	"fmt"
	"math"
)

// checkHelix validates a helix radius and a number of turns: both finite and > 0 (fractional turns allowed).
func checkHelix(radius, turns float64) error {
	if err := checkPositive(radius, "radius"); err != nil {
		return err
	}
	return checkPositive(turns, "turns")
}

// HelixArcLength returns the length of a circular helix of the given radius and pitch (the rise per turn)
// over the given number of turns: L = n·√((2πr)² + p²). The helix (r cos t, r sin t, (p/2π)t) has speed
// √(r² + (p/2π)²), so n turns (t = 2πn) unroll to the hypotenuse of a right triangle with legs 2πrn and pn.
// It is computed as 4πn·hypot(r/2, p/4π), which cannot overflow before the length does. A pitch of 0 is a
// circle (n·CircleCircumference).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the length is beyond float64.
//
// ID:      NUM-GEOM-081
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.20.1.3 "Example", p. 375 (the circular helix (a cos t, a sin t, bt): ds/dt = √(a² + b²), arc length s = t√(a² + b²))
// Tags:    geometry, space-curve, helix, arc-length, garland
// Units:   radius and pitch in one length unit; turns is a count (fractional allowed); result in that unit
// Example: 1, 2, 3 → 3·√(4π² + 4) ≈ 19.781449856854 (hand-derived)
func HelixArcLength(radius, pitch, turns float64) (float64, error) {
	if err := checkHelix(radius, turns); err != nil {
		return 0, err
	}
	if err := checkPositive(pitch, "pitch"); err != nil {
		return 0, err
	}
	return checkResult(product(4*math.Pi, turns, math.Hypot(radius/2, pitch/(4*math.Pi)), 1), "length")
}

// HelixHeight returns the total height of a circular helix of the given radius whose given arc length
// makes the given number of turns: the other leg of the unrolled right triangle,
// h = √(L² − (2πrn)²), computed as √(L − 2πrn)·√(L + 2πrn): L − 2πrn is exact when the two are close, so
// nothing is lost beyond the input's own rounding, and no square overflows. (For a tightly wound helix the
// height is inherently sensitive to L: a relative change δ in L moves h by about δ·(L/h)².) An arc length
// of exactly 2πrn is a flat coil and gives 0. The pitch is h/n.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOutOfRange when
// arcLength < 2π·radius·turns (too short to make that many turns).
//
// ID:      NUM-GEOM-082
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.20.1.3 "Example", p. 375 (arc length s = t√(a² + b²) of the circular helix, solved for the height bt)
// Tags:    geometry, space-curve, helix, pitch
// Units:   radius and arcLength in one length unit; turns is a count (fractional allowed); result in that unit
// Example: 1, 19.781449856854, 3 → 6 (hand-derived; the inverse of NUM-GEOM-081's example)
func HelixHeight(radius, arcLength, turns float64) (float64, error) {
	if err := checkHelix(radius, turns); err != nil {
		return 0, err
	}
	if err := checkPositive(arcLength, "arcLength"); err != nil {
		return 0, err
	}
	around := product(2*math.Pi, radius, turns, 1)
	if arcLength < around {
		return 0, fmt.Errorf("%w: arcLength < 2π·radius·turns", ErrOutOfRange)
	}
	return legFromHypot(arcLength, around), nil
}
