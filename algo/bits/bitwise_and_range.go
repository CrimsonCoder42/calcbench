package bits

import (
	"fmt"
	mathbits "math/bits"
)

// BitwiseAndRange returns lo & (lo+1) & … & hi without visiting the range: the common high-order bits
// of lo and hi, with every bit from their highest differing bit downward cleared. Proof in
// doc/algorithms_core_2026-11-12/catalog/bits.md (ALG-BITS-012).
// Errors: ErrRange when lo > hi.
//
// ID:      ALG-BITS-012
// Time:    O(1)   one leading-zero count
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques" (bitwise AND; the leftmost bit of x ⊕ y marks the first position where x and y differ); proof of the range identity in catalog/bits.md
// Tags:    bits, and, range, common-prefix
func BitwiseAndRange(lo, hi uint64) (uint64, error) {
	if lo > hi {
		return 0, fmt.Errorf("%w: lo = %d, hi = %d", ErrRange, lo, hi)
	}
	shift := uint(64 - mathbits.LeadingZeros64(lo^hi)) // 0 when lo == hi; a shift of 64 yields 0
	return lo >> shift << shift, nil
}
