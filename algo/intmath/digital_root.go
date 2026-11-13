package intmath

// DigitalRoot returns the digital root of n: the single digit reached by summing the decimal digits of n
// repeatedly (38 → 11 → 2). Since 10 ≡ 1 (mod 9), n is congruent to its digit sum modulo 9, so the root
// is 0 for n = 0 and 1 + (n − 1) mod 9 otherwise; no digit is read.
//
// ID:      ALG-INTMTH-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Gauss, Disquisitiones Arithmeticae, 1801, art. 12 (a number and the sum of its digits have the same least residue modulo 9: casting out nines)
// Tags:    integer, digits, digital-root, digit-sum, casting-out-nines, modular
func DigitalRoot(n uint64) uint64 {
	if n == 0 {
		return 0
	}
	return 1 + (n-1)%9
}
