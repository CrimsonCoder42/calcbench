package sample

import (
	"errors"
	"fmt"
	"math/rand/v2"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrNilRand: the *rand.Rand argument is nil.
	ErrNilRand = errors.New("sample: random source is nil")
	// ErrNilSeq: the sequence argument is nil.
	ErrNilSeq = errors.New("sample: sequence is nil")
	// ErrNilFunc: a function argument is nil.
	ErrNilFunc = errors.New("sample: function is nil")
	// ErrNegative: a count is negative.
	ErrNegative = errors.New("sample: count is negative")
	// ErrNonPositive: a count that must be at least 1 is not.
	ErrNonPositive = errors.New("sample: count must be positive")
	// ErrEmpty: the sequence has no items.
	ErrEmpty = errors.New("sample: sequence is empty")
	// ErrMaxTries: rejection sampling rejected every one of its allowed tries.
	ErrMaxTries = errors.New("sample: no candidate accepted within maxTries")
	// ErrTooLarge: a requested length would exceed maxAllocBytes.
	ErrTooLarge = errors.New("sample: length too large for a slice")
)

// ShuffleInPlace puts values in a uniformly random order, each of the n! orders equally likely given a
// uniform source, and modifies values in place. A source with s bits of state reaches at most 2^s
// orders (PCG has 128, so not every order of 35 or more items is reachable).
// Errors: ErrNilRand when rng is nil (values is then unchanged).
//
// ID:      ALG-SAMPLE-001
// Time:    O(n)   n = len(values); n−1 random draws and swaps
// Space:   O(1) extra, no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §3.4.2, Algorithm P (shuffling)
// Source:  Durstenfeld, "Algorithm 235: Random permutation", Communications of the ACM 7(7):420, 1964, doi:10.1145/364520.364540
// Tags:    shuffle, fisher-yates, knuth-shuffle, random-permutation, in-place
func ShuffleInPlace[T any](rng *rand.Rand, values []T) error {
	if rng == nil {
		return fmt.Errorf("%w: rng", ErrNilRand)
	}
	for j := len(values) - 1; j > 0; j-- {
		k := rng.IntN(j + 1)
		values[j], values[k] = values[k], values[j]
	}
	return nil
}
