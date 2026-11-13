# sample-size-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform cl select {70, 75, 80, 85, 90, 95, 98, 99, 99.9, 99.99, 99.999} (95), ci (5, margin %), pp (50, population proportion %), ps (population size, empty = infinite); calform2 cl2 (95), ss2 (100, sample size), pc2 (60, proportion %), ps2 (empty) → margin of error.
MATCH: NONE (NUM-PROB-061 NormalQuantile supplies z).
RESEARCH: Cochran, Sampling Techniques, 3rd ed., Wiley 1977, §4.4 "The Formula for n in Sampling for Proportions" (n₀ = t²PQ/d²; n = n₀/(1 + (n₀ − 1)/N); d = t√((N − n)/(N − 1))√(PQ/n)).
BUILD: NUM-STATS-047 SampleSizeProportion, NUM-STATS-048 MarginOfErrorProportion, num/stats/sample_size_proportion.go (z from prob.NormalQuantile: num/stats now imports num/prob, same group, no cycle); tests sample_size_proportion_test.go.

COMPARE (scratch/sample-size-calculator.go; 6 submissions). Ours is unrounded n; the sample is ⌈n⌉.

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 95 %, d 5 %, p 50 %, N ∞ | 384.146 → 385 | 385 | 0 | match |
| 95 %, d 5 %, p 50 %, N 1000 | 277.73 → 278 | 278 | 0 | match |
| 99 %, d 2 %, p 30 %, N ∞ | 3483.32 → 3484 | 3495 | 11 | THEIRS DIFFERS |
| 90 %, d 5 %, p 50 %, N 10000 | 263.45 → 264 | 266 | 2 | THEIRS DIFFERS |
| MOE 95 %, n 100, p 60 %, N ∞ | 0.0960182 (9.60 %) | 9.60 % | 0 | match |
| MOE 95 %, n 100, p 60 %, N 1000 | 0.0911365 (9.11 %) | 9.11 % | 0 | match |

DISCREPANCY: THEIRS DIFFERS. Their numbers are reproduced exactly with two-decimal z values: z = 2.58 at 99 % (2.58²·0.21/0.0004 = 3494.6 → 3495) and z = 1.65 at 90 % (n₀ = 272.25, n = 265.06 → 266). The 0.995 and 0.95 standard normal quantiles are 2.5758293… and 1.6448536… (Wichura, Algorithm AS 241, Applied Statistics 1988, 37(3):477–484, as NUM-PROB-061); Cochran §4.4 uses t, the normal deviate at the chosen confidence, not a rounded table value. Ours stays.
