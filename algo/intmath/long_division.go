package intmath

import (
	"fmt"
	"math/bits"
)

// signedInteger is the set of Go's built-in signed integer types and the types defined on them.
type signedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// integer is the set of Go's built-in integer types and the types defined on them.
type integer interface {
	signedInteger | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// magnitude returns |x| as a uint64; |minimum of the type| is representable (2^(w−1)).
func magnitude[T signedInteger](x T) uint64 {
	if x < 0 {
		return -uint64(int64(x))
	}
	return uint64(x)
}

// LongDivision returns the quotient and remainder of dividend / divisor computed by binary long division
// (shift and subtract), without the division or remainder operators. The semantics are Go's for every
// signed type: the quotient is truncated toward zero and the remainder has the sign of the dividend, so
// dividend = quotient·divisor + remainder and |remainder| < |divisor|.
// Errors: ErrDivideByZero when divisor = 0; ErrOverflow for the one quotient that does not fit, the
// minimum of the type divided by −1 (Go's / would wrap it to the minimum; this function refuses).
//
// ID:      ALG-INTMTH-003
// Time:    O(w)   w = bit length of |dividend| (≤ 64); one compare-and-subtract per bit
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.3.1 "The Classical Algorithms", Algorithm D (long division), here in radix 2 where each quotient digit is found by one comparison
// Source:  Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, ch. 9 "Integer Division" (shift-and-subtract restoring division; truncating signed division from unsigned)
// Tags:    integer, division, long-division, shift-subtract, restoring-division, overflow-checked
func LongDivision[T signedInteger](dividend, divisor T) (quotient, remainder T, err error) {
	if divisor == 0 {
		return 0, 0, fmt.Errorf("%w: divisor = 0", ErrDivideByZero)
	}
	if divisor == -1 && dividend < 0 && -dividend < 0 { // only the minimum is its own negation
		return 0, 0, fmt.Errorf("%w: %d / -1", ErrOverflow, dividend)
	}
	n, d := magnitude(dividend), magnitude(divisor)
	var q, r uint64
	for i := bits.Len64(n) - 1; i >= 0; i-- {
		r = r<<1 | (n>>uint(i))&1 // r < d ≤ 2^63, so the shift cannot overflow
		if r >= d {
			r -= d
			q |= 1 << uint(i)
		}
	}
	quotient, remainder = T(q), T(r)
	if (dividend < 0) != (divisor < 0) {
		quotient = -quotient
	}
	if dividend < 0 {
		remainder = -remainder
	}
	return quotient, remainder, nil
}
