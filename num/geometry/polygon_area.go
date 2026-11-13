package geometry

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// verticesScale checks every vertex for NaN and ±Inf and returns the binary exponent e of the largest
// coordinate magnitude (0 when every coordinate is 0). Scaling by 2⁻ᵉ is exact for normal numbers and
// brings every coordinate into [−1, 1], so no difference, product or sum below can overflow.
func verticesScale(vertices []Point) (int, error) {
	m := 0.0
	for i, v := range vertices {
		if !reduce.AllFinite(v.X) || !reduce.AllFinite(v.Y) {
			return 0, fmt.Errorf("%w: vertices[%d]", ErrNotFinite, i)
		}
		m = math.Max(m, math.Max(math.Abs(v.X), math.Abs(v.Y)))
	}
	_, e := math.Frexp(m)
	return e, nil
}

// PolygonArea returns the area of a simple polygon from its vertices in order (either orientation; the
// first vertex is not repeated at the end), by the shoelace (surveyor's) formula
// A = ½|Σ (xᵢyᵢ₊₁ − xᵢ₊₁yᵢ)|, indices mod n. Every vertex is first translated by vertex 0, which removes
// the cancellation when the polygon is far from the origin, and scaled by a power of two (exact), which
// rules out overflow; the cross terms are summed with Neumaier's compensated summation in the same pass.
// Each coordinate is resolved relative to the largest one (to 2⁻¹⁰⁷⁴ of it). The input is only read.
// Self-intersection is not detected (it would cost O(n log n)): for a
// self-intersecting polygon the result is the magnitude of the signed sum, not an area.
// Errors: ErrDegenerate when there are fewer than 3 vertices or the computed area is exactly 0 (all
// vertices collinear); ErrNotFinite when a coordinate is NaN or ±Inf; ErrOverflow when the area is beyond
// float64.
//
// ID:      NUM-GEOM-032
// Time:    O(n)   n = len(vertices)
// Space:   O(1) extra, no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5 "Polygons", eq. (4.5.1), p. 318
// Source:  Braden, "The Surveyor's Area Formula", The College Mathematics Journal, 1986, 17(4):326–337, doi:10.2307/2686282
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 (the compensated sum)
// Tags:    geometry, plane, polygon, shoelace, compensated-summation
// Units:   coordinates in one length unit; result in that unit squared
// Example: (0,0), (4,0), (4,3), (2,5), (0,3) → 16 (hand-derived: rectangle 12 + triangle 4)
func PolygonArea(vertices []Point) (float64, error) {
	n := len(vertices)
	if n < 3 {
		return 0, fmt.Errorf("%w: vertices: fewer than 3", ErrDegenerate)
	}
	e, err := verticesScale(vertices)
	if err != nil {
		return 0, err
	}
	x0, y0 := math.Ldexp(vertices[0].X, -e), math.Ldexp(vertices[0].Y, -e)
	px, py := math.Ldexp(vertices[1].X, -e)-x0, math.Ldexp(vertices[1].Y, -e)-y0
	var sum, comp float64
	// With vertex 0 at the origin, the two cross terms that involve it are 0.
	for _, v := range vertices[2:] {
		qx, qy := math.Ldexp(v.X, -e)-x0, math.Ldexp(v.Y, -e)-y0
		sum, comp = reduce.NeumaierAdd(sum, comp, px*qy-qx*py)
		px, py = qx, qy
	}
	twice := sum + comp
	if twice == 0 {
		return 0, fmt.Errorf("%w: vertices: zero area", ErrDegenerate)
	}
	return checkResult(math.Ldexp(math.Abs(twice)/2, 2*e), "area")
}
