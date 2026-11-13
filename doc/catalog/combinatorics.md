# Catalog: combinatorics (`num/combinatorics`, SHORT `COMBI`)

Status: CATALOG, awaiting Devin's approval. Nothing is built until it is approved.
Work folder: [../combinatorics_2026-11-12/](../combinatorics_2026-11-12/README.md). Rows for shared files:
[PENDING_SHARED.md](../combinatorics_2026-11-12/PENDING_SHARED.md).

The input list is vCalc's Combinatorics names only (`.claude/rules/sources.md`); no vCalc page was read, so
where a name is ambiguous ("Comparison Permutations", "PC!") the reading is ours and marked. Read so as not to
duplicate: `doc/catalog/algebra.md` (proposes `num/integer` for primes), `doc/catalog/stats.md` (routes rows 12
Binomial Coefficient, 27 Circular r-Permutation, 66 Gamma Function here). `doc/catalog/prob.md` and a `num/prob`
work folder did not exist when this pass ran (2026-11-12), so there are no `num/prob` requests to confirm yet.

Conventions (each is a question in the report until Devin decides):
- Counts in are `int` (Go's idiom, and what `math/big.Int.Binomial` takes as `int64`); a negative count is
  `ErrNegative`. Exact counts out are `uint64`; log forms are `float64` (natural log).
- A result above 2⁶⁴ − 1 = 18 446 744 073 709 551 615 is `ErrOverflow`, never a wrapped value and never a panic.
  Each entry states the exact threshold, computed here with Python's exact integers (`doc/combinatorics_2026-11-12/`).
- k > n where the count is defined as 0 (binomial, permutations, Stirling) returns 0, no error, as DLMF 26.3.2
  does for C(m,n). Only a negative input or a malformed slice is an error.
- Sentinels: `ErrNegative`, `ErrOverflow`, `ErrEmpty` (multinomial with no parts).

## STEP 1: the input list, cleaned

| # | Name as given | Decision | Home | Reason |
|---|---|---|---|---|
| 1 | Combinatorics (collection) | DROP | — | A collection page. Bundles the counting rows below (factorial, nPr, nCr, repetition, circular, Stirling): NUM-COMBI-001…016 |
| 2 | Poker Calculator (five card) (calculator) | DROP | — | A calculator. Bundles the ten five-card hand probabilities; the formulas MOVE to `life/games` with the hand rows |
| 3 | Stats Calc (calculator) | DROP | — | A calculator. Bundles descriptive statistics; `num/stats` owns that subject |
| 4 | Circular Permutations | MERGE | `num/combinatorics` | Rows 4–6 are one function plus its turn-over variant: NUM-COMBI-009 CircularPermutations (k = n gives (n−1)!) |
| 5 | Circular r-Permutation of n Elements | MERGE | `num/combinatorics` | NUM-COMBI-009 CircularPermutations(n, k) |
| 6 | Circular r-Permutation of n Elements (second listing) | MERGE | `num/combinatorics` | Duplicate of row 5 (listed twice). NUM-COMBI-009; turn-over variant NUM-COMBI-010 |
| 7 | Combinations C(n,k) | KEEP | `num/combinatorics` | NUM-COMBI-003 Binomial (with row 10 merged in) |
| 8 | Combinations with Repetitions | KEEP | `num/combinatorics` | NUM-COMBI-007 CombinationsWithRepetition |
| 9 | Comparison Permutations | MERGE | `num/combinatorics` | OUR READING of the name (vCalc not read): the count of pairwise comparisons among n items, C(n,2) unordered or P(n,2) ordered, so NUM-COMBI-003 / -005 with k = 2. Not a game. Devin to confirm |
| 10 | Conditional Probability of B given A | MOVE | `num/prob` | A probability rule, P(B|A) = P(A∩B)/P(A) |
| 11 | Count of Combinations (nCr) | MERGE | `num/combinatorics` | Same function as row 7: NUM-COMBI-003 Binomial |
| 12 | factorial | KEEP | `num/combinatorics` | NUM-COMBI-001 Factorial; log form NUM-COMBI-002 LogFactorial |
| 13 | Gamma Function | DROP-STDLIB | — | `math.Gamma`, `math.Lgamma`. A wrapper adds only an error policy; the calcbench forms that add something are LogFactorial (-002) and LogBinomial (-004). Also routed here by `doc/catalog/stats.md` row 66; this is the answer to that route |
| 14 | Guess 1 to 31 | MOVE | `life/games` | A game (the binary card trick: five cards, a number 1–31). Not a counting formula |
| 15 | Law of Total Probability P(A+B) | MOVE | `num/prob` | A probability rule. (The name's "P(A+B)" reads as a union; `num/prob` owns both) |
| 16 | PC! | MERGE | `num/combinatorics` | Read as factorial (the request's default; the name gives nothing else): NUM-COMBI-001 |
| 17 | Permutations | KEEP | `num/combinatorics` | NUM-COMBI-005 Permutations (with row 18 merged in) |
| 18 | Permutations nPr | MERGE | `num/combinatorics` | Same function as row 17: NUM-COMBI-005 |
| 19 | Permutations with repetition | KEEP | `num/combinatorics` | NUM-COMBI-006 PermutationsWithRepetition (nᵏ); the multiset reading is NUM-COMBI-008 Multinomial |
| 20 | Powerball Chances | MOVE | `life/games` | A lottery's odds; the game's rules (69 white, 26 red) are data for `life/games`, which calls NUM-COMBI-003 |
| 21 | Prime Number Theorem | MOVE | `num/integer` | Number theory. `num/integer` is PROPOSED by the algebra pass (doc/catalog/algebra.md Q1), not yet approved; if it is refused, the home is Devin's call |
| 22 | Probability of A and B | MOVE | `num/prob` | A probability rule |
| 23 | Probability of A or B | MOVE | `num/prob` | A probability rule |
| 24 | Probability of Non-independent Events (A, B) | MOVE | `num/prob` | A probability rule (the multiplication rule) |
| 25 | Probability of Flush losing (fraction) | MOVE | `life/games` | A poker hand probability; counts come from NUM-COMBI-003 |
| 26 | Probability of Flush losing (percentage) | MERGE | `life/games` | Same value as the fraction row ×100; one function in `life/games` |
| 27 | Probability of Four of a Kind losing (fraction) | MOVE | `life/games` | A poker hand probability; counts come from NUM-COMBI-003 |
| 28 | Probability of Four of a Kind losing (percentage) | MERGE | `life/games` | Same value as the fraction row ×100; one function in `life/games` |
| 29 | Probability of Full House losing (fraction) | MOVE | `life/games` | A poker hand probability; counts come from NUM-COMBI-003 |
| 30 | Probability of Full House losing (percentage) | MERGE | `life/games` | Same value as the fraction row ×100; one function in `life/games` |
| 31 | Probability of High Card losing (fraction) | MOVE | `life/games` | A poker hand probability; counts come from NUM-COMBI-003 |
| 32 | Probability of High Card losing (percentage) | MERGE | `life/games` | Same value as the fraction row ×100; one function in `life/games` |
| 33 | Probability of One Pair losing (fraction) | MOVE | `life/games` | A poker hand probability; counts come from NUM-COMBI-003 |
| 34 | Probability of One Pair losing (percentage) | MERGE | `life/games` | Same value as the fraction row ×100; one function in `life/games` |
| 35 | Probability of Royal Flush losing (fraction) | MOVE | `life/games` | A poker hand probability; counts come from NUM-COMBI-003 |
| 36 | Probability of Royal Flush losing (percentage) | MERGE | `life/games` | Same value as the fraction row ×100; one function in `life/games` |
| 37 | Probability of Straight Flush losing (fraction) | MOVE | `life/games` | A poker hand probability; counts come from NUM-COMBI-003 |
| 38 | Probability of Straight Flush losing (percentage) | MERGE | `life/games` | Same value as the fraction row ×100; one function in `life/games` |
| 39 | Probability of Straight losing (fraction) | MOVE | `life/games` | A poker hand probability; counts come from NUM-COMBI-003 |
| 40 | Probability of Straight losing (percentage) | MERGE | `life/games` | Same value as the fraction row ×100; one function in `life/games` |
| 41 | Probability of Three of a kind losing (fraction) | MOVE | `life/games` | A poker hand probability; counts come from NUM-COMBI-003 |
| 42 | Probability of Three of a kind losing (percentage) | MERGE | `life/games` | Same value as the fraction row ×100; one function in `life/games` |
| 43 | Probability of Two Pair losing (fraction) | MOVE | `life/games` | A poker hand probability; counts come from NUM-COMBI-003 |
| 44 | Probability of Two Pair losing (percentage) | MERGE | `life/games` | Same value as the fraction row ×100; one function in `life/games` |
| 45 | Stirling Numbers | KEEP | `num/combinatorics` | NUM-COMBI-013 StirlingFirstUnsigned, -014 StirlingSecond (Q3: in or LATER) |
| 46 | Stirling's Formula | KEEP | `num/combinatorics` | NUM-COMBI-016 StirlingApproximation |
| 47 | Data item: five card hands | MOVE | `life/games` | Data (the hand counts); `life/games` derives them with NUM-COMBI-003 |
| 48 | Data item: Powerball Possibilities | MOVE | `life/games` | Data; C(69,5)·26 via NUM-COMBI-003 |
| 49 | Data item: Probability of Flush (fraction) | MOVE | `life/games` | Poker data |
| 50 | Data item: Probability of Flush (percentage) | MERGE | `life/games` | Same value as the fraction row ×100 |
| 51 | Data item: Probability of Four of a Kind (fraction) | MOVE | `life/games` | Poker data |
| 52 | Data item: Probability of Four of a Kind (percentage) | MERGE | `life/games` | Same value as the fraction row ×100 |
| 53 | Data item: Probability of Full House (fraction) | MOVE | `life/games` | Poker data |
| 54 | Data item: Probability of Full House (percentage) | MERGE | `life/games` | Same value as the fraction row ×100 |
| 55 | Data item: Probability of High Card (fraction) | MOVE | `life/games` | Poker data |
| 56 | Data item: Probability of High Card (percentage) | MERGE | `life/games` | Same value as the fraction row ×100 |
| 57 | Data item: Probability of One Pair (fraction) | MOVE | `life/games` | Poker data |
| 58 | Data item: Probability of One Pair (percentage) | MERGE | `life/games` | Same value as the fraction row ×100 |
| 59 | Data item: Probability of Royal Flush (fraction) | MOVE | `life/games` | Poker data |
| 60 | Data item: Probability of Royal Flush (percentage) | MERGE | `life/games` | Same value as the fraction row ×100 |
| 61 | Data item: Probability of Straight Flush (fraction) | MOVE | `life/games` | Poker data |
| 62 | Data item: Probability of Straight Flush (percentage) | MERGE | `life/games` | Same value as the fraction row ×100 |
| 63 | Data item: Probability of Straight (fraction) | MOVE | `life/games` | Poker data |
| 64 | Data item: Probability of Straight (percentage) | MERGE | `life/games` | Same value as the fraction row ×100 |
| 65 | Data item: Probability of Three of a kind (fraction) | MOVE | `life/games` | Poker data |
| 66 | Data item: Probability of Three of a kind (percentage) | MERGE | `life/games` | Same value as the fraction row ×100 |
| 67 | Data item: Probability of Two Pair (fraction) | MOVE | `life/games` | Poker data |
| 68 | Data item: Probability of Two Pair (percentage) | MERGE | `life/games` | Same value as the fraction row ×100 |
Counts (68 rows; the "for each of ten hands" lines are expanded, two rows per hand per list): KEEP 7 ·
MERGE 27 (7 into combinatorics functions, 20 percentage twins into their `life/games` fraction rows) ·
MOVE 30 (`num/prob` 5, `life/games` 24, `num/integer` 1) · DROP 3 · DROP-STDLIB 1 · LATER 0.
The request expected six probability-rule rows; the list has five (rows 10, 15, 22, 23, 24).

## STEP 2: the standard set

| ID | Function | What | Covers (STEP 1 rows) |
|---|---|---|---|
| NUM-COMBI-001 | Factorial | n!, exact, n ≤ 20 | 12, 16 |
| NUM-COMBI-002 | LogFactorial | ln n! | 12, 13 (the useful face of Gamma) |
| NUM-COMBI-003 | Binomial | C(n,k), exact, overflow detected | 7, 11, 9 (k = 2, our reading); stats row 12 |
| NUM-COMBI-004 | LogBinomial | ln C(n,k), for probability | 13 |
| NUM-COMBI-005 | Permutations | nPr = n!/(n−k)!, exact | 17, 18, 9 (ordered reading) |
| NUM-COMBI-006 | PermutationsWithRepetition | nᵏ (k-tuples from n) | 19 |
| NUM-COMBI-007 | CombinationsWithRepetition | C(n+k−1, k) (multisets) | 8 |
| NUM-COMBI-008 | Multinomial | (Σcᵢ)!/Πcᵢ! (also: permutations of a multiset) | 19 (multiset reading) |
| NUM-COMBI-009 | CircularPermutations | n!/((n−k)!·k) | 4, 5, 6; stats row 27 |
| NUM-COMBI-010 | CircularPermutationsUndirected | half of -009 for k ≥ 3 (necklace that may be turned over) | 4, 5, 6 |
| NUM-COMBI-011 | Derangements | !n | — |
| NUM-COMBI-012 | Catalan | Cₙ | — |
| NUM-COMBI-013 | StirlingFirstUnsigned | [n k], permutations of n with k cycles | 63 |
| NUM-COMBI-014 | StirlingSecond | {n k}, partitions of n into k blocks | 63 |
| NUM-COMBI-015 | Bell | Bₙ | — |
| NUM-COMBI-016 | StirlingApproximation | √(2πn)(n/e)ⁿ with Robbins's bound | 64 |
| NUM-COMBI-017…022 | MultinomialBig, DerangementsBig, CatalanBig, StirlingFirstUnsignedBig, StirlingSecondBig, BellBig | `*big.Int` twins, reserved ONLY if Q1 = "uint64 + big twin" | — |

Not functions here:
- **FactorialBig, BinomialBig, PermutationsBig: DROP-STDLIB.** `new(big.Int).MulRange(1, n)`,
  `new(big.Int).Binomial(n, k)`, `new(big.Int).MulRange(n−k+1, n)` (pkg.go.dev/math/big, read 2026-11-12). Each
  uint64 function's doc comment names the stdlib call for callers past the threshold.
- **Gamma function: DROP-STDLIB.** `math.Gamma`, `math.Lgamma` (with its sign result). A wrapper would add only
  NaN-in → error; LogFactorial and LogBinomial are the wrappers that add something (domain, accuracy).
- **Integer partitions p(n): LATER.** The source exists (DLMF 27.14.6, Euler's pentagonal recurrence; p(100) =
  190 569 292, DLMF §27.14), but no catalog in the tree uses it yet. It is number theory as much as counting:
  `num/integer` (if approved) is an equally good home. Devin's call when a user appears.
- **Enumeration** (generating the permutations, combinations, subsets themselves, lexicographic, as an iterator):
  recommended as a new `algo/enumerate` (SHORT `ENUM`), sources Knuth TAOCP vol. 4A §7.2.1.2 (all permutations,
  Algorithm L) and §7.2.1.3 (all combinations, Algorithm T): see PENDING_SHARED.md §3 and Q2. Not catalogued here.

## STEP 3: one entry per function

Source short forms (each read or confirmed 2026-11-12):
- **DLMF**: NIST Digital Library of Mathematical Functions, https://dlmf.nist.gov (read 2026-11-12). Chapter 26
  "Combinatorial Analysis" (Bressoud), chapter 5 "Gamma Function" (Askey, Roy), chapter 27 (Apostol).
- **Robbins1955**: H. Robbins, "A Remark on Stirling's Formula", American Mathematical Monthly 62(1), 1955,
  26–29, doi:10.2307/2308012 (JSTOR 2308012; read in scan).
- **OEIS**: The On-Line Encyclopedia of Integer Sequences, https://oeis.org (read 2026-11-12), second source for values.
- Knuth TAOCP vol. 1 and Graham–Knuth–Patashnik were NOT seen in this pass; they are cited only where the
  locator was confirmed (Concrete Mathematics §6.1 "Stirling Numbers", from the publisher's contents).

Shared: every function is pure, allocates nothing (except where stated) and never panics. A DLMF equation
read through a summarising fetcher is marked "(check by eye)" where the fetcher garbled it once.

### Factorial      (ID reserved: NUM-COMBI-001)
Formula: n! = 1·2·…·n, 0! = 1.
Inputs and units: `n int`, 0 ≤ n.
Output: `uint64`. 20! = 2 432 902 008 176 640 000 is the largest that fits; 21! = 51 090 942 171 709 440 000 > 2⁶⁴.
Errors: `ErrNegative` (n < 0); `ErrOverflow` (n > 20).
Primary source: DLMF §26.1 (n! and the falling factorial); OEIS A000142 (values).
Worked example: 10! = 3 628 800 (OEIS A000142); 20! as above, HAND-CHECKED with Python `math.factorial`.
Variants: big.Int → `big.Int.MulRange(1, n)` (DROP-STDLIB); log → NUM-COMBI-002.
Best approach: a 21-entry constant table: O(1) time, O(1) space, no multiply at run time.
Notes: must NOT compute in float64 and convert (exact only to 18!).

### LogFactorial      (ID reserved: NUM-COMBI-002)
Formula: ln n! = ln Γ(n+1).
Inputs and units: `n int`, 0 ≤ n.
Output: `float64`, natural log; finite for every int (ln((2⁶³−1)!) ≈ 3.9e20).
Errors: `ErrNegative`.
Primary source: DLMF §5.4 (Γ(n+1) = n!), §5.11 eq. 5.11.1 (Stirling series, used for the accuracy argument).
Worked example: ln 100! = 363.739375555563 (HAND-CHECKED: `math.Lgamma(101)` vs Python's exact 100! then log).
Variants: none.
Best approach: n ≤ 20: `math.Log` of the exact table entry (correctly rounded input); otherwise `math.Lgamma(n+1)`.
O(1) time and space.
Notes: Lgamma's sign result is always +1 here; the function discards it.

### Binomial      (ID reserved: NUM-COMBI-003)
Formula: C(n,k) = n!/(k!(n−k)!) (DLMF 26.3.1); C(n,k) = 0 for k > n (26.3.2); symmetry C(n,k) = C(n,n−k).
Inputs and units: `n, k int`, 0 ≤ n, 0 ≤ k.
Output: `uint64`. Every C(n,k) with n ≤ 67 fits (C(67,33) = 14 226 520 737 620 288 370); C(68,34) =
28 453 041 475 240 576 740 does not. Per k, the largest n that fits: k=2: 6 074 001 000; k=3: 4 801 280;
k=4: 145 056; k=5: 18 580; k=10: 386; k=20: 86; k=30: 68 (HAND-CHECKED, exact integers).
Errors: `ErrNegative`; `ErrOverflow` when the exact result exceeds 2⁶⁴ − 1.
Primary source: DLMF §26.3, eqs. 26.3.1, 26.3.2, 26.3.6 (multiplicative form), Table 26.3.1; OEIS A007318.
Worked example: DLMF Table 26.3.1: C(6,3) = 20, C(8,2) = 28, C(10,5) = 252. Overflow edge: C(67,33) fits,
C(68,34) → `ErrOverflow` (HAND-CHECKED).
Variants: big.Int → `big.Int.Binomial` (DROP-STDLIB); log → -004; table of a whole row → Pascal recurrence
26.3.5, not offered (O(n) per row is the caller's loop over this function or a later `BinomialRow`).
Best approach: k ← min(k, n−k); r = 1; for i = 1…k: r ← r·(n−k+i)/i, which is exact at every step because
r·(n−k+i)/i = C(n−k+i, i). To stay in 64 bits: g = gcd(r, i); r ← (r/g)·((n−k+i)/(i/g)), where i/g divides
n−k+i because i divides r·(n−k+i) and gcd(r/g, i/g) = 1. The multiply uses `math/bits.Mul64`; a non-zero
high word is `ErrOverflow`. The gcd step is our derivation from 26.3.6 (proved in the line above), so every
intermediate is ≤ the final result: no false overflow. O(min(k, n−k)) multiplies and gcds (gcd O(log n)), O(1) space.
Notes: must NOT use the Pascal triangle (O(nk)) or float64 (exact only below 2⁵³).

### LogBinomial      (ID reserved: NUM-COMBI-004)
Formula: ln C(n,k) = ln Γ(n+1) − ln Γ(k+1) − ln Γ(n−k+1).
Inputs and units: `n, k int`, 0 ≤ k ≤ n (k > n has no log: C = 0).
Output: `float64` natural log.
Errors: `ErrNegative`; `ErrOutOfRange` for k > n (log of 0; the caller tests k ≤ n or uses -003).
Primary source: DLMF 26.3.1 with DLMF §5.4 (Γ(n+1) = n!) and eq. 5.11.1 (Stirling series for ln Γ).
Worked example: ln C(100,50) = 66.783841652017 (HAND-CHECKED: Python exact C(100,50), then log). Accuracy
case: n = 10⁹, k = 3: the naive three-Lgamma difference gives 60.378036499 vs exact 60.378038039 (relative
error 2.5e-8 in the log, 1.5e-6 in C), so the naive form is NOT the approach.
Variants: exact → -003.
Best approach: m = min(k, n−k). m ≤ 64: Σᵢ₌₁..ₘ ln((n−m+i)/i) via `math.Log1p((n−m)/i)` terms, O(m). Larger m: the
Stirling-series form (DLMF 5.11.1): write ln Γ(x+1) = (x+½)ln x − x + ½ln 2π + δ(x) with δ the series remainder,
cancel the large terms analytically (k·ln(n/k) + (n−k+½)·log1p(k/(n−k)) − ½ln(2πk) form) and add
δ(n) − δ(k) − δ(n−k); O(1). The switch point 64 is proposed; verify measures the error on both sides.
Notes: the result is what `num/prob` binomial and hypergeometric probabilities need (offered to `num/prob`).

### Permutations      (ID reserved: NUM-COMBI-005)
Formula: P(n,k) = n!/(n−k)! = n(n−1)…(n−k+1), the falling factorial (DLMF §26.1); 0 for k > n.
Inputs and units: `n, k int`, 0 ≤ n, 0 ≤ k.
Output: `uint64`. P(20,20) = 20! fits; P(100,9) = 690 281 878 632 192 000 fits, P(100,10) does not (HAND-CHECKED).
Errors: `ErrNegative`; `ErrOverflow` (high word of `bits.Mul64` non-zero).
Primary source: DLMF §26.1 (falling factorial); DLMF §26.17, Table 26.17.1 (one-to-one maps of labelled
balls into labelled boxes = k-permutations).
Worked example: P(10,3) = 720; P(5,5) = 120 (HAND-CHECKED; equals -001 at k = n).
Variants: big.Int → `big.Int.MulRange(n−k+1, n)` (DROP-STDLIB).
Best approach: k multiplies with overflow check; stop early on overflow. O(k) time, O(1) space.
Notes: none.

### PermutationsWithRepetition      (ID reserved: NUM-COMBI-006)
Formula: nᵏ, the number of length-k sequences over n symbols.
Inputs and units: `n, k int`, ≥ 0 (0⁰ = 1, the empty sequence).
Output: `uint64`. Overflow when nᵏ > 2⁶⁴ − 1 (e.g. 2⁶³ fits, 2⁶⁴ does not; 3⁴⁰ fits, 3⁴¹ does not).
Errors: `ErrNegative`; `ErrOverflow`.
Primary source: DLMF §26.17, Table 26.17.1 (unrestricted maps of labelled balls into labelled boxes).
Worked example: n 26, k 3 → 17 576 (three-letter codes). HAND-CHECKED.
Variants: big.Int → `big.Int.Exp` (stdlib). Multiset reading → -008.
Best approach: binary exponentiation with `bits.Mul64` overflow checks. O(log k) time, O(1) space.
Notes: none.

### CombinationsWithRepetition      (ID reserved: NUM-COMBI-007)
Formula: C(n+k−1, k), the number of size-k multisets from n kinds.
Inputs and units: `n, k int`, ≥ 0. n = 0, k = 0 → 1; n = 0, k > 0 → 0.
Output: `uint64`, threshold that of -003 at (n+k−1, k).
Errors: `ErrNegative`; `ErrOverflow` (including n+k−1 overflowing int).
Primary source: DLMF §26.17, Table 26.17.1 (unlabelled balls into labelled boxes, unrestricted).
Worked example: n 5, k 3 → C(7,3) = 35. HAND-CHECKED.
Variants: as -003.
Best approach: calls -003. O(min(k, n−1)).
Notes: none.

### Multinomial      (ID reserved: NUM-COMBI-008)
Formula: (c₁ + … + cₘ)! / (c₁!·…·cₘ!) (DLMF 26.4.2); also the number of distinct permutations of a multiset.
Inputs and units: `counts []int`, each ≥ 0, len ≥ 1.
Output: `uint64`; overflow detected as in -003.
Errors: `ErrEmpty`; `ErrNegative`; `ErrOverflow` (also if Σcᵢ overflows int).
Primary source: DLMF §26.4, eq. 26.4.2, Table 26.4.1.
Worked example: (5,3,2) → 2520; MISSISSIPPI (1,4,4,2) → 34 650. HAND-CHECKED.
Variants: big.Int twin -017 if Q1 = twin.
Best approach: product of binomials Πⱼ C(sⱼ, cⱼ) with sⱼ = c₁+…+cⱼ, each by -003's method with overflow check.
O(m + Σ min(cⱼ, sⱼ₋₁)) time, O(1) space. The input slice is not modified.
Notes: must NOT compute (Σc)! (overflows at 21 even when the result is small).

### CircularPermutations      (ID reserved: NUM-COMBI-009)
Formula: arrangements of k of n distinct objects around a circle, rotations equal: P(n,k)/k = n!/((n−k)!·k);
k = n gives (n−1)!.
Inputs and units: `n, k int`, 1 ≤ k ≤ n (k = 0 is `ErrOutOfRange`: an empty circle is not defined).
Output: `uint64`, equal to C(n,k)·(k−1)!; overflow detected on that product.
Errors: `ErrNegative`; `ErrOutOfRange` (k = 0 or k > n); `ErrOverflow`.
Primary source: SOURCE TO CONFIRM. Brualdi, Introductory Combinatorics, 5th ed., Pearson, 2010, §2.2
(circular permutations, n!/(r·(n−r)!)) per a search result; the theorem number and text were NOT seen. No
DLMF or OEIS equation found for the r-subset form. Must be confirmed in a copy before build.
Worked example: n 6, k 4 → 90; n 6, k 6 → 120 = 5!. HAND-CHECKED (P(6,4) = 360, ÷ 4).
Variants: turn-over (necklace / key ring) → -010.
Best approach: P(n,k)/k = C(n,k)·(k−1)!: -003 then k−2 multiplies with `bits.Mul64` checks. Every
intermediate ≤ the result, so no false overflow and no division. O(k) time, O(1) space.
Notes: covers vCalc rows 4–6 and `doc/catalog/stats.md` row 27.

### CircularPermutationsUndirected      (ID reserved: NUM-COMBI-010)
Formula: as -009, with a circle and its mirror image equal (it can be turned over): P(n,k)/(2k) for k ≥ 3;
equal to -009 for k ≤ 2.
Inputs and units, Output, Errors: as -009.
Primary source: SOURCE TO CONFIRM, as -009 (Brualdi 5th ed. §2.2, not seen).
Worked example: n 6, k 6 → 60 (necklace of six distinct beads). HAND-CHECKED.
Variants: none.
Best approach: -009 then halve. Exact: for k ≥ 3, -009 = C(n,k)·(k−1)! and (k−1)! is even. O(k), O(1).
Notes: if the source cannot be confirmed, -009 and -010 become LATER (no source, no ship).

### Derangements      (ID reserved: NUM-COMBI-011)
Formula: d(n) = n!·Σⱼ₌₀..ₙ (−1)ʲ/j! (DLMF 26.13.4, check by eye); recurrence d(n) = (n−1)(d(n−1) + d(n−2)), d(0) = 1, d(1) = 0.
Inputs and units: `n int`, ≥ 0.
Output: `uint64`. d(20) = 895 014 631 192 902 121 fits; d(21) = 18 795 307 255 050 944 540 also fits (< 2⁶⁴);
d(22) = 413 496 759 611 120 779 881 does not. So n ≤ 21.
Errors: `ErrNegative`; `ErrOverflow` (n > 21).
Primary source: DLMF §26.13, eq. 26.13.4; OEIS A000166 (values). The recurrence is checked against OEIS values.
Worked example: d(4) = 9, d(6) = 265 (OEIS A000166); d(10) = 1 334 961 (HAND-CHECKED).
Variants: big.Int twin -018 if Q1 = twin.
Best approach: 22-entry constant table. O(1).
Notes: none.

### Catalan      (ID reserved: NUM-COMBI-012)
Formula: Cₙ = C(2n, n)/(n+1) (DLMF 26.5.1); recurrence C(n+1) = 2(2n+1)/(n+2)·C(n) (26.5.4).
Inputs and units: `n int`, ≥ 0.
Output: `uint64`. C₃₆ = 11 959 798 385 860 453 492 fits; C₃₇ = 45 950 804 324 621 742 364 does not. So n ≤ 36.
Errors: `ErrNegative`; `ErrOverflow` (n > 36).
Primary source: DLMF §26.5, eqs. 26.5.1, 26.5.4, Table 26.5.1; OEIS A000108.
Worked example: C₁₀ = 16 796 (DLMF Table 26.5.1).
Variants: big.Int twin -019.
Best approach: 37-entry constant table. O(1).
Notes: none.

### StirlingFirstUnsigned      (ID reserved: NUM-COMBI-013)
Formula: [n k] = |s(n,k)|, the number of permutations of n elements with exactly k cycles; DLMF's s(n,k) is
signed, s(n,k) = (−1)ⁿ⁻ᵏ[n k] (26.8.3). Recurrence (unsigned form of 26.8.18): [n k] = [n−1 k−1] + (n−1)[n−1 k].
Inputs and units: `n, k int`, ≥ 0; k > n → 0; [0 0] = 1.
Output: `uint64`. Every [n k] with n ≤ 21 fits; for n ≥ 22 depends on k (checked per call).
Errors: `ErrNegative`; `ErrOverflow`.
Primary source: DLMF §26.8, eqs. 26.8.3, 26.8.18, Table 26.8.1; Graham, Knuth, Patashnik, Concrete
Mathematics, 2nd ed., Addison-Wesley, 1994, §6.1 "Stirling Numbers"; OEIS A132393 (unsigned), A008275 (signed).
Worked example: DLMF Table 26.8.1 row n = 6: s(6,3) = −225 → [6 3] = 225; s(10,5) = −269 325 → 269 325.
Variants: signed → the caller applies (−1)ⁿ⁻ᵏ (int64 would halve the range). big.Int twin -020.
Best approach: one row of the triangle, updated in place right-to-left, k+1 entries: O(n·k) time, O(k) space
(one allocation of k+1 uint64, or none for k ≤ 32 using a fixed array). The explicit sums alternate in sign and
are not used.
Notes: must NOT return the signed value in a uint64.

### StirlingSecond      (ID reserved: NUM-COMBI-014)
Formula: {n k}, partitions of an n-set into k non-empty blocks. Recurrence {n k} = k{n−1 k} + {n−1 k−1}
(DLMF 26.8.22). Explicit: {n k} = (1/k!)Σⱼ (−1)ᵏ⁻ʲ C(k,j) jⁿ (26.8.6), not used (alternating).
Inputs and units: `n, k int`, ≥ 0; k > n → 0; {0 0} = 1.
Output: `uint64`. Every {n k} with n ≤ 26 fits; larger n checked per call.
Errors: `ErrNegative`; `ErrOverflow`.
Primary source: DLMF §26.8, eqs. 26.8.6, 26.8.22, Table 26.8.2; Concrete Mathematics 2nd ed. §6.1; OEIS A008277.
Worked example: DLMF Table 26.8.2 row 6: {6 3} = 90; {10 5} = 42 525.
Variants: big.Int twin -021.
Best approach: as -013. O(n·k) time, O(k) space.
Notes: none.

### Bell      (ID reserved: NUM-COMBI-015)
Formula: Bₙ = Σₖ {n k} (DLMF 26.7.2); B(n+1) = Σₖ C(n,k)B(k) (26.7.6).
Inputs and units: `n int`, ≥ 0.
Output: `uint64`. B₂₅ = 4 638 590 332 229 999 353 fits; B₂₆ = 49 631 246 523 618 756 274 does not. So n ≤ 25.
Errors: `ErrNegative`; `ErrOverflow` (n > 25).
Primary source: DLMF §26.7, eqs. 26.7.2, 26.7.6, Table 26.7.1; OEIS A000110.
Worked example: B₅ = 52, B₁₀ = 115 975 (DLMF Table 26.7.1).
Variants: big.Int twin -022.
Best approach: 26-entry constant table, generated and checked against DLMF/OEIS in the test. O(1).
Notes: none.

### StirlingApproximation      (ID reserved: NUM-COMBI-016)
Formula: n! ≈ √(2πn)·(n/e)ⁿ. Bound (Robbins1955 eqs. 1–2): n! = √(2π)·n^(n+½)·e^(−n)·e^(rₙ) with
1/(12n+1) < rₙ < 1/(12n), n ≥ 1. So the approximation is low by a factor between e^(1/(12n+1)) and e^(1/(12n)).
Inputs and units: `n int`, 1 ≤ n ≤ 170.
Output: `float64` approximation of n!. Documented relative error: n!/approx ∈ (e^(1/(12n+1)), e^(1/(12n))):
n = 1: 1.0844 (bounds 1.0800, 1.0869); n = 10: 1.008365 (1.008299, 1.008368); n = 100: 1.000834.
Errors: `ErrOutOfRange` (n < 1); `ErrOverflow` (n > 170: 171! > max float64).
Primary source: Robbins1955 eqs. (1), (2); DLMF §5.11 eq. 5.11.1 (the series) and §5.6 eq. 5.6.1
(1 < (2π)^(−½) x^(½−x) eˣ Γ(x) < e^(1/(12x)), the same bound for real x > 0).
Worked example: n = 10: 3 598 695.619 vs 10! = 3 628 800; ratio 1.008365, inside Robbins's (1.008299, 1.008368).
HAND-CHECKED.
Variants: for an accurate n!, use -001 (exact) or -002 (log). This function exists to teach and to bound.
Best approach: exp(½ln(2πn) + n·ln n − n) in log space (no overflow of nⁿ before 170). O(1).
Notes: must NOT be used inside -002 or -004 (they use Lgamma and the full series).

### Big twins      (IDs reserved only if Q1 = "uint64 + big.Int twin": NUM-COMBI-017…022)
MultinomialBig (017), DerangementsBig (018), CatalanBig (019), StirlingFirstUnsignedBig (020), StirlingSecondBig
(021), BellBig (022). Each: same formula, source and worked example as its uint64 function; signature
`(n …int) (*big.Int, error)`; no overflow error; time as the uint64 form times the cost of a big multiply on
O(n log n)-bit numbers; space O(result bits) (StirlingFirst/Second: O(k) big ints). Derangements, Catalan and Bell
use the recurrences 26.13 / 26.5.4 / Bell triangle instead of the tables. The other uint64 functions need no
twin: their big forms are one stdlib call (`big.Int.MulRange`, `Binomial`, `Exp`).

## STEP 4: order and dependencies

Build order (IDs are already in it): 001 Factorial → 002 LogFactorial → 003 Binomial → 004 LogBinomial →
005 Permutations → 006 PermutationsWithRepetition → 007 CombinationsWithRepetition (calls 003) → 008 Multinomial
(003's step) → 009, 010 Circular (005) → 011 Derangements, 012 Catalan → 013, 014 Stirling → 015 Bell
(014's identity checks it) → 016 StirlingApproximation → 017–022 if approved.

Dependencies: none on another subject. gcd for 003 is an unexported 5-line binary gcd here, not
`num/integer` (proposed, unapproved; and `num/integer` importing `num/combinatorics` or the reverse is a
same-group cross-import to avoid). Standard library only: `math`, `math/bits`, `math/big` (twins).

`num/prob` requests: none exist yet (no `doc/catalog/prob.md`, no prob work folder at 2026-11-12). Offers to
`num/prob`, in PENDING_SHARED.md: 003 Binomial and 004 LogBinomial (binomial and hypergeometric pmf), 001/002,
005. Offers to `life/games`: 003 for poker hand counts (C(52,5) = 2 598 960) and Powerball (C(69,5)·26 =
292 201 338). Answer to `num/stats` routes: row 12 → 003, row 27 → 009, row 66 → DROP-STDLIB (002, 004 instead).
