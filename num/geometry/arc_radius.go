package geometry

import "math"

// ArcRadiusFromChordSagitta returns the radius of a circular arc from its chord c and its sagitta s (the
// rise of the arc above the chord's midpoint): R = c²/(8s) + s/2, the source's c = √(4h(2R − h)) solved
// for R. Both terms are positive, so there is no cancellation; c²/(8s) is formed from the mantissas and
// exponents separately, so it overflows only when it is truly beyond float64. A sagitta above c/2 is a
// major arc, not an error.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when an input ≤ 0; ErrOverflow when
// the radius is beyond float64.
//
// ID:      NUM-GEOM-041
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.6.6 "Circles", p. 334 (c = √(4h(2R − h)), h the height of the segment; Figure 4.18, p. 335)
// Tags:    geometry, plane, circle, arc, sagitta, arch
// Units:   chord and sagitta in one length unit; result in that unit
// Example: chord 8, sagitta 2 → 5 (hand-derived: half-chord 4 and R − s = 3 make a 3-4-5 triangle)
func ArcRadiusFromChordSagitta(chord, sagitta float64) (float64, error) {
	if err := checkPositive(chord, "chord"); err != nil {
		return 0, err
	}
	if err := checkPositive(sagitta, "sagitta"); err != nil {
		return 0, err
	}
	mc, ec := math.Frexp(chord)
	ms, es := math.Frexp(sagitta)
	return checkResult(math.Ldexp(mc*mc/ms, 2*ec-es-3)+sagitta/2, "radius")
}
