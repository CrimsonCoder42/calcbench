package reduce

import (
	"fmt"
	"math"
)

// Prod returns the product of values. The running product is held as a mantissa in [0.5, 1) and a
// separate integer exponent (math.Frexp), so no intermediate product overflows or underflows: the result
// is ±Inf-free whenever the true product is representable, whatever the order of the factors. The result
// is rounded once per factor, as a plain loop is. An empty slice has product 1. A result below the
// smallest subnormal is 0 (gradual underflow, IEEE 754).
// Errors: ErrNotFinite (with the index) when a value is NaN or ±Inf; ErrOverflow when the product exceeds
// ±MaxFloat64.
//
// ID:      ALG-RED-004
// Time:    O(n)   n = len(values)
// Space:   O(1) (no allocation)
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §3.1 (error of a product, γₙ bound) and §27.8 (avoiding overflow by scaling)
// Tags:    product, prod, reduction, scaling, overflow-safe, no-alloc
func Prod(values []float64) (float64, error) {
	mant, exp := 1.0, 0
	for i, x := range values {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return 0, fmt.Errorf("%w: values[%d] = %v", ErrNotFinite, i, x)
		}
		mant, exp = scaledMul(mant, exp, x)
	}
	p := math.Ldexp(mant, exp)
	if math.IsInf(p, 0) {
		return 0, fmt.Errorf("%w: product", ErrOverflow)
	}
	return p, nil
}

// scaledMul multiplies the scaled value mant·2^exp by a finite x and renormalises, returning the new
// mantissa (|mant| in [0.5, 1), or 0) and exponent. The exponent cannot overflow an int: it moves by at
// most 1075 per factor and a slice of 2^53 factors is far beyond memory.
func scaledMul(mant float64, exp int, x float64) (float64, int) {
	fx, ex := math.Frexp(x)
	m, em := math.Frexp(mant * fx) // |mant·fx| is in [0.25, 1): no overflow, no underflow.
	if m == 0 {
		return m, 0
	}
	return m, exp + ex + em
}
