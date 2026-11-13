# Knuth, The Art of Computer Programming, Vol. 2

**Citation**: Knuth, *The Art of Computer Programming*, Vol. 2: *Seminumerical Algorithms*, 3rd ed., Addison-Wesley, 1997.

| Taken | Locator | Used by | How checked |
|---|---|---|---|
| Strong pseudoprime test (Algorithm P) | §4.5.4, Algorithm P | NUM-INTEG-001 | `TestIsPrime_oracle` vs math/big ProbablyPrime(0) |
| Trial division and Pollard rho (Algorithms A and B) | §4.5.4, Algorithms A and B | NUM-INTEG-002, -003 | Factor oracles vs math/big |
| Binary gcd (internal `gcd`) | §4.5.2, Algorithm B | internal to num/integer | `TestModularArithmetic_oracle` vs math/big |
| Right-to-left binary powering (internal `powMod`) | §4.6.3, Algorithm A | internal to num/integer | `TestModularArithmetic_oracle` vs math/big |
| Horner's rule; n multiplications and n additions are optimal for a general polynomial (Theorems A, C there) | §4.6.4, pp. 486–488 | NUM-ALGEB-001 | Hand-checked [−1, 2, −6, 2] at 3 → 5; oracle vs math/big within Higham's γ₂ₙ |

Status: §4.5.4 and the helper sections CHECKED in the integer catalog pass, not re-checked in the build. §4.6.4 pp. 486–488 confirmed by two secondary sources (Wikipedia "Horner's method", citing Knuth 1997 §4.6.4 pp. 486–488; Erickson's notes, jeffe.cs.illinois.edu/teaching/497/08-polynomials.pdf). §4.6.4 eq. (2) NOT confirmed and not cited; the §4.6 introduction for the polynomial product was replaced by CLRS §30.1.


## Used by S10 (algorithms core, 2026-11-12)

IDs: ALG-DP-010, ALG-INTMTH-001, ALG-INTMTH-003, ALG-INTMTH-004, ALG-INTMTH-005, ALG-INTMTH-006, ALG-INTMTH-013, ALG-INTMTH-014, ALG-INTMTH-015, ALG-INTMTH-016, ALG-LIST-001, ALG-LIST-002, ALG-LIST-003, ALG-LIST-004, ALG-LIST-005, ALG-LIST-006, ALG-LIST-007, ALG-LIST-008, ALG-LIST-009, ALG-LIST-010, ALG-LIST-012, ALG-LIST-013, ALG-LIST-014, ALG-LIST-017, ALG-MATRIX-001, ALG-QUEUE-002, ALG-SAMPLE-001, ALG-SAMPLE-002, ALG-SAMPLE-007, ALG-SEQ-008, ALG-STACK-001, ALG-STACK-002, ALG-STACK-003, ALG-STACK-004, ALG-STACK-008, ALG-STACK-009, ALG-STACK-010, ALG-STREAM-001, ALG-STREAM-002, ALG-TREE-006, ALG-TREE-007, ALG-TREE-008, ALG-TREE-014, ALG-TREE-015, ALG-TREE-017.

Citations as written in the doc comments:
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1, Exercise 6 (Floyd's cycle test, the guard)
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1, Exercise 6 (Floyd's one-and-two-step walk; when the fast pointer ends after k steps the slow one is at index k = ⌊n/2⌋)
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1, Exercise 6 (Floyd's two-speed walk, R. W. Floyd's method as Knuth ascribes it)
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1, Exercise 6 and its answer (X_n = X_2n; the method Knuth ascribes to R. W. Floyd)
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1, Exercise 6 and its answer (μ by restarting one pointer at X_0, λ by one lap; Floyd's method)
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1, exercise 6 (Floyd's cycle-finding method: find n > 0 with X_n = X_2n using O(1) memory)
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.3.1 "The Classical Algorithms", Algorithm A (addition of nonnegative integers), radix b = 2
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.3.1 "The Classical Algorithms", Algorithm A (addition with carry), radix 10, with v = 1
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.3.1 "The Classical Algorithms", Algorithm D (long division), here in radix 2 where each quotient digit is found by one comparison
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.3.1 "The Classical Algorithms", Algorithm M (multiplication of nonnegative integers), radix b = 10
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.3.1, Algorithm A (addition of nonnegative integers in radix b, carry k ∈ {0, 1})
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.4 "Radix Conversion", Method (1a) (digits by repeated division by 10) and Method (1b) (Horner's rule to rebuild the reversed half)
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.4 "Radix Conversion", Method (1a) (the digits of an integer from least significant by repeated division by the radix) and Method (1b) (rebuilding an integer by Horner's rule u·10 + digit)
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.4 "Radix Conversion", Method (2a) (fraction digits by repeated multiplication by the radix)
- Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.6.3 "Evaluation of Powers", Algorithm A (right-to-left binary method for exponentiation)
- Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §3.4.2, Algorithm P (shuffling)
- Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §3.4.2, Algorithm P, via ShuffleInPlace (ALG-SAMPLE-001)
- Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §3.4.2, Algorithm R (reservoir sampling)

Checked: see the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ (confirmation status per citation; not first-hand items are listed in PROGRESS.md and REPORT.md).

## Matrix gaps (2026-11-12)

§4.6.3 "Evaluation of Powers", Algorithm A (right-to-left binary method), applied to matrix powers. Used by NUM-LINALG-122; test: Aᵏ equals k repeated products, and [1 1; 1 0]⁵ = [8 5; 5 3].
