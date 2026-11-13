package geometry

import "math"

// EllipsePerimeterRamanujan returns Ramanujan's second approximation to the perimeter of an ellipse from
// its two semi-axes (either order): π(a + b)·(1 + 3λ²/(10 + √(4 − 3λ²))), λ = (a − b)/(a + b). It
// underestimates; the relative error is about 3λ¹⁰/2¹⁷ for small λ (4.6·10⁻¹⁰ at a = 2b) and at most
// 4.1·10⁻⁴, reached as b/a → 0. The sum a + b is formed as 2(a/2 + b/2) so it cannot overflow.
// EllipsePerimeter is exact and is the default; this closed form is for callers that want it.
// Errors: ErrNotFinite when a semi-axis is NaN or ±Inf; ErrNonPositive when a semi-axis ≤ 0; ErrOverflow
// when the perimeter is beyond float64.
//
// ID:      NUM-GEOM-047
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ramanujan, "Modular equations and approximations to π", Quarterly Journal of Mathematics, 1914, 45:350–372, §15 (p = π{(a+b) + 3(a−b)²/[10(a+b) + √(a²+14ab+b²)]}, "obtained empirically", error about 3a·e²⁰/2³⁶), https://ramanujan.sirinudi.org/Volumes/published/ram06.html (read 2026-11-12)
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 19.9.10 (the same approximation in λ, with its leading error term 3λ¹⁰/2¹⁷), https://dlmf.nist.gov/19.9 (read 2026-11-12)
// Tags:    geometry, plane, ellipse, conic, perimeter, approximation, ramanujan
// Units:   semi-axes in one length unit; result in that unit
// Example: 2, 1 → 9.688448216130 (exact 9.688448220548; hand-derived, 60 digits)
func EllipsePerimeterRamanujan(semiAxisA, semiAxisB float64) (float64, error) {
	if err := checkPositive(semiAxisA, "semiAxisA"); err != nil {
		return 0, err
	}
	if err := checkPositive(semiAxisB, "semiAxisB"); err != nil {
		return 0, err
	}
	half := semiAxisA/2 + semiAxisB/2
	lambda := (semiAxisA/2 - semiAxisB/2) / half
	l2 := 3 * lambda * lambda
	factor := 1 + l2/(10+math.Sqrt(4-l2))
	return checkResult(product(2*math.Pi, half, factor, 1), "perimeter")
}
