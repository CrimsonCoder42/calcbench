# algo/window (SHORT = WINDOW) — S2 additions (builder W1-B2)

ALG-WINDOW-001…004 are S1 (unchanged; algo/window/monotonic_deque.go's unexported appendWindowExtremes, which
RollingMin/RollingMax use, is untouched). New: the exported MonotonicDeque (005–010), three variable-window
functions (011–013) and the exponential moving average (014). doc.go: lines added, existing lines unchanged.

New errors: `ErrNotPositive` (min_length_subarray_sum_at_least.go:9), `ErrNegativeK`
(longest_with_at_most_k_distinct.go:9), `ErrAlpha` (exponential_moving_average.go:11). Reused: `ErrNotFinite`,
`ErrOverflow`.

## MonotonicDeque API (the one home of the monotonic queue; algo/queue must not duplicate it)
```go
type MonotonicDeque[T cmp.Ordered] struct{ /* ring buffer */ }  // zero value = empty min-deque
func NewMonotonicDeque[T cmp.Ordered](wantMax bool) *MonotonicDeque[T]  // ALG-WINDOW-005
func (d *MonotonicDeque[T]) PushBack(v T) int                          // ALG-WINDOW-006 (returns evicted count)
func (d *MonotonicDeque[T]) Front() (T, bool)                          // ALG-WINDOW-007
func (d *MonotonicDeque[T]) PopFront() (T, bool)                       // ALG-WINDOW-008
func (d *MonotonicDeque[T]) PopFrontIfEqual(leaving T) bool            // ALG-WINDOW-009
func (d *MonotonicDeque[T]) Len() int                                  // ALG-WINDOW-010
```
(monotonic_deque.go:82, 97, 110, 151, 167, 189, 204)

### ALG-WINDOW-005…010 MonotonicDeque
- **Source**: Lemire, Nordic Journal of Computing 13(4):328–339, 2006, arXiv:cs/0610046 (S1 source,
  doc/sources/lemire_2006.md). Laaksonen, *Competitive Programmer's Handbook*, 2018, §8.3 "Sliding window
  minimum" pp. 81–82 (read first-hand).
- **Best known**: O(1) amortized push (each element evicted at most once), O(1) front/pop; a window of n
  values costs ≤ 2n comparisons. **Lost**: a heap (O(log w)); a two-stack queue with running minima (also
  O(1) amortized, but two arrays and no "evicted" notion).
- **Decision for the lead (API)**: `cmp.Ordered`, value-only, min/max by a flag, zero value usable. PushBack
  evicts only elements strictly worse than the new one; equal elements are kept (Laaksonen's queue is
  strictly increasing and evicts equals). Keeping equals is what makes `PopFrontIfEqual(leaving)` correct
  without storing indexes (argument in the doc comment; oracle against the naive rescan with heavy ties).
  **Lost**: `T any` + a compare func (zero value unusable, nil func must become an error) and (index, value)
  pairs. If wave 2 needs struct elements (time-stamped points), a func-ordered variant gets a new ID.
  Ring buffer of power-of-two capacity, doubling; a nil receiver is empty for every method but PushBack.
- **Worked example**: Laaksonen §8.3: 2 1 4 5 3 4 1 2, window 4 → queues 1 4 5 / 1 3 / 3 4 / 1 / 1 2, minima
  1 1 3 1 1 (queue contents checked step by step).
- **Consumer evidence**: the consumer backend/knowledgebase/series/moments.go:132–146 (`RollingSum`, a time-based window
  by two pointers) — the same left-pointer pattern a time-windowed min/max would drive with this deque.
- **Tags**: monotonic-deque, monotonic-queue, sliding-window. LeetCode 239, 1438, 862.

### ALG-WINDOW-011 MinLengthSubarraySumAtLeast (min_length_subarray_sum_at_least.go:31)
`func MinLengthSubarraySumAtLeast[T integer](values []T, target T) (start, length int, err error)`
- **Source**: Laaksonen §8.1 "Subarray sum" p. 77 (positive values, two pointers, O(n)).
- **Best known**: Θ(n), O(1) for positive values. **Lost**: prefix sums + binary search (O(n log n), O(n));
  with negative values a monotonic deque over prefix sums is needed (not this function).
- **Approach**: never forms a sum: tracks `need` (target − sum) while short and `spare` (sum − target) once
  reached, both bounded by target or the entering value, so no overflow for any T. Ties → earliest start.
  target ≤ 0 → (0, 1) for non-empty input. No run → (0, 0).
