package bits

import (
	"fmt"
	"slices"
)

// maxAllocBytes is the largest single allocation this package makes: 2⁴⁷ bytes on 64-bit, 2³¹ on 32-bit,
// half the runtime's maxAlloc, above which makeslice panics (runtime/malloc.go; lead rule L41).
const maxAllocBytes = 1 << (31 + 16*(^uint(0)>>63))

// PopCountTable appends to dst the number of set bits of each of 0, 1, …, n (n+1 entries) and returns
// the extended slice. Each entry follows from one already written: popcount(i) = popcount(i >> 1) +
// (i & 1), so the table costs one addition per entry and no per-entry bit count.
// Errors: ErrNegative when n < 0; ErrTooLarge when the table (len(dst)+n+1 bytes)
// would exceed maxAllocBytes (2⁴⁷ bytes on 64-bit, 2³¹ on 32-bit; rule L41), never a makeslice panic. On error dst is returned with its
// length unchanged.
//
// ID:      ALG-BITS-009
// Time:    O(n)   n = the argument n; n+1 entries, one addition each
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques", sideways addition ν(x) and ν(x) = ν(⌊x/2⌋) + (x mod 2)
// Source:  Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, ch. 5 "Counting Bits", §5-1 "Counting 1-Bits"
// Tags:    bits, popcount, counting-bits, table, dynamic-programming
func PopCountTable(dst []uint8, n int) ([]uint8, error) {
	if n < 0 {
		return dst, fmt.Errorf("%w: n = %d", ErrNegative, n)
	}
	if int64(n) >= int64(maxAllocBytes)-int64(len(dst)) { // uint8 entries: bytes = len(dst) + n + 1; int64 so 32-bit builds compile and nothing overflows
		return dst, fmt.Errorf("%w: n = %d, %d existing + n+1 entries exceed %d bytes", ErrTooLarge, n, len(dst), int64(maxAllocBytes))
	}
	out := slices.Grow(dst, n+1)
	base := len(out)
	out = append(out, 0)
	for i := 1; i <= n; i++ {
		out = append(out, out[base+i>>1]+uint8(i&1))
	}
	return out, nil
}
