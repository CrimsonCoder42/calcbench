package combinatorics

import (
	"fmt"
	"math/big"
	"math/bits"
)

// StirlingSecond returns the Stirling number of the second kind {n k}: the number of ways to partition an
// n-element set into k non-empty blocks. {0 0} = 1; {n k} = 0 when k > n, or k = 0 < n. Past uint64 use
// StirlingSecondBig.
// Errors: ErrNegative when n < 0 or k < 0; ErrOverflow when {n k} > 2⁶⁴ − 1.
//
// ID:      NUM-COMBI-014
// Time:    O(n·w)   w = min(n−k+1, 64) (one band of the triangle; the loop also stops at the first overflow)
// Space:   O(1), no allocation (a fixed 64-entry array)
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.8 eqs. 26.8.6, 26.8.22, Table 26.8.2
// Source:  Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §6.1 "Stirling Numbers" (table of {n k}, p. 244)
// Source:  OEIS A008277, https://oeis.org/A008277 (read 2026-11-12)
// Tags:    combinatorics, stirling, set-partitions
// Units:   n: a count of elements; k: a count of blocks; result: a count of partitions
// Example: StirlingSecond(6, 3) = 90, StirlingSecond(10, 5) = 42 525 (DLMF Table 26.8.2)
func StirlingSecond(n, k int) (uint64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n", ErrNegative)
	}
	if k < 0 {
		return 0, fmt.Errorf("%w: k", ErrNegative)
	}
	switch {
	case k > n:
		return 0, nil
	case k == n:
		return 1, nil
	case k == 0:
		return 0, nil
	case k == 1:
		return 1, nil
	case k == n-1:
		// One block of two, the rest singletons: C(n,2).
		r, ok := binomialU64(uint64(n), 2)
		if !ok {
			return 0, fmt.Errorf("%w: {n k}", ErrOverflow)
		}
		return r, nil
	case n-k >= stirlingBandMax:
		return 0, fmt.Errorf("%w: {n k}", ErrOverflow)
	}
	// Band indexed by u = i − j ∈ [0, d], d = n − k. DLMF 26.8.22, {i j} = j·{i−1 j} + {i−1 j−1}, becomes
	// a[u] += (i−u)·a[u−1]. Each band entry is ≤ {n k}, so an overflow is real.
	d := n - k
	var a [stirlingBandMax]uint64
	a[0] = 1
	for i := 1; i <= n; i++ {
		lo, hi := max(1, i-k), min(d, i-1)
		for u := hi; u >= lo; u-- {
			p1, p0 := bits.Mul64(uint64(i-u), a[u-1])
			s, carry := bits.Add64(a[u], p0, 0)
			if p1 != 0 || carry != 0 {
				return 0, fmt.Errorf("%w: {n k}", ErrOverflow)
			}
			a[u] = s
		}
	}
	return a[d], nil
}

// StirlingSecondBig returns {n k}, the number of partitions of an n-element set into k non-empty blocks,
// as a *big.Int: StirlingSecond without the uint64 limit.
// Errors: ErrNegative when n < 0 or k < 0.
//
// ID:      NUM-COMBI-021
// Time:    O(n·w) big multiply-adds   w = min(k, n−k) + 1; each on numbers of at most b bits, b = bit length of the result
// Space:   O(w·b) bits (w big integers, one band of the triangle)
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.8 eqs. 26.8.6, 26.8.22, Table 26.8.2
// Source:  Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §6.1 "Stirling Numbers" (table of {n k}, p. 244)
// Source:  OEIS A008277, https://oeis.org/A008277 (read 2026-11-12)
// Tags:    combinatorics, stirling, set-partitions, big-int
// Units:   n: a count of elements; k: a count of blocks; result: a count of partitions
// Example: StirlingSecondBig(10, 5) = 42 525 (DLMF Table 26.8.2)
func StirlingSecondBig(n, k int) (*big.Int, error) {
	if n < 0 {
		return nil, fmt.Errorf("%w: n", ErrNegative)
	}
	if k < 0 {
		return nil, fmt.Errorf("%w: k", ErrNegative)
	}
	if k > n || (k == 0 && n > 0) {
		return new(big.Int), nil
	}
	if k == n {
		return big.NewInt(1), nil
	}
	return stirlingBand(n, k, true), nil
}
