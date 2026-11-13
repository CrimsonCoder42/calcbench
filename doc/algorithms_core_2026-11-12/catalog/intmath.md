# Catalog: algo/intmath (SHORT = INTMTH), new package

Builder W3-B10, 2026-11-12. Integer algorithms on machine words and digit strings. Inputs never modified.
Every result that can exceed its type returns ErrOverflow. Signed results follow Go's truncated division.

## Not built here (one function, one home)
- GCD/LCM, ExtendedGCD, ModPow, ModInverse, SieveOfEratosthenes, IsPrimeDeterministic64 → num/integer (L4, D-005;
  IsPrime is NUM-INTEG-001).
- **SqrtNewton DROPPED**: math.Sqrt is correctly rounded (IEEE 754-2008 §5.4.1 requires squareRoot to be
  correctly rounded; Go's math.Sqrt documents IEEE behaviour and compiles to the hardware instruction). A
  Newton iteration on float64 can at best equal it and is slower. Integer ⌊√n⌋ is built (INTMTH-002).
- **NumberOfDigitOne HELD (no source)**: the per-position counting formula for "how many digit 1s appear in
  1…n" was found only on problem sites and blogs (search 2026-11-12); no textbook or paper located.

Sources used (confirmed 2026-11-12 by search / table of contents unless marked):
- Knuth, *TAOCP* Vol. 2, 3rd ed., 1997: §3.1 ex. 6 (Floyd's cycle finding, confirmed); §4.3.1 Algorithms A, D, M;
  §4.4 radix conversion Methods (1a), (1b), (2a); §4.6.3 Algorithm A (right-to-left binary exponentiation).
- Cohen, *A Course in Computational Algebraic Number Theory*, Springer GTM 138, 1993, Algorithm 1.7.1 (integer
  square root; the book is confirmed, the algorithm number is from the standard reference list, not seen).
- Crandall, Pomerance, *Prime Numbers: A Computational Perspective*, 2nd ed., 2005, Algorithm 9.2.11 (not seen; second source).
- Dijkstra, *A Discipline of Programming*, Prentice-Hall, 1976, ch. 17 "An exercise attributed to R. W. Hamming" (confirmed).
- Graham, Knuth, Patashnik, *Concrete Mathematics*, 2nd ed., 1994, §4.4 "Factorial Factors" (ε_p(n!); the
  equation number (4.24) is from memory of the text, section confirmed).
- Gauss, *Disquisitiones Arithmeticae*, 1801, art. 12 (digit sum ≡ number mod 9, confirmed) and art. 291–292
  (three-square theorem, confirmed).
- Hardy, Wright, *An Introduction to the Theory of Numbers*, 6th ed., OUP, 2008: Theorem 2 (fundamental theorem
  of arithmetic); §9.2 "Terminating and recurring decimals" (confirmed); §20.5 Theorem 369 (four squares) and
  §20.10 (g(2) = 4 via 8m + 7) (confirmed in the table of contents).
- Guy, *Unsolved Problems in Number Theory*, 3rd ed., Springer, 2004, E34 "Happy numbers".
- Cormen et al., *Introduction to Algorithms*, 3rd ed., 2009, §31.6 (repeated squaring); Warren, *Hacker's
  Delight*, 2nd ed., 2013, ch. 9 "Integer Division".

### ALG-INTMTH-001 PowerBySquaring
`func PowerBySquaring(base, exp int64) (int64, error)`
- **Source**: Knuth Vol. 2 §4.6.3 Algorithm A; CLRS §31.6. **Best known**: O(log e) multiplications
  (addition chains can save a constant, not the order; Knuth §4.6.3).
- **Approach**: right-to-left binary on |base| in uint64 with bits.Mul64; the limit is 2^63 − 1, or 2^63
  when the result is negative. The square of the running base is checked only when a set bit of e remains,
  because it is then a factor of the result. **Lost**: left-to-right (same count; needs the top bit first);
  float math.Pow (inexact above 2^53).
- **Worked example**: 2^23 = 8388608 (Knuth's x^23). (−2)^63 = MinInt64; 0^0 = 1.
- **Errors**: ErrNegative (exp < 0); ErrOverflow. **Tags**: integer, power, binary-exponentiation. LeetCode 50 (integer form).

### ALG-INTMTH-002 ISqrt
`func ISqrt(n uint64) uint64`
- **Source**: Cohen Algorithm 1.7.1; Crandall–Pomerance Algorithm 9.2.11.
- **Termination proof**: let s = ⌊√n⌋. (i) Never below s: since x is an integer, ⌊(x + ⌊n/x⌋)/2⌋ =
  ⌊(x + n/x)/2⌋ (nested floors with an integer term, Concrete Mathematics §3.2), and (x + n/x)/2 ≥ √n by AM–GM, so y ≥ s.
  (ii) Strict decrease above s: if x ≥ s + 1 > √n then n/x < x, so (x + n/x)/2 < x and y < x. The iterates
  are integers that fall strictly while > s and never drop below s, so they reach s; there (i) gives y ≥ s = x
  and the loop returns s.
  Start x₀ = 2^⌈b/2⌉ > √n (b = bit length). Overflow: x ≤ 2^32 and ⌊n/x⌋ ≤ 2^32 + 2.
- **Best known**: O(1) word operations for 64 bits; Newton from within a factor 2 converges quadratically
  (≤ 7 steps). **Lost**: float math.Sqrt then correct by ±1 (correct but depends on float rounding near
  2^64; the integer proof is simpler to audit); bit-by-bit digit method (32 steps).
- **Worked example**: 99 → 9, 100 → 10, 2^64 − 1 → 2^32 − 1. **Errors**: none. **Tags**: integer, isqrt. LeetCode 69, 367.

### ALG-INTMTH-003 LongDivision
`func LongDivision[T signedInteger](dividend, divisor T) (quotient, remainder T, err error)`
- **Source**: Knuth §4.3.1 Algorithm D in radix 2 (each quotient digit is one comparison); Warren ch. 9.
- **Best known**: O(w) bit steps without the hardware divide (w = bit length of |dividend|).
- **Approach**: restoring shift-and-subtract on magnitudes in uint64, then Go's signs (quotient toward zero,
  remainder with the dividend's sign). **Lost**: doubling the divisor and subtracting from the top (same
  O(w) but needs overflow care on the doubled divisor); repeated subtraction (O(quotient)).
- **Worked example**: Go spec: 5/3 = 1 r 2, −5/3 = −1 r −2, 5/−3 = −1 r 2, −5/−3 = 1 r −2.
- **Errors**: ErrDivideByZero; ErrOverflow for MinT / −1 (documented rule: refused, not wrapped).
- **Tags**: integer, division, shift-subtract. LeetCode 29.

### ALG-INTMTH-004 ReverseDigits
`func ReverseDigits[T signedInteger](x T) (T, error)`
- **Source**: Knuth §4.4 Methods (1a) and (1b). **Best known**: O(d), d ≤ 19 digits.
- **Approach**: reverse |x| in uint64 (cannot overflow: < 10^19 < 2^64), then range-check against T.
  **Lost**: reverse in T with a pre-multiply bound check per digit (needs T's max, which a generic body
  cannot name without reflection).
- **Worked example**: 1234 → 4321, −560 → −65. **Errors**: ErrOverflow. **Tags**: integer, digits, reverse. LeetCode 7.

### ALG-INTMTH-005 IsPalindromeNumber
`func IsPalindromeNumber[T integer](x T) bool`
- **Source**: Knuth §4.4 Methods (1a), (1b). **Best known**: O(d). **Approach**: reverse only the low half
  and compare (no overflow, d/2 divisions). **Lost**: format to a string (allocates); full reversal (can overflow).
- **Worked example**: 121, 1221, 0 → true; −121, 10 → false. **Errors**: none. **Tags**: integer, palindrome. LeetCode 9.

### ALG-INTMTH-006 IsHappyNumber
`func IsHappyNumber(n uint64) bool`
- **Source**: Knuth Vol. 2 §3.1 ex. 6 (Floyd); Guy E34 (definition, the unhappy cycle 4 → … → 4).
- **Best known**: O(log n) (read the digits once; after one step every value ≤ 1620). **Lost**: a visited set
  (O(cycle) memory, a map); hard-coding the cycle {4, …} (correct, but rests on a fact the code would not check).
- **Worked example** (Guy E34): happy numbers ≤ 50 are 1 7 10 13 19 23 28 31 32 44 49. **Errors**: none; 0 → false.
- **Tags**: integer, happy-number, cycle-detection. LeetCode 202.

### ALG-INTMTH-007 DigitalRoot
`func DigitalRoot(n uint64) uint64`
- **Source**: Gauss art. 12 (n ≡ digit sum mod 9). Since the root is a single digit 1…9 for n ≥ 1 and is
  ≡ n mod 9, it is 1 + (n − 1) mod 9. **Best known**: O(1). **Lost**: repeated digit sums (O(log n)).
- **Worked example**: 38 → 2. **Errors**: none; 0 → 0. **Tags**: integer, digital-root. LeetCode 258.

### ALG-INTMTH-008 FactorialTrailingZeros
`func FactorialTrailingZeros(n uint64) uint64`
- **Source**: Concrete Mathematics §4.4 (Legendre: ε_p(n!) = Σ ⌊n/p^k⌋); zeros = min(ε_2, ε_5) = ε_5.
- **Best known**: O(log n). **Lost**: computing n! (impossible beyond 20!).
- **Worked example**: 100! ends in 24 zeros. **Errors**: none. **Tags**: integer, factorial, legendre-formula. LeetCode 172.

### ALG-INTMTH-009 IsHammingNumber · ALG-INTMTH-010 NthHammingNumber
`func IsHammingNumber(n uint64) bool` · `func NthHammingNumber(n int) (uint64, error)`
- **Source**: Dijkstra 1976 ch. 17 (set definition; generation in order by three indexes into the sequence).
- **Best known**: membership O(log n). n-th: O(n) time and space by Dijkstra's merge; O(n^(2/3))-ish counting
  methods exist but need floating logarithms and lose exactness; n is at most 13282 for a uint64 result.
- **Lost** for 010: a min-heap of candidates (O(n log n)); enumerate triples and sort (O(n log n)).
- **Overflow**: 13282 Hamming numbers are ≤ 2^64 − 1 (counted exactly with math/big; the test recounts);
  the largest is 18_432_000_000_000_000_000 = 2^26·3^2·5^15. n > 13282 → ErrOverflow before allocating.
- **Worked example** (Dijkstra): 1 2 3 4 5 6 8 9 10 12 15 16 18 20 24 25 27 30.
- **Errors**: 010: ErrOutOfRange (n < 1), ErrOverflow. **Tags**: integer, hamming-number, 5-smooth, ugly-number. LeetCode 263, 264.

### ALG-INTMTH-011 MinSquaresToSum
`func MinSquaresToSum(n uint64) int`
- **Source**: Hardy–Wright §20.5 Theorem 369 (four squares); Gauss art. 291–292 (n is a sum of three
  squares iff n ≠ 4^a(8b+7)); H–W §20.10 (necessity). Two squares: 4m = x² + y² forces x, y even (squares are
  0 or 1 mod 4), so stripping factors 4 preserves the answer; m ≡ 3 (mod 4) is never a sum of two squares.
- **Time stated**: O(√n), worst case about √m two-pointer steps: ≈4.3e9 for n near 2⁶⁴ (seconds of CPU; the
  oracle test caps its sums below 2⁴¹).
- **Best known**: O(√n) with this method (a factorisation-based two-square test, Fermat/H–W Theorem 366,
  would be faster only with a fast factoriser; num/integer has Pollard rho but algo may not import num).
  **Lost**: the DP over all m ≤ n (O(n√n) time, O(n) space) and BFS (same).
- **Worked example**: 7, 15, 28 → 4; 12 → 3; 13 → 2; 25 → 1. **Errors**: none; 0 → 0.
- **Tags**: integer, sum-of-squares, lagrange, legendre. LeetCode 279.

### ALG-INTMTH-012 IsPowerOfThree
`func IsPowerOfThree[T integer](x T) bool`
- **Source**: Hardy–Wright Theorem 2 (unique factorisation ⇒ divisors of 3^40 are exactly 3^j). 3^40 =
  12157665459056928801 is the largest power of 3 ≤ 2^64 − 1 (checked with math/big). **Best known**: O(1).
  **Lost**: divide by 3 in a loop (O(log x)); float log (inexact).
- **Worked example**: 27 → true, 45 → false. **Errors**: none. **Tags**: integer, power-of-three. LeetCode 326.

### ALG-INTMTH-013 FractionToDecimal
`func FractionToDecimal(numerator, denominator int64, maxFractionDigits int) (string, error)`
- **Source**: Hardy–Wright §9.2 (a rational has a terminating or recurring decimal); Knuth §4.4 Method (2a).
- **Approach**: long division; remainders r ← 10r mod d in 128-bit (bits.Mul64/Div64, so |d| up to 2^63 is
  exact); a map from remainder to the position of its digit finds the first repeat; "(" is inserted there.
  **maxFractionDigits** is a caller cap because the period can be |d| − 1 digits (the output alone could be
  9·10^18 bytes); ErrTooLong past it. **Best known**: O(L) for L output digits. **Lost**: compute the
  pre-period max(v₂, v₅) and the period ord_q(10) first (needs gcd and order, i.e. num/integer, and still
  emits L digits); Floyd/Brent cycle detection (O(1) memory but ~3× the divisions, and the result is O(L) anyway).
- **Worked example** (H–W §9.2): 1/3 = 0.(3), 3/7 = 0.(428571); 1/6 = 0.1(6); −1/6 = "-0.1(6)"; 22/7 = 3.(142857).
- **Errors**: ErrDivideByZero; ErrNegative (cap < 0); ErrTooLong. **Tags**: fraction, repeating-decimal. LeetCode 166.

### ALG-INTMTH-014 AddBinaryStrings
`func AddBinaryStrings(a, b string) (string, error)`
- **Source**: Knuth §4.3.1 Algorithm A, radix 2. **Best known**: O(m + n). **Lost**: math/big parse and add
  (three allocations and conversions; this is one pass).
- **Worked example**: "1011" + "11" = "1110". Leading zeros accepted, none in the output.
- **Errors**: ErrEmpty; ErrInvalidDigit. **Tags**: binary, addition, digit-string. LeetCode 67.

### ALG-INTMTH-015 MultiplyDecimalStrings
`func MultiplyDecimalStrings(a, b string) (string, error)`
- **Source**: Knuth §4.3.1 Algorithm M, radix 10. **Best known** for schoolbook: O(m·n); asymptotically
  faster methods exist (Karatsuba O(n^1.585), Knuth §4.3.3; Schönhage–Strassen / Harvey–van der Hoeven
  O(n log n)) and are what math/big uses for large operands. Schoolbook kept: the request named it, it is
  exact and simple, and a caller with long operands should use math/big. **Lost**: Karatsuba (faster above
  a few hundred digits; not requested, and math/big already ships it).
- **Worked example**: "123" × "45" = "5535". **Errors**: ErrEmpty; ErrInvalidDigit. **Tags**: multiplication, digit-string. LeetCode 43.

### ALG-INTMTH-016 IncrementDigits
`func IncrementDigits(dst, digits []uint8) ([]uint8, error)`
- **Source**: Knuth §4.3.1 Algorithm A with v = 1, radix 10. **Best known**: O(n) worst case (all nines);
  O(1) amortised over a run of increments. **Approach**: find the last non-9, copy, bump, zero the tail.
- **Worked example**: [1 2 9] → [1 3 0]; [9 9] → [1 0 0]; [0 9] → [1 0] (leading zeros kept).
- **Errors**: ErrEmpty; ErrInvalidDigit (> 9). **Input modified**: no; appends to dst. **Tags**: digits, plus-one. LeetCode 66.
