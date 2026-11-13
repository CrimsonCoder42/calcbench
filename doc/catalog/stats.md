# Catalog: statistics

Status: DRAFT for Devin's approval (2026-11-12). Nothing here is built. Work folder: `../catalog_stats_2026-11-12/`.
Input: the names in vCalc's Statistics category, pasted by Devin. vCalc pages were not fetched, and nothing
but the names was taken from it. Formulas, examples and sources come from the primary sources cited below,
read 2026-11-12.

Source status of each KEEP:
- **[E]**: NIST/SEMATECH e-Handbook of Statistical Methods, or a paper confirmed by Crossref.
- **[D]**: NIST Dataplot Reference Manual only. This is NIST's documentation for its own software. Whether it
  counts as acceptable is question Q4.
- **SOURCE NOT FOUND**: no acceptable source was readable, so no formula is given.

A quote marked "[WF]" came through a web summarizer: it is near-verbatim and must be checked character by
character before it goes into a doc comment.

## Step 1: Clean

208 rows: 10 subcategory or collection names, 194 equations and 4 data items, in the order given. The
descriptive paragraph and the section headers are not rows. A collection that bundles other rows is MERGEd
into the row that covers it.

Counts: **KEEP 29 · MOVE 68 · MERGE 57 · DROP 29 · LATER 25.**

