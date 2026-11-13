# algo/stream (SHORT = STREAM)

Builder W4-B15. New package. Not yet in `doc/taxonomy.md` (lead to add: `ALG | algo/stream | STREAM |
one-pass summaries the caller drives: moving average, hit counter, disjoint intervals, running statistics`).
Pure types: the caller passes `now time.Time` (HitCounter) and owns locking; no clock, no goroutine.
num/stats has BATCH variance (two passes over a slice, another subject); RunningStats is the one-pass,
mergeable streaming type: distinct, not a duplicate.

| ID | Name | Signature |
|---|---|---|
| ALG-STREAM-001 | MovingAverageStream | `func NewMovingAverageStream(windowLen int) (*MovingAverageStream, error)`; `Add(value float64) (float64, error)`, `Mean() (float64, bool)`, `Len() int`, `WindowLen() int` |
| ALG-STREAM-002 | HitCounter | `func NewHitCounter(window, resolution time.Duration) (*HitCounter, error)`; `Hit(now time.Time) error`, `Add(now time.Time, hits int64) error`, `Count(now time.Time) (int64, error)` |
| ALG-STREAM-003 | DisjointIntervals | `type DisjointIntervals[T integer]` (zero value ready); `Add(v T)`, `Contains(v T) bool`, `Len() int`, `Intervals(dst []interval.Interval[T]) []interval.Interval[T]` |
| ALG-STREAM-004 | RunningStats | `type RunningStats` (zero value ready); `Add(x float64) error`, `Merge(other RunningStats) error`, `Count() int64`, `Mean`, `PopulationVariance`, `SampleVariance`, `Min`, `Max` `() (float64, bool)` |

One ID per type (L23: the methods are accessors/mutators of one structure).

### ALG-STREAM-001 MovingAverageStream
- **Source**: NIST/SEMATECH e-Handbook of Statistical Methods, §6.4.2.1 "Single Moving Average" (read
  2026-11-12; already cited by ALG-WINDOW-001/002); Neumaier, ZAMM 54(1):39–51, 1974,
  doi:10.1002/zamm.19740540106 (compensated sum); Knuth TAOCP Vol. 1 §2.2.2 via `queue.CircularBuffer`
  (ALG-QUEUE-002, imported, Overwrite policy).
- **Best known**: O(1) per value, O(w) space. Met as O(1) amortized: the running sum is updated by
  subtract-leaving / add-entering (Neumaier), and recomputed from the buffer every w additions (O(w) each), so
  the error depends only on the last two windows, never on stream length (tested over 1e5 values).
  **Lost**: a plain running sum (error grows with stream length: catastrophic after a large value leaves);
  re-summing the window on every Add (O(w) per value).
