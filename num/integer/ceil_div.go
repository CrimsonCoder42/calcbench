package integer

import (
	"fmt"
	"math"
)

// CeilDiv returns ⌈a/b⌉, the least integer not less than the exact quotient a/b, for every int64 pair whose
// result fits int64. It takes Go's truncated quotient a/b and adds 1 when the remainder is non-zero and has
// the same sign as b (the exact quotient is then positive and was rounded down).
// Errors: ErrZeroDenominator when b = 0; ErrOverflow when a = math.MinInt64 and b = −1 (the result 2⁶³ does
// not fit int64).
//
// ID:      NUM-INTEG-013
// Time:    O(1)
// Space:   O(1) (no allocation)
// Source:  Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §3.1 eq. (3.4) (⌊−x⌋ = −⌈x⌉); Exercise 3.12 (⌈n/m⌉ = ⌊(n + m − 1)/m⌋ for m > 0)
// Source:  The Go Programming Language Specification, "Integer operators", https://go.dev/ref/spec#Integer_operators (read 2026-11-13): "x = q*y + r  and  |r| < |y|" "with x / y truncated towards zero"; "if the dividend x is the most negative value for the int type of x, the quotient q = x / -1 is equal to x (and r = 0) due to two's-complement integer overflow" (so r is 0 or has the sign of a, derived)
// Tags:    ceil-div, ceiling-division, integer-division, rounding, overflow-checked, no-alloc
// Units:   a, b and the result in one unit or dimensionless; the result is a count of b in a, rounded up
// Example: CeilDiv(10, 3) = 4; CeilDiv(9, 3) = 3; CeilDiv(−7, 2) = −3 (⌈−3.5⌉)
func CeilDiv(a, b int64) (int64, error) {
	if b == 0 {
		return 0, fmt.Errorf("%w: ceil(%d / 0)", ErrZeroDenominator, a)
	}
	if a == math.MinInt64 && b == -1 {
		return 0, fmt.Errorf("%w: ceil(%d / -1) = 2^63 > MaxInt64", ErrOverflow, a)
	}
	q, r := a/b, a%b
	if r != 0 && (r > 0) == (b > 0) {
		q++
	}
	return q, nil
}
