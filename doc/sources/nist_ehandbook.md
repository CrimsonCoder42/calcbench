# NIST/SEMATECH e-Handbook of Statistical Methods

**Citation**: NIST/SEMATECH, *e-Handbook of Statistical Methods*, §6.4.2.1 "Single Moving Average", https://www.itl.nist.gov/div898/handbook/pmc/section4/pmc421.htm, and §6.4.2 (the supplier data table), .../pmc42.htm. Read 2026-11-12.

| Taken | Used by | How checked |
|---|---|---|
| M_t = (X_t + X_{t-1} + … + X_{t−N+1}) / N | ALG-WINDOW-001, -002 | Formula quoted from the page |
| Supplier data 9, 8, 9, 12, 9, 12, 11, 7, 13, 9, 11, 10 (1000-dollar units); N = 3; printed averages 8.667, 9.667, 10.000, 11.000, 10.667, 10.000, 10.333, 9.667, 11.000, 10.000 | `TestRollingMean_nistExample` (to ±0.0005 of the printed values), `TestRollingSum_nistExample`; RollingMin/Max examples are hand-derived from the same data | Test passes |

## Chapters 1, 4 and 7: statistics (num/stats, math build 2026-11-12)

All pages under https://www.itl.nist.gov/div898/handbook/, each fetched 2026-11-12.

| Section | Taken | Used by | How checked |
|---|---|---|---|
| §1.3.5.1 "Measures of Location" (eda351.htm) | mean, median, mode, mid-range and its "not robust" note | NUM-STATS-001, -012, -013, -020, -024 | Resistivity data examples |
| §1.3.5.6 "Measures of Scale" (eda356.htm) | s² with N − 1, range, AAD = Σ\|Y − Ȳ\|/N, MAD = median\|Y − Ỹ\| unscaled | NUM-STATS-011, -018, -019, -020, -021, -031 | Worked examples |
| §1.3.3.14 "Histogram" (histogra.htm) | classes, normalized count | NUM-STATS-002, -008 | Worked examples |
| §4.4.3.1 "Least Squares" (pmd431.htm) | β̂1, β̂0 | NUM-STATS-009 | StRD Norris (see nist_strd.md) |
| §7.2.6.2 "Percentiles" (prc262.htm) | p(N+1) = k + d; 12 resistivity values; Y(90) 95.1981, R7 95.1957, R8 95.1972 | NUM-STATS-015; the resistivity data also in the -011, -012, -013, -018, -025, -031 tests | `TestPercentile_eHandbook` |
| §1.3.5.3 "Two-Sample t-Test for Equal Means" (eda353.htm) | s_p; AUTO83B N1 249, s1 6.41470, N2 79, s2 6.10771, s_p 6.34260 | NUM-STATS-016; NUM-STATS-004 test | Worked example |
| §1.3.5.2 "Confidence Limits for the Mean" (eda352.htm) | s/√N inside the CI; ZARR13 N 195, s 0.022789, interval 9.258242–9.264679 | NUM-STATS-023 | Worked example |

Status: confirmed (each page fetched; examples reproduced in the `_worked`/`_eHandbook` tests).

## §1.3.6: probability distributions (num/prob, math build 2026-11-12)

| Section | Taken | Used by | How checked |
|---|---|---|---|
| §1.3.6.6.1 Normal (eda3661.htm) | PDF; CDF and PPF "computed numerically" | NUM-PROB-047…050 | Normal table below |
| §1.3.6.6.2 Uniform (eda3662.htm) | general PDF 1/(B − A), mean (A+B)/2, SD √((B−A)²/12); CDF and PPF in the standard form only | NUM-PROB-044…046, -059 | Hand-checked; Ross §5.3 is the second source for the general CDF |
| §1.3.6.6.7 Exponential (eda3667.htm) | CDF 1 − e^(−x/β), PPF −β ln(1 − p), mean and SD β, λ = 1/β | NUM-PROB-054…056, -060 | Hand-checked |
| §1.3.6.6.9 Lognormal (eda3669.htm) | CDF Φ(ln x/σ), PPF exp(σΦ⁻¹(p)), mean e^(σ²/2), SD √(e^(σ²)(e^(σ²) − 1)) for m = 1 | NUM-PROB-051…053, -062 | Hand-checked |
| §1.3.6.6.18 Binomial (eda366i.htm) | PMF, CDF as the sum, mean np, SD √(np(1−p)) | NUM-PROB-026…028 | Exact big.Rat oracle |
| §1.3.6.6.19 Poisson (eda366j.htm) | PMF, mean λ, SD √λ | NUM-PROB-041…043 | Exact oracle; λ = 1000 value |
| §1.3.6.7.1 Normal table (eda3671.htm) | row 1.9, col 0.06 = 0.47500; critical values 0.975 → 1.960, 0.995 → 2.576, 0.999 → 3.090 | NUM-PROB-048, -049, -050, -061 | 0.47500 to 5 decimals; critical values to 3 decimals |

Status: every page opened 2026-11-12. Also found: §1.3.6.6.3 Cauchy (eda3663.htm) gives the PDF, the standard CDF 0.5 + arctan(x)/π and the PPF −cot(πp); Cauchy is sourced but LATER (D-013), no ID.


## Used by S10 (algorithms core, 2026-11-12)

IDs: ALG-STREAM-001, ALG-WINDOW-014.

Citations as written in the doc comments:
- NIST/SEMATECH, e-Handbook of Statistical Methods, §6.4.2.1 "Single Moving Average", https://www.itl.nist.gov/div898/handbook/pmc/section4/pmc421.htm (read 2026-11-12)
- NIST/SEMATECH, e-Handbook of Statistical Methods, §6.4.3.1 "Single Exponential Smoothing", https://www.itl.nist.gov/div898/handbook/pmc/section4/pmc431.htm (read 2026-11-12)

Checked: see the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ (confirmation status per citation; not first-hand items are listed in PROGRESS.md and REPORT.md).

## Matrix gaps (2026-11-12)

§6.5.4.1 "Mean Vector and Covariance Matrix" read 2026-11-12: S = (1/(n − 1)) Σ (Xᵢ − X̄)(Xᵢ − X̄)′; its example S (3×3, n = 5) is reproduced by the covariance ⇄ correlation round-trip test. Used by NUM-LINALG-167, 168.
