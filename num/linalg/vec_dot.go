package linalg

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// twoProduct returns x = fl(a·b) and the exact rounding error y, so a·b = x + y exactly
// (Ogita, Rump, Oishi 2005, Algorithm 3.5, with a fused multiply-add).
func twoProduct(a, b float64) (x, y float64) {
	// The explicit conversion rounds a·b and forbids the compiler from fusing it into the caller's next
	// addition; without it TwoSum(p, a·b) could see the unrounded product and the transformation would no
	// longer be error-free.
	x = float64(a * b)
	y = math.FMA(a, b, -x)
	return x, y
}

// dot2 accumulates a dot product as if in twice the working precision (Ogita, Rump, Oishi 2005,
// Algorithm 5.3, Dot2): p carries the ordinary sum, s the summed rounding errors.
type dot2 struct{ p, s float64 }

func (d *dot2) add(a, b float64) {
	h, r := twoProduct(a, b)
	var q float64
	d.p, q = reduce.TwoSum(d.p, h)
	d.s += q + r
}

func (d *dot2) result() float64 { return d.p + d.s }

// dotInputError reports the first NaN or ±Inf element of a or b, or ErrOverflow when every element is
// finite (the sum overflowed).
func dotInputError(a, b []float64) error {
	for i := range a {
		if !reduce.AllFinite(a[i]) {
			return fmt.Errorf("%w: a[%d]", ErrNotFinite, i)
		}
		if !reduce.AllFinite(b[i]) {
			return fmt.Errorf("%w: b[%d]", ErrNotFinite, i)
		}
	}
	return errOverflow("a·b")
}

// VecDot returns the dot product Σ aᵢbᵢ, accumulated with fused multiply-adds. Its error grows with n
// and with cancellation; VecDotCompensated is as accurate as twice the working precision.
// Errors: ErrDimension when len(a) ≠ len(b); ErrNotFinite for a NaN or ±Inf element; ErrOverflow when
// the sum is outside the float64 range. The dot product of two empty vectors is 0.
//
// ID:      NUM-LINALG-039
// Time:    O(n)   n = len(a)
// Space:   O(1), no allocation
// Source:  Strang, Introduction to Linear Algebra, 5th ed., Wellesley-Cambridge Press, 2016, §1.2 "Lengths and Dot Products"
// Tags:    vector, n-dimensional, dot-product, inner-product, linear-algebra
// Units:   a in unit A, b in unit B; the result in A·B
// Example: (1, 2, 3, 4)·(4, 3, 2, 1) → 20 (hand-checked)
func VecDot(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, lengthError("a", len(a), "b", len(b))
	}
	s := 0.0
	for i := range a {
		s = math.FMA(a[i], b[i], s)
	}
	if !reduce.AllFinite(s) {
		return 0, dotInputError(a, b)
	}
	return s, nil
}

// VecDotCompensated returns the dot product Σ aᵢbᵢ as accurately as if it were computed in twice the
// working precision and then rounded: |res − aᵀb| ≤ u|aᵀb| + γ²ₙ|a|ᵀ|b| (u = 2⁻⁵³, γₙ = nu/(1 − nu)),
// by Ogita, Rump and Oishi's Dot2 with error-free products from fused multiply-adds. Use it for long
// vectors and for dot products with heavy cancellation.
// Errors: as VecDot.
//
// ID:      NUM-LINALG-040
// Time:    O(n)   n = len(a); about 6n floating-point operations
// Space:   O(1), no allocation
// Source:  Ogita, Rump, Oishi, "Accurate Sum and Dot Product", SIAM Journal on Scientific Computing, 2005, 26(6):1955–1988, doi:10.1137/030601818, Algorithm 5.3 (Dot2), Algorithms 3.1 (TwoSum) and 3.5 (TwoProductFMA), Proposition 5.5 (error bound)
// Tags:    vector, n-dimensional, dot-product, compensated, accurate, error-free-transformation, linear-algebra
// Units:   a in unit A, b in unit B; the result in A·B
// Example: (1e16, 1, −1e16)·(1, 1, 1) → 1 (the plain sum gives 0; hand-checked)
func VecDotCompensated(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return 0, lengthError("a", len(a), "b", len(b))
	}
	var d dot2
	for i := range a {
		d.add(a[i], b[i])
	}
	r := d.result()
	if !reduce.AllFinite(r) {
		return 0, dotInputError(a, b)
	}
	return r, nil
}
