# S8: Statistics (descriptive)

Status: SHIPPED 2026-11-12
Work folder: `doc/math_build_2026-11-12/`
Decisions: D-010 · Assumptions: proposed in `doc/math_build_2026-11-12/notes_stats.md` (A-? for the lead)
Criteria: `doc/catalog/stats.md`, Step 2 (one entry per ID). Open questions there are settled by D-010.

## Functions

| ID | Name | Signature | Units |
|---|---|---|---|
| `NUM-STATS-001` | `ArithmeticMean` | `func ArithmeticMean(values []float64) (float64, error)` | values any unit; result same |
| `NUM-STATS-002` | `RelativeFrequency` | `func RelativeFrequency(counts []int) ([]float64, error)` | counts and result dimensionless |
| `NUM-STATS-003` | `CoefficientOfVariation` | `func CoefficientOfVariation(values []float64) (float64, error)` | ratio-scale unit; result dimensionless |
| `NUM-STATS-004` | `CohenD` | `func CohenD(meanA, meanB, sd float64) (float64, error)` | one unit; result dimensionless |
| `NUM-STATS-005` | `GroupedMean` | `func GroupedMean(edges []float64, counts []int) (float64, error)` | edges in data unit; result same |
| `NUM-STATS-006` | `GroupedMedian` | `func GroupedMedian(edges []float64, counts []int) (float64, error)` | edges in data unit; result same |
| `NUM-STATS-007` | `GroupedMode` | HELD BACK (no confirmed source) | |
| `NUM-STATS-008` | `FrequencyDistribution` | `func FrequencyDistribution(values, edges []float64) ([]int, error)` | one unit; counts |
| `NUM-STATS-009` | `LinearRegression` | `func LinearRegression(x, y []float64) (intercept, slope float64, err error)` | intercept unit Y; slope Y per X |
| `NUM-STATS-010` | `MatthewsCorrelation` | `func MatthewsCorrelation(truePositives, trueNegatives, falsePositives, falseNegatives int) (float64, error)` | counts; result in [−1, 1] |
| `NUM-STATS-011` | `MeanAbsoluteDeviation` | `func MeanAbsoluteDeviation(values []float64) (float64, error)` | values unit |
| `NUM-STATS-012` | `Median` | `func Median(values []float64) (float64, error)` | values unit |
| `NUM-STATS-013` | `Midrange` | `func Midrange(values []float64) (float64, error)` | values unit |
| `NUM-STATS-014` | `PearsonCorrelation` | `func PearsonCorrelation(x, y []float64) (float64, error)` | result in [−1, 1] |
| `NUM-STATS-015` | `Percentile` | `func Percentile(values []float64, p float64, method QuantileMethod) (float64, error)` | p in [0, 1]; result values unit |
| `NUM-STATS-016` | `PooledStandardDeviation` | `func PooledStandardDeviation(n1 int, sd1 float64, n2 int, sd2 float64) (float64, error)` | sd unit |
| `NUM-STATS-017` | `PopulationVariance` | `func PopulationVariance(values []float64) (float64, error)` | unit² |
| `NUM-STATS-018` | `Range` | `func Range(values []float64) (float64, error)` | values unit |
| `NUM-STATS-019` | `SampleStandardDeviation` | `func SampleStandardDeviation(values []float64) (float64, error)` | values unit |
| `NUM-STATS-020` | `SummaryStatistics` | `func SummaryStatistics(values []float64) (Summary, error)` | fields in values unit |
| `NUM-STATS-021` | `SampleVariance` | `func SampleVariance(values []float64) (float64, error)` | unit² |
| `NUM-STATS-022` | `PopulationStandardDeviation` | `func PopulationStandardDeviation(values []float64) (float64, error)` | values unit |
| `NUM-STATS-023` | `StandardErrorOfMean` | `func StandardErrorOfMean(values []float64) (float64, error)` | values unit |
| `NUM-STATS-024` | `Mode` | `func Mode[T cmp.Ordered](values []T) (modes []T, count int, err error)` | type T |
| `NUM-STATS-025` | `SumOfSquares` | `func SumOfSquares(values []float64) (float64, error)` | unit² |
| `NUM-STATS-026` | `WeightedMean` | `func WeightedMean(values, weights []float64) (float64, error)` | values unit; weights dimensionless |
| `NUM-STATS-027` | `WeightedGeometricMean` | NOT BUILT (D-010: out, no source found) | |
| `NUM-STATS-028` | `ZScore` | `func ZScore(x, mean, sd float64) (float64, error)` | dimensionless |
| `NUM-STATS-029` | `RawScoreFromZScore` | `func RawScoreFromZScore(z, mean, sd float64) (float64, error)` | mean unit |
| `NUM-STATS-030` | `PercentRelativeStandardDeviation` | `func PercentRelativeStandardDeviation(values []float64) (float64, error)` | percent |
| `NUM-STATS-031` | `MedianAbsoluteDeviation` | `func MedianAbsoluteDeviation(values []float64) (float64, error)` | values unit |

