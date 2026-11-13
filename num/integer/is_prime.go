package integer

import "math/bits"

// millerRabinBases are the 12 primes ≤ 37: as strong-probable-prime bases they decide primality exactly
// for every n < ψ₁₂ = 318665857834031151167461 > 2⁶⁴ (Sorenson and Webster 2017).
var millerRabinBases = [12]uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}

// IsPrime reports whether n is prime. The answer is exact for every uint64, not probabilistic: 0 and 1
// are not prime, 2 is.
//
// ID:      NUM-INTEG-001
// Time:    O(log n) modular multiplications   n = the input; 12 bases, a fixed number
// Space:   O(1), no allocation
// Source:  Sorenson, Webster, "Strong pseudoprimes to twelve prime bases", Mathematics of Computation, 2017, 86(304):985–1003, doi:10.1090/mcom/3134 (ψ₁₂ > 2⁶⁴: the 12 prime bases ≤ 37 are exact on uint64)
// Source:  Miller, "Riemann's hypothesis and tests for primality", Journal of Computer and System Sciences, 1976, 13(3):300–317, doi:10.1016/S0022-0000(76)80043-8
// Source:  Rabin, "Probabilistic algorithm for testing primality", Journal of Number Theory, 1980, 12(1):128–138, doi:10.1016/0022-314X(80)90084-0
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.5.4, Algorithm P
// Tags:    prime, primality, miller-rabin, number-theory
// Units:   n dimensionless (a natural number); result bool
// Example: IsPrime(3215031751) = false: 151·751·28351, a strong pseudoprime to bases 2, 3, 5, 7 (Pomerance, Selfridge, Wagstaff, "The pseudoprimes to 25·10⁹", Mathematics of Computation, 1980, 35(151):1003–1026, doi:10.1090/S0025-5718-1980-0572872-7)
func IsPrime(n uint64) bool {
	if n < 2 {
		return false
	}
	for _, p := range millerRabinBases {
		if n%p == 0 {
			return n == p
		}
	}
	// No prime factor ≤ 37, so a composite n is at least 41².
	if n < 41*41 {
		return true
	}
	// n − 1 = 2^s · d with d odd.
	s := bits.TrailingZeros64(n - 1)
	d := (n - 1) >> s
	for _, a := range millerRabinBases {
		if !strongProbablePrime(n, a, d, s) {
			return false
		}
	}
	return true
}

// strongProbablePrime reports whether odd n > a passes the strong test to base a, where n − 1 = 2^s·d,
// d odd: a^d ≡ 1, or a^(2^r·d) ≡ −1 (mod n) for some 0 ≤ r < s (Knuth §4.5.4, Algorithm P).
func strongProbablePrime(n, a, d uint64, s int) bool {
	x := powMod(a, d, n)
	if x == 1 || x == n-1 {
		return true
	}
	for r := 1; r < s; r++ {
		x = mulMod(x, x, n)
		if x == n-1 {
			return true
		}
	}
	return false
}
