# permutation-and-combination-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform cnv (n, 6), crv (r, 2).
MATCH: HAVE NUM-COMBI-005 Permutations, NUM-COMBI-003 Binomial (exact uint64, overflow detected).
COMPARE (scratch/permutation-and-combination-calculator.go; 4 submissions)

| inputs n, r | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 6, 2 | P 30, C 15 | 30, 15 | 0 | match |
| 52, 5 | P 311875200, C 2598960 | 311,875,200, 2,598,960 | 0 | match |
| 66, 33 | P ErrOverflow (> 2⁶⁴), C 7219428434016265740 | P 6.27e61 (exact big integer), C 7,219,428,434,016,265,740 | C 0; P by design | match (C); P out of uint64 range by contract |
| 20, 20 | P 2432902008176640000, C 1 | 2,432,902,008,176,640,000, 1 | 0 | match |

Gap: their nPr is arbitrary precision; ours stops at 2⁶⁴ − 1 with ErrOverflow (documented contract). NUM-COMBI-002/-004 give the logarithms beyond that. Noted in LEAD.md.
