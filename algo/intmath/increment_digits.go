package intmath

import (
	"fmt"
	"slices"
)

// IncrementDigits appends to dst the decimal digits of the number one greater than digits and returns
// the extended slice. digits holds one value 0…9 per element, most significant first ([1 2 9] is 129,
// giving [1 3 0]). The result has the length of digits, or one more when every digit is 9 ([9 9] →
// [1 0 0]); leading zeros in the input are kept ([0 9] → [1 0]). digits is read, never modified.
// Errors: ErrEmpty when digits is empty; ErrInvalidDigit when an element is above 9. On error dst is
// returned with its length unchanged.
//
// ID:      ALG-INTMTH-016
// Time:    O(n)   n = len(digits)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.3.1 "The Classical Algorithms", Algorithm A (addition with carry), radix 10, with v = 1
// Tags:    integer, digits, increment, plus-one, carry, big-integer
func IncrementDigits(dst, digits []uint8) ([]uint8, error) {
	if len(digits) == 0 {
		return dst, fmt.Errorf("%w: digits is empty", ErrEmpty)
	}
	lastNonNine := -1
	for i, d := range digits {
		if d > 9 {
			return dst, fmt.Errorf("%w: digits[%d] = %d", ErrInvalidDigit, i, d)
		}
		if d != 9 {
			lastNonNine = i
		}
	}
	out := slices.Grow(dst, len(digits)+1)
	if lastNonNine < 0 { // all nines: the carry leaves the top digit
		out = append(out, 1)
		for range digits {
			out = append(out, 0)
		}
		return out, nil
	}
	out = append(out, digits[:lastNonNine]...)
	out = append(out, digits[lastNonNine]+1)
	for range len(digits) - lastNonNine - 1 {
		out = append(out, 0)
	}
	return out, nil
}
