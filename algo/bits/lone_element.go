package bits

// LoneElementAmongPairs returns the XOR of all values: when every value but one occurs an even number of
// times (typically twice), that is the value occurring an odd number of times. Because v ⊕ v = 0 and ⊕
// is associative and commutative, each pair cancels wherever it stands. Empty input returns 0.
// The precondition is not checked: on other input the result is the XOR of everything (never a panic).
// values is read, never modified.
//
// ID:      ALG-BITS-010
// Time:    O(n)   n = len(values)
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques" (x ⊕ x = 0, x ⊕ 0 = x; ⊕ is associative and commutative)
// Tags:    bits, xor, single-number, lone-element, pairs
func LoneElementAmongPairs[T integer](values []T) T {
	var acc T
	for _, v := range values {
		acc ^= v
	}
	return acc
}

// LoneElementAmongTriples returns the value that occurs exactly once when every other value occurs
// exactly three times. For each bit position, the number of values with that bit set is 3m + b, where b
// is the lone value's bit; two words (ones, twos) hold that count modulo 3 for all positions at once,
// cycling 00 → 01 → 10 → 00 on each set bit, and ones ends equal to the lone value. Empty input returns
// 0. The precondition is not checked: on other input the result is the bitwise count modulo 3 read as
// "count = 1" (never a panic). values is read, never modified.
//
// ID:      ALG-BITS-011
// Time:    O(n)   n = len(values); four bitwise operations per value
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques" (bitwise operations act on all bit positions in parallel); the per-position counter modulo 3 and the proof are in doc/algorithms_core_2026-11-12/catalog/bits.md (ALG-BITS-011)
// Tags:    bits, single-number, lone-element, triples, modular-counting, bit-counting
func LoneElementAmongTriples[T integer](values []T) T {
	var ones, twos T
	for _, v := range values {
		ones = (ones ^ v) &^ twos
		twos = (twos ^ v) &^ ones
	}
	return ones
}
