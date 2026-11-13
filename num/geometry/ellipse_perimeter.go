package geometry

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ErrNoConvergence: an iteration reached its cap. Never expected for valid input; a library must not loop.
var ErrNoConvergence = errors.New("geometry: iteration did not converge")

const (
	// agmMaxSteps caps the arithmetic-geometric mean. From (1, t) with t ≥ 2⁻³⁰ it meets agmTolerance in
	// at most about 7 steps (the convergence is quadratic once the means are close).
	agmMaxSteps = 64
	// agmTolerance is Carlson's stopping rule |xₙ − yₙ| < 2.7·√r·|xₙ| with r = 2⁻⁵⁶, a relative error
	// bound eight times below double precision's 2⁻⁵³.
	agmTolerance = 2.7 * 0x1p-28
	// flatEllipseRatio: below this b/a, E(k) rounds to exactly 1. By DLMF 19.12.2,
	// E(k) = 1 + ½k′²(ln(4/k′) − ½) + O(k′⁴ ln k′) with k′ = b/a, and at k′ = 2⁻³⁰ the correction is
	// under 10⁻¹⁷, below half an ulp of 1. The AGM's final subtraction loses digits in proportion to
	// ln(1/k′) there, so the limit is both faster and more accurate.
	flatEllipseRatio = 0x1p-30
)

// completeEllipticE returns E(k), Legendre's complete elliptic integral of the second kind, for
// k² = 1 − t², 0 < t ≤ 1, as 2R_G(t², 1, 0) by Carlson's arithmetic-geometric mean algorithm:
// x₀ = 1, y₀ = t, xₘ₊₁ = (xₘ + yₘ)/2, yₘ₊₁ = √xₘ·√yₘ, stopping at n with |xₙ − yₙ| < 2.7√r·xₙ; then
// R_F = π/(xₙ + yₙ) and 2R_G = (((x₀ + y₀)/2)² − Σₘ₌₁ⁿ 2ᵐ⁻²(xₘ − yₘ)²)·R_F (Carlson 1995, eqs. (42)–(46)).
// ok is false if the cap is reached.
func completeEllipticE(t float64) (e float64, ok bool) {
	x, y := 1.0, t
	half := (x + y) / 2
	bracket := half * half
	var sum, comp float64
	weight := 0.5 // 2^(m−2) for m = 1
	for range agmMaxSteps {
		if math.Abs(x-y) < agmTolerance*x {
			return (bracket - (sum + comp)) * (math.Pi / (x + y)), true
		}
		// √x·√y, not √(xy): xy underflows to 0 when t is near the smallest subnormal.
		x, y = (x+y)/2, math.Sqrt(x)*math.Sqrt(y)
		d := x - y
		sum, comp = reduce.NeumaierAdd(sum, comp, weight*d*d)
		weight *= 2
	}
	return 0, false
}

// EllipsePerimeter returns the exact perimeter of an ellipse from its two semi-axes (either order):
// L = 4a·E(k), k² = 1 − b²/a², a ≥ b, with E evaluated by the arithmetic-geometric mean. The semi-axes are
// scaled to (1, b/a) first (L is homogeneous of degree 1), which keeps the iteration in range for every
// ratio. Below b/a = 2⁻³⁰ the result is 4a: there E(k) = 1 to double precision (DLMF 19.12.2), while the
// AGM's last subtraction would lose digits in proportion to ln(a/b). A circle (a == b) gives 2πa.
// Accurate to a few ulps; EllipsePerimeterRamanujan is the closed-form approximation.
// Errors: ErrNotFinite when a semi-axis is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrOverflow
// when the perimeter is beyond float64; ErrNoConvergence if the AGM reaches its cap (not expected).
//
// ID:      NUM-GEOM-046
// Time:    O(1); at most 64 AGM steps, about 7 in practice
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 19.9.9 (L(a,b) = 4aE(k), k² = 1 − b²/a²), https://dlmf.nist.gov/19.9 (read 2026-11-12)
// Source:  Carlson, "Numerical computation of real or complex elliptic integrals", Numerical Algorithms, 1995, 10(1):13–26, doi:10.1007/BF02198293, §2, eqs. (42)–(46) (E(k) = 2R_G(1 − k², 1, 0) by arithmetic and geometric means); preprint arXiv:math/9409227 (read 2026-11-12)
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 19.12.2 (E(k) as k′ → 0, used for b/a < 2⁻³⁰), https://dlmf.nist.gov/19.12 (read 2026-11-12)
// Source:  Zwillinger (ed.), CRC Standard Mathematical Tables and Formulae, 31st ed., Chapman & Hall/CRC, 2003, §4.6.3 "Additional properties of ellipses", eq. (4.6.11), p. 330 (perimeter of C = 4aE(π/2, e))
// Tags:    geometry, plane, ellipse, conic, perimeter, agm, elliptic-integral
// Units:   semi-axes in one length unit; result in that unit
// Example: 2, 1 → 9.688448220547676; 5, 3 → 25.526998863398128 (hand-derived, 60-digit AGM)
func EllipsePerimeter(semiAxisA, semiAxisB float64) (float64, error) {
	if err := checkPositive(semiAxisA, "semiAxisA"); err != nil {
		return 0, err
	}
	if err := checkPositive(semiAxisB, "semiAxisB"); err != nil {
		return 0, err
	}
	a, b := semiAxisA, semiAxisB
	if a < b {
		a, b = b, a
	}
	e := 1.0
	if t := b / a; t >= flatEllipseRatio {
		var ok bool
		if e, ok = completeEllipticE(t); !ok {
			return 0, fmt.Errorf("%w: arithmetic-geometric mean", ErrNoConvergence)
		}
	}
	return checkResult(product(4, a, e, 1), "perimeter")
}
