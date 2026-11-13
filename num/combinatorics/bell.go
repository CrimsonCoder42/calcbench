package combinatorics

import (
	"fmt"
	"math/big"
)

// bells holds B(n) for n = 0…25; B(26) = 49 631 246 523 618 756 274 exceeds 2⁶⁴ − 1.
// Values: OEIS A000110, checked in the tests against DLMF 26.7.2 in math/big.
var bells = [...]uint64{
	1,
	1,
	2,
	5,
	15,
	52,
	203,
	877,
	4140,
	21147,
	115975,
	678570,
	4213597,
	27644437,
	190899322,
	1382958545,
	10480142147,
	82864869804,
	682076806159,
	5832742205057,
	51724158235372,
	474869816156751,
	4506715738447323,
	44152005855084346,
	445958869294805289,
	4638590332229999353,
}

// Bell returns the Bell number B(n) = Σₖ {n k}: the number of partitions of an n-element set into
// non-empty blocks. B(25) = 4 638 590 332 229 999 353 is the largest that fits a uint64; past it use
// BellBig.
// Errors: ErrNegative when n < 0; ErrOverflow when n > 25.
//
// ID:      NUM-COMBI-015
// Time:    O(1)   (a 26-entry table)
// Space:   O(1), no allocation
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.7 eqs. 26.7.2, 26.7.6, Table 26.7.1
// Source:  OEIS A000110, https://oeis.org/A000110 (read 2026-11-12)
// Tags:    combinatorics, bell, set-partitions
// Units:   n: a count of elements; result: a count of partitions
// Example: Bell(5) = 52, Bell(10) = 115 975 (DLMF Table 26.7.1)
func Bell(n int) (uint64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n", ErrNegative)
	}
	if n >= len(bells) {
		return 0, fmt.Errorf("%w: n", ErrOverflow)
	}
	return bells[n], nil
}

// BellBig returns the Bell number B(n) as a *big.Int: Bell without the uint64 limit. It sums one row of
// the Stirling triangle, B(n) = Σₖ {n k} (DLMF 26.7.2), built with {i j} = j·{i−1 j} + {i−1 j−1}
// (DLMF 26.8.22).
// Errors: ErrNegative when n < 0.
//
// ID:      NUM-COMBI-022
// Time:    O(n²) big multiply-adds by a small integer, each on numbers of at most b bits   b = bit length of the result, O(n log n)
// Space:   O(n·b) bits (one row of n+1 big integers)
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.7 eqs. 26.7.2, 26.7.6, Table 26.7.1
// Source:  OEIS A000110, https://oeis.org/A000110 (read 2026-11-12)
// Tags:    combinatorics, bell, set-partitions, big-int
// Units:   n: a count of elements; result: a count of partitions
// Example: BellBig(26) = 49 631 246 523 618 756 274 (OEIS A000110)
func BellBig(n int) (*big.Int, error) {
	if n < 0 {
		return nil, fmt.Errorf("%w: n", ErrNegative)
	}
	row := make([]big.Int, n+1)
	row[0].SetInt64(1)
	var t, f big.Int
	for i := 1; i <= n; i++ {
		for j := i; j >= 1; j-- {
			t.Mul(f.SetInt64(int64(j)), &row[j])
			row[j].Add(&t, &row[j-1])
		}
		row[0].SetInt64(0)
	}
	sum := new(big.Int)
	for j := range row {
		sum.Add(sum, &row[j])
	}
	return sum, nil
}
