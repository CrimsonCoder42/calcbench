# math: notes for the lead (A-MATH)

## New IDs (for INDEX.md regeneration)
- NUM-INTEG-007..011 FractionReduce, FractionAdd, FractionSub, FractionMul, FractionDiv: num/integer/fraction.go (new sentinel ErrZeroDenominator in the same file). doc.go list updated.
- NUM-ALGEB-029..033 PercentOf, PercentFromPart, TotalFromPercent, PercentChange, ApplyPercentChange: num/algebra/percent.go. NUM-ALGEB-034 PercentDifference: num/algebra/percent_difference.go. NUM-ALGEB-035 PercentError: num/algebra/percent_error.go. Tests for all seven in num/algebra/percent_test.go. doc.go list updated.
- NUM-LOGS-005..010 ExponentialDecayRemaining, ExponentialDecayInitial, ExponentialDecayHalfLife, ExponentialDecayElapsed, DecayConstantFromHalfLife, MeanLifetimeFromHalfLife: num/logs/exponential_decay.go (new sentinels ErrOutOfRange, ErrOverflow there). doc.go list updated. Taxonomy row for num/logs says "logarithms, exponentials": exponential decay fits; the lead may want the "Holds" text to mention it.
- NUM-ALGEB-036 ProportionFourthTerm: num/algebra/proportion.go (+ proportion_test.go). doc.go list updated.
- NUM-INTEG-012 Divisors: num/integer/divisors.go (+ divisors_test.go). doc.go list updated.
- NUM-ALGEB-037 RoundDecimalsMode and exported type RoundMode (six IEEE 754-2019 constants): num/algebra/round_decimals_mode.go. New exported type = API shape for Devin to confirm (alternative: one function per mode, lost because nine near-identical functions). doc.go list updated (the type is not an ID).
- NUM-SERIES-020..023 ArithmeticSequenceTerm, ArithmeticSeriesSum (num/series/arithmetic_sequence.go), GeometricSequenceTerm, GeometricSeriesSum (num/series/geometric_sequence.go); tests num/series/sequence_test.go. doc.go list updated.
- NUM-STATS-047 SampleSizeProportion, NUM-STATS-048 MarginOfErrorProportion: num/stats/sample_size_proportion.go (+ _test.go). num/stats now imports num/prob (NormalQuantile) for the first time; prob does not import stats, so no cycle. doc.go list updated.
- NUM-STATS-049 GeometricMean: num/stats/geometric_mean.go (+ _test.go). doc.go list updated. Policy question for Devin: with a 0 and negatives present, return 0 (product is 0) or ErrNegative (current)?
- NUM-STATS-050 MeanConfidenceIntervalZ, NUM-STATS-051 MeanConfidenceIntervalT: num/stats/mean_confidence_interval.go (+ _test.go). doc.go list updated.
- Sources to record in doc/sources/ (shared, lead): Knuth TAOCP Vol. 2 §4.5.1; Kirk Contemporary Mathematics §6.1; Marecek Prealgebra 2e §6.2; Törnqvist/Vartia/Vartia 1985; Burden/Faires/Burden Numerical Analysis 10th ed. §1.2; Abramson College Algebra 2e §6.7; Ling/Sanny/Moebs University Physics Vol. 3 §10.3; Rosen Elementary Number Theory 6th ed. §7.2; Abramson College Algebra 2e §9.2–9.4; IEEE 754-2019 §4.3, §9.5; Cochran Sampling Techniques 3rd ed. §4.4; NIST Dataplot "GEOMETRIC MEAN" (geommean.htm); Montgomery & Runger 6th ed. §8.1–8.2; NIST e-Handbook §7.2.2.1.

