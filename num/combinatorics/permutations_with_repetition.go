package combinatorics

import (
	"fmt"
	"math/bits"
)

// PermutationsWithRepetition returns nᵏ exactly: the number of length-k sequences over n symbols, each
// symbol usable any number of times. 0⁰ = 1 (the empty sequence). Past uint64 use new(big.Int).Exp.
// Errors: ErrNegative when n < 0 or k < 0; ErrOverflow when nᵏ > 2⁶⁴ − 1.
//
// ID:      NUM-COMBI-006
// Time:    O(log k)   (square and multiply)
// Space:   O(1), no allocation
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.17 Table 26.17.1 (unrestricted maps of labelled balls into labelled boxes)
// Tags:    combinatorics, permutations, repetition, counting
// Units:   n: number of symbols; k: sequence length; result: a count of sequences
// Example: PermutationsWithRepetition(26, 3) = 17 576 (three-letter codes)
func PermutationsWithRepetition(n, k int) (uint64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n", ErrNegative)
	}
	if k < 0 {
		return 0, fmt.Errorf("%w: k", ErrNegative)
	}
	if n <= 1 || k == 0 {
		if k == 0 {
			return 1, nil
		}
		return uint64(n), nil // 0ᵏ = 0 and 1ᵏ = 1 for k ≥ 1
	}
	result, base, e := uint64(1), uint64(n), uint64(k)
	for {
		if e&1 == 1 {
			hi, lo := bits.Mul64(result, base)
			if hi != 0 {
				return 0, fmt.Errorf("%w: n^k", ErrOverflow)
			}
			result = lo
		}
		e >>= 1
		if e == 0 {
			return result, nil
		}
		// A higher bit of k is still set, so the result will include a factor ≥ base²: an overflow here
		// is real.
		hi, lo := bits.Mul64(base, base)
		if hi != 0 {
			return 0, fmt.Errorf("%w: n^k", ErrOverflow)
		}
		base = lo
	}
}