Exported types: `QuantileMethod` with constants `HyndmanFan1` … `HyndmanFan9` and `DefaultQuantileMethod = HyndmanFan7`
(percentile.go); `Summary` (summary_statistics.go). Sentinels (arithmetic_mean.go): `ErrEmpty`, `ErrNotFinite`,
`ErrOverflow`, `ErrTooFew`, `ErrLengthMismatch`, `ErrNegative`, `ErrZeroTotal`, `ErrZeroMean`, `ErrZeroVariance`,
`ErrOutOfRange`, `ErrNotAscending`.

Shared internals (compensated_sum.go, unexported): `compensatedSum` (Neumaier), `checkFinite`, `checkScalar`,
`meanOf` (rescales by 2^-k on overflow, so a mean of finite values is always returned), `centredSquares` and
`centredProducts` (corrected two-pass, Chan–Golub–LeVeque; rescale by 2^-600 on overflow, x and y separately),
`midpoint`, `lerp`.

## Per function

Common to every slice function: input read, never modified; empty → `ErrEmpty`; NaN/±Inf → `ErrNotFinite` with the
index. Full source citations are in each doc comment; worked examples are the first test in each `_test.go`.

### Location: 001 ArithmeticMean, 012 Median, 013 Midrange, 024 Mode, 026 WeightedMean
- **Source**: e-Handbook §1.3.5.1 (001, 012, 013, 024); Dataplot WEIGHTED MEAN (ch2/weigmean.pdf) and WEIGHTED TRIMMED MEAN (026).
- **Best known**: Ω(n) to read the input. 001/013/026 O(n), O(1). 012 O(n) worst case by introselect (ALG-SELECT-002), one copy. 024 O(n log n): sort a copy and scan runs.
- **Approach**: compensated sums; median via `algo/selection.Median` then an overflow-safe midpoint. 026 redoes its sums on weights normalised by the largest and values scaled by 2^-k (k the binary exponent of the largest |value|) when a sum overflows or a nonzero product falls below 2^-1022 (fix 2026-11-12; was silent 0 for weights and values near 1e-300). **Lost**: naive Σ/n (fails StRD NumAcc); sorting for the median (O(n log n)); a hash map for Mode (O(n) expected but random order, ties need a sort anyway, and NaN keys never match).
- **Worked example**: StRD NumAcc1 mean 10000002; resistivity median 95.1579, midrange 95.1300; mode 2,3,3,5,5,5 → [5]×3; weights 1,3 on 80,90 → 87.5.

### Scale: 017/021 variances, 019/022 SDs, 025 SumOfSquares, 023 SEM, 003 CV, 030 %RSD, 011 AAD, 031 MAD, 018 Range
- **Source**: e-Handbook §1.3.5.6 (021, 019, 011, 018, 031); Hogg–McKean–Craig 8th ed. §1.9 Def. 1.9.2 and Casella–Berger §2.3 (017, 022); Dataplot SUM OF SQUARES FROM MEAN (025), COEFFICIENT OF VARIATION (003), RELATIVE STANDARD DEVIATION (030), MEDIAN ABSOLUTE DEVIATION (031), STANDARD DEVIATION OF THE MEAN (023); Chan, Golub, LeVeque 1983 for the corrected two-pass method.
- **Best known**: O(n), O(1) for the moment-based ones (three passes); 031 O(n) by two selections, O(n) extra (one buffer: a copy of the values selected in place for the median, then overwritten by the deviations and selected in place again; fix 2026-11-12, was three allocations).
- **Approach**: corrected two-pass Σd² − (Σd)²/n with compensated sums. Range (fix 2026-11-12): the pass also tracks the largest |d|; when the sum overflows, or the largest |d| is below 2^-480 (squares would underflow), it is redone with values and mean scaled by 2^-k (k the binary exponent of the largest |d|, clamped at −1000; 600 when a deviation itself overflows), exact, and the scale is returned with the sum. A corrected sum that rounds to ≤ 0 while the deviations differ is floored at its lower bound 2·((max d − min d)/2)², so a zero sum, and `ErrZeroVariance`, means every value is equal. CV/%RSD divide by the scaled mean when k ≤ 0. **Lost**: a separate max-|d| pass before every sum (one extra pass always, for a case that is rare); per-element `math.Ldexp` (slow). **Lost (as before)**: one-pass Σx² − (Σx)²/n (catastrophic cancellation, can go negative; fails NumAcc2–4); Welford's update (one pass but less accurate on a slice; it is the later streaming accumulator).
- **Worked example**: StRD certified SDs (NumAcc1–4, Michelso, Mavro, Lew); Michelso CV 2.6349813e-04.
- **Errors**: `ErrTooFew` for n < 2 on sample forms; `ErrZeroMean` for CV/%RSD; `ErrOverflow` only when the result itself is out of range. A variance or sum of squares below the float64 range rounds to a subnormal or 0 (no error); the SDs, CV, %RSD and SEM stay accurate down to subnormal spreads.

