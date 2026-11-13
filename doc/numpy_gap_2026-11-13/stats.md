# numpy gap: statistics, sorting/searching/counting, sets (agent A2, 2026-11-13)

Pages read (WebFetch, 2026-11-13, numpy v2.5 "stable"):
- https://numpy.org/doc/stable/reference/routines.statistics.html (order statistics, averages and variances, correlating, histograms)
- https://numpy.org/doc/stable/reference/routines.sort.html (sorting, searching, counting)
- https://numpy.org/doc/stable/reference/routines.set.html (making proper sets, boolean operations)

Pages that failed: none.

Calcbench sources checked: `calcbench_functions.md` (grep by meaning), `num/stats/percentile.go:10-100`,
`num/stats/frequency_distribution.go:9-25`, `algo/selection/kth_smallest_into.go` (doc comment),
`algo/hashing/set_intersection.go:3-4`, `algo/sorting/sort_indices.go:18`, `algo/search/lower_bound.go:15`.

Convention: calcbench works on one `[]float64` at a time; a 1-D match is HAVE. numpy's `axis`, `keepdims`, `out`
and `dtype` arguments are ndarray machinery and are not counted as missing. Where a row is PARTIAL the missing
mathematics is named.

**Counts (51 routines): HAVE 16 / PARTIAL 8 / MISSING, IN SCOPE 18 / OUT OF SCOPE 9**

By page: statistics 23 (HAVE 7, PARTIAL 5, MISSING 10, OUT 1) · sorting 7 (HAVE 3, PARTIAL 1, MISSING 1, OUT 2) ·
searching/counting 11 (HAVE 1, MISSING 4, OUT 6) · sets 10 (HAVE 5, PARTIAL 2, MISSING 3).

