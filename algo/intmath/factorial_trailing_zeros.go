package intmath

// FactorialTrailingZeros returns the number of zeros at the end of n! written in decimal: the exponent
// of 5 in n!, since every factor 10 needs one 5 and the 2s are more plentiful. By Legendre's formula that
// exponent is ⌊n/5⌋ + ⌊n/25⌋ + ⌊n/125⌋ + …; n! itself is never formed. 0! = 1 has none.
//
// ID:      ALG-INTMTH-008
// Time:    O(log n)   ⌊log5 n⌋ + 1 divisions
// Space:   O(1), no allocation
// Source:  Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §4.4 "Factorial Factors", eq. (4.24): ε_p(n!) = Σ_{k≥1} ⌊n/p^k⌋
// Tags:    integer, factorial, trailing-zeros, legendre-formula, prime-exponent
func FactorialTrailingZeros(n uint64) uint64 {
	var zeros uint64
	for q := n / 5; q > 0; q /= 5 {
		zeros += q
	}
	return zeros
}
