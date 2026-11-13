package bits

import (
	"fmt"
	mathbits "math/bits"
)

// ReverseBits returns the low width bits of x in reverse order: bit i of the result is bit width−1−i of
// x, for i < width, and every bit at or above width is zero. ReverseBits(0b0011, 4) = 0b1100. For width
// = 64 it equals math/bits.Reverse64.
// Errors: ErrWidth when width is outside 1…64; ErrOutOfRange when x has a set bit at or above width
// (it would otherwise be lost silently).
//
// ID:      ALG-BITS-004
// Time:    O(1)   one 64-bit reversal (math/bits) and a shift
// Space:   O(1), no allocation
// Source:  Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §7-1 "Reversing Bits and Bytes" (reverse the word, then shift right by w − width to reverse a field)
// Source:  Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques", bit reversal by swaps
// Tags:    bits, reverse, bit-reversal, fft-index
func ReverseBits(x uint64, width int) (uint64, error) {
	if width < 1 || width > 64 {
		return 0, fmt.Errorf("%w: width = %d", ErrWidth, width)
	}
	if width < 64 && x>>uint(width) != 0 {
		return 0, fmt.Errorf("%w: x = %#x, width = %d", ErrOutOfRange, x, width)
	}
	return mathbits.Reverse64(x) >> uint(64-width), nil
}
