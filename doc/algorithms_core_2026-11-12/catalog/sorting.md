# Catalog: algo/sorting (SHORT = SORT)

Builder W1-B1, 2026-11-12. New package `algo/sorting`, named `sorting` and not `sort` because callers
import the standard `sort` beside ours.

Package conventions (stated once in `doc.go`):
- Every sort writes its result into `dst`: `dst` is reused when `cap(dst) >= len(values)`, otherwise one
  slice is allocated. `dst` may be `values` itself (same first element) to sort in place; any other overlap
  between the two is not supported. `values` is never modified unless it is `dst`.
- Ordered functions compare with `cmp.Compare` (NaN first, -0 equals +0; D-002). Functions that take a
  `compare func(a, b T) int` follow the `slices.SortFunc` contract (negative, zero, positive); a nil compare
  is `ErrNilCompare`.
- Errors are sentinels wrapped with the parameter name.

Sources confirmed 2026-11-12: Sedgewick & Wayne booksite pages algs4.cs.princeton.edu/22mergesort,
/23quicksort, /25applications (quotes below). CLRS 3rd ed. ch. 8 section numbers (8.2 counting, 8.3 radix,
8.4 bucket) and Knuth Vol. 3 §5.2 numbering are from the books' tables of contents, not re-read here; the
reviewer should treat page-level locators as NOT verified against a copy (as S1 did, TODO.md).

---

### ALG-SORT-001 MergeSortStable
`func MergeSortStable[T any](dst, values []T, compare func(a, b T) int) ([]T, error)`
- **Source**: Knuth, TAOCP Vol. 3, 2nd ed., 1998, §5.2.4 (sorting by merging; straight two-way merge).
  Cormen et al., *Introduction to Algorithms*, 3rd ed., MIT Press, 2009, §2.3.1 (MERGE). Sedgewick & Wayne,
  *Algorithms*, 4th ed., 2011, §2.2: "Bottom-up mergesort uses between 1/2 N lg N and N lg N compares".
- **Best known**: O(n log n) time; Ω(n log n) comparisons lower bound for any comparison sort (CLRS §8.1).
  Stable in O(n log n) time needs O(n) extra space for a plain merge (in-place stable merges exist but cost
  more moves).
- **Approach**: bottom-up (iterative, no recursion) merge sort. Runs of `smallRun` = 12 elements are first
  sorted by insertion sort (stable), then merged pairwise, ping-ponging between `dst` and one scratch buffer.
  Ties take the left element, so the sort is stable.
  **Lost**: `slices.SortStableFunc` (insertion blocks + SymMerge, in place: O(n log n) compares but
  O(n log² n) moves); top-down recursive merge sort (same bound, recursion); a stable sort by appending the
  index as a tie key to an unstable sort (extra O(n) and a second compare).
- **Worked example**: Sedgewick & Wayne §2.5 stability example in our own data: records (city, time) sorted
  by city keep time order inside each city. Plus KNUTH16 → its sorted order.
- **Errors**: `ErrNilCompare`. **Input modified**: no (unless dst is values).
- **Consumer evidence**: `the consumer backend's api/handlers/…/feed.go:725-732` `sortItems` sorts with
  `sort.Slice` and an ID tiebreak because "Go's sort is not stable". Note: the ID tiebreak already makes the
  order total and deterministic, so the swap is a simplification (newest first, input order kept for equal
  times) only if the input order is itself meaningful; listed for Devin as a consumer change, never edited.
- **Tags**: sorting, merge-sort, stable, comparison-sort. LeetCode 912, 148.

### ALG-SORT-002 InsertionSort
`func InsertionSort[T cmp.Ordered](dst, values []T) []T`
- **Source**: Knuth, TAOCP Vol. 3, §5.2.1, Algorithm S (straight insertion). CLRS §2.1 (INSERTION-SORT).
- **Best known for its niche**: O(n + I) time, I = number of inversions; O(1) extra. Worst O(n²).
- **When to use (the n bound)**: small inputs and nearly sorted inputs. The crossover against an
  O(n log n) sort is NOT measured in this run. The bound comes from the source: Sedgewick & Wayne §2.3,
  "it pays to switch to insertion sort for tiny arrays. The optimum value of the cutoff is system-dependent,
  but any value between 5 and 15 is likely to work well"; §2.2: insertion sort for small subarrays improves
  mergesort "by 10 to 15 percent". This package uses 12 (inside 5–15) for MergeSortStable and BucketSort.
- **Approach**: straight insertion, shifting (not swapping). **Lost**: binary insertion (fewer compares,
  same moves, worse on nearly sorted input).
- **Worked example**: KNUTH16 → sorted order. **Errors**: none. **Input modified**: no (unless dst is values).
- **Tags**: sorting, insertion-sort, stable, small-n, adaptive.

