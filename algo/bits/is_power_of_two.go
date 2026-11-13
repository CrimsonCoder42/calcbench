package bits

// integer is the set of Go's built-in integer types and the types defined on them.
type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// IsPowerOfTwo reports whether x is 2^k for some k ≥ 0 (1, 2, 4, 8, …). Zero and negative values are
// not powers of two. For a signed type the minimum value (−2^(w−1)) is negative and so returns false.
// x & (x−1) clears the lowest set bit of x; the result is zero exactly when x has one set bit.
//
// ID:      ALG-BITS-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §2-1 "Manipulating Rightmost Bits" (x & (x−1) turns off the rightmost 1-bit; "can be used to determine if an unsigned integer is a power of 2")
// Source:  Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques", "Working with the rightmost bits"
// Tags:    bits, power-of-two, bit-trick
func IsPowerOfTwo[T integer](x T) bool {
	return x > 0 && x&(x-1) == 0
}

// powerOfFourBits has a 1 in every even bit position 0, 2, 4, …, 62: the positions 4^k occupies.
const powerOfFourBits = 0x5555_5555_5555_5555

// IsPowerOfFour reports whether x is 4^k for some k ≥ 0 (1, 4, 16, 64, …): x is a power of two and its
// one set bit is in an even position. Zero and negative values return false.
//
// ID:      ALG-BITS-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §2-1 "Manipulating Rightmost Bits" (power-of-two test); 4^k = 2^(2k), so the single bit sits at an even index (mask test)
// Tags:    bits, power-of-four, bit-trick
func IsPowerOfFour[T integer](x T) bool {
	return IsPowerOfTwo(x) && uint64(x)&powerOfFourBits != 0
}
