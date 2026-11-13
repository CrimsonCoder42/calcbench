package bits

import "iter"

// Submasks returns an iterator over every submask s of mask (every s with s & ^mask == 0), in
// descending numeric order from mask itself down to 0: 2^k values for a mask with k set bits. The next
// submask after s is (s − 1) & mask: subtracting 1 borrows through the low zeros of s and the AND drops
// the bits outside mask. Stopping early (break) is supported.
//
// ID:      ALG-BITS-014
// Time:    O(2^k)   k = number of set bits in mask; O(1) per submask
// Space:   O(1), no allocation per submask
// Source:  Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques" (running through all subsets of a set represented by a mask with x ← (x − χ) & χ)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, ch. 10 "Bit manipulation", §10.3 "Representing sets", "Iterating through subsets", https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    bits, submask, subset-enumeration, iterator, bitmask-dp
func Submasks(mask uint64) iter.Seq[uint64] {
	return func(yield func(uint64) bool) {
		s := mask
		for {
			if !yield(s) || s == 0 {
				return
			}
			s = (s - 1) & mask
		}
	}
}
