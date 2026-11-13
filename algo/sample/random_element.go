package sample

import (
	"fmt"
	"iter"
	"math/rand/v2"
)

// RandomElement returns one item of seq, read once, chosen uniformly: with N items each is returned
// with probability 1/N, and N need not be known in advance (a linked list's values, the lines of a
// stream). It is reservoir sampling with a reservoir of one.
// Errors: ErrNilRand, ErrNilSeq, ErrEmpty (seq yields nothing).
//
// ID:      ALG-SAMPLE-006
// Time:    O(N)   N = number of items in seq; N − 1 random draws
// Space:   O(1) extra
// Source:  Vitter, "Random sampling with a reservoir", ACM Transactions on Mathematical Software 11(1):37–57, 1985, doi:10.1145/3147.3165, Algorithm R with k = 1, via ReservoirSample (ALG-SAMPLE-002)
// Tags:    reservoir-sampling, streaming, random-element, linked-list, iterator
func RandomElement[T any](rng *rand.Rand, seq iter.Seq[T]) (T, error) {
	var zero T
	var slot [1]T
	got, err := ReservoirSample(slot[:0], rng, seq, 1)
	if err != nil {
		return zero, err
	}
	if len(got) == 0 {
		return zero, fmt.Errorf("%w: seq", ErrEmpty)
	}
	return got[0], nil
}
