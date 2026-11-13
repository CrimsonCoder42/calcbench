package algebra

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkCoeffs reports ErrEmpty for an empty slice and ErrNotFinite for a NaN or ±Inf coefficient.
func checkCoeffs(name string, c []float64) error {
	if len(c) == 0 {
		return fmt.Errorf("%w: %s", ErrEmpty, name)
	}
	for i, v := range c {
		if !reduce.AllFinite(v) {
			return fmt.Errorf("%w: %s[%d] = %v", ErrNotFinite, name, i, v)
		}
	}
	return nil
}

// addSub returns a + sign·b, coefficient by coefficient, with len max(len(a), len(b)).
func addSub(a, b []float64, sign float64) ([]float64, error) {
	if err := checkCoeffs("a", a); err != nil {
		return nil, err
	}
	if err := checkCoeffs("b", b); err != nil {
		return nil, err
	}
	c := make([]float64, max(len(a), len(b)))
	copy(c, a)
	for i, v := range b {
		c[i] += sign * v
		if !reduce.AllFinite(c[i]) {
			return nil, fmt.Errorf("%w: coefficient %d", ErrOverflow, i)
		}
	}
	return c, nil
}

// PolynomialAdd returns the coefficients of p + q (ascending order; a[i] multiplies xⁱ). The result has
// max(len(a), len(b)) coefficients; trailing zeros are kept (PolynomialTrim removes them).
// Errors: ErrEmpty when a or b is empty; ErrNotFinite for a NaN or ±Inf coefficient; ErrOverflow when a
// sum is outside the float64 range.
//
// ID:      NUM-ALGEB-011
// Time:    O(n + m)   n = len(a), m = len(b)
// Space:   O(max(n, m)) for the result (one allocation), O(1) extra
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6 "Polynomial Arithmetic" (addition coefficient by coefficient)
// Tags:    polynomial, add, coefficients, algebra
// Units:   coefficients are dimensionless
// Example: (1 + 2x) + (3 + 4x + 5x²) → [4, 6, 5] (hand-checked)
func PolynomialAdd(a, b []float64) ([]float64, error) { return addSub(a, b, 1) }

// PolynomialSubtract returns the coefficients of p − q (ascending order), as PolynomialAdd.
// Errors: as PolynomialAdd.
//
// ID:      NUM-ALGEB-012
// Time:    O(n + m)   n = len(a), m = len(b)
// Space:   O(max(n, m)) for the result (one allocation), O(1) extra
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.6 "Polynomial Arithmetic"
// Tags:    polynomial, subtract, difference, coefficients, algebra
// Units:   coefficients are dimensionless
// Example: (1 + 2x) − (3 + 4x + 5x²) → [−2, −2, −5] (hand-checked)
func PolynomialSubtract(a, b []float64) ([]float64, error) { return addSub(a, b, -1) }
