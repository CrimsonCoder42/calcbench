# notes_stats: num/stats builder hand-off (2026-11-12)

Package `num/stats` (SHORT `STATS`): 29 functions built (NUM-STATS-001…006, 008…026, 028…031). Held back:
NUM-STATS-007 GroupedMode (no confirmed source; file moved to `_delete/num/stats/grouped_mode.go`). Not built:
NUM-STATS-027 WeightedGeometricMean (D-010: out; no confirmed locator found). Spec: `doc/specs/s8_stats.md`
(Status BUILDING). Checks run: `gofmt -l num/stats` (clean), `go build ./...` (clean), `go vet ./num/stats/`
(clean). Tests and benchmarks NOT run (lead verifies). Escape analysis (`go build -gcflags=-m`): heap only on
error paths, plus the documented copies/results (Median, MAD, Percentile, Mode, Summary, FrequencyDistribution,
RelativeFrequency).

Files: `num/stats/` doc.go · compensated_sum.go (unexported shared core: Neumaier sum, finiteness checks, meanOf,
centredSquares, centredProducts, midpoint, lerp, bitLen) · arithmetic_mean.go (all sentinel errors) · one file
per function (snake_case of its name) · a `_test.go` beside each. StRD fixtures live in arithmetic_mean_test.go.
Dependencies: `algo/selection` (Median, KthSmallest), `algo/search` (UpperBound); standard library otherwise.
No local helper was needed in place of an algo function; Percentile gets X(j+1) with one O(n) pass after
KthSmallest instead of a new ALG id.

## (a) INDEX.md rows (final form)

