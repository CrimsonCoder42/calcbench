package integer

import "fmt"

// SmallestPrimeFactor returns the least prime factor of n; it is n itself when n is prime.
// Errors: ErrOutOfRange when n < 2.
//
// ID:      NUM-INTEG-003
// Time:    O(n^(1/4)) expected modular multiplications (heuristic, not proven)   n = the input
// Space:   O(1), no allocation
// Source:  Pollard, "A Monte Carlo method for factorization", BIT, 1975, 15(3):331–334, doi:10.1007/BF01933667
// Source:  Brent, "An improved Monte Carlo factorization algorithm", BIT, 1980, 20(2):176–184, doi:10.1007/BF01933190
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.5.4, Algorithms A and B
// Tags:    prime, factor, smallest-prime-factor, number-theory
// Units:   n dimensionless (a natural number); result dimensionless
// Example: SmallestPrimeFactor(4294967291·4294967279 = 18446743979220271189) = 4294967279 (hand-checked)
func SmallestPrimeFactor(n uint64) (uint64, error) {
	if n < 2 {
		return 0, fmt.Errorf("%w: n = %d, want ≥ 2", ErrOutOfRange, n)
	}
	for _, p := range smallPrimes {
		if p*p > n {
			return n, nil
		}
		if n%p == 0 {
			return p, nil
		}
	}
	// No prime factor below 1031.
	if n < trialBound || IsPrime(n) {
		return n, nil
	}
	// Pollard's rho finds a factor, not the least one: factor completely, then take the minimum.
	var buf [maxFactors]uint64
	k := factorCofactor(n, &buf, 0)
	least := buf[0] // k ≥ 2 here: n is composite.
	for _, f := range buf[1:k] {
		least = min(least, f)
	}
	return least, nil
}
