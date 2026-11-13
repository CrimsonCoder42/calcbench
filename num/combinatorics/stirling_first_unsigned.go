package combinatorics

import (
	"fmt"
	"math/big"
	"math/bits"
)

// stirlingBandMax bounds n−k for the uint64 Stirling numbers: both kinds are ≥ k^(n−k) (put 1…k in
// separate blocks or cycles, then the other n−k anywhere), so k ≥ 2 and n−k ≥ 64 is an overflow.
const stirlingBandMax = 64

// StirlingFirstUnsigned returns the unsigned Stirling number of the first kind [n k] = |s(n,k)|: the
// number of permutations of n objects with exactly k cycles. [0 0] = 1; [n k] = 0 when k > n, or k = 0 < n.
// DLMF's s(n,k) is signed, s(n,k) = (−1)^(n−k)·[n k]; the sign is the caller's (an int64 would halve the
// range). Past uint64 use StirlingFirstUnsignedBig.
// Errors: ErrNegative when n < 0 or k < 0; ErrOverflow when [n k] > 2⁶⁴ − 1.
//
// ID:      NUM-COMBI-013
// Time:    O(n·w)   w = min(n−k+1, 64) (one band of the triangle; the loop also stops at the first overflow)
// Space:   O(1), no allocation (a fixed 64-entry array)
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.8 eqs. 26.8.3, 26.8.18, Table 26.8.1
// Source:  Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §6.1 "Stirling Numbers" (table of [n k], p. 245)
// Source:  OEIS A132393, https://oeis.org/A132393 (read 2026-11-12)
// Tags:    combinatorics, stirling, cycles, permutations
// Units:   n: a count of objects; k: a count of cycles; result: a count of permutations
// Example: StirlingFirstUnsigned(6, 3) = 225, StirlingFirstUnsigned(10, 5) = 269 325 (DLMF Table 26.8.1: s = −225, −269 325)
func StirlingFirstUnsigned(n, k int) (uint64, error) {
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
	case k == n-1:
		// One 2-cycle and n−2 fixed points: C(n,2).
		r, ok := binomialU64(uint64(n), 2)
		if !ok {
			return 0, fmt.Errorf("%w: [n k]", ErrOverflow)
		}
		return r, nil
	case k == 1:
		// One n-cycle: (n−1)!.
		r, ok := mulFactorial(1, uint64(n-1), 1)
		if !ok {
			return 0, fmt.Errorf("%w: [n k]", ErrOverflow)
		}
		return r, nil
	case n-k >= stirlingBandMax:
		return 0, fmt.Errorf("%w: [n k]", ErrOverflow)
	}
	// Band of the triangle indexed by u = i − j ∈ [0, d], d = n − k. The recurrence (unsigned form of
	// DLMF 26.8.18) [i j] = [i−1 j−1] + (i−1)·[i−1 j] becomes a[u] += (i−1)·a[u−1]. Each band entry is
	// ≤ [n k] (both steps of the recurrence are non-decreasing on the band), so an overflow is real.
	d := n - k
	var a [stirlingBandMax]uint64
	a[0] = 1
	for i := 1; i <= n; i++ {
		lo, hi := max(1, i-k), min(d, i-1)
		for u := hi; u >= lo; u-- {
			p1, p0 := bits.Mul64(uint64(i-1), a[u-1])
			s, carry := bits.Add64(a[u], p0, 0)
			if p1 != 0 || carry != 0 {
				return 0, fmt.Errorf("%w: [n k]", ErrOverflow)
			}
			a[u] = s
		}
	}
	return a[d], nil
}

// StirlingFirstUnsignedBig returns [n k], the number of permutations of n objects with exactly k cycles,
// as a *big.Int: StirlingFirstUnsigned without the uint64 limit.
// Errors: ErrNegative when n < 0 or k < 0.
//
// ID:      NUM-COMBI-020
// Time:    O(n·w) big multiply-adds   w = min(k, n−k) + 1; each on numbers of at most b bits, b = bit length of the result
// Space:   O(w·b) bits (w big integers, one band of the triangle)
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.8 eqs. 26.8.3, 26.8.18, Table 26.8.1
// Source:  Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §6.1 "Stirling Numbers" (table of [n k], p. 245)
// Source:  OEIS A132393, https://oeis.org/A132393 (read 2026-11-12)
// Tags:    combinatorics, stirling, cycles, permutations, big-int
// Units:   n: a count of objects; k: a count of cycles; result: a count of permutations
// Example: StirlingFirstUnsignedBig(10, 5) = 269 325 (DLMF Table 26.8.1)
func StirlingFirstUnsignedBig(n, k int) (*big.Int, error) {
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
	return stirlingBand(n, k, false), nil
}

// stirlingBand returns [n k] (second == false) or {n k} (second == true) for 1 ≤ k < n, computing one
// band of the triangle in w = min(k, n−k) + 1 big integers. When k ≤ n−k the band is indexed by j (the
// column); otherwise by u = i − j. Recurrences: [i j] = [i−1 j−1] + (i−1)·[i−1 j] (DLMF 26.8.18,
// unsigned) and {i j} = j·{i−1 j} + {i−1 j−1} (DLMF 26.8.22).
func stirlingBand(n, k int, second bool) *big.Int {
	d := n - k
	var t, f big.Int
	if k <= d {
		row := make([]big.Int, k+1)
		row[0].SetInt64(1)
		for i := 1; i <= n; i++ {
			lo, hi := max(1, i-d), min(i, k)
			for j := hi; j >= lo; j-- {
				// row[j] holds column j of row i−1, row[j−1] column j−1.
				if second {
					f.SetInt64(int64(j))
				} else {
					f.SetInt64(int64(i - 1))
				}
				t.Mul(&f, &row[j])
				row[j].Add(&t, &row[j-1])
			}
			if i == 1 {
				row[0].SetInt64(0) // [i 0] = {i 0} = 0 for i ≥ 1
			}
		}
		return new(big.Int).Set(&row[k])
	}
	a := make([]big.Int, d+1)
	a[0].SetInt64(1)
	for i := 1; i <= n; i++ {
		lo, hi := max(1, i-k), min(d, i-1)
		for u := hi; u >= lo; u-- {
			// a[u] holds (i−1, j−1), a[u−1] holds (i−1, j), with j = i − u.
			if second {
				f.SetInt64(int64(i - u))
			} else {
				f.SetInt64(int64(i - 1))
			}
			t.Mul(&f, &a[u-1])
			a[u].Add(&a[u], &t)
		}
	}
	return new(big.Int).Set(&a[d])
}
