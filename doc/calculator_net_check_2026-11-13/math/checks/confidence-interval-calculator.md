# confidence-interval-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform size (50), mean (20.6), sd (3.2), cl (95, percent).
MATCH: PARTIAL NUM-STATS-023 StandardErrorOfMean, NUM-PROB-061 NormalQuantile, NUM-PROB-079 StudentTQuantile; no interval function.
RESEARCH: Montgomery & Runger, Applied Statistics and Probability for Engineers, 6th ed., Wiley 2014, §8.1 (variance known, z) and §8.2 (variance unknown, t); NIST/SEMATECH e-Handbook §7.2.2.1 (t limits).
BUILD: NUM-STATS-050 MeanConfidenceIntervalZ, NUM-STATS-051 MeanConfidenceIntervalT, num/stats/mean_confidence_interval.go; tests mean_confidence_interval_test.go.

COMPARE (scratch/confidence-interval-calculator.go; 3 submissions). They use z at every n (they print a note below n = 30 but still use z), so theirs is compared with MeanConfidenceIntervalZ.

| inputs n, x̄, sd, cl | ours (Z) | theirs | diff | verdict |
|---|---|---|---|---|
| 50, 20.6, 3.2, 95 % | 20.6 ± 0.886978 [19.713022, 21.486978] | ± 0.887 [19.713, 21.487] | 0 at 3 dp | match |
| 10, 10, 2, 95 % | 10 ± 1.239590 [8.760410, 11.239590] | ± 1.240 [8.760, 11.240] | 0 at 3 dp | match |
| 30, −4.5, 12.25, 99 % | −4.5 ± 5.760929 [−10.260929, 1.260929] | ± 5.761 [−10.261, 1.261] | 0 at 3 dp | match |

Note: with a sample s and n = 10 the t interval (ours T: ± 1.430714) is the one Montgomery §8.2 prescribes; their page uses z. Not a discrepancy in the formula they state (z with σ known), so no DISCREPANCIES row; noted here only.
