# algo/window — MovingMedian (builder W2-B7, lead decision L7)

Files: algo/window/moving_median.go, moving_median_test.go; one line added to algo/window/doc.go (all other
lines unchanged). Last window ID before this was ALG-WINDOW-014 (doc.go confirmed). Imports algo/heap
(NewHeap, Push, Pop, Peek, Len). Reuses ErrWindowLen and ErrNotFinite; benchSizes/benchFloats/benchWindow
from rolling_sum_test.go.

### ALG-WINDOW-015 MovingMedian (moving_median.go:49)
`func MovingMedian(dst, values []float64, windowLen int) ([]float64, error)`
- **Sources**: Härdle, Steiger, "Algorithm AS 296: Optimal Median Smoothing", JRSS Series C (Applied
  Statistics), 1995, 44(2):258–264, doi:10.2307/2986349 — metadata confirmed on Crossref 2026-11-12 (paper not
  read; that it uses a double heap with O(log w) per step is from its standard description, **lead: confirm
  or keep S&W as the method source**). Sedgewick & Wayne, Algorithms 4th ed., §2.4 Creative Problems,
  "Dynamic median-finding" (max-heap below, min-heap above) — the source ALG-HEAP-010…013 already cite.
- **Even window (decision)**: the arithmetic mean of the two middle values, computed as (a+b)/2 and, only on
  overflow, a/2+b/2 (correctly rounded, never overflows). Lost: lower median (integer-friendly but not the
  statistical median of an even sample) and returning (lower, upper) pairs (two dst slices, awkward API);
  a caller who needs the lower median uses an odd window or MedianStream.
- **Numbers**: float64 only, NaN/±Inf rejected with ErrNotFinite (as RollingSum/RollingMean). Lost: generic
  cmp.Ordered (an even window's mean needs arithmetic; NaN would need D-002's NaN-least order, which makes a
  median of NaN-containing data meaningless).
- **Best known**: O(n log w) (Härdle–Steiger; also indexable skiplist / order-statistic tree). Lower bound
  Ω(n log w) comparisons for comparison-based running medians (cited in Härdle–Steiger's abstract as
  "optimal"; not verified first-hand). O(w) space.
- **Approach**: two heaps of (value, position) — a total order with no equal entries, so one comparison with
  the lower top tells which half a leaving value is in. Leaving values are counted out (liveLower/liveUpper)
  and discarded lazily when they surface at a top. A heap grown by stale entries past 2w+2 is rebuilt from its
  live entries (O(w log w), amortized over ≥ w+1 pushes), so memory is O(w), not O(n). Fixed number of
  allocations (two heaps presized to 2w+3, one rebuild buffer allocated at most once).
  Lost: (1) sort each window, O(n·w log w); (2) a sorted window with binary insert/delete, O(n·w) moves;
  (3) lazy deletion without rebuild — O(n) memory on monotone input; (4) an order-statistic tree (O(log w)
  worst case, but a new balanced tree in algo for one caller).
- **Worked example**: hand-derived span-3 smoother on 2 80 6 3 1 9 4 → 6 6 3 3 4 (the paper's data were not
  available; stated in the test).
- **Tests**: table (empty, window > n, one, w=1 identity, two/even, duplicates, sorted, reversed, negative and
  zero, extremes ±MaxFloat64, subnormal), problem-site shape in our words, dst append, input unmodified, every
  error, oracle (400 trials, heavy ties, sorted inputs forcing rebuilds, n up to 40w) and a float oracle.
  Benchmarks: BenchmarkMovingMedian/1e3,1e4,1e5 (w=50) and BenchmarkMovingMedianWindow/w=10,1e3,1e4 (n=1e5).
- **Errors**: ErrWindowLen (windowLen < 1), ErrNotFinite. dst unchanged on error. values not modified.
- **Tags**: moving-median, running-median, median-filter, two-heaps, lazy-deletion. LeetCode 480.
