# standard-deviation-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform numberinputs (textarea; default "10, 12, 23, 23, 16, 23, 21, 16"), ctype radio {p population, s sample}.
MATCH: HAVE NUM-STATS-019 SampleStandardDeviation, NUM-STATS-021 SampleVariance, NUM-STATS-022 PopulationStandardDeviation, PopulationVariance, NUM-STATS-001 ArithmeticMean.
COMPARE (scratch/standard-deviation-calculator.go; 5 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| default 8 values, population | σ 4.89897948556636, σ² 24, mean 18 | 4.8989794855664, 24, 18 | 0 at 13 dp | match |
| default 8 values, sample | s 5.23722936566382, s² 27.4285714285714 | 5.2372293656638, 27.428571428571 | 0 at 13 dp | match |
| 10000001, 10000003, 10000002 (NIST StRD NumAcc1: mean 10000002, s 1) | s 1, mean 10000002 | 1, 10000002 | 0 | match |
| 2, 4, 4, 4, 5, 5, 7, 9, population | σ 2 | 2 | 0 | match |
| 1000000000.2, .1, .3, .1, .3 (+1e9), sample | s 0.0999999642372167 | 0.099999964237231 | 0 at 15 dp | match |

Note on the last row: the decimal values are not representable at 1e9 (ulp 1.19e−7), so both sides compute the exact statistic of the stored binary values (0.09999996…), not of the decimals (0.1). Not a discrepancy; it is input representation.
