package geometry

import (
	"fmt"
	"math"
)

const (
	// carlsonMaxSteps caps the duplication loop of carlsonRD. With arguments in [0, 1] it stops in under
	// 20 steps, even when one argument is near the smallest subnormal (the ratio of the extreme arguments
	// goes to its square root each step until they are comparable).
	carlsonMaxSteps = 100
	// carlsonRDQ is (r/4)^(−1/6) for the relative error bound r = 2⁻⁵⁶ in Carlson's stopping rule (35):
	// 2^(58/6). The duplication stops when 4⁻ⁿ·Q < Aₙ.
	carlsonRDQ = 812.7493386077178
)

// carlsonRD returns Carlson's symmetric elliptic integral of the second kind R_D(x, y, z) for finite
// x, y ≥ 0 (at most one 0) and z > 0, by Carlson's duplication algorithm (Carlson 1995, §2, eqs.
// (35)–(41)), with relative error below 2⁻⁵⁶ before roundoff. ok is false if the cap is reached.
func carlsonRD(x, y, z float64) (rd float64, ok bool) {
	a0 := (x + y + 3*z) / 5
	q := carlsonRDQ * math.Max(math.Abs(a0-x), math.Max(math.Abs(a0-y), math.Abs(a0-z)))
	x0, y0 := x, y
	a, fac, sum := a0, 1.0, 0.0
	for range carlsonMaxSteps {
		if fac*q < a {
			bx := fac * (a0 - x0) / a
			by := fac * (a0 - y0) / a
			bz := -(bx + by) / 3
			xy, z2 := bx*by, bz*bz
			e2 := xy - 6*z2
			e3 := (3*xy - 8*z2) * bz
			e4 := 3 * (xy - z2) * z2
			e5 := xy * z2 * bz
			series := 1 - 3.0/14*e2 + e3/6 + 9.0/88*e2*e2 - 3.0/22*e4 - 9.0/52*e2*e3 + 3.0/26*e5
			return fac/(a*math.Sqrt(a))*series + 3*sum, true
		}
		sx, sy, sz := math.Sqrt(x), math.Sqrt(y), math.Sqrt(z)
		lambda := sx*sy + sx*sz + sy*sz
		sum += fac / (sz * (z + lambda))
		fac /= 4
		x, y, z = (x+lambda)/4, (y+lambda)/4, (z+lambda)/4
		a = (a + lambda) / 4
	}
	return 0, false
}

// carlsonRG returns Carlson's completely symmetric elliptic integral R_G(x, y, z) for finite x, y, z ≥ 0,
// from 6R_G = x(y+z)R_D(y,z,x) + y(z+x)R_D(z,x,y) + z(x+y)R_D(x,y,z) (DLMF 19.21.11). Every term is
// non-negative, so nothing cancels, unlike 19.21.10. A term whose leading factor is 0 is 0 and is skipped
// (its R_D would be undefined); with two arguments 0, R_G(0, 0, z) = √z/2 (Carlson 1995, §2). ok is false
// if an R_D reaches its cap.
func carlsonRG(x, y, z float64) (rg float64, ok bool) {
	zeros := 0
	for _, v := range [...]float64{x, y, z} {
		if v == 0 {
			zeros++
		}
	}
	if zeros >= 2 {
		return math.Sqrt(x+y+z) / 2, true
	}
	sum := 0.0
	for _, t := range [...][3]float64{{x, y, z}, {y, z, x}, {z, x, y}} {
		lead, p, r := t[0], t[1], t[2]
		if lead == 0 {
			continue
		}
		rd, ok := carlsonRD(p, r, lead)
		if !ok {
			return 0, false
		}
		sum += lead * (p + r) * rd
	}
	return sum / 6, true
}

// EllipsoidSurfaceArea returns the exact surface area of an ellipsoid from its three semi-axes (any
// order): S = 3V·R_G(a⁻², b⁻², c⁻²), V = 4πabc/3 (DLMF 19.33.1). With the axes sorted a ≥ b ≥ c and the
// homogeneity of R_G this is S = 4πab·R_G(c²/a², c²/b², 1), whose arguments lie in (0, 1]: nothing
// overflows before the area does. R_G is summed from three non-negative R_D terms (DLMF 19.21.11), each by
// Carlson's duplication algorithm. A sphere gives 4πa²; when c is so small against a and b that both ratios
// underflow, R_G(0, 0, 1) = ½ gives the flat limit 2πab.
// Errors: ErrNotFinite when a semi-axis is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrOverflow
// when the area is beyond float64; ErrNoConvergence if a duplication loop reaches its cap (not expected).
//
// ID:      NUM-GEOM-070
// Time:    O(1); three R_D evaluations of at most 100 duplication steps, under 20 in practice
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 19.33.1 (S = 3V R_G(a⁻², b⁻², c⁻²)), https://dlmf.nist.gov/19.33, and eq. 19.21.11 (R_G from R_D), https://dlmf.nist.gov/19.21 (read 2026-11-12)
// Source:  Carlson, "Numerical computation of real or complex elliptic integrals", Numerical Algorithms, 1995, 10(1):13–26, doi:10.1007/BF02198293, §2 "Algorithm for R_D", eqs. (35)–(41), and §3 check values; preprint arXiv:math/9409227 (read 2026-11-12)
// Tags:    geometry, solid, quadric, ellipsoid, elliptic-integral, carlson
// Units:   semi-axes in one length unit; result in that unit squared
// Example: 3, 2, 1 → 48.882146302582; 2, 1, 1 → 21.478435327884 (hand-derived, 60 digits; the second equals the prolate-spheroid closed form, CRC 31st ed. eq. (4.18.1))
func EllipsoidSurfaceArea(semiAxisA, semiAxisB, semiAxisC float64) (float64, error) {
	if err := checkSemiAxes3(semiAxisA, semiAxisB, semiAxisC); err != nil {
		return 0, err
	}
	a, b, c := semiAxisA, semiAxisB, semiAxisC
	if a < b {
		a, b = b, a
	}
	if b < c {
		b, c = c, b
	}
	if a < b {
		a, b = b, a
	}
	ra, rb := c/a, c/b
	rg, ok := carlsonRG(ra*ra, rb*rb, 1)
	if !ok {
		return 0, fmt.Errorf("%w: Carlson duplication", ErrNoConvergence)
	}
	return checkResult(product(4*math.Pi, a, b, rg), "area")
}
