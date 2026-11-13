package integer

import (
	"fmt"
	"slices"
)

// Divisors returns every positive divisor of n in ascending order, 1 and n included. It factors n
// (PrimeFactors) and forms each product p₁^b₁ ⋯ p_k^b_k with 0 ≤ bᵢ ≤ aᵢ, so the cost is the
// factorisation plus the d(n) divisors themselves, not a trial division up to √n.
// Errors: ErrOutOfRange when n = 0 (every integer divides 0).
//
// ID:      NUM-INTEG-012
// Time:    O(n^(1/4) + d log d) expected   n = the input, d = d(n) the number of divisors (d(n) ≤ 103 680 for n < 2⁶⁴)
// Space:   O(d) for the result (one allocation sized ∏(aᵢ + 1)); O(log n) extra for the prime factors
// Source:  Rosen, Elementary Number Theory and Its Applications, 6th ed., Addison-Wesley, 2011, §7.2 "The Sum and Number of Divisors" (the divisors of n = ∏ pᵢ^aᵢ are exactly ∏ pᵢ^bᵢ, 0 ≤ bᵢ ≤ aᵢ; τ(n) = ∏(aᵢ + 1))
// Tags:    divisors, factors, factor-pairs, number-theory
// Units:   n dimensionless (a natural number); result dimensionless
// Example: Divisors(120) = [1 2 3 4 5 6 8 10 12 15 20 24 30 40 60 120] (τ(120) = 4·2·2 = 16, hand-checked)
func Divisors(n uint64) ([]uint64, error) {
	if n == 0 {
		return nil, fmt.Errorf("%w: n = 0, want ≥ 1", ErrOutOfRange)
	}
	if n == 1 {
		return []uint64{1}, nil
	}
	primes, err := PrimeFactors(n)
	if err != nil {
		return nil, err
	}
	count := 1
	for i := 0; i < len(primes); {
		j := i
		for j < len(primes) && primes[j] == primes[i] {
			j++
		}
		count *= j - i + 1
		i = j
	}
	divs := make([]uint64, 1, count)
	divs[0] = 1
	for i := 0; i < len(primes); {
		p := primes[i]
		j := i
		for j < len(primes) && primes[j] == p {
			j++
		}
		// Multiply every divisor so far by p, p², …, p^(j−i); every product divides n, so none overflows.
		base := len(divs)
		pk := uint64(1)
		for range j - i {
			pk *= p
			for _, d := range divs[:base] {
				divs = append(divs, d*pk)
			}
		}
		i = j
	}
	slices.Sort(divs)
	return divs, nil
}
