// Package bloom holds the Bloom filter, an m-bit set summary that answers "possibly present" or
// "definitely absent", and the formulas that size it.
//
// ID short form: BLOOM (IDs are ALG-BLOOM-NNN).
//
// Hashing: each item is hashed once with SHA-256 (FIPS 180-4); the first two 64-bit words h₁, h₂ of the
// digest give the k bit indices (h₁ + i·h₂ + (i³ − i)/6) mod m, i = 0 … k−1: enhanced double hashing
// (Dillinger and Manolios 2004, §5.2; Kirsch and Mitzenmacher 2008, §5.2, any m). The
// hash is fixed and keyless, so a filter's bits are the same on every machine and in every process:
// filters can be stored, sent and unioned. An adversary who chooses the items can search for ones that
// raise the false-positive rate.
//
// Guarantees: no false negatives (an added item always tests present); false positives at the rate
// FalsePositiveRate predicts. Items are read, never modified or retained.
//
// Concurrency: MayContain, EstimatedFPR and the accessors only read and may run concurrently with each
// other; Add and UnionWith write and need the caller's lock. The zero value holds no bits: Add does
// nothing and MayContain returns false; build filters with NewBloomFilter.
//
// Size: at most MaxBits = 1<<30 bits (128 MiB) per filter (ErrTooLarge); 1 ≤ hashCount ≤ min(bitCount,
// MaxHashCount = 1074) (ErrHashCount).
//
// Functions:
//
//	ALG-BLOOM-001  BloomFilter            NewBloomFilter, Add, MayContain, EstimatedFPR, BitCount, HashCount, SetBits
//	ALG-BLOOM-002  FalsePositiveRate      (1 − e^(−kn/m))^k for m bits, n items, k hashes
//	ALG-BLOOM-003  OptimalParams          m and k for n items at a target false-positive rate
//	ALG-BLOOM-004  BloomFilter.UnionWith  OR another filter's bits into this one
package bloom
