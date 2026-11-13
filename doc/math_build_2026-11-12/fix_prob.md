# Fix pass: num/prob (2026-11-12)

Fixer: prob. Owned paths: `num/prob/**`, `doc/specs/s9_prob.md`, this file. No public signature changed.
Final: `gofmt -l num/prob` clean · `go vet ./num/prob/` ok · `go test -count=1 ./num/prob/` →
`ok  github.com/CrimsonCoder42/calcbench/num/prob  8.000s`.

## 1. HIGH: odds sources (D-015), NUM-PROB-014…017

Štrumbelj 2014 still unreadable (ScienceDirect 403; no preprint on arXiv/RePEc/the author's pages found). It is
removed as a Source. Each formula is now carried by peer-reviewed papers read first-hand, with locators:

| ID | Confirmed by (read 2026-11-12) |
|---|---|
| 014 DecimalOddsFromProbability | Wheatcroft preprint §2.1 (rᵢ = 1/Oᵢ); Wunderlich–Memmert, Method, "Transferring betting odds to probabilities" ("If no bookmaker margin was contained in the betting odds, the inverse betting odds … could be interpreted as its probability") |
| 015 ImpliedProbability | Wheatcroft §2.1 (rᵢ = 1/Oᵢ, example O = 3, 1.4 → 0.33, 0.71); Egidi et al. §2 ("The inverse odd … corresponds to the unfair probability") |
| 016 Overround | Wheatcroft §2.1 eq. (1) π = (Σ 1/Oᵢ) − 1, "the overround"; Egidi et al. §2 eq. (2.1) booksum β = Σ inverse odds |
| 017 NormalizedImpliedProbabilities | Egidi et al. §2 eq. (2.1) πᵢ = oᵢ/β ("basic normalisation"); Wunderlich–Memmert, same subsection (basic normalisation) |

Published versions checked on Crossref: Wheatcroft, J. Sports Analytics 7(2):77–97, 2021; Egidi, Pauli, Torelli,
Statistical Modelling 18(5–6):436–459, 2018. Locators are from the author preprints (arXiv:2001.09097v1,
arXiv:1802.08848, text extracted with pdftotext) and say so. **HELD: none. Moved to _delete/: nothing.**
017's prose now says Shin's method is better "(Štrumbelj 2014, as reported by Egidi et al. 2018, §1 of the
preprint)"; the favourite–longshot claim, which was not confirmed, is gone.
Files: decimal_odds_from_probability.go:13-14, implied_probability.go:13-14, overround.go:30-31,
normalized_implied_probabilities.go:7-9,15-16; test comments in the four _test.go files.
Tests added: implied_probability_test.go:19 `TestImpliedProbability_wheatcroftExample` and overround_test.go:18
`TestOverround_wheatcroftExample` (π = 1/3 + 1/1.4 − 1 = 1/21), both from Wheatcroft §2.1.

## 2. Altman & Bland 1994, Casella–Berger numbers

- Altman–Bland 1994 is a page scan only (PMC2540558, PDF not fetchable). The Source line is dropped from
  positive_predictive_value.go:15 and negative_predictive_value.go:14; Ross §3.3 remains. The example is now an
  unattributed hand-checked 2×2 table (TP 231, FP 32, FN 27, TN 54), and the test comments were changed to match.
- Casella–Berger result numbers were not confirmed and are dropped. Sections only, with titles confirmed from
  the TOC: complement_probability.go:12, union_probability.go:14, total_probability.go:101 → §1.2;
  conditional_probability.go:13, independent_joint_probability.go:11, bayes_theorem.go:13, bayes_posterior.go:20
  → §1.3. Six test comments were changed the same way.

## 3. LOW: DiscreteVariance correction term

discrete_variance.go:33-46 now uses the corrected two-pass algorithm, Σ pᵢdᵢ² − (Σ pᵢdᵢ)². dᵢ² also keeps its
FMA rounding error. The locator was read first-hand: Chan, Golub, LeVeque 1983, eq. (1.7). Higham §1.9 is kept
for the two-pass context only.
Test: discrete_variance_test.go:39 `TestDiscreteVariance_correctionTerm` ({1e9+0.1, 1e9+0.3}, ½ each, exact big.Rat,
tolerance 4.5e-16 relative; the old code is off by 3.6e-13). The oracle tolerance went from 1e-11 to 1e-13, as its
comment already stated, and it passes.

## 4. LOW: representable results reported as ErrOverflow

- Uniform: uniform_pdf.go:5-20 checkUniform no longer fails when b − a overflows. When it does, the functions
  work on the halved endpoints, which are exact there:
  - PDF ½/(b/2 − a/2) (uniform_pdf.go:47). The result is ≥ 2⁻¹⁰²⁵ and never 0. ErrOverflow stays only when the
    density overflows (b − a subnormal).
  - CDF (uniform_cdf.go:29).
  - Quantile 2·FMA(u, b/2 − a/2, a/2) (uniform_quantile.go:28). (0.5, −1e308, 1e308) → 0.
  - UniformMeanVariance keeps ErrOverflow because the variance itself is not representable. Its doc and message
    say so.
- DiscreteUniform*: the message already said "int range". The real defect was that the limit rejected
  representable answers. Now discrete_uniform_pmf.go:9 returns the exact uint64 span and m = b − a + 1 as a
  float64 (up to 2⁶⁴), so every int range is accepted and ErrOverflow is never returned.
  - Mean: now a + ⌊d/2⌋ (+½) in integers (discrete_uniform_mean_variance.go:23). The old a/2 + b/2 gave 0
    instead of −½ on MinInt…MaxInt.
  - Quantile: clamps in uint64 (discrete_uniform_quantile.go:35). Its doc states the result works at float64-CDF
    resolution above 2⁵³ points.
