# C10 (pending weighted.go:107) → NUM-STATS-037 FrequencyWeightedMedian (PORTED)
- Citation: NIST/SEMATECH e-Handbook §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm: median "if N is odd" Y((N+1)/2), "if N is even" (Y(N/2) + Y(N/2+1))/2 — applied to the data expanded by the frequencies (the source's definition: "a value counted N times gives what N copies in a slice would").
- Selection: CLRS 3rd ed. §9.2, Problem 9-2 (weighted median); Musser 1997, SP&E 27(8):983-993 (work budget, fallback).
- Checks: OpenStax 2e Example 2.32 table → 10.5; source-shaped cases 1, 5, 3; oracle vs NUM-STATS-012 on expanded data. Read 2026-11-12.
