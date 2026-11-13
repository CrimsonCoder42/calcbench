// Package integer answers questions about integers: primality, prime factorisation, gcd and lcm, exact for
// every uint64.
//
// ID short form: INTEG (IDs are NUM-INTEG-NNN).
//
// Every function is deterministic: the Miller–Rabin bases are the 12 primes ≤ 37 (exact below 2⁶⁴ by
// Sorenson and Webster 2017) and Pollard's rho uses the fixed sequence c = 1, 2, 3, …, never random
// numbers. Products modulo n are formed in 128 bits (math/bits), so no step overflows.
//
// Functions:
//
//	NUM-INTEG-001  IsPrime              whether n is prime, exact for every uint64
//	NUM-INTEG-002  PrimeFactors         prime factors of n with multiplicity, ascending
//	NUM-INTEG-003  SmallestPrimeFactor  least prime factor of n
//	NUM-INTEG-004  GCD                  greatest common divisor (binary algorithm)
//	NUM-INTEG-005  LCM                  least common multiple, overflow-checked
//	NUM-INTEG-006  ExtendedGCD          gcd and Bézout coefficients (extended Euclid)
//	NUM-INTEG-007  FractionReduce       a fraction in lowest terms, positive denominator
//	NUM-INTEG-008  FractionAdd          sum of two fractions in lowest terms, overflow-checked
//	NUM-INTEG-009  FractionSub          difference of two fractions in lowest terms
//	NUM-INTEG-010  FractionMul          product of two fractions in lowest terms, overflow-checked
//	NUM-INTEG-011  FractionDiv          quotient of two fractions in lowest terms
//	NUM-INTEG-012  Divisors             every positive divisor of n, ascending
//	NUM-INTEG-013  CeilDiv              ⌈a/b⌉ for int64, exact, overflow-checked
package integer