- Tests: `TestUniform{PDF,CDF,Quantile}_wideInterval` (uniform_pdf_test.go:27, uniform_cdf_test.go:25,
  uniform_quantile_test.go:26), a new MeanVariance error case, and `TestDiscreteUniform{PMF,CDF,MeanVariance,
  Quantile}_fullIntRange` (pmf_test:43, cdf_test:45, mean_variance_test:35, quantile_test:57). They replace the
  four old ErrOverflow assertions.

## Benchmarks (out/bench_prob_fix.txt; baseline out/bench_prob.txt), all 0 B/op, 0 allocs/op

DiscreteVariance/1e3 30168 ns (was 29104) · /1e5 3034869 ns (was 2875847) · UniformPDF 6.04 ns (6.47) ·
UniformCDF 5.87 (5.87) · UniformQuantile 5.94 (5.95) · UniformMeanVariance 4.53 (4.77) · DiscreteUniformPMF
2.62 (2.76) · CDF 2.73 (2.89) · MeanVariance 2.82 (2.95) · Quantile 5.14 (5.34).

## Changed Source lines (verbatim)

    // Source:  Wheatcroft, "Forecasting football matches by predicting match statistics", Journal of Sports Analytics, 2021, 7(2):77–97, doi:10.3233/JSA-200462; author preprint arXiv:2001.09097v1 (read 2026-11-12), §2.1 "Betting odds" (odds-implied probability rᵢ = 1/Oᵢ, so fair odds are 1/p)   [014]
    // Source:  Wunderlich, Memmert, "The Betting Odds Rating System: Using soccer forecasts to forecast soccer", PLOS ONE, 2018, 13(6):e0198668, doi:10.1371/journal.pone.0198668 (read 2026-11-12), Method, "Transferring betting odds to probabilities" (with no bookmaker margin the inverse odds are the probability)   [014]
    // Source:  Wheatcroft, "Forecasting football matches by predicting match statistics", Journal of Sports Analytics, 2021, 7(2):77–97, doi:10.3233/JSA-200462; author preprint arXiv:2001.09097v1 (read 2026-11-12), §2.1 "Betting odds" (odds-implied probability rᵢ = 1/Oᵢ)   [015]
    // Source:  Egidi, Pauli, Torelli, "Combining historical data and bookmakers' odds in modelling football scores", Statistical Modelling, 2018, 18(5–6):436–459, doi:10.1177/1471082X18798414; author preprint arXiv:1802.08848 (read 2026-11-12), §2 "Transforming the betting odds into probabilities" (the inverse odd is the unfair probability of the event)   [015]
    // Source:  Wheatcroft, "Forecasting football matches by predicting match statistics", Journal of Sports Analytics, 2021, 7(2):77–97, doi:10.3233/JSA-200462; author preprint arXiv:2001.09097v1 (read 2026-11-12), §2.1 "Betting odds", eq. (1) (overround π = Σᵢ 1/Oᵢ − 1)   [016]
    // Source:  Egidi, Pauli, Torelli, "Combining historical data and bookmakers' odds in modelling football scores", Statistical Modelling, 2018, 18(5–6):436–459, doi:10.1177/1471082X18798414; author preprint arXiv:1802.08848 (read 2026-11-12), §2 "Transforming the betting odds into probabilities", eq. (2.1) (booksum β = Σᵢ of the inverse odds)   [016]
    // Source:  Egidi, Pauli, Torelli, "Combining historical data and bookmakers' odds in modelling football scores", Statistical Modelling, 2018, 18(5–6):436–459, doi:10.1177/1471082X18798414; author preprint arXiv:1802.08848 (read 2026-11-12), §2 "Transforming the betting odds into probabilities", eq. (2.1) (basic normalisation πᵢ = oᵢ/β, oᵢ the inverse odds, β the booksum)   [017]
    // Source:  Wunderlich, Memmert, "The Betting Odds Rating System: Using soccer forecasts to forecast soccer", PLOS ONE, 2018, 13(6):e0198668, doi:10.1371/journal.pone.0198668 (read 2026-11-12), Method, "Transferring betting odds to probabilities" (basic normalisation: each inverse odd divided by their sum)   [017]
    // Source:  Chan, Golub, LeVeque, "Algorithms for Computing the Sample Variance: Analysis and Recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115, eq. (1.7) (the corrected two-pass algorithm)   [019, added]
    // Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §1.2 "Basics of Probability Theory"   [002, 005, 007]
    // Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §1.3 "Conditional Probability and Independence"   [003, 006]
    // Source:  Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §1.3 "Conditional Probability and Independence" (Bayes' Rule)   [008, 009]
    Removed: the Štrumbelj 2014 line (014–017); the Altman–Bland 1994 line (010, 011).

## For the lead (shared files, not edited here)

- INDEX.md: the Source cells of 002, 003, 005–011, 014–017 and 019 change to the lines above.
- doc/sources/: retire `strumbelj_2014.md` and `altman_bland_1994.md` as sources (or mark them "not read, not
  cited"); add `wheatcroft_2021.md`, `egidi_2018.md`, `wunderlich_memmert_2018.md` and `chan_golub_leveque_1983.md`;
  change `casella_berger_statistical_inference.md` to section level.
- No ID was held or removed. doc.go is unchanged.
- Question for Devin (sources disagree on a name): Wunderlich–Memmert call the booksum itself the
  "overround" (an "average bookmaker overround of 1.064"), but Wheatcroft eq. (1) defines the overround as
  Σ 1/Oᵢ − 1. Overround (016) follows Wheatcroft and its doc says "booksum minus 1", so Wunderlich is not cited
  for 016. No change unless Devin prefers the booksum convention, which would be a new ID.
