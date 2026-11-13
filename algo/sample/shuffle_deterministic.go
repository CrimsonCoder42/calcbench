package sample

import (
	"fmt"
	"math/rand/v2"
	"slices"
)

// maxAllocBytes is the largest buffer this package allocates from a size argument: half the runtime's
// maxAlloc (runtime/malloc.go), above which makeslice panics. 2⁴⁷ on 64-bit, 2³¹ on 32-bit (rule L41).
const maxAllocBytes = 1 << (31 + 16*(^uint(0)>>63))

// pcgStream is the fixed second PCG word for ShuffleDeterministic ("calcbench" in ASCII hex, padded).
const pcgStream = 0x63616c6362656e63

// ShuffleDeterministic appends to dst a permutation of 0, 1, …, n−1 fixed by seed: the same seed gives the
// same permutation (within one Go toolchain; math/rand/v2 does not promise identical draws across Go
// releases), different seeds give independent-looking ones. Meant for tests and reproducible fixtures.
// NOT cryptographic: anyone who knows the seed knows the permutation.
// Errors: ErrNegative (n < 0), ErrTooLarge ((len(dst)+n)·8 bytes exceed maxAllocBytes: 2⁴⁷ on 64-bit,
// 2³¹ on 32-bit, above which make would panic). On error dst is returned unchanged.
//
// ID:      ALG-SAMPLE-007
// Time:    O(n)   n = the argument n, the length of the permutation
// Space:   O(1) extra beyond the generator; O(n) for the result, one allocation if dst lacks capacity
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §3.4.2, Algorithm P, via ShuffleInPlace (ALG-SAMPLE-001)
// Source:  O'Neill, "PCG: A Family of Simple Fast Space-Efficient Statistically Good Algorithms for Random Number Generation", Harvey Mudd College, HMC-CS-2014-0905, 2014 (the generator behind math/rand/v2 NewPCG)
// Tags:    shuffle, seeded, deterministic, testing, random-permutation
func ShuffleDeterministic(dst []int, n int, seed uint64) ([]int, error) {
	if n < 0 {
		return dst, fmt.Errorf("%w: n = %d", ErrNegative, n)
	}
	if n > maxAllocBytes/8-len(dst) {
		return dst, fmt.Errorf("%w: n = %d", ErrTooLarge, n)
	}
	base := len(dst)
	dst = slices.Grow(dst, n)
	for i := range n {
		dst = append(dst, i)
	}
	// rng is never nil here, so ShuffleInPlace cannot fail.
	_ = ShuffleInPlace(rand.New(rand.NewPCG(seed, pcgStream)), dst[base:])
	return dst, nil
}
