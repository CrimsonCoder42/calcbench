package bits

// LowestSetBit returns x with every bit cleared except its lowest set bit: x & −x (two's complement).
// LowestSetBit(0b10110) = 0b00010; LowestSetBit(0) = 0. For the bit's index use
// math/bits.TrailingZeros64(x) (64 for x = 0).
//
// ID:      ALG-BITS-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques", "Working with the rightmost bits" (x & −x extracts the rightmost 1)
// Source:  Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §2-1 "Manipulating Rightmost Bits" (x & (−x) isolates the rightmost 1-bit)
// Tags:    bits, lowest-set-bit, rightmost-one, isolate, fenwick
func LowestSetBit(x uint64) uint64 {
	return x & -x
}
