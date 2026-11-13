package bits

import mathbits "math/bits"

// HighestSetBit returns x with every bit cleared except its highest set bit: the greatest power of two
// that is ≤ x, or 0 for x = 0. HighestSetBit(0b10110) = 0b10000. For the bit's index use
// math/bits.Len64(x) − 1.
//
// ID:      ALG-BITS-008
// Time:    O(1)   one leading-zero count
// Space:   O(1), no allocation
// Source:  Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §3-2 "Rounding Up/Down to the Next Power of 2" (flp2: the greatest power of 2 ≤ x, 0 for x = 0)
// Source:  Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques", the leftmost bit 2^⌊lg x⌋
// Tags:    bits, highest-set-bit, leftmost-one, floor-power-of-two
func HighestSetBit(x uint64) uint64 {
	if x == 0 {
		return 0
	}
	return 1 << (63 - mathbits.LeadingZeros64(x))
}
