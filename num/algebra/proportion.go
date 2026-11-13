package algebra

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// ProportionFourthTerm returns d such that a : b = c : d, by cross products (a·d = b·c, so d = b·c/a).
// To solve for another position, reorder the proportion: the missing term of x : b = c : d is d·b/c
// with this function as ProportionFourthTerm(c, d, b).
// Errors: ErrNotFinite for NaN or ±Inf input; ErrDegenerate when a = 0; ErrOverflow when d is outside
// the float64 range.
//
// ID:      NUM-ALGEB-036
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Marecek, Anthony-Smith and Honeycutt Mathis, Prealgebra 2e, OpenStax, 2020, §6.5 "Solve Proportions and their Applications" (cross products: a/b = c/d means a·d = b·c)
// Tags:    ratio, proportion, cross-multiplication, rule-of-three, scalar, algebra
// Units:   a and c in one unit, b and d in one unit
// Example: 3 : 4 = 600 : d → d = 800 (hand-checked)
func ProportionFourthTerm(a, b, c float64) (float64, error) {
	if !reduce.AllFinite(a, b, c) {
		return 0, fmt.Errorf("%w: a = %v, b = %v, c = %v", ErrNotFinite, a, b, c)
	}
	if a == 0 {
		return 0, fmt.Errorf("%w: a = 0", ErrDegenerate)
	}
	// Divide first when b·c alone would overflow but the quotient fits.
	d := b * c / a
	if !reduce.AllFinite(d) {
		d = b * (c / a)
	}
	return finiteResult(d, "fourth term")
}