| # | Name as given | Decision | Home | Reason |
|---|---|---|---|---|
| | **Subcategories and collections** | | | |
| c1 | Small Sample Statistics | LATER | later inference slice | Small-sample methods are t-based inference |
| c2 | Baseball Statistics | MOVE | `life/sports` | Baseball rows |
| c3 | Basketball Statistics Calculator | MOVE | `life/sports` | Basketball rows |
| c4 | Central Tendacy Calculator (Grouped Data) | MERGE | → #47 Estimated Mean equation | A bundle of the three grouped-data estimates (#47–#49) |
| c5 | Chi-Square Calc | LATER | later inference slice | Chi-square test |
| c6 | Psychology and Stats Calc | MERGE | → #139 Sample Statistics | A bundle of descriptive measures |
| c7 | SMCM Psychology and Statstics Collection | MERGE | → #139 Sample Statistics | A bundle of descriptive measures |
| c8 | Statistical Decision Tree | LATER | later inference slice | Chooses an inferential test; a guide, not a formula |
| c9 | Statistics Calculator | MERGE | → #139 Sample Statistics | A bundle of descriptive measures |
| c10 | Stats Calc | MERGE | → #139 Sample Statistics | A bundle of descriptive measures |
| | **Equations** | | | |
| 1 | Anova | LATER | later inference slice | ANOVA |
| 2 | Arithmetic Mean | KEEP | `num/stats` NUM-STATS-001 | Location |
| 3 | Arithmetic Mean (average) | MERGE | → #2 | Duplicate |
| 4 | average (mean) | MERGE | → #2 | Duplicate |
| 5 | Average number of kernels on five ears | MERGE | → #2 | Fixed n = 5 with a domain label |
| 6 | Basketball Scoring Efficiency | MOVE | `life/sports` | Basketball |
| 7 | Batting Average | MOVE | `life/sports` | Baseball |
| 8 | Batting Average Against | MOVE | `life/sports` | Baseball |
| 9 | Bayes' Theorem for Disease Testing | MOVE | `num/prob` | Probability (a medical wording does not make it `med/`) |
| 10 | Bayes' Theorem for Event Probability | MOVE | `num/prob` | Probability |
| 11 | Big Decimal Round | MOVE | `num/algebra` | Rounding, not statistics |
| 12 | Binomial Coefficient | MOVE | `num/combinatorics` | Combinatorics |
| 13 | Business Stats | MERGE | → #139 Sample Statistics | A bundle of descriptive measures |
| 14 | Chebyshev's Inequality(1) | MOVE | `num/prob` | Probability bound |
| 15 | Chebyshev's Inequality(1) | MERGE | → #14 | Duplicate (second author) |
| 16 | Chebyshev's Inequality(2) | MERGE | → #14 | The same inequality in its other form (k vs ε) |
| 17 | Chebyshev's Inequality(2) | MERGE | → #14 | Duplicate |
| 18 | Chi-Square (Stats) | LATER | later inference slice | Chi-square |
| 19 | Chi-square Degrees of Freedom | LATER | later inference slice | Part of the chi-square test |
| 20 | Chi-square Expected Result | LATER | later inference slice | Part of the chi-square test |
| 21 | Chi-square goodness-of-fit test (C=3) | LATER | later inference slice | Chi-square |
| 22 | Chi-Square Test | LATER | later inference slice | Chi-square |
| 23 | Chi-squared test for association | LATER | later inference slice | Chi-square |
| 24 | Choosing a Non-Parametric Test | LATER | later inference slice | Test-selection guide, not a formula |
| 25 | Choosing a Parametric Test | LATER | later inference slice | Test-selection guide |
| 26 | Choosing a Relationship Test | LATER | later inference slice | Test-selection guide |
| 27 | Circular r-Permutation of n Elements | MOVE | `num/combinatorics` | Permutations |
| 28 | Class Relative Frequency | KEEP | `num/stats` NUM-STATS-002 | Frequency |
| 29 | Coefficient of variation | KEEP | `num/stats` NUM-STATS-003 | Relative dispersion |
| 30 | Coin Flip | DROP | — | Random: not a pure function |
| 31 | Competition Analysis | DROP | — | The name does not say what it computes, and vCalc is not read; question for Devin |
| 32 | Conditional Probability of B given A | MOVE | `num/prob` | Probability |
| 33 | Count | DROP | — | `len` |
| 34 | Defense-Independent Component ERA | MOVE | `life/sports` | Baseball |
| 35 | Descriptive Statistics (n=10) | MERGE | → #139 Sample Statistics | Fixed-n variant |
| 36 | Descriptive Statistics (n=11) | MERGE | → #139 | Fixed-n variant |
| 37 | Descriptive Statistics (n=12) | MERGE | → #139 | Fixed-n variant |
| 38 | Descriptive Statistics (n=6) | MERGE | → #139 | Fixed-n variant |
| 39 | Descriptive Statistics (n=7) | MERGE | → #139 | Fixed-n variant |
| 40 | Descriptive Statistics (n=8) | MERGE | → #139 | Fixed-n variant |
| 41 | Descriptive Statistics (n=9) | MERGE | → #139 | Fixed-n variant |
| 42 | Displaced Gaussian | MOVE | `num/prob` | Distribution (normal pdf with a location) |
| 43 | Earned Run Average | MOVE | `life/sports` | Baseball |
| 44 | Effect Size (Cohen's d) | KEEP | `num/stats` NUM-STATS-004 | Descriptive effect size from two samples |
| 45 | Effect Size (r-squared) | LATER | later inference slice | Computed from a t statistic and its degrees of freedom |
| 46 | Equivalent Average | MOVE | `life/sports` | Baseball |
| 47 | Estimated Mean equation | KEEP | `num/stats` NUM-STATS-005 | Mean of grouped data |
| 48 | Estimated Median equation | KEEP | `num/stats` NUM-STATS-006 | Median of grouped data |
| 49 | Estimated Mode equation | KEEP | `num/stats` NUM-STATS-007 | Mode of grouped data |
| 50 | Expected Value(Bernoulli) | MOVE | `num/prob` | Expected value |
| 51 | Expected Value(Binomial) | MOVE | `num/prob` | Expected value |
| 52 | Expected Value(Discrete Uniform) | MOVE | `num/prob` | Expected value |
| 53 | Expected Value(Geometric) | MOVE | `num/prob` | Expected value |
| 54 | Expected Value(Geometric) | MERGE | → #53 | Duplicate |
| 55 | Expected Value(Hypergeometric) | MOVE | `num/prob` | Expected value |
| 56 | Expected Value(Pascal Distribution) | MOVE | `num/prob` | Expected value |
| 57 | Expected Value(Pascal) | MERGE | → #56 | Duplicate |
| 58 | Factorial Analysis of Variance (ANOVA) - 2x2, n=24 | LATER | later inference slice | ANOVA |
| 59 | Fielding Independent Pitching | MOVE | `life/sports` | Baseball |
| 60 | Fielding Percentage | MOVE | `life/sports` | Baseball |
| 61 | Five Ear Sample Average | MERGE | → #2 | Fixed n = 5 with a domain label |
| 62 | Free Throw Percentage | MOVE | `life/sports` | Basketball |
| 63 | Frequency Distribution | KEEP | `num/stats` NUM-STATS-008 | Counts per class |
| 64 | Frequency Distribution | MERGE | → #63 | Duplicate |
| 65 | Frequency Distribution (10 bin) | MERGE | → #63 | Fixed bin count |
| 66 | Gamma Function | MOVE | `num/combinatorics` | Special function (Γ(n+1) = n!); the home is a question for Devin |
| 67 | Gaussian | MOVE | `num/prob` | Distribution |
| 68 | Gaussian Distribution | MERGE | → #67 | Duplicate |
| 69 | Gaussian Distribution - Displaced | MERGE | → #42 | Duplicate |
| 70 | Grades Card Frequency Distribution | MERGE | → #63 | Fixed classes with a domain label |
| 71 | Last-Squares Trend Line | KEEP | `num/stats` NUM-STATS-009 | Least-squares line (descriptive fit) |
| 72 | Law of Total Probability P(A+B) | MOVE | `num/prob` | Probability |
| 73 | Least-squares Trend Inference | LATER | later inference slice | Inference on the slope |
| 74 | Lognormal | MOVE | `num/prob` | Distribution |
| 75 | Lognormal Distribution | MERGE | → #74 | Duplicate |
| 76 | Mann-Whitney U test | LATER | later inference slice | Mann-Whitney |
| 77 | Matthews Correlation Coefficient | KEEP | `num/stats` NUM-STATS-010 | Association in a 2×2 table |
| 78 | Max | DROP | — | `slices.Max` |
| 79 | Maximum | DROP | — | `slices.Max` |
| 80 | Mean | MERGE | → #2 | Duplicate |
| 81 | Mean Absolute Deviation | KEEP | `num/stats` NUM-STATS-011 | Dispersion |
| 82 | Median | KEEP | `num/stats` NUM-STATS-012 | Location |
| 83 | Median | MERGE | → #82 | Duplicate |
| 84 | Mid Point | KEEP | `num/stats` NUM-STATS-013 | Read as the midrange (max+min)/2; question for Devin |
| 85 | Min | DROP | — | `slices.Min` |
| 86 | Minimum | DROP | — | `slices.Min` |
| 87 | Multiple Values from Gaussian Distribution | DROP | — | Random draws: not a pure function |
| 88 | Multiple Values from Gaussian Distribution (as string) | DROP | — | Random draws |
| 89 | Numeric Sort | DROP | — | `slices.Sort` |
| 90 | Observational Statistics | MERGE | → #139 Sample Statistics | A bundle of descriptive measures |
| 91 | On-base Percentage | MOVE | `life/sports` | Baseball |
| 92 | On-base Plus Slugging | MOVE | `life/sports` | Baseball |
| 93 | One Way ANOVA | LATER | later inference slice | ANOVA |
| 94 | One-way Analysis of Variance (ANOVA) - 3 Groups | LATER | later inference slice | ANOVA |
| 95 | P(event) = F / T | MOVE | `num/prob` | Probability |
| 96 | P(n,S) - Bionomial Probability | MOVE | `num/prob` | Distribution |
| 97 | P(y) from z | MOVE | `num/prob` | Normal CDF |
| 98 | P(y) Left and Right of z | MOVE | `num/prob` | Normal CDF |
| 99 | P(y) right of z | MOVE | `num/prob` | Normal CDF |
| 100 | Paired Samples t-test | LATER | later inference slice | t-test |
| 101 | Paired Samples t-test (n=6) | LATER | later inference slice | t-test |
| 102 | Pearson correlation | KEEP | `num/stats` NUM-STATS-014 | Correlation (the coefficient, not its test) |
| 103 | Percent Relative Standard Deviation | MERGE | → #29 | CV × 100 |
| 104 | Percentile | KEEP | `num/stats` NUM-STATS-015 | Quantile |
| 105 | Performance Index Rating | MOVE | `life/sports` | Basketball |
| 106 | Pooled Standard Deviation | KEEP | `num/stats` NUM-STATS-016 | Dispersion of two samples |
| 107 | Population Variance | KEEP | `num/stats` NUM-STATS-017 | Dispersion |
| 108 | Probability between raw scores | MOVE | `num/prob` | Normal CDF |
| 109 | Probability between z SCORES | MOVE | `num/prob` | Normal CDF |
| 110 | Probability of A and B | MOVE | `num/prob` | Probability |
| 111 | Probability of A or B | MOVE | `num/prob` | Probability |
| 112 | Probability of Non-independent Events (A, B) | MOVE | `num/prob` | Probability |
| 113 | Pythagorean Expectation | MOVE | `life/sports` | Baseball |
| 114 | Random 0 OR 1 | DROP | — | Random |
| 115 | Random Integer | DROP | — | Random |
| 116 | Random Integer [0 , Upper] | DROP | — | Random |
| 117 | Random Real Number in Range | DROP | — | Random |
| 118 | Random Real [0.0, 1.0] | DROP | — | Random |
| 119 | Random Real [0.0, Upper] | DROP | — | Random |
| 120 | Random Sample (k items) | DROP | — | Random |
| 121 | Range | KEEP | `num/stats` NUM-STATS-018 | Dispersion (max − min) |
| 122 | Range Factor | MOVE | `life/sports` | Baseball |
| 123 | Rank Order Centroid | MOVE | home not in taxonomy | Decision-analysis weights, not statistics; question for Devin |
| 124 | Regression Line (n=6) | MERGE | → #71 | Fixed-n variant |
| 125 | regression test -- data set access (Weighted Geometric Mean) | MERGE | → #182 | A vCalc test page for the same formula |
| 126 | Relative Standard Deviation | MERGE | → #29 | Same quantity as CV |
| 127 | Repeated-measures Analysis of Variance (ANOVA) - 3 Groups | LATER | later inference slice | ANOVA |
| 128 | Roll a 100-sided Virtual Dice | DROP | — | Random |
| 129 | Roll a 12-sided Virtual Dice | DROP | — | Random |
| 130 | Roll a 20-sided Virtual Dice | DROP | — | Random |
| 131 | Roll a 4-sided Virtual Dice | DROP | — | Random |
| 132 | Roll a 50-sided Virtual Dice | DROP | — | Random |
| 133 | Roll a 6-sided Virtual Dice | DROP | — | Random |
| 134 | Roll a 8-sided Virtual Dice | DROP | — | Random |
| 135 | Roll a D-100 | DROP | — | Random |
| 136 | Roll Five Virtual Cube Dice | DROP | — | Random |
| 137 | Roll Four Virtual Sided Dice | DROP | — | Random |
| 138 | Sample Standard Deviation | KEEP | `num/stats` NUM-STATS-019 | Dispersion |
| 139 | Sample Statistics | KEEP | `num/stats` NUM-STATS-020 | One call for the usual summary set |
| 140 | Sample Variance | KEEP | `num/stats` NUM-STATS-021 | Dispersion |
| 141 | Simple Corn Yield | MOVE | `life/garden` | Agriculture yield estimate |
| 142 | Simple Stats | MERGE | → #139 | A bundle of descriptive measures |
| 143 | Single-Sample t-test | LATER | later inference slice | t-test |
| 144 | Slugging Percentage | MOVE | `life/sports` | Baseball |
| 145 | Spearman's rho correlation | LATER | later inference slice | Spearman |
| 146 | Standard Deviation | KEEP | `num/stats` NUM-STATS-022 | Taken as the population SD (the sample SD has its own rows) |
| 147 | Standard Deviation (Sample) | MERGE | → #138 | Duplicate |
| 148 | Standard Deviation (SD) | MERGE | → #146 | Same unqualified name |
| 149 | Standard Deviation of Mean | MERGE | → #150 | Same quantity as SEM |
| 150 | Standard Error of the Mean (SEM) | KEEP | `num/stats` NUM-STATS-023 | Precision of the mean |
| 151 | stat tool | MERGE | → #139 | A bundle of descriptive measures |
| 152 | Statistical Mode | KEEP | `num/stats` NUM-STATS-024 | Location |
| 153 | Stolen Base Percentage | MOVE | `life/sports` | Baseball |
| 154 | Sum | DROP | — | A loop; the compensated sum `num/stats` needs inside is unexported |
| 155 | Sum of Squares | KEEP | `num/stats` NUM-STATS-025 | Σ(x − x̄)², the building block of variance |
| 156 | Sum Squared | MERGE | → #155 | Ambiguous name (Σx² or (Σx)²); both are internal steps of #155 |
| 157 | t Distribution Lookup | MOVE | `num/prob` | Distribution |
| 158 | t-Dist | MERGE | → #157 | Same distribution |
| 159 | t-test critical value lookup | MOVE | `num/prob` | t quantile |
| 160 | The Sorensen Coefficient | MOVE | `sci/ecology` | Similarity of species lists |
| 161 | Times on Base | MOVE | `life/sports` | Baseball |
| 162 | Total Average | MOVE | `life/sports` | Baseball |
| 163 | Total Bases | MOVE | `life/sports` | Baseball |
| 164 | Variance (Population) | MERGE | → #107 | Duplicate |
| 165 | Variance (Sample) | MERGE | → #140 | Duplicate |
| 166 | Variance(Bernoulli) | MOVE | `num/prob` | Distribution variance |
| 167 | Variance(Bernoulli) | MERGE | → #166 | Duplicate |
| 168 | Variance(Binomial) | MOVE | `num/prob` | Distribution variance |
| 169 | Variance(Binomial) | MERGE | → #168 | Duplicate |
| 170 | Variance(Geometric) | MOVE | `num/prob` | Distribution variance |
| 171 | Variance(Geometric) | MERGE | → #170 | Duplicate |
| 172 | Variance(Hypergeometric) | MOVE | `num/prob` | Distribution variance |
| 173 | Variance(Hypergeometric) | MERGE | → #172 | Duplicate |
| 174 | Variance(Pascal) | MOVE | `num/prob` | Distribution variance |
| 175 | Variance(Pascal) | MERGE | → #174 | Duplicate |
| 176 | Variance(Uniform) | MOVE | `num/prob` | Distribution variance |
| 177 | Variance(Uniform) | MERGE | → #176 | Duplicate |
| 178 | Virtual Dice | DROP | — | Random |
| 179 | Watts and Strogatz Degree Distribution | MOVE | `num/prob` | Degree distribution of a random-graph model |
| 180 | Weighted Average | KEEP | `num/stats` NUM-STATS-026 | Weighted mean |
| 181 | Weighted Geometric Mean (bounded criteria) | MERGE | → #182 | Same formula with input bounds |
| 182 | Weighted Geometric Mean (DataSet criteria) | KEEP | `num/stats` NUM-STATS-027 | Weighted geometric mean |
| 183 | Weighted Geometric Mean (ten input criteria) | MERGE | → #182 | Fixed n = 10 |
| 184 | Weighted Grade Average | MERGE | → #180 | Weighted mean with a domain label |
| 185 | Wilcoxon signed rank test | LATER | later inference slice | Wilcoxon |
| 186 | x-Dist | MOVE | `num/prob` | A distribution lookup; the name does not say which one |
| 187 | y from P(y) Mean and SD | MOVE | `num/prob` | Normal quantile |
| 188 | z from P(y) | MOVE | `num/prob` | Standard normal quantile |
| 189 | z SCORE | KEEP | `num/stats` NUM-STATS-028 | Standardized value |
| 190 | z SCORE formula | MERGE | → #189 | Duplicate |
| 191 | z-score | MERGE | → #189 | Duplicate |
| 192 | z-score equation (find x) | KEEP | `num/stats` NUM-STATS-029 | Inverse of #189 |
| 193 | z-score equation (find z) | MERGE | → #189 | Duplicate |
| 194 | z-score from known mean & SD | MERGE | → #189 | Duplicate |
| | **Data items** | | | |
| d1 | 1 sigma | MOVE | `num/prob` | Normal coverage within ±1σ (a constant of the distribution) |
| d2 | 2 sigma | MOVE | `num/prob` | Normal coverage within ±2σ |
| d3 | 3 sigma | MOVE | `num/prob` | Normal coverage within ±3σ |
| d4 | 6 sigma | MOVE | `num/prob` | Normal coverage within ±6σ |

### LATER: the inferential list (for a later slice)
ANOVA (#1, #58, #93, #94, #127) · chi-square (#18–#23, c5) · t-tests (#100, #101, #143) · Mann-Whitney (#76) ·
Wilcoxon (#185) · Spearman (#145) · inference on the least-squares slope (#73) · effect size r² from t (#45) ·
test-choice guides (#24–#26, c8) · small-sample statistics (c1).

## Step 2: Research, one entry per KEEP

Common to every entry unless it says otherwise:
- Inputs are `[]float64` in the caller's unit. A result keeps that unit; a variance is in unit².
- The input is read, never modified.
- Every function returns `ErrEmpty` on empty input and `ErrNotFinite` (with the index) on NaN or ±Inf.
- Sums use the compensated (Neumaier) summation already shipped in S1 (doc/sources/neumaier_1974.md),
  as an unexported helper.

### ArithmeticMean      (ID reserved: NUM-STATS-001)
Formula: Ȳ = (Σ_{i=1..N} Y_i) / N, where Y_i is the i-th value and N the count.
Inputs and units: values []float64, any unit, N ≥ 1.
Output: float64, the mean, in the unit of values.
Errors: ErrEmpty; ErrNotFinite; ErrOverflow if the sum leaves the float64 range.
Primary source: [E] NIST/SEMATECH e-Handbook §1.3.5.1 "Measures of Location",
  https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm: "the sum of the data points divided by the number of data points".
Worked example: NIST StRD Univariate dataset NumAcc1: 10000001, 10000003, 10000002 → certified mean 10000002 (exact).
  https://www.itl.nist.gov/div898/strd/univ/data/NumAcc1.dat. NumAcc4 (n = 1001, certified 10000000.2) is the cancellation stress case.
Variants: none (weighted: NUM-STATS-026; grouped: NUM-STATS-005).
Best approach: one pass, compensated sum, then divide. O(n) time, O(1) space.
Notes: no naive Σ then /N without compensation; the StRD NumAcc cases exist to catch exactly that.

### RelativeFrequency      (ID reserved: NUM-STATS-002)
Formula: r_j = c_j / N, where c_j is the count in class j and N = Σ c_j.
Inputs and units: counts []int (from NUM-STATS-008), dimensionless, each ≥ 0, Σ > 0.
Output: []float64 fractions that sum to 1 (×100 for percent is the caller's job).
Errors: ErrEmpty; ErrNegative for a negative count; ErrZeroTotal when Σ c_j = 0.
Primary source: [E] e-Handbook §1.3.3.14 "Histogram", https://www.itl.nist.gov/div898/handbook/eda/section3/histogra.htm:
  "The normalized count is the count in a class divided by the total number of observations."
Worked example: HAND-CHECKED: counts 1, 3, 1 → 0.2, 0.6, 0.2.
Variants: density form, count / (N × class width), also on that page: a second function if wanted.
Best approach: one pass for the total, one to divide. O(b) for b classes.
Notes: the total is an integer sum: check it for overflow.

### CoefficientOfVariation      (ID reserved: NUM-STATS-003)
Formula: cv = s / x̄, where s is the sample standard deviation (divisor N−1) and x̄ the mean. Percent RSD = 100·s/|x̄|.
Inputs and units: values []float64, ratio scale (a unit with a true zero), N ≥ 2.
Output: float64, dimensionless.
Errors: ErrEmpty; ErrTooFew (N < 2); ErrZeroMean when x̄ = 0; ErrNotFinite.
Primary source: [D] NIST Dataplot Reference Manual, "COEFFICIENT OF VARIATION",
  https://itl.nist.gov/div898/software/dataplot/refman2/auxillar/coefvari.htm: "cv = s/x̄"; "relsd = 100 s/|x̄|". The e-Handbook §1.3.5.6 does not define it.
Worked example: HAND-CHECKED from StRD certified values: NumAcc1 → 1/10000002 = 9.9999980e-08;
  Michelso → 0.0790105478190518 / 299.8524 = 2.6349813e-04.
Variants: percent RSD (rows #103, #126) as a second function or a ×100 by the caller (Q5); Dataplot's
  small-sample unbiased form cv·(1 + 1/(4n)) is noted, not proposed.
Best approach: two passes (mean, then the sum of squares), O(n) time, O(1) space.
Notes: must NOT return ±Inf when the mean is 0, and must not silently take |x̄| (only the RSD form uses the absolute value).

### CohenD      (ID reserved: NUM-STATS-004)
SOURCE NOT FOUND (no formula given). The citation is confirmed: Cohen, *Statistical Power Analysis for the
Behavioral Sciences*, 2nd ed., Lawrence Erlbaum, 1988; Routledge reissue 2013, doi:10.4324/9780203771587.
No page of the book could be read; secondary statements of d do not count. It needs a library copy
(section and equation number) before it is specified. Dependency when sourced: PooledStandardDeviation.

### GroupedMean      (ID reserved: NUM-STATS-005)
SOURCE NOT FOUND (no formula given). No acceptable source states the mean of grouped data from class midpoints;
only school and revision sites came up. It needs a textbook section checked in a physical copy.

### GroupedMedian      (ID reserved: NUM-STATS-006)
Formula: p = l + ((p%·N/100) − F)/f · C. For the median p% = 50: l is the lower limit of the class holding
  the median, N the total frequency, F the cumulative frequency below that class, f its frequency, C its width.
Inputs and units: lower class limits []float64 (ascending, the value's unit), class width float64 > 0 (same unit) or upper limits, counts []int ≥ 0.
Output: float64 in the value's unit.
Errors: ErrEmpty; length mismatch; ErrNegative count; ErrZeroTotal; limits not ascending or overlapping.
Primary source: [D] NIST Dataplot Reference Manual, "GROUPED PERCENTILE",
  https://itl.nist.gov/div898/software/dataplot/refman2/auxillar/groupper.htm. It cites Selby, *CRC Standard
  Mathematical Tables*, 22nd ed., CRC Press, 1974, p. 572 (not checked). It "implicitly assumes that the frequencies are uniformly distributed within a class".
Worked example: HAND-CHECKED: classes [0,10) [10,20) [20,30) with counts 2, 5, 3 (N = 10). The median class is
  [10,20) (cumulative 2 < 5 ≤ 7), so 10 + (5 − 2)/5·10 = 16.
Variants: any percentile, the same formula with p% in place of 50.
Best approach: one cumulative pass, plus `algo/search.LowerBound` on the cumulative counts. O(b) time.
Notes: an estimate under the uniform-within-class assumption; it must say so and never claim to be the data's median.

### GroupedMode      (ID reserved: NUM-STATS-007)
SOURCE NOT FOUND (no formula given). The interpolated mode L + d1/(d1+d2)·h appeared only on school and
revision sites.

### FrequencyDistribution      (ID reserved: NUM-STATS-008)
Formula: for class edges e_0 < e_1 < … < e_b, c_j = #{ i : e_{j−1} ≤ Y_i < e_j } (the last class closed on the right).
Inputs and units: values []float64; edges []float64 ascending, in the same unit, len ≥ 2.
Output: []int, b counts.
Errors: ErrEmpty; ErrNotFinite; edges not strictly ascending; ErrOutOfRange for a value outside [e_0, e_b] (Q: error or count it?).
Primary source: [E] e-Handbook §1.3.3.14 "Histogram": "splitting the range of the data into equal-sized bins (called classes).
  Then for each bin, the number of points from the data set that fall into each bin are counted." No bin-count rule is named.
Worked example: HAND-CHECKED: values 1, 2, 2, 3, 5 with edges 0, 2, 4, 6 → 1, 3, 1.
Variants: equal-width bins from (min, max, b) as a helper; the cumulative form ("counts for that bin plus all bins for smaller values" [WF]).
Best approach: `algo/search.UpperBound` per value on the edges: O(n log b) time, O(b) for the result. Equal-width bins can be computed directly in O(n).
Notes: no automatic bin rule (Sturges, Scott) without its own source; the page names Scott but gives no formula.

### LinearRegression      (ID reserved: NUM-STATS-009)
Formula: β̂1 = Σ(x_i − x̄)(y_i − ȳ) / Σ(x_i − x̄)², β̂0 = ȳ − β̂1·x̄ (the least-squares straight line y = β0 + β1 x).
Inputs and units: x []float64 (unit X), y []float64 (unit Y), equal length, n ≥ 2, x not all equal.
Output: intercept float64 (unit Y), slope float64 (unit Y per unit X).
Errors: ErrEmpty; length mismatch; ErrTooFew (n < 2); ErrZeroVariance when all x are equal; ErrNotFinite.
Primary source: [E] e-Handbook §4.4.3.1 "Least Squares", https://www.itl.nist.gov/div898/handbook/pmd/section4/pmd431.htm
  ("yields the following estimators for the parameters").
Worked example: NIST StRD Linear Least Squares dataset Norris (36 observations), https://www.itl.nist.gov/div898/strd/lls/data/LINKS/DATA/Norris.dat:
  certified B0 = −0.262323073774029, B1 = 1.00211681802045, R² = 0.999993745883712.
Variants: R² and residual SD as later functions (the same StRD file certifies them).
Best approach: two passes (means, then the centred cross-sums). O(n) time, O(1) space. The naive Σxy − n x̄ȳ form fails StRD's harder sets.
Notes: a descriptive fit only. No standard errors or p-values here (those are LATER #73).

### MatthewsCorrelation      (ID reserved: NUM-STATS-010)
SOURCE NOT FOUND (no formula given). The citation is confirmed: Matthews, "Comparison of the predicted and
observed secondary structure of T4 phage lysozyme", *Biochimica et Biophysica Acta – Protein Structure*,
1975, 405(2):442–451, doi:10.1016/0005-2795(75)90109-9. The text was not reachable (paywall), so neither the
formula nor the zero-margin case can be quoted. It needs a library copy.

### MeanAbsoluteDeviation      (ID reserved: NUM-STATS-011)
Formula: AAD = Σ|Y_i − Ȳ| / N, where Ȳ is the mean.
Inputs and units: values []float64, any unit, N ≥ 1.
Output: float64 in the unit of values.
Errors: ErrEmpty; ErrNotFinite.
Primary source: [E] e-Handbook §1.3.5.6 "Measures of Scale", https://www.itl.nist.gov/div898/handbook/eda/section3/eda356.htm
  ("average absolute deviation", "AAD = Σ(|Y_i − Ȳ|)/N").
Worked example: HAND-CHECKED: NumAcc1 → deviations 1, 1, 0 → 2/3. The NIST §7.2.6.2 resistivity data (12 values) → 0.035876.
Variants: the median absolute deviation, MAD = median|Y_i − Ỹ| (same page, no 1.4826 factor), is a separate function if wanted (Q5).
Best approach: two passes (mean, then absolute deviations), compensated. O(n) time, O(1) space.
Notes: the name "MAD" is ambiguous in practice. The function name must say "mean", and the doc must say "about the mean".

### Median      (ID reserved: NUM-STATS-012)
Formula: odd N: Ỹ = Y_((N+1)/2); even N: Ỹ = (Y_(N/2) + Y_(N/2+1)) / 2, where Y_(i) is the i-th smallest value.
Inputs and units: values []float64, any unit, N ≥ 1.
Output: float64 in the unit of values.
Errors: ErrEmpty; ErrNotFinite.
Primary source: [E] e-Handbook §1.3.5.1 "Measures of Location" (eda351.htm).
Worked example: HAND-CHECKED on the 12 resistivity values printed in e-Handbook §7.2.6.2 → (95.1567 + 95.1591)/2 = 95.1579.
Variants: none (the generic lower/upper medians are ALG-SELECT-002).
Best approach: `algo/selection.Median` (ALG-SELECT-002) for the two middles, then their mean computed as
  lower + (upper − lower)/2 so it cannot overflow. O(n) worst case, one O(n) copy.
Notes: must NOT sort the caller's slice.

### Midrange      (ID reserved: NUM-STATS-013)
Formula: (smallest + largest) / 2.
Inputs and units: values []float64, any unit, N ≥ 1.
Output: float64 in the unit of values.
Errors: ErrEmpty; ErrNotFinite.
Primary source: [E] e-Handbook §1.3.5.1 (eda351.htm): "Mid-range = (smallest + largest)/2", "not robust".
Worked example: HAND-CHECKED on the §7.2.6.2 resistivity data → (95.0610 + 95.1990)/2 = 95.1300.
Variants: none. Row #84 "Mid Point" is read as the midrange; if it meant a class midpoint, that is (lower + upper)/2 inside the grouped functions (Q6).
Best approach: one pass for min and max; min + (max − min)/2 avoids overflow. O(n), O(1).
Notes: the doc must say "not robust", quoting NIST.

### PearsonCorrelation      (ID reserved: NUM-STATS-014)
Formula: r = S_xy / (√S_xx · √S_yy), with S_xx = Σ(X_i − X̄)², S_yy = Σ(Y_i − Ȳ)², S_xy = Σ(X_i − X̄)(Y_i − Ȳ).
Inputs and units: x, y []float64 of equal length (any units), n ≥ 2.
Output: float64 in [−1, 1], dimensionless.
Errors: ErrEmpty; length mismatch; ErrTooFew; ErrZeroVariance when x or y is constant; ErrNotFinite.
Primary source: [D] NIST Dataplot Reference Manual, "CORRELATION",
  https://itl.nist.gov/div898/software/dataplot/refman2/auxillar/correlat.htm. The e-Handbook has no formula (checked §1.3.5, §1.3.3.26, §4.4.4).
Worked example: derived from StRD Norris: r = +√R² = 0.999996872936967 (the sign is the slope's). HAND-CHECKED from certified values.
Variants: none (Spearman is LATER #145).
Best approach: two passes (means, then the three centred sums), sharing an unexported helper with LinearRegression. O(n), O(1).
Notes: clamp to [−1, 1] after rounding, and document the clamp. No significance test here.

### Percentile      (ID reserved: NUM-STATS-015)
Formula (Hyndman & Fan type 6, the e-Handbook's): p(N+1) = k + d with k an integer and 0 ≤ d < 1, Y_[i] sorted ascending.
  For 0 < k < N: Y(p) = Y[k] + d(Y[k+1] − Y[k]). For k = 0: Y[1]. For k ≥ N: Y[N].
Inputs and units: values []float64 (any unit), p float64 in [0, 1] (dimensionless), method (Q1).
Output: float64 in the unit of values.
Errors: ErrEmpty; ErrNotFinite; ErrOutOfRange for p ∉ [0, 1].
Primary source: [E] e-Handbook §7.2.6.2 "Percentiles", https://www.itl.nist.gov/div898/handbook/prc/section2/prc262.htm.
  [E] Hyndman, Fan, "Sample quantiles in statistical packages", *The American Statistician*, 1996, 50(4):361–365,
  doi:10.1080/00031305.1996.10473566: nine definitions, eq. (1) and Table 3; they recommend type 8.
Worked example: e-Handbook §7.2.6.2, 12 resistivity values (ohm·cm), p = 0.90: printed "Y(90) … 95.1981" (type 6),
  "R7 gives 95.1957", "R8 gives 95.1972". All three re-derived by hand here (95.19807, 95.19568, 95.19724).
Variants: H&F types 1–9 (Q1).
Best approach: two adjacent order statistics, by `algo/selection` on one copy. O(n) worst case, one O(n) copy.
  It needs Y[k] and Y[k+1] together: select k, then the minimum of the part to its right (the Median trick). Exposing
  that as an ALG function would be a new ID in `algo/selection`.
Notes: the method must be visible at the call site; never a silent default.

### PooledStandardDeviation      (ID reserved: NUM-STATS-016)
Formula: s_p² = ((N1 − 1)s1² + (N2 − 1)s2²) / (N1 + N2 − 2); s_p = √s_p².
Inputs and units: n1, n2 int (counts, each ≥ 1, n1 + n2 ≥ 3); s1, s2 float64 ≥ 0 (same unit). A slice form on two samples is a variant.
Output: float64 in the unit of the data.
Errors: ErrTooFew (n1 + n2 < 3, or either n < 1); ErrNegative SD; ErrNotFinite.
Primary source: [E] e-Handbook §1.3.5.3 "Two-Sample t-Test", https://www.itl.nist.gov/div898/handbook/eda/section3/eda353.htm.
Worked example: same page, AUTO83B.DAT: N1 = 249, s1 = 6.41470; N2 = 79, s2 = 6.10771 → printed "s_p = 6.34260" (reproduced by hand: 6.34260).
Variants: k samples (Σ(n_i − 1)s_i² / Σ(n_i − 1)): needs its own source.
Best approach: scalar arithmetic, O(1). Weight with float64 to avoid integer overflow.
Notes: the t statistic itself is LATER.

### PopulationVariance      (ID reserved: NUM-STATS-017)
SOURCE NOT FOUND (no formula given). The e-Handbook §1.3.5.6 defines only the N−1 form; no page with the
divisor-N variance was readable from an acceptable source. Casella & Berger and Hogg-McKean-Craig need checking in a physical copy (Q2).

### Range      (ID reserved: NUM-STATS-018)
Formula: largest value − smallest value.
Inputs and units: values []float64, any unit, N ≥ 1.
Output: float64 ≥ 0 in the unit of values.
Errors: ErrEmpty; ErrNotFinite; ErrOverflow if max − min exceeds float64 (e.g. ±MaxFloat64).
Primary source: [E] e-Handbook §1.3.5.6 (eda356.htm): "the largest value minus the smallest value in a data set".
Worked example: HAND-CHECKED on the §7.2.6.2 resistivity data → 95.1990 − 95.0610 = 0.1380.
Variants: none.
Best approach: one pass for min and max. O(n), O(1).
Notes: none beyond the overflow check.

### SampleStandardDeviation      (ID reserved: NUM-STATS-019)
Formula: s = √( Σ(Y_i − Ȳ)² / (N − 1) ).
Inputs and units: values []float64, any unit, N ≥ 2.
Output: float64 ≥ 0 in the unit of values.
Errors: ErrEmpty; ErrTooFew (N < 2); ErrNotFinite.
Primary source: [E] e-Handbook §1.3.5.6 (eda356.htm). NIST StRD certifies "Sample Standard Deviation (denom. = n-1)".
Worked example: StRD NumAcc1 → certified 1 (exact); Michelso (n = 100) → 0.0790105478190518; Mavro (n = 50) → 0.000429123454003053.
Variants: population SD is NUM-STATS-022.
Best approach: the square root of SampleVariance (NUM-STATS-021).
Notes: NumAcc2–4 (certified SD 0.1) must pass. They defeat the one-pass Σx² − (Σx)²/N formula.

### SummaryStatistics      (ID reserved: NUM-STATS-020)
Formula: no formula of its own. It returns N, mean, sample SD, min, max, median and range, each defined by its own entry above.
Inputs and units: values []float64, any unit, N ≥ 2 (for the SD).
Output: a struct with those fields, in the unit of values (N dimensionless).
Errors: the union of the parts' errors.
Primary source: [E] e-Handbook §1.3.5.1 (location) and §1.3.5.6 (scale). §1.3.5 has no list named "summary
  statistics"; the set is ours and must be stated as such. StRD Univariate certifies mean and SD.
Worked example: StRD NumAcc1 → N 3, mean 10000002, SD 1, min 10000001, max 10000003, median 10000002, range 2 (mean and SD certified; the rest HAND-CHECKED).
Variants: none (MERGEs rows c6, c7, c9, c10, #13, #35–#41, #90, #142, #151).
Best approach: a single O(n) copy reused by median selection; mean and sum of squares in two passes. O(n) time, one allocation.
Notes: it must call the public functions' cores so the numbers can never disagree with them.

### SampleVariance      (ID reserved: NUM-STATS-021)
Formula: s² = Σ(Y_i − Ȳ)² / (N − 1).
Inputs and units: values []float64, any unit, N ≥ 2.
Output: float64 ≥ 0 in unit².
Errors: ErrEmpty; ErrTooFew (N < 2); ErrNotFinite; ErrOverflow.
Primary source: [E] e-Handbook §1.3.5.6 (eda356.htm). For the method: Welford, "Note on a method for calculating
  corrected sums of squares and products", *Technometrics*, 1962, 4(3):419–420, doi:10.1080/00401706.1962.10490022.
  Knuth, TAOCP Vol. 2, 3rd ed., §4.2.2, is NOT VERIFIED (no copy); it is believed to be eqs. (15)–(16), to check.
Worked example: StRD NumAcc1 → 1 (the certified SD squared); NumAcc4 (n = 1001) → 0.01 (from the certified SD 0.1).
Variants: population (NUM-STATS-017, blocked on its source).
Best approach: two-pass corrected sum of squares (NUM-STATS-025) on a slice, which is the most accurate. Welford's
  one-pass update is for streaming, a later accumulator type. O(n), O(1).
Notes: never the textbook one-pass Σx² − (Σx)²/N, which can even go negative.

### PopulationStandardDeviation      (ID reserved: NUM-STATS-022)
SOURCE NOT FOUND (no formula given). It is the square root of NUM-STATS-017 and is blocked on the same source (Q2).

### StandardErrorOfMean      (ID reserved: NUM-STATS-023)
Formula: SE = s / √N, where s is the sample SD (N−1).
Inputs and units: values []float64 (any unit), N ≥ 2; or (s float64 ≥ 0, n int ≥ 2) as a variant.
Output: float64 in the unit of values.
Errors: ErrEmpty; ErrTooFew; ErrNotFinite.
Primary source: [E] e-Handbook §1.3.5.2 "Confidence Limits for the Mean" (eda352.htm) uses s/√N inside Ȳ ± t·s/√N.
  [D] Dataplot names it ("s_mean = s/√n", from a search-result excerpt, page not opened).
Worked example: e-Handbook §1.3.5.2, ZARR13.DAT: N = 195, s = 0.022789 → s/√N = 0.0016320. This agrees with the printed
  95% interval (9.258242, 9.264679): half-width / t (1.9723) = 0.0016321. HAND-CHECKED from printed values.
Variants: from summary values (s, n).
Best approach: SampleStandardDeviation, then divide. O(n).
Notes: the confidence interval itself is LATER.

### Mode      (ID reserved: NUM-STATS-024)
Formula: the value occurring with the greatest frequency; "not necessarily unique".
Inputs and units: values []T, T any ordered type (Q3), N ≥ 1.
Output: the mode(s), ascending, and their count.
Errors: ErrEmpty; NaN follows the total order of D-002 (it equals itself for counting).
Primary source: [E] e-Handbook §1.3.5.1 (eda351.htm): "the value of the random sample that occurs with the greatest frequency. It is not necessarily unique."
Worked example: HAND-CHECKED: 2, 3, 3, 5, 5, 5 → [5] ×3. Ties: 1, 1, 2, 2, 3 → [1, 2] ×2.
Variants: none. Grouped data is NUM-STATS-007.
Best approach: sort a copy and scan runs: O(n log n) worst case, O(n) copy, deterministic. A map would be O(n)
  expected but with a random iteration order, so it would need a sort of the ties anyway.
Notes: the tie policy is Q7; all values distinct returns every value (count 1), as the definition implies.

### SumOfSquares      (ID reserved: NUM-STATS-025)
Formula: SSQ = Σ_{i=1..n} (X_i − X̄)² (the sum of squared deviations about the mean).
Inputs and units: values []float64, any unit, n ≥ 1.
Output: float64 ≥ 0 in unit².
Errors: ErrEmpty; ErrNotFinite; ErrOverflow.
Primary source: [D] NIST Dataplot Reference Manual, "SUM OF SQUARES FROM MEAN",
  https://itl.nist.gov/div898/software/dataplot/refman2/auxillar/ssqmean.htm. It is the numerator of the e-Handbook §1.3.5.6 variance.
Worked example: StRD NumAcc1 → 2 (that is, (N−1)·s² = 2·1); the resistivity data → 0.02164115 (HAND-CHECKED).
Variants: the raw Σx² (Dataplot "SUM OF SQUARES", ssq.htm) is not proposed. Row #156 "Sum Squared" is ambiguous (Q6).
Best approach: two passes. Pass 1 gives the compensated mean. Pass 2 sums squared deviations with the correction term Σ(x−x̄) (the corrected two-pass algorithm). O(n), O(1).
Notes: the core of 019, 021, 017, 003, 023.

### WeightedMean      (ID reserved: NUM-STATS-026)
Formula: W̄ = Σ W_i X_i / Σ W_i.
Inputs and units: values []float64 (any unit); weights []float64 (dimensionless, each ≥ 0, at least one > 0); equal length.
Output: float64 in the unit of values.
Errors: ErrEmpty; length mismatch; ErrNegative weight; ErrZeroTotal (all weights 0); ErrNotFinite.
Primary source: [D] NIST Dataplot Reference Manual, "WEIGHTED TRIMMED MEAN" (weigtmea.htm) and "WEIGHTED ORDER STATISTIC
  MEAN" (wos.htm): the same ratio form, and "none of the weights can be negative". The plain WEIGHTED MEAN page redirects and
  could not be read; finding a cleaner source is Q4.
Worked example: HAND-CHECKED: values 80, 90 with weights 1, 3 → 87.5.
Variants: none (grades, row #184, is just this).
Best approach: one pass with two compensated sums. O(n), O(1).
Notes: the weights are not normalised by the caller; the function divides.

### WeightedGeometricMean      (ID reserved: NUM-STATS-027)
SOURCE NOT FOUND (no formula given). Dataplot gives only the unweighted geometric mean, (∏X_i)^(1/n)
(geommean.htm). No acceptable source for the weighted form was found.

### ZScore      (ID reserved: NUM-STATS-028)
Formula: z = (x − x_ref) / σ, where x_ref is the reference mean and σ the standard deviation (the caller supplies both).
Inputs and units: x, mean float64 (same unit); sd float64 > 0 (same unit).
Output: float64, dimensionless.
Errors: ErrZeroVariance (sd = 0); ErrNegative (sd < 0); ErrNotFinite.
Primary source: [D] NIST Dataplot "ISO 13528 ZSCORE" (zscore.htm): "Z(i) = (X(i) − Xref)/sigma", citing ISO 13528:2005, pp. 25–26 (not checked);
  and Dataplot "STANDARDIZE" [WF]. No e-Handbook section defines z.
Worked example: HAND-CHECKED: x = 12, mean 10, sd 2 → 1.
Variants: a slice form (standardize every value with its own sample mean and SD) as a separate function if wanted.
Best approach: scalar, O(1).
Notes: ISO 13528's |z| > 2 "warning" and > 3 "action" labels are interpretation and do NOT belong in this function.

### RawScoreFromZScore      (ID reserved: NUM-STATS-029)
Formula: x = x_ref + z·σ (NUM-STATS-028 solved for x; no separate source states it).
Inputs and units: z float64 (dimensionless); mean, sd float64 (same unit; sd ≥ 0).
Output: float64 in the unit of mean.
Errors: ErrNegative (sd < 0); ErrNotFinite; ErrOverflow.
Primary source: as NUM-STATS-028 [D], by algebra (Q4 covers whether that is acceptable).
Worked example: HAND-CHECKED: z = 1, mean 10, sd 2 → 12 (the inverse of the 028 example).
Variants: none.
Best approach: scalar, O(1); math.FMA(z, sd, mean) for one rounding.
Notes: none.

## Step 3: Build order and dependencies

Lowest first. "→" means "calls". The S1 algorithms used are in brackets.

0. Unexported helpers in `num/stats`: the finiteness check; compensated summation (Neumaier, as in ALG-WINDOW-001's core,
   re-implemented here; `algo/window` does not export it).
1. 001 ArithmeticMean.
2. 025 SumOfSquares → 001.
3. 021 SampleVariance → 025. 019 SampleStandardDeviation → 021. (017 and 022, population forms → 025: blocked on source.)
4. 023 StandardErrorOfMean → 019. 003 CoefficientOfVariation → 001, 019. 011 MeanAbsoluteDeviation → 001.
5. 018 Range and 013 Midrange: one min/max pass, no dependencies.
6. 012 Median → [ALG-SELECT-002 Median]. 015 Percentile → [ALG-SELECT-001 KthSmallest, plus a two-order-statistic select; maybe a new ALG ID].
7. 024 Mode: sort a copy (`slices.Sort`) and scan.
8. 008 FrequencyDistribution → [ALG-SEARCH-002 UpperBound]. 002 RelativeFrequency → 008's output.
9. 006 GroupedMedian → [ALG-SEARCH-001 LowerBound] on the cumulative counts. (005 and 007: blocked on source.)
10. 026 WeightedMean → helper 0. (027: blocked on source.)
11. 014 PearsonCorrelation and 009 LinearRegression → 001 and one shared unexported centred-cross-sum helper.
12. 016 PooledStandardDeviation (scalar). 028 ZScore and 029 RawScoreFromZScore (scalar). (004 CohenD → 016, 001; 010 Matthews: blocked on source.)
13. 020 SummaryStatistics → 001, 019, 012, 018 (the composite comes last).
Window: no KEEP is a rolling form. A rolling mean already exists (ALG-WINDOW-002). Rolling variance or SD would build on
`algo/window` plus 025's method; propose it as a later function if wanted.
