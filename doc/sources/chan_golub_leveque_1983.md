# Chan, Golub, LeVeque 1983

**Citation**: Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", *The American Statistician*, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115.

| Taken | Locator | Used by |
|---|---|---|
| The corrected two-pass algorithm Σd² − (Σd)²/n | (no locator in the stats Source lines) | NUM-STATS-017, -019, -021, -022, -025 (and internally 003, 009, 014, 020, 023, 030) |
| The corrected two-pass algorithm, weighted form Σ pᵢdᵢ² − (Σ pᵢdᵢ)² | eq. (1.7) | NUM-PROB-019 (discrete_variance.go:21, added 2026-11-12) |

How checked: stats, StRD NumAcc fixtures and oracle tests vs 512-bit big.Float. NUM-PROB-019, `TestDiscreteVariance_correctionTerm` (discrete_variance_test.go:39; {1e9+0.1, 1e9+0.3}, ½ each, exact big.Rat, 4.5e-16 relative).

Status: citation confirmed (Crossref). Eq. (1.7) read first-hand in the fix pass of 2026-11-12 (fix_prob.md §3) and cited by NUM-PROB-019; the NUM-STATS Source lines still cite the paper without an equation number.