- **Errors**: ErrWindowLen (w < 1 or w > 1<<24 — the run's 128 MiB silent-allocation cap, L13/L35),
  ErrNotFinite (NaN/±Inf), ErrOverflow (window sum beyond float64; state unchanged), ErrUninitialized (zero
  value). Mean on empty → (0, false). Partial windows: the mean of what has arrived.
- **Worked example**: NIST §6.4.2 supplier data, N = 3 → 26/3, 29/3, … (the handbook's averages).
- **Tags**: LeetCode 346 (moving average from a data stream).

### ALG-STREAM-002 HitCounter
- **Source**: Datar, Gionis, Indyk, Motwani, "Maintaining stream statistics over sliding windows", SIAM J.
  Comput. 31(6):1794–1813, 2002, doi:10.1137/S0097539701398363 (verified on Crossref 2026-11-12: the
  sliding-window model; exact counting over N time units needs Θ(N) bits); Knuth TAOCP Vol. 1 §2.2.2
  (circular array indexed modulo M).
- **Best known**: exact counts at a fixed resolution: O(1) per hit, O(B) per count, O(B) space, B =
  window/resolution buckets; a total maintained incrementally gives O(1) count only when time never goes
  backward and costs O(min(Δ, B)) per time jump anyway. **Lost**: a queue of every hit's timestamp (exact to
  the nanosecond but O(hits) memory: unbounded under load); DGIM exponential histograms (O(log² N) bits but
  approximate counts — the brief asks for a counter, approximation is a different contract).
- **Resolution (documented)**: tick = ⌊UnixNano / resolution⌋ (floor also before 1970). A hit at x counts at
  now iff tick(now) − B < tick(x) ≤ tick(now): the effective window is between window − resolution and
  window. Window 300 s at 1 s: Count at second 301 drops second 1. Buckets are stamped with their tick, so a
  bucket from an older lap reads as empty. Future hits (after now) are not counted. A hit older than the
  newest tick held in its slot is dropped. Times must be in UnixNano's range (years 1678–2262).
- **Errors**: ErrSpan (window or resolution ≤ 0, or resolution not dividing window), ErrBuckets (B > 1<<23:
  two int64 per bucket = 128 MiB, as L13), ErrHits (hits < 0), ErrTimeRange, ErrOverflow (bucket or total
  beyond MaxInt64; state unchanged), ErrUninitialized.
- **Worked example**: window 300 s, resolution 1 s: hits at 1, 2, 3 → Count(4) = 3; 3 hits at 300 →
  Count(300) = 6; Count(301) = 5. Hand-derived.
- **Tags**: LeetCode 362 (design hit counter), rate limiting.

### ALG-STREAM-003 DisjointIntervals
- **Name**: `DisjointIntervals` (practitioner term), not "DataStreamDisjointIntervals" (problem-site
  title; L29). Output reuses `interval.Interval[T]` (ALG-INTVL-001's type).
- **Source**: Knuth TAOCP Vol. 3, 2nd ed., 1998, §6.2.1 Algorithm B (binary search in an ordered table), via
  `search.UpperBound` (ALG-SEARCH-002, imported); Cormen et al., 3rd ed., 2009, ch. 13 (red-black trees, the
  O(log k) alternative).
- **Best known**: O(log k) per Add with a balanced search tree keyed by range start (find the predecessor
  and successor ranges, merge, delete at most one). **Chosen**: two parallel sorted slices (starts, ends):
  O(log k) search + at most one slices.Insert/Delete, O(k) memmove worst case; Contains O(log k);
  Intervals O(k) with no tree walk. **Lost**: a balanced BST — calcbench has none (algo/tree's BST is
  unbalanced and has no predecessor/successor query); for the k a stream of integers produces in practice the
  memmove is a few cache lines per Add, and Intervals is a straight copy. If a consumer needs k ≫ 1e5 ranges
  with heavy churn, a balanced-tree version gets its own ID. Benchmark at 1e3/1e4/1e5 adds shows the k term.
- **Overflow**: adjacency tests are `ends[i-1] == v−1` (only after ends[i-1] < v) and `starts[i]−1 == v`
  (starts[i] > v): no wrap at the extremes of any integer type (tested on int, int8, uint8 extremes).
- **Errors**: none (every integer is valid). Zero value is an empty set.
- **Worked example**: add 1, 3, 7, 2, 6 → [1,3], [6,7]. Hand-derived. Oracle: a presence table.
- **Tags**: LeetCode 352 (data stream as disjoint intervals).

### ALG-STREAM-004 RunningStats
- **Source**: Welford, "Note on a method for calculating corrected sums of squares and products",
  Technometrics 4(3):419–420, 1962, doi:10.1080/00401706.1962.10490022 (verified on Crossref 2026-11-12).
  Chan, Golub, LeVeque, "Algorithms for computing the sample variance: analysis and recommendations", The
  American Statistician 37(3):242–247, 1983, doi:10.1080/00031305.1983.10483115 (verified on Crossref
  2026-11-12; the updating formula for combining two subsamples). Chan, Golub, LeVeque, "Updating formulae and
  a pairwise algorithm for computing sample variances", Stanford CS Technical Report STAN-CS-79-773, November
  1979 (confirmed by its catalog record; DTIC ADA083170; the full text was not readable — 403 — so the
  equation number of the merge formula is not cited; the reviewer should check it against §2 of the 1983
  paper).
- **Formulas**: Welford: δ = x − mean; mean += δ/k; M2 += δ(x − mean_new). Merge: n = nA + nB;
  δ = meanB − meanA; mean = meanA + δ·nB/n; M2 = M2A + M2B + δ²·nA·nB/n. Variances M2/n and M2/(n − 1).
- **Best known**: O(1) per value and per merge, O(1) space. Met. **Lost**: running Σx and Σx² (one pass but
  catastrophic cancellation: the 1e9 + {4,7,13,16} test would lose every digit); keeping the values for a
  two-pass at the end (O(n) memory: that is num/stats).
- **Errors**: ErrNotFinite (NaN/±Inf), ErrOverflow (count beyond MaxInt64; deviation, mean difference or M2
  beyond float64; state unchanged). Accessors return (0, false) when undefined (no values; n < 2 for the
  sample variance). Min/Max use Go's min/max built-ins (−0 < +0).
- **Worked example**: 2 4 4 4 5 5 7 9 → mean 5, population variance 4, sample variance 32/7. Hand-derived.
  Oracle: two-pass on 300 random sets, and Merge of a random split vs one summary.
- **Tags**: online variance, parallel variance, Welford.
