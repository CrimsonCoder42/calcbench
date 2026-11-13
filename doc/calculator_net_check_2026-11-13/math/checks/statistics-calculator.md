# statistics-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform numberinputs (textarea; default "10, 2, 38, 23, 38, 23, 21, 23"). Outputs: count, sum, mean, median, mode, largest, smallest, range, geometric mean, σ, σ², s, s², sorted data.
MATCH: HAVE NUM-STATS-001 ArithmeticMean, -012 Median, -018 Range, -019 SampleStandardDeviation, -021 SampleVariance, -022 PopulationStandardDeviation, PopulationVariance, -024 Mode, -042 Percentiles. PARTIAL: no plain geometric mean (only the out-of-scope NUM-STATS-027 WeightedGeometricMean reservation).
RESEARCH: NIST, Dataplot Reference Manual, "GEOMETRIC MEAN", XGM = (∏Xᵢ)^(1/n), https://itl.nist.gov/div898/software/dataplot/refman2/auxillar/geommean.htm (read 2026-11-13).
BUILD: NUM-STATS-049 GeometricMean, num/stats/geometric_mean.go (exp of the compensated mean of logs; zero → 0; negative → ErrNegative); tests geometric_mean_test.go.

COMPARE (scratch/statistics-calculator.go; 2 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 10, 2, 38, 23, 38, 23, 21, 23 | n 8, mean 22.25, median 23, mode 23 ×3, range 36, GM 17.119851726053064, σ 11.5081492864839, s 12.3027290816771 | 8, 22.25, 23, 23 ×3, 36, 17.119851726053, 11.508149286484, 12.302729081677 | 0 at 12 dp | match |
| 1.5, −2.25, 7, 7, 3.125, −2.25, 0 | mean 2.01785714285714, median 1.5, modes −2.25, 7 ×2, range 9.25, σ 3.61817619521104, s 3.90807695371911, GM ErrNegative | 2.0178571428571, 1.5, −2.25, 7 ×2, 9.25, 3.618176195211, 3.9080769537191, GM 0 | GM only | THEIRS DIFFERS (GM) |

DISCREPANCY: THEIRS DIFFERS on the geometric mean of data with negative values: they print 0. The NIST definition (∏Xᵢ)^(1/n) has no real value for a product of mixed-sign factors in general (here the product is 0 because of the 0, so 0 is right for this set by the formula; but ours returns ErrNegative first because a negative value is present). Ours keeps the domain check (non-negative values); recorded for Devin as a policy point: return 0 when any value is 0 even if negatives are present? Kept as is.
