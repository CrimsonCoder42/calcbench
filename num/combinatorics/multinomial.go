package combinatorics

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
)

// Multinomial returns the multinomial coefficient (c₁+…+cₘ)!/(c₁!·…·cₘ!) exactly: the number of distinct
// orderings of a multiset with cᵢ copies of item i, or of ways to split Σcᵢ labelled objects into boxes
// of those sizes. It is formed as the product of binomials Πⱼ C(sⱼ, cⱼ), sⱼ = c₁+…+cⱼ (DLMF 26.4.2), so
// (Σcᵢ)! is never formed and every partial product is ≤ the result: an overflow is never false. Past
// uint64 use MultinomialBig. counts is read, never modified.
// Errors: ErrEmpty when counts is empty; ErrNegative when a count is negative (reported before any
// overflow); ErrOverflow when the result > 2⁶⁴ − 1.
//
// ID:      NUM-COMBI-008
// Time:    O(m + Σⱼ min(cⱼ, sⱼ₋₁, 64)·log s)   m = len(counts), sⱼ = c₁+…+cⱼ, s = sₘ
// Space:   O(1), no allocation
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.4 eq. 26.4.2
// Tags:    combinatorics, multinomial, multiset, permutations, counting
// Units:   counts: copies of each item; result: a count of orderings
// Example: Multinomial([]int{1, 4, 4, 2}) = 34 650 (the letters of MISSISSIPPI)
func Multinomial(counts []int) (uint64, error) {
	if len(counts) == 0 {
		return 0, fmt.Errorf("%w: counts", ErrEmpty)
	}
	for i, c := range counts {
		if c < 0 {
			return 0, fmt.Errorf("%w: counts[%d]", ErrNegative, i)
		}
	}
	var sum uint64
	r := uint64(1)
	for _, c := range counts {
		var carry uint64
		sum, carry = bits.Add64(sum, uint64(c), 0)
		if carry != 0 {
			return 0, fmt.Errorf("%w: counts", ErrOverflow)
		}
		b, ok := binomialU64(sum, uint64(c))
		if !ok {
			return 0, fmt.Errorf("%w: counts", ErrOverflow)
		}
		hi, lo := bits.Mul64(r, b)
		if hi != 0 {
			return 0, fmt.Errorf("%w: counts", ErrOverflow)
		}
		r = lo
	}
	return r, nil
}

// MultinomialBig returns the multinomial coefficient (c₁+…+cₘ)!/(c₁!·…·cₘ!) as a *big.Int: Multinomial
// without the uint64 limit, formed as Πⱼ C(sⱼ, cⱼ) (DLMF 26.4.2). counts is read, never modified.
// Errors: ErrEmpty when counts is empty; ErrNegative when a count is negative; ErrOverflow when Σcᵢ
// exceeds math.MaxInt64 (the result would not fit in memory).
//
// ID:      NUM-COMBI-017
// Time:    O(m) big.Int.Binomial calls and multiplies   m = len(counts); each on numbers of at most b bits, b = bit length of the result
// Space:   O(b) bits
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.4 eq. 26.4.2
// Tags:    combinatorics, multinomial, multiset, permutations, counting, big-int
// Units:   counts: copies of each item; result: a count of orderings
// Example: MultinomialBig([]int{1, 4, 4, 2}) = 34 650 (the letters of MISSISSIPPI)
func MultinomialBig(counts []int) (*big.Int, error) {
	if len(counts) == 0 {
		return nil, fmt.Errorf("%w: counts", ErrEmpty)
	}
	var sum int64
	for i, c := range counts {
		if c < 0 {
			return nil, fmt.Errorf("%w: counts[%d]", ErrNegative, i)
		}
		if int64(c) > math.MaxInt64-sum {
			return nil, fmt.Errorf("%w: sum of counts", ErrOverflow)
		}
		sum += int64(c)
	}
	r := big.NewInt(1)
	var b big.Int
	sum = 0
	for _, c := range counts {
		sum += int64(c)
		r.Mul(r, b.Binomial(sum, int64(c)))
	}
	return r, nil
}
