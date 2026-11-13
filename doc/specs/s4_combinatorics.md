# S4: Combinatorics (`num/combinatorics`, SHORT `COMBI`)

Status: SHIPPED 2026-11-12
Work folder: `doc/math_build_2026-11-12/` (hand-off: `notes_combinatorics.md`)
Decisions: Devin, 2026-11-12 (math build request): uint64 results with `ErrOverflow` plus a `*big.Int` twin
(017–022); enumeration is a later `algo` package; Stirling numbers in if sourced; no Gamma wrapper.
Catalog: [`doc/catalog/combinatorics.md`](../catalog/combinatorics.md) STEP 3 holds the full entry per function;
this spec records only what the build fixes or changes.

## Functions

| ID | Name | Signature | Units |
|---|---|---|---|
| NUM-COMBI-001 | `Factorial` | `func Factorial(n int) (uint64, error)` | n: count; result: count |
| NUM-COMBI-002 | `LogFactorial` | `func LogFactorial(n int) (float64, error)` | n: count; result: natural log |
| NUM-COMBI-003 | `Binomial` | `func Binomial(n, k int) (uint64, error)` | counts; result: count |
| NUM-COMBI-004 | `LogBinomial` | `func LogBinomial(n, k int) (float64, error)` | counts; result: natural log |
| NUM-COMBI-005 | `Permutations` | `func Permutations(n, k int) (uint64, error)` | counts; result: count |
| NUM-COMBI-006 | `PermutationsWithRepetition` | `func PermutationsWithRepetition(n, k int) (uint64, error)` | counts; result: count |
| NUM-COMBI-007 | `CombinationsWithRepetition` | `func CombinationsWithRepetition(n, k int) (uint64, error)` | counts; result: count |
| NUM-COMBI-008 | `Multinomial` | `func Multinomial(counts []int) (uint64, error)` | counts; result: count |
| NUM-COMBI-009 | `CircularPermutations` | `func CircularPermutations(n, k int) (uint64, error)` | counts; result: count |
| NUM-COMBI-010 | `CircularPermutationsUndirected` | `func CircularPermutationsUndirected(n, k int) (uint64, error)` | counts; result: count |
| NUM-COMBI-011 | `Derangements` | `func Derangements(n int) (uint64, error)` | count; result: count |
| NUM-COMBI-012 | `Catalan` | `func Catalan(n int) (uint64, error)` | index; result: count |
| NUM-COMBI-013 | `StirlingFirstUnsigned` | `func StirlingFirstUnsigned(n, k int) (uint64, error)` | counts; result: count |
| NUM-COMBI-014 | `StirlingSecond` | `func StirlingSecond(n, k int) (uint64, error)` | counts; result: count |
| NUM-COMBI-015 | `Bell` | `func Bell(n int) (uint64, error)` | count; result: count |
| NUM-COMBI-016 | `StirlingApproximation` | `func StirlingApproximation(n int) (float64, error)` | count; result: approximation of n! |
| NUM-COMBI-017 | `MultinomialBig` | `func MultinomialBig(counts []int) (*big.Int, error)` | as 008 |
| NUM-COMBI-018 | `DerangementsBig` | `func DerangementsBig(n int) (*big.Int, error)` | as 011 |
| NUM-COMBI-019 | `CatalanBig` | `func CatalanBig(n int) (*big.Int, error)` | as 012 |
| NUM-COMBI-020 | `StirlingFirstUnsignedBig` | `func StirlingFirstUnsignedBig(n, k int) (*big.Int, error)` | as 013 |
| NUM-COMBI-021 | `StirlingSecondBig` | `func StirlingSecondBig(n, k int) (*big.Int, error)` | as 014 |
| NUM-COMBI-022 | `BellBig` | `func BellBig(n int) (*big.Int, error)` | as 015 |

Sentinels: `ErrNegative`, `ErrOverflow`, `ErrOutOfRange`, `ErrEmpty`, wrapped with the parameter name.
No function modifies its input. Every uint64/float64 function allocates nothing.

## Per function (changes from the catalog only; everything else as cataloged)

