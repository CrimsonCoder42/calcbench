# numpy gap analysis (2026-11-13)

READ-ONLY. What numpy has that calcbench lacks, routine by routine, and which of it belongs in calcbench.

## SUMMARY

1. **Counts (317 rows):** HAVE 123 (about 52 of them are Go `math` scalars) · PARTIAL 39 · MISSING, IN SCOPE 117 · OUT OF SCOPE 38. [GAPS.md](GAPS.md) de-duplicates these into 142 rows (108 missing, 34 partial) in 16 groups.
2. **Five biggest groups:** num/prob distributions 29 (gamma, beta, t, F, Weibull, quantiles) · NaN family 17 (nanmean, nanpercentile, fmax) · NEW num/signal 17 (fft, rfft, fftfreq, windows, correlate) · num/linalg 13 (nonsymmetric eig, minimum-norm lstsq, polyfit, slogdet) · polynomial arithmetic 12 (polyadd, polydiv, polyder, polyint, polycompanion).
3. **Vec layer:** the evidence supports exported reductions in algo, not a new num/vec. num/linalg already exports the elementwise Vec* set over `[]float64`; the reductions (Sum, Prod, Diff, CumProd, ArgMax, ArgMin) are what is missing.
4. **The duplicated code:** 10 packages each carry their own private compensated sum, `neumaierAdd` is defined 5 times, and a private finiteness check sits in about 25 files. Two of the sum copies are in algo, which cannot import num, so the shared Sum has to live in algo.
5. **Decisions for Devin:** NEW num/signal, num/special, num/orthopoly and num/float · the NaN policy · isclose symmetric or asymmetric · percentile's four extra methods.
6. **Unreachable pages:** only routines.polynomials-poly1d.html (a 404 from a wrong URL in the brief); A3 read the correct page instead.

## Hand loops in calcbench where numpy has one routine (no edits; candidates for a shared layer later)

- **A private compensated sum** (the kernel a shared Sum would hold):
  - ALG-STREAM-001 MovingAverage (moving_average_stream.go:134) and ALG-WINDOW-001 RollingSum (rolling_sum.go:86);
  - NUM-SERIES-003 TrailingTimeWindowSum (trailing_time_window_sum.go:72) and NUM-SERIES-006 FixedStrideBuckets
    (fixed_stride_buckets.go:85, an inline compensated mean);
  - NUM-LINALG-163/164/165 MatNormFrobenius/1/Inf (mat_norms.go:92) and NUM-GEOM-032 PolygonArea (polygon_area.go:10);
  - num/stats, which shares its own private kernel across 7 files (compensated_sum.go), and num/prob, which shares
    one across 9 files (total_probability.go:19);
  - FIN-ECON-004 HerfindahlHirschmanIndex, LIF-FIT-011 TotalVolumeLoad and MED-RESP-051
    PEFDiurnalVariabilityPercent, each inline.
- **Private kernels in num/stats that could become the shared layer:** NUM-STATS-001 ArithmeticMean, -025
  SumOfSquares and -026 WeightedMean already share the package's private kernels (`meanOf` at compensated_sum.go:60,
  `centredSquares` at :130). -013 Midrange and -018 Range scan for min/max by hand.
- **A mean of differences:** NUM-SERIES-009 AverageRealVariability is mean(abs(diff(x))).
- **A Horner loop where NUM-ALGEB-001/002 exist:** FIN-INT-038 NPV, FIN-INT-039 IRR, FIN-INVEST-004, -020 and -021
  (MIRR, profitability index and Macaulay duration; IDs from the files, matched to names by A3).
- **A private normal CDF:** fin/invest/black_scholes_call.go:9 instead of NUM-PROB-048 NormalCDF (A4; the put file
  has not been checked).
- **A private gcd:** num/integer/modular_arithmetic.go:38, num/combinatorics/binomial.go:60 and
  algo/geometry2d/max_points_on_line.go:62. No GCD is exported.
- **An unexported gamma CDF:** num/prob/chi_square_cdf.go:45 `regularizedGamma`, which is all of numpy's gamma
  distribution apart from the quantile.

## Files

- [BRIEF.md](BRIEF.md): the inventory source, the scope rule and the agent split.
- [calcbench_functions.md](calcbench_functions.md): every live exported function in INDEX.md (1,297).
- [math.md](math.md) (A1), [stats.md](stats.md) (A2), [linalg.md](linalg.md) (A3) and
  [fft_special.md](fft_special.md) (A4): one row per numpy routine.
- [GAPS.md](GAPS.md): every MISSING, IN SCOPE row and every PARTIAL row, de-duplicated and grouped by proposed home.

Pages that could not be fetched: only routines.polynomials-poly1d.html (404, a wrong URL in the brief). A3 read the
correct page, routines.polynomials.poly1d.html, instead.
