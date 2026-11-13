package bloom

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrBitCount: bitCount is less than 1.
	ErrBitCount = errors.New("bloom: bit count must be at least 1")
	// ErrHashCount: hashCount is less than 1.
	ErrHashCount = errors.New("bloom: hash count must be in [1, min(bit count, MaxHashCount)]")
	// ErrTooLarge: the filter would need more than MaxBits bits.
	ErrTooLarge = errors.New("bloom: filter larger than MaxBits")
	// ErrItems: expectedItems is less than 1.
	ErrItems = errors.New("bloom: expected items must be at least 1")
	// ErrProbability: a false-positive rate is not strictly between 0 and 1, or is NaN.
	ErrProbability = errors.New("bloom: false-positive rate must be in (0, 1)")
	// ErrIncompatible: two filters differ in bit count or hash count, or one is nil.
	ErrIncompatible = errors.New("bloom: filters are not compatible")
)

// MaxBits is the largest bitCount a filter may have: 1<<30 bits, 128 MiB, the largest allocation a
// library call makes silently in this run (lead decisions L13, L26).
const MaxBits = 1 << 30

// MaxHashCount is the largest hashCount a filter may have. At the optimum the false-positive rate is
// (1/2)^k (Broder–Mitzenmacher §2.1), so k > 1074 would aim below 2^-1074, the smallest positive float64;
// a larger k only costs time. OptimalParams never returns more.
const MaxHashCount = 1074

// BloomFilter is a Bloom filter of m bits and k hash functions: Add sets the k bits an item hashes to,
// and MayContain reports whether all k are set, so an added item always tests present and an absent one
// tests present with probability about (1 − e^(−kn/m))^k after n distinct items. The k indices come from
// one SHA-256 digest by enhanced double hashing, gᵢ = (h₁ + i·h₂ + (i³ − i)/6) mod m (Dillinger–Manolios
// §5.2): unlike plain double hashing, a step h₂ mod m that is 0 or shares a factor with m still gives
// distinct indices, so m need not be prime and h₂ is not restricted. Kirsch–Mitzenmacher Theorem 5.2
// assumes no prime m (h₁, h₂ uniform on [m], any m) and covers any added f(i), so the asymptotic
// false-positive rate is that of k independent hashes.
//
// Methods: Add(item) inserts item. MayContain(item) is false when item was certainly never added.
// EstimatedFPR() is (SetBits/m)^k, the false-positive probability given the bits now set (it ignores
// duplicate Adds and stays right after UnionWith). BitCount, HashCount and SetBits report m, k and the
// number of 1 bits. Add and UnionWith need the caller's lock; the rest only read.
//
// ID:      ALG-BLOOM-001
// Time:    Add and MayContain O(L + k)   L = len(item) in bytes (one SHA-256 pass), k = hash count; EstimatedFPR and accessors O(1); NewBloomFilter O(m/64)
// Space:   O(m/64) words, allocated once by NewBloomFilter; O(1) extra per call, no allocation
// Source:  Bloom, "Space/time trade-offs in hash coding with allowable errors", Communications of the ACM 13(7):422–426, 1970, doi:10.1145/362686.362692
// Source:  Kirsch and Mitzenmacher, "Less hashing, same performance: building a better Bloom filter", Random Structures & Algorithms 33(2):187–218, 2008, doi:10.1002/rsa.20208, §2 (standard Bloom filter; FPR (1 − ρ)^k given the fraction ρ of 0 bits) and §5.2 Theorem 5.2 ((enhanced) double hashing h₁ + i·h₂ + f(i) mod m, any m, any f)
// Source:  Dillinger and Manolios, "Bloom filters in probabilistic verification", FMCAD 2004, LNCS 3312, pp. 367–381, 2004, doi:10.1007/978-3-540-30494-4_26, §5.2 Algorithm 2 (enhanced double hashing)
// Source:  NIST FIPS 180-4, Secure Hash Standard, 2015, §6.2 (SHA-256)
// Tags:    bloom-filter, probabilistic, set-membership, approximate, double-hashing, sketch
type BloomFilter struct {
	words     []uint64
	bitCount  uint64 // m
	hashCount int    // k
	setBits   uint64 // number of 1 bits
}

