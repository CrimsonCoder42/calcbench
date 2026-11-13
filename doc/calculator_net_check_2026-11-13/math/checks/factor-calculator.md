# factor-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform cvar (default 120). Output: all factors, factor pairs, prime factorisation.
MATCH: PARTIAL NUM-INTEG-002 PrimeFactors (prime factorisation); all-divisors list missing. Factor pairs are (d, n/d) for the first half of the divisor list (no new maths).
RESEARCH: Rosen, Elementary Number Theory and Its Applications, 6th ed., Addison-Wesley 2011, §7.2 "The Sum and Number of Divisors" (divisors of ∏pᵢ^aᵢ are ∏pᵢ^bᵢ, 0 ≤ bᵢ ≤ aᵢ; τ(n) = ∏(aᵢ + 1)).
BUILD: NUM-INTEG-012 Divisors, num/integer/divisors.go (from the factorisation, one allocation sized τ(n), then sorted); tests num/integer/divisors_test.go (cases, 2⁶⁴ − 1 with 128 divisors, trial-division oracle, benchmark at 3 sizes).

COMPARE (scratch/factor-calculator.go; 5 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 120 | 16 divisors 1…120; primes 2·2·2·3·5 | same list; 2 × 2 × 2 × 3 × 5 | 0 | match |
| 97 | [1 97]; prime | (1, 97); "prime" | 0 | match |
| 1 | [1] | (1, 1) | 0 | match |
| 963761198400 (= 2⁶·3⁴·5²·7·11·13·17·19·23) | 6720 divisors; prime factors as shown | 6720 divisors, byte-identical list (cmp); same prime factors | 0 | match |
