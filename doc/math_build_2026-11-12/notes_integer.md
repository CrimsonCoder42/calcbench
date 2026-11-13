# notes_integer: num/integer builder hand-off (2026-11-12)

Package `num/integer` (SHORT `INTEG`), 3 functions, all built; none held back. Spec: `doc/specs/s3_integer.md`
(Status BUILDING). Checks run: `gofmt -l num/integer` (clean), `go build ./...` (clean), `go vet ./num/integer/`
(clean). Tests and benchmarks NOT run (lead verifies). Escape analysis (`go build -gcflags=-m`): only the
result clone in PrimeFactors and the error values reach the heap.

Files: `num/integer/` doc.go · is_prime.go · prime_factors.go (ErrOutOfRange, smallPrimes, trialDivide,
factorCofactor) · smallest_prime_factor.go · pollard_rho.go (unexported rho/Brent, trial fallback) ·
modular_arithmetic.go (unexported mulMod, addMod, powMod, gcd) · a `_test.go` beside each.

## (a) INDEX.md rows (final form)

| ID | Function | Package | Does (one line) | Time | Space | Source | Tags |
|---|---|---|---|---|---|---|---|
| NUM-INTEG-001 | IsPrime | num/integer | whether n is prime, exact for every uint64 | O(log n) modular multiplications   n = the input; 12 bases, a fixed number | O(1), no allocation | Sorenson, Webster, "Strong pseudoprimes to twelve prime bases", Mathematics of Computation, 2017, 86(304):985–1003, doi:10.1090/mcom/3134 (ψ₁₂ > 2⁶⁴: the 12 prime bases ≤ 37 are exact on uint64); Miller, "Riemann's hypothesis and tests for primality", Journal of Computer and System Sciences, 1976, 13(3):300–317, doi:10.1016/S0022-0000(76)80043-8; Rabin, "Probabilistic algorithm for testing primality", Journal of Number Theory, 1980, 12(1):128–138, doi:10.1016/0022-314X(80)90084-0; Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.5.4, Algorithm P | prime, primality, miller-rabin, number-theory |
| NUM-INTEG-002 | PrimeFactors | num/integer | prime factors of n with multiplicity, ascending | O(n^(1/4)) expected modular multiplications (heuristic, not proven)   n = the input | O(log n) for the result (one allocation, at most log₂ n factors); O(1) extra | Pollard, "A Monte Carlo method for factorization", BIT, 1975, 15(3):331–334, doi:10.1007/BF01933667; Brent, "An improved Monte Carlo factorization algorithm", BIT, 1980, 20(2):176–184, doi:10.1007/BF01933190; Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.5.4, Algorithms A and B | prime, factorization, pollard-rho, number-theory |
| NUM-INTEG-003 | SmallestPrimeFactor | num/integer | least prime factor of n | O(n^(1/4)) expected modular multiplications (heuristic, not proven)   n = the input | O(1), no allocation | Pollard, "A Monte Carlo method for factorization", BIT, 1975, 15(3):331–334, doi:10.1007/BF01933667; Brent, "An improved Monte Carlo factorization algorithm", BIT, 1980, 20(2):176–184, doi:10.1007/BF01933190; Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §4.5.4, Algorithms A and B | prime, factor, smallest-prime-factor, number-theory |

Note vs the draft rows in doc/algebra_2026-11-12/PENDING_SHARED.md: Time/Space reworded (variables defined,
003 written out instead of "as NUM-INTEG-002"), Pollard/Brent citations gained issue numbers (3) and (2),
IsPrime gained the Miller and Rabin lines (now confirmed). Use the rows above, not the draft.

## (b) doc/sources/ records

- **sorenson_webster_2017.md**: Sorenson, Webster, "Strong pseudoprimes to twelve prime bases", Mathematics of
  Computation, 2017 (online 2016-06-02), 86(304):985–1003, doi:10.1090/mcom/3134. Taken: ψ₁₂ =
  318665857834031151167461 > 2⁶⁴, so strong tests to the 12 primes ≤ 37 decide primality exactly on uint64.
  IDs: NUM-INTEG-001. Checked: ψ₁₂ re-verified by script to be a spsp to all 12 bases and > 2⁶⁴; tests
  TestIsPrime_strongPseudoprimes and TestIsPrime_oracle (vs math/big ProbablyPrime(0)).