- **Worked example**: Laaksonen's array 1 3 2 5 1 1 2 3 with target 8 → (1, 3) (3+2+5), hand-derived.
- **Errors**: ErrNotPositive (value ≤ 0). **Tags**: two-pointers, variable-window. LeetCode 209.

### ALG-WINDOW-012 LongestWithAtMostKDistinct (longest_with_at_most_k_distinct.go:24)
`func LongestWithAtMostKDistinct[T comparable](values []T, k int) (start, length int, err error)`
- **Source**: Laaksonen §8.1 (one-way pointers, O(n) moves); CLRS ch. 11 (hash map). The specific problem has
  no textbook page found; it is the §8.1 pattern with a count map. **Lead: confirm "pattern + textbook".**
- **Best known**: O(n) expected, O(min(n, k+1)) space. **Lost**: binary search on length (O(n log n));
  O(n²) scans.
- **Approach**: NaN (v != v) counted as its own distinct value without the map (it cannot be deleted from a
  Go map); −0 and +0 are one key. Ties → earliest start. k = 0 → (0, 0).
- **Worked example**: hand-made 1 2 1 3 4 3 5 3, k = 2 → (0, 3). **Errors**: ErrNegativeK.
- **Tags**: two-pointers, distinct, hash-map. LeetCode 340, 159, 904.

### ALG-WINDOW-013 LongestWithoutRepeat (longest_without_repeat.go:16)
`func LongestWithoutRepeat[T comparable](values []T) (start, length int)`
- **Source**: as 012 (Laaksonen §8.1; CLRS ch. 11). Same confirmation needed.
- **Best known**: O(n) expected, O(d) space (d distinct values). **Lost**: count map with a shrinking loop
  (same O(n), more map writes); a fixed 256-entry table (bytes only; strings belong to algo/text).
- **Worked example**: hand-made 4 1 4 2 3 1 2 → (1, 4). NaN never repeats. **Errors**: none.
- **Tags**: two-pointers, distinct, unique. LeetCode 3 (on slices).

### ALG-WINDOW-014 ExponentialMovingAverage (exponential_moving_average.go:29)
`func ExponentialMovingAverage(dst, values []float64, alpha float64) ([]float64, error)`
- **Source**: NIST/SEMATECH e-Handbook §6.4.3.1 "Single Exponential Smoothing",
  https://www.itl.nist.gov/div898/handbook/pmc/section4/pmc431.htm, read 2026-11-12: Sₜ = α·yₜ₋₁ + (1−α)·Sₜ₋₁,
  0 < α ≤ 1, t ≥ 3; initialisation S₂ = y₁ ("one method"; the handbook also mentions a target value or the
  mean of the first 4–5 observations).
- **Index mapping (documented)**: result i = S_{i+2}, i.e. result 0 = values[0] and result i = α·values[i] +
  (1−α)·result[i−1]: the smoothed value after observation i+1, which NIST uses as the forecast for the next
  period. **Decision for the lead**: initialisation S₂ = y₁ only (NIST's first method). Lost: a caller-given
  initial value (can be a later variant with a new ID).
- **Best known**: Θ(n), O(1) extra. **Approach**: α·v + (1−α)·s (a convex combination of finite values, so it
  cannot overflow except by rounding within an ulp of MaxFloat64, which returns ErrOverflow). Lost: s + α(v − s)
  (v − s overflows for opposite extremes; tested).
- **Worked example**: NIST's 12 observations 71 70 69 68 64 65 72 78 75 75 75 70, α = 0.1 → printed S₂…S₁₂
  71, 70.9, 70.71, 70.44, 69.80, 69.32, 69.58, 70.43, 70.88, 71.29, 71.67, checked to ±0.005 (half the
  printed unit).
- **Errors**: ErrAlpha (α ∉ (0,1], NaN), ErrNotFinite, ErrOverflow; dst unchanged on error.
- **Consumer evidence**: none in the consumer backend/knowledgebase/series (trend.go uses split halves, Theil–Sen and
  OLS, no smoothing). Note taxonomy also lists "trend and smoothing of a number series" under num/series;
  the brief placed EMA here. **Lead: confirm the home** (it is a recursive filter, not a window).
- **Tags**: ema, exponential-smoothing, time-series.

## Consumer evidence checked, no replacement
- days.go:62 `BucketFixed` (fixed-stride time buckets with min/mean/max) — a group-by, not a sliding window;
  no function here replaces it.
- glucose.go:11 `TimeInRange` (count per band by binary search over edges) — a histogram; belongs with
  algo/search, not window or prefix.
- fitness.go:269 `RollingUsualRange` (last-w values, Tukey fences) — a single trailing window; no
  replacement here.
