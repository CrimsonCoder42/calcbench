# The On-Line Encyclopedia of Integer Sequences

**Citation**: OEIS Foundation, *The On-Line Encyclopedia of Integer Sequences*, https://oeis.org (each sequence page read 2026-11-12).

| Sequence | Taken | Used by | How checked |
|---|---|---|---|
| A000142 | Factorials | NUM-COMBI-001 | Table constants vs math/big |
| A007318 | Pascal's triangle | NUM-COMBI-003 | Table 26.3.1 values; oracle |
| A000166 | Derangements; d(22) = 413 496 759 611 120 779 881 | NUM-COMBI-011, -018 | math/big recurrence; limit n ≤ 20 (d(21) overflows uint64) |
| A000108 | Catalan numbers | NUM-COMBI-012, -019 | math/big oracle |
| A000110 | Bell numbers; B(26) | NUM-COMBI-015, -022 | math/big oracle |
| A132393 | Unsigned Stirling numbers of the first kind | NUM-COMBI-013, -020 | math/big recurrence |
| A008277 | Stirling numbers of the second kind | NUM-COMBI-014, -021 | math/big recurrence |
| A111492 | a(n,k) = (k−1)!·C(n,k) = P(n,k)/k, k-cycles of S_n; row 5: 5 10 20 30 24 | NUM-COMBI-009 | Row 5 is the first test |
| A144151 | T(n,k) = C(n,k) for k ≤ 2, else C(n,k)·(k−1)!/2, undirected k-cycles on n labelled nodes (A. P. Heinz 2008); row 5: 1 5 10 10 15 12 | NUM-COMBI-010 | Row 5 is the first test |
| A001622 | Decimal expansion of the golden ratio φ = (1+√5)/2 | NUM-LOGS-001 (second source) | `TestGoldenRatio_matchesMathPhi` |

Status: A111492 and A144151 pages read 2026-11-12 (name, formula, data rows); they replaced Brualdi, *Introductory Combinatorics* 5th ed. §2.2, which could not be confirmed (the 5th ed.'s ch. 2 is "The Pigeonhole Principle") and gets no file. A001622 read 2026-11-12 by the lead.
