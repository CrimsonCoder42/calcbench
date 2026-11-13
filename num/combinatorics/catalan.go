package combinatorics

import (
	"fmt"
	"math"
	"math/big"
)

// catalans holds C(n) for n = 0…36; C(37) = 45 950 804 324 621 742 364 exceeds 2⁶⁴ − 1.
// Values: OEIS A000108, checked in the tests against DLMF 26.5.1 in math/big.
var catalans = [...]uint64{
	1,
	1,
	2,
	5,
	14,
	42,
	132,
	429,
	1430,
	4862,
	16796,
	58786,
	208012,
	742900,
	2674440,
	9694845,
	35357670,
	129644790,
	477638700,
	1767263190,
	6564120420,
	24466267020,
	91482563640,
	343059613650,
	1289904147324,
	4861946401452,
	18367353072152,
	69533550916004,
	263747951750360,
	1002242216651368,
	3814986502092304,
	14544636039226909,
	55534064877048198,
	212336130412243110,
	812944042149730764,
	3116285494907301262,
	11959798385860453492,
}

// Catalan returns the n-th Catalan number C(n) = C(2n,n)/(n+1): among much else, the number of binary
// trees with n internal nodes and of balanced strings of n pairs of parentheses. C(36) =
// 11 959 798 385 860 453 492 is the largest that fits a uint64; past it use CatalanBig.
// Errors: ErrNegative when n < 0; ErrOverflow when n > 36.
//
// ID:      NUM-COMBI-012
// Time:    O(1)   (a 37-entry table)
// Space:   O(1), no allocation
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.5 eqs. 26.5.1, 26.5.4, Table 26.5.1
// Source:  OEIS A000108, https://oeis.org/A000108 (read 2026-11-12)
// Tags:    combinatorics, catalan, counting
// Units:   n: index (a count of pairs or nodes); result: a count of structures
// Example: Catalan(10) = 16 796 (DLMF Table 26.5.1)
func Catalan(n int) (uint64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n", ErrNegative)
	}
	if n >= len(catalans) {
		return 0, fmt.Errorf("%w: n", ErrOverflow)
	}
	return catalans[n], nil
}

// CatalanBig returns the n-th Catalan number C(2n,n)/(n+1) as a *big.Int: Catalan without the uint64
// limit.
// Errors: ErrNegative when n < 0; ErrOverflow when 2n exceeds math.MaxInt64 (the result would not fit in
// memory).
//
// ID:      NUM-COMBI-019
// Time:    O(n) big multiplies (one big.Int.Binomial(2n, n)) and one division, on numbers of at most b bits   b = bit length of the result, O(n)
// Space:   O(b) bits
// Source:  NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.5 eqs. 26.5.1, 26.5.4, Table 26.5.1
// Source:  OEIS A000108, https://oeis.org/A000108 (read 2026-11-12)
// Tags:    combinatorics, catalan, counting, big-int
// Units:   n: index (a count of pairs or nodes); result: a count of structures
// Example: CatalanBig(10) = 16 796 (DLMF Table 26.5.1)
func CatalanBig(n int) (*big.Int, error) {
	if n < 0 {
		return nil, fmt.Errorf("%w: n", ErrNegative)
	}
	if int64(n) > math.MaxInt64/2 {
		return nil, fmt.Errorf("%w: n", ErrOverflow)
	}
	r := new(big.Int).Binomial(2*int64(n), int64(n))
	return r.Quo(r, big.NewInt(int64(n)+1)), nil
}