- **miller_1976.md**: Miller, "Riemann's hypothesis and tests for primality", J. Computer and System Sciences,
  1976, 13(3):300–317, doi:10.1016/S0022-0000(76)80043-8. Taken: the strong (Miller) test. IDs: NUM-INTEG-001.
  Checked: test reproduced in TestIsPrime_pomeranceSelfridgeWagstaffExample (bases 2, 3, 5, 7 pass 3215031751).
- **rabin_1980.md**: Rabin, "Probabilistic algorithm for testing primality", J. Number Theory, 1980,
  12(1):128–138, doi:10.1016/0022-314X(80)90084-0. Taken: the test as a repeated-base algorithm (we use fixed,
  not random, bases). IDs: NUM-INTEG-001. Checked: as miller_1976.
- **knuth_taocp_vol2.md** (also listed by algebra): Knuth, TAOCP Vol. 2, 3rd ed., Addison-Wesley, 1997.
  Taken: §4.5.4 Algorithm P (strong test; NUM-INTEG-001), §4.5.4 Algorithms A and B (trial division, rho;
  NUM-INTEG-002/003); internal helpers §4.6.3 Algorithm A (binary powering, powMod) and §4.5.2 Algorithm B
  (binary gcd). Checked: TestModularArithmetic_oracle vs math/big; factor oracles.
- **pomerance_selfridge_wagstaff_1980.md**: Pomerance, Selfridge, Wagstaff, "The pseudoprimes to 25·10⁹",
  Mathematics of Computation, 1980, 35(151):1003–1026, doi:10.1090/S0025-5718-1980-0572872-7. Taken: worked
  example 3215031751 = 151·751·28351, spsp to bases 2, 3, 5, 7. IDs: NUM-INTEG-001. Checked:
  TestIsPrime_pomeranceSelfridgeWagstaffExample (factorisation, the four passes, IsPrime = false).
- **pollard_1975.md**: Pollard, "A Monte Carlo method for factorization", BIT, 1975, 15(3):331–334,
  doi:10.1007/BF01933667. Taken: rho iteration x ← x² + c mod m, gcd of differences. IDs: NUM-INTEG-002, 003.
  Checked: TestPrimeFactors_example (2⁶⁴−1, hand-checked), oracles vs math/big.
- **brent_1980.md**: Brent, "An improved Monte Carlo factorization algorithm", BIT, 1980, 20(2):176–184,
  doi:10.1007/BF01933190. Taken: power-of-two cycle detection, batching products before one gcd, backtrack on
  overshoot. IDs: NUM-INTEG-002, 003. Checked: TestFindFactor_nontrivial, TestPrimeFactors_oracle (100
  semiprimes of random 32-bit primes).

## (c) Citations

Confirmed 2026-11-12 via Crossref (api.crossref.org/works/<DOI>): Miller 1976 (title, author, JCSS 13(3)
300–317) and Rabin 1980 (title, author, JNT 12(1) 128–138), the two marked CONFIRM; also re-confirmed
Sorenson–Webster (Math. Comp. 86(304) 985–1003), Pollard (BIT 15(3) 331–334), Brent (BIT 20(2) 176–184).
Pomerance–Selfridge–Wagstaff: Crossref rate-limited (429); doi.org resolved the DOI to
ams.org/mcom/1980-35-151/S0025-5718-1980-0572872-7 (vol 35, no 151, matches). Knuth sections were CHECKED
in the catalog and not re-checked here. NOT confirmed: none. Held back: none.

## (d) FOUND, NOT FIXED

- doc/algebra_2026-11-12/PENDING_SHARED.md §1 draft integer rows differ from the final rows above (see note).
- PENDING_SHARED.md §3 taxonomy row and the D-id for Q1 still need applying by the lead (not owned here).
- Performance note, not a defect: mulMod uses bits.Rem64 (catalog choice); Montgomery multiplication would be
  faster where 128/64 division is slow (arm64). An internal swap, no API change; recorded in the spec as the
  alternative that lost for now.
