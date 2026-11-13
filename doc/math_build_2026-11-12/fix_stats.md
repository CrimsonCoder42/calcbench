# fix_stats: num/stats reviewer fixes (2026-11-12)

Owned paths only: `num/stats/**`, `doc/specs/s8_stats.md`, this file. Public API unchanged (no exported
name, parameter or result changed; no new exported identifier). Final checks: `gofmt -l num/stats` clean,
`go vet ./num/stats/` clean, `go test -count=1 ./num/stats/` →
`ok  	github.com/CrimsonCoder42/calcbench/num/stats	0.392s`. Benchmarks: `out/bench_stats_fix.txt`.

## 1. HIGH: silent underflow of squared deviations (FIXED)

Shared helpers, `num/stats/compensated_sum.go`:
- :87–103 constants `overflowScaleExp` (600, as before), `smallDeviation` = 2^-480, `minScaleExp` = −1000.
- :106 `deviationScaleExp(dMax)`: k = binary exponent of the largest |deviation| (math.Frexp), clamped at
  −1000; 600 when a deviation itself overflowed.
- :117 `needsRescale(ss, dMax)`: redo when the sum is not finite or 0 < dMax < 2^-480.
- :130 `centredSquares` / :144 `correctedSquares`: the unscaled pass now also tracks min/max deviation;
  a rescale multiplies values and mean by 2^-k (exact), so Σd² = ss·4^k with k now possibly negative.
- :161 `floorSquares`: a corrected sum that rounds to ≤ 0 while the deviations differ is floored at its
  true lower bound 2·((max d − min d)/2)², so ss = 0 (hence ErrZeroVariance) only when every value is equal.
- :180 `centredProducts` / :198 `correctedProducts`: same rule, x and y rescaled independently.
- `coefficient_of_variation.go:35–44` `sdOverMean` (CV, %RSD): when k ≤ 0 it divides by the scaled mean,
  sd/(mean·2^-k), since sd/mean could overflow before unscaling (test: mean 2^-1074, CV 2^74).
- `weighted_mean.go:48–80` (same class, from notes (d)): a nonzero product below 2^-1022 triggers the
  rescaled pass; the retry now scales values by the largest |value|'s exponent instead of 2^-bitLen(n).
- Callers needed no change: every one unscales with math.Ldexp, which takes negative exponents. Pooled SD
  (math.Hypot), ZScore, CohenD, RawScoreFromZScore form no squares; a tiny-SD test was added for pooled SD.
- Docs: `doc.go:11–16` Accuracy paragraph (range, ErrZeroVariance meaning, sub-range variance rounds);
  SampleVariance, PopulationVariance, SumOfSquares: "a result below the smallest float64 rounds to a
  subnormal or 0" sentence; Sample/PopulationStandardDeviation: accurate "even when s² would overflow or
  underflow". `doc/specs/s8_stats.md` Scale, Location and Association sections updated.

Tests added (expected values by hand from power-of-two scaling of hand-checked data, or exact `big.Rat`
on the float64 inputs; helpers `exactCentred`, `ratSqrt`, `ratFloat` in `sample_standard_deviation_test.go`):
`TestSampleStandardDeviation_tinySpread` (2^-700, 2^-1000 ≈ 1e-301, 2^-1070 subnormal, two smallest
subnormals, and 1e-200 decimal vs big.Rat), `TestPopulationStandardDeviation_tinySpread`,
`TestCoefficientOfVariation_tinySpread` (0.5 exact at 2^-700 and 2^-1070; 1e-200; mean 2^-1074 → 2^74),
`TestPercentRelativeStandardDeviation_tinySpread` (50), `TestStandardErrorOfMean_tinySpread`,
`TestSumOfSquares_tinySpread`, `TestSummaryStatistics_tinySpread`, `TestPearsonCorrelation_tinySpread`
(r = 1/2 for x = 1e-200·(1,2,3), y = 2,1,3, the reviewer's data reconstructed: S_xy ∝ 1, S_xx ∝ 2,
S_yy = 2; equal tiny x still ErrZeroVariance), `TestLinearRegression_tinySpread` (2^-700 case exact:
intercept 1, slope 2^699; 1e-200: intercept 1, slope 5e199), `TestWeightedMean_tinyProducts`,
`TestPooledStandardDeviation_tinySD`.

## 2. LOW: MedianAbsoluteDeviation allocations (FIXED)

- New `num/stats/median_in_place.go`: unexported `medianInPlace` (:7) and a float64 in-place introselect
  `selectFloatInPlace`/`selectFloatRange` (:35, :41), the same algorithm as ALG-SELECT-001 (Musser 1997;
  Blum et al. 1973), O(n) worst case, O(log n) stack. `algo/selection` exports only copying selectors.
- `median_absolute_deviation.go:27–42`: one buffer: copy of values → median in place → overwritten by
  the deviations → median in place. Benchmark: 3 allocs/op before (`out/bench_stats.txt:28–30`), 1 now.
- Tests: `TestMedianAbsoluteDeviation_allocs` (testing.AllocsPerRun == 1),
  `TestMedianAbsoluteDeviation_inputUnchanged`, `TestMedianInPlace_oracle` (400 trials, random/ties/
  sorted/reversed, with and without forced median-of-medians).
- Suggested for later (not done, not owned): an exported in-place selector in `algo/selection` would
  let this file go.

## 3. LOW: doc-comment label order (FIXED)

All 29 exported stats functions now read ID, Time, Space, Source(s), Tags, Units, Example; content
unchanged (a script moved the `Tags:` line only; a checker found 0 blocks out of order).

## 4. NOTE: reserved IDs in doc.go

`doc.go:54–56` now says, in prose outside the function list: "Reserved, not shipped: NUM-STATS-007
GroupedMode is held … and NUM-STATS-027 WeightedGeometricMean is out (D-010 …). Both IDs stay taken and
are never reused."

## Changed index lines (verbatim, for the lead to regenerate INDEX.md)

- NUM-STATS-025 SumOfSquares Time: `O(n)   n = len(values); two passes (four only when a plain sum overflows or every deviation is below 2^-480)`
- NUM-STATS-026 WeightedMean Time: `O(n)   n = len(values); two passes (four only when a plain sum overflows or a product underflows)`
- NUM-STATS-031 MedianAbsoluteDeviation Space: `O(n) extra (one allocation: a buffer holding a copy of values, then the deviations, each selected in place) plus O(log n) stack`

No Source line changed. Other Time/Space lines unchanged.

## Benchmarks (`out/bench_stats_fix.txt`, -benchmem -benchtime=200ms, Apple M1 Pro)

| Benchmark (1e5) | before ns/op (bench_stats.txt) | after ns/op | allocs/op after |
|---|---|---|---|
| SampleVariance | 1439178 | 1434463 | 0 |
| SumOfSquares | 1584363 | 1573245 | 0 |
| PearsonCorrelation | 2201078 | 2181142 | 0 |
| LinearRegression | 2219506 | 2187510 | 0 |
| WeightedMean | 887977 | 945588 | 0 |
| MedianAbsoluteDeviation | 2969410 (3 allocs, 2408451 B) | 2258790 | 1 (802816 B) |

Growth is linear at 1e3/1e4/1e5 for every function (see the file). WeightedMean +6%: the per-term
underflow check.

## Notes (d) status

`notes_stats.md` (d) items 1 (MAD two allocations), 2 (squares below ~1e-154 underflow) and 3
(WeightedMean product underflow) are fixed here; the notes file is not owned by this pass and was not
edited.
