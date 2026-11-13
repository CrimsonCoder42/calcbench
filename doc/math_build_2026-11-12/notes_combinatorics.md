# Hand-off: num/combinatorics (NUM-COMBI-001…022), 2026-11-12

Builder: combinatorics. Owned paths only: `num/combinatorics/**`, `doc/specs/s4_combinatorics.md`, this file.
Checks run while building: `gofmt -l num/combinatorics` (clean), `go build ./...` (ok), `go vet ./num/combinatorics/` (ok).
Tests and benchmarks NOT run (lead verifies). 22 functions built, 0 held back.

## (a) INDEX.md rows (final form, generated from the doc comments and doc.go)

| ID | Function | Package | Does (one line) | Time | Space | Source | Tags |
|---|---|---|---|---|---|---|---|
| NUM-COMBI-001 | Factorial | num/combinatorics | n! exactly, n ≤ 20 | O(1)   (a 21-entry table) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §5.4 eq. 5.4.1 (n! = Γ(n+1)); OEIS A000142, https://oeis.org/A000142 (read 2026-11-12) | combinatorics, factorial, counting |
| NUM-COMBI-002 | LogFactorial | num/combinatorics | natural log of n! | O(1) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §5.4 eq. 5.4.1 (n! = Γ(n+1)), §5.11 eq. 5.11.1 | combinatorics, factorial, log-space, gamma |
| NUM-COMBI-003 | Binomial | num/combinatorics | C(n,k) exactly, overflow detected | O(m log n)   m = min(k, n−k, 64) steps (the partial result C(n−k+i, i) ≥ 2ⁱ, so a result that fits needs fewer than 64), each with one O(log n) binary gcd | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.3 eqs. 26.3.1, 26.3.2, 26.3.6, Table 26.3.1; OEIS A007318, https://oeis.org/A007318 (read 2026-11-12) | combinatorics, binomial, combinations, nCr, choose, counting |
| NUM-COMBI-004 | LogBinomial | num/combinatorics | natural log of C(n,k) | O(min(k, n−k, 64)) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.3 eq. 26.3.1, §5.4 eq. 5.4.1, §5.11 eq. 5.11.1 | combinatorics, binomial, log-space, probability |
| NUM-COMBI-005 | Permutations | num/combinatorics | nPr = n!/(n−k)! exactly | O(min(k, 64))   (every factor but the last is ≥ 2, so overflow stops the loop within 64 multiplies) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.1 (falling factorial), §26.17 Table 26.17.1 (one-to-one maps) | combinatorics, permutations, nPr, falling-factorial, counting |
| NUM-COMBI-006 | PermutationsWithRepetition | num/combinatorics | nᵏ sequences of length k over n symbols | O(log k)   (square and multiply) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.17 Table 26.17.1 (unrestricted maps of labelled balls into labelled boxes) | combinatorics, permutations, repetition, counting |
| NUM-COMBI-007 | CombinationsWithRepetition | num/combinatorics | multisets of size k from n kinds, C(n+k−1,k) | O(m log n)   m = min(k, n−1, 64), as Binomial on (n+k−1, k) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.17 Table 26.17.1 (unlabelled balls into labelled boxes, unrestricted) | combinatorics, combinations, repetition, multiset, counting |
| NUM-COMBI-008 | Multinomial | num/combinatorics | (Σcᵢ)!/Πcᵢ!, permutations of a multiset | O(m + Σⱼ min(cⱼ, sⱼ₋₁, 64)·log s)   m = len(counts), sⱼ = c₁+…+cⱼ, s = sₘ | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.4 eq. 26.4.2 | combinatorics, multinomial, multiset, permutations, counting |
| NUM-COMBI-009 | CircularPermutations | num/combinatorics | arrangements of k of n around a circle, rotations equal | O(m log n)   m = min(k, n−k, 64), as Binomial | O(1), no allocation | OEIS A111492, https://oeis.org/A111492 (read 2026-11-12), a(n,k) = (k−1)!·C(n,k) = P(n,k)/k | combinatorics, permutations, circular, cycles |
| NUM-COMBI-010 | CircularPermutationsUndirected | num/combinatorics | circular arrangements, rotations and reflections equal | O(m log n)   m = min(k, n−k, 64), as Binomial | O(1), no allocation | OEIS A144151, https://oeis.org/A144151 (read 2026-11-12), T(n,k) = C(n,k) for k ≤ 2, else C(n,k)·(k−1)!/2 | combinatorics, permutations, circular, necklace, cycles |
| NUM-COMBI-011 | Derangements | num/combinatorics | permutations of n with no fixed point, n ≤ 20 | O(1)   (a 21-entry table) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.13 eq. 26.13.4; OEIS A000166, https://oeis.org/A000166 (read 2026-11-12) | combinatorics, derangements, subfactorial, counting |
| NUM-COMBI-012 | Catalan | num/combinatorics | n-th Catalan number, n ≤ 36 | O(1)   (a 37-entry table) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.5 eqs. 26.5.1, 26.5.4, Table 26.5.1; OEIS A000108, https://oeis.org/A000108 (read 2026-11-12) | combinatorics, catalan, counting |
| NUM-COMBI-013 | StirlingFirstUnsigned | num/combinatorics | permutations of n with exactly k cycles | O(n·w)   w = min(n−k+1, 64) (one band of the triangle; the loop also stops at the first overflow) | O(1), no allocation (a fixed 64-entry array) | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.8 eqs. 26.8.3, 26.8.18, Table 26.8.1; Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §6.1 "Stirling Numbers" (table of [n k], p. 245); OEIS A132393, https://oeis.org/A132393 (read 2026-11-12) | combinatorics, stirling, cycles, permutations |
| NUM-COMBI-014 | StirlingSecond | num/combinatorics | partitions of an n-set into k blocks | O(n·w)   w = min(n−k+1, 64) (one band of the triangle; the loop also stops at the first overflow) | O(1), no allocation (a fixed 64-entry array) | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.8 eqs. 26.8.6, 26.8.22, Table 26.8.2; Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §6.1 "Stirling Numbers" (table of {n k}, p. 244); OEIS A008277, https://oeis.org/A008277 (read 2026-11-12) | combinatorics, stirling, set-partitions |
| NUM-COMBI-015 | Bell | num/combinatorics | number of partitions of an n-set, n ≤ 25 | O(1)   (a 26-entry table) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.7 eqs. 26.7.2, 26.7.6, Table 26.7.1; OEIS A000110, https://oeis.org/A000110 (read 2026-11-12) | combinatorics, bell, set-partitions |
| NUM-COMBI-016 | StirlingApproximation | num/combinatorics | √(2πn)(n/e)ⁿ, with Robbins's error bound | O(1) | O(1), no allocation | Robbins, "A Remark on Stirling's Formula", American Mathematical Monthly, 1955, 62(1):26–29, doi:10.2307/2308012 (the bound 1/(12n+1) < rₙ < 1/(12n)); NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §5.6 eq. 5.6.1, §5.11 eq. 5.11.1 | combinatorics, factorial, stirling, approximation, education |
| NUM-COMBI-017 | MultinomialBig | num/combinatorics | Multinomial as a *big.Int | O(m) big.Int.Binomial calls and multiplies   m = len(counts); each on numbers of at most b bits, b = bit length of the result | O(b) bits | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.4 eq. 26.4.2 | combinatorics, multinomial, multiset, permutations, counting, big-int |
| NUM-COMBI-018 | DerangementsBig | num/combinatorics | Derangements as a *big.Int | O(n·b) bit operations   b = bit length of the result, O(n log n) | O(b) bits | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.13 eq. 26.13.4; OEIS A000166, https://oeis.org/A000166 (read 2026-11-12) | combinatorics, derangements, subfactorial, counting, big-int |
| NUM-COMBI-019 | CatalanBig | num/combinatorics | Catalan as a *big.Int | O(n) big multiplies (one big.Int.Binomial(2n, n)) and one division, on numbers of at most b bits   b = bit length of the result, O(n) | O(b) bits | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.5 eqs. 26.5.1, 26.5.4, Table 26.5.1; OEIS A000108, https://oeis.org/A000108 (read 2026-11-12) | combinatorics, catalan, counting, big-int |
| NUM-COMBI-020 | StirlingFirstUnsignedBig | num/combinatorics | StirlingFirstUnsigned as a *big.Int | O(n·w) big multiply-adds   w = min(k, n−k) + 1; each on numbers of at most b bits, b = bit length of the result | O(w·b) bits (w big integers, one band of the triangle) | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.8 eqs. 26.8.3, 26.8.18, Table 26.8.1; Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §6.1 "Stirling Numbers" (table of [n k], p. 245); OEIS A132393, https://oeis.org/A132393 (read 2026-11-12) | combinatorics, stirling, cycles, permutations, big-int |
| NUM-COMBI-021 | StirlingSecondBig | num/combinatorics | StirlingSecond as a *big.Int | O(n·w) big multiply-adds   w = min(k, n−k) + 1; each on numbers of at most b bits, b = bit length of the result | O(w·b) bits (w big integers, one band of the triangle) | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.8 eqs. 26.8.6, 26.8.22, Table 26.8.2; Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §6.1 "Stirling Numbers" (table of {n k}, p. 244); OEIS A008277, https://oeis.org/A008277 (read 2026-11-12) | combinatorics, stirling, set-partitions, big-int |
| NUM-COMBI-022 | BellBig | num/combinatorics | Bell as a *big.Int | O(n²) big multiply-adds by a small integer, each on numbers of at most b bits   b = bit length of the result, O(n log n) | O(n·b) bits (one row of n+1 big integers) | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.7 eqs. 26.7.2, 26.7.6, Table 26.7.1; OEIS A000110, https://oeis.org/A000110 (read 2026-11-12) | combinatorics, bell, set-partitions, big-int |

