package geometry

import (
	"errors"
	"fmt"
	"math"
)

// ErrNegative: a value that may be 0 is negative.
var ErrNegative = errors.New("geometry: value is negative")

// PolarToCartesian returns the point with polar coordinates (r, θ): (r cos θ, r sin θ), from one
// math.Sincos call. The result cannot overflow (|cos θ|, |sin θ| ≤ 1). θ is not reduced: for very large
// |θ| the accuracy of math.Sincos limits the result, which is not an error.
// Errors: ErrNotFinite when r or thetaRad is NaN or ±Inf; ErrNegative when r < 0.
//
// ID:      NUM-GEOM-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.1.4.1 "Relations between Cartesian and polar coordinates", p. 303 (x = r cos θ, y = r sin θ)
// Tags:    geometry, plane, coordinates, polar
// Units:   r in a length unit; thetaRad in radians; result coordinates in r's unit
// Example: r 2, θ π/6 → (√3, 1) (hand-derived)
func PolarToCartesian(r, thetaRad float64) (Point, error) {
	if err := checkFinite(r, "r"); err != nil {
		return Point{}, err
	}
	if err := checkFinite(thetaRad, "thetaRad"); err != nil {
		return Point{}, err
	}
	if r < 0 {
		return Point{}, fmt.Errorf("%w: r", ErrNegative)
	}
	sin, cos := math.Sincos(thetaRad)
	return Point{X: r * cos, Y: r * sin}, nil
}

// CartesianToPolar returns the polar coordinates of p: r = √(x²+y²) by math.Hypot and θ = atan2(y, x) in
// (−π, π] (the −π that atan2 gives for y = −0 is returned as π). The origin gives (0, 0) by the atan2
// convention. θ in [0, 2π) is the caller's adjustment.
// Errors: ErrNotFinite when a coordinate is NaN or ±Inf; ErrOverflow when r is beyond float64.
//
// ID:      NUM-GEOM-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.1.4.1 "Relations between Cartesian and polar coordinates", p. 303 (r = √(x²+y²), θ = tan⁻¹(y/x))
// Tags:    geometry, plane, coordinates, polar
// Units:   coordinates of p in one length unit; r in that unit; thetaRad in radians
// Example: (−1, 0) → (1, π) (hand-derived)
func CartesianToPolar(p Point) (r, thetaRad float64, err error) {
	if err := checkPoint(p, "p"); err != nil {
		return 0, 0, err
	}
	if r, err = checkResult(math.Hypot(p.X, p.Y), "r"); err != nil {
		return 0, 0, err
	}
	return r, halfOpenAngle(math.Atan2(p.Y, p.X)), nil
}