## Statistics: order statistics

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| ptp | max − min | NUM-STATS-018 Range | HAVE | |
| percentile | q-th percentile, 13 methods, optional weights | NUM-STATS-015 Percentile (all nine Hyndman–Fan types 1–9, `percentile.go:29-37`; default type 7, D-010) | PARTIAL | Missing: (a) numpy's four non-HF methods `lower`, `higher`, `nearest`, `midpoint` (the type-7 position j+γ rounded down/up/to nearest/averaged); (b) `weights=` (numpy supports it for `inverted_cdf` = HF type 1 only); (c) several q in one call (each call re-sorts; a `Percentiles(values, ps)` sorting once is O(n log n + m)). Home num/stats. Source: Hyndman & Fan 1996, Am. Stat. 50(4):361–365, §2–3 (types 1–9); the four extra methods have no HF number and need their own cite or a decision (they are numpy conventions, not in HF). |
| nanpercentile | percentile ignoring NaN | none (Percentile rejects NaN with ErrNotFinite) | MISSING, IN SCOPE | num/stats. Same sources as Percentile; NaN semantics per IEEE 754-2019 §6.2. Recommend one shared policy for the whole nan family (see nanmean). |
| quantile | q-th quantile, q in [0,1] | NUM-STATS-015 Percentile (takes p in [0,1], so it is numpy's quantile) | PARTIAL | Same gaps as percentile. |
| nanquantile | quantile ignoring NaN | none | MISSING, IN SCOPE | num/stats; as nanpercentile. |

## Statistics: averages and variances

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| median | middle value / mean of two middles | NUM-STATS-012 Median (O(n) selection via ALG-SELECT-002) | HAVE | |
| average | weighted mean Σw·x/Σw | NUM-STATS-026 WeightedMean | HAVE | (numpy's `returned=True` sum of weights is a trivial extra) |
| mean | arithmetic mean | NUM-STATS-001 ArithmeticMean; ALG-STREAM-004 RunningStats (streaming) | HAVE | |
| std | standard deviation, divisor N − ddof | NUM-STATS-022 (ddof 0), NUM-STATS-019 (ddof 1); FrequencyWeighted 035/036 | HAVE | Arbitrary ddof not offered; only 0 and 1 have a statistical meaning in practice. |
| var | variance, divisor N − ddof | NUM-STATS-017 (ddof 0), NUM-STATS-021 (ddof 1); ALG-STREAM-004 (Welford) | HAVE | |
| nanmedian | median ignoring NaN | none | MISSING, IN SCOPE | num/stats. |
| nanmean | mean ignoring NaN | none | MISSING, IN SCOPE | num/stats. Decision for Devin: (A) one exported filter (e.g. `FiniteValues(dst, values)`) that callers compose with the existing functions, or (B) a `Nan*` twin per function. Recommend A: one function, no doubling of the API; B loses on surface size. numpy's docs are not the citation; the semantics (skip NaN, all-NaN → NaN/error) is a policy to log in DECISIONS. |
| nanstd | std ignoring NaN | none | MISSING, IN SCOPE | num/stats, as nanmean. |
| nanvar | var ignoring NaN | none | MISSING, IN SCOPE | num/stats, as nanmean. |

## Statistics: correlating

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| corrcoef | Pearson correlation matrix of k variables | NUM-STATS-014 PearsonCorrelation (two variables only) | PARTIAL | Missing the k×k matrix form over `[][]float64`. num/stats `CorrelationMatrix`. Source: Press et al., Numerical Recipes 3rd ed., §14.5 (linear correlation); one-pass co-moments per Chan, Golub & LeVeque, "Algorithms for computing the sample variance", Am. Stat. 37(3), 1983. |
| correlate | discrete cross-correlation of two 1-D sequences (valid/same/full) | none | MISSING, IN SCOPE | num/series (existing, "sequences and series") or NEW `num/signal` SIGNAL if A4 proposes one for FFT/windows: decide once for correlate, convolve, fft. Source: Oppenheim & Schafer, Discrete-Time Signal Processing 3rd ed., §2.3 (convolution sum) and §2.8/Problem 2.x on correlation; Press et al. NR3 §13.2. Direct O(nm); FFT route O((n+m) log(n+m)) belongs with A4. |
| cov | covariance matrix (ddof, fweights, aweights) | none: no covariance at all, not even of two variables | MISSING, IN SCOPE | num/stats `SampleCovariance(x, y)` + `CovarianceMatrix(rows)`. Source: NIST/SEMATECH e-Handbook §6.5.4.1 (sample covariance matrix); numerically, Welford 1962 Technometrics 4(3):419–420 extended to co-moments (Chan, Golub & LeVeque 1983). Weighted forms: Price 1972 / Cochran, Sampling Techniques 3rd ed. — cite before shipping. |

## Statistics: histograms

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| histogram | counts (or density) per bin, bins given or chosen | NUM-STATS-008 FrequencyDistribution (explicit edges; last bin closed; `frequency_distribution.go:9-16`) | PARTIAL | Missing: `density=True` (count/(N·width)); `weights=`; automatic bins (needs histogram_bin_edges). Also a policy difference: calcbench errors on a value outside the edges (ErrOutOfRange), numpy drops it silently — calcbench's choice is deliberate and should stay; a dropping variant would be a new function. Home num/stats. Source: Scott, Multivariate Density Estimation, 2nd ed., Wiley 2015, §3.1 (histogram density). |
| histogram2d | 2-D bin counts of (x, y) pairs | none | MISSING, IN SCOPE | num/stats `FrequencyDistribution2D(x, y, xEdges, yEdges) [][]int`. Source: Scott 2015 §3.4 (bivariate histogram). |
| histogramdd | d-dimensional bin counts | none | OUT OF SCOPE | Result is a d-dimensional array; `[][]float64` stops at 2 dims. |
| bincount | count (or weight sum) of each non-negative int value 0..max | ALG-HASH-001 FrequencyCount (map, any comparable); counting inside ALG-SORT-003 | PARTIAL | Missing the dense `[]int` (and weighted `[]float64`) result indexed by value, with `minlength`. algo/hashing or algo/sorting `CountByValue`. Source: Knuth TAOCP Vol. 3 2nd ed. §5.2 Algorithm D (distribution counting). |
| histogram_bin_edges | bin edges from a rule (sturges, scott, fd, doane, rice, sqrt, stone, auto) | none | MISSING, IN SCOPE | num/stats `BinEdges(values, rule)` / `BinCount*`. Sources: Sturges, JASA 21(153):65–66, 1926; Scott, Biometrika 66(3):605–610, 1979; Freedman & Diaconis, Z. Wahrsch. 57:453–476, 1981; Doane, Am. Stat. 30(4):181–183, 1976. "rice", "sqrt", "stone" and numpy's "auto" need their own cites or are left out. |
| digitize | bin index of each value for given edges | ALG-SEARCH-002 UpperBound (right=False) / ALG-SEARCH-001 LowerBound (right=True), one value per call | HAVE | Slice form is the vec-layer question. |

## Sorting

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| sort | sorted copy | Go stdlib: slices.Sort / slices.SortStableFunc; ALG-SORT-001 MergeSortStable, 003/004/005 specialised sorts | HAVE (Go stdlib: slices.Sort) | |
| lexsort | indirect stable sort by several key arrays (last key primary) | ALG-SORT-009 SortIndices (stable, compare func) can do it with a composite compare | PARTIAL | No keys-`[][]float64` form. algo/sorting `LexSortIndices(keys)`. Source: Knuth TAOCP Vol. 3 2nd ed. §5.2 (lexicographic order, stable multi-key sorting; LSD via stability). |
| argsort | permutation that sorts | ALG-SORT-009 SortIndices (stable) | HAVE | |
| ndarray.sort | in-place sort method | Go stdlib: slices.Sort | OUT OF SCOPE | ndarray method. |
| sort_complex | sort complex by real then imaginary | none | OUT OF SCOPE | complex dtype; calcbench is `[]float64`. |
| partition | copy with k-th element in place, smaller before, larger after | ALG-SELECT-004 KthSmallestInto (leaves work[:k] = k smallest, work[k-1] = k-th); ALG-SORT-010 ThreeWayPartition | HAVE | |
| argpartition | indices of a partition | none | MISSING, IN SCOPE | algo/selection `KthSmallestIndex` / `ArgPartition`. Source: CLRS 3rd ed. §9.2–9.3; Musser 1997 (introselect), as ALG-SELECT-004. |

## Searching and counting

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| argmax | index of the maximum | none (Range/RunningStats return values, not indices) | MISSING, IN SCOPE | algo/selection `ArgMax` (first index on ties, as numpy). Source: Knuth TAOCP Vol. 1 3rd ed. §1.2.10 Algorithm M (find the maximum). |
| nanargmax | argmax ignoring NaN | none | MISSING, IN SCOPE | algo/selection; same source, NaN policy per the nan-family decision. |
| argmin | index of the minimum | none | MISSING, IN SCOPE | algo/selection `ArgMin`; Knuth Vol. 1 §1.2.10 (mirror of Algorithm M). |
| nanargmin | argmin ignoring NaN | none | MISSING, IN SCOPE | as nanargmax. |
| argwhere | indices of non-zero elements, grouped | none | OUT OF SCOPE | boolean-mask machinery; a Go loop. |
| nonzero | indices of non-zero elements | none | OUT OF SCOPE | as argwhere. |
| flatnonzero | non-zero indices of flattened array | none | OUT OF SCOPE | as argwhere (flattening is ndarray). |
| where | elementwise select x or y by condition | none | OUT OF SCOPE | ufunc/broadcast select. |
| searchsorted | insertion index keeping order (left/right) | ALG-SEARCH-001 LowerBound (side=left), ALG-SEARCH-002 UpperBound (side=right); Go stdlib slices.BinarySearch | HAVE | `sorter=` argument missing (compose with SortIndices); slice-of-targets form is the vec-layer question. |
| extract | elements where a condition holds | none | OUT OF SCOPE | mask; Go slices.DeleteFunc on a copy. |
| count_nonzero | number of non-zero elements | none | OUT OF SCOPE | mask count; no mathematics. |

## Sets

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| unique | sorted distinct values (+ optional index, inverse, counts) | slices.Sort + ALG-SORT-007 DedupeSorted | HAVE | Optional outputs: see unique_all. NaN: numpy collapses all NaNs into one (`equal_nan=True`); calcbench's hashing sets treat NaN as in no set (`set_intersection.go:4`); DedupeSorted's NaN behaviour should be stated. |
| unique_all | distinct values + first indices + inverse + counts | none combined | PARTIAL | Missing first-occurrence indices and inverse. algo/sorting `UniqueAll`. Source: Knuth TAOCP Vol. 3 §5.2 (sorting then scanning equal runs); stable sort keeps first occurrence. |
| unique_counts | distinct values + counts | ALG-HASH-001 FrequencyCount (unordered map) | PARTIAL | Missing the sorted values/counts pair of slices (deterministic order). algo/sorting `UniqueCounts`; Knuth Vol. 3 §5.2. |
| unique_inverse | distinct values + index into them for each input | none | MISSING, IN SCOPE | algo/sorting `UniqueInverse` (the "factorize"/label-encode step). Source: Knuth TAOCP Vol. 3 §5.2 (via SortIndices, ALG-SORT-009). |
| unique_values | distinct values | ALG-SORT-007 DedupeSorted (after sort); ALG-HASH-007 SetUnion(a, nil) (first-appearance order) | HAVE | |
| intersect1d | sorted common values | ALG-HASH-006 SetIntersection (first-appearance order, not sorted) | HAVE | `return_indices` not offered. |
| isin | membership mask of each element in a test set | none (Go stdlib slices.Contains is per element O(m); ALG-BLOOM-001 is approximate) | MISSING, IN SCOPE | algo/hashing `IsIn(values, test) []bool`, O(n + m). Source: Knuth TAOCP Vol. 3 2nd ed. §6.4 (hashing). |
| setdiff1d | values of a not in b | ALG-HASH-008 SetDifference | HAVE | |
| setxor1d | values in exactly one of a, b | none | MISSING, IN SCOPE | algo/hashing `SetSymmetricDifference`; Knuth TAOCP Vol. 3 §6.4 (hashing). Composable today from two SetDifference calls + SetUnion. |
| union1d | sorted union | ALG-HASH-007 SetUnion (first-appearance order) | HAVE | |

## Out of scope, listed once

ndarray.sort (method), sort_complex (complex dtype), histogramdd (d-dimensional result), argwhere / nonzero /
flatnonzero / where / extract / count_nonzero (boolean-mask and broadcast machinery), and every `axis=`, `keepdims=`,
`out=`, `dtype=`, `kind=` / `order=` argument on the in-scope routines.

## Hand loops in calcbench

- `num/series/fixed_stride_buckets.go:85` and `:91`: a compensated mean of a bucket computed inline
  (`(sum + comp) / float64(len(run))`) instead of calling NUM-STATS-001 ArithmeticMean (if ArithmeticMean is
  compensated, this is a candidate to reuse it; if not, ArithmeticMean is the one to upgrade).
- `algo/stream` RunningStats (ALG-STREAM-004) and num/stats variance (NUM-STATS-017/021) are two implementations of
  the same variance; not a defect (streaming vs batch) but worth one shared test oracle.
- No domain package (fin, life, sci, eng) was found computing a mean, variance, percentile or correlation of a
  slice by hand (grep of `/ float64(len(` outside num/stats found only the num/series lines above). The many
  `med/*` score functions loop to sum item scores; that is a questionnaire total, not a numpy-shaped gap.

## Top gaps in scope (for the lead)

1. cov (covariance of two variables and the matrix): calcbench has none at all.
2. argmax / argmin (and nan variants): a basic reduction with no home today.
3. nan family (nanmean, nanvar, nanstd, nanmedian, nanpercentile, nanquantile): one policy decision unlocks six rows.
4. histogram_bin_edges (Sturges, Scott, Freedman–Diaconis, Doane): needed before histogram can choose its own bins.
5. correlate (and its sibling convolve on A1/A4's pages): needs a home decision (num/series vs NEW num/signal).