## (b) doc/sources/ records (for the lead to fold in)

### `nist_dlmf.md` (add or create a "Chapters 5 and 26" section)
- Citation: NIST Digital Library of Mathematical Functions, F. W. J. Olver et al., eds., https://dlmf.nist.gov
  (read 2026-11-12). Ch. 5 Gamma Function (Askey, Roy); ch. 26 Combinatorial Analysis (Bressoud).
- Taken: eq. 5.4.1 (n! = Γ(n+1)) → 001, 002, 004; eq. 5.6.1 (1 < (2π)^(−½)x^(½−x)eˣΓ(x) < e^(1/(12x))) → 016;
  eq. 5.11.1 (Stirling series for ln Γ, coefficients 1/12, −1/360, 1/1260, −1/1680) → 002, 004, 016;
  §26.1 falling factorial → 005; eqs. 26.3.1, 26.3.2, 26.3.6, Table 26.3.1 → 003; eq. 26.4.2 → 008, 017;
  eqs. 26.5.1, 26.5.4, Table 26.5.1 → 012, 019; eqs. 26.7.2, 26.7.6, Table 26.7.1 → 015, 022;
  eqs. 26.8.3, 26.8.6, 26.8.18, 26.8.22, Tables 26.8.1, 26.8.2 → 013, 014, 020, 021, 022;
  eq. 26.13.4 → 011, 018; §26.17 Table 26.17.1 (twelvefold way: kⁿ, one-to-one, C(k+n−1, n)) → 005, 006, 007.
