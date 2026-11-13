package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// QuadraticRealRoots returns the real roots of a·x² + b·x + c = 0 as lo ≤ hi, and n, the number of
// distinct real roots: 2 when D = b² − 4ac > 0, 1 when D = 0 (a double root, lo = hi), 0 when D < 0
// (no real roots: lo = hi = 0 and err = nil; that is a result, not an error). D is decided exactly as
// NUM-ALGEB-006 computes it, with no tolerance. The roots use q = −½(b + sgn(b)·√D), sgn(0) = +1, and
// x₁ = q/a, x₂ = c/q, so no two nearly equal numbers are ever subtracted; when q = 0 (b = c = 0) both
// roots are 0. All work is done on power-of-two-scaled values, so b² and 4ac never overflow.
// Errors: ErrDegenerate when a = 0 (a linear equation: use LinearSolve); ErrNotFinite when a, b or c is
// NaN or ±Inf; ErrOverflow when a root is outside the float64 range. lo, hi and n are 0 with any error.
//
// ID:      NUM-ALGEB-007
// Time:    O(1)
// Space:   O(1)
// Source:  Press, Teukolsky, Vetterling, Flannery, Numerical Recipes in C: The Art of Scientific Computing, 2nd ed., Cambridge University Press, 1992, §5.6, eqs. (5.6.4)–(5.6.5)
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §1.8 "Solving a Quadratic Equation"
// Tags:    quadratic, roots, quadratic-formula, algebra
// Units:   a, b, c and the roots are dimensionless
// Example: x² − 3x + 2 → (1, 2, 2); x² − 2x + 1 → (1, 1, 1); x² − 10⁸x + 1 → (1e-8, 1e8, 2) to full precision
func QuadraticRealRoots(a, b, c float64) (lo, hi float64, n int, err error) {
	if err := checkQuadratic(a, b, c); err != nil {
		return 0, 0, 0, err
	}
	if a == 0 {
		return 0, 0, 0, fmt.Errorf("%w: a = 0", ErrDegenerate)
	}
	d, exp, bScaled := scaledDiscriminant(a, b, c)
	if d < 0 {
		return 0, 0, 0, nil
	}
	x1, x2, err := quadraticRootsFromScaled(a, c, d, exp, bScaled)
	if err != nil {
		return 0, 0, 0, err
	}
	if d == 0 {
		return x1, x1, 1, nil
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	return x1, x2, 2, nil
}

// quadraticRootsFromScaled returns x₁ = q/a and x₂ = c/q for d ≥ 0, where b² − 4ac = d·2^exp and
// b = bScaled·2^(exp/2) (from scaledDiscriminant). With q = q′·2^(exp/2), q′ = −½(bScaled + sgn·√d):
// x₁ = (q′/fa)·2^(exp/2−ka) and x₂ = (fc/q′)·2^(kc−exp/2), with a = fa·2^ka, c = fc·2^kc. bScaled and √d
// have the same sign in the sum, so q′ has no cancellation; q′ = 0 only when b = 0 and c = 0.
func quadraticRootsFromScaled(a, c, d float64, exp int, bScaled float64) (x1, x2 float64, err error) {
	fa, ka := math.Frexp(a)
	fc, kc := math.Frexp(c)
	sq := math.Sqrt(d)
	if bScaled < 0 {
		sq = -sq
	}
	q := -0.5 * (bScaled + sq)
	if q == 0 {
		return 0, 0, nil
	}
	x1 = math.Ldexp(q/fa, exp/2-ka)
	x2 = math.Ldexp(fc/q, kc-exp/2)
	if !reduce.AllFinite(x1) || !reduce.AllFinite(x2) {
		return 0, 0, fmt.Errorf("%w: a root of a·x² + b·x + c", ErrOverflow)
	}
	return x1, x2, nil
}
