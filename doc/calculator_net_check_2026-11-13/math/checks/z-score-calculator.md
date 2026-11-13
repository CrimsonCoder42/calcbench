# z-score-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform c1raw (5), c1mean (3), c1sd (2); calform2 c2z (2) or one of c2p, c2pg, c2p0, c2pin, c2pout (probability → z); calform3 c3z1 (−1), c3z2 (0).
MATCH: HAVE NUM-STATS-028 ZScore, NUM-PROB-048 NormalCDF, -049 NormalSurvival, -050 NormalIntervalProbability (-061 NormalQuantile for the inverse forms).
COMPARE (scratch/z-score-calculator.go; 4 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| x −12.5, μ 4, σ 7.25 | z −2.2758620689655173; Φ(z) 0.0114270 (cross-check: 0.5·erfc) | −2.27586; 0.011427 | 0 at 5 dp | match |
| z 2 | P(<z) 0.9772498681, P(>z) 0.02275013195 | 0.97725, 0.02275 | 0 at 5 dp | match |
| z −3.5 | P(<z) 0.000232629079, P(>z) 0.9997673709 | 0.00023263, 0.99977 | 0 at displayed sig. figs | match |
| −1 < z < 0 | 0.3413447460685429 | 0.34134 | 0 at 5 dp | match |