## Gaps noted, not built
- right-triangle: solving from altitude, area or perimeter not built. circle: radius from circumference or area not built (stdlib one-liners).
- permutation-and-combination: exact nPr/nCr beyond 2⁶⁴ (math/big) not built; ours returns ErrOverflow (log forms NUM-COMBI-002/-004 exist).
- number-sequence: infinite geometric sum a₁/(1 − r), |r| < 1, not built.
- rounding: "half up", "half down", "half to odd" ties: NEEDS SOURCE (not IEEE 754 attributes). Rounding to the nearest 1/2…1/64 not built.
- root-calculator: a signed real n-th root (odd n of a negative x; math.Pow gives NaN) could be one small num/algebra function if Devin wants it.
- binary/hex: only addition of digit strings exists (ALG-INTMTH-014). A radix-general AddDigitStrings(a, b, radix) would cover hex; base conversion left to strconv. Devin to say whether stdlib-only calculators (exponent real powers, log, root, hex, big number) should get calcbench wrappers; I closed them as MATCH-on-what-we-have or SKIPPED.
- fraction-calculator: decimal-to-fraction (e.g. 1.375 → 11/8) and arbitrary-size fractions (math/big) not built.

## Driver gaps
- (session 3) percent-calculator with form `calc`: `sent` {} — the c21–c23 fields are still not found under that form name. The 8-submission limit is reached; nothing more to try from this side.
- (session 2) matrix-calculator: label "Determinant" exists on the page but the 4th-argument match (input value / <button> text) does not find it; fell back to live mode, exit 4. Match other clickable elements by visible text.
- (session 2) volume-calculator and percent-calculator `calc`: with form "" the driver still clicks the first Calculate on the page, not the one that belongs to the filled fields; results stay empty. Pick the button inside the nearest container of the first filled field.
- DRIVER: volume-calculator: eleven unnamed forms; run without a form name returns the result only for the first form (sphere); cone, frustum, ellipsoid, cap came back with empty result text although `sent` was right.
- DRIVER: matrix-calculator (form calcf): run exits 4, the operations are separate script buttons (det, inverse, AB, …); the driver needs a button-name argument.
- DRIVER: percent-calculator form `calc` (fields c21par1/2, c22par1/2, c23par1/2): run returns sent {} and every field missing. The page appears to have the three phrase calculators under one form name `calc` (or several forms sharing it).

## Handoffs
- roman-numeral-converter (from A-OTHER): NEEDS SOURCE, row 39. Unicode 15.0 §22.3 read: no composition rules. Cajori 1928 Vol. 1 pp. 30–37 located but not read (download denied). Devin: accept a read citation or record the standard-form convention as a D-id; home would be algo/intmath (ALG-INTMTH-017 next).
- distance-calculator → A-OTHER (sci/geoscience): Earth surface distance (WGS 84 geodesic, Vincenty 1975 or Karney 2013; spherical wrapper on NUM-TRIG-031 with a cited radius). Their pair 38.8976,−77.0366 → 39.9496,−75.1503 gives 200.0 km / 124.3 mi. See math/checks/distance-calculator.md.
- half-life-calculator: pure exponential decay built in num/logs (NUM-LOGS-005..010); nothing for sci/** unless A-OTHER wants nuclide data (out of scope).

## _delete moves
(none)

## FOUND, NOT FIXED
- (none in owned paths) Tests for PercentDifference and PercentError sit in num/algebra/percent_test.go rather than their own _test files; the lead may want them split. No go test was run (per BRIEF): every new test runs for the first time at VERIFY.
- Processes at hand-back: none of mine alive. PIDs 3340/3459 (pregnancy-calculator run) belong to another agent, not touched.

## Reviewer fixes (2026-11-13)
- NUM-STATS-049 GeometricMean: exp(mean ln x) replaced by a scaled Frexp product with an exact integer exponent sum (num/stats/geometric_mean.go:8–58); Time/Space unchanged, second Source line IEEE 754-2019 §5.3.3; Tags compensated-summation → scaled-product. Test tolerance now relative 2e-15 with its reason, two cases added (subnormal factors, negative odd exponent sum). Spot check (scratch/geometric-mean-recheck.go, go run): 2⁶⁰⁰ case exact, others ≤ 1.6e-16.
- BenchmarkFractionSub, BenchmarkFractionDiv added (num/integer/fraction_test.go:145, :152).