- Checked: each equation re-read on dlmf.nist.gov on 2026-11-12 during the build. Tests reproduce Table 26.3.1
  (C(6,3)=20, C(8,2)=28, C(10,5)=252), Table 26.5.1 (C(10)=16 796), Table 26.7.1 (B(5)=52, B(10)=115 975),
  Table 26.8.1 (s(6,3)=−225, s(10,5)=−269 325), Table 26.8.2 (S(6,3)=90, S(10,5)=42 525); 26.13.4, 26.5.4, 26.7.6,
  26.8.6 are coded as independent math/big oracles in the tests.

### `oeis.md`
- Citation: The On-Line Encyclopedia of Integer Sequences, https://oeis.org (read 2026-11-12).
- Taken: A000142 (factorials) → 001; A007318 (Pascal) → 003; A000166 (derangements; d(22) = 413 496 759 611 120 779 881)
  → 011, 018; A000108 (Catalan) → 012, 019; A000110 (Bell; B(26)) → 015, 022; A132393 (unsigned Stirling 1st) →
  013, 020; A008277 (Stirling 2nd) → 014, 021; **A111492** (a(n,k) = (k−1)!·C(n,k) = P(n,k)/k, k-cycles of S_n;
  row 5: 5 10 20 30 24) → 009; **A144151** (T(n,k) = C(n,k) for k ≤ 2, else C(n,k)·(k−1)!/2, undirected k-cycles
  on n labelled nodes, A. P. Heinz 2008; row 5: 1 5 10 10 15 12) → 010.
- Checked: A111492 and A144151 pages read 2026-11-12 (name, formula, data rows); their row 5 is the first test of
  009/010. Table constants checked in tests against math/big recurrences.

### `robbins_1955.md`
- Citation: H. Robbins, "A Remark on Stirling's Formula", The American Mathematical Monthly, 1955, 62(1):26–29,
  doi:10.2307/2308012.
