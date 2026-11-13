package geometry

import (
	"fmt"
	"math"
)

// checkSemiAxes validates the two semi-axes of an ellipse chord: finite and > 0.
func checkSemiAxes(semiAxisAlong, semiAxisAcross float64) error {
	if err := checkPositive(semiAxisAlong, "semiAxisAlong"); err != nil {
		return err
	}
	return checkPositive(semiAxisAcross, "semiAxisAcross")
}

// EllipseChordFromCenter returns the length of the chord of an ellipse that runs parallel to one axis at a
// given offset from the centre. From x²/a² + y²/b² = 1 with a the semi-axis parallel to the chord
// (semiAxisAlong) and b the other (semiAxisAcross), the chord at offset y is L = 2a·√(1 − y²/b²), computed
// as 2a·(√(b−y)·√(b+y))/b: no cancellation near the vertex, and nothing to overflow. A horizontal or a
// vertical chord is only a choice of which semi-axis is "along". Offset 0 gives the full axis 2a; offset b
// gives 0.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrNegative
// when offsetFromCenter < 0; ErrOutOfRange when offsetFromCenter > semiAxisAcross; ErrOverflow when the
// chord is beyond float64.
//
// ID:      NUM-GEOM-044
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.6.3 "Additional properties of ellipses", p. 330 (the ellipse x²/a² + y²/b² = 1)
// Tags:    geometry, plane, ellipse, conic, chord
// Units:   semi-axes and offsetFromCenter in one length unit; result in that unit
// Example: a 5, b 3, offset 1.8 → 8 (hand-derived: √(1 − 0.36) = 0.8)
func EllipseChordFromCenter(semiAxisAlong, semiAxisAcross, offsetFromCenter float64) (float64, error) {
	if err := checkSemiAxes(semiAxisAlong, semiAxisAcross); err != nil {
		return 0, err
	}
	if err := checkNonNegative(offsetFromCenter, "offsetFromCenter"); err != nil {
		return 0, err
	}
	if offsetFromCenter > semiAxisAcross {
		return 0, fmt.Errorf("%w: offsetFromCenter > semiAxisAcross", ErrOutOfRange)
	}
	ratio := legFromHypot(semiAxisAcross, offsetFromCenter) / semiAxisAcross
	return checkResult(product(2, semiAxisAlong, ratio, 1), "chord")
}

// EllipseChordFromEdge returns the length of the chord of an ellipse that runs parallel to one axis at a
// given depth h measured inward from the vertex on the other axis: with y = b − h in EllipseChordFromCenter,
// L = (2a/b)·√(h·(2b − h)), computed from h and 2b − h directly so a chord near the vertex keeps every digit
// (b − y would cancel). Depth 0 and depth 2b give 0; depth b gives the full axis 2a.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrNegative
// when depthFromEdge < 0; ErrOutOfRange when depthFromEdge > 2·semiAxisAcross; ErrOverflow when the chord
// is beyond float64.
//
// ID:      NUM-GEOM-045
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.6.3 "Additional properties of ellipses", p. 330 (the ellipse x²/a² + y²/b² = 1)
// Tags:    geometry, plane, ellipse, conic, chord
// Units:   semi-axes and depthFromEdge in one length unit; result in that unit
// Example: a 5, b 3, depth 1.2 → 8 (hand-derived; the same chord as NUM-GEOM-044's example)
func EllipseChordFromEdge(semiAxisAlong, semiAxisAcross, depthFromEdge float64) (float64, error) {
	if err := checkSemiAxes(semiAxisAlong, semiAxisAcross); err != nil {
		return 0, err
	}
	if err := checkNonNegative(depthFromEdge, "depthFromEdge"); err != nil {
		return 0, err
	}
	b, h := semiAxisAcross, depthFromEdge
	if h/2 > b {
		return 0, fmt.Errorf("%w: depthFromEdge > 2·semiAxisAcross", ErrOutOfRange)
	}
	var root float64
	if w := 2*b - h; !math.IsInf(w, 0) {
		root = math.Sqrt(h) * math.Sqrt(w)
	} else {
		root = math.Sqrt(h) * (math.Sqrt2 * math.Sqrt(b-h/2))
	}
	return checkResult(product(2, semiAxisAlong, root/b, 1), "chord")
}
