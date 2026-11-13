# NIST Statistical Reference Datasets (StRD)

**Citation**: NIST, *Statistical Reference Datasets*, Information Technology Laboratory, https://www.itl.nist.gov/div898/strd/ (read 2026-11-12).

| Taken | Used by (tests) | How checked |
|---|---|---|
| Univariate Summary Statistics: NumAcc1–4, Michelso, Mavro, Lew (certified mean and sample SD) | NUM-STATS-001, -003, -019, -020, -021, -023, -025, -030 | Certified values reproduced; e.g. Michelso (n = 100) sample SD 0.0790105478190518; NumAcc1 sum of squares 2 |
| Linear Least Squares: Norris (B0 −0.262323073774029, B1 1.00211681802045, R² 0.999993745883712) | NUM-STATS-009, -014 | r = +√R² = 0.999996872936967 (hand-checked) |

Fixtures are Go literals in num/stats/arithmetic_mean_test.go. Status: data transcribed (the raw download was blocked; NumAcc2–4 built from the page's stated pattern) and verified in exact rationals against every certified value. PiDigits NOT used: the transcription could not reproduce NIST's certified values.