| ID | Function | Package | Does (one line) | Time | Space | Source | Tags |
|---|---|---|---|---|---|---|---|
| NUM-STATS-001 | ArithmeticMean | num/stats | sum of the values divided by their count | O(n)   n = len(values); one pass (two only when the plain sum overflows) | O(1) extra, no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm (read 2026-11-12); Neumaier, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 | mean, average, location, central-tendency, descriptive, compensated-summation |
| NUM-STATS-002 | RelativeFrequency | num/stats | each class count divided by the total count | O(b)   b = len(counts); two passes | O(b) for the result (one allocation); O(1) extra | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.3.14 "Histogram", https://www.itl.nist.gov/div898/handbook/eda/section3/histogra.htm (read 2026-11-12) | relative-frequency, frequency, proportion, histogram, descriptive |
| NUM-STATS-003 | CoefficientOfVariation | num/stats | sample standard deviation divided by the mean | O(n)   n = len(values); three passes (mean, then the centred sums) | O(1) extra, no allocation | NIST, Dataplot Reference Manual, "COEFFICIENT OF VARIATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/coefvari.htm (read 2026-11-12) | coefficient-of-variation, cv, relative-dispersion, dispersion, descriptive |
| NUM-STATS-004 | CohenD | num/stats | difference of two means in standard deviations (Cohen's d) | O(1) | O(1), no allocation | Cohen, Statistical Power Analysis for the Behavioral Sciences, 2nd ed., Lawrence Erlbaum, 1988, §2.2, eq. 2.2.1, p. 20; Lakens, "Calculating and reporting effect sizes to facilitate cumulative science: a practical primer for t-tests and ANOVAs", Frontiers in Psychology, 2013, 4:863, doi:10.3389/fpsyg.2013.00863 (eq. 1, d_s with the pooled standard deviation); Goulet-Pelletier, Cousineau, "A review of effect sizes and their confidence intervals, Part I: The Cohen's d family", The Quantitative Methods for Psychology, 2018, 14(4):242–265, doi:10.20982/tqmp.14.4.p242 (eq. 1) | cohen-d, effect-size, standardized-mean-difference, two-sample, descriptive |
| NUM-STATS-005 | GroupedMean | num/stats | estimated mean of grouped data from class marks | O(b)   b = len(counts); two passes | O(1) extra, no allocation | Illowsky, Dean, Introductory Statistics, 2nd ed., OpenStax, 2023, §2.5 "Measures of the Center of the Data", "Calculating the Mean of Grouped Frequency Tables", https://openstax.org/books/introductory-statistics-2e/pages/2-5-measures-of-the-center-of-the-data (read 2026-11-12) | grouped-data, mean, class-mark, frequency, location, central-tendency, descriptive |
| NUM-STATS-006 | GroupedMedian | num/stats | estimated median of grouped data by interpolation | O(b)   b = len(counts); two passes (total, then cumulative) | O(1) extra, no allocation | NIST, Dataplot Reference Manual, "GROUPED PERCENTILE", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/groupper.htm (read 2026-11-12) | grouped-data, median, interpolation, frequency, location, central-tendency, descriptive |
| NUM-STATS-008 | FrequencyDistribution | num/stats | count of values in each class | O(n log b)   n = len(values), b = len(edges) − 1 classes; one binary search per value (ALG-SEARCH-002) | O(b) for the result (one allocation); O(1) extra | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.3.14 "Histogram", https://www.itl.nist.gov/div898/handbook/eda/section3/histogra.htm (read 2026-11-12) | frequency-distribution, histogram, counts, classes, bins, descriptive |
| NUM-STATS-009 | LinearRegression | num/stats | intercept and slope of the least-squares line | O(n)   n = len(x); three passes (means, then the centred sums) | O(1) extra, no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §4.4.3.1 "Least Squares", https://www.itl.nist.gov/div898/handbook/pmd/section4/pmd431.htm (read 2026-11-12) | linear-regression, least-squares, trend-line, slope, intercept, fit, descriptive |
| NUM-STATS-010 | MatthewsCorrelation | num/stats | Matthews correlation coefficient of a 2×2 table | O(1) | O(1), no allocation | Matthews, "Comparison of the predicted and observed secondary structure of T4 phage lysozyme", Biochimica et Biophysica Acta (BBA) – Protein Structure, 1975, 405(2):442–451, doi:10.1016/0005-2795(75)90109-9; Baldi, Brunak, Chauvin, Andersen, Nielsen, "Assessing the accuracy of prediction algorithms for classification: an overview", Bioinformatics, 2000, 16(5):412–424, doi:10.1093/bioinformatics/16.5.412 (eq. 11; undefined when a margin is 0, p. 415); Chicco, Jurman, "The advantages of the Matthews correlation coefficient (MCC) over F1 score and accuracy in binary classification evaluation", BMC Genomics, 2020, 21:6, doi:10.1186/s12864-019-6413-7 (eq. 2) | matthews-correlation, mcc, phi-coefficient, confusion-matrix, classification, association |
| NUM-STATS-011 | MeanAbsoluteDeviation | num/stats | average absolute deviation about the mean | O(n)   n = len(values); two passes (mean, then deviations) | O(1) extra, no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.6 "Measures of Scale", https://www.itl.nist.gov/div898/handbook/eda/section3/eda356.htm (read 2026-11-12) | mean-absolute-deviation, average-absolute-deviation, aad, dispersion, scale, descriptive |
| NUM-STATS-012 | Median | num/stats | middle value, or mean of the two middle values | O(n) worst case   n = len(values) | O(n) extra (one allocation, the copy) plus O(log n) stack | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm (read 2026-11-12) | median, location, central-tendency, robust, order-statistic, descriptive |
| NUM-STATS-013 | Midrange | num/stats | mean of the smallest and largest values | O(n)   n = len(values); one pass | O(1) extra, no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm (read 2026-11-12) | midrange, mid-range, location, central-tendency, descriptive |
| NUM-STATS-014 | PearsonCorrelation | num/stats | Pearson correlation coefficient of paired samples | O(n)   n = len(x); three passes (means, then the centred sums) | O(1) extra, no allocation | NIST, Dataplot Reference Manual, "CORRELATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/correlat.htm (read 2026-11-12) | correlation, pearson, product-moment, association, descriptive |
| NUM-STATS-015 | Percentile | num/stats | p-quantile by a chosen Hyndman–Fan method | O(n) worst case   n = len(values); one linear-time selection and one pass | O(n) extra (one allocation, the selection's copy) plus O(log n) stack | Hyndman, Fan, "Sample quantiles in statistical packages", The American Statistician, 1996, 50(4):361–365, doi:10.1080/00031305.1996.10473566 (eq. 1; definitions 1–3 §2, 4–9 §3 and Table 3); NIST/SEMATECH, e-Handbook of Statistical Methods, §7.2.6.2 "Percentiles", https://www.itl.nist.gov/div898/handbook/prc/section2/prc262.htm (read 2026-11-12) | percentile, quantile, order-statistic, hyndman-fan, descriptive |
| NUM-STATS-016 | PooledStandardDeviation | num/stats | pooled standard deviation of two samples | O(1) | O(1), no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.3 "Two-Sample t-Test for Equal Means", https://www.itl.nist.gov/div898/handbook/eda/section3/eda353.htm (read 2026-11-12) | pooled-standard-deviation, standard-deviation, two-sample, dispersion, descriptive |
| NUM-STATS-017 | PopulationVariance | num/stats | variance with divisor N | O(n)   n = len(values); three passes (mean, then the centred sums) | O(1) extra, no allocation | Hogg, McKean, Craig, Introduction to Mathematical Statistics, 8th ed., Pearson, 2019, §1.9 "Some Special Expectations", Definition 1.9.2 (variance σ² = E[(X − μ)²]); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §2.3 "Moments and Moment Generating Functions" (variance); Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the corrected two-pass algorithm) | variance, population-variance, dispersion, scale, descriptive, two-pass |
| NUM-STATS-018 | Range | num/stats | largest value minus smallest value | O(n)   n = len(values); one pass | O(1) extra, no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.6 "Measures of Scale", https://www.itl.nist.gov/div898/handbook/eda/section3/eda356.htm (read 2026-11-12) | range, spread, dispersion, scale, descriptive, min, max |
| NUM-STATS-019 | SampleStandardDeviation | num/stats | standard deviation with divisor n − 1 | O(n)   n = len(values); three passes (mean, then the centred sums) | O(1) extra, no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.6 "Measures of Scale", https://www.itl.nist.gov/div898/handbook/eda/section3/eda356.htm (read 2026-11-12); Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the corrected two-pass algorithm) | standard-deviation, sample-standard-deviation, dispersion, scale, descriptive |
| NUM-STATS-020 | SummaryStatistics | num/stats | count, mean, sample SD, min, max, median and range | O(n) worst case   n = len(values); four passes and one linear-time selection | O(n) extra (one allocation, the median's copy) plus O(log n) stack | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm (read 2026-11-12); NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.6 "Measures of Scale", https://www.itl.nist.gov/div898/handbook/eda/section3/eda356.htm (read 2026-11-12) | summary-statistics, descriptive, mean, standard-deviation, median, range, min, max |
| NUM-STATS-021 | SampleVariance | num/stats | variance with divisor n − 1 | O(n)   n = len(values); three passes (mean, then the centred sums) | O(1) extra, no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.6 "Measures of Scale", https://www.itl.nist.gov/div898/handbook/eda/section3/eda356.htm (read 2026-11-12); Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the corrected two-pass algorithm) | variance, sample-variance, dispersion, scale, descriptive, two-pass |
| NUM-STATS-022 | PopulationStandardDeviation | num/stats | standard deviation with divisor N | O(n)   n = len(values); three passes (mean, then the centred sums) | O(1) extra, no allocation | Hogg, McKean, Craig, Introduction to Mathematical Statistics, 8th ed., Pearson, 2019, §1.9 "Some Special Expectations", Definition 1.9.2 and p. 69 (standard deviation σ, the positive square root of the variance); Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §2.3 "Moments and Moment Generating Functions" (variance and standard deviation); Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the corrected two-pass algorithm) | standard-deviation, population-standard-deviation, dispersion, scale, descriptive |
| NUM-STATS-023 | StandardErrorOfMean | num/stats | sample standard deviation divided by √n | O(n)   n = len(values); three passes (mean, then the centred sums) | O(1) extra, no allocation | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.2 "Confidence Limits for the Mean", https://www.itl.nist.gov/div898/handbook/eda/section3/eda352.htm (read 2026-11-12); NIST, Dataplot Reference Manual, "STANDARD DEVIATION OF THE MEAN", https://www.itl.nist.gov/div898/software/dataplot/refman2/ch2/sd_mean.pdf (read 2026-11-12) | standard-error, sem, mean, precision, descriptive |
| NUM-STATS-024 | Mode | num/stats | most frequent values, all ties, ascending | O(n log n) worst case   n = len(values); one sort of the copy, then two linear scans | O(n) extra (the copy) plus O(k) for the result, k = number of modes; two allocations | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm (read 2026-11-12) | mode, location, central-tendency, frequency, descriptive |
| NUM-STATS-025 | SumOfSquares | num/stats | sum of squared deviations about the mean | O(n)   n = len(values); two passes (four only when a plain sum overflows) | O(1) extra, no allocation | NIST, Dataplot Reference Manual, "SUM OF SQUARES FROM MEAN", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/ssqmean.htm (read 2026-11-12); Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115 (the corrected two-pass algorithm) | sum-of-squares, deviation, dispersion, variance, descriptive, two-pass |
| NUM-STATS-026 | WeightedMean | num/stats | weighted mean Σw·x / Σw | O(n)   n = len(values); two passes (four only when a plain sum overflows) | O(1) extra, no allocation | NIST, Dataplot Reference Manual, "WEIGHTED MEAN", https://www.itl.nist.gov/div898/software/dataplot/refman2/ch2/weigmean.pdf (read 2026-11-12); NIST, Dataplot Reference Manual, "WEIGHTED TRIMMED MEAN", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/weigtmea.htm (read 2026-11-12) (weights non-negative, at least one positive) | weighted-mean, weighted-average, mean, location, central-tendency, descriptive |
| NUM-STATS-028 | ZScore | num/stats | standard score (x − mean)/sd | O(1) | O(1), no allocation | NIST, Dataplot Reference Manual, "ISO 13528 ZSCORE", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/zscore.htm (read 2026-11-12); ISO 13528:2015, Statistical methods for use in proficiency testing by interlaboratory comparison, 2nd ed., ISO, §9.4 "z scores" | z-score, standard-score, standardize, normalization, descriptive |
| NUM-STATS-029 | RawScoreFromZScore | num/stats | value with a given standard score, mean + z·sd | O(1) | O(1), no allocation | NIST, Dataplot Reference Manual, "ISO 13528 ZSCORE", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/zscore.htm (read 2026-11-12) (Z(i) = (X(i) − Xref)/sigma, solved for X(i)); ISO 13528:2015, Statistical methods for use in proficiency testing by interlaboratory comparison, 2nd ed., ISO, §9.4 "z scores" | z-score, raw-score, standard-score, inverse, descriptive |
| NUM-STATS-030 | PercentRelativeStandardDeviation | num/stats | 100·s/abs(mean), the percent relative standard deviation | O(n)   n = len(values); three passes (mean, then the centred sums) | O(1) extra, no allocation | NIST, Dataplot Reference Manual, "RELATIVE STANDARD DEVIATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/ch2/relat_sd.pdf (read 2026-11-12); NIST, Dataplot Reference Manual, "COEFFICIENT OF VARIATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/coefvari.htm (read 2026-11-12) (relsd, the percent relative standard deviation) | relative-standard-deviation, rsd, percent-rsd, coefficient-of-variation, dispersion, descriptive |
| NUM-STATS-031 | MedianAbsoluteDeviation | num/stats | median absolute deviation about the median | O(n) worst case   n = len(values); two linear-time selections | O(n) extra (two allocations: the deviations and the selection's copy) plus O(log n) stack | NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.6 "Measures of Scale", https://www.itl.nist.gov/div898/handbook/eda/section3/eda356.htm (read 2026-11-12); NIST, Dataplot Reference Manual, "MEDIAN ABSOLUTE DEVIATION", https://www.itl.nist.gov/div898/software/dataplot/refman2/auxillar/mad.htm (read 2026-11-12) | median-absolute-deviation, mad, robust, dispersion, scale, descriptive |

Rows for unused IDs (if the index keeps reserved rows): `| NUM-STATS-007 | GroupedMode | num/stats | NOT BUILT:
held back, no confirmed source | | | | |` and `| NUM-STATS-027 | WeightedGeometricMean | num/stats | NOT BUILT: out
(D-010), no source | | | | |`.

## (b) doc/sources/ records

- **nist_ehandbook.md** (exists; add): §1.3.5.1 Measures of Location (eda351.htm; mean, median, mode, mid-range
  and its "not robust" note) → 001, 012, 013, 020, 024; §1.3.5.6 Measures of Scale (eda356.htm; s² with N−1,
  range, AAD = Σ|Y−Ȳ|/N, MAD = median|Y−Ỹ| unscaled) → 011, 018, 019, 020, 021, 031; §1.3.3.14 Histogram
  (histogra.htm; classes, normalized count) → 002, 008; §4.4.3.1 Least Squares (pmd431.htm; β̂1, β̂0) → 009;
  §7.2.6.2 Percentiles (prc262.htm; p(N+1) = k + d, 12 resistivity values, Y(90) 95.1981 / R7 95.1957 / R8
  95.1972) → 015, and the resistivity data used by 012, 013, 018, 011, 031, 025 tests; §1.3.5.3 Two-Sample t-Test
  (eda353.htm; s_p, AUTO83B N1 249 s1 6.41470 N2 79 s2 6.10771 s_p 6.34260) → 016, 004 test; §1.3.5.2 (eda352.htm;
  s/√N inside the CI, ZARR13 N 195 s 0.022789, interval 9.258242–9.264679) → 023. Checked: each page fetched
  2026-11-12; examples reproduced in the `_worked`/`_eHandbook` tests.
- **nist_strd.md** (new): NIST Statistical Reference Datasets, https://www.itl.nist.gov/div898/strd/ (read
  2026-11-12). Univariate NumAcc1–4, Michelso, Mavro, Lew (certified mean and sample SD) → 001, 003, 019, 020, 021,
  023, 025, 030; Linear LS Norris (B0 −0.262323073774029, B1 1.00211681802045, R² 0.999993745883712) → 009, 014.
  Checked: data transcribed (raw download was blocked; NumAcc2–4 built from the page's stated pattern) and
  verified in exact rationals against every certified value; embedded as Go literals; PiDigits NOT used (the
  transcription could not reproduce NIST's certified values).
- **nist_dataplot.md** (new): NIST Dataplot Reference Manual, https://www.itl.nist.gov/div898/software/dataplot/
  refman2/ (read 2026-11-12). auxillar/coefvari.htm (cv = s/x̄; relsd = 100(s/|x̄|)) → 003, 030; ch2/relat_sd.pdf
  (relsd = 100·(sd/|mean|), cites Snedecor & Cochran 6th ed. pp. 62–64, not checked) → 030; auxillar/groupper.htm
  (p_i = l + ((iN/100) − F)/f·C, uniform-within-class note; cites Selby, CRC Standard Mathematical Tables, 22nd
  ed., 1974, p. 572, not checked) → 006; auxillar/correlat.htm (r = S_xy/(√S_xx·√S_yy)) → 014; auxillar/ssqmean.htm
  (SSQ = Σ(X−X̄)²) → 025; ch2/weigmean.pdf (x̄w = Σw x/Σw) and auxillar/weigtmea.htm ("none of the weights can be
  negative", at least one positive) → 026; auxillar/zscore.htm (Z = (X − Xref)/sigma; cites ISO 13528:2005 pp.
  25–26) → 028, 029; auxillar/mad.htm (MAD = median|x − x̃|, MADN = MAD/0.6745) → 031; ch2/sd_mean.pdf (s/√N) → 023.
  Checked: pages fetched; formulas compared; worked examples in tests.
- **hyndman_fan_1996.md** (new): Hyndman, Fan, "Sample quantiles in statistical packages", The American
  Statistician, 1996, 50(4):361–365, doi:10.1080/00031305.1996.10473566. Taken: eq. (1), definitions 1–3 (§2),
  4–9 (§3, Table 3); type 8 recommended (§5). → 015. Checked: Crossref; author's PDF; R ?quantile agrees;
  TestPercentile_allTypes (table computed in exact rationals) and _eHandbook.
- **chan_golub_leveque_1983.md** (new): Chan, Golub, LeVeque, "Algorithms for computing the sample variance:
  analysis and recommendations", The American Statistician, 1983, 37(3):242–247, doi:10.1080/00031305.1983.10483115.
  Taken: the corrected two-pass algorithm Σd² − (Σd)²/n. → 017, 019, 021, 022, 025 (and internally 003, 020, 023,
  030, 009, 014). Checked: Crossref; StRD NumAcc and oracle tests vs 512-bit big.Float.
- **neumaier_1974.md** (exists; add IDs): compensated summation → 001 (and every compensated sum in num/stats).
- **hogg_mckean_craig.md** (new): Hogg, McKean, Craig, Introduction to Mathematical Statistics, 8th ed., Pearson,
  2019, §1.9 "Some Special Expectations", Definition 1.9.2 (variance E[(X−μ)²]), p. 69 (σ, the positive square root).
  → 017, 022. Checked: text read from a copy of the 8th ed. (pp. 68–69); hand-checked examples in tests.
- **casella_berger_2002.md** (new): Casella, Berger, Statistical Inference, 2nd ed., Duxbury, 2002, §2.3 "Moments
  and Moment Generating Functions". → 017, 022. Checked: section title via table of contents only.
- **openstax_introductory_statistics_2e.md** (new): Illowsky, Dean, Introductory Statistics, 2nd ed., OpenStax,
  2023, §2.5, "Calculating the Mean of Grouped Frequency Tables": mean = Σfm/Σf, m the midpoint. → 005. Checked:
  page fetched; TestGroupedMean_worked and _oracle.
- **cohen_1988.md** (new): Cohen, Statistical Power Analysis for the Behavioral Sciences, 2nd ed., Lawrence Erlbaum,
  1988, §2.2, eq. 2.2.1, p. 20: d = (m_A − m_B)/σ. Secondaries: Lakens 2013, Front. Psychol. 4:863,
  doi:10.3389/fpsyg.2013.00863, eq. 1; Goulet-Pelletier, Cousineau 2018, TQMP 14(4):242–265,
  doi:10.20982/tqmp.14.4.p242, eq. 1. → 004. Checked: secondaries read; the eq./page locator from Devin (D-010)
  and a GSU memo (Bakeman 2021), not from the book.
- **matthews_1975.md** (new): Matthews, BBA – Protein Structure, 1975, 405(2):442–451,
  doi:10.1016/0005-2795(75)90109-9 (Crossref). Secondaries: Baldi, Brunak, Chauvin, Andersen, Nielsen 2000,
  Bioinformatics 16(5):412–424, doi:10.1093/bioinformatics/16.5.412, eq. (11), p. 415 (undefined at a zero
  margin); Chicco, Jurman 2020, BMC Genomics 21:6, doi:10.1186/s12864-019-6413-7, eq. 2. → 010. Checked: Crossref
  for all three; secondaries read; TestMatthewsCorrelation_oracle vs 256-bit big.Float.
- **iso_13528.md** (new): ISO 13528:2015, Statistical methods for use in proficiency testing by interlaboratory
  comparison, 2nd ed., §9.4 "z scores" (TOC confirmed via DIN ISO 13528:2020 and the 2022 edition's §9.4 p. 26);
  2005 1st ed. pp. 25–26 per Dataplot. → 028, 029 (second source). Checked: TOC only; formula via Dataplot.

## (c) Citations

Confirmed 2026-11-12: all e-Handbook and Dataplot quotes above (fetched; the catalog's [WF] quotes match, except
"WEIGHTED MEAN": weigtmea.htm is the weighted TRIMMED mean page; the plain weighted mean is ch2/weigmean.pdf,
now cited); Crossref: Hyndman–Fan, Chan–Golub–LeVeque, Matthews, Baldi et al., Chicco–Jurman, Welford (not
cited: no streaming function), Hampel 1974 (not cited, see below); Hogg–McKean–Craig 8th ed. §1.9 read in text.
Changed from Devin's list: Hogg–McKean–Craig cited as the 8th ed. (2019), where §1.9 / Def. 1.9.2 was read; the
7th-ed. numbering could not be confirmed (believed identical).
NOT confirmed, and therefore dropped: Spiegel & Stephens, Schaum's Outline of Statistics, 4th ed., McGraw-Hill,
2008, ISBN 978-0-07-148584-5 (edition and ch. 3 title confirmed; the grouped formulas and equation numbers could
not be read); Hampel 1974 as defining the MAD; "eq. 1.7" of Chan–Golub–LeVeque; Knuth TAOCP §4.2.2 (not used);
Bullen, Handbook of Means (2003), the weighted geometric mean's page in ch. II; Casella–Berger Def. 2.3.2 (kept at
§2.3 level only); Cohen's own page (via Devin + memo only).
Held back: **NUM-STATS-007 GroupedMode**: Spiegel unread, and no scholarly secondary reproduces L + Δ1/(Δ1+Δ2)·c
(only revision sites). The code was written (no tests); the file is in `_delete/num/stats/grouped_mode.go`. Releasing it needs one library check of Spiegel ch. 3.
IDs left unused: NUM-STATS-007 (held), NUM-STATS-027 (out).

Proposed ASSUMPTIONS (for the lead to number): (1) FrequencyDistribution: a value outside the edges is
ErrOutOfRange, never dropped or counted. (2) Percentile: a fractional part within 4ε·max(1, |np + m|) of an
integer is that integer (R-style fuzz). (3) MatthewsCorrelation: a zero margin is ErrZeroTotal (Baldi: undefined),
not Chicco–Jurman's limit 0. (4) CohenD takes the common σ as an argument (eq. 2.2.1 as written); a two-sample
convenience form is not built. (5) Mode compares with cmp.Compare (NaN = NaN, −0 = +0).

## (d) FOUND, NOT FIXED

> Superseded by fix_stats.md (2026-11-12): the first three items (MAD allocations, squared deviations below ~1e-154 underflowing, WeightedMean product underflow) are fixed there.

- `algo/selection` exports only cloning selectors, so MedianAbsoluteDeviation makes two O(n) allocations (its
  deviation buffer, then selection's copy of it). An in-place `SelectInPlace` in algo/selection would cut one.
- Limitation, documented nowhere yet: squares of deviations below ~1e-154 underflow, so a sample whose spread is
  that small (e.g. x spread 1e-200 in PearsonCorrelation) reports ErrZeroVariance or a variance of 0. A scale-up
  retry, symmetric to the overflow retry, would fix it; not done.
- WeightedMean: a weight·value product that underflows to 0 while the weight sum does not (weights and values
  both ~1e-300) loses that term; the overflow retry does not cover underflow.
- Catalog `doc/catalog/stats.md`: the NUM-STATS-026 entry names weigtmea.htm as the weighted mean page; it is the
  weighted trimmed mean (the plain one is ch2/weigmean.pdf). The %RSD/MAD IDs 030/031 are not in the catalog.
- PiDigits StRD fixture could not be obtained (download blocked); not used.
