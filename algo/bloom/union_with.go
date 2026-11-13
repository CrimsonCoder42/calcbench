package bloom

import (
	"fmt"
	"math/bits"
)

// UnionWith sets in f every bit set in other, so f then represents the union of both filters' items
// (no false negatives for either). Both must have the same bit count and hash count; the hash function
// is the same for every filter of this package. f is modified; other is only read. f.UnionWith(f)
// changes nothing.
// Errors: ErrIncompatible when other is nil or its bit count or hash count differs from f's; f is then
// unchanged.
//
// ID:      ALG-BLOOM-004
// Time:    O(m/64)   m = bit count
// Space:   O(1) extra
// Source:  Broder and Mitzenmacher, "Network applications of Bloom filters: a survey", Internet Mathematics 1(4):485–509, 2004, doi:10.1080/15427951.2004.10129096, §2.4 (the OR of two filters with the same bits and hash functions represents the union)
// Tags:    bloom-filter, union, set-operations, merge
func (f *BloomFilter) UnionWith(other *BloomFilter) error {
	if other == nil {
		return fmt.Errorf("%w: other is nil", ErrIncompatible)
	}
	if other.bitCount != f.bitCount || other.hashCount != f.hashCount {
		return fmt.Errorf("%w: (m, k) = (%d, %d) and (%d, %d)", ErrIncompatible, f.bitCount, f.hashCount, other.bitCount, other.hashCount)
	}
	var ones uint64
	for i, w := range other.words {
		f.words[i] |= w
		ones += uint64(bits.OnesCount64(f.words[i]))
	}
	f.setBits = ones
	return nil
}
