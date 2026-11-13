package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// QuadraticRoots returns both roots of a·x² + b·x + c = 0 as complex numbers, r1 first: by smaller real
// part, then by smaller imaginary part. When D = b² − 4ac ≥ 0 (computed as NUM-ALGEB-006) the roots are
// those of QuadraticRealRoots, with imaginary part exactly 0 (a double root is returned twice). When
// D < 0 they are re ∓ i·im with re = −b/(2a) and im = √(−D)/(2|a|), each computed separately, so there
// is no cancellation; r1 has the negative imaginary part. All work is done on power-of-two-scaled values.
// Errors: ErrDegenerate when a = 0; ErrNotFinite when a, b or c is NaN or ±Inf; ErrOverflow when a
// root's real or imaginary part is outside the float64 range. Both roots are 0 with any error.
//
// ID:      NUM-ALGEB-008
// Time:    O(1)
// Space:   O(1)
// Source:  Press, Teukolsky, Vetterling, Flannery, Numerical Recipes in C: The Art of Scientific Computing, 2nd ed., Cambridge University Press, 1992, §5.6, eqs. (5.6.4)–(5.6.5)
// Source:  NIST Digital Library of Mathematical Functions, §1.11(iii) (quadratic equations), https://dlmf.nist.gov/1.11 (read 2026-11-12)
// Tags:    quadratic, roots, complex, quadratic-formula, algebra
// Units:   a, b, c and the roots are dimensionless
// Example: x² + 2x + 5 → (−1 − 2i, −1 + 2i); x² − 3x + 2 → (1 + 0i, 2 + 0i) (hand-checked)
func QuadraticRoots(a, b, c float64) (r1, r2 complex128, err error) {
	if err := checkQuadratic(a, b, c); err != nil {
		return 0, 0, err
	}
	if a == 0 {
		return 0, 0, fmt.Errorf("%w: a = 0", ErrDegenerate)
	}
	d, exp, bScaled := scaledDiscriminant(a, b, c)
	if d >= 0 {
		x1, x2, err := quadraticRootsFromScaled(a, c, d, exp, bScaled)
		if err != nil {
			return 0, 0, err
		}
		if d == 0 {
			x2 = x1
		}
		if x2 < x1 {
			x1, x2 = x2, x1
		}
		return complex(x1, 0), complex(x2, 0), nil
	}
	// a = fa·2^ka and b = bScaled·2^(exp/2), so b/(2a) = (bScaled/fa)·2^(exp/2−ka−1) and
	// √(−D)/(2|a|) = (√(−d)/|fa|)·2^(exp/2−ka−1).
	fa, ka := math.Frexp(a)
	shift := exp/2 - ka - 1
	re := -math.Ldexp(bScaled/fa, shift)
	im := math.Ldexp(math.Sqrt(-d)/math.Abs(fa), shift)
	if !reduce.AllFinite(re) || !reduce.AllFinite(im) {
		return 0, 0, fmt.Errorf("%w: a root of a·x² + b·x + c", ErrOverflow)
	}
	return complex(re, -im), complex(re, im), nil
}
