# Pending rows for shared files (combinatorics catalog pass, 2026-11-12)

The combinatorics agent does not edit shared files. Each block is in final form. Nothing is pasted before Devin
approves the catalog; INDEX rows are pasted only when the code ships.

## 1. INDEX.md rows (paste at build)

-009 and -010 carry SOURCE TO CONFIRM (Brualdi §2.2 not seen): their Source cell is final only after a copy is checked.

| ID | Function | Package | Does (one line) | Time | Space | Source | Tags |
|---|---|---|---|---|---|---|---|
| NUM-COMBI-001 | Factorial | num/combinatorics | n! exactly, n ≤ 20 | O(1) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.1; OEIS A000142, https://oeis.org/A000142 (read 2026-11-12) | combinatorics, factorial, counting |
| NUM-COMBI-002 | LogFactorial | num/combinatorics | natural log of n! | O(1) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §5.4, §5.11 eq. 5.11.1 | combinatorics, factorial, log-space, gamma |
| NUM-COMBI-003 | Binomial | num/combinatorics | C(n,k) exactly, overflow detected | O(min(k, n−k)) multiplies and gcds | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.3 eqs. 26.3.1, 26.3.2, 26.3.6, Table 26.3.1; OEIS A007318, https://oeis.org/A007318 (read 2026-11-12) | combinatorics, binomial, combinations, nCr, choose, counting |
| NUM-COMBI-004 | LogBinomial | num/combinatorics | natural log of C(n,k) | O(min(k, n−k, 64)) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.3 eq. 26.3.1, §5.4, §5.11 eq. 5.11.1 | combinatorics, binomial, log-space, probability |
| NUM-COMBI-005 | Permutations | num/combinatorics | nPr = n!/(n−k)! exactly | O(k) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.1 (falling factorial), §26.17 Table 26.17.1 | combinatorics, permutations, nPr, falling-factorial, counting |
| NUM-COMBI-006 | PermutationsWithRepetition | num/combinatorics | nᵏ sequences of length k over n symbols | O(log k) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.17 Table 26.17.1 | combinatorics, permutations, repetition, counting |
| NUM-COMBI-007 | CombinationsWithRepetition | num/combinatorics | multisets of size k from n kinds, C(n+k−1,k) | O(min(k, n−1)) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.17 Table 26.17.1 | combinatorics, combinations, repetition, multiset, counting |
| NUM-COMBI-008 | Multinomial | num/combinatorics | (Σcᵢ)!/Πcᵢ!, permutations of a multiset | O(m + Σ min(cⱼ, sⱼ₋₁))   m = len(counts) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.4 eq. 26.4.2, Table 26.4.1 | combinatorics, multinomial, multiset, permutations, counting |
| NUM-COMBI-009 | CircularPermutations | num/combinatorics | arrangements of k of n around a circle, rotations equal | O(k) | O(1), no allocation | Brualdi, Introductory Combinatorics, 5th ed., Pearson, 2010, §2.2 | combinatorics, permutations, circular |
| NUM-COMBI-010 | CircularPermutationsUndirected | num/combinatorics | circular arrangements, rotations and reflections equal | O(k) | O(1), no allocation | Brualdi, Introductory Combinatorics, 5th ed., Pearson, 2010, §2.2 | combinatorics, permutations, circular, necklace |
| NUM-COMBI-011 | Derangements | num/combinatorics | permutations of n with no fixed point, n ≤ 21 | O(1) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.13 eq. 26.13.4; OEIS A000166, https://oeis.org/A000166 (read 2026-11-12) | combinatorics, derangements, subfactorial, counting |
| NUM-COMBI-012 | Catalan | num/combinatorics | n-th Catalan number, n ≤ 36 | O(1) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.5 eqs. 26.5.1, 26.5.4, Table 26.5.1; OEIS A000108, https://oeis.org/A000108 (read 2026-11-12) | combinatorics, catalan, counting |
| NUM-COMBI-013 | StirlingFirstUnsigned | num/combinatorics | permutations of n with exactly k cycles | O(n·k) | O(k), no allocation for k ≤ 32 | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.8 eqs. 26.8.3, 26.8.18, Table 26.8.1; Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §6.1; OEIS A132393, https://oeis.org/A132393 (read 2026-11-12) | combinatorics, stirling, cycles, permutations |
| NUM-COMBI-014 | StirlingSecond | num/combinatorics | partitions of an n-set into k blocks | O(n·k) | O(k), no allocation for k ≤ 32 | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.8 eqs. 26.8.6, 26.8.22, Table 26.8.2; Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §6.1; OEIS A008277, https://oeis.org/A008277 (read 2026-11-12) | combinatorics, stirling, set-partitions |
| NUM-COMBI-015 | Bell | num/combinatorics | number of partitions of an n-set, n ≤ 25 | O(1) | O(1), no allocation | NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §26.7 eqs. 26.7.2, 26.7.6, Table 26.7.1; OEIS A000110, https://oeis.org/A000110 (read 2026-11-12) | combinatorics, bell, set-partitions |
| NUM-COMBI-016 | StirlingApproximation | num/combinatorics | √(2πn)(n/e)ⁿ, with Robbins's error bound | O(1) | O(1), no allocation | Robbins, "A Remark on Stirling's Formula", American Mathematical Monthly, 1955, 62(1):26–29, doi:10.2307/2308012, eqs. (1)–(2); NIST DLMF, https://dlmf.nist.gov (read 2026-11-12), §5.6 eq. 5.6.1, §5.11 eq. 5.11.1 | combinatorics, factorial, stirling, approximation, education |