### ALG-SORT-003 CountingSort
`func CountingSort(dst, values []int, maxKey int) ([]int, error)`
- **Source**: CLRS §8.2 (COUNTING-SORT); Knuth, TAOCP Vol. 3, §5.2, Algorithm D (distribution counting).
- **Best known**: Θ(n + k) time, k = maxKey + 1; O(k) extra. Beats Ω(n log n) because it does not compare.
- **Approach**: count each key, then write each key count times in key order. Plain ints carry no
  satellite data, so the prefix-sum (stable placement) pass of CLRS is not needed.
  **Lost**: the CLRS prefix-sum form (one more pass and an n-length output scan, same result for bare ints).
- **Policy (assumption, for Devin)**: `maxKey` above `maxCountingKey` = 1<<24 is rejected with
  `ErrKeyRange`, so a caller's mistake cannot request a multi-gigabyte count array. Use RadixSortUint64 for
  wide keys. Alternative that lost: no cap (one bad argument can exhaust memory, and make panics on
  overflowing lengths).
- **Worked example**: our own CLRS-style data [2 5 3 0 2 3 0 3], maxKey 5 → [0 0 2 2 3 3 3 5]
  (CLRS Figure 8.2 uses an array of this shape; the numbers here are re-derived by hand, not quoted).
- **Errors**: `ErrKeyRange` (maxKey < 0, maxKey > 1<<24, or a value outside [0, maxKey]); checked before
  dst is written. **Input modified**: no.
- **Tags**: sorting, counting-sort, integer-sort, linear-time, non-comparison.

### ALG-SORT-004 RadixSortUint64
`func RadixSortUint64(dst, values []uint64) []uint64`
- **Source**: CLRS §8.3 (RADIX-SORT, Lemma 8.4: d passes of a stable sort on b-bit keys split into r-bit
  digits take Θ((b/r)(n + 2^r))); Knuth, TAOCP Vol. 3, §5.2.5 (radix sorting).
- **Best known**: O(n·w/r) time for w-bit keys and r-bit digits; here w = 64, r = 8: 8 passes, O(n) for fixed
  w. O(n + 2^r) extra.
- **Approach**: LSD, 8-bit digits, one counting pass builds all 8 histograms, a stable scatter per digit,
  ping-pong between dst and one scratch buffer. A digit on which every key agrees is skipped.
  **Lost**: MSD radix (recursion, better on short common prefixes, worse constant on random keys);
  16-bit digits (4 passes but 65,536-entry tables that thrash L1 for n below ~1e5).
- **Worked example**: CLRS Figure 8.3 idea (sort by least significant digit first) on our own data:
  [329 457 657 839 436 720 355] → [329 355 436 457 657 720 839]; plus 0 and MaxUint64.
- **Errors**: none. **Input modified**: no (unless dst is values).
- **Tags**: sorting, radix-sort, lsd, integer-sort, linear-time.

### ALG-SORT-005 BucketSort
`func BucketSort(dst, values []float64) ([]float64, error)`
- **Source**: CLRS §8.4 (BUCKET-SORT: n buckets, bucket ⌊n·x⌋, insertion-sort each, concatenate;
  expected Θ(n) for uniform input). CLRS Exercise 8.4-2: worst case Θ(n²), fixed by a worst-case
  O(n log n) sort inside each bucket.
- **Best known**: expected O(n) for independent uniform inputs; worst O(n log n) as written here.
- **Approach**: count bucket sizes, prefix sums, scatter into dst (no per-bucket slices), then insertion
  sort each bucket of at most 12 and `slices.Sort` (pdqsort, O(m log m)) for larger ones.
  The index ⌊n·x⌋ can round up to n when x is just below 1; it is clamped to n-1.
  **Lost**: CLRS's linked-list buckets (n small allocations); insertion sort in every bucket (Θ(n²) worst).
- **Worked example**: CLRS Figure 8.4 shape on our own ten values in [0,1).
- **Errors**: `ErrOutOfRange` (a value < 0, ≥ 1, or NaN; -0 is accepted); checked before dst is written.
  **Input modified**: no.
- **Tags**: sorting, bucket-sort, distribution-sort, floats, linear-time-expected.

### ALG-SORT-006 IsSorted
`func IsSorted[T cmp.Ordered](values []T) (sorted bool, firstDescent int)`
- **Source**: Knuth, TAOCP Vol. 3, ch. 5 introduction (sorted order K_p(1) ≤ … ≤ K_p(n)).
- **Best known**: Θ(n) worst (every adjacent pair must be read); stops at the first descent.
- **Approach**: one scan. Differs from `slices.IsSorted` by also returning the first index i with
  values[i] < values[i-1] (len(values) when sorted), which callers use to report where order broke.
