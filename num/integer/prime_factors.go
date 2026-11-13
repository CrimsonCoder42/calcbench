package integer

import (
	"errors"
	"fmt"
	"slices"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrOutOfRange: an input is outside the values the function accepts.
	ErrOutOfRange = errors.New("integer: out of range")
)

// smallPrimes are the 172 primes below 1024, used for trial division before Pollard's rho. The next prime
// is 1031, so a cofactor left after dividing them all out that is below 1031² is itself prime.
var smallPrimes = [172]uint64{
	2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53,
	59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131,
	137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223,
	227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311,
	313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409,
	419, 421, 431, 433, 439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, 503,
	509, 521, 523, 541, 547, 557, 563, 569, 571, 577, 587, 593, 599, 601, 607, 613,
	617, 619, 631, 641, 643, 647, 653, 659, 661, 673, 677, 683, 691, 701, 709, 719,
	727, 733, 739, 743, 751, 757, 761, 769, 773, 787, 797, 809, 811, 821, 823, 827,
	829, 839, 853, 857, 859, 863, 877, 881, 883, 887, 907, 911, 919, 929, 937, 941,
	947, 953, 967, 971, 977, 983, 991, 997, 1009, 1013, 1019, 1021,
}

// trialBound is 1031², the least composite with no prime factor in smallPrimes.
const trialBound = 1031 * 1031

// maxFactors is the most prime factors (with multiplicity) a uint64 can have: 2⁶³ has 63.
const maxFactors = 64

// PrimeFactors returns the prime factors of n with multiplicity, in ascending order: their product is n.
// Errors: ErrOutOfRange when n < 2.
//
// ID:      NUM-INTEG-002
// Time:    O(n^(1/4)) expected modular multiplications (heuristic, not proven)   n = the input
// Space:   O(log n) for the result (one allocation, at most log₂ n factors); O(1) extra
// Source:  Pollard, "A Monte Carlo method for factorization", BIT, 1975, 15(3):331–334, doi:10.1007/BF01933667
// Source:  Brent, "An improved Monte Carlo factorization algorithm", BIT, 1980, 20(2):176–184, doi:10.1007/BF01933190
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.5.4, Algorithms A and B
// Tags:    prime, factorization, pollard-rho, number-theory
// Units:   n dimensionless (a natural number); result dimensionless
// Example: PrimeFactors(2⁶⁴−1) = [3 5 17 257 641 65537 6700417], since 2⁶⁴−1 = (2³²−1)(2³²+1) and 2³²+1 = 641·6700417 (hand-checked)
func PrimeFactors(n uint64) ([]uint64, error) {
	if n < 2 {
		return nil, fmt.Errorf("%w: n = %d, want ≥ 2", ErrOutOfRange, n)
	}
	var buf [maxFactors]uint64
	k, m := trialDivide(n, &buf)
	if m > 1 {
		start := k
		k = factorCofactor(m, &buf, k)
		// Trial-division factors are ascending and all below the cofactor's; sort only the cofactor's.
		slices.Sort(buf[start:k])
	}
	return slices.Clone(buf[:k]), nil
}

// trialDivide divides out every prime of smallPrimes from n, storing each in buf with multiplicity, and
// returns the count and the cofactor m, which has no prime factor below 1031 (m = 1 when none is left).
func trialDivide(n uint64, buf *[maxFactors]uint64) (k int, m uint64) {
	m = n
	for _, p := range smallPrimes {
		if p*p > m {
			break
		}
		for m%p == 0 {
			buf[k] = p
			k++
			m /= p
		}
	}
	// Stopped early with p² > m: m is 1 or a prime, which the caller's trialBound check accepts.
	return k, m
}

// factorCofactor appends the prime factors of m > 1, which has no prime factor below 1031, to buf from
// index k (unordered) and returns the new count. An explicit stack replaces recursion: every entry is at
// least 1031 and their product is at most m < 2⁶⁴, so it never holds more than 6 (1031⁷ > 2⁶⁴).
func factorCofactor(m uint64, buf *[maxFactors]uint64, k int) int {
	var stack [8]uint64
	stack[0] = m
	top := 1
	for top > 0 {
		top--
		x := stack[top]
		if x < trialBound || IsPrime(x) {
			buf[k] = x
			k++
			continue
		}
		d := findFactor(x)
		stack[top] = d
		stack[top+1] = x / d
		top += 2
	}
	return k
}
