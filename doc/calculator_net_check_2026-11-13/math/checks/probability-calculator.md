# probability-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform cal1pa (0.5), cal1pb (0.4) — two independent events; calform3 cal3pa, cal3pb, cal3na, cal3nb, cal3and (0.4), cal3or, cal3xor, cal3nab — solve from any two; calform4 cal4pa (0.6), cal4par (5), cal4pb (0.3), cal4pbr (3) — repeated independent trials; calform2 val2mean (0), val2deviation (1), val2lb (−1), val2rb (1) — normal interval.
MATCH: HAVE NUM-PROB-002 ComplementProbability, -003 IndependentJointProbability, -005 UnionProbability, -006 ConditionalProbability, -050 NormalIntervalProbability (and -048 NormalCDF). Repeated trials (pᵏ) and the solve-from-any-two form (calform3) are compositions, not compared.
COMPARE (scratch/probability-calculator.go; 3 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| P(A) 0.37, P(B) 0.81, independent | A' 0.63, A∩B 0.2997, A∪B 0.8803, AΔB 0.5806, neither 0.1197 | 0.63, 0.2997, 0.8803, 0.5806, 0.1197 | 0 | match |
| N(0, 1), −1.96 ≤ X ≤ 1.96 | 0.950004209703559 | 0.95000 | 0 at 5 dp | match |
| N(100, 15), 110 ≤ X ≤ 130 | 0.2297424055987437 | 0.22974 | 0 at 5 dp | match |