Only if Q1 = "uint64 + big.Int twin" (same Source and Tags as the uint64 row, plus tag `big-int`; Time and
Space stated at build after measurement): NUM-COMBI-017 MultinomialBig · 018 DerangementsBig · 019 CatalanBig ·
020 StirlingFirstUnsignedBig · 021 StirlingSecondBig · 022 BellBig.

## 2. doc/scope/SCOPE.md slice row (paste on approval; slice number is Devin's)

| S? | NUM | Combinatorics: factorial, binomial (exact and log), permutations, repetition, multinomial, circular, derangements, Catalan, Stirling, Bell, Stirling's approximation (16 functions, NUM-COMBI-001…016; +6 big twins on Q1) | [combinatorics catalog](../catalog/combinatorics.md); spec to write from it | SPEC |

## 3. doc/taxonomy.md (proposed)

Add a row to "Groups and planned packages", after `algo/dp`:

| ALG | `algo/enumerate` | ENUM | generating permutations, combinations, subsets and partitions in order, as iterators |

Why: enumeration is an algorithm (no domain words), not a count; Knuth TAOCP vol. 4A §7.2.1 ("Generating basic
combinatorial patterns") is its source. Lost: `algo/dp` (enumeration is not dynamic programming; a reader
searching for "next permutation" would not look there); `num/combinatorics` (would mix counting with iteration
and break "algo imports nothing here" if it grew helpers).

## 4. Offers and answers to other subjects

| To | What | Note |
|---|---|---|
| `num/prob` | NUM-COMBI-003 Binomial, -004 LogBinomial, -001, -002, -005 | for binomial / hypergeometric pmf and counting rules. No request from prob existed at this pass |
| `num/prob` | MOVE rows 10, 15, 22, 23, 24 (conditional, total probability, A and B, A or B, non-independent) | probability rules |
| `life/games` | MOVE: Guess 1 to 31, Powerball (C(69,5)·26 = 292 201 338), the ten five-card hand probabilities and data items (C(52,5) = 2 598 960) | calls NUM-COMBI-003 |
| `num/integer` (proposed by algebra) | MOVE: Prime Number Theorem | if `num/integer` is refused, home is Devin's call |
| `num/stats` | answer to routes: row 12 → NUM-COMBI-003; row 27 → NUM-COMBI-009; row 66 Gamma → DROP-STDLIB (`math.Gamma`/`Lgamma`; NUM-COMBI-002, -004 for the log forms) | |

## 5. doc/sources/ files to create at build

`nist_dlmf.md` (add ch. 5, 26, 27 sections; exists only if geometry or another subject created it) ·
`robbins_1955.md` · `oeis.md` · `graham_knuth_patashnik_concrete_math.md` (§6.1 to confirm in a copy) ·
`brualdi_introductory_combinatorics_5.md` (§2.2 to confirm in a copy before -009/-010 ship).

APPLIED 2026-11-12 (lockdown pass; see doc/lockdown_2026-11-12/APPLY_LOG.md)
