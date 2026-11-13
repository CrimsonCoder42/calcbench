// Package intmath holds integer algorithms on machine words and digit strings: exponentiation by
// squaring, the exact integer square root, binary long division, digit manipulations (reversal,
// palindromes, digital root, increment), Hamming (5-smooth) numbers, the least number of squares with
// a given sum, repeating decimals, and schoolbook addition and multiplication of digit strings.
//
// ID short form: INTMTH (IDs are ALG-INTMTH-NNN).
//
// What is not here: primes, factorisation, gcd, lcm and modular arithmetic live in num/integer
// (doc/DECISIONS.md D-005; IsPrime is NUM-INTEG-001); factorials and binomials live in
// num/combinatorics. A floating-point square root is the standard library's math.Sqrt, which is
// correctly rounded; no Newton iteration on float64 is offered.
//
// Numbers: every result that can exceed its type returns ErrOverflow instead of wrapping. Signed
// results follow Go's truncated division (the quotient rounds toward zero; the remainder has the sign
// of the dividend). Digit strings are ASCII decimal or binary, most significant digit first, without a
// sign; results carry no leading zeros except the single digit "0". No input is modified.
//
// Functions:
//
//	ALG-INTMTH-001  PowerBySquaring         base^exp for int64, overflow checked
//	ALG-INTMTH-002  ISqrt                   ⌊√n⌋ for every uint64, exact
//	ALG-INTMTH-003  LongDivision            quotient and remainder by shift-and-subtract, no division operator
//	ALG-INTMTH-004  ReverseDigits           the decimal digits of x in reverse order, overflow checked
//	ALG-INTMTH-005  IsPalindromeNumber      whether x reads the same in decimal both ways
//	ALG-INTMTH-006  IsHappyNumber           whether iterating the sum of squared digits reaches 1
//	ALG-INTMTH-007  DigitalRoot             the repeated digit sum of n, in O(1)
//	ALG-INTMTH-008  FactorialTrailingZeros  the number of trailing zeros of n! in decimal
//	ALG-INTMTH-009  IsHammingNumber         whether n = 2^i·3^j·5^k (a regular, 5-smooth number)
//	ALG-INTMTH-010  NthHammingNumber        the n-th Hamming number, ascending from 1
//	ALG-INTMTH-011  MinSquaresToSum         the least number of perfect squares that sum to n
//	ALG-INTMTH-012  IsPowerOfThree          whether x is 3^k for some k ≥ 0
//	ALG-INTMTH-013  FractionToDecimal       the decimal expansion of a fraction, the repetend in parentheses
//	ALG-INTMTH-014  AddBinaryStrings        the sum of two binary digit strings
//	ALG-INTMTH-015  MultiplyDecimalStrings  the product of two decimal digit strings
//	ALG-INTMTH-016  IncrementDigits         a decimal digit slice plus one
package intmath
