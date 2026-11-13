package geometry

import (
	"fmt"
	"math"
)

// checkRadii validates an outer and an inner radius: both finite and > 0, inner < outer.
func checkRadii(outerRadius, innerRadius float64) error {
	if err := checkPositive(outerRadius, "outerRadius"); err != nil {
		return err
	}
	if err := checkPositive(innerRadius, "innerRadius"); err != nil {
		return err
	}
	if innerRadius >= outerRadius {
		return fmt.Errorf("%w: innerRadius >= outerRadius", ErrOutOfRange)
	}
	return nil
}

// AnnulusArea returns the area of an annulus (ring) between concentric circles: π(R² − r²), computed as
// 2π(R − r)·(R/2 + r/2). The factored form has no cancellation for a thin ring and no R² to overflow.
// Errors: ErrNotFinite when a radius is NaN or ±Inf; ErrNonPositive when a radius ≤ 0; ErrOutOfRange
// when innerRadius ≥ outerRadius; ErrOverflow when the area is beyond float64.
//
// ID:      NUM-GEOM-036
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.6.6 "Circles", p. 334 (area πr², of the outer circle less the inner)
// Tags:    geometry, plane, circle, annulus, ring
// Units:   radii in one length unit; result in that unit squared
// Example: 5, 3 → 16π ≈ 50.265482457437 (hand-derived)
func AnnulusArea(outerRadius, innerRadius float64) (float64, error) {
	if err := checkRadii(outerRadius, innerRadius); err != nil {
		return 0, err
	}
	return checkResult(product(2*math.Pi, outerRadius-innerRadius, outerRadius/2+innerRadius/2, 1), "area")
}