### Frequency and grouped data: 008, 002, 005, 006
- **Source**: e-Handbook §1.3.3.14 (008, 002); OpenStax Introductory Statistics 2e §2.5 (005); Dataplot GROUPED PERCENTILE (006).
- **Best known**: 008 O(n log b) with binary search per value (`algo/search.UpperBound`); O(n) only for equal-width bins (a later variant). 002/005/006 O(b).
- **Approach**: 006 scans the cumulative counts once. **Lost**: `algo/search.LowerBound` on a cumulative array (same O(b) to build, plus an allocation).
- **Decision**: a value outside the edges is `ErrOutOfRange`, not dropped (catalog Q left open; proposed assumption).
- **Worked example**: 1,2,2,3,5 on edges 0,2,4,6 → 1,3,1; counts 2,5,3 on [0,10,20,30] → mean 16, median 16.

### Association: 009 LinearRegression, 014 PearsonCorrelation, 010 MatthewsCorrelation
- **Source**: e-Handbook §4.4.3.1 (009); Dataplot CORRELATION (014); Matthews 1975 with Baldi et al. 2000 eq. (11) and Chicco–Jurman 2020 eq. 2 (010).
- **Best known**: O(n), O(1) (009, 014); O(1) (010).
- **Approach**: shared corrected two-pass centred sums, x and y rescaled independently by the same rule as the scale group (tiny or overflowing spreads); 010 numerator exact in 128 bits (`math/bits.Mul64`). **Lost**: Σxy − n·x̄·ȳ (fails StRD's harder sets).
- **Worked example**: StRD Norris certified B0, B1 and R² (r = +√R²); MCC 90/80/20/10 → 0.70353.
- **Errors**: `ErrZeroVariance` for constant x (or y for 014); `ErrZeroTotal` for a zero margin in 010 (Baldi: undefined).

### Quantiles: 015 Percentile
- **Source**: Hyndman & Fan 1996, eq. (1), §2–3, Table 3; e-Handbook §7.2.6.2.
- **Best known**: O(n) worst case selection; one O(n) copy.
- **Approach**: `algo/selection.KthSmallest` for X(j), then one pass for X(j+1) (a tie straddling the rank, else the least larger value). **Lost**: sorting (O(n log n)); two selections (two copies); a new ALG function (not needed).
- **Numeric policy**: g within 4ε·max(1, |np + m|) of an integer is taken as that integer (as R's fuzz), so p = 0.3, n = 10 stays at X(3) for types 1–3.
- **Worked example**: resistivity p = 0.9: type 6 95.1981, 7 95.1957, 8 95.1972; nine-type table on 3,1,4,1,5,9,2,6,5,3.

### Two-sample and standardized: 016, 004, 028, 029
- **Source**: e-Handbook §1.3.5.3 (016); Cohen 1988 eq. 2.2.1 with Lakens 2013 and Goulet-Pelletier & Cousineau 2018 (004); Dataplot ISO 13528 ZSCORE and ISO 13528:2015 §9.4 (028, 029).
- **Best known**: O(1).
- **Approach**: 016 via `math.Hypot` (never overflows, clamped to max(sd1, sd2)); 004/028 divide each term when the difference overflows; 029 one FMA.
- **Worked example**: AUTO83B s_p = 6.34260; d = 1.62968 on AUTO83B means; z(12; 10, 2) = 1 and back.

### Composite: 020 SummaryStatistics
- Count, mean, sample SD, min, max, median, range from the same internals as the single functions (a test asserts exact equality). One allocation (the median's copy).

## Out of scope

GroupedMode (007, held: no confirmed source), WeightedGeometricMean (027, D-010 out), the k-sample pooled SD,
equal-width O(n) binning, a rolling/streaming variance (Welford, later accumulator), the inferential LATER list of
the catalog (t, chi-square, ANOVA, Spearman), Gamma (to `num/combinatorics`), rank order centroid (LATER).
