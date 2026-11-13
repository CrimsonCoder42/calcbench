package geometry

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// Errors shared by every function of the package. They are wrapped with the offending parameter or the
// quantity that overflowed; test with errors.Is. The other sentinels are declared beside the function
// that first needs them: ErrDegenerate and ErrVertical (slope.go), ErrNegative (polar.go), ErrOutOfRange
// (complementary_angle.go), ErrTooFewSides (polygon_angles.go), ErrNonPositive and ErrNotTriangle
// (triangle_perimeter.go), ErrNoConvergence (ellipse_perimeter.go).
var (
	// ErrNotFinite: an input is NaN or ±Inf.
	ErrNotFinite = errors.New("geometry: value is NaN or infinite")
	// ErrOverflow: the true result of finite inputs is beyond the float64 range.
	ErrOverflow = errors.New("geometry: result exceeds the float64 range")
)

// checkFinite returns ErrNotFinite wrapped with name when x is NaN or ±Inf.
func checkFinite(x float64, name string) error {
	if !reduce.AllFinite(x) {
		return fmt.Errorf("%w: %s", ErrNotFinite, name)
	}
	return nil
}

// checkPoint returns ErrNotFinite wrapped with name when a coordinate of p is NaN or ±Inf.
func checkPoint(p Point, name string) error {
	if !reduce.AllFinite(p.X) || !reduce.AllFinite(p.Y) {
		return fmt.Errorf("%w: %s", ErrNotFinite, name)
	}
	return nil
}

// checkResult returns v, or ErrOverflow wrapped with name when v is ±Inf.
func checkResult(v float64, name string) (float64, error) {
	if math.IsInf(v, 0) {
		return 0, fmt.Errorf("%w: %s", ErrOverflow, name)
	}
	return v, nil
}

// PointDistance returns the distance between points p and q, √((x₂−x₁)² + (y₂−y₁)²), computed with
// math.Hypot so the squares never overflow or underflow. A coordinate difference that overflows means
// the distance does too.
// Errors: ErrNotFinite when a coordinate is NaN or ±Inf; ErrOverflow when the distance is beyond float64.
//
// ID:      NUM-GEOM-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.4.2 "Distances", eq. (4.4.5), p. 316
// Tags:    geometry, plane, coordinates, distance
// Units:   coordinates of p and q in one length unit; result in that unit
// Example: (0, 0), (3, 4) → 5 (hand-derived from eq. (4.4.5))
func PointDistance(p, q Point) (float64, error) {
	if err := checkPoint(p, "p"); err != nil {
		return 0, err
	}
	if err := checkPoint(q, "q"); err != nil {
		return 0, err
	}
	return checkResult(math.Hypot(q.X-p.X, q.Y-p.Y), "distance")
}
