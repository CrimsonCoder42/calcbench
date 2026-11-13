package geometry

import (
	"errors"
	"fmt"
	"math"
)

// ErrTooFewSides: a polygon needs at least 3 sides.
var ErrTooFewSides = errors.New("geometry: a polygon needs at least 3 sides")

// checkSides returns ErrTooFewSides wrapped with "sides" when sides < 3.
func checkSides(sides int) error {
	if sides < 3 {
		return fmt.Errorf("%w: sides", ErrTooFewSides)
	}
	return nil
}

// PolygonInteriorAngleSumRad returns the sum of the interior angles of a simple polygon with the given
// number of sides, in radians: (n − 2)π. n converts to float64 exactly for n < 2⁵³; above that it is
// rounded to the nearest float64.
// Errors: ErrTooFewSides when sides < 3.
//
// ID:      NUM-GEOM-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5 "Polygons", p. 318 ("In any k-gon, the sum of the angles equals 2(k − 2) right angles")
// Tags:    geometry, plane, angle
// Units:   sides is a count; result in radians
// Example: 5 → 3π (hand-derived)
func PolygonInteriorAngleSumRad(sides int) (float64, error) {
	if err := checkSides(sides); err != nil {
		return 0, err
	}
	return float64(sides-2) * math.Pi, nil
}

// RegularPolygonInteriorAngleRad returns one interior angle of a regular polygon with the given number of
// sides, in radians: (n − 2)π/n, computed as π − 2π/n (one rounding fewer).
// Errors: ErrTooFewSides when sides < 3.
//
// ID:      NUM-GEOM-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5.3 "Regular polygons", p. 324 (θ = ((k − 2)/k)·180°)
// Tags:    geometry, plane, angle
// Units:   sides is a count; result in radians
// Example: 6 → 2π/3 (hand-derived)
func RegularPolygonInteriorAngleRad(sides int) (float64, error) {
	if err := checkSides(sides); err != nil {
		return 0, err
	}
	return math.Pi - 2*math.Pi/float64(sides), nil
}

// RegularPolygonExteriorAngleRad returns one exterior angle of a regular polygon with the given number of
// sides, in radians: 2π/n, the supplement of the interior angle.
// Errors: ErrTooFewSides when sides < 3.
//
// ID:      NUM-GEOM-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5 "Polygons", p. 318 (the exterior angle: between one side and the extension of the adjacent side), with θ of §4.5.3 "Regular polygons", p. 324
// Tags:    geometry, plane, angle
// Units:   sides is a count; result in radians
// Example: 8 → π/4 (hand-derived)
func RegularPolygonExteriorAngleRad(sides int) (float64, error) {
	if err := checkSides(sides); err != nil {
		return 0, err
	}
	return 2 * math.Pi / float64(sides), nil
}
