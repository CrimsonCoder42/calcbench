package geometry

import "math"

// segmentSeriesBelowRad is where SegmentArea switches from the sine series to θ − sin θ. At θ = 0.5 the
// direct difference has lost about 5 bits to cancellation (θ/(θ − sin θ) ≈ 24), and the series needs about
// 9 terms (each term is at most θ²/20 = 1/80 of the one before).
const segmentSeriesBelowRad = 0.5

// thetaMinusSin returns θ − sin θ for 0 < θ ≤ 2π without catastrophic cancellation: below
// segmentSeriesBelowRad it sums the Maclaurin series θ³/3! − θ⁵/5! + θ⁷/7! − … (the sine series with its
// first term removed) until a term is below 2⁻⁵³ of the sum; the loop is capped at 30 terms.
func thetaMinusSin(theta float64) float64 {
	if theta >= segmentSeriesBelowRad {
		return theta - math.Sin(theta)
	}
	t2 := theta * theta
	term := theta * t2 / 6
	sum := term
	for k := 4; k < 64; k += 2 {
		term *= -t2 / float64(k*(k+1))
		sum += term
		if math.Abs(term) <= 0x1p-53*math.Abs(sum) {
			break
		}
	}
	return sum
}

// SegmentArea returns the area of a circular segment (the region between a chord and its arc) of the given
// radius and central angle: ½·r²·(θ − sin θ). For θ < 0.5 the difference θ − sin θ is summed from the sine
// series instead (at θ = 10⁻⁴ the direct form keeps only about 8 digits).
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrNonPositive when radius ≤ 0; ErrOutOfRange when
// centralAngleRad is outside (0, 2π]; ErrOverflow when the area is beyond float64.
//
// ID:      NUM-GEOM-040
// Time:    O(1); at most 30 series terms
// Space:   O(1), no allocation
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.6.6 "Circles", p. 334 (area of segment = ½R²(θ − sin θ))
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.19.1 (sin z = z − z³/3! + z⁵/5! − …), https://dlmf.nist.gov/4.19 (read 2026-11-12)
// Tags:    geometry, plane, circle, segment
// Units:   radius in a length unit; centralAngleRad in radians; result in that unit squared
// Example: 2, π/3 → 2(π/3 − √3/2) ≈ 0.362344294824; 1, π → π/2 (hand-derived)
func SegmentArea(radius, centralAngleRad float64) (float64, error) {
	if err := checkPositive(radius, "radius"); err != nil {
		return 0, err
	}
	if err := checkCentralAngle(centralAngleRad, "centralAngleRad"); err != nil {
		return 0, err
	}
	return checkResult(product(0.5, radius, radius, thetaMinusSin(centralAngleRad)), "area")
}
