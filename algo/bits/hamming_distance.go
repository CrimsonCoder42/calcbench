package bits

import mathbits "math/bits"

// HammingDistance returns the number of bit positions in which x and y differ: the population count of
// x ⊕ y, from 0 to 64. The Hamming weight of one word is math/bits.OnesCount64 (not repeated here).
//
// ID:      ALG-BITS-013
// Time:    O(1)   one XOR and one population count
// Space:   O(1), no allocation
// Source:  Hamming, R. W., "Error detecting and error correcting codes", Bell System Technical Journal, 1950, 29(2):147–160, doi:10.1002/j.1538-7305.1950.tb00463.x (the distance between two code points)
// Source:  Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques", sideways addition ν(x ⊕ y)
// Tags:    bits, hamming-distance, popcount, xor, error-correcting-codes
func HammingDistance(x, y uint64) int {
	return mathbits.OnesCount64(x ^ y)
}
