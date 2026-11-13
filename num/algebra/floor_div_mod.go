package algebra

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// FloorDivMod returns q = ⌊x/y⌋ and r = x − y·q, the floored quotient and the remainder with the sign of y
// (0 ≤ r < y for y > 0, y < r ≤ 0 for y < 0), as Knuth defines x mod y. r is computed exactly with
// math.Mod and adjusted by y; q is (x − r)/y rounded to the nearest integer, so x ≈ y·q + r even when x/y
// rounds across an integer. A zero r takes the sign of y.
// Errors: ErrNotFinite for NaN or ±Inf x or y; ErrDegenerate when y = 0.
//
// ID:      NUM-ALGEB-023
// Time:    O(1)
// Space:   O(1)
// Source:  Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §1.2.4 "Integer Functions and Elementary Number Theory", eq. (1) (x mod y = x − y⌊x/y⌋)
// Tags:    floor-division, modulo, divmod, remainder, scalar, algebra
// Units:   x and y in one unit; q dimensionless; r in that unit
// Example: (7, 3) → (2, 1); (−7, 3) → (−3, 2); (7, −3) → (−3, −2) (Knuth §1.2.4 sign convention; hand-checked)
func FloorDivMod(x, y float64) (q, r float64, err error) {
	if !reduce.AllFinite(x) || !reduce.AllFinite(y) {
		return 0, 0, fmt.Errorf("%w: x = %v, y = %v", ErrNotFinite, x, y)
	}
	if y == 0 {
		return 0, 0, fmt.Errorf("%w: y = 0", ErrDegenerate)
	}
	r = math.Mod(x, y)
	d := (x - r) / y
	if r != 0 {
		if (y < 0) != (r < 0) {
			r += y
			d--
		}
	} else {
		r = math.Copysign(0, y)
	}
	if d != 0 {
		q = math.Floor(d)
		if d-q > 0.5 {
			q++
		}
	} else {
		q = math.Copysign(0, x/y)
	}
	return q, r, nil
}

// FloorMod returns x mod y = x − y⌊x/y⌋, the remainder with the sign of y (math.Mod truncates and takes the
// sign of x instead), as FloorDivMod.
// Errors: as FloorDivMod.
//
// ID:      NUM-ALGEB-024
// Time:    O(1)
// Space:   O(1)
// Source:  Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §1.2.4, eq. (1)
// Tags:    modulo, floor-mod, remainder, scalar, algebra
// Units:   x, y and the result in one unit
// Example: (−7, 3) → 2; (7, −3) → −2; (5.5, 2) → 1.5 (hand-checked)
func FloorMod(x, y float64) (float64, error) {
	_, r, err := FloorDivMod(x, y)
	return r, err
}