### NUM-COMBI-001 `Factorial`
- **Source**: DLMF eq. 5.4.1 (n! = Γ(n+1); §26.1 does not define n!, so the catalog's §26.1 is replaced); OEIS A000142.
- **Approach**: 21-entry table; **lost**: a multiply loop (O(n) for nothing).
- **Worked example**: 10! = 3 628 800 (OEIS A000142). **Errors**: n < 0, n > 20.

### NUM-COMBI-002 `LogFactorial`
- **Approach**: log of the exact table for n ≤ 20, else `math.Lgamma(n+1)`. **lost**: Σ ln i (O(n), error grows).

### NUM-COMBI-003 `Binomial`
- **Source**: DLMF 26.3.1, 26.3.2, 26.3.6, Table 26.3.1. **Worked example**: C(6,3)=20, C(8,2)=28, C(10,5)=252.
- **Approach**: multiplicative 26.3.6 with a gcd step; every intermediate is C(n−k+i, i) ≤ C(n,k), so no false
  overflow; C(2i,i) ≥ 2ⁱ bounds the loop at 64 steps. **lost**: Pascal triangle O(nk); float64.

### NUM-COMBI-004 `LogBinomial`
- **Approach**: m = min(k, n−k) ≤ 64: Σ log1p((n−m)/i); else Stirling series DLMF 5.11.1 with the large terms
  cancelled analytically (four correction terms, truncation < 1e-16 for x ≥ 65). **lost**: three-Lgamma
  difference (relative error 2.5e-8 at n = 1e9, k = 3). **Errors**: ErrNegative; ErrOutOfRange for k > n.

### NUM-COMBI-005 `Permutations`, 006, 007
- As cataloged. 005 stops within 64 multiplies (every factor but the last is ≥ 2). 006: square-and-multiply.
  007 computes C(n+k−1, k) in uint64, so n+k−1 never overflows int.

### NUM-COMBI-008 `Multinomial`
- **Source**: DLMF eq. 26.4.2 (which is itself the product-of-binomials form used). Table 26.4.1 not cited
  (its layout was not read reliably). Validation pass first, so a negative count is reported before overflow.

### NUM-COMBI-009 `CircularPermutations`, 010 `CircularPermutationsUndirected`
- **Source**: Brualdi §2.2 NOT confirmed (in the 5th ed. permutations are chapter 3; the section and
  theorem were not seen) and is dropped. Confirmed instead: OEIS A111492, a(n,k) = (k−1)!·C(n,k) = P(n,k)/k
  (k-cycles of S_n = circular arrangements of k of n); OEIS A144151, T(n,k) = C(n,k) for k ≤ 2, else
  C(n,k)·(k−1)!/2 (undirected cycles of length k on n labelled nodes).
- **Approach**: C(n,k)·(k−1)! (010: C(n,k)·((k−1)!/2) for k ≥ 3, so no false overflow at the boundary).
- **Worked example**: A111492 row 5: 5 10 20 30 24; A144151 row 5: 5 10 10 15 12.
- **Errors**: k = 0 or k > n is ErrOutOfRange (A144151 gives T(n,0) = 1; we keep the catalog's choice).

### NUM-COMBI-011 `Derangements`
- **Change**: n ≤ 20, not 21. d(21) = 18 795 307 255 050 944 540 > 2⁶⁴ − 1 (the catalog was wrong).

### NUM-COMBI-012 `Catalan`, 015 `Bell`
- Tables (37 and 26 entries) as cataloged; worked examples DLMF Table 26.5.1 C(10) = 16 796, Table 26.7.1
  B(5) = 52, B(10) = 115 975.

### NUM-COMBI-013 `StirlingFirstUnsigned`, 014 `StirlingSecond`
- **Approach**: fast paths k = n, k = n−1 (C(n,2)), k = 1; then since both are ≥ k^(n−k), k ≥ 2 and n−k ≥ 64
  is ErrOverflow at once; otherwise one band of the triangle indexed by u = i − j (≤ 64 entries, fixed
  array, no allocation). Every band entry is ≤ the result, so overflow is never false.
  Time O(n·(n−k+1)), and the loop stops by overflow when the result does not fit. **lost**: explicit
  alternating sums (26.8.6) and a full row of k+1 entries (needs allocation and computes entries larger than
  the result).

### NUM-COMBI-016 `StirlingApproximation`
- **Source**: Robbins 1955 (the bound 1/(12n+1) < rₙ < 1/(12n)); DLMF 5.6.1, 5.11.1. Equation numbers inside
  Robbins's paper are not cited (not seen in this build).

### NUM-COMBI-017…022 big twins
- 017: Π big.Int.Binomial (26.4.2). 018: d(n) = n·d(n−1) + (−1)ⁿ (from 26.13.4). 019: C(2n,n)/(n+1) (26.5.1).
  020/021: band of the triangle, indexed by j when k ≤ n−k, by u otherwise (O(min(k, n−k)+1) big ints).
  022: row of S(n,·) summed (26.7.2, 26.8.22). Inputs whose sums exceed MaxInt64 return ErrOverflow.

## Out of scope

Enumeration (later `algo` package), integer partitions p(n) (LATER), Gamma wrapper (`math.Gamma`/`Lgamma`),
FactorialBig/BinomialBig/PermutationsBig (one stdlib call each: `big.Int.MulRange`, `Binomial`).
