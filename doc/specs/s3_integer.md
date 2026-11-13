# S3: Integer (primes and factorisation)

Status: SHIPPED 2026-11-12
Work folder: `doc/math_build_2026-11-12/`
Decisions: num/integer approved by Devin 2026-11-12 (Q1 of `doc/catalog/algebra.md`; D-id assigned at merge) · Assumptions: none
Criteria: `doc/catalog/algebra.md`, section "Routed: researched here for the proposed `num/integer`".

## Functions

| ID | Name | Signature | Units |
|---|---|---|---|
| `NUM-INTEG-001` | `IsPrime` | `func IsPrime(n uint64) bool` | n dimensionless; result bool |
| `NUM-INTEG-002` | `PrimeFactors` | `func PrimeFactors(n uint64) ([]uint64, error)` | n dimensionless; prime factors, ascending, with multiplicity |
| `NUM-INTEG-003` | `SmallestPrimeFactor` | `func SmallestPrimeFactor(n uint64) (uint64, error)` | n dimensionless; least prime factor |

Shared internals (unexported): `mulMod`, `addMod`, `powMod`, `gcd` (modular_arithmetic.go), `findFactor`,
`rhoBrent`, `trialFactor` (pollard_rho.go), `smallPrimes`, `factorCofactor` (prime_factors.go).

## Per function

### `NUM-INTEG-001` `IsPrime`
- **Source**: Sorenson, Webster 2017 (ψ₁₂ > 2⁶⁴, so the 12 prime bases ≤ 37 are exact on uint64); Miller 1976;
  Rabin 1980; Knuth TAOCP Vol. 2, 3rd ed., §4.5.4, Algorithm P. Full citations in the doc comment.
- **Best known**: O(log n) modular multiplications (12 bases, fixed); O(1) space. Reading the input is O(1).
- **Approach**: trial division by the 12 bases, then the strong-probable-prime test to each base, 128-bit
  products via `math/bits.Mul64` and `bits.Rem64`. **Lost**: `math/big.Int.ProbablyPrime(0)` (exact below 2⁶⁴
  but allocates a big.Int per call); Montgomery multiplication (faster on platforms without hardware 128/64
  division, but more code and a further source; can replace the internal `mulMod` later without API change);
  random bases (forbidden: no randomness).
- **Input modified**: no (scalar).
- **Worked example**: 3215031751 = 151·751·28351, a strong pseudoprime to bases 2, 3, 5, 7 (Pomerance,
  Selfridge, Wagstaff 1980) → false. 2⁶⁴ − 59 (largest prime below 2⁶⁴, hand-checked) → true.
- **Errors**: none; every uint64 has an answer (0 and 1 → false).

### `NUM-INTEG-002` `PrimeFactors`
- **Source**: Pollard 1975 (rho); Brent 1980 (cycle detection, batched gcds); Knuth §4.5.4 Algorithms A and B.
- **Best known**: for 64-bit inputs, Pollard rho with Brent: O(n^(1/4)) expected multiplications (heuristic,
  not proven). Result O(log n).
- **Approach**: trial division by the 172 primes < 1024; a cofactor m < 1031² is prime; otherwise IsPrime,
  and a composite is split by rho, f(x) = x² + c mod m, x₀ = 2, c = 1, 2, 3, … (deterministic), batches of
  128 products per gcd. An explicit stack (≤ 6 entries, since every factor is ≥ 1031) replaces recursion. If
  rho fails for c = 1..255 (never observed) the cofactor falls back to trial division, so the call always
  ends. **Lost**: trial division to √n (up to 2³² divisions for a 64-bit semiprime); SQUFOF (same order, more
  code, no advantage on uint64).
- **Input modified**: no.
- **Worked example**: 2⁶⁴ − 1 → [3 5 17 257 641 65537 6700417] (hand-checked: (2³²−1)(2³²+1), 641·6700417 = 2³²+1).
- **Errors**: `ErrOutOfRange` when n < 2.

### `NUM-INTEG-003` `SmallestPrimeFactor`
- **Source**: as NUM-INTEG-002.
- **Best known**: as NUM-INTEG-002; Pollard finds a factor, not the least, so a full factorisation is needed
  in the worst case. O(1) space, no allocation.
- **Approach**: trial division returns a small factor at once; a prime cofactor is returned; otherwise the
  shared factorisation into a stack buffer, then its minimum. **Lost**: trial division to √n (as above).
- **Input modified**: no.
- **Worked example**: 2⁶⁴ − 1 → 3; 4294967291 · 4294967279 = 18446743979220271189 → 4294967279 (hand-checked).
- **Errors**: `ErrOutOfRange` when n < 2.

## Out of scope

gcd, lcm, modular inverse and modular exponentiation as exported functions (future `num/integer` IDs);
integers above 2⁶⁴ (`math/big` territory); prime sieves and prime counting (future `num/integer` IDs).
