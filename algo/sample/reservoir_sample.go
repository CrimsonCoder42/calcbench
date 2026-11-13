package sample

import (
	"fmt"
	"iter"
	"math/rand/v2"
)

// ReservoirSample appends to dst a uniform random sample of k items from seq, read once, whose length
// need not be known: every k-subset of the items is equally likely. When seq has fewer than k items,
// all of them are appended. The appended items are in reservoir-slot order, not a random order (shuffle
// them if order matters). k = 0 appends nothing and does not read seq.
// Errors: ErrNilRand, ErrNilSeq, ErrNegative (k < 0). On error dst is returned unchanged.
//
// ID:      ALG-SAMPLE-002
// Time:    O(N)   N = number of items in seq; N − k random draws
// Space:   O(1) extra; O(min(k, N)) for the result, grown by append (never sized from k, so a huge k allocates nothing extra)
// Source:  Vitter, "Random sampling with a reservoir", ACM Transactions on Mathematical Software 11(1):37–57, 1985, doi:10.1145/3147.3165, Algorithm R
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §3.4.2, Algorithm R (reservoir sampling)
// Tags:    reservoir-sampling, streaming, random-sample, iterator
func ReservoirSample[T any](dst []T, rng *rand.Rand, seq iter.Seq[T], k int) ([]T, error) {
	switch {
	case rng == nil:
		return dst, fmt.Errorf("%w: rng", ErrNilRand)
	case seq == nil:
		return dst, fmt.Errorf("%w: seq", ErrNilSeq)
	case k < 0:
		return dst, fmt.Errorf("%w: k = %d", ErrNegative, k)
	case k == 0:
		return dst, nil
	}
	base := len(dst)
	t := 0 // items seen so far
	for v := range seq {
		t++
		if t <= k {
			dst = append(dst, v)
			continue
		}
		// The t-th item enters with probability k/t, into a uniform slot.
		if j := rng.IntN(t); j < k {
			dst[base+j] = v
		}
	}
	return dst, nil
}
