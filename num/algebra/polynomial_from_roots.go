package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// PolynomialFromRoots returns the real monic polynomial ∏(x − rᵢ) (ascending coefficients, len(roots)+1
// of them) whose roots are the given complex numbers. Each non-real root must have its conjugate in the
// list (matched within tol·max(1, |r|); tol < 0 selects 2⁻⁴⁰); a conjugate pair r, r̄ is multiplied in as the
// real quadratic x² − 2Re(r)·x + |r|², so the coefficients stay exactly real. An empty list gives [1].
// Errors: ErrNotFinite for a NaN or ±Inf part; ErrDegenerate when a non-real root has no conjugate;
// ErrOverflow when a coefficient leaves the float64 range.
//
// ID:      NUM-ALGEB-017
// Time:    O(n²)   n = len(roots) (one linear or quadratic factor at a time)
// Space:   O(n) for the result and a used-flag per root
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6 "Polynomial Arithmetic" (product of linear factors); real quadratic factor of a conjugate pair
// Tags:    polynomial, roots, from-roots, monic, conjugate-pairs, algebra
// Units:   roots and coefficients are dimensionless
// Example: roots 1, 2 → [2, −3, 1] (x² − 3x + 2); roots ±i → [1, 0, 1] (hand-checked)
func PolynomialFromRoots(roots []complex128, tol float64) ([]float64, error) {
	if tol < 0 {
		tol = 0x1p-40
	}
	for i, r := range roots {
		if !reduce.AllFinite(real(r)) || !reduce.AllFinite(imag(r)) {
			return nil, fmt.Errorf("%w: roots[%d] = %v", ErrNotFinite, i, r)
		}
	}
	n := len(roots)
	c := make([]float64, n+1)
	c[0] = 1
	deg := 0
	used := make([]bool, n)
	for i, r := range roots {
		if used[i] {
			continue
		}
		used[i] = true
		if imag(r) == 0 {
			// multiply by (x − r)
			for k := deg + 1; k >= 1; k-- {
				c[k] = c[k-1] - real(r)*c[k]
			}
			c[0] = -real(r) * c[0]
			deg++
			continue
		}
		mag := max(1, math.Hypot(real(r), imag(r)))
		pair := -1
		for j := i + 1; j < n; j++ {
			s := roots[j]
			if !used[j] && math.Abs(real(s)-real(r)) <= tol*mag && math.Abs(imag(s)+imag(r)) <= tol*mag {
				pair = j
				break
			}
		}
		if pair < 0 {
			return nil, fmt.Errorf("%w: roots[%d] = %v has no conjugate", ErrDegenerate, i, r)
		}
		used[pair] = true
		p := -2 * real(r)
		q := real(r)*real(r) + imag(r)*imag(r)
		// multiply by (x² + p·x + q)
		for k := deg + 2; k >= 0; k-- {
			v := 0.0
			if k-2 >= 0 {
				v += c[k-2]
			}
			if k-1 >= 0 && k-1 <= deg {
				v += p * c[k-1]
			}
			if k <= deg {
				v += q * c[k]
			}
			c[k] = v
		}
		deg += 2
	}
	for i, v := range c {
		if !reduce.AllFinite(v) {
			return nil, fmt.Errorf("%w: coefficient %d", ErrOverflow, i)
		}
	}
	return c, nil
}

// PolynomialValueFromRoots returns ∏(x − rᵢ), the monic polynomial with real roots rᵢ evaluated at x,
// without forming its coefficients (more accurate near the roots than expanding then evaluating). An empty
// list gives 1.
// Errors: ErrNotFinite for a NaN or ±Inf root or x; ErrOverflow when the product leaves the float64 range.
//
// ID:      NUM-ALGEB-018
// Time:    O(n)   n = len(roots)
// Space:   O(1)
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6 (a polynomial as the product of its linear factors; definitional)
// Tags:    polynomial, roots, evaluate, product, algebra
// Units:   roots, x and the result are dimensionless
// Example: roots 1, 2 at x = 4 → (4 − 1)(4 − 2) = 6 (hand-checked)
func PolynomialValueFromRoots(roots []float64, x float64) (float64, error) {
	if !reduce.AllFinite(x) {
		return 0, fmt.Errorf("%w: x = %v", ErrNotFinite, x)
	}
	v := 1.0
	for i, r := range roots {
		if !reduce.AllFinite(r) {
			return 0, fmt.Errorf("%w: roots[%d] = %v", ErrNotFinite, i, r)
		}
		v *= x - r
	}
	if !reduce.AllFinite(v) {
		return 0, fmt.Errorf("%w: product", ErrOverflow)
	}
	return v, nil
}