// NewBloomFilter returns an empty BloomFilter of bitCount bits using hashCount hash functions (see
// OptimalParams to choose them).
// Errors: ErrBitCount when bitCount < 1; ErrHashCount when hashCount < 1, hashCount > bitCount or
// hashCount > MaxHashCount; ErrTooLarge when
// bitCount > MaxBits.
func NewBloomFilter(bitCount uint64, hashCount int) (*BloomFilter, error) {
	if bitCount < 1 {
		return nil, fmt.Errorf("%w: bitCount = %d", ErrBitCount, bitCount)
	}
	if hashCount < 1 {
		return nil, fmt.Errorf("%w: hashCount = %d", ErrHashCount, hashCount)
	}
	if bitCount > MaxBits {
		return nil, fmt.Errorf("%w: bitCount = %d > %d", ErrTooLarge, bitCount, MaxBits)
	}
	if uint64(hashCount) > bitCount || hashCount > MaxHashCount {
		return nil, fmt.Errorf("%w: hashCount = %d with bitCount = %d (max %d)", ErrHashCount, hashCount, bitCount, MaxHashCount)
	}
	return &BloomFilter{
		words:     make([]uint64, (bitCount+63)/64),
		bitCount:  bitCount,
		hashCount: hashCount,
	}, nil
}

// indexStart returns g₀ = h₁ mod m and the first step h₂ mod m for item.
func (f *BloomFilter) indexStart(item []byte) (first, step uint64) {
	digest := sha256.Sum256(item)
	h1 := binary.BigEndian.Uint64(digest[0:8])
	h2 := binary.BigEndian.Uint64(digest[8:16])
	return h1 % f.bitCount, h2 % f.bitCount
}

// nextProbe advances enhanced double hashing (Dillinger–Manolios §5.2, Algorithm 2) after probe i:
// g ← (g + step) mod m, then step ← (step + i + 1) mod m. g, step < m ≤ 2³⁰ and i + 1 ≤ k ≤ m, so each
// sum is below 2m: one conditional subtraction reduces it, and nothing overflows.
func nextProbe(g, step, i, m uint64) (uint64, uint64) {
	if g += step; g >= m {
		g -= m
	}
	if step += i + 1; step >= m {
		step -= m
	}
	return g, step
}

// Add inserts item: it sets the k bits item hashes to. The zero value ignores Add.
func (f *BloomFilter) Add(item []byte) {
	if f.bitCount == 0 {
		return
	}
	g, step := f.indexStart(item)
	for i := range f.hashCount {
		word, mask := g>>6, uint64(1)<<(g&63)
		if f.words[word]&mask == 0 {
			f.words[word] |= mask
			f.setBits++
		}
		g, step = nextProbe(g, step, uint64(i), f.bitCount)
	}
}

// MayContain reports whether item may have been added: false means it certainly was not; true means it
// was, or is a false positive. The zero value returns false.
func (f *BloomFilter) MayContain(item []byte) bool {
	if f.bitCount == 0 {
		return false
	}
	g, step := f.indexStart(item)
	for i := range f.hashCount {
		if f.words[g>>6]&(uint64(1)<<(g&63)) == 0 {
			return false
		}
		g, step = nextProbe(g, step, uint64(i), f.bitCount)
	}
	return true
}

// EstimatedFPR returns (SetBits/BitCount)^HashCount: the probability that an item never added tests
// present, given the bits set now (Kirsch–Mitzenmacher §2). 0 for an empty or zero-value filter.
func (f *BloomFilter) EstimatedFPR() float64 {
	if f.bitCount == 0 {
		return 0
	}
	return math.Pow(float64(f.setBits)/float64(f.bitCount), float64(f.hashCount))
}

// BitCount returns m, the number of bits.
func (f *BloomFilter) BitCount() uint64 { return f.bitCount }

// HashCount returns k, the number of bit indices per item.
func (f *BloomFilter) HashCount() int { return f.hashCount }

// SetBits returns the number of bits set to 1.
func (f *BloomFilter) SetBits() uint64 { return f.setBits }
