package enumerate

import (
	"fmt"
	"iter"
)

// GrayCodeSequence returns an iterator over the 2^bits words of the binary reflected Gray code on bits
// bits, g(k) = k XOR ⌊k/2⌋ for k = 0, 1, …, 2^bits − 1: consecutive words differ in exactly one bit, and
// the last differs from the first in one bit. For 3: 0 1 3 2 6 7 5 4. bits = 0 yields 0 once.
// Errors: ErrBitsRange when bits < 0 or bits > 64.
//
// ID:      ALG-ENUM-009
// Time:    O(1) per word, O(2^b) in all   b = bits
// Space:   O(1) extra, no allocation per word
// Source:  Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.1 (the binary reflected Gray code, g(k) = k ⊕ ⌊k/2⌋)
// Source:  Gray, "Pulse code communication", US Patent 2,632,058, 1953
// Tags:    gray-code, bits, enumeration, iterator, reflected-binary
func GrayCodeSequence(bits int) (iter.Seq[uint64], error) {
	if bits < 0 || bits > 64 {
		return nil, fmt.Errorf("%w: bits = %d", ErrBitsRange, bits)
	}
	last := ^uint64(0) >> (64 - bits) // 2^bits − 1; a shift by 64 gives 0
	return func(yield func(uint64) bool) {
		for k := uint64(0); ; k++ {
			if !yield(k^(k>>1)) || k == last {
				return
			}
		}
	}, nil
}
