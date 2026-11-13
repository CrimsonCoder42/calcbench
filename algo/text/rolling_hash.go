package text

import (
	"errors"
	"fmt"
	"math/bits"
)

// Errors returned by NewRollingHash and the functions built on it.
var (
	// ErrHashParams: the modulus is below 2, or the base is not in [1, modulus).
	ErrHashParams = errors.New("text: rolling hash needs modulus >= 2 and 1 <= base < modulus")
	// ErrWindowLen: the window length is less than 1.
	ErrWindowLen = errors.New("text: window length must be at least 1")
)

// Default fingerprint parameters for callers with no reason to choose: the Mersenne prime 2⁶¹−1 and a
// fixed base. A fixed base can be attacked with crafted input; every caller here verifies matches, so
// an attack costs time, never correctness. Pass a base drawn from your own random source to avoid that.
const (
	DefaultRollingModulus uint64 = 1<<61 - 1
	DefaultRollingBase    uint64 = 1_000_003
)

// RollingHash is the Karp–Rabin fingerprint of a window of bytes: the polynomial
// c₀·b^(w−1) + c₁·b^(w−2) + … + c_(w−1) mod q, kept up to date as the window slides one byte at a time.
// Build it with NewRollingHash. The zero value (modulus 0) is inert: Append and Roll do nothing and Sum
// stays 0, so a forgotten constructor never panics on a division by zero.
// Append extends the window (use it to fill the first window), Roll slides a full window by one byte,
// Sum reads the fingerprint, Len the number of bytes appended, Reset empties it.
// Equal windows have equal fingerprints; unequal windows can collide, so a caller that needs certainty
// compares the bytes on a match.
// Append and Roll cost O(1) each (one 64×64→128-bit product and remainder); nothing allocates.
// The type and its methods share ALG-TEXT-004, documented on NewRollingHash.
type RollingHash struct {
	base, modulus uint64
	top           uint64 // base^(windowLen-1) mod modulus: the weight of the byte leaving the window
	value         uint64
	n             int
}

// NewRollingHash returns an empty RollingHash for windows of windowLen bytes, with the given base and
// modulus. A prime modulus gives the collision bound; any modulus >= 2 is accepted.
// Errors: ErrHashParams when modulus < 2 or base is not in [1, modulus); ErrWindowLen when windowLen < 1.
//
// ID:      ALG-TEXT-004
// Time:    O(log w) here (square-and-multiply for base^(w−1)), then O(1) per Append and Roll   w = windowLen
// Space:   O(1), no allocation
// Source:  Karp, Rabin, "Efficient randomized pattern-matching algorithms", IBM Journal of Research and Development, 1987, 31(2):249–260, doi:10.1147/rd.312.0249
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §32.2 (Horner's rule, the rolling update, Figure 32.5)
// Tags:    string, hash, rolling-hash, fingerprint, rabin-karp, polynomial-hash
func NewRollingHash(base, modulus uint64, windowLen int) (RollingHash, error) {
	if modulus < 2 {
		return RollingHash{}, fmt.Errorf("modulus %d: %w", modulus, ErrHashParams)
	}
	if base < 1 || base >= modulus {
		return RollingHash{}, fmt.Errorf("base %d: %w", base, ErrHashParams)
	}
	if windowLen < 1 {
		return RollingHash{}, fmt.Errorf("windowLen %d: %w", windowLen, ErrWindowLen)
	}
	return RollingHash{base: base, modulus: modulus, top: powMod(base, uint64(windowLen-1), modulus)}, nil
}

// Append adds c at the end of the window: value = value·base + c (mod modulus). No-op on the zero value.
func (h *RollingHash) Append(c byte) {
	if h.modulus < 2 { // the zero value: inert, see the type comment
		return
	}
	h.value = addMod(mulMod(h.value, h.base, h.modulus), uint64(c)%h.modulus, h.modulus)
	h.n++
}

// Roll slides a full window by one byte: out is the first byte of the window, in the byte entering it.
// Call it only after windowLen Appends; on a window of another length the fingerprint is meaningless.
// No-op on the zero value.
func (h *RollingHash) Roll(out, in byte) {
	if h.modulus < 2 { // the zero value: inert, see the type comment
		return
	}
	v := subMod(h.value, mulMod(uint64(out)%h.modulus, h.top, h.modulus), h.modulus)
	h.value = addMod(mulMod(v, h.base, h.modulus), uint64(in)%h.modulus, h.modulus)
}

// Sum returns the fingerprint of the current window, in [0, modulus).
func (h *RollingHash) Sum() uint64 { return h.value }

// Len returns the number of bytes appended since construction or the last Reset (Roll does not change it).
func (h *RollingHash) Len() int { return h.n }

// Reset empties the window, keeping base, modulus and window length.
func (h *RollingHash) Reset() { h.value, h.n = 0, 0 }

// mulMod returns a·b mod m without overflow, through the 128-bit product.
func mulMod(a, b, m uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return bits.Rem64(hi, lo, m)
}

// addMod returns a+b mod m for a, b < m, correct even when a+b overflows 64 bits.
func addMod(a, b, m uint64) uint64 {
	s := a + b
	if s < a || s >= m {
		s -= m
	}
	return s
}

// subMod returns a−b mod m for a, b < m.
func subMod(a, b, m uint64) uint64 {
	if a >= b {
		return a - b
	}
	return a + (m - b)
}

// powMod returns base^exp mod m by square-and-multiply (Knuth, TAOCP Vol. 2, 3rd ed., §4.6.3, Algorithm A).
func powMod(base, exp, m uint64) uint64 {
	result := uint64(1) % m
	base %= m
	for exp > 0 {
		if exp&1 == 1 {
			result = mulMod(result, base, m)
		}
		base = mulMod(base, base, m)
		exp >>= 1
	}
	return result
}
