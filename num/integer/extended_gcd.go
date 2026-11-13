package integer

import (
	"fmt"
	"math"
)

// ExtendedGCD returns g = gcd(a, b) ≥ 0 and Bézout coefficients x, y with a·x + b·y = g. The coefficients
// are the ones the extended Euclidean algorithm produces, so |x| ≤ max(1, |b|/g) and |y| ≤ max(1, |a|/g), and no
// intermediate value overflows int64. Signs follow the inputs: ExtendedGCD(-a, b) gives (g, -x, y).
// ExtendedGCD(0, 0) = (0, 0, 0); ExtendedGCD(a, 0) = (|a|, sign(a), 0).
// Errors: ErrOverflow when a or b is math.MinInt64, whose magnitude 2⁶³ does not fit int64.
//
// ID:      NUM-INTEG-006
// Time:    O(log min(|a|, |b|)) division steps (Lamé's theorem)
// Space:   O(1) (no allocation)
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1998, §4.5.2, Algorithm X (extended Euclid's algorithm)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §31.2 (EXTENDED-EUCLID)
// Tags:    gcd, extended-gcd, bezout, extended-euclid, modular-inverse, number-theory, no-alloc
func ExtendedGCD(a, b int64) (g, x, y int64, err error) {
	if a == math.MinInt64 || b == math.MinInt64 {
		return 0, 0, 0, fmt.Errorf("%w: a = %d, b = %d, |MinInt64| does not fit int64", ErrOverflow, a, b)
	}
	ua, ub := a, b
	if ua < 0 {
		ua = -ua
	}
	if ub < 0 {
		ub = -ub
	}
	if ua == 0 && ub == 0 {
		return 0, 0, 0, nil
	}
	// Algorithm X: (u1, u2, u3) ← (1, 0, |a|), (v1, v2, v3) ← (0, 1, |b|); invariant u1·|a| + u2·|b| = u3.
	u1, u2, u3 := int64(1), int64(0), ua
	v1, v2, v3 := int64(0), int64(1), ub
	for v3 != 0 {
		q := u3 / v3
		u1, v1 = v1, u1-q*v1
		u2, v2 = v2, u2-q*v2
		u3, v3 = v3, u3-q*v3
	}
	if a < 0 {
		u1 = -u1
	}
	if b < 0 {
		u2 = -u2
	}
	return u3, u1, u2, nil
}
