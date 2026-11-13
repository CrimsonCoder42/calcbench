package intmath

import (
	"fmt"
	"math"
)

// ReverseDigits returns the integer whose decimal digits are those of x in reverse order, keeping the
// sign: 1234 → 4321, −560 → −65 (trailing zeros become leading zeros and vanish), 0 → 0.
// Errors: ErrOverflow when the reversal does not fit in T (for int32, 1_000_000_009 → 9_000_000_001).
//
// ID:      ALG-INTMTH-004
// Time:    O(d)   d = number of decimal digits of x (≤ 19)
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.4 "Radix Conversion", Method (1a) (the digits of an integer from least significant by repeated division by the radix) and Method (1b) (rebuilding an integer by Horner's rule u·10 + digit)
// Tags:    integer, digits, reverse, decimal, overflow-checked
func ReverseDigits[T signedInteger](x T) (T, error) {
	m := magnitude(x)
	var r uint64 // at most 19 digits, so r < 10^19 < 2^64: the accumulation itself never overflows
	for ; m > 0; m /= 10 {
		r = r*10 + m%10
	}
	var v int64
	if x < 0 {
		if r > 1<<63 {
			return 0, fmt.Errorf("%w: reverse of %d", ErrOverflow, x)
		}
		v = -int64(r)
	} else {
		if r > math.MaxInt64 {
			return 0, fmt.Errorf("%w: reverse of %d", ErrOverflow, x)
		}
		v = int64(r)
	}
	if int64(T(v)) != v {
		return 0, fmt.Errorf("%w: reverse of %d", ErrOverflow, x)
	}
	return T(v), nil
}