- Taken: n! = √(2π)n^(n+½)e^(−n)e^(rₙ), 1/(12n+1) < rₙ < 1/(12n) → 016.
- Checked: bibliographic data and bound confirmed via the Wikipedia "Stirling's approximation" reference list
  (2026-11-12); the catalog pass read the JSTOR scan. Equation numbers inside the paper are NOT cited (not seen in
  this build). Test: the bound holds for every n = 1…170 against exact n!.

### `graham_knuth_patashnik_concrete_math.md`
- Citation: R. L. Graham, D. E. Knuth, O. Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994,
  §6.1 "Stirling Numbers" (table of {n k} p. 244, table of [n k] p. 245).
- Taken: second, independent source for the Stirling definitions and recurrences → 013, 014, 020, 021.
- Checked (not seen in a copy): §6.1 title from a search result of the book's contents (ch. 6 "Special Numbers",
  §6.1 "Stirling Numbers", §6.2 "Eulerian Numbers"); Pearson product page (2nd ed., 1994, ch. 6 "Special Numbers");
  OEIS wiki "Sequences from Graham, Knuth, Patashnik Concrete Math" (A008277 on p. 244, A008275 on p. 245);
  Wikipedia "Stirling numbers of the second kind" (cites "section 6.1 of Concrete Mathematics", p. 244).

### `brualdi_introductory_combinatorics_5.md`: NOT to be created (citation dropped, see (c)).

## (c) Citations

Confirmed (how):
- DLMF 5.4.1, 5.6.1, 5.11.1, §26.1, 26.3.1/2/6 + Table 26.3.1, 26.4.2, 26.5.1/4 + Table 26.5.1, 26.7.2/6 + Table
  26.7.1, 26.8.3/6/18/22 + Tables 26.8.1/2, 26.13.4, Table 26.17.1: fetched from dlmf.nist.gov 2026-11-12.
- OEIS A111492, A144151: fetched 2026-11-12.
- Concrete Mathematics §6.1, pp. 244–245: two independent secondary sources (OEIS wiki page list; Wikipedia) plus
  the contents listing, as above.
- Robbins 1955: bibliographic data and bound via Wikipedia's reference; JSTOR page itself did not load.

NOT confirmed, dropped:
- Brualdi, Introductory Combinatorics, 5th ed., §2.2 (catalog's source for 009/010). Pearson's contents show
  chapter 2 is "The Pigeonhole Principle" and chapter 3 "Permutations and Combinations", so §2.2 is wrong for the
  5th ed.; the circular-permutation section/theorem number could not be seen. Replaced by OEIS A111492 (009) and
  A144151 (010), both confirmed.
- DLMF §26.1 for n! (catalog): §26.1 defines the falling factorial only; 001 cites DLMF 5.4.1 instead.
- DLMF Table 26.4.1 for 008 (catalog): the table layout could not be read reliably; 008 cites eq. 26.4.2 only.
- Robbins eqs. (1)–(2) locators: not seen; the article is cited without equation numbers.

Held back: none.

## (d) FOUND, NOT FIXED

- `doc/catalog/combinatorics.md` NUM-COMBI-011 says d(21) = 18 795 307 255 050 944 540 fits a uint64 and n ≤ 21.
  It does not (2⁶⁴ − 1 = 18 446 744 073 709 551 615). Built as n ≤ 20; the catalog and PENDING_SHARED.md row
  ("n ≤ 21") are wrong.
- `doc/combinatorics_2026-11-12/PENDING_SHARED.md` §1 rows differ from the built rows (Time/Space/Source cells,
  009/010 sources); use the rows in (a), not those.
- Catalog STEP 1 rows 45/46 point at "63"/"64" in STEP 2's Covers column (should be 45/46); cosmetic.
- API notes for the lead: 009/010 return ErrOutOfRange for k = 0 (A144151 gives T(n,0) = 1; the catalog's choice
  was kept). The big twins return ErrOverflow only when an input sum (Σcounts, 2n) exceeds MaxInt64.
- Complexity notes differing from the catalog: Binomial is O(min(k, n−k, 64)·log n) (overflow bounds the loop);
  Permutations O(min(k, 64)); Stirling uint64 forms O(n·min(n−k+1, 64)) with no allocation at any k (k ≥ 2 and
  n−k ≥ 64 is an immediate ErrOverflow since both kinds are ≥ k^(n−k)); the catalog said O(n·k) with allocation
  for k > 32.
- Benchmarks to watch in verify: BenchmarkMultinomialBig, DerangementsBig, CatalanBig, BellBig,
  Stirling*Big at growing sizes (growth per the Time lines); every uint64/float64 function has an AllocsPerRun = 0 test.
