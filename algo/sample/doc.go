// Package sample shuffles and samples with a random source supplied by the caller: Fisher–Yates
// shuffles, reservoir sampling of a stream, weighted picks, rejection sampling and a seeded shuffle for
// tests.
//
// ID short form: SAMPLE (IDs are ALG-SAMPLE-NNN).
//
// Randomness: every function takes a *rand.Rand from math/rand/v2 (or, for ShuffleDeterministic, a
// seed). There is no global source, no clock and no goroutine; a nil *rand.Rand is ErrNilRand, never a
// panic. The results are uniform exactly when the source is (IntN and Uint64N are unbiased). Nothing here
// is cryptographic: for secrets, use crypto/rand through rand.New with a crypto-backed Source.
//
// Inputs are never modified except by ShuffleInPlace, whose name says so. Functions that build a slice
// append to dst and return the extended slice (doc/DECISIONS.md D-003); on error dst is returned with
// its length unchanged. Sequences are iter.Seq values, read once.
//
// Functions:
//
//	ALG-SAMPLE-001  ShuffleInPlace          uniform random permutation in place (Fisher–Yates)
//	ALG-SAMPLE-002  ReservoirSample         k uniform items from a sequence of unknown length (Algorithm R)
//	ALG-SAMPLE-003  NewWeightedPicker       prepares integer weights for weighted picks (prefix sums)
//	ALG-SAMPLE-004  WeightedPicker.Pick     an index chosen with probability weight/total, O(log n)
//	ALG-SAMPLE-005  RejectionSample         a draw from a proposal accepted by the caller's test
//	ALG-SAMPLE-006  RandomElement           one uniform item from a sequence of unknown length
//	ALG-SAMPLE-007  ShuffleDeterministic    a permutation of 0…n−1 fixed by a seed, for tests
package sample
