package geometry

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolygonPerimeter returns the perimeter of a closed polygon from its vertices in order (the first vertex
// is not repeated; the closing side back to it is included): Σ |Vᵢ₊₁ − Vᵢ|, indices mod n. Each side is a
// math.Hypot of coordinates scaled by a power of two (exact, so nothing overflows), and the sides are
// summed with Neumaier's compensated summation. With two vertices the result is twice their distance
// (there and back). The input is only read.
// Errors: ErrDegenerate when there are fewer than 2 vertices or all vertices coincide; ErrNotFinite when a
// coordinate is NaN or ±Inf; ErrOverflow when the perimeter is beyond float64.
//
// ID:      NUM-GEOM-033
// Time:    O(n)   n = len(vertices)
// Space:   O(1) extra, no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.5 "Polygons", p. 318 (the sides A_iA_{i+1}, the last vertex joined to the first), with the side lengths of §4.4.2, eq. (4.4.5), p. 316
// Source:  Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 (the compensated sum)
// Tags:    geometry, plane, polygon, shoelace, compensated-summation
// Units:   coordinates in one length unit; result in that unit
// Example: (0,0), (4,0), (4,3), (2,5), (0,3) → 10 + 4√2 ≈ 15.656854249492 (hand-derived: 4 + 3 + 2√2 + 2√2 + 3)
func PolygonPerimeter(vertices []Point) (float64, error) {
	n := len(vertices)
	if n < 2 {
		return 0, fmt.Errorf("%w: vertices: fewer than 2", ErrDegenerate)
	}
	e, err := verticesScale(vertices)
	if err != nil {
		return 0, err
	}
	last := vertices[n-1]
	px, py := math.Ldexp(last.X, -e), math.Ldexp(last.Y, -e)
	var sum, comp float64
	for _, v := range vertices {
		qx, qy := math.Ldexp(v.X, -e), math.Ldexp(v.Y, -e)
		sum, comp = reduce.NeumaierAdd(sum, comp, math.Hypot(qx-px, qy-py))
		px, py = qx, qy
	}
	total := sum + comp
	if total == 0 {
		return 0, fmt.Errorf("%w: vertices: all coincide", ErrDegenerate)
	}
	return checkResult(math.Ldexp(total, e), "perimeter")
}