- **Errors**: none. **Input modified**: no. **Tags**: sorting, check, sorted, precondition.

### ALG-SORT-007 DedupeSorted
`func DedupeSorted[T cmp.Ordered](dst, sorted []T) []T`
- **Source**: Sedgewick & Wayne §2.5 "Duplicates": "first sort the array, then make a pass through the
  sorted array, taking note of duplicate values that appear consecutively".
- **Best known**: Θ(n) time, O(1) extra. **Approach**: keep the first of each run of cmp.Compare-equal
  values. dst may be sorted itself (write index never passes read index). **Lost**: a hash set (O(n) extra,
  loses order) and `slices.Compact` (in place only, uses ==, so NaNs never compare equal).
- **Worked example**: [1 1 2 3 3 3 4] → [1 2 3 4]. **Errors**: none. **Input modified**: no (unless dst is sorted).
- **Tags**: sorting, dedupe, unique, sorted, distinct. LeetCode 26.

### ALG-SORT-008 PartialSort
`func PartialSort[T cmp.Ordered](dst, values []T, k int) ([]T, error)`
- **Source**: Chambers, "Algorithm 410: Partial sorting", Communications of the ACM, 1971, 14(5):357–358,
  doi:10.1145/362588.362602 (confirmed on Crossref 2026-11-12); Musser 1997 (introselect,
  doi:10.1002/(SICI)1097-024X(199708)27:8<983::AID-SPE117>3.0.CO;2-#).
- **Best known**: O(n + k log k) time; lower bound Ω(n + k log k) (read every element; sort k).
- **Approach**: `selection.KthSmallestInto` (ALG-SELECT-004) puts the k smallest in dst[:k] in O(n) worst
  case, then `slices.Sort(dst[:k])`. dst is the work buffer of length n; the result is dst[:k].
  **Lost**: a bounded max-heap of k (O(n log k) time, O(k) extra); already in the library as TopK
  (ALG-HEAP-001) for callers who need O(k) memory.
- **Worked example**: KNUTH16, k = 5 → [61 87 154 170 275]. **Errors**: `ErrNegativeK`; k > n is clamped
  to n. **Input modified**: no (unless dst is values).
- **Tags**: sorting, partial-sort, selection, introselect, smallest-k. LeetCode 215, 973.

### ALG-SORT-009 SortIndices
`func SortIndices[T any](dst []int, values []T, compare func(a, b T) int) ([]int, error)`
- **Source**: Sedgewick & Wayne §2.5 (pointer sorting: sort references, do not move the data); stability
  by MergeSortStable (ALG-SORT-001) sources.
- **Best known**: O(n log n) time, O(n) extra. **Approach**: fill dst with 0…n-1 and run the same bottom-up
  stable merge sort on the indexes, comparing values[i] with values[j]; equal values keep ascending index.
  **Lost**: sorting (value, index) pairs (2n copies of T, a T may be large).
- **Worked example**: [30 10 20 10] → [1 3 2 0]. **Errors**: `ErrNilCompare`. **Input modified**: no.
- **Consumer evidence**: as ALG-SORT-001 (feed.go:725-732): sorting indexes leaves FeedItem in place.
- **Tags**: sorting, argsort, indirect-sort, stable, permutation.

### ALG-SORT-010 ThreeWayPartition
`func ThreeWayPartition[T cmp.Ordered](dst, values []T, pivot T) (partitioned []T, less, notGreater int)`
- **Source**: Dijkstra, *A Discipline of Programming*, Prentice-Hall, 1976, ch. 14 "The problem of the
  Dutch national flag". Sedgewick & Wayne §2.3 (3-way partitioning, invariant a[lo..lt-1] < v,
  a[lt..i-1] = v, a[gt+1..hi] > v; confirmed on the booksite).
- **Best known**: Θ(n) time, one pass, at most n swaps, O(1) extra. **Approach**: Dijkstra's lt/i/gt
  invariant. **Lost**: Bentley–McIlroy fast 3-way (equal keys parked at both ends then swapped in; fewer
  swaps when duplicates are rare, more code). Not stable.
- **Worked example**: [2 0 2 1 1 0] pivot 1 → blocks [0 0 | 1 1 | 2 2], less = 2, notGreater = 4.
- **Errors**: none. **Input modified**: no (unless dst is values).
- **Note**: `algo/selection/introselect.go` has an unexported `partitionThreeWay` (same algorithm). selection
  cannot import sorting (sorting imports selection for PartialSort), so this is a second implementation;
  reported to the lead.
- **Tags**: sorting, partition, three-way, dutch-national-flag, quicksort. LeetCode 75.
