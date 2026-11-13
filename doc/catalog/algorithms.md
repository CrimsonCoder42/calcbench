# Catalog: algorithms — the core set (S2), 2026-11-12

One entry per function in the S1 shape, written by the builder of each package (sections below, in the
taxonomy's order) and assembled by the lead. Lead decisions L1–L40 with the alternatives that lost:
[../algorithms_core_2026-11-12/README.md](../algorithms_core_2026-11-12/README.md). The code's doc comments
are the source of truth; where a section and the code differ, the code wins.

DROPPED from the request's LeetCode-derived scope, one line each: sudoku and n-queens (constraint search; no
reusable primitive beyond enumerate's backtracking) · text justification and zigzag conversion (formatting
puzzles, no source) · SQL and shell problems (not Go library material) · design-only problems (APIs, not
algorithms) · brainteasers and game theory (one-off puzzles) · SimplifyPath (= path.Clean) · StrStr
(= strings.Index) · SqrtNewton (= math.Sqrt) · single math/bits calls (L30).

HELD (no acceptable source): MedianOfTwoSorted · MaxProductSubarray · ContainerMostWater · TrappingRainWater ·
NumberOfDigitOne · IntegerToRoman / RomanToInteger · plain CompareVersions · DecodeWays · MaxProfit with k
trades · MaxProfit with cooldown · MaximalSquare · SpiralOrder · SetZeroesInPlace · PerfectRectangle ·
SlidingWindowCounter · Bloom intersection. ROUTED (L4): GCD/LCM, ExtendedGCD, ModPow, ModInverse,
SieveOfEratosthenes, IsPrimeDeterministic64 → num/integer.

The code-heim channel patterns (errgroup, or-done, tee, generators, fan-out/fan-in) are NOT library material:
they own goroutines (see the seq section).


---

# Catalog: algo/search (SHORT = SEARCH), S2 additions

Builder W1-B1, 2026-11-12. Continues after ALG-SEARCH-004. cmp.Compare order; input never modified;
preconditions (sortedness, rotation) are not checked, but nothing panics on input that breaks them.

Consumer evidence for the existing ALG-SEARCH-001/002: `red_lite_cache/internal/cache/ordered.go:45`
hand-rolls `sort.Search(len(n.keys), func(i int) bool { return n.keys[i] > k })` = UpperBound
(ALG-SEARCH-002); `ordered.go:64,118,201` call `sort.SearchStrings(n.keys, k)` = LowerBound
(ALG-SEARCH-001). Consumer change, listed for Devin, not edited.

### ALG-SEARCH-005 SearchRotated
`func SearchRotated[T cmp.Ordered](rotated []T, target T) (index int, found bool)`
- **Source**: the reduction to two binary searches. Knuth, TAOCP Vol. 3, 2nd ed., §6.2.1 (binary search);
  Bentley, *Programming Pearls*, 2nd ed., Column 9.3 (search for the first index where a monotone test
  turns true). For distinct keys rotated by any r, p(i) = rotated[i] < rotated[0] is false then true, so
  the rotation point is the first true index (SearchPredicate, ALG-SEARCH-008); the target then lies in
  rotated[:p] when target ≥ rotated[0], else in rotated[p:], each ascending (LowerBound, ALG-SEARCH-001).
  No textbook section names "search in a rotated array"; the method is this composition of cited steps.
- **Best known**: O(log n), Ω(log n) (Knuth §6.2.1 comparison-tree bound applies to the unrotated case).
- **Approach**: as above, two binary searches. **Lost**: the one-loop variant that decides at each mid
  which half is sorted (same bound, harder invariant, same compares within a constant).
- **Precondition**: distinct keys (a strictly ascending slice rotated). Duplicates: SearchRotatedDup.
- **Worked example**: [40 50 60 10 20 30], 20 → (4, true); 35 → (-1, false).
- **Errors**: none; not found is (-1, false). **Tags**: search, binary-search, rotated, cyclic-shift. LeetCode 33, 153.

### ALG-SEARCH-006 SearchRotatedDup
`func SearchRotatedDup[T cmp.Ordered](rotated []T, target T) (index int, found bool)`
- **Source**: as ALG-SEARCH-005. With duplicates, rotated[i] < rotated[0] stops being monotone only
  through elements equal to rotated[0] at the tail; removing that tail restores it.
- **Best known**: O(n) worst case: [1 1 … 1 0 1 … 1] hides the one different key at any position, so any
  algorithm must read Ω(n) entries in the worst case (adversary argument). O(log n) when
  rotated[0] ≠ rotated[n-1].
- **Approach**: if target equals rotated[0], answer 0. Else trim the trailing run equal to rotated[0]
  (linear, the worst-case term), then SearchRotated's two binary searches on what is left.
  **Lost**: shrink-both-ends inside the loop (same worst case, harder to prove).
- **Worked example**: [2 2 2 3 4 2 2], 4 → (4, true); [1 1 1 0 1 1], 0 → (3, true).
- **Errors**: none; not found is (-1, false). **Tags**: search, binary-search, rotated, duplicates. LeetCode 81, 154.

### ALG-SEARCH-007 PeakElement
`func PeakElement[T cmp.Ordered](values []T) (index int, ok bool)`
- **Source**: Kleinberg, Tardos, *Algorithm Design*, Addison-Wesley, 2006, ch. 5, Solved Exercise 1 (peak
  of a unimodal array with O(log n) probes: if the next entry is larger, the peak is to the right; confirmed
  by search summary 2026-11-12). Kiefer, "Sequential minimax search for a maximum", Proceedings of the AMS,
  1953, 4(3):502–506, doi:10.1090/S0002-9939-1953-0055639-3 (confirmed on Crossref 2026-11-12).
- **Best known**: O(log n) probes; Ω(log n) for unimodal input.
- **Approach**: binary search on p(i) = values[i] > values[i+1] over [0, n-1). Returns the first index the
  search settles on with values[i] ≥ values[i-1] and values[i] ≥ values[i+1] (outside the slice = -∞).
  Proof: the invariant "values[lo-1] < values[lo] (or lo = 0) and values[hi] > values[hi+1] (or hi = n-1)"
  keeps a local maximum inside [lo, hi]. On unimodal input it is the maximum; with adjacent distinct
  values it is a strict local maximum.
  **Lost**: linear scan (O(n)).
- **Worked example**: [1 3 7 9 6 2] → 3. **Errors**: none; empty input is (0, false).
- **Tags**: search, binary-search, peak, local-maximum, unimodal. LeetCode 162, 852.

### ALG-SEARCH-008 SearchPredicate
`func SearchPredicate(lo, hi int, pred func(int) bool) (int, error)`
- **Source**: Knuth §6.2.1 and Bentley Column 9.3 (as ALG-SEARCH-001), generalised from "element ≥ target"
  to any monotone predicate: "binary search on the answer". Same form as Go's `sort.Search` (standard
  library docs), with explicit bounds.
- **Best known**: ⌈log2(hi-lo+1)⌉ predicate calls, O(1) extra.
- **Approach**: half-open bisection; midpoint lo + int(uint(hi-lo)/2) so lo = MinInt, hi = MaxInt does not
  overflow. Returns hi when pred is false everywhere.
- **Errors**: `ErrBounds` (lo > hi), `ErrNilPredicate`. **Worked example**: smallest x in [0, 100) with
  x·x ≥ 50 → 8. **Tags**: search, binary-search, predicate, monotone, answer-search. LeetCode 875, 1011.

### ALG-SEARCH-009 SearchSortedMatrix · ALG-SEARCH-010 SearchSortedMatrixSaddleback
`func SearchSortedMatrix[T cmp.Ordered](matrix [][]T, target T) (row, col int, found bool)`
`func SearchSortedMatrixSaddleback[T cmp.Ordered](matrix [][]T, target T) (row, col int, found bool)`
- **SearchSortedMatrix** (rows ascending and each row's first element > previous row's last: the matrix
  read row by row is one sorted list). Source: Knuth §6.2.1 binary search. O(log m + log n) = O(log(mn)):
  binary search over rows on the first element (LowerBound-style), then LowerBound (ALG-SEARCH-001) in the
  row. Ragged rows are fine; an empty row breaks the precondition (result unspecified, no panic).
  **Lost**: flat index i → (i / cols, i % cols) (needs equal-length rows, a division per probe).
- **SearchSortedMatrixSaddleback** (each row ascending and each column ascending). Source: Gries, *The
  Science of Programming*, Springer, 1981 (saddleback search, an exercise in the loop-development part;
  exact page NOT confirmed); Bird, "Improving saddleback search: a lesson in algorithm design",
  Mathematics of Program Construction (MPC 2006), LNCS 4014, Springer, 2006, pp. 82–89,
  doi:10.1007/11783596_8 (confirmed on Crossref 2026-11-12; the brief's _7 is Backhouse & Michaelis,
  pp. 69–81, a different paper). O(m + n), start top-right, step left or down. Bird gives the lower bound
  Ω(m log(n/m)) for m ≤ n and a divide-and-conquer that meets it.
  **Lost** (for now): Bird's O(m log(n/m)) method (better when m ≪ n; more code, recursion); kept the
  staircase as the brief asks.
  Rows must have equal length; a shorter row is skipped safely (no panic).
- **Worked examples**: 3×4 [[1 3 5 7] [10 11 16 20] [23 30 34 60]], 16 → (1, 2, true), 13 → not found.
  Saddleback 4×4 [[1 4 7 11] [2 5 8 12] [3 6 9 16] [10 13 14 17]], 5 → (1, 1, true), 15 → not found.
  (Our own numbers.)
- **Errors**: none; not found is (-1, -1, false). **Tags**: search, binary-search, matrix, saddleback,
  staircase. LeetCode 74, 240.

---

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

---

# Catalog: algo/selection (SHORT = SELECT), S2 additions

Builder W1-B1, 2026-11-12. Continues after ALG-SELECT-002. Conventions of the package doc.go apply
(cmp.Compare order, D-002; introselect, D-001).

### ALG-SELECT-003 KthLargest
`func KthLargest[T cmp.Ordered](values []T, k int) (T, error)`
- **Source**: as ALG-SELECT-001 (Hoare 1961; Blum et al. 1973; Musser 1997; CLRS §9.2–9.3). The k-th
  largest of n is the (n-k+1)-th smallest (CLRS ch. 9 introduction, order statistics).
- **Best known**: O(n) worst case; Ω(n). **Approach**: validate k against 1..n, then call KthSmallest with
  n-k+1. **Lost**: a size-k min-heap (O(n log k)); that is TopK (ALG-HEAP-001) when all k are wanted.
- **Worked example**: KNUTH16, k = 1 → 908, k = 3 → 765, k = 16 → 61 (read off the sorted keys).
- **Errors**: `ErrEmpty`, `ErrRank` (k < 1 or k > n). **Input modified**: no (a copy is reordered).
- **Tags**: selection, order-statistic, kth-largest, introselect. LeetCode 215.

### MedianOfTwoSorted (no ID): HELD (no confirmed primary source for the method)
`func MedianOfTwoSorted[T cmp.Ordered](a, b []T) (lower, upper T, err error)` (proposed, NOT built)
- **Method** (proposed): binary search on the cut i in the shorter array, j = ⌈(m+n)/2⌉ - i in the longer,
  until a[i-1] ≤ b[j] and b[j-1] ≤ a[i]; O(log min(m, n)) time, O(1) extra.
- **Why held**: the brief requires the source of the partition method. Candidates checked 2026-11-12:
  - Frederickson, Johnson, "The complexity of selection and ranking in X+Y and matrices with sorted
    columns", JCSS 24(2):197–208, 1982, doi:10.1016/0022-0000(82)90048-4 (DOI confirmed on Crossref).
    Two sorted arrays are a two-column matrix with sorted columns, and a secondary summary credits the
    paper with O(m + p·log(k/p)) selection, p = min(k, m), i.e. O(log k) for two arrays. The paper itself
    was not readable here (paywall), so neither the bound nor the partition method is confirmed from it,
    and O(log k) is not the O(log min(m, n)) asked for.
  - CLRS 3rd ed., Exercise 9.3-8: median of two sorted arrays of EQUAL length n in O(lg n). An exercise,
    no method given, equal lengths only.
  - Everything else found was LeetCode 4 and blogs (not acceptable).
- **To unblock**: Devin (or a reader with library access) confirms Frederickson–Johnson 1982 §2–3, or
  accepts "binary search on a monotone predicate (Knuth §6.2.1) over the cut position" as the cited
  reduction. It takes the next free ALG-SELECT ID when built (none is reserved).
- **Tags** (proposed): selection, median, two-sorted, binary-search, partition. LeetCode 4.

### ALG-SELECT-004 KthSmallestInto
`func KthSmallestInto[T cmp.Ordered](dst, values []T, k int) (kth T, work []T, err error)`
- **Source**: as ALG-SELECT-001.
- **Best known**: O(n) worst case; O(1) extra when the caller supplies the buffer.
- **Approach**: copy values into dst (grown only when cap(dst) < n), then call the package's unexported
  `selectInPlace` (algo/selection/introselect.go), the same introselect KthSmallest uses: called, not
  copied. work = dst[:n] reordered so that work[:k] are the k smallest and work[k-1] is the k-th.
  dst may be values itself to select in place. Used by PartialSort (ALG-SORT-008).
  **Lost**: an `InPlace` variant only (forces callers to clone); a pool (a global, forbidden).
- **Worked example**: KNUTH16 into a 16-cap buffer, k = 8 → 509, 0 allocations (AllocsPerRun).
- **Errors**: `ErrEmpty`, `ErrRank`, as KthSmallest; on error dst is returned unchanged (not written).
- **Input modified**: no (unless dst is values).
- **Tags**: selection, order-statistic, quickselect, introselect, no-alloc, dst.

---

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

---

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

---

# algo/prefix (SHORT = PREFIX) — S2, new package (builder W1-B2)

Package in `doc/taxonomy.md` ("prefix sums and difference arrays"). Integer element types only (any built-in
signed or unsigned type ≤ 64 bits, unexported constraint `integer`); floats lose small values to
cancellation in a prefix difference (S1 spec, RollingSum "Lost"), and algo/window already has a compensated
float sum. Ranges are half-open. Overflow policy: results that must fit T return `ErrOverflow` when they do
not; intermediates that need not fit T are carried exactly in a private two's-complement `int128`
(`algo/prefix/int128.go`), so errors are exact (no false alarms).

Errors: `ErrOverflow`, `ErrRange`, `ErrRagged`, `ErrEmpty` (algo/prefix/prefix_sums.go:11–20).

## Sources (verification state)
- **Laaksonen**, *Competitive Programmer's Handbook*, draft of 3 July 2018 (author's free edition of the
  Springer *Guide to Competitive Programming*), https://cses.fi/book/book.pdf, read first-hand 2026-11-12:
  §2.4 "Maximum subarray sum" pp. 21–23 (Algorithm 3, the example −1 2 4 −3 5 2 −5 2 → 10, footnote
  attributing it to Kadane via Bentley); §9.1 "Static array queries" p. 84 (prefix sum array 1 3 4 8 6 1 4 2
  → 1 4 8 16 22 23 27 29; sum(3,6) = 27 − 8 = 19) and p. 85 (2D prefix sums); §9.4 "Range updates" p. 93
  (difference array 3 3 1 1 1 5 2 2 → 3 0 −2 0 0 4 −3 0; range increase by changing two entries).
  **Lead: confirm a competitive-programming textbook is an acceptable primary source** (not in the
  BRIEF_COMMON list; it is a textbook, not a problem site). Every function below also has an independent
  classic source.
- **Blelloch**, "Prefix Sums and Their Applications", in Reif (ed.), *Synthesis of Parallel Algorithms*,
  Morgan Kaufmann, 1993, ch. 1, pp. 35–60; tech report CMU-CS-90-190, Nov. 1990. Read first-hand (PDF,
  cs.cmu.edu/~guyb/papers/Ble93.pdf): §1.1 defines the scan (any binary associative ⊕) and the prescan
  (exclusive, ⊕ with identity); §1.2 uses a max-prescan.
- **Crow** 1984, doi:10.1145/964965.808600 — Crossref-confirmed: ACM SIGGRAPH Computer Graphics 18(3):207–212.
- **Viola & Jones** 2004, doi:10.1023/B:VISI.0000013087.49260.fb — Crossref-confirmed: IJCV 57(2):137–154.
  Section §2.1 "Integral Image" cited from the paper's structure, not re-read here.
- **Bentley** CACM 27(9):865–873, 1984, doi:10.1145/358234.381162 — Crossref-confirmed ("Programming pearls:
  algorithm design techniques"). Programming Pearls 2nd ed., Column 8 "Algorithm Design Techniques"
  (title confirmed from the O'Reilly table of contents; it was Column 7 in the 1st ed., where §7.2 is "Two
  Quadratic Algorithms"); §8.2 (cumarr) and §8.4 "A Scanning Algorithm" are from the table of contents, page
  not read. The brief's "column 7" is the 1st edition; the 2nd edition is Column 8.
- **CLRS** 3rd ed., ch. 11 (hash tables, O(1) expected) — standard; not re-read.

### ALG-PREFIX-001 PrefixSums (prefix_sums.go:35)
`func PrefixSums[T integer](dst, values []T) ([]T, error)`
- **Source**: Laaksonen §9.1 p. 84; Blelloch §1.1 (prescan + total); Bentley Column 8 §8.2 cumarr.
- **Best known**: Θ(n) (must read every value). O(1) extra. **Approach**: one running checked add, append to
  dst (D-003). **Lost**: inclusive n-element output (needs a special case for lo = 0 in every range query);
  wrapping arithmetic with a final check (cannot tell which prefix overflowed).
- **Worked example**: Laaksonen p. 84 → 0 1 4 8 16 22 23 27 29. **Errors**: ErrOverflow (index in message),
  dst unchanged. **Input modified**: no.
- **Consumer evidence**: the consumer backend/knowledgebase/series/moments.go:119–130 (`NewPrefixSums`, `RangeSum`,
  float64, plain `+`/`−`). Not a drop-in: that one is float; it is the same pattern and would use
  PrefixSums/RangeSum on integer counts or cents.
- **Tags**: prefix-sum, cumulative-sum, scan. LeetCode 303, 1480.

### ALG-PREFIX-002 RangeSum (range_sum.go:17)
`func RangeSum[T integer](prefix []T, lo, hi int) (T, error)`
- **Source**: Laaksonen §9.1 p. 84 (sum(a,b) = sum(0,b) − sum(0,a−1)); Bentley §8.2.
- **Best known**: O(1) per query after O(n) preprocessing (optimal). **Lost**: storing values and summing
  (O(hi−lo)); Fenwick tree (O(log n), only needed with updates).
- **Worked example**: Laaksonen sum(3,6) = 19 (our [3,7)). **Errors**: ErrRange; ErrOverflow (two in-range
  prefixes can differ by up to twice the range). **Input modified**: no.
- **Tags**: prefix-sum, range-sum. LeetCode 303.

### ALG-PREFIX-003 PrefixSums2D · 004 SummedAreaTable.RectangleSum · 013 Rows · 014 Cols (summed_area_table.go)
`func PrefixSums2D[T integer](grid [][]T) (SummedAreaTable[T], error)` ·
`func (t SummedAreaTable[T]) RectangleSum(row0, col0, row1, col1 int) (T, error)` · `Rows() int` · `Cols() int`
- **Source**: Crow 1984 (summed-area table); Viola & Jones 2004 §2.1 (integral image, a rectangle from four
  references); Laaksonen §9.1 pp. 84–85.
- **Best known**: Θ(r·c) build, O(1) query, Θ(r·c) table. **Approach**: S(r+1,c+1) = S(r,c+1) + rowprefix(r,c),
  so every checked addition is a true rectangle or row-prefix sum and an overflow means a real sum does not
  fit. The query combines four entries in int128. **Lost**: the textbook four-term recurrence
  S = g + S(up) + S(left) − S(diag) (an intermediate can overflow when the result fits); a [][]T table (r+1
  allocations; ours is one flat slice).
- **Worked example**: none numeric in the sources; hand-derived 3×4 grid 1…12 (full table and 7 rectangles in
  the test). **Errors**: ErrRagged, ErrOverflow (build); ErrRange, ErrOverflow (query). The zero
  SummedAreaTable is a 0×0 table. **Input modified**: no.
- **Tags**: summed-area-table, integral-image, 2d. LeetCode 304.
- **Decision for the lead (API)**: a value type with methods (Rows/Cols/RectangleSum, each with an ID, as
  ALG-HEAP gives methods IDs). Lost: returning the raw (r+1)×(c+1) slice and a free function
  `RangeSum2D(table, …)` (callers would have to know the layout).

### ALG-PREFIX-005 NewDifferenceArray · 006 RangeAdd · 007 Materialize · 012 Len (difference_array.go)
`func NewDifferenceArray[T integer](initial []T) *DifferenceArray[T]` ·
`func (d *DifferenceArray[T]) RangeAdd(lo, hi int, delta T) error` ·
`func (d *DifferenceArray[T]) Materialize(dst []T) ([]T, error)` · `Len() int`
- **Source**: Laaksonen §9.4 p. 93; Blelloch §1.1 (the scan inverts the difference).
- **Best known**: O(1) per range add, Θ(n) to read all values back (optimal for offline updates). **Lost**:
  Fenwick tree with range update (O(log n) each; only needed when reads interleave with updates);
  naive per-element add (O(hi−lo)).
- **Approach**: differences kept in int128 (exact for < 2⁶³ calls), so only final values that do not fit T
  are errors; a value that goes out of range and comes back is fine.
- **Worked example**: Laaksonen p. 93 (difference array checked entry by entry; +5 on closed [1,4] = our
  RangeAdd(1,5,5) → 3 8 6 6 6 5 2 2, hand-derived). **Errors**: ErrRange (RangeAdd, array unchanged);
  ErrOverflow (Materialize, dst unchanged). Zero value and nil are empty. **Input modified**: no (copies).
- **Tags**: difference-array, range-update. LeetCode 370, 1109.

### ALG-PREFIX-008 CountSubarraysSumEquals (count_subarrays_sum_equals.go:15)
`func CountSubarraysSumEquals[T integer](values []T, target T) int`
- **Source**: Laaksonen §9.1 p. 84 (range sum = prefix difference); CLRS 3rd ed. ch. 11 (hash table).
  The combination (count earlier prefixes equal to P − target) is ours from those two; no single textbook
  page found for it. **Lead: confirm this counts as "pattern + textbook".**
- **Best known**: O(n) expected with a hash map; O(n log n) worst-case with a sorted structure; Ω(n).
  Negative values rule out two pointers. **Lost**: O(n²) all pairs; sorting prefixes (O(n log n)).
- **Approach**: prefix sums in int128 as map keys, so no input overflows and no error is returned.
- **Worked example**: Laaksonen p. 77 array 1 3 2 5 1 1 2 3, target 8 → 1 (prefixes hand-derived).
- **Input modified**: no. **Tags**: prefix-sum, hash-map, subarray. LeetCode 560, 930.

### ALG-PREFIX-009 MaxSubarray (max_subarray.go:20)
`func MaxSubarray[T integer](values []T) (sum T, start, end int, err error)`
- **Source**: Bentley CACM 1984 (scanning algorithm, credited to Kadane); Programming Pearls 2nd ed. Column 8
  §8.4 (1st ed. Column 7); Laaksonen §2.4 Algorithm 3.
- **Best known**: Θ(n) time, O(1) space (Bentley shows any algorithm must read every element).
  **Lost**: divide and conquer O(n log n) (Bentley §8.3); cumulative-array O(n²).
- **Approach**: Kadane with restart when the running sum is ≤ 0 (gives the shortest subarray among those
  ending first). Non-empty subarray (Laaksonen allows empty; on his example both give 10). Overflow is exact:
  the running sum only grows past the type when the true maximum does.
- **Worked example**: Bentley's vector 31 −41 59 26 −53 58 97 −93 −23 84 → 187 at [2,7); Laaksonen → 10 at
  [1,6). **Errors**: ErrEmpty, ErrOverflow. **Input modified**: no.
- **Tags**: kadane, maximum-subarray. LeetCode 53.

### ALG-PREFIX-010 ProductExceptSelf (product_except_self.go:25)
`func ProductExceptSelf[T integer](dst, values []T) ([]T, error)`
- **Source**: Blelloch §1.1 (prescan under × with identity 1, run forwards and backwards).
- **Best known**: Θ(n) time, O(1) extra besides the output, no division. **Lost**: total product ÷ value
  (division forbidden by the brief; breaks on zeros); separate prefix and suffix arrays (O(n) extra).
- **Approach**: count zeros first (≥ 2 → all zero; 1 → only the zero's slot is non-zero). Otherwise products
  run as 64-bit magnitude + sign (math/bits.Mul64), the prefix is stored in the output as its negation (so
  |MinInt| is storable), and the sign comes from the parity of the other negatives. Overflow is exact: e.g.
  int8 64 2 −1 −1 → 2 64 −128 −128 although the prefix 128 exceeds MaxInt8.
- **Worked example**: none in the source; 1 2 3 4 → 24 12 8 6 hand-derived from the prescans.
- **Errors**: ErrOverflow (index in message), dst unchanged. **Input modified**: no.
- **Tags**: prefix-product, no-division. LeetCode 238.

### ALG-PREFIX-011 EquilibriumIndex (equilibrium_index.go:15)
`func EquilibriumIndex[T integer](values []T) int`
- **Source**: Laaksonen §9.1 (prefix difference); Bentley §8.2. Returns −1 when none.
- **Best known**: Θ(n), O(1). **Lost**: storing a prefix array (O(n) extra). Sums in int128: no error.
- **Worked example**: hand-derived −7 1 5 2 −4 3 0 → 3 (6 also balances; first wins).
- **Tags**: prefix-sum, equilibrium, pivot-index. LeetCode 724, 1991.

## HELD (no source)
- **MaxProductSubarray**: the method (track the largest and smallest product ending at each index, since a
  negative factor swaps them) is a Kadane variant found only on problem and interview sites; no textbook or
  paper located (searched 2026-11-12). Not built, no ID given. LeetCode 152.

## For the lead (outside my paths)
- `integer` constraint is now defined privately in algo/prefix, algo/twopointer, algo/window, and (by other
  builders) algo/hashing (exported `Integer`) and algo/heap (`number`). One home would be better (decision).
- A private int128 lives in algo/prefix/int128.go and a smaller copy in algo/twopointer/two_sum_sorted.go.
  Candidate shared home: algo/bits (taxonomy) or num/integer, exported. Decision for Devin.
- doc/sources entries needed: laaksonen_cphb_2018.md (new), blelloch_1993.md (new), crow_1984.md,
  viola_jones_2004.md, bentley (add Column 8 / CACM 1984 rows).

---

# algo/twopointer (SHORT = TWOPTR) — S2, new package (builder W1-B2)

**Not in doc/taxonomy.md**: the lead's brief created it ("NEW, SHORT TWOPTR"). Taxonomy needs the row
`| ALG | algo/twopointer | TWOPTR | two-index scans: pair and triple sums, in-place dedup, merge, partition |`
(lead's file). Variable-length windows live in algo/window, string forms in algo/text.

Error: `ErrLength` (merge_sorted_in_place.go:10).

## Sources (verification state)
- **Laaksonen**, *Competitive Programmer's Handbook*, 2018, §8.1 "Two pointers method" pp. 77–79 (read
  first-hand): pointers that move one way only, O(n); "Subarray sum" (positive values) and "2SUM problem"
  (1 4 5 6 7 9 9 10, x = 12 → 5 + 7).
- **Sedgewick & Wayne**, *Algorithms*, 4th ed., 2011, Exercise 1.4.15 "Faster 3-sum" (linear two-sum on a
  sorted array, then quadratic 3-sum). Exercise number and content from search results; the book's
  exercise text was not read first-hand (the booksite lists only some exercises). §2.5 "Duplicates" read
  first-hand on algs4.cs.princeton.edu/25applications.
- **Gajentaan & Overmars** 1995, doi:10.1016/0925-7721(95)00022-2 — Crossref-confirmed: Computational
  Geometry 5(3):165–185.
- **Knuth** TAOCP Vol. 3, 2nd ed., §5.2.4 "Sorting by Merging", Algorithm M, p. 158 (section and page from a
  bibliographic search; not re-read). **CLRS** 3rd ed. §2.3.1 MERGE (Figure 2.3 example), §7.1 PARTITION
  (Lomuto) — standard; not re-read.
- **Manacher** 1975, doi:10.1145/321892.321896 — Crossref-confirmed: JACM 22(3):346–351.
- **Go standard library**, slices.Compact, https://pkg.go.dev/slices#Compact.

### ALG-TWOPTR-001 TwoSumSorted (two_sum_sorted.go:23)
`func TwoSumSorted[T integer](sorted []T, target T) (i, j int, ok bool)`
- **Source**: Laaksonen §8.1 "2SUM" pp. 78–79; S&W Ex. 1.4.15.
- **Best known**: Θ(n) on sorted input (must read the data in the worst case), O(1). **Lost**: binary search
  per element (O(n log n), S&W TwoSumFast); hash set (O(n) expected, O(n) extra; needed only when unsorted).
- **Approach**: pointers from both ends; sums compared in int128, so wrap-around never matches (uint8
  100 + 200 ≠ 44). Result: smallest i, then largest j (proved in the oracle against an O(n²) scan).
- **Worked example**: Laaksonen → (2, 4). **Errors**: none (absent → (−1, −1, false)). Unsorted input:
  unspecified result, no panic. **Input modified**: no.
- **Tags**: two-pointers, two-sum. LeetCode 167, 1.

### ALG-TWOPTR-002 ThreeSum (three_sum.go:19)
`func ThreeSum[T integer](dst [][3]T, values []T, target T) [][3]T`
- **Source**: S&W Ex. 1.4.15; Gajentaan & Overmars 1995; Laaksonen §8.1.
- **Best known**: O(n²) is the practical bound; subquadratic word-RAM algorithms exist (Grønlund & Pettie,
  FOCS 2014; not verified further here, impractical); output can be Θ(n²) triples. **Lost**: hash per pair
  (O(n²) expected, O(n) extra, harder dedup); O(n³) enumeration.
- **Approach**: sorted copy (input not modified, one O(n) allocation), fixed first element, two pointers,
  duplicate skipping on all three positions; exact int128 sums. Output: each triple ascending, triples
  lexicographic, appended to dst.
- **Worked example**: none numeric in the sources; hand-derived 30 −40 −20 −10 40 0 10 5 → 4 triples.
- **Errors**: none. **Tags**: two-pointers, 3sum. LeetCode 15, 16, 18.

### ALG-TWOPTR-003 RemoveDuplicatesInPlace (remove_duplicates_in_place.go:17)
`func RemoveDuplicatesInPlace[T comparable](values []T) int`
- **Source**: S&W §2.5 "Duplicates"; Go slices.Compact.
- **Best known**: Θ(n), O(1). **Lost**: calling slices.Compact (returns a slice, and clears the tail: an extra
  pass of writes). The tail values[n:] is left as is (documented). Equality ==: NaN kept, −0/+0 merged.
- **Note for the lead**: this overlaps slices.Compact closely; kept because the brief asks for the
  length-returning in-place form.
- **Worked example**: hand-made 1 1 2 3 3 3 7 → 4. **Input modified**: yes (name says so).
- **Tags**: dedup, in-place. LeetCode 26, 80.

### ALG-TWOPTR-004 MergeSortedInPlace (merge_sorted_in_place.go:26)
`func MergeSortedInPlace[T cmp.Ordered](buf []T, m int, other []T) error`
- **Source**: Knuth §5.2.4 Algorithm M; CLRS §2.3.1 MERGE.
- **Best known**: Θ(m + k), O(1) extra (the buffer already has room). **Lost**: front-to-back merge (needs
  O(m) scratch or overwrites unread elements); append then sort (O((m+k) log(m+k))).
- **Approach**: back to front, take from buf only when strictly greater, so equal elements of buf stay first
  (stable; visible with −0/+0 in the test). cmp.Compare order (NaN first).
- **Worked example**: CLRS Figure 2.3, 2 4 5 7 + 1 2 3 6 → 1 2 2 3 4 5 6 7. **Errors**: ErrLength (m < 0 or no
  room), buf unchanged. other must not overlap buf (documented, not detectable). **Input modified**: buf yes.
- **Tags**: merge, in-place, stable. LeetCode 88.

### ALG-TWOPTR-005 IsPalindrome (is_palindrome.go:14)
`func IsPalindrome[T comparable](values []T) bool`
- **Source**: Manacher 1975 (definition); Laaksonen §8.1 (pointer pattern). Slices only; strings → algo/text.
- **Best known**: Θ(n) worst case, O(1). **Worked example**: hand-made 3 1 4 1 3. NaN off-centre → false.
- **Tags**: palindrome. LeetCode 125, 234 (on slices).

### ALG-TWOPTR-006 MoveZeroesInPlace (move_zeroes_in_place.go:14)
`func MoveZeroesInPlace[T comparable](values []T) int`
- **Source**: CLRS 3rd ed. §7.1 PARTITION (Lomuto: one boundary index, swap kept elements to it).
- **Best known**: Θ(n), O(1), ≤ n swaps. **Lost**: copy non-zeros forward then fill zeros (fewer writes when
  zeros are few, but rewrites −0 as +0, so the multiset changes); stable partition by rotations (O(n log n)).
- **Decision (lead to log)**: named **MoveZeroesInPlace**, not MoveZeroes: the project rule says a function
  that modifies its input says so in its name (…InPlace). Returns the non-zero count. Zero = T's zero value
  (== compare): −0 counts as zero (sign bit preserved by swapping), NaN does not.
- **Worked example**: hand-made 4 0 0 −2 7 0 → 4 −2 7 0 0 0. **Tags**: partition, in-place. LeetCode 283.

## HELD (no source)
- **ContainerMostWater**: the "move the shorter side inward" two-pointer and its correctness argument appear
  only on problem and interview sites; no textbook or paper found (searched 2026-11-12). LeetCode 11.
- **TrappingRainWater**: the formula water(i) = min(max left, max right) − h(i) and its two-pointer form have
  no textbook or paper source found (searched 2026-11-12); only problem/interview sites. LeetCode 42.
  (If Devin accepts the formula as a definition, it is a max-prescan (Blelloch §1.1) each way plus a sum.)

## Consumer evidence
the consumer backend/knowledgebase/series/ has no hand-rolled pair-sum, merge or partition; rank.go's
MergeIntervals belongs to algo/interval. None to cite.

---

# algo/hashing (SHORT = HASH)

Builder W1-B3. Package name `hashing`, not `hash`: callers import the standard `hash` package beside ours.
Not yet in `doc/taxonomy.md` (lead to add the row: `ALG | algo/hashing | HASH | counting, grouping and set
operations with hash tables`).

Common to the package:
- **Hash-table costs**: CLRS = Cormen, Leiserson, Rivest, Stein, *Introduction to Algorithms*, 3rd ed., MIT
  Press, 2009, ch. 11 (Theorems 11.1 and 11.2: under simple uniform hashing, a search costs Θ(1 + α) expected).
  Go's map is a hash table whose operations are expected O(1); every "Time" below is expected (average) time,
  and the worst case of each map operation is O(m) (all keys colliding). Second source: Sedgewick & Wayne,
  *Algorithms*, 4th ed., Addison-Wesley, 2011, §3.4 (hash tables, Proposition K: constant expected cost under
  the uniform hashing assumption). Knuth, TAOCP Vol. 3, 2nd ed., 1998, §6.4 (hashing).
- **Order**: the Go language specification, "For statements with range clause": "The iteration order over maps
  is not specified and is not guaranteed to be the same from one iteration to the next." No output of this
  package depends on map iteration order: groups and set results follow the order of first appearance in
  the input; ranked results have a total order on (count, value).
- **Equality** is Go's `==`. A float NaN is not equal to itself, so each NaN is its own key: it never matches
  and is never merged with another NaN. -0 and +0 are one key.
- No function modifies its input.

| ID | Name | Signature |
|---|---|---|
| ALG-HASH-001 | FrequencyCount | `func FrequencyCount[T comparable](values []T) map[T]int` |
| ALG-HASH-002 | TopKFrequent | `func TopKFrequent[T cmp.Ordered](values []T, k int) ([]Count[T], error)` |
| ALG-HASH-003 | TopKByCount | `func TopKByCount[T cmp.Ordered](counts map[T]int, k int) ([]Count[T], error)` |
| ALG-HASH-004 | GroupBy | `func GroupBy[T any, K comparable](values []T, key func(T) K) ([]Group[K, T], error)` |
| ALG-HASH-005 | FirstUnique | `func FirstUnique[T comparable](values []T) int` |
| ALG-HASH-006 | SetIntersection | `func SetIntersection[T comparable](a, b []T) []T` |
| ALG-HASH-007 | SetUnion | `func SetUnion[T comparable](a, b []T) []T` |
| ALG-HASH-008 | SetDifference | `func SetDifference[T comparable](a, b []T) []T` |
| ALG-HASH-009 | LongestConsecutive | `func LongestConsecutive[T Integer](values []T) (start T, length int)` |
| ALG-HASH-010 | GroupAnagrams | `func GroupAnagrams(words []string) [][]string` |
| ALG-HASH-011 | GroupAnagramsByCount | `func GroupAnagramsByCount(words []string, alphabet string) ([][]string, error)` |
| ALG-HASH-012 | TwoSum | `func TwoSum[T Integer](values []T, target T) (i, j int)` |
| ALG-HASH-013 | ContainsNearbyDuplicate | `func ContainsNearbyDuplicate[T comparable](values []T, maxGap int) (bool, error)` |

`Count[T]` is `{Value T; Count int}`. `Group[K, T]` is `{Key K; Items []T}`. `Integer` is the set of Go's
built-in integer types (`~int | ~int8 | … | ~uintptr`).

### ALG-HASH-001 FrequencyCount
- **Source**: CLRS ch. 11 (direct-address and hash tables as a dictionary); Sedgewick & Wayne §3.1, client
  FrequencyCounter (a symbol table counting occurrences), §3.4.
- **Best known**: O(n) expected time, Ω(n) to read the input; O(m) space for the m distinct values (the
  result). A sort-and-scan is O(n log n) and needs an ordered type.
- **Approach**: one map, not pre-sized (pre-sizing to n wastes memory when m ≪ n); it grows as keys arrive. **Lost**:
  sort a copy then count runs (O(n log n), O(n) extra, needs `cmp.Ordered`).
- **Input modified**: no. **Errors**: none. Empty input → empty, non-nil map.
- **Worked example**: Sedgewick & Wayne §3.1's FrequencyCounter idea on our own words: "to be or not to be"
  → to 2, be 2, or 1, not 1. Hand-derived.

### ALG-HASH-002 TopKFrequent · ALG-HASH-003 TopKByCount
- **Source**: counting by hash table (CLRS ch. 11) followed by a bounded min-heap of the k best (Sedgewick &
  Wayne §2.4, client TopM; Williams, "Algorithm 232: Heapsort", CACM 7(6):347–348, 1964,
  doi:10.1145/512274.3734138; Floyd, "Algorithm 245: Treesort", CACM 7(12):701, 1964,
  doi:10.1145/355588.365103). Same pattern as ALG-HEAP-001, keyed on (count, value).
- **Order (deterministic)**: count descending; equal counts by value ascending (`cmp.Compare`). Keys of a map
  are distinct, so (count, value) is a total order and map iteration order never reaches the output.
- **Best known**: O(n + m) expected (counting plus selection, CLRS §9.3) then O(k log k) to order the k;
  the lower bound is Ω(n + k log k). Bounded heap: O(n + m log k) expected time, O(m) for the counts (TopKFrequent
  only) and O(k) for the result.
- **Approach**: the bounded heap. **Lost**: bucket sort by count (CLRS §8.4): O(n + m) time, but n+1 buckets
  of O(n) extra space, and the deterministic tie-break needs every bucket sorted by value, which is
  O(m log m) when all counts are equal: no better than the heap in the worst case and more memory.
  Also lost: sorting all m (count, value) pairs (O(m log m)).
- **Input modified**: no. **Errors**: `ErrNegativeK` (k < 0). k > m returns all m, ranked. k = 0 → empty,
  non-nil slice.
- **Worked example**: counts of "to be or not to be" with k = 2 → (be 2), (to 2): equal counts, "be" < "to".
  Hand-derived.
- **Consumer evidence**: `the consumer backend/knowledgebase/series/rank.go:19-37` hand-rolls a `countHeap` over
  `container/heap` and a `TopK(counts map[K]int, k)` with exactly this order (count desc, key asc); its caller
  is `knowledgebase/analysis/symptoms/symptoms.go:661`. Consumer change: replace `series.TopK(counts, k)` with
  `hashing.TopKByCount(counts, k)` (field `Key` becomes `Value`; k < 0 is an error instead of nil; k = 0 gives
  an empty slice instead of nil).
- **Tags**: LeetCode 347, 692.

### ALG-HASH-004 GroupBy
- **Source**: CLRS ch. 11 (dictionary operations INSERT and SEARCH); Sedgewick & Wayne §3.4.
- **Best known**: O(n) expected time (n calls of key, n map operations), Ω(n). O(g) extra for the map of g
  distinct keys; O(n) for the result.
- **Approach**: a map from key to group index into a result slice; groups in order of their key's first
  appearance, items in input order (stable). **Lost**: sort by key then split runs (O(n log n), needs an
  ordered key, loses input order unless stable).
- **Input modified**: no (key is called once per value). **Errors**: `ErrNilKeyFunc` (key == nil).
- **Worked example**: 1…7 grouped by parity → [odd: 1 3 5 7], [even: 2 4 6]. Hand-derived.

### ALG-HASH-005 FirstUnique
- **Source**: counting with a hash table (CLRS ch. 11), then a second scan.
- **Best known**: O(n) expected, Ω(n) (an unread element might be the unique one). O(m) extra.
- **Approach**: count, then the first index whose count is 1. **Lost**: a linked hash map of candidates
  (one pass, but more bookkeeping and no gain on a slice that can be re-read); O(n²) pairwise scan.
- **Returns**: the index, or −1 when every value repeats or the slice is empty. **Errors**: none.
- **Worked example**: the bytes of "statistics" (s 3, t 3, a 1, i 2, c 1) → index 2 ('a'). Hand-derived.
- **Tags**: LeetCode 387.

### ALG-HASH-006 SetIntersection · ALG-HASH-007 SetUnion · ALG-HASH-008 SetDifference
- **Source**: set operations on hash tables, CLRS ch. 11; Knuth, TAOCP Vol. 3, 2nd ed., §6.4.
- **Order (deterministic)**: the result holds distinct values. Intersection: the values of a that occur in b,
  in order of first appearance in a. Union: the distinct values of a in order of first appearance, then the
  values of b not in a, in order of first appearance in b. Difference: the values of a not in b, in order of
  first appearance in a.
- **Best known**: O(|a| + |b|) expected time, Ω(|a| + |b|). O(|a| + |b|) extra (a set of b and a seen-set).
- **Approach**: a set of b and a seen-set for de-duplication. **Lost**: sort both and merge (O(N log N), needs
  `cmp.Ordered`, and the output order would be sorted rather than input order).
- **Input modified**: no. **Errors**: none. Empty results are non-nil.
- **Worked example**: a = 1 2 2 3 4, b = 4 3 3 5 → ∩ = 3 4, ∪ = 1 2 3 4 5, a∖b = 1 2. Hand-derived.
- **Tags**: LeetCode 349.

### ALG-HASH-009 LongestConsecutive
- **Source**: membership in a hash set, CLRS ch. 11 (expected O(1) SEARCH); the run is counted only from a
  value whose predecessor is absent, so each value is visited O(1) times.
- **Best known**: O(n) expected time, O(n) extra; the comparison lower bound for the sorting approach is
  Ω(n log n) (element distinctness reduces to it in the comparison model), which hashing sidesteps.
- **Approach**: set of values; for each value v with v−1 absent, walk v+1, v+2, … while present. Overflow
  safe: the walk stops at the type's maximum and the predecessor test skips the type's minimum.
  Ties broken by the smallest start. **Lost**: sort and scan (O(n log n), O(n) copy).
- **Returns**: (start, length); empty input → (0, 0). Duplicates count once. **Errors**: none.
- **Worked example**: 10 4 20 1 3 2 11 → start 1, length 4 (1 2 3 4). Hand-derived.
- **Tags**: LeetCode 128.

### ALG-HASH-010 GroupAnagrams · ALG-HASH-011 GroupAnagramsByCount
- **Source**: the anagram signature: Bentley, *Programming Pearls*, 2nd ed., Addison-Wesley, 2000, Column 2,
  problem C (posed in §2.1; signatures in §2.8, pp. 18–20): sort the letters of each word to get its signature; words with equal
  signatures are anagrams. Grouping by signature: a hash table (CLRS ch. 11).
- **Definition**: two words are anagrams when they hold the same multiset of runes (Unicode code points).
  Groups are in order of first appearance; words keep input order; duplicates are kept.
- **Best known**: GroupAnagrams: O(Σ Lᵢ log Lᵢ) time (sorting each word), O(Σ Lᵢ) extra for the keys.
  GroupAnagramsByCount: O(Σ Lᵢ + n·σ) time with σ = number of runes in the alphabet, O(n·σ) extra for the
  keys; it is faster when words are long and σ is small (the classic 26-letter count key).
- **Approach**: sorted-rune key (Bentley's signature); count key over a caller-given alphabet, encoded as
  uvarint counts. **Lost**: a product-of-primes signature (overflows past a few letters); pairwise anagram
  tests (O(n²)).
- **Input modified**: no. **Errors**: ByCount: `ErrAlphabet` (a rune repeated in alphabet, or invalid
  UTF-8 in it), `ErrNotInAlphabet` (a word holds a rune outside alphabet).
- **Worked example**: Bentley's own (Column 2): "deposit", "dopiest", "posited", "topside" share the signature
  "deiopst". Test: those four plus "pots", "stop", "tops" → 2 groups.
- **Tags**: LeetCode 49.

### ALG-HASH-012 TwoSum
- **Source**: complement lookup in a hash table (CLRS ch. 11); the problem is CLRS 3rd ed., Exercise 2.3-7
  (a pair summing to x, solved there in Θ(n lg n) by sorting), here in O(n) expected with hashing.
- **Best known**: O(n) expected time, O(n) extra. Sorting plus two pointers is O(n log n) and loses the
  original indices without an O(n) index array.
- **Approach**: scan j = 0…n−1; look up target − values[j] in a map of value → first index. Returns the pair
  with the smallest j, and for it the smallest i < j. The subtraction is overflow-checked: if target − v is
  out of the type's range, no partner exists for v. **Lost**: sort + two pointers (above).
- **Returns**: (i, j) with i < j, or (−1, −1). **Errors**: none.
- **Worked example**: 3 8 11 −2 7, target 9 → (2, 3), since 11 + (−2) = 9 and no earlier pair sums to 9.
  Hand-derived.
- **Tags**: LeetCode 1.

### ALG-HASH-013 ContainsNearbyDuplicate
- **Source**: a sliding window held in a hash set (CLRS ch. 11 dictionary; the window as in ALG-WINDOW).
- **Best known**: O(n) expected time; O(min(n, maxGap)) extra, the window's set.
- **Approach**: a set of the last maxGap values; insert, and delete the one leaving. **Lost**: a map of last
  index per value (O(m) extra, m distinct, which is larger than the window when maxGap ≪ m).
- **Returns**: true when some i < j has values[i] == values[j] and j − i ≤ maxGap. maxGap = 0 → false.
- **Errors**: `ErrNegativeGap` (maxGap < 0).
- **Worked example**: 5 1 7 5 with maxGap 3 → true; with maxGap 2 → false. Hand-derived.
- **Tags**: LeetCode 219.

---

# algo/stack (SHORT = STACK)

Builder W1-B3. Not yet in `doc/taxonomy.md` (lead to add: `ALG | algo/stack | STACK | monotonic stacks,
nearest greater/smaller, stack evaluators`).

Common to the package: ordering is `cmp.Compare` / `cmp.Less` (D-002: NaN is the least value, -0 equals +0).
Functions that return one index per element append to `dst` and return the extended slice (D-003): nil
allocates once, `buf[:0]` reuses a buffer. No function modifies its input.

| ID | Name | Signature |
|---|---|---|
| ALG-STACK-001 | NextGreater | `func NextGreater[T cmp.Ordered](dst []int, values []T) []int` |
| ALG-STACK-002 | NextSmaller | `func NextSmaller[T cmp.Ordered](dst []int, values []T) []int` |
| ALG-STACK-003 | PreviousGreater | `func PreviousGreater[T cmp.Ordered](dst []int, values []T) []int` |
| ALG-STACK-004 | PreviousSmaller | `func PreviousSmaller[T cmp.Ordered](dst []int, values []T) []int` |
| ALG-STACK-005 | LargestRectangleHistogram | `func LargestRectangleHistogram(heights []int64) (int64, error)` |
| ALG-STACK-006 | MaximalRectangle | `func MaximalRectangle(grid [][]bool) (int, error)` |
| ALG-STACK-007 | DailyTemperatures | `func DailyTemperatures[T cmp.Ordered](dst []int, readings []T) []int` |
| ALG-STACK-008 | MinStack | `type MinStack[T cmp.Ordered]` with `Push(T)`, `Pop() (T, bool)`, `Top() (T, bool)`, `Min() (T, bool)`, `Len() int` |
| ALG-STACK-009 | EvaluateRPN | `func EvaluateRPN(tokens []string) (int64, error)` |
| ALG-STACK-010 | BalancedBrackets | `func BalancedBrackets(text, pairs string) (bool, error)` |
| — | SimplifyPath | DROPPED: duplicates `path.Clean` (below) |

### ALG-STACK-001 NextGreater · -002 NextSmaller · -003 PreviousGreater · -004 PreviousSmaller
- **Source**: the all-nearest-smaller-values (ANSV) problem and its sequential stack algorithm: Berkman,
  Schieber, Vishkin, "Optimal doubly logarithmic parallel algorithms based on finding all nearest smaller
  values", Journal of Algorithms 14(3):344–370, 1993, doi:10.1006/jagm.1993.1018 (verified on Crossref
  2026-11-12). Earlier: Gabow, Bentley, Tarjan, "Scaling and related techniques for geometry problems",
  STOC '84, pp. 135–143, 1984, doi:10.1145/800057.808675 (verified on Crossref). The stack held as a linked
  list: Knuth, TAOCP Vol. 1, 3rd ed., 1997, §2.2.3 (linked allocation of a stack).
- **Definition**: NextGreater[i] = the least j > i with values[j] > values[i] (strict), else −1. NextSmaller
  uses <. PreviousGreater/PreviousSmaller: the greatest j < i, strict, else −1.
- **Best known**: O(n) time, Ω(n) to read the input; the textbook stack is O(n) extra.
- **Approach**: the monotonic stack kept as a linked list through the result itself: scanning from the far
  end, the candidates for i are i±1, then result[i±1], then result[result[i±1]], …; every skipped index is
  skipped for good, so the total work is O(n) amortized and no stack is allocated: O(1) extra.
  **Lost**: an explicit index stack (same O(n) time, O(n) extra and a second allocation); the O(n²) scan.
- **Input modified**: no. **Errors**: none (no input is invalid).
- **Worked example**: 2 7 3 5 4 6 8 (hand-derived). NextGreater 1 6 3 5 5 6 −1; NextSmaller −1 2 −1 4 −1 −1 −1;
  PreviousGreater −1 −1 1 1 3 1 −1; PreviousSmaller −1 0 0 2 2 4 5.
- **Tags**: LeetCode 496, 503, 739, 901.

### ALG-STACK-005 LargestRectangleHistogram
- **Source**: every maximal rectangle under a histogram has the height of some bar i and spans exactly from
  its previous smaller bar to its next smaller bar (exclusive) — the ANSV application, Berkman, Schieber,
  Vishkin 1993 (above), with Gabow, Bentley, Tarjan 1984 (above).
- **Best known**: O(n) time, Ω(n); O(n) extra (the stack). Lower bound on extra space not known to us.
- **Approach**: one left-to-right pass with an index stack of non-decreasing heights; a bar's rectangle is
  closed when it is popped (its right limit is the popping index, its left limit the new top). Popping on
  ≥ makes the last of a run of equal bars compute the full width. **Lost**: PreviousSmaller + NextSmaller
  as two calls (O(n) extra twice, two result arrays); divide and conquer on the minimum (O(n log n), or O(n)
  only with a Cartesian tree, more memory).
- **Units**: heights in any unit u; the area is in u × bar widths.
- **Input modified**: no. **Errors**: `ErrNegativeHeight` (a height < 0); `ErrOverflow` (height × width
  exceeds int64). Empty → 0.
- **Worked example**: heights 2 1 5 6 2 3 → 10 (bars 5 and 6, height 5, width 2). Hand-derived.
- **Tags**: LeetCode 84.

### ALG-STACK-006 MaximalRectangle
- **Source**: reduction to the histogram, row by row: the largest all-true rectangle has a bottom row r;
  above row r, column c holds a run of h_r[c] trues ending at r, and a rectangle with bottom row r fits
  exactly when its height ≤ min h_r over its columns, which is the histogram problem on h_r. This proof is
  written here; the histogram step is ALG-STACK-005 (ANSV, Berkman, Schieber, Vishkin 1993). The reduction
  is folklore in the largest-empty-rectangle literature (Naamad, Lee, Hsu, "On the maximum empty rectangle
  problem", Discrete Applied Mathematics 8(3):267–277, 1984); a primary source stating this exact reduction
  was not found and is FLAGGED for Devin.
- **Best known**: O(R·C) time, Ω(R·C) to read the grid; O(C) extra.
- **Approach**: update h_r in place per row; run the histogram pass with one reused stack. **Lost**: the
  dynamic-programming left/right/height arrays (also O(R·C), three O(C) arrays); enumerating rectangles
  (O(R²C²)).
- **Input modified**: no. **Errors**: `ErrRagged` (rows of different lengths). Empty grid → 0.
- **Worked example**: our own 4×5 grid (rows 10100, 10111, 11111, 10010) → 6. Hand-derived.
- **Tags**: LeetCode 85.

### ALG-STACK-007 DailyTemperatures
- **Source**: NextGreater (ALG-STACK-001), by a call; the wait is next − i, or 0 when none.
- **Best known**: O(n), O(1) extra (as NextGreater).
- **Units**: readings in any unit, compared only; the result counts positions (days when there is one
  reading per day).
- **Worked example**: 73 74 75 71 69 72 76 73 (our own readings) → 1 1 4 2 1 1 0 0. Hand-derived.
- **Tags**: LeetCode 739.

### ALG-STACK-008 MinStack
- **Source**: Gajewska, Tarjan, "Deques with heap order", Information Processing Letters 22(4):197–200, 1986,
  doi:10.1016/0020-0190(86)90028-1 (find-min in O(1) worst case on a deque; the stack case keeps, with
  each element, the minimum of the elements below it).
- **Best known**: O(1) per operation; Push O(1) amortized here (slice growth), O(n) space.
- **Approach**: each slot stores (value, minimum so far). **Lost**: an auxiliary stack of minima pushed only
  on ≤ (saves space when minima are rare, but two slices and a second comparison per pop).
- **Errors**: none; Pop, Top and Min on an empty stack return (zero, false). Zero value is an empty stack.
- **Worked example**: push 5 3 7 3 → Min 3; pop, pop → Min 3; pop → Min 5. Hand-derived.
- **Tags**: LeetCode 155.

### ALG-STACK-009 EvaluateRPN
- **Source**: Hamblin, "Translation to and from Polish notation", The Computer Journal 5(3):210–213, 1962,
  doi:10.1093/comjnl/5.3.210 (verified on Crossref); Knuth, TAOCP Vol. 1, 3rd ed., 1997, §2.2.1 (stacks);
  Sedgewick & Wayne, *Algorithms*, 4th ed., 2011, §1.3, Exercise 1.3.11 (EvaluatePostfix) beside Dijkstra's
  two-stack evaluator in the same section.
- **Best known**: O(n) time for n tokens (each token once), O(n) extra for the operand stack.
- **Semantics**: tokens are base-10 int64 literals (an optional sign) or one of + − * /. Division truncates
  toward zero (the Go specification, "Arithmetic operators"). Results are exact or an error.
- **Errors**: `ErrMalformed` (no tokens, an operator with fewer than two operands, or more than one value
  left); `ErrBadToken` (neither a literal nor an operator); `ErrOverflow` (a literal or a result outside
  int64, including MinInt64 / −1); `ErrDivideByZero`.
- **Worked example**: Hamblin's reverse Polish order on our own expression (5 + 1) × (12 − 4) / 3 → tokens
  5 1 + 12 4 − * 3 / → 16. Hand-derived.
- **Tags**: LeetCode 150.

### ALG-STACK-010 BalancedBrackets
- **Source**: Sedgewick & Wayne, *Algorithms*, 4th ed., 2011, §1.3, Exercise 1.3.4 (Parentheses: a stack
  client that checks bracket balance); Knuth, TAOCP Vol. 1, 3rd ed., 1997, §2.2.1.
- **Semantics**: pairs is a string of opener/closer rune pairs, e.g. "()[]{}"; other runes are ignored.
  Balanced: every closer matches the most recent unmatched opener, and none is left open.
- **Best known**: O(n + p) time (n runes of text, p runes of pairs), O(n) extra worst case (the stack).
- **Errors**: `ErrPairs` (odd rune count, an opener equal to its closer, a rune in two roles, U+FFFD, or
  invalid UTF-8).
- **Worked example**: Sedgewick & Wayne's own: "[()]{}{[()()]()}" → true, "[(])" → false.
- **Tags**: LeetCode 20.

### SimplifyPath — DROPPED
Unix path canonicalisation (collapse "//", drop ".", resolve ".." lexically, root stays root) is exactly the
standard library's `path.Clean` on a rooted path: Go standard library, package path, func Clean (rules 1–4,
read 2026-11-12), which cites Pike, "Lexical File Names in Plan 9 or Getting Dot-Dot Right". POSIX pathname
resolution (IEEE Std 1003.1, Base Definitions, Pathname Resolution) resolves ".." through the file system
(symbolic links), which no pure function can do. Callers use `path.Clean`. No ID given.

---

# algo/queue (SHORT = QUEUE)

Builder W1-B3. Not yet in `doc/taxonomy.md` (lead to add: `ALG | algo/queue | QUEUE | FIFO queues: two-stack
queue, ring buffer`). The monotonic deque has one home, `algo/window` (`window.MonotonicDeque`, another
builder); it is not built here.

| ID | Name | Signature |
|---|---|---|
| ALG-QUEUE-001 | QueueFromTwoStacks | `type QueueFromTwoStacks[T any]` with `Enqueue(T)`, `Dequeue() (T, bool)`, `Front() (T, bool)`, `Len() int` |
| ALG-QUEUE-002 | CircularBuffer | `func NewCircularBuffer[T any](capacity int, policy FullPolicy) (*CircularBuffer[T], error)` with `Push(T) (overwritten bool, err error)`, `Pop() (T, bool)`, `Front() (T, bool)`, `Back() (T, bool)`, `At(i int) (T, bool)`, `Len() int`, `Cap() int` |

### ALG-QUEUE-001 QueueFromTwoStacks
- **Source**: Okasaki, "Simple and efficient purely functional queues and deques", Journal of Functional
  Programming 5(4):583–592, 1995, doi:10.1017/S0956796800001489 (verified 2026-11-12; §2 recalls the
  two-list "batched" queue: a front list and a reversed rear list, reversed when the front empties, O(1)
  amortized). Okasaki, *Purely Functional Data Structures*, Cambridge University Press, 1998, §5.2.
  Cormen et al., *Introduction to Algorithms*, 3rd ed., 2009, Exercise 10.1-6 (a queue from two stacks).
  Hood and Melville, "Real-time queue operations in pure LISP", Information Processing Letters 13(2), 1981,
  give the worst-case O(1) version (citation not re-verified; not used).
- **Best known**: O(1) amortized per operation (each element is pushed and popped on each stack once); O(1)
  worst case needs Hood–Melville's incremental reversal (or a ring buffer, ALG-QUEUE-002, when the
  capacity is bounded). O(n) space.
- **Approach**: two Go slices as stacks; Dequeue/Front move the whole input stack, reversed, to the output
  stack only when the output is empty. **Lost**: Hood–Melville real-time queue (worst-case O(1), several
  times the bookkeeping, no gain for a non-persistent Go value); a slice with re-slicing from the front
  (`q = q[1:]`) leaks the underlying array's head until growth copies.
- **Errors**: none; Dequeue/Front on an empty queue return (zero, false). Zero value is an empty queue.
  Not safe for concurrent use.
- **Worked example**: enqueue 1 2 3, dequeue → 1, enqueue 4, dequeue all → 2 3 4. Hand-derived.
- **Tags**: LeetCode 232.

### ALG-QUEUE-002 CircularBuffer
- **Source**: Knuth, TAOCP Vol. 1, 3rd ed., 1997, §2.2.2 (sequential allocation; a queue in a circular array
  of M slots, the rear and front pointers advancing modulo M); Cormen et al., 3rd ed., §10.1 (ENQUEUE and
  DEQUEUE on an array with wrap-around).
- **Best known**: O(1) worst case per operation, O(capacity) space, allocated once at construction.
- **Approach**: a slice of fixed length, a head index and a count; the tail is (head + count) mod capacity.
  Policy chosen by the caller: `Reject` (Push on a full buffer returns `ErrFull` and changes nothing) or
  `Overwrite` (Push on a full buffer replaces the oldest element and reports overwritten = true).
  **Lost**: head and tail indices with one slot left empty to tell full from empty (wastes a slot);
  power-of-two capacity with a mask (faster modulo but forces the capacity).
- **Errors**: `ErrCapacity` (capacity < 1), `ErrPolicy` (an unknown policy), `ErrFull` (Push under Reject
  when full). Pop/Front/Back/At on missing elements return (zero, false).
- **Worked example**: capacity 3, Overwrite: push 1 2 3 4 → 2 3 4 (1 overwritten); Reject: push 1 2 3 4 →
  ErrFull on 4, holds 1 2 3. Hand-derived.
- **Tags**: LeetCode 622, 641.

---

# algo/heap (SHORT = HEAP) — S2 additions (builder W1-B4)

ALG-HEAP-001 TopK is S1 (shipped); its sift is unexported and stays as is. `Heap[T]` below is the
reusable exported type that wave 2 imports (graph Dijkstra/Prim, list LFU, window MovingMedian).

## Heap API (for wave 2)

```go
var ErrNilLess = errors.New("heap: less function is nil")
type Heap[T any] struct{ /* unexported */ }
func NewHeap[T any](less func(a, b T) bool, capacityHint int) (*Heap[T], error)   // ALG-HEAP-002
func Heapify[T any](less func(a, b T) bool, items []T) (*Heap[T], error)          // ALG-HEAP-003 (copies items)
func (h *Heap[T]) Push(v T) bool         // ALG-HEAP-004 (false only on a zero/nil Heap)
func (h *Heap[T]) Pop() (v T, ok bool)   // ALG-HEAP-005
func (h *Heap[T]) Peek() (v T, ok bool)  // ALG-HEAP-006
func (h *Heap[T]) Len() int              // ALG-HEAP-007
```
less(a, b) true means a comes out first (min-heap for `<`, max-heap for `>`).

**Decision needed (lead to log as a D-id)**: API shape.
- Chosen: constructor with the less function (`NewHeap`, `Heapify`), error on nil less.
  **Lost**: a zero-usable value plus `SetLess`. A zero Heap has no order, so Push must either panic, mis-order
  silently, or report failure; a constructor makes the order part of the value from the start. The zero value
  is still safe: Push returns false and stores nothing, Pop/Peek report empty, nil receivers never panic.
  Also lost: wrapping `container/heap` (interface-boxing of every element through `any` in Push/Pop, and a
  caller-written Len/Less/Swap type per use; not benchmarked here) and a `cmp.Ordered`-only heap (Dijkstra/LFU need
  struct elements).
- No `Fix(i)` / `Remove(i)`: an element's slot moves on every Push and Pop, so an index is useless without a
  swap callback that reports moves. Decrease-key users (Dijkstra, Prim) push the new key and skip stale
  entries on Pop (lazy deletion; O((V+E) log V) as with decrease-key). An indexed heap (Sedgewick & Wayne
  §2.4 IndexMinPQ) is the future home if needed, under a new ID.

### ALG-HEAP-002 NewHeap · 003 Heapify · 004 Push · 005 Pop · 006 Peek · 007 Len
- **Source**: Williams, "Algorithm 232: Heapsort", CACM 7(6):347–348, 1964, doi:10.1145/512274.3734138
  (INHEAP/OUTHEAP: sift-up insert, sift-down removal). Floyd, "Algorithm 245: Treesort", CACM 7(12):701, 1964,
  doi:10.1145/355588.365103 (bottom-up construction). Cormen et al., *Introduction to Algorithms*, 3rd ed.,
  MIT Press, 2009, §6.1–6.3, §6.5 (BUILD-MAX-HEAP's O(n) bound; MAX-HEAP-INSERT, HEAP-EXTRACT-MAX, HEAP-MAXIMUM).
- **Best known**: binary heap: insert and delete-min O(log n), find-min O(1), build O(n). Lower bound: n
  inserts + n delete-mins sort, so Ω(n log n) comparisons in total (Ω(log n) amortized per pair). Fibonacci /
  pairing heaps give O(1) amortized insert but worse constants and pointer-heavy layout; for a general
  library the array heap wins on cache and allocation. **Lost**: d-ary heap (fewer levels, more compares per
  level; a tuning choice that can be a later variant).
- **Approach**: array heap, iterative sift with a hole (one move per level, not a swap). Heapify is Floyd's
  bottom-up pass over a copy (≤ 2n compares).
- **Input modified**: Heapify copies its slice.
- **Worked example**: CLRS §6.3 Figure 6.3: BUILD-MAX-HEAP of 4 1 3 2 16 9 10 14 8 7 → 16 14 10 8 7 9 3 2 4 1
  (checked by hand against this code's sift too). Pop order of KNUTH16 = sorted KNUTH16.
- **Errors**: `ErrNilLess`. Empty Pop/Peek → (zero, false). Zero/nil Heap: Push false.
- **Tests**: layout check, cases, struct elements, zero/nil, random Push/Pop oracle against a sorted model,
  AllocsPerRun 0 within capacity. Benchmarks NewHeap, Heapify/1e3…1e5, HeapPush/…, HeapPop/…, HeapPeek, HeapLen.
- **Tags**: heap, priority-queue, binary-heap, heapify. LeetCode: 703, 1046 (pattern only).

### ALG-HEAP-008 KSmallestPairs
`func KSmallestPairs[T number](dst []IndexPair, a, b []T, k int) ([]IndexPair, error)`; `IndexPair{I, J int}`;
`number` = integer and float kinds (unexported constraint).
- **Source**: k-way merge with a min-heap, CLRS 3rd ed., §6.5, Exercise 6.5-9; row i is a[i]+b[0…]. The X+Y
  matrix view: Frederickson, Johnson, "Generalized selection and ranking: sorted matrices", SIAM J. Comput.
  13(1):14–30, 1984, doi:10.1137/0213002 (verified via Crossref-indexed listings 2026-11-12).
- **Best known**: output Ω(k). Heap frontier O(m + k log m), m = min(k, len a), O(m) space. Frederickson–Johnson
  select the k-th sum in O(m) and the k smallest then in O(k) more by a threshold scan, but they need
  arithmetic-free selection machinery and give no sorted order; sorting the k output costs k log k anyway,
  so for sorted output the heap is within a constant of the bound.
- **Approach**: heap of (sum, i, j) seeded with column 0 of the first m rows (Floyd heapify), pop and push
  (i, j+1). Ties ordered by i then j, so the output is exactly "all pairs sorted by (sum, i, j), first k".
  **Lost**: the two-successor frontier from (0,0) with a visited set (needs a hash set, same bound); full
  enumeration and sort (O(nm log nm)).
- **Input modified**: no. Output by index pairs: the caller recovers values and gets positions too.
- **Worked example**: hand-derived: a = 0 3 10, b = 1 2 20, k = 5 → (0,0) (0,1) (1,0) (1,1) (2,0).
- **Errors**: `ErrNegativeK`; `ErrNotFinite` (NaN/±Inf float read); `ErrOverflow` (integer sum formed overflows).
  Unsorted input: precondition, unspecified result, no panic.
- **Tags**: heap, pairs, k-smallest, x-plus-y, sorted-matrix. LeetCode: 373.

### ALG-HEAP-009 KthSmallestInSortedMatrix
`func KthSmallestInSortedMatrix[T cmp.Ordered](matrix [][]T, k int) (T, error)`
- **Source**: CLRS 3rd ed., §6.5, Exercise 6.5-9 (k-way merge, stopped at the k-th). Problem: Frederickson,
  Johnson 1984 (above), which selects in an n×n sorted matrix in O(n).
- **Best known**: O(r) for an r×r sorted matrix (Frederickson–Johnson 1984). Lower bound Ω(r) (Ω(min(k, r))
  in general).
- **Approach (decision for the lead to log)**: heap over the first m = min(k, r) rows: O(r + k log m) time
  (O(r) is the rectangularity check), O(m) space; generic over `cmp.Ordered` (strings work) and exact.
  **Lost**: (1) binary search on the value range, O(r·log(max−min)) — needs arithmetic on T (no strings),
  terminates only by integer steps (floats need a bit-level bisection), and its count is O(r+c) per step;
  faster when k ≈ r·c/2 on integers, so it is the candidate for an integer-only variant under a new ID.
  (2) Frederickson–Johnson O(r): best bound, but a long, constant-heavy algorithm; not justified until a
  consumer needs large matrices.
- **Input modified**: no.
- **Worked example**: hand-derived: rows 1 4 7 / 2 5 8 / 3 6 9 → k-th smallest = k for k = 1…9.
- **Errors**: `ErrRagged` (rows differ in length); `ErrRank` (k < 1 or k > r·c; empty matrix). k 1-based (A-001).
- **Tags**: heap, selection, k-th-smallest, sorted-matrix, young-tableau. LeetCode: 378.

### ALG-HEAP-010 NewMedianStream · 011 MedianStream.Add · 012 MedianStream.Median · 013 MedianStream.Len
`type MedianStream[T cmp.Ordered]` (zero value ready) · `NewMedianStream[T](capacityHint int) *MedianStream[T]` ·
`Add(v T) bool` · `Median() (lower, upper T, err error)` · `Len() int`.
- **Source**: Sedgewick, Wayne, *Algorithms*, 4th ed., Addison-Wesley, 2011, §2.4 Creative Problems,
  "Dynamic median-finding": insert O(log n), find the median O(1), using a max-heap for keys below the median
  and a min-heap for keys above (booksite algs4.cs.princeton.edu/24pq, read 2026-11-12; problem text and
  hint confirmed there). Lower/upper median: CLRS ch. 9 introduction.
- **Best known**: O(log n) per insert, O(1) median, O(n) space (the textbook bound above). No lower bound is
  claimed here (none cited).
- **Approach**: two heaps, lower half max-heap holds the extra value. **Lost**: an order-statistic tree
  (O(log n) both, more memory and code); re-selecting each time (O(n) per query); a sorted slice (O(n) insert).
- **Median shape**: (lower, upper, err) as ALG-SELECT-002, so callers average or not as their domain says.
- **Worked example**: KNUTH16 in order, medians after 1…5 keys: (503,503) (87,503) (503,503) (87,503)
  (503,503); after 16: (509,512), matching ALG-SELECT-002. Hand-derived.
- **Errors**: `ErrEmpty`. nil *MedianStream: Add false.
- **Tags**: heap, median, running-median, streaming, two-heaps. LeetCode: 295 (480 is window MovingMedian, wave 2).

## Not here
- IntervalSchedulingMaxNonOverlap was listed under heap in the request; it needs no heap (sort + greedy) and
  its one home is algo/interval (ALG-INTVL-008).

---

# algo/interval (SHORT = INTVL) — S2 additions (builder W1-B4)

Convention (ALG-INTVL-001, A-005): `Interval[T]` is closed `[Start, End]`. Sharing an endpoint is overlap:
merges, intersects in a point, needs two rooms, conflicts in scheduling. Half-open callers pass End−1 unit.
Exception: Skyline (buildings raise the outline on Left ≤ x < Right, the key-point reading). "Sorted and
disjoint" = each starts strictly after the previous ends (what MergeOverlapping returns).

## Consumer evidence
`the consumer backend/knowledgebase/series/rank.go:70` `MergeIntervals(iv []Interval, gap time.Duration)`: half-open
`[Start, End)` of `time.Time`, merge when next.Start ≤ last.End + gap. Replaced by **ALG-INTVL-001
MergeOverlapping**, not InsertInterval: convert to `Interval[int64]{Start.UnixNano(), End.UnixNano() + gap}`
(the gap folded into End; equal-start/touching then merge exactly as rank.go does), call MergeOverlapping,
subtract gap from each End and convert back. Their `End > Start` check stays in the consumer (ours allows
points). Consumer change, listed for the lead; no edit made.

### ALG-INTVL-002 InsertInterval
`func InsertInterval[T cmp.Ordered](sorted []Interval[T], added Interval[T]) ([]Interval[T], error)`
- **Source**: Knuth, TAOCP Vol. 3, 2nd ed., 1998, §5.2.4 (one-pass merge of sorted sequences); Klee 1977
  (union by a sweep in Start order), doi:10.1080/00029890.1977.11994336.
- **Best known**: O(n) (the result is a new list of up to n+1; Ω(n) to write it). Finding the spot is
  O(log n) by binary search, but the copy dominates; validation is free in the same pass.
- **Approach**: copy those ending before, merge those touching, copy the rest. **Lost**: append + MergeOverlapping
  (O(n log n)); binary search for the bounds (same O(n) because of the copy; adds code).
- **Worked example**: hand: [1,2] [5,7] [9,9] [12,15] + [6,10] → [1,2] [5,10] [12,15].
- **Errors**: `ErrInverted`, `ErrNotSortedDisjoint`. Input not modified.
- **Tags**: interval, insert, merge. LeetCode: 57.

### ALG-INTVL-003 IntersectLists
`func IntersectLists[T cmp.Ordered](a, b []Interval[T]) ([]Interval[T], error)`
- **Source**: Knuth, TAOCP Vol. 3, 2nd ed., §5.2.4, Algorithm M (two-way merge).
- **Best known**: Θ(n+m) (must read both lists). Space O(1) extra; result ≤ n+m pieces, one allocation.
- **Approach**: two pointers; emit [max starts, min ends] when non-empty; advance the one ending first.
  **Lost**: marking a grid (unbounded domain); interval tree queries (O((n+m) log)).
- **Worked example**: hand: [1,3] [5,9] ∩ [2,6] [8,10] → [2,3] [5,6] [8,9]. [1,3]∩[3,5] = [3,3].
- **Errors**: `ErrInverted`, `ErrNotSortedDisjoint` (named a or b). **Tags**: intersection, two-pointer. LeetCode: 986.

### ALG-INTVL-004 MinMeetingRooms
`func MinMeetingRooms[T cmp.Ordered](intervals []Interval[T]) (int, error)`
- **Source**: Kleinberg, Tardos, *Algorithm Design*, Addison-Wesley, 2006, §4.1, "A Related Problem: Scheduling
  All Intervals" (interval partitioning; rooms = depth; confirmed via course notes citing §4.1, 2026-11-12).
  CLRS 3rd ed., §16.1, Exercise 16.1-4 (lecture halls / interval-graph coloring; confirmed 2026-11-12).
- **Best known**: O(n log n); Ω(n log n) in the algebraic decision-tree model (element distinctness reduces
  to depth ≥ 2 with point intervals). Space O(n).
- **Approach**: sort starts and ends separately (one 2n buffer), sweep. **Lost**: K–T's greedy with a min-heap
  of room end times (same bound, more work per step; it also assigns rooms, which we do not return).
- **Worked example**: hand: [1,4] [2,5] [6,8] [7,9] [3,3] → 3 (point 3).
- **Errors**: `ErrInverted`. **Tags**: sweep-line, depth, partitioning. LeetCode: 253 (shape only; half-open there).

### ALG-INTVL-005 MissingRanges · ALG-INTVL-006 SummaryRanges
`func MissingRanges[T integer](values []T, lower, upper T) ([]Interval[T], error)` ·
`func SummaryRanges[T integer](values []T) ([]Interval[T], error)` (`integer` unexported: all int/uint kinds)
- **Source**: Klee 1977 (union of intervals by one sweep in sorted order; here the unit intervals [v, v+1):
  SummaryRanges is their union, MissingRanges its complement in [lower, upper]); Knuth §5.2.4 (one-pass merge).
- **Best known**: Θ(n) on sorted input (read every value). O(1) extra.
- **Approach**: single scan, sortedness checked in the same pass (ErrNotSorted), no v+1 at the type's maximum
  (overflow-safe; tested at int8/uint8/int extremes). **Lost**: sort inside (O(n log n), hides a precondition
  violation); a bitmap over the range (O(range)).
- **Worked examples**: hand: 2 3 7 9 in [0,10] → [0,1] [4,6] [8,8] [10,10]; 0 1 2 4 5 7 → [0,2] [4,5] [7,7].
- **Errors**: `ErrInverted` (upper < lower), `ErrNotSorted`. Values outside [lower, upper] and duplicates ignored.
- **Tags**: ranges, gaps, runs. LeetCode: 163, 228.

### ALG-INTVL-007 Covers
`func Covers[T cmp.Ordered](intervals []Interval[T], target Interval[T]) (bool, error)`
- **Source**: Klee 1977 (sort by left end, sweep the union); Fredman, Weide, CACM 21(7):540–544, 1978,
  doi:10.1145/359545.359553 (Ω(n log n) for the union's measure; a bound for this decision problem is not
  claimed).
- **Best known**: O(n log n) by sort and sweep; O(n) extra for the sorted copy.
- **Approach**: clone, sort by Start, extend reach; fail at the first gap. Continuous semantics: [1,2] ∪ [3,4]
  does not cover [1,4] (A-005). **Lost**: MergeOverlapping then binary search (same bound, builds the whole union).
- **Worked example**: hand: [3,8] [1,4] [10,12] covers [2,7], not [2,11].
- **Errors**: `ErrInverted`. **Tags**: cover, coverage, union. LeetCode: 1893, 1024 (shape).

### ALG-INTVL-008 IntervalSchedulingMaxNonOverlap
`func IntervalSchedulingMaxNonOverlap[T cmp.Ordered](intervals []Interval[T]) ([]int, error)`
- Listed under heap in the request; one home is interval (no heap used).
- **Source**: CLRS 3rd ed., §16.1 (activity selection, Theorem 16.1, GREEDY-ACTIVITY-SELECTOR); Kleinberg–Tardos
  §4.1 (earliest finish time; "stays ahead").
- **Best known**: O(n log n) (Θ(n) if sorted by End); Ω(n log n) unsorted (decision-tree). O(log n) extra.
- **Approach**: sort indices by (End, index), keep each whose Start > last kept End (closed). Returns indices.
  **Lost**: DP over sorted intervals (weighted version; needless here); returning only the count.
- **Worked example**: CLRS Figure 16.1, 11 activities → {a1, a4, a8, a11} = indices 0 3 7 10 (same under closed).
- **Errors**: `ErrInverted`. **Tags**: greedy, activity-selection. LeetCode: 435, 646, 452 (shapes).

### ALG-INTVL-009 Skyline
`func Skyline[T cmp.Ordered](buildings []Building[T]) ([]KeyPoint[T], error)`; `Building{Left, Right, Height}`,
`KeyPoint{X, Height}`.
- **Source**: Manber, *Introduction to Algorithms: A Creative Approach*, Addison-Wesley, 1989, §5.6 "The Skyline
  Problem" (problem and O(n log n) divide and conquer; section number confirmed by TOC listing 2026-11-12).
  Sweep: de Berg, Cheong, van Kreveld, Overmars, *Computational Geometry*, 3rd ed., Springer, 2008, §2.1
  (plane-sweep paradigm), with algo/heap's Heap as the status structure (lazy deletion).
- **Best known**: O(n log n); Ω(n log n) (sorting reduces to it). O(n) space.
- **Approach**: sweep distinct x of all edges; push buildings starting, drop ended ones from the top only,
  emit when the top height changes. **Lost**: Manber's divide and conquer (same bound, recursion depth log n,
  more allocation for merged halves); a balanced multiset (no stdlib type).
- **Worked example**: hand: (1,5,3) (2,4,6) (4,7,2) (9,10,4) → (1,3) (2,6) (4,3) (5,2) (7,0) (9,4) (10,0).
- **Errors**: `ErrInverted` (Right < Left), `ErrNegativeHeight` (below zero or NaN). Zero width/height ignored.
- **Tags**: skyline, sweep-line, heap. LeetCode: 218.

---

# algo/list (SHORT = LIST) — S2 additions (builder W2-B5)

New package; `algo/list` is in `doc/taxonomy.md` ("linked lists, LRU"). Every ID below is new (001…017).

## Decisions for the lead (each with the alternative that lost)

| # | Decision | Alternative | Why it lost |
|---|---|---|---|
| B5-1 | Functions that relink a caller's nodes carry `InPlace` (`ReverseInPlace`, `ReverseKGroupInPlace`, `MergeTwoSortedInPlace`, `RemoveNthFromEndInPlace`, `ReorderInPlace`, `PartitionInPlace`), as L21 | the bare names of the brief (`Reverse`, …) | BRIEF_COMMON: "inputs never modified unless the name says so (…InPlace)"; a copying `Reverse` would allocate n nodes to do what the pattern exists to avoid |
| B5-2 | Every singly linked walk is guarded by Floyd's two-pointer test (`walk`, unexported): a cyclic input returns `ErrCycle` and is left unchanged; nil is the empty list | trust the caller (no check) / a visited-set map | an unguarded walk hangs on a cycle (the brief forbids it); a map costs O(n) space and allocation; Floyd costs ~1.5 extra pointer reads per node and O(1) space |
| B5-3 | `DoublyLinkedList[T]` is built although `container/list` exists | use `container/list` | container/list stores `any`: every Value is boxed (an extra allocation for non-pointer T) and type-asserted at each read; no `iter.Seq`; ours is generic, unboxed, has `All`/`Backward` iterators, and the same O(1) ownership check (an Element from another list is refused, never corrupts). LRU and LFU are built on it |
| B5-4 | Package name `list` (taxonomy-approved) | `linkedlist` | a caller importing `container/list` beside ours must rename one; `container/list` is not in CLAUDE.md's forbidden list and our `DoublyLinkedList` replaces its use. Flag for Devin |
| B5-5 | Types (Node, DoublyLinkedList, LRUCache, LFUCache) take one ID each, methods listed in the type's doc comment (ALG-QUEUE-002 shape) | an ID per method (ALG-HEAP-004…007 shape) | fewer rows for one concept; both shapes exist in S2 — lead to pick one for the index |
| B5-6 | A key not equal to itself (a float NaN, or a struct holding one) is refused by the caches' Put with `ErrKeyNotComparable` | store it | Go map lookups and deletes never find a NaN key, so it could never be read or evicted: the map would grow past capacity (L17 treats NaN as never equal) |
| B5-7 | LFU ties (equal use count) broken by least recent use; Put on an existing key counts as a use | arbitrary within a frequency (the paper's SET) | deterministic eviction; the order costs nothing with a linked node list |

## Consumer evidence

`/Users/devinanderson/Desktop/red_lite_cache/internal/cache` (read 2026-11-12): **no hand-rolled LRU or
eviction list found.** `store.go:15-16` "the store itself does not evict" (bounded by MaxKeys and TTL);
`groups_bench_test.go:15` and `:163` show the eviction it benchmarks is a uniform random draw from a dense slot
array (its D-043), not recency. `ordered.go:30` is a B-tree leaf chain (`next *bnode`), not a cache list.
So LRUCache/LFUCache answer a stated need ("LRU needs") but replace no existing code there.

| ID | Name | Signature |
|---|---|---|
| ALG-LIST-001 | Node | `type Node[T any] struct { Value T; Next *Node[T] }` |
| ALG-LIST-002 | FromSlice | `func FromSlice[T any](values []T) *Node[T]` |
| ALG-LIST-003 | AppendValues | `func AppendValues[T any](dst []T, head *Node[T]) ([]T, error)` |
| ALG-LIST-004 | Length | `func Length[T any](head *Node[T]) (int, error)` |
| ALG-LIST-005 | DoublyLinkedList | `type DoublyLinkedList[T any]`; `type Element[T any]`; methods PushFront, PushBack, InsertBefore, InsertAfter, Remove, MoveToFront, MoveToBack, Front, Back, Len, All, Backward; Element.Next, Element.Prev |
| ALG-LIST-006 | ReverseInPlace | `func ReverseInPlace[T any](head *Node[T]) (*Node[T], error)` |
| ALG-LIST-007 | ReverseKGroupInPlace | `func ReverseKGroupInPlace[T any](head *Node[T], k int) (*Node[T], error)` |
| ALG-LIST-008 | DetectCycle | `func DetectCycle[T any](head *Node[T]) bool` |
| ALG-LIST-009 | CycleStart | `func CycleStart[T any](head *Node[T]) (start *Node[T], cycleLen int)` |
| ALG-LIST-010 | Middle | `func Middle[T any](head *Node[T]) (*Node[T], error)` |
| ALG-LIST-011 | MergeTwoSortedInPlace | `func MergeTwoSortedInPlace[T cmp.Ordered](a, b *Node[T]) (*Node[T], error)` |
| ALG-LIST-012 | RemoveNthFromEndInPlace | `func RemoveNthFromEndInPlace[T any](head *Node[T], fromEnd int) (*Node[T], error)` |
| ALG-LIST-013 | ReorderInPlace | `func ReorderInPlace[T any](head *Node[T]) (*Node[T], error)` |
| ALG-LIST-014 | PartitionInPlace | `func PartitionInPlace[T cmp.Ordered](head *Node[T], pivot T) (*Node[T], error)` |
| ALG-LIST-015 | LRUCache | `func NewLRUCache[K comparable, V any](capacity int) (*LRUCache[K, V], error)`; Get, Peek, Put, Delete, Len, Cap |
| ALG-LIST-016 | LFUCache | `func NewLFUCache[K comparable, V any](capacity int) (*LFUCache[K, V], error)`; Get, Put, Delete, Len, Cap |
| ALG-LIST-017 | AddTwoNumbers | `func AddTwoNumbers(a, b *Node[int], radix int) (*Node[int], error)` |

Errors (package `list`): `ErrCycle`, `ErrSharedNodes`, `ErrGroupSize`, `ErrPosition`, `ErrCapacity`,
`ErrKeyNotComparable`, `ErrRadix`, `ErrDigit`. nil is the empty list everywhere; no function panics on a
nil list, a cyclic list, a nil receiver or a foreign Element.

### ALG-LIST-001 Node · 002 FromSlice · 003 AppendValues · 004 Length
- **Source**: Knuth, *The Art of Computer Programming*, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (linked
  allocation: INFO and LINK fields, Λ the null link). Cormen et al., *Introduction to Algorithms*, 3rd ed.,
  MIT Press, 2009, §10.2 (singly linked lists).
- **Best known**: building or reading n nodes Θ(n). Length of a list that may be cyclic: Θ(n) time, O(1) space
  with Floyd's test (ALG-LIST-008's source).
- **Approach**: FromSlice allocates the n nodes as one block (one allocation) and links them. AppendValues and
  Length walk with the Floyd guard (B5-2). **Lost**: one allocation per node (n allocations; slower to build,
  nothing gained for a library type the caller may relink anyway; documented that the block stays alive
  while any of its nodes is reachable).
- **Input modified**: never. **Errors**: `ErrCycle` (AppendValues returns dst unchanged).
- **Worked example**: FromSlice([1 2 3]) → 1→2→3→nil; Length 3; AppendValues → [1 2 3]. Hand-derived.
- **Tags**: linked-list, singly-linked.

### ALG-LIST-005 DoublyLinkedList
- **Source**: Knuth Vol. 1, 3rd ed., §2.2.5 (doubly linked lists with a list head: LLINK/RLINK, deletion and
  insertion in O(1) given the node). CLRS 3rd ed., §10.2 (LIST-INSERT, LIST-DELETE, sentinels, Figure 10.4).
- **Best known**: O(1) worst case per insert/remove/move given an element; O(n) to iterate; O(n) space.
- **Approach**: circular list around a sentinel element embedded in the list value (lazy init: the zero
  value is an empty list). Each Element records its list; operations given an Element of another list (or a
  removed one, or nil) do nothing and report false/nil. `All`/`Backward` are `iter.Seq[T]`; the next element
  is read before yield, so removing the current element during iteration is safe. **Lost**: `container/list`
  (B5-3); a non-circular list with nil ends (two more branches per operation).
- **Errors**: none; Remove/Move report ok=false on a foreign element. Not safe for concurrent use; must not be
  copied after first use (the sentinel's address is the list's identity).
- **Worked example**: PushBack 1,2,3; PushFront 0; MoveToBack(front) → 1 2 3 0; Remove(back) → 1 2 3. Hand-derived.
- **Tags**: linked-list, doubly-linked, deque, data-structure. LeetCode 707 (pattern only).

### ALG-LIST-006 ReverseInPlace
- **Source**: Knuth Vol. 1, 3rd ed., §2.2.3, Exercise 7 ("invert" a linked linear list by changing its links;
  confirmed 2026-11-12). CLRS 3rd ed., Exercise 10.2-7 (Θ(n) nonrecursive reversal using constant storage).
- **Best known**: Θ(n) time (every link changes), O(1) extra space.
- **Approach**: one pass with three pointers after the Floyd guard. **Lost**: recursion (O(n) stack depth,
  unbounded on a long list); copying into a slice and relinking (O(n) space).
- **Input modified**: yes (links). **Errors**: `ErrCycle` (list untouched). nil → nil.
- **Worked example**: 1→2→3→4 → 4→3→2→1. Hand-derived from Knuth's description.
- **Tags**: linked-list, reverse, in-place. LeetCode 206.

### ALG-LIST-007 ReverseKGroupInPlace
- **Source**: same as 006 (inversion applied to each consecutive block of k nodes); Knuth Vol. 1 §2.2.3 Ex. 7.
- **Best known**: Θ(n), O(1) extra.
- **Approach**: length first (guarded), then ⌊n/k⌋ block reversals splicing each block's tail to the next
  block; a final block of fewer than k nodes keeps its order. **Lost**: recursion per block (O(n/k) stack);
  a stack of k nodes (O(k) space).
- **Errors**: `ErrGroupSize` (k < 1), `ErrCycle`. k = 1 or k > n: unchanged.
- **Worked example**: 1→2→3→4→5, k = 2 → 2→1→4→3→5; k = 3 → 3→2→1→4→5. Hand-derived.
- **Tags**: linked-list, reverse, in-place. LeetCode 25, 24 (k = 2).

### ALG-LIST-008 DetectCycle · ALG-LIST-009 CycleStart
- **Source**: Knuth, *The Art of Computer Programming*, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1,
  Exercise 6 (find μ and λ in O(μ + λ) steps with a bounded number of memory locations; Knuth ascribes the
  X_n = X_2n method to R. W. Floyd) and its answer. Confirmed 2026-11-12 (exercise statement). Brent's
  power-of-two variant is Exercise 7 (R. P. Brent, 1977).
- **Best known**: Θ(μ + λ) time, O(1) space; μ = nodes before the cycle, λ = cycle length. Lower bound:
  every node must be read (Ω(μ + λ)).
- **Approach**: Floyd's tortoise and hare; CycleStart then restarts one pointer at head and steps both by one
  until they meet (the entry, μ steps), then walks the cycle once (λ). **Lost**: Brent's method (fewer
  evaluations of f, which matters for an expensive f; here f is a pointer read and Floyd's meeting point is
  needed for the entry anyway); a visited set (O(n) space).
- **Errors**: none. DetectCycle(nil) false; CycleStart on an acyclic list → (nil, 0).
- **Worked example**: 1→2→3→4→5→(back to 3): μ = 2, λ = 3, start = node 3. Self loop: λ = 1. Hand-derived.
- **Tags**: linked-list, cycle-detection, floyd, two-pointer, tortoise-hare. LeetCode 141, 142, 287 (pattern).

### ALG-LIST-010 Middle
- **Source**: Floyd's two-speed walk (Knuth Vol. 2 §3.1 Ex. 6, as above): when the fast pointer runs off the
  end after k double steps, the slow pointer is at index k = ⌊n/2⌋. Laaksonen, *Guide to Competitive
  Programming*, Springer (L18) — not needed; Knuth suffices.
- **Best known**: Θ(n) (the end must be found), O(1) space; one pass.
- **Approach**: slow/fast; returns the node at index ⌊n/2⌋ (0-based) — for even n the second of the two
  middles. A meeting of the pointers proves a cycle (`ErrCycle`), so no separate guard pass. **Lost**: count
  then walk n/2 (1.5 passes as well, but needs its own cycle guard).
- **Errors**: `ErrCycle`. nil → nil.
- **Worked example**: 1→2→3→4→5 → 3; 1→2→3→4 → 3. Hand-derived.
- **Tags**: linked-list, middle, two-pointer. LeetCode 876.

### ALG-LIST-011 MergeTwoSortedInPlace
- **Source**: Knuth, *TAOCP* Vol. 3, 2nd ed., Addison-Wesley, 1998, §5.2.4, Algorithm L (list merge sort; its
  merge step relinks two sorted linked lists). CLRS 3rd ed., §2.3.1 (MERGE).
- **Best known**: Θ(m + n) time, ≤ m + n − 1 comparisons, O(1) extra space; lower bound for merging equal
  sizes is 2n − 1 comparisons (Knuth Vol. 3 §5.3.2), met.
- **Approach**: a dummy head, take from a while !(b < a): stable (equal values: a's before b's).
  Values compared with `cmp.Less` (NaN first, D-002). **Lost**: recursion (O(m + n) stack).
- **Input modified**: both lists are relinked into the result. Unsorted input gives a merged permutation, no
  error. **Errors**: `ErrCycle`, `ErrSharedNodes` (the two lists share a node: their last nodes coincide —
  merging would build a cycle).
- **Worked example**: 1→3→5 and 1→2→4 → 1(a)→1(b)→2→3→4→5. Hand-derived.
- **Tags**: linked-list, merge, sorted, stable, in-place. LeetCode 21.

### ALG-LIST-012 RemoveNthFromEndInPlace
- **Source**: Knuth Vol. 1 §2.2.3 (deleting a node by relinking its predecessor) + the length of the list.
- **Best known**: Θ(L) time, O(1) space. One pass with a gap of n between two pointers is possible, but the
  cycle guard reads the list once anyway, so length-then-walk is the same order with simpler code. **Lost**:
  the gap pointer (it alone cannot detect a cycle).
- **Input modified**: yes; the removed node's Next is set to nil. `fromEnd` is 1-based (1 = last node).
- **Errors**: `ErrPosition` (fromEnd < 1 or > length), `ErrCycle`.
- **Worked example**: 1→2→3→4→5, fromEnd 2 → 1→2→3→5. Hand-derived.
- **Tags**: linked-list, two-pointer, in-place. LeetCode 19.

### ALG-LIST-013 ReorderInPlace
- **Source**: composition of cited steps: the middle (ALG-LIST-010's walk), inversion of the second half
  (Knuth Vol. 1 §2.2.3 Ex. 7), and alternating merge of two lists by relinking (Knuth Vol. 1 §2.2.3).
  Proof of the order: first half L0…L⌈n/2⌉−1, reversed second half Ln−1…L⌈n/2⌉; alternating them gives
  L0, Ln−1, L1, Ln−2, …, which is the target order.
- **Best known**: Θ(n), O(1) extra. **Lost**: copying nodes into a slice and relinking by index (O(n) space).
- **Errors**: `ErrCycle`. nil, one or two nodes: unchanged.
- **Worked example**: 1→2→3→4 → 1→4→2→3; 1→2→3→4→5 → 1→5→2→4→3. Hand-derived.
- **Tags**: linked-list, reorder, two-pointer, in-place. LeetCode 143.

### ALG-LIST-014 PartitionInPlace
- **Source**: stable partition of a linked list into two sublists by relinking (Knuth Vol. 1 §2.2.3, linked
  allocation; CLRS 3rd ed. §7.1 PARTITION for the predicate; a linked list allows the stable form in O(1)
  space, unlike the array form).
- **Best known**: Θ(n), O(1) extra.
- **Approach**: two dummy-headed tails (values < pivot, values ≥ pivot by `cmp.Less`), each keeping input
  order, then concatenated. NaN pivot: nothing is less than NaN under cmp.Less, so all go to the second part.
- **Errors**: `ErrCycle`.
- **Worked example**: 1→4→3→2→5→2, pivot 3 → 1→2→2→4→3→5. Hand-derived.
- **Tags**: linked-list, partition, stable, in-place. LeetCode 86.

### ALG-LIST-015 LRUCache
- **Source**: Mattson, Gecsei, Slutz, Traiger, "Evaluation techniques for storage hierarchies", IBM Systems
  Journal 9(2):78–117, 1970, doi:10.1147/sj.92.0078 (LRU replacement as a stack algorithm; verified on
  Crossref 2026-11-12). Implementation: hash table (CLRS 3rd ed., ch. 11) + doubly linked list in recency
  order (CLRS §10.2; this package's ALG-LIST-005).
- **Best known**: O(1) per Get/Put/Delete (expected, hash map), O(capacity) space. **Lost**: a timestamp per
  entry with an O(n) scan or a heap (O(log n)) for the victim; `container/list` (B5-3).
- **Approach**: `map[K]*Element[entry]` + DoublyLinkedList, front = most recent. Put on a full cache reuses the
  evicted element (no allocation in steady state). Get marks the key most recent; Peek does not.
  `Put(key, value) (evicted Entry[K, V], didEvict bool, err error)`.
- **Errors**: `ErrCapacity` (capacity < 1), `ErrKeyNotComparable` (B5-6). Zero value or nil: capacity 0,
  Get misses, Put stores nothing and reports the given pair as evicted. No clock, no goroutines; not safe for
  concurrent use.
- **Worked example**: capacity 2: Put(1,a) Put(2,b) Get(1)=a Put(3,c) evicts 2; Get(2) miss; Put(4,d)
  evicts 1. Hand-derived (Mattson's LRU stack order).
- **Tags**: cache, lru, eviction, hash-map, linked-list. LeetCode 146.

### ALG-LIST-016 LFUCache
- **Source**: Shah, Mitra, Matani, "An O(1) algorithm for implementing the LFU cache eviction scheme",
  technical report, 16 August 2010, arXiv:2110.11602 (posted 2021; read 2026-11-12): §5, a doubly linked
  *frequency list* of nodes, each holding the items with that access count, a hash table by key, and every
  item pointing to its frequency node; insert, access and evict O(1).
- **Best known**: O(1) per operation (the paper), O(capacity) space. **Lost**: min-heap on (count, recency)
  with a map (O(log n) per op, the paper's §4 baseline; algo/heap has no decrease-key).
- **Approach**: the paper's structure, with the per-frequency SET replaced by a recency-ordered linked list
  (B5-7): the victim is the least recently used item of the lowest count. A bucket whose only item moves up
  to an absent count+1 is renumbered in place (no allocation). The count saturates at 2^64−1.
- **Errors**: `ErrCapacity`, `ErrKeyNotComparable`. Zero value/nil as LRUCache. Put on an existing key replaces
  the value and counts as a use.
- **Worked example**: capacity 2: Put(1) Put(2) Get(1) Put(3) evicts 2 (count 1 vs 2); Get(3) (count 2);
  Put(4) evicts 1 (counts tie at 2, 1 used less recently). Hand-derived.
- **Tags**: cache, lfu, eviction, hash-map, linked-list. LeetCode 460.

### ALG-LIST-017 AddTwoNumbers
- **Source**: Knuth, *TAOCP* Vol. 2, 3rd ed., 1997, §4.3.1, Algorithm A (addition of nonnegative integers in
  radix b, least significant digit first, carry k ∈ {0, 1}).
- **Best known**: Θ(max(m, n)) time; output max(m, n) or +1 digits. O(1) extra besides the result.
- **Approach**: lengths first (guard), result allocated as one block of max(m, n)+1 nodes; a shorter operand is
  read as zeros (Knuth assumes equal length n; padding with zeros is the same number). Digit sums in uint
  (no overflow for any radix ≤ MaxInt). nil is zero; nil + nil = nil.
- **Input modified**: never. **Errors**: `ErrRadix` (radix < 2), `ErrDigit` (digit < 0 or ≥ radix), `ErrCycle`.
- **Worked example**: 342 + 465 in radix 10, as 2→4→3 + 5→6→4 → 7→0→8 (807); 99 + 1 → 0→0→1. Hand-derived.
- **Tags**: linked-list, big-integer, addition, carry. LeetCode 2, 445 (pattern).

---

# algo/tree (SHORT = TREE) — S2 wave 2 (builder W2-B6)

New package (planned in doc/taxonomy.md: `algo/tree`, TREE, "balanced trees, traversals"). No balanced
(red-black/AVL) type is built; `SortedToBalancedBST` gives a perfectly balanced BST from sorted keys.

Sources confirmed 2026-11-12:
- **CLRS**: Cormen, Leiserson, Rivest, Stein, *Introduction to Algorithms*, 3rd ed., MIT Press, 2009, ch. 12
  (§12.1 BST property and INORDER-TREE-WALK, Theorem 12.1; §12.2 TREE-SEARCH, TREE-MINIMUM; §12.3
  TREE-INSERT, TRANSPLANT, TREE-DELETE), §22.2 (breadth-first search), Appendix B.5.3 (height of a tree).
- **Hibbard** 1962: T. N. Hibbard, "Some combinatorial properties of certain trees with applications to
  searching and sorting", J. ACM 9(1):13–28, doi:10.1145/321105.321108 (Crossref checked).
- **Knuth** TAOCP Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.3.1: Algorithm T (inorder with a stack),
  exercise 12 (preorder analogue), exercise 13 (postorder: a second state per stack entry), exercise 7
  (preorder + inorder determine the tree). Exercise numbers checked against an index of §2.3.1 exercises.
- **Laaksonen**, *Competitive Programmer's Handbook*, draft of 3 July 2018, https://cses.fi/book/book.pdf
  (read 2026-11-12; textbook per lead decision L18): §9.2 Binary indexed tree pp. 86–88; §9.3 Segment tree
  pp. 89–93; §9.4 "Range updates" pp. 93–94; §14.2 Diameter, Algorithm 1 pp. 135–136; §14.4 Binary trees
  p. 139; §18.3 Lowest common ancestor p. 167; §28.1 Lazy propagation pp. 258–260.
- **Fenwick** 1994: P. M. Fenwick, "A new data structure for cumulative frequency tables", Software:
  Practice and Experience 24(3):327–336, doi:10.1002/spe.4380240306 (Crossref checked).
- **Wirth**, *Algorithms and Data Structures*, Oberon version, 2004 (author's PDF AD2012), §4.4.1,
  pp. 151–152 (perfectly balanced tree: root, n DIV 2 nodes left, n − n DIV 2 − 1 right).
- **Aho, Hopcroft, Ullman**, "On finding lowest common ancestors in trees", SIAM J. Comput. 5(1):115–132,
  1976, doi:10.1137/0205011 (Crossref checked).
- **Bernstein**, "Netstrings", 1997, https://cr.yp.to/proto/netstrings.txt (also IETF
  draft-bernstein-netstrings-02): `len ":" bytes ","`, no leading zeros.
- **Go spec**, "For statements with range clause" (go1.27 doc/go_spec.html, #For_range): "If the loop body
  terminates (such as by a break statement), yield returns false and must not be called again."

Shared worked tree (Laaksonen §14.4 p. 139): root 1; 1 → (2, 3); 2 → (4, 5); 5 → (6, nil); 3 → (nil, 7).
Preorder 1 2 4 5 6 3 7 · inorder 4 2 6 5 1 3 7 · postorder 4 6 5 2 7 3 1.

Conventions: every traversal is iterative (explicit stack or queue), so a degenerate tree of 1e5 nodes is
safe; segment-tree recursion is over the index range, depth ⌈log2 n⌉ + 1 ≤ 64. The input to every tree
function must be a tree (no node reachable twice); a cyclic structure is a precondition breach (documented,
not detected). Keys compare with cmp.Compare (NaN first, D-002). Integer sums are overflow-checked
(ErrOverflow); floats reject NaN/±Inf (ErrNotFinite). Every mutating method is atomic: on error the
structure is unchanged.

## Types
`Node[T any]{Value T; Left, Right *Node[T]}` — the general binary tree. `BST[T cmp.Ordered]` — zero value is
an empty set (no duplicates). `FenwickTree[T]`, `FenwickRangeUpdate[T]` (T signed integer or float),
`SegmentTree[T any]`, `LazySegmentTree[T, U any]` with `LazyOps[T, U]`.

## BST — ALG-TREE-001 Insert · 002 Search · 003 Delete · 004 Len · 005 Root
- **Source**: CLRS §12.2–12.3; Hibbard 1962 (deletion by successor). Iterative.
- **Best known**: O(h) per operation for an unbalanced BST (h = height; n worst, O(log n) expected for random
  insertion order, CLRS §12.4). Lower bound Ω(log n) comparisons in the worst case for any comparison
  search tree. **Lost**: a self-balancing tree (red-black, CLRS ch. 13) — not requested; recorded for later.
- **Approach**: CLRS TREE-DELETE relinks the successor node (no value copy, so outstanding *Node pointers
  stay valid). **Lost**: Hibbard's original copy-successor-key form (invalidates node identity).
- **Worked example**: CLRS Figure 12.2 shape built by inserting 15 6 18 3 7 17 20 2 4 13 9 (hand-checked
  inorder 2 3 4 6 7 9 13 15 17 18 20), then deleting 6 (two children: successor 7), 13 (one child), 2 (leaf).
- **Errors**: none; Insert false for a duplicate, Delete false when absent; nil *BST safe.

## Traversals — ALG-TREE-006 Inorder · 007 Preorder · 008 Postorder · 009 LevelOrder (iter.Seq[T])
- **Source**: Knuth §2.3.1 Algorithm T, ex. 12, ex. 13; CLRS §22.2 (BFS by levels); Go spec #For_range.
- **Best known**: Θ(n) time; O(h) extra for depth-first with a stack, O(w) for level order (w = max level
  width). Morris threading gives O(1) extra but writes into the tree temporarily — **lost**: the input must
  not be modified and an early break would leave threads behind.
- **Early break**: each iterator returns immediately when yield returns false (test counts calls).
- **Worked example**: the shared tree's three orders; level order 1 2 3 4 5 7 6.

## ALG-TREE-010 Height · 011 IsBalanced
- Height = edges on the longest root-to-leaf path (CLRS B.5.3): empty −1, one node 0 (MaxDepth = Height+1).
- IsBalanced: every node's subtree heights differ by at most 1 (Adel'son-Vel'skii–Landis condition as
  stated in Knuth TAOCP Vol. 3, 2nd ed., 1998, §6.2.3). Computed by one iterative postorder fold, O(n) time,
  O(h) extra; stops at the first unbalanced node. **Lost**: recursive height per node (O(n²) on a chain).
- Worked example (hand-checked): shared tree height 3 (path 1 2 5 6); balanced (node 1: 2 vs 1, node 2:
  0 vs 1, node 5: 0 vs −1, node 3: −1 vs 0); the chain 1→2→3 is not (node 1: −1 vs 1).

## ALG-TREE-012 LowestCommonAncestorBST · 013 LowestCommonAncestor
- **Source**: definition AHU 1976; Laaksonen §18.3. General tree: one postorder fold counting p and q;
  the first node whose subtree holds both is the LCA. O(n), O(h) extra. **Lost**: Euler tour + RMQ or
  binary lifting (O(n log n) preprocessing, worth it only for many queries; a future type).
- BST: the LCA of keys p ≤ q is the first node on the search path with p ≤ key ≤ q (proof: above it both
  keys lie on one side, so the node is a common ancestor's descendant on the shared path; at it they split
  or one equals it, so no deeper node is an ancestor of both). Search via CLRS §12.2; both keys then
  confirmed present. O(h), O(1) extra.
- Worked (hand-checked on the insertion tree 15 6 18 3 7 17 20 2 4 13 9): LCA(2, 4) = 3, LCA(9, 17) = 15,
  LCA(13, 9) = 13, LCA(4, 9) = 6; the shared tree (general): LCA(4, 6) = 2, LCA(6, 7) = 1.

## ALG-TREE-014 FromPreIn · 015 FromInPost
- **Source**: Knuth §2.3.1 ex. 7; Laaksonen §14.4 p. 139 (pre+in and post+in determine the tree; pre+post
  does not). T comparable, values distinct.
- **Best known**: Θ(n). **Approach**: one stack pass (a node becomes a right child when the inorder cursor
  matches the stack top), nodes in one allocation, then both traversals of the result are compared with the
  inputs (O(n)), so any inconsistent input is rejected. **Lost**: recursive split with an index map
  (O(n) but recursion depth n on a chain).
- **Errors**: ErrLength (lengths differ), ErrDuplicate, ErrInconsistent (not the traversals of one tree;
  a NaN never matches itself, so float NaN input is ErrInconsistent).
- Worked example: shared tree from (pre, in) and from (in, post).

## ALG-TREE-016 SortedToBalancedBST
- **Source**: Wirth §4.4.1 pp. 151–152 (perfectly balanced: left gets n DIV 2 nodes). Θ(n) time, one
  allocation of n nodes, O(log n) extra stack (explicit). Errors: ErrNotSorted (not strictly ascending).
- Worked: 1..7 → root 4, (2: 1, 3), (6: 5, 7).

## ALG-TREE-017 KthSmallestBST · 018 ValidateBST
- **Source**: CLRS §12.1 Theorem 12.1 (inorder walk lists a BST's keys in sorted order). Kth: O(h + k),
  O(h) extra, stops early. Errors ErrRange (k < 1 or k > n). **Lost**: size-augmented nodes (O(h) per
  query; needs a different node type). Validate: strictly ascending inorder, O(n), O(h); ErrNotBST names
  the first pair out of order.

## ALG-TREE-019 Serialize · 020 Deserialize
- **Format (ours)**: the preorder sequence with every empty subtree written as `#` and every node's value
  as a netstring (Bernstein 1997): `len:bytes,`. Empty tree `#`; a leaf 7 encoded as "7" is `1:7,##`.
  Uniqueness: preorder with explicit empty markers is a prefix code of extended binary trees (each node
  token is followed by exactly two subtree encodings), so parsing is unambiguous; netstring framing makes
  any byte value safe. Values are encoded by the caller's appendValue / parseValue.
- Θ(n + total value bytes); O(h) extra stack, iterative both ways. Errors: ErrNilFunc, ErrMalformed
  (bad length, leading zero, missing comma, truncated, trailing bytes), parse errors wrapped.

## ALG-TREE-021 MaxPathSum
- **Source**: Laaksonen §14.2 Algorithm 1 (every path has a highest node; per node the best downward
  path and the best path through it), weighted, a downward gain clipped at 0. Θ(n), O(h) extra.
- Overflow: only a positive overflow is possible (a node plus non-negative gains), which means the true
  maximum does not fit T → ErrOverflow. ErrEmptyTree, ErrNotFinite.
- Worked: −10 → (9, 20 → (15, 7)) gives 42 (15+20+7), hand-checked.

## ALG-TREE-022 RightSideView
- Last value of each level in level order (CLRS §22.2). Θ(n), O(w) extra; appends to dst.
- Worked: shared tree → 1 3 7 6.

## Fenwick — ALG-TREE-023 NewFenwickTree · 024 Add · 025 PrefixSum · 026 RangeSum · 027 Len
## ALG-TREE-028 NewFenwickRangeUpdate · 029 RangeAdd · 030 Get · 031 Len
- **Source**: Fenwick 1994; Laaksonen §9.2 (tree[k] = sum of (k − p(k), k], p(k) = k & −k) and §9.4
  "Range updates" (the tree over the difference array: range add = two point adds, point value = prefix).
- **Best known**: O(log n) update and prefix query; Ω(log n / log log n) cell-probe lower bound for
  dynamic prefix sums (Fredman–Saks 1989, not re-read; recorded as a pointer). Build O(n): each node adds
  itself to its parent k + p(k) once (follows from the definition).
- **Decision**: range-update + point-query variant built (sourced, §9.4). Range-update + range-query by two
  Fenwick trees **not built** (no primary source read for it); that need is met by the lazy segment tree
  (ALG-TREE-039…045), sourced in §28.1.
- RangeSum walks the two index chains until they meet (equivalent to PrefixSum(hi) − PrefixSum(lo), fewer
  terms, fewer spurious overflows).
- T = signed integers and floats (the range variant subtracts). ErrRange, ErrNotFinite, ErrOverflow
  (a partial sum the tree forms does not fit T; can happen for a range whose total fits — use a wider T).
- Worked: Laaksonen p. 87: [1 3 4 8 6 1 4 2] → tree 1 4 4 16 6 7 4 29; sum of the first 7 = 27.
  §9.4 p. 93–94: [3 3 1 1 1 5 2 2], add 5 to positions 1..4 → 3 8 6 6 6 5 2 2.

## Segment trees — ALG-TREE-032 NewSegmentTree · 033 NewSumSegmentTree · 034 NewMinSegmentTree
## 035 NewMaxSegmentTree · 036 Set · 037 Query · 038 Len
## Lazy — ALG-TREE-039 NewLazySegmentTree · 040 NewRangeAddSumTree · 041 NewRangeAddMinTree
## 042 NewRangeAddMaxTree · 043 RangeUpdate · 044 Query · 045 Len
- **Source**: Laaksonen §9.3 (segment tree; any associative operation), §28.1 (lazy propagation: node
  value plus a pending update, combined updates compose). Generic: a monoid (associative combine +
  identity; combine need not commute) and, for lazy, an update type U with Apply(u, aggregate, length),
  Compose(newer, older) and NoUpdate.
- **Best known**: O(log n) point update and range query, O(log n) range update with lazy propagation;
  O(n) build. Layout: 2n − 1 nodes (left child i+1, right child i + 2·(mid − lo)), not 4n.
  **Lost**: the bottom-up iterative tree over 2n slots (needs commutativity or padding to a power of two
  plus an identity for the padding; min/max ready-mades have no natural identity for every T) and the
  4n heap layout (twice the memory).
- Query on the lazy tree is read-only (pending updates are composed on the way down, not pushed), so
  concurrent readers do not race with each other. Set and RangeUpdate journal the nodes they write and
  roll back on overflow (atomic).
- Empty range [lo, lo): the identity for NewSegmentTree/NewLazySegmentTree/sum trees; ErrEmptyRange for
  min/max ready-mades (no identity for every T).
- ErrRange, ErrNilFunc, ErrNotFinite (number ready-mades), ErrOverflow, ErrEmptyRange.
- Worked: Laaksonen §28.1 pp. 258–259: [5 8 6 3 2 7 2 6 7 1 7 5 6 2 3 2] sums to 72; adding 2 to positions
  5..13 makes the total 90 (root 90/0 in the book's figure) and the left half 45, right half 45.

Tags (LeetCode, pattern only): 94, 144, 145, 102, 104, 110, 235, 236, 105, 106, 108, 230, 98, 297, 124, 199,
307, 315, 450, 701.

---

# algo/trie (SHORT = TRIE) — new package (builder W2-B7)

Planned in doc/taxonomy.md:20 (`algo/trie`, TRIE, "prefix trees"). Files: doc.go, trie.go, autocomplete_top_k.go,
wildcard_search.go, replace_words_by_root.go, prefix_suffix_index.go, each with its _test.go. Imports algo/heap
(AutocompleteTopK). New error: `ErrNegativeK` (autocomplete_top_k.go:12). Exported types: `Trie[V]`
(trie.go:17), `Entry[V]` (autocomplete_top_k.go:15), `PrefixSuffixIndex` (prefix_suffix_index.go:18), constant
`Wildcard = '.'` (wildcard_search.go:6).

## Sources (confirmed)
- Fredkin, "Trie memory", Communications of the ACM, 1960, 3(9):490–499, doi:10.1145/367390.367400 —
  metadata confirmed on Crossref 2026-11-12 (paper itself not read).
- Knuth, TAOCP Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching",
  pp. 492–512 (section and pages confirmed via a search listing; book not opened in this run).
- Manning, Raghavan, Schütze, *Introduction to Information Retrieval*, Cambridge UP, 2008, §3.2.1 "General
  wildcard queries", Permuterm indexes — READ FIRST-HAND in the 2007 preliminary draft (§3.1.1, pp. 48–49):
  "the dictionary consists of all rotations of each term (with the $ terminating symbol appended)"; the query
  m*n is rotated to n$m* and looked up as a prefix, finding man, men and (via n$moro) moron. The same page
  describes the losing alternative: a normal and a reverse B-tree, intersecting the prefix set W and the suffix
  set R. Section number in the 2008 print edition is §3.2.1 (TOC); **lead: confirm the print page**.
- CLRS 3rd ed. §6.5 for the bounded priority queue in AutocompleteTopK.

## Design decisions (each with the alternative that lost)
- **Children**: sorted byte labels + child pointers per node, binary search (≤ 8 comparisons, σ ≤ 256).
  Lost: `[256]*node` per node (O(1) step but 2 KiB per node: 2 GiB for 1e6 nodes); a map per node (no order,
  so every walk would sort). Sorted children give byte-order iteration for free → every output deterministic.
- **Byte semantics** everywhere (keys, '.', prefix/suffix). Lost: rune semantics (needs UTF-8 decoding on
  every step and an invalid-UTF-8 policy); a caller wanting runes can store one key per normalised string.
- **Iteration**: `iter.Seq2[string, V]` for WithPrefix and WildcardSearch (lazy, stops early); iterative walks
  with an explicit stack, so a 1 MB key never overflows the goroutine stack.
- **Zero value usable**; nil *Trie acts empty (Insert stores nothing), as Heap.
- **ReplaceWordsByRoot ignores empty roots** (ASSUMPTION: "" is a prefix of everything and would erase every
  word). **Lead: log in doc/ASSUMPTIONS.md** or tell me the other policy.

### ALG-TRIE-001…006 Trie.Insert / Search / HasPrefix / Delete / Len / WithPrefix (trie.go:65, 96, 112, 131, 173, 192)
- **Best known**: Θ(m) per operation for a key of m bytes (must read the key); our O(m·log σ) = O(m) with
  σ ≤ 256. WithPrefix O(m + s), s = nodes below the prefix (output-sensitive up to key lengths).
  Lost: sorted slice + binary search (O(m log n) search, O(n) insert); hash map (no prefix queries).
- **Delete**: one walk down remembering the deepest ancestor that must stay (root, holds a key, or has
  another child); cut that one edge. O(m), O(1) extra; no parent pointers, no recursion. Invariant: every
  non-root node leads to a key, so HasPrefix is "node exists". Tests check node counts after pruning
  (trie_test.go TestTrie_deletePrunes, and the oracle compares node count with the number of distinct
  prefixes of the remaining keys).
- **Worked example**: Knuth's 31 common English words (§6.2.2/§6.3); expected results derived by hand
  (the list was not re-checked against the book; noted in trie_test.go).
- **Errors**: none. Input not modified. **Tags**: trie, prefix-tree. LeetCode 208.

### ALG-TRIE-007 Trie.AutocompleteTopK (autocomplete_top_k.go:38)
`func (t *Trie[V]) AutocompleteTopK(dst []Entry[V], prefix string, k int, rank func(a, b V) int) ([]Entry[V], error)`
- **Decision**: ranked by a caller comparator on the payload (cmp.Compare shape), ties broken by byte order
  of the key (smaller first); rank nil = byte order alone with early stop after k. Lost: a `weight func(V)
  float64` (forces floats, NaN policy) and storing per-node precomputed top-k lists (O(1) query, but
  O(k·depth) extra per insert and fixed k).
- **Best known**: O(m + s log k) with a k-bounded heap over the subtree; O(m + k·L) is possible only with
  per-node cached lists (lost, above). Because the walk is in byte order, a later key loses every tie, so it
  enters the heap only if strictly better: one rank call per key and a string allocation only on entry.
- **Errors**: ErrNegativeK (k < 0). k = 0 → nothing. Huge k never sizes the heap (hint min(k, Len)).
- **Tags**: autocomplete, typeahead, top-k. LeetCode 1268, 642.

### ALG-TRIE-008 Trie.WildcardSearch (wildcard_search.go:23)
- '.' (exported `Wildcard`) matches exactly one BYTE; no escape (a literal '.' in a key is matched by '.'
  anyway); match is exact length. Iterator in byte order; existence = range once and stop.
- **Best known**: O(m) without '.', O(nodes at depth ≤ m) worst — the trie walk is optimal among
  index-free methods. Lost: scanning all keys (O(total bytes)); recursion (depth = m).
- **Tags**: wildcard, dont-care, partial-match. LeetCode 211.

### ALG-TRIE-009 Trie.ShortestPrefixOf, ALG-TRIE-010 ReplaceWordsByRoot (replace_words_by_root.go:15, 47)
- ShortestPrefixOf stops at the first node holding a key: O(p), returns word[:p] (no allocation, tested).
- ReplaceWordsByRoot(dst, roots, words) builds a Trie[struct{}] of roots: O(R + W). Lost: for each word,
  test every root (O(W·|roots|)); hash set of roots probing each prefix length (O(Σ L²) hashing).
- **Tags**: replace-words, stemming. LeetCode 648.

### ALG-TRIE-011…013 NewPrefixSuffixIndex / Lookup / LookupLast (prefix_suffix_index.go:49, 117, 145)
- **Method**: permuterm (Manning et al. §3.2.1) with one change, proved here: the book stores rotations
  w[j:] $ w[:j] and answers X*Y by the prefix Y$X, which requires len(w) ≥ |X|+|Y| (no overlap). We store
  w[j:] $ w (the whole word after $), so for a query (p, s): a key s' $ w has prefix s $ p ⇔ s' = s (the
  separator, symbol 256, is not a byte, so it pins the end of s) and w starts with p; s' = w[j:] = s ⇔ s is a
  suffix of w. Hence matches = {w : HasPrefix(w,p) ∧ HasSuffix(w,s)}, overlap allowed, and each match has
  exactly one key below the query node (j = len(w) − len(s)), so no deduplication.
  Same O(Σ L²) space as the rotations. The separator is a 257th symbol (uint16 labels), not '#' or '$', so
  words containing any byte are safe (tested).
- **Best known / lost**: query O(p + q) for LookupLast (each node keeps the largest word index through it);
  Lookup O(p + q + s + r log r). Lost: two tries (prefix and reversed) and intersecting W ∩ R — the book's
  own alternative — costs O(|W| + |R|) per query even when the answer is empty; brute force O(Σ L) per query.
- **Worked example**: the book's m*n → man, men, moron (prefix_suffix_index_test.go).
- **Errors**: none; nil/zero index is empty; LookupLast returns -1 when nothing matches.
- **Tags**: permuterm, prefix-and-suffix, word-filter. LeetCode 745.

## Consumer evidence (read-only)
- /Users/devinanderson/Desktop/red_lite_cache/internal/cache/ordered.go:209-231 `appendPrefix`/`cutPrefix`:
  a prefix group walk over a sorted B+-tree of keys. Trie.WithPrefix gives the same ordered prefix
  enumeration, but **does not replace it**: the cache's ordered tree also serves range scans, ordered
  eviction and cut-in-place of a prefix run; a trie would be a second index. A Trie-backed DeletePrefix
  would need a "cut subtree" method (not built; would be a new ID if asked).
- /Users/devinanderson/Desktop/red_lite_cache/internal/sheet/sweeplog.go:17-24 `covers()`: checks one key
  against one logged prefix or suffix (strings.HasPrefix / HasSuffix). One check per entry is already O(len);
  **not replaced**. If the log grows large and the question becomes "which logged suffix-deletes cover this
  key", a Trie of reversed suffixes with ShortestPrefixOf (or a walk) answers it in O(len(key)).
- a consumer's cache key-prefix constants (cache/keys.go:60),
  :159 and :311 (two more prefix constants): a ':'-delimited key grammar swept by DeletePrefix. Motivates
  HasPrefix/WithPrefix (find what a sweep would cover before running it); no code change proposed.

---

# algo/graph (SHORT = GRAPH)

Builder W2-B8. New package (in `doc/taxonomy.md`: `ALG | algo/graph | GRAPH | traversal, shortest path,
topological order`). Imports `algo/unionfind` (ConnectedComponents, Kruskal, and the Cell type) and
`algo/heap` (Dijkstra, Prim).

| ID | Name | Signature |
|---|---|---|
| ALG-GRAPH-001 | NewDirected | `func NewDirected(n int, edges []Edge) (*Graph, error)`; `type Edge struct{ From, To int; Weight int64 }` |
| ALG-GRAPH-002 | NewUndirected | `func NewUndirected(n int, edges []Edge) (*Graph, error)` |
| ALG-GRAPH-003 | Graph.AddEdge | `func (g *Graph) AddEdge(from, to int, weight int64) error` |
| ALG-GRAPH-004 | Graph.VertexCount | `func (g *Graph) VertexCount() int` |
| ALG-GRAPH-005 | Graph.EdgeCount | `func (g *Graph) EdgeCount() int` |
| ALG-GRAPH-006 | Graph.Directed | `func (g *Graph) Directed() bool` |
| ALG-GRAPH-007 | Graph.Neighbors | `func (g *Graph) Neighbors(v int) iter.Seq2[int, int64]` |
| ALG-GRAPH-008 | Graph.Edges | `func (g *Graph) Edges() []Edge` |
| ALG-GRAPH-009 | BFS | `func BFS(g *Graph, source int) (iter.Seq[int], error)` |
| ALG-GRAPH-010 | DFS | `func DFS(g *Graph, source int) (iter.Seq[int], error)` |
| ALG-GRAPH-011 | ShortestPathUnweighted | `func ShortestPathUnweighted(g *Graph, source, target int) ([]int, error)` |
| ALG-GRAPH-012 | Dijkstra | `func Dijkstra(g *Graph, source int) (dist []int64, parent []int, err error)` |
| ALG-GRAPH-013 | BellmanFord | `func BellmanFord(g *Graph, source int) (dist []int64, parent []int, err error)` |
| ALG-GRAPH-014 | FloydWarshall | `func FloydWarshall(g *Graph) ([][]int64, error)` |
| ALG-GRAPH-015 | TopologicalSort | `func TopologicalSort(g *Graph) ([]int, error)` |
| ALG-GRAPH-016 | ConnectedComponents | `func ConnectedComponents(g *Graph) (component []int, count int, err error)` |
| ALG-GRAPH-017 | IsBipartite | `func IsBipartite(g *Graph) (side []bool, ok bool, err error)` |
| ALG-GRAPH-018 | Kruskal | `func Kruskal(g *Graph) (tree []Edge, totalWeight int64, err error)` |
| ALG-GRAPH-019 | Prim | `func Prim(g *Graph) (tree []Edge, totalWeight int64, err error)` |
| ALG-GRAPH-020 | CloneGraph | `func CloneGraph(g *Graph) *Graph` |
| ALG-GRAPH-021 | WordLadder | `func WordLadder(start, target string, dictionary []string) []string` |
| ALG-GRAPH-022 | GridBFS | `func GridBFS[T any](grid [][]T, start Cell, passable func(T) bool, nb Neighborhood) (iter.Seq[Cell], error)` |
| ALG-GRAPH-023 | GridDFS | `func GridDFS[T any](grid [][]T, start Cell, passable func(T) bool, nb Neighborhood) (iter.Seq[Cell], error)` |
| ALG-GRAPH-024 | NumberOfIslands | `func NumberOfIslands[T any](grid [][]T, isLand func(T) bool, nb Neighborhood) (int, error)` |
| ALG-GRAPH-025 | FloodFillInPlace | `func FloodFillInPlace[T comparable](grid [][]T, start Cell, replacement T, nb Neighborhood) (int, error)` |
| ALG-GRAPH-026 | SurroundedRegionsInPlace | `func SurroundedRegionsInPlace[T comparable](grid [][]T, region, fill T, nb Neighborhood) (int, error)` |

`type Arc` is not exported: Neighbors yields (to, weight). `type Cell = unionfind.Cell` (alias); `type
Neighborhood int` with `FourNeighbors` (up, right, down, left) and `EightNeighbors` (clockwise from up).
`const Unreachable int64 = math.MaxInt64` (the distance of a vertex with no path).

## Sources (DOIs confirmed on Crossref, api.crossref.org, 2026-11-12)
- **CLRS** = Cormen, Leiserson, Rivest, Stein, *Introduction to Algorithms*, 3rd ed., MIT Press, 2009.
- Dijkstra, "A note on two problems in connexion with graphs", Numerische Mathematik 1(1):269–271, 1959,
  doi:10.1007/BF01386390 (problem 2 is single-source shortest paths; problem 1 is the spanning tree).
- Bellman, "On a routing problem", Quarterly of Applied Mathematics 16(1):87–90, 1958, doi:10.1090/qam/102435.
  Ford, "Network flow theory", RAND Corporation paper P-923, 1956 (no DOI; cited through CLRS ch. 24 notes).
- Floyd, "Algorithm 97: Shortest path", Communications of the ACM 5(6):345, 1962, doi:10.1145/367766.368168.
  Warshall, "A theorem on Boolean matrices", Journal of the ACM 9(1):11–12, 1962, doi:10.1145/321105.321107.
- Kahn, "Topological sorting of large networks", Communications of the ACM 5(11):558–562, 1962,
  doi:10.1145/368996.369025.
- Kruskal, "On the shortest spanning subtree of a graph and the traveling salesman problem", Proceedings of
  the American Mathematical Society 7(1):48–50, 1956, doi:10.1090/S0002-9939-1956-0078686-7.
- Prim, "Shortest connection networks and some generalizations", Bell System Technical Journal
  36(6):1389–1401, 1957, doi:10.1002/j.1538-7305.1957.tb01515.x.
- Tarjan, "Depth-first search and linear graph algorithms", SIAM Journal on Computing 1(2):146–160, 1972,
  doi:10.1137/0201010.
- Smith, "Tint fill", Proceedings of SIGGRAPH '79 (6th annual conference on computer graphics and
  interactive techniques), 1979, pp. 276–283, doi:10.1145/800249.807456.
- Knuth, *The Stanford GraphBase: A Platform for Combinatorial Computing*, ACM Press / Addison-Wesley, 1993:
  GB_WORDS (adjacent words differ in one position; built with "five hash tables, one for each of the strings
  of four letters obtained by suppressing a single letter", section "The output phase") and LADDERS (shortest
  word ladders). Checked against the CWEB source gb_words.w on CTAN, 2026-11-12.
- The Go Programming Language Specification, "For statements with range clause" (range over function
  iterators, Go 1.23+), go.dev/ref/spec; package iter (iter.Seq, iter.Seq2).

## Decisions (lead to log)
- **Weights: int64**, exact. Every path sum is checked: a sum that does not fit, or that equals
  `Unreachable`, is `ErrOverflow`. **Lost**: float64 with NaN/Inf rejected (sums lose exactness above 2⁵³,
  results depend on summation order, and a "no path" sentinel of +Inf mixes with real values). Callers with
  fractional weights scale to integers (cents, millimetres).
- **One Graph type, two constructors** (`NewDirected`, `NewUndirected`); unweighted graphs use weight 0
  (BFS-based functions ignore weights). Edges are validated before anything is stored (all or nothing).
  Parallel edges and self-loops are allowed (an undirected self-loop is stored once). The adjacency rows
  share one backing array, each clipped to its degree, so AddEdge never aliases another row. **Lost**: a
  `directed bool` parameter (unreadable at the call site); separate weighted/unweighted types (doubles every
  algorithm's API).
- **Huge n (open, for the lead)**: NewDirected/NewUndirected allocate n list headers; an n beyond what the
  runtime can allocate makes `make` panic (as heap.NewHeap's capacityHint and unionfind.NewDisjointSet(n)
  do). No cap was invented (no source for one); the same policy question as L13.
- **Nil graph** → `ErrNilGraph` from every algorithm; methods on a nil *Graph return zero values.
- **Directed vs undirected**: TopologicalSort needs directed (`ErrUndirected`); IsBipartite, Kruskal and Prim
  need undirected (`ErrDirected`); ConnectedComponents on a directed graph returns WEAKLY connected
  components (edge direction ignored), stated in the doc.
- **ConnectedComponents calls unionfind** (one Union per edge, O(V + E α(V))). **Lost**: BFS per component
  (O(V + E) but needs reverse arcs for the weak components of a directed graph; union-find handles both
  kinds with one code path and is the CLRS §21.1 method).
- **Grids**: `[][]T` must be rectangular (`ErrNotRectangular`); rows×cols that overflow int is `ErrTooLarge`.
  **FloodFillInPlace and SurroundedRegionsInPlace modify the grid** (names say so). **Lost**: dst copies
  (O(rows×cols) extra per call on images that are usually edited in place; a caller who wants a copy clones
  first). Both return the number of cells changed.
- **Traversals return iter.Seq** built lazily: each range over the sequence runs a fresh traversal, so a
  sequence may be ranged twice; breaking out stops the work. Argument errors are returned up front, before
  any iteration.
- Dijkstra and Prim use `algo/heap` with lazy deletion (push the new key, skip stale entries on pop):
  O((V + E) log E) = O((V + E) log V) time, O(V + E) heap entries. **Lost**: decrease-key heap (heap has
  none; an indexed heap would give O((V+E) log V) with O(V) heap space — same bound); Fibonacci heap
  (O(E + V log V), Fredman–Tarjan 1987, not cited here, high constants).

### ALG-GRAPH-001…008 Graph (adjacency list)
- **Source**: CLRS §22.1 (adjacency-list representation; Θ(V + E) memory; undirected edges stored in both lists).
- **Best known**: O(V + E) build and memory (reading the input is Ω(V + E)).
- **Worked example**: CLRS §22.1 Figure 22.1(a,b): the undirected graph on 1..5 (here 0..4) with edges 1–2,
  1–5, 2–5, 2–3, 2–4, 3–4, 4–5: degrees 2 4 2 3 3, 7 edges. Test TestNewUndirected_clrsFigure22_1.
- **Errors**: `ErrNegativeVertexCount` (n < 0), `ErrVertexOutOfRange` (an endpoint outside 0..n−1, wrapped
  with the edge index), `ErrNilGraph` (AddEdge on nil). Input slices not modified.
- **Tests**: Figure 22.1; directed/undirected, self-loop, parallel edges, n 0; Edges round trip; AddEdge
  does not alias other rows; errors. Benchmarks NewDirected/NewUndirected/AddEdge/Neighbors/Edges at
  1e3/1e4/1e5 edges; VertexCount/EdgeCount/Directed single case.
- **Tags**: graph, adjacency-list.

### ALG-GRAPH-009 BFS · 010 DFS
- **Source**: CLRS §22.2 (BFS), §22.3 (DFS); Tarjan 1972 (DFS, linear time). Go spec range-over-func.
- **Best known**: Θ(V + E) time, O(V) extra.
- **Approach**: BFS with a slice queue (head index); DFS iterative with an explicit stack of (vertex, next
  arc) frames — the same preorder as the recursive DFS-VISIT visiting arcs in list order, with no recursion
  depth limit. **Lost**: recursive DFS (stack depth V; a path of 10⁷ vertices would grow the goroutine
  stack to GBs); the "push all neighbours" stack DFS (different order than DFS-VISIT, O(E) stack).
- **Worked example**: CLRS Figure 22.3 (BFS from s on the 8-vertex undirected graph r s t u v w x y):
  distances r1 s0 t2 u3 v2 w1 x2 y3. CLRS Figure 22.4 (DFS on the directed graph u v w x y z, edges
  u→v u→x v→y w→y w→z x→v y→x z→z): DFS-VISIT(u) discovers u v y x (then the book restarts at w: w z);
  our single-source DFS from u yields u v y x, from w yields w y x v z. Tests TestBFS_clrsFigure22_3, TestDFS_clrsFigure22_4.
- **Errors**: `ErrNilGraph`, `ErrVertexOutOfRange`.
- **Tests**: figures; one vertex; disconnected; early break; re-ranging; oracle against a recursive DFS
  and a level-by-level BFS on 300 random graphs. Benchmarks at V = 1e3/1e4/1e5 (E = 4V).
- **Tags**: bfs, dfs, traversal, iterator. LeetCode: 797 (DFS paths, pattern only).

### ALG-GRAPH-011 ShortestPathUnweighted
- **Source**: CLRS §22.2 (BFS computes shortest-path distances, Theorem 22.5; PRINT-PATH reconstruction).
- **Best known**: O(V + E), O(V) extra. Bidirectional BFS has the same worst case (lost: two frontiers, more code).
- **Output**: vertex path source…target; source == target → [source]; unreachable → nil, nil.
- **Worked example**: CLRS Figure 22.3: the shortest s → y path is s w x y (3 edges, δ(s, y) = 3; it is unique).
- **Errors**: `ErrNilGraph`, `ErrVertexOutOfRange`. **Tags**: bfs, shortest-path, unweighted. LeetCode: 1091.

### ALG-GRAPH-012 Dijkstra
- **Source**: Dijkstra 1959 (problem 2); CLRS §24.3.
- **Best known**: O(E + V log V) with a Fibonacci heap; binary heap O((V + E) log V). Lower bound Ω(V log V)
  for comparison-based (it sorts). Chosen: binary heap with lazy deletion, O((V + E) log E), O(V + E) extra.
- **Worked example**: CLRS Figure 24.6 (s t x y z; edges s→t10 s→y5 t→x1 t→y2 y→t3 y→x9 y→z2 x→z4 z→x6
  z→s7): δ = s0 t8 x9 y5 z7. Test TestDijkstra_clrsFigure24_6.
- **Errors**: `ErrNilGraph`, `ErrVertexOutOfRange`, `ErrNegativeWeight` (any arc in the graph), `ErrOverflow`.
- **Tests**: figure; unreachable; zero weights; parallel edges; oracle vs BellmanFord on 300 random graphs.
  Benchmark V = 1e3/1e4/1e5, E = 4V.
- **Tags**: dijkstra, shortest-path, weighted, heap. LeetCode: 743, 1514.

### ALG-GRAPH-013 BellmanFord
- **Source**: Bellman 1958; Ford 1956; CLRS §24.1 (negative-cycle check after V−1 passes).
- **Best known**: O(V·E) for general weights; chosen O(V·(V + E)) with early exit when a pass relaxes
  nothing. **Lost**: SPFA queue variant (same worst case, no cleaner source).
- **Worked example**: CLRS Figure 24.4 (s t x y z; s→t6 s→y7 t→x5 t→y8 t→z−4 x→t−2 y→x−3 y→z9 z→s2 z→x7):
  δ = s0 t2 x4 y7 z−2. Test TestBellmanFord_clrsFigure24_4.
- **Errors**: `ErrNilGraph`, `ErrVertexOutOfRange`, `ErrNegativeCycle` (reachable from source), `ErrOverflow`
  (an intermediate distance outside int64; can precede detection of a negative cycle).
- **Tags**: bellman-ford, shortest-path, negative-weights, negative-cycle. LeetCode: 787 (pattern only).

### ALG-GRAPH-014 FloydWarshall
- **Source**: Floyd 1962; Warshall 1962; CLRS §25.2.
- **Best known**: Θ(V³) for dense all-pairs (Johnson O(V E + V² log V) on sparse — lost: needs Bellman-Ford
  + V Dijkstras, more memory traffic; a separate function later if needed). Space Θ(V²) result, O(1) extra.
- **Worked example**: CLRS Figure 25.4 (1..5: 1→2 3, 1→3 8, 1→5 −4, 2→4 1, 2→5 7, 3→2 4, 4→1 2, 4→3 −5,
  5→4 6): D⁽⁵⁾ rows 0 1 −3 2 −4 / 3 0 −4 1 −1 / 7 4 0 5 3 / 2 −1 −5 0 −2 / 8 5 1 6 0. Test TestFloydWarshall_clrsFigure25_4.
- **Errors**: `ErrNilGraph`, `ErrNegativeCycle` (any d[i][i] < 0, detected the moment it appears), `ErrOverflow`.
  A guard returns `ErrTooLarge` when V² overflows int (no graph that fits in memory reaches it; untested).
  V is the caller's to bound: V = 10⁴ allocates 800 MB.
- **Tags**: floyd-warshall, all-pairs, shortest-path, dense. LeetCode: 1334.

### ALG-GRAPH-015 TopologicalSort
- **Source**: Kahn 1962; CLRS §22.4 (the DFS method — lost: Kahn detects a cycle without a 3-colour DFS
  and its output order is the BFS-like "sources first" order).
- **Best known**: Θ(V + E), O(V) extra (the result slice doubles as the queue).
- **Order**: initial sources in increasing vertex order, then FIFO; deterministic.
- **Worked example**: CLRS Figure 22.7 (Professor Bumstead's clothes: undershorts→pants, undershorts→shoes,
  pants→belt, pants→shoes, belt→jacket, shirt→belt, shirt→tie, tie→jacket, socks→shoes; watch isolated):
  any output must respect every edge; test checks validity and Kahn's exact order.
- **CourseSchedule** (prerequisite pairs → can all be finished) is a TEST of TestTopologicalSort_courseSchedule,
  not a function.
- **Errors**: `ErrNilGraph`, `ErrUndirected`, `ErrCycle` (including a self-loop).
- **Tags**: topological-sort, kahn, dag, cycle-detection, scheduling. LeetCode: 207, 210.

### ALG-GRAPH-016 ConnectedComponents
- **Source**: CLRS §21.1 CONNECTED-COMPONENTS (disjoint sets); calls algo/unionfind.
- **Best known**: Θ(V + E) (BFS); chosen O(V + E α(V)), O(V) extra. Labels 0..count−1 in order of each
  component's smallest vertex.
- **Worked example**: CLRS Figure 21.1 (a..j, 4 components). **Errors**: `ErrNilGraph`.
- **Tags**: connected-components, union-find. LeetCode: 323, 547.

### ALG-GRAPH-017 IsBipartite
- **Source**: CLRS Exercise 22.2-7 (split wrestlers into two kinds with every rivalry across: BFS 2-colouring,
  O(V + E)); Sedgewick & Wayne, *Algorithms*, 4th ed., 2011, §4.1 (two-colourability).
- **Best known**: Θ(V + E), O(V) extra. side[v] false/true for the two parts (each component's smallest
  vertex is false); side nil when not bipartite.
- **Worked example**: hand: a 4-cycle is bipartite {0,2}/{1,3}; a triangle is not; a self-loop is not.
- **Errors**: `ErrNilGraph`, `ErrDirected`. **Tags**: bipartite, two-coloring, bfs. LeetCode: 785, 886.

### ALG-GRAPH-018 Kruskal · 019 Prim
- **Source**: Kruskal 1956; Prim 1957; CLRS §23.2 (MST-KRUSKAL with disjoint sets; MST-PRIM).
- **Best known**: O(E log V) for both with binary heaps / sorting; O(E α(V)) after sorting. Lower bound
  Ω(E). (Randomised linear-time Karger–Klein–Tarjan, lost: complex, not needed.)
- **Approach**: Kruskal sorts a copy of the edges stably by weight and unions with algo/unionfind, stopping
  at V − 1 tree edges. Prim restarts from each unvisited vertex, lazy heap. Both return a minimum spanning
  FOREST when the graph is disconnected; negative weights are allowed.
- **Worked example**: CLRS Figure 23.1/23.4 (a..i; a–b4 a–h8 b–c8 b–h11 c–d7 c–f4 c–i2 d–e9 d–f14 e–f10
  f–g2 g–h1 g–i6 h–i7): total weight 37, 8 edges. Tests TestKruskal_clrsFigure23_4, TestPrim_clrsFigure23_5.
- **Errors**: `ErrNilGraph`, `ErrDirected`, `ErrOverflow` (total weight).
- **Tests**: oracle Kruskal total == Prim total == brute force over all edge subsets on small random graphs.
- **Tags**: mst, kruskal, prim, union-find, greedy. LeetCode: 1584.

### ALG-GRAPH-020 CloneGraph
- **Source**: CLRS §22.1 (adjacency lists; copying each list). Deep copy: the clone shares no memory with g.
- **Best known**: Θ(V + E). Two allocations plus the struct. nil → nil.
- **Tags**: graph, clone, deep-copy. LeetCode: 133.

### ALG-GRAPH-021 WordLadder
- **Source**: Knuth, *The Stanford GraphBase*, 1993 (GB_WORDS adjacency by "suppressing a single letter"
  hash tables; LADDERS shortest ladders); BFS on the implicit graph, CLRS §22.2.
- **Best known**: O(N·L²) time and space with the masked-key buckets (N words of length L: N·L keys of L−1
  bytes); each bucket is scanned once then dropped. **Lost**: trying all 26 (or 256) letters per position
  (O(N·L·Σ·L), alphabet-dependent); bidirectional BFS (same worst case; optional in the brief, not built).
- **Output**: one shortest ladder start…target (deterministic for the same input); nil when none; [start]
  when start == target. target must be in dictionary; start need not be. Letters are bytes (multi-byte UTF-8
  letters are not one-letter changes). Words of another length are ignored.
- **Worked example**: Knuth's demo pair is not reproduced; hand example: cold → warm differ in all 4
  letters, so ≥ 5 words; with dictionary cord card ward warm word worm wore core (in that order) BFS returns
  cold cord word ward warm.
- **Tags**: word-ladder, bfs, implicit-graph, strings. LeetCode: 127, 126.

### ALG-GRAPH-022 GridBFS · 023 GridDFS · 024 NumberOfIslands · 025 FloodFillInPlace · 026 SurroundedRegionsInPlace
- **Source**: BFS/DFS on the grid graph (vertices = cells, edges = 4- or 8-neighbours), CLRS §22.2, §22.3;
  Smith 1979 (flood "tint fill" of a 4- or 8-connected region of one value). NumberOfIslands = connected
  components of the passable subgraph (CLRS §22.3 / §21.1). SurroundedRegions = cells of `region` whose
  component touches no border cell: BFS from all border cells of that value at once (a BFS from a virtual
  source joined to every border cell, CLRS §22.2), then fill the rest.
- **Best known**: Θ(R·C) time, O(R·C) extra (the visited bitmap / queue). FloodFillInPlace needs no visited
  array: changed cells no longer match.
- **Worked examples** (hand): a 4×5 grid with 4 islands under 4-neighbours that join into 1 under
  8-neighbours; GridBFS/GridDFS orders on a 3×3 maze (4 and 8 neighbours); a 3×3 flood fill (5 cells from
  the centre under 4, 3 cells from a corner under 8); a 4×4 O/X board where the inner O block (3 cells) is
  filled under 4-neighbours and nothing under 8 (a diagonal reaches the border O).
- **Errors**: `ErrNotRectangular`, `ErrTooLarge`, `ErrCellOutOfRange`, `ErrNilPredicate`, `ErrNeighborhood`.
- **Tests**: hand cases, 1×1, empty grid, all water, all land, 8 vs 4, oracle vs unionfind.IslandsOnline's
  last count for NumberOfIslands (4-neighbours) on 300 random grids. Benchmarks on square grids of
  1e4/1e5/1e6 cells (the n of the 1e3/1e4/1e5 rule, ×10 because one cell is cheap).
- **Tags**: grid, bfs, dfs, flood-fill, islands, connected-components. LeetCode: 200, 733, 130, 1020.

---

# algo/unionfind (SHORT = UNION)

Builder W2-B8. New package (in `doc/taxonomy.md`: `ALG | algo/unionfind | UNION | disjoint sets`).
Built before `algo/graph`: graph's ConnectedComponents and Kruskal call it.

| ID | Name | Signature |
|---|---|---|
| ALG-UNION-001 | NewDisjointSet | `func NewDisjointSet(n int) (*DisjointSet, error)` |
| ALG-UNION-002 | DisjointSet.MakeSet | `func (s *DisjointSet) MakeSet() int` |
| ALG-UNION-003 | DisjointSet.Find | `func (s *DisjointSet) Find(x int) (int, error)` |
| ALG-UNION-004 | DisjointSet.Union | `func (s *DisjointSet) Union(a, b int) (merged bool, err error)` |
| ALG-UNION-005 | DisjointSet.Connected | `func (s *DisjointSet) Connected(a, b int) (bool, error)` |
| ALG-UNION-006 | DisjointSet.Count | `func (s *DisjointSet) Count() int` |
| ALG-UNION-007 | DisjointSet.Len | `func (s *DisjointSet) Len() int` |
| ALG-UNION-008 | IslandsOnline | `func IslandsOnline(dst []int, rows, cols int, additions []Cell) ([]int, error)`; `type Cell struct{ Row, Col int }` |

## Sources (all DOIs confirmed on Crossref, api.crossref.org, 2026-11-12)
- Galler, Fischer, "An improved equivalence algorithm", Communications of the ACM 7(5):301–303, 1964,
  doi:10.1145/364099.364331 (the parent-pointer forest; union of trees by root). Crossref lists the second
  author as "Michael J. Fisher"; the paper is usually cited as Fischer.
- Tarjan, "Efficiency of a good but not linear set union algorithm", Journal of the ACM 22(2):215–225, 1975,
  doi:10.1145/321879.321884 (union by rank (weight) with path compression: O(m α(m, n)) for m operations on n
  elements, and the matching lower bound for this class of algorithms).
- Cormen, Leiserson, Rivest, Stein, *Introduction to Algorithms*, 3rd ed., MIT Press, 2009, ch. 21:
  §21.1 (operations MAKE-SET, UNION, FIND-SET; CONNECTED-COMPONENTS / SAME-COMPONENT, Figure 21.1),
  §21.3 (union by rank, path compression; MAKE-SET, UNION, LINK, FIND-SET pseudocode), §21.4 (the α bound).

## Decisions (lead to log)
- **Element type**: `int` 0..n−1, parent `[]int`, rank `[]uint8` (rank ≤ ⌊log₂ n⌋ ≤ 63, CLRS Lemma 21.4/
  Corollary 21.5). **Lost**: generic `DisjointSet[T comparable]` over a map (a hash per Find, 5–10× slower;
  callers with keys map them to indices once). **Lost**: int32 parents (halves memory, caps n at 2³¹).
- **Union by rank + full path compression (two-pass)**. **Lost**: path halving / splitting (same α bound,
  one pass; Tarjan & van Leeuwen 1984 — not cited here because not used); union by size (same bound; size
  is useful only when callers ask for set sizes, which no brief asked for).
- **Errors, not panics**: Find/Union/Connected return `ErrOutOfRange` for an element outside 0..Len()−1.
  **Lost**: (root, ok bool) — the wrapped error names the bad element.
- **MakeSet** (CLRS's MAKE-SET) grows the structure by one element; the zero DisjointSet is an empty,
  usable structure. Needed by IslandsOnline (elements created as land appears).
- **IslandsOnline storage**: a map from Cell to element id, elements created only for added land cells:
  O(k) memory for k additions whatever rows×cols is. **Lost**: a dense rows×cols DisjointSet (faster per
  operation, but a 10⁵×10⁵ grid would allocate 10¹⁰ entries for a few additions; a size cap would be an
  invented constant).
- **Cell** lives here (the lower layer); `algo/graph` re-exports it as `graph.Cell = unionfind.Cell`
  (a type alias: one type, one home).

### ALG-UNION-001…007 DisjointSet
- **Best known**: m operations on n elements in O(m α(m, n)) with union by rank + path compression; Tarjan
  1975 proves this tight for the class of separable pointer algorithms (Fredman–Saks 1989 extend the lower
  bound to the cell-probe model — pointer only, not used). Find/Union/Connected: O(α(n)) amortized;
  NewDisjointSet O(n); MakeSet O(1) amortized; Count/Len O(1). Space O(n).
- **Approach**: CLRS §21.3 pseudocode, iterative Find (root walk, then a second walk pointing every node
  on the path at the root), so no recursion depth.
- **Worked example**: CLRS §21.1 Figure 21.1: vertices a..j (0..9), edges (b,d) (e,g) (a,c) (h,i) (a,b) (e,f)
  (b,c) → components {a,b,c,d} {e,f,g} {h,i} {j}: Count 4; Connected(a,d) true; Connected(e,h) false;
  Union(b,c) reports false (already joined). Test TestDisjointSet_clrsFigure21_1.
- **Errors**: `ErrNegativeSize` (n < 0), `ErrOutOfRange` (element outside 0..Len()−1). nil *DisjointSet:
  Count/Len 0, MakeSet −1, Find/Union/Connected `ErrOutOfRange`.
- **Input modified**: no inputs; Find compresses paths inside the structure (not safe for concurrent use).
- **Tests**: Figure 21.1; cases (n 0, 1, 2, self-union, repeated union); errors; oracle against a naive
  label array (relabel on union) over 300 random op sequences; rank never exceeds log₂ n; AllocsPerRun 0 for
  Find/Union/Connected. Benchmarks NewDisjointSet, Union, Find, Connected at 1e3/1e4/1e5; MakeSet; Count; Len.
- **Tags**: union-find, disjoint-set, dsu, path-compression, union-by-rank. LeetCode: 547, 684, 721, 990.

### ALG-UNION-008 IslandsOnline
- **Problem**: an m×n grid starts all water; cells become land one at a time; after each addition report the
  number of islands (4-neighbour connected land groups). Adding land that is already land changes nothing.
- **Source**: dynamic connected components with disjoint sets, CLRS §21.1 ("…when the edges of the graph
  are dynamic … the connected components can be maintained faster with disjoint sets"; CONNECTED-COMPONENTS);
  Tarjan 1975 for the bound. The grid is the 4-neighbour grid graph; each addition is MAKE-SET then up to 4
  UNIONs.
- **Best known**: Ω(k) to read k additions; O(k α(k)) with union-find (expected, because of hashing). Space
  O(k). **Lost**: re-running a flood fill after each addition (O(k·m·n)); dense array union-find (see
  Decisions).
- **Worked example** (hand-derived): 3×3, additions (0,0) (0,1) (1,2) (2,1) (1,1) (0,0) → 1 1 2 3 1 1.
  (1,1) joins all three groups; re-adding (0,0) keeps 1.
- **Errors**: `ErrNegativeSize` (rows < 0 or cols < 0), `ErrOutOfRange` (an addition outside the grid,
  wrapped with its index). rows×cols is never multiplied, so no overflow case.
- **Output**: appended to dst[:0]; dst == nil allocates once (len(additions)). Inputs not modified.
- **Tests**: worked example; cases (no additions, 1×1, duplicate additions, row, column, checkerboard
  diagonal stays separate); errors; oracle against a BFS count over the grid after each step (200 random
  runs). Benchmark IslandsOnline/1e3,1e4,1e5 additions on a 1000×1000 grid.
- **Tags**: union-find, grid, islands, online, dynamic-connectivity. LeetCode: 305.

---

# Catalog: algo/dp (SHORT = DP), new package

Builder W3-B9, 2026-11-12. Package in doc/taxonomy.md (`algo/dp`, DP). Conventions (doc.go): inputs never
modified; ordered comparisons by cmp.Compare, matches by == (NaN matches nothing); values int64, counts
uint64, sizes/weights int; every integer accumulation checked (ErrOverflow); partial-state overflow policy
as L28; tables beyond O(input) capped at 128 MiB (ErrTooLarge before allocating; L13, L26). Sentinels live
in checked_arith.go. Citations were checked 2026-11-12 against Crossref and author PDFs (a research agent);
locators marked (print) could not be opened online and should be checked in print by the reviewer.

## Decisions (each with the alternative that lost)

| # | Decision | Alternative | Why it lost |
|---|---|---|---|
| D1 | Names are practitioner terms: MaxNonAdjacentSum(Circular), GridPathCount, TextSegmentation, CanJumpToEnd, MinJumpsToEnd, MaxProfitOneTrade, MaxProfitUnlimitedTrades, SubsequenceOccurrences | the brief's problem-site titles (HouseRobber, UniquePaths, WordBreak, JumpGame, BestTimeToBuySell, DistinctSubsequences) | BRIEF_COMMON: never copy identifiers; "WordBreak" also clashes with Unicode word-boundary breaking |
| D2 | ClimbingStairs not a separate function: it is Fibonacci(n+1), said in the doc comment, doc.go and tags | a second exported function | one function, one home |
| D3 | Fibonacci overflows at n ≥ 94 (uint64), not 93 as the brief says | error at 93 | F(93) = 12200160415121876738 < 2⁶⁴; 93 is the int64 limit's n+1 |
| D4 | CoinChangeCount returns uint64 + ErrOverflow | a big.Int variant | one allocation per cell; no consumer; a Big variant can be added later with its own ID |
| D5 | CoinChangeCount rejects repeated denominations (ErrDuplicate); CoinChangeMin allows them | dedupe silently | a repeat would double-count combinations; for the minimum it is harmless |
| D6 | One 128 MiB table cap for every DP table (1-D and 2-D) | per-function caps | one rule, same as L13; also stops make() panicking (L26) |
| D7 | Knapsack01 returns the items (n·W choice bits) | value only in O(W) | callers want the set; bits are 1/64 of an int table; capacity is first capped at Σ weights |
| D8 | ALG-DP-016 unused: LongestPalindromicSubstring removed (lead L33; home algo/text ALG-TEXT-009); dp/palindrome.go and its test moved to _delete/algo/dp/ | keep an O(1)-space variant | one home |

## Entries

### ALG-DP-001 LongestIncreasingSubsequence
`func LongestIncreasingSubsequence[T cmp.Ordered](values []T) []int`
- **Source**: Fredman, Discrete Mathematics 11(1):29–35, 1975, doi:10.1016/0012-365X(75)90103-X (confirmed);
  Aldous, Diaconis, Bull. AMS 36(4):413–432, 1999, doi:10.1090/S0273-0979-99-00796-X, §1.1 Lemma 1 (piles =
  LIS length; preprint numbering, published numbering (print)).
- **Best known**: O(n log n); Fredman: n log n − n log log n + O(n) comparisons needed. Reading input Ω(n).
- **Approach**: patience sorting with search.LowerBound (ALG-SEARCH-001) on pile tops, back-pointers for one
  subsequence. **Lost**: O(n²) recurrence (slower).
- **Worked example**: deck 7 2 8 1 3 4 10 6 9 5 → 5 piles; indices [3 4 5 7 8] (1 3 4 6 9).
- **Errors**: none. **Tags**: dp, lis, patience-sorting. LeetCode 300, 673.

### ALG-DP-002 LongestCommonSubsequence · ALG-DP-003 LCSLength
`func LongestCommonSubsequence[T comparable](first, second []T) []T` · `func LCSLength[T comparable](first, second []T) int`
- **Source**: Hirschberg, CACM 18(6):341–343, 1975, doi:10.1145/360825.360861 (confirmed); CLRS 3rd ed. §15.4
  eq. (15.9), Exercise 15.4-4 (print).
- **Best known**: O(mn) time (no strongly subquadratic algorithm under SETH); O(min(m, n)) space for the
  length, O(m + n) for the sequence (Hirschberg). **Approach**: one row (length); Hirschberg divide and
  conquer, halving the longer input, depth ⌈log2⌉. **Lost**: full table + traceback (O(mn) memory);
  Hunt–Szymanski (better only when matches are sparse; worse worst case).
- **Worked example**: CLRS Fig. 15.8: ABCBDAB, BDCABA → length 4 (BCBA is one).
- **Errors**: none. **Tags**: dp, lcs, diff. LeetCode 1143, 1092.

### ALG-DP-004 EditDistance · ALG-DP-005 OneEditDistance
`func EditDistance[T comparable](from, to []T) int` · `func OneEditDistance[T comparable](from, to []T) bool`
- **Source**: Wagner, Fischer, JACM 21(1):168–173, 1974, doi:10.1145/321796.321811 (confirmed); Gusfield
  1997 §11.3 (print); Ukkonen, Information and Control 64(1–3):100–118, 1985, doi:10.1016/S0019-9958(85)80046-2
  (testing distance ≤ t in O(t·min(m, n)): confirmed from the abstract).
- **Best known**: O(mn) / O(min) space; for "distance = 1" O(min(m, n)) (Ukkonen band, t = 1), O(1) space.
- **Approach**: one row, shorter input spans it. OneEditDistance: first mismatch, then compare the rest
  (slices.Equal). **Lost**: full table; for OneEditDistance, running EditDistance (O(mn)).
- **Worked examples**: kitten→sitting 3; intention→execution 5 (hand, checked offline).
- **Errors**: none. **Tags**: edit-distance, levenshtein. LeetCode 72, 161.

### ALG-DP-006 Knapsack01 · ALG-DP-007 UnboundedKnapsack
`func Knapsack01(weights []int, values []int64, capacity int) (int64, []int, error)` ·
`func UnboundedKnapsack(weights []int, values []int64, capacity int) (int64, []int, error)`
- **Source**: Kellerer, Pferschy, Pisinger, Knapsack Problems, Springer 2004, doi:10.1007/978-3-540-24777-7,
  §2.3 Dynamic Programming and ch. 8 (§8.3 DP for the unbounded problem) (TOC confirmed); CLRS Exercise 16.2-2.
- **Best known**: NP-hard; O(nW) pseudo-polynomial. **Approach**: Bellman recursion over capacities, one row;
  0-1 records choices in n·W bits; unbounded keeps the last item per capacity. **Lost**: value-indexed DP
  (O(n·V), better only when values are small); branch and bound (exponential worst case).
- **Worked examples**: CLRS Fig. 16.2: (10,60),(20,100),(30,120), W 50 → 220, items [1 2]. Unbounded with the
  CLRS rod prices as items, W 4 → 10 (two of length 2).
- **Errors**: ErrLengthMismatch, ErrNegative, ErrOverflow, ErrTooLarge; ErrUnbounded (zero weight, positive
  value). **Tags**: knapsack. LeetCode 416, 474.

### ALG-DP-008 CoinChangeMin · ALG-DP-009 CoinChangeCount
- **Source**: Wright, "The change-making problem", JACM 22(1):125–128, 1975, doi:10.1145/321864.321874
  (confirmed); CLRS Problem 16-1(d) (assumes a penny; ours does not need one). Count: Graham, Knuth, Patashnik,
  Concrete Mathematics 2nd ed. §7.1 "Domino Theory and Change" (∏ 1/(1 − z^c); 292 ways for a dollar:
  confirmed via secondary sources and recomputed offline); Pólya 1956 doi:10.2307/2309555.
- **Best known**: O(kA) for arbitrary denominations (change-making is weakly NP-hard). **Lost**: greedy
  (wrong for {1,3,4}); BFS over amounts (same bound, more memory).
- **Worked examples**: {1,3,4}, 6 → 2; {1,5,10,25,50}, 100 → 292 ways.
- **Errors**: ErrNegative, ErrNotPositive, ErrTooLarge; count also ErrDuplicate, ErrOverflow (partial, L28).
- **Tags**: coin-change. LeetCode 322, 518.

### ALG-DP-010 Fibonacci
- **Source**: Knuth TAOCP Vol. 1 3rd ed. §1.2.8 (definition; addition identity is eq. (6), confirmed);
  Erickson, Algorithms (2019) §3.1 "Mātrāvṛtta" (F(n+1) meters of total length n: the climbing-stairs count;
  confirmed). **Best known**: O(log n) by doubling, but n ≤ 93, so the loop is O(1) in practice (lost: doubling,
  no gain at n ≤ 93; closed form, float error). **Errors**: ErrNegative, ErrOverflow (n ≥ 94).
- **Worked example**: F(0…12) table; F(93) = 12200160415121876738. **Tags**: fibonacci, climbing-stairs. LeetCode 509, 70.

### ALG-DP-011 MaxNonAdjacentSum · ALG-DP-012 MaxNonAdjacentSumCircular
- **Source**: Roughgarden, Algorithms Illuminated Part 3 (2019) ch. 16 §16.1–16.2 (WIS in path graphs) (print);
  Kleinberg–Tardos §6.1 (same recurrence). **Circular, proof**: a cycle independent set cannot contain both
  values[0] and values[n−1], so it is independent in values[1:] or values[:n−1]; conversely each path set is
  a cycle set. So the optimum is the max of two path optima. (Composition of the cited recurrence, as L15.)
- **Best known**: O(n), O(1). **Worked example**: path 1 4 5 4 → 8. **Errors**: ErrOverflow (exact: an
  overflowing candidate is a feasible sum > MaxInt64). **Tags**: independent-set. LeetCode 198, 213.

### ALG-DP-013 GridPathCount · ALG-DP-014 MinPathSumGrid · ALG-DP-024 MinPathSumTriangle
- **Source**: CLRS Exercise 22.4-2 (number of paths in a DAG, linear time) (print); CLRS §24.2 (DAG shortest
  paths); Concrete Mathematics §5.1 (addition formula: the open grid gives C(r+c−2, r−1)). Krattenthaler,
  "Lattice Path Enumeration" (Handbook of Enumerative Combinatorics, 2015, ch. 10, arXiv:1503.05930) §10.7,
  Lemma 10.7.2 counts paths avoiding points by a determinant — the alternative that lost (O(k³) for k
  obstacles, needs big binomials; the DP is exact and O(rc)). The grid with right/down moves is a DAG whose
  row-major order is topological; the triangle likewise bottom-up.
- **Best known**: O(cells) time, O(row) space; reading the grid Ω(cells).
- **Worked examples**: 3×7 open → 28; [[2 9 1][4 1 8][7 3 2]] → 12; triangle 3/6 2/1 9 4/8 5 7 2 → 11.
- **Errors**: ErrNotRectangular, ErrOverflow (partial, L28); ErrEmpty (sums); ErrNotTriangle.
- **Tags**: grid, lattice-path, dag. LeetCode 62, 63, 64, 120.

### ALG-DP-015 TextSegmentation
- **Source**: Erickson, Algorithms (2019) §2.5 "Text Segmentation (Interpunctio Verborum)" and §3.3
  "Interpunctio Verborum Redux" (confirmed); Kleinberg–Tardos ch. 6 Exercise 5 (print).
- **Best known**: O(D + n·L) with a dictionary trie. **Approach**: trie in one hash map (edge → child), forward
  reachability with back-pointers. **Lost**: substring hashing per length (O(n·L²)); algo/trie's Trie (no
  exported single-step walk). **Worked example**: Erickson's BOTHEARTHANDSATURNSPIN splits.
- **Errors**: none; (nil, false) when no split. **Tags**: word-break. LeetCode 139, 140.

### ALG-DP-017 PalindromePartitionMinCuts
- **Source**: Fici, Gagie, Kärkkäinen, Kempa, J. Discrete Algorithms 28:41–48, 2014,
  doi:10.1016/j.jda.2014.08.001, §2 "A simple quadratic algorithm" (confirmed).
- **Best known**: O(n log n) time, O(n) space (same paper; also I et al. 2014). **Built**: O(n²) worst, O(n)
  space, expanding centres in increasing order so cuts[i] is final when read. **Why not the best**: the
  O(n log n) method needs a palindromic tree with series links; LEFT FOR DEVIN (go-complexity rule 3: needs a
  benchmark-backed reason or the faster build). **Worked examples**: racecarxyz → 3; noonabbad → 2 (hand,
  checked offline). **Tags**: palindrome, factorization. LeetCode 132.

### ALG-DP-018 RodCutting
- **Source**: CLRS §15.1 eq. (15.2), BOTTOM-UP-CUT-ROD, EXTENDED-BOTTOM-UP-CUT-ROD, Fig. 15.1 (print).
- **Best known**: O(n·min(n, m)). **Worked example**: r1…r10 = 1 5 8 10 13 17 18 22 25 30; 4 → 2+2; 7 → 1+6.
- Lengths beyond the price table are waste (earn 0). **Errors**: ErrNegative, ErrTooLarge, ErrOverflow.

### ALG-DP-019 MatrixChainOrder
- **Source**: CLRS §15.2 eq. (15.7), MATRIX-CHAIN-ORDER, PRINT-OPTIMAL-PARENS, Fig. 15.5 (print).
- **Best known**: O(n log n), Hu–Shing, SIAM J. Comput. 11(2):362–373, 1982, doi:10.1137/0211028 (confirmed) and
  Part II 13(2):228–251, 1984. **Built**: O(n³)/O(n²) CLRS. **Why not the best**: Hu–Shing is a long,
  intricate polygon-partition algorithm; LEFT FOR DEVIN. Overflowing splits are +∞ (costs ≥ 0), so only a true
  overflow of the answer errors. **Worked example**: 30 35 15 5 10 20 25 → 15125, ((A1(A2A3))((A4A5)A6)).
- **Errors**: ErrEmpty, ErrNotPositive, ErrTooLarge (n ≳ 3340), ErrOverflow. LeetCode none (pattern: interval DP).

### ALG-DP-020 CanJumpToEnd · ALG-DP-021 MinJumpsToEnd
- **Source**: CLRS §22.2 (BFS gives fewest-edge distances). **Proof (for the one-pass form)**: edges go
  i → i+1 … i+a_i. Claim: each BFS layer L_d is a contiguous run (e_{d−1}, e_d]. L_0 = {0}. If every index
  ≤ e_d is reached by layer d, the indices reachable in one more move are exactly (e_d, max_{i ≤ e_d}(i+a_i)],
  a run, since every j ≤ the max is covered by the edge from the index achieving it. So layer boundaries are
  the running maxima, and reachability is "farthest ≥ n−1". (Composition of a cited method, as L15; flagged.)
- **Best known**: O(n), O(1). **Lost**: explicit BFS (O(n + Σa_i)). **Worked example**: [2 0 1 0] → 2 moves.
- **Errors**: ErrNegative. **Tags**: bfs, greedy. LeetCode 55, 45.

### ALG-DP-022 MaxProfitOneTrade · ALG-DP-023 MaxProfitUnlimitedTrades
- **Source**: Bentley, Programming Pearls 2nd ed., Column 8 §8.4 "A Scanning Algorithm" (print): p[s] − p[b]
  is the sum of the changes b+1 … s, so one trade is the maximum subvector of the changes; the scan state is
  the lowest price so far. **Unlimited, proof**: a trade's profit is a sum of changes ≤ the sum of its positive
  changes; trades are disjoint, so any set earns ≤ Σ positive changes; trading every rise attains it.
  (Flagged as a short proof on a cited identity.)
- **Best known**: O(n), O(1). **Worked example**: [9 4 6 2 8 5] → one trade 6 (buy 3, sell 4); unlimited 8.
- **Errors**: ErrNegative (prices ≥ 0, so differences fit), ErrOverflow (unlimited). LeetCode 121, 122.

### ALG-DP-025 SubsequenceOccurrences
- **Source**: Lothaire, Combinatorics on Words, Addison-Wesley 1983 (CUP 1997), ch. 6 "Subwords" (Sakarovitch,
  Simon) §6.3: (ua choose vb) = (u choose vb) + δ_{a,b}(u choose v), (u choose ε) = 1 (section confirmed via
  Rigo et al., orbi.uliege.be; the book page (print)). Elzinga–Rahmann–Wang 2008 was checked and does NOT
  cover this count; not cited. **Best known**: O(nm) time, O(m) space. **Worked example**: "see" in
  "subsequence" → 6. **Errors**: ErrOverflow (partial, L28). LeetCode 115.

## HELD (no acceptable source)
- DecodeWays (digit-string decodings 1–26): no textbook or paper found.
- MaxProfit with at most k trades: Csűrös, IEEE/ACM TCBB 1(4):139–150, 2004, doi:10.1109/TCBB.2004.43 gives
  O(n log n) for k maximum-scoring segments, but whether it states the O(nk) state DP is unverified, and its
  O(n log n) method was not read in full. HELD until the paper is read.
- MaxProfit with cooldown (state-machine DP): no source found.
- MaximalSquare (largest all-true square): no citable source for the min-of-three recurrence found
  (Rosenfeld–Pfaltz 1966 distance transforms are related but not this recurrence without a proof of ours).

## Consumer evidence
None given in the brief.

---

# Catalog: algo/bits (SHORT = BITS), new package

Builder W3-B10, 2026-11-12. Package name `bits` is the one approved in doc/taxonomy.md; the standard
library's math/bits is imported beside it under an alias (same reasoning as L24 for algo/list). Inputs are
never modified; nothing allocates except PopCountTable when dst lacks capacity.

Sources used throughout (confirmed by table of contents / search, 2026-11-12):
- Knuth, *The Art of Computer Programming*, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and
  techniques" (rightmost bits, x & −x, sideways addition ν, subset enumeration x ← (x − χ) & χ,
  bit interleaving); §7.2.1.1 "Generating all n-tuples" (Gray binary code g(k) = k ⊕ ⌊k/2⌋ and its inverse).
- Warren, *Hacker's Delight*, 2nd ed., Addison-Wesley, 2013: §2-1 "Manipulating Rightmost Bits", §3-2
  "Rounding Up/Down to the Next Power of 2" (clp2, flp2), ch. 5 "Counting Bits", §7-1 "Reversing Bits and
  Bytes", §7-2 "Shuffling Bits", ch. 13 "Gray Code".
- Gray, F., "Pulse code communication", U.S. Patent 2,632,058 (filed 1947, issued 17 Mar 1953).
- Hamming, R. W., "Error detecting and error correcting codes", Bell System Technical Journal 29(2):147–160,
  1950, doi:10.1002/j.1538-7305.1950.tb00463.x.
- Morton, G. M., "A computer oriented geodetic data base and a new technique in file sequencing", IBM Ltd.,
  Ottawa, 1966 (WorldCat OCLC 52181262; confirmed 2026-11-12).
- Laaksonen, *Competitive Programmer's Handbook* (L18), ch. 10 "Bit manipulation" (submask iteration).

## Dropped (exactly one math/bits call; logged)
- **HammingWeight / popcount** = `math/bits.OnesCount64`.
- **LowestSetBit index** = `math/bits.TrailingZeros64`; **HighestSetBit index** = `math/bits.Len64(x) − 1`.
  The VALUE forms (x & −x; the greatest power of two ≤ x) are built (BITS-007, -008): they are not one call.
- **Full 64-bit ReverseBits** = `math/bits.Reverse64`; the field form with a width is built (BITS-004).

### ALG-BITS-001 IsPowerOfTwo
`func IsPowerOfTwo[T integer](x T) bool`
- **Source**: Warren §2-1 (x & (x−1) turns off the rightmost 1); Knuth §7.1.3. **Best known**: O(1).
- **Approach**: x > 0 && x&(x−1) == 0. **Lost**: OnesCount64(x) == 1 (needs the unsigned conversion and a
  sign test anyway; same O(1)). Signed minimum is negative → false.
- **Worked example**: 2^k for k = 0…63 → true; 0, 6, −4 → false. **Errors**: none. **Tags**: bits, power-of-two. LeetCode 231.

### ALG-BITS-002 IsPowerOfFour
`func IsPowerOfFour[T integer](x T) bool`
- **Source**: as 001; 4^k = 2^(2k) has its one bit at an even index (mask 0x5555…). **Best known**: O(1).
- **Lost**: loop dividing by 4 (O(log x)); x % 3 == 1 test (correct, but a division instead of an AND).
- **Worked example**: 1, 4, 16, 2^62 → true; 2, 8, 2^63 → false. **Errors**: none. **Tags**: bits, power-of-four. LeetCode 342.

### ALG-BITS-003 NextPowerOfTwo
`func NextPowerOfTwo(x uint64) (uint64, error)`
- **Source**: Warren §3-2, clp2(x) = 2^(64 − nlz(x − 1)). **Best known**: O(1).
- **Lost**: the smear-right sequence x |= x>>1 … >>32 (six steps; LeadingZeros64 is one instruction).
- **Worked example**: 0 → 1, 1 → 1, 5 → 8, 1000 → 1024, 2^63 → 2^63.
- **Errors**: ErrOverflow for x > 2^63. **Tags**: bits, power-of-two, round-up.

### ALG-BITS-004 ReverseBits
`func ReverseBits(x uint64, width int) (uint64, error)`
- **Source**: Warren §7-1 (reverse the word, shift right by w − width); Knuth §7.1.3. **Best known**: O(1).
- **Lost**: bit-by-bit loop (O(width)); a 256-entry byte table (a global table and 8 lookups; math/bits is
  one instruction on arm64, a few on amd64).
- **Worked example**: (0b0011, 4) → 0b1100; (1, 64) → 2^63.
- **Errors**: ErrWidth (width ∉ 1…64); ErrOutOfRange (x has bits at or above width). **Tags**: bits, reverse. LeetCode 190.

### ALG-BITS-005 GrayEncode · ALG-BITS-006 GrayDecode
`func GrayEncode(x uint64) uint64` · `func GrayDecode(g uint64) uint64`
- **Source**: Gray 1953 patent; Knuth §7.2.1.1 (g(k) = k ⊕ ⌊k/2⌋, inverse k = g ⊕ ⌊g/2⌋ ⊕ ⌊g/4⌋ ⊕ …); Warren ch. 13.
- **Best known**: O(1) encode; decode O(log w) by prefix XOR (6 steps for w = 64). **Lost** for decode:
  the bit-serial loop (O(w)).
- **Worked example** (Knuth §7.2.1.1): 0…7 → 000 001 011 010 110 111 101 100, and back.
- **Errors**: none. **Tags**: bits, gray-code. LeetCode 89 (sequence = GrayEncode(0…2^n−1)).

### ALG-BITS-007 LowestSetBit · ALG-BITS-008 HighestSetBit
`func LowestSetBit(x uint64) uint64` · `func HighestSetBit(x uint64) uint64`
- **Source**: Knuth §7.1.3 "rightmost bits" (x & −x); Warren §2-1; Warren §3-2 flp2 (greatest power of 2 ≤ x).
- **Best known**: O(1). **Lost** for 008: smear-right then x − (x >> 1) (six steps vs one nlz).
- **Worked example**: 0b10110 → 0b00010 and 0b10000; 0 → 0 for both. **Errors**: none.
- **Tags**: bits, lowest-set-bit, highest-set-bit, fenwick.

### ALG-BITS-009 PopCountTable
`func PopCountTable(dst []uint8, n int) ([]uint8, error)`
- **Source**: Knuth §7.1.3 sideways addition with ν(x) = ν(⌊x/2⌋) + (x mod 2); Warren ch. 5.
- **Best known**: O(n) (the output has n+1 entries: Ω(n)). **Approach**: table recurrence, one add per entry.
  **Lost**: OnesCount per entry (also O(n) with POPCNT, but not on every target; the recurrence is
  branch-free and uses only the table).
- **Worked example**: n = 8 → 0 1 1 2 1 2 2 3 1. **Errors**: ErrNegative (n < 0); ErrTooLarge (n+1
  overflows int; beyond that, a table too large for memory is the runtime's, per L26).
- **Input modified**: no; appends to dst. **Tags**: bits, popcount, counting-bits, dp. LeetCode 338.

### ALG-BITS-010 LoneElementAmongPairs
`func LoneElementAmongPairs[T integer](values []T) T`
- **Source**: Knuth §7.1.3 (⊕ associative, commutative, x ⊕ x = 0). **Best known**: O(n) time, O(1) space
  (every value must be read). **Lost**: hash-map count (O(n) space); sort (O(n log n)).
- **Worked example**: [3 8 5 3 5] → 8. **Errors**: none; empty → 0; precondition unchecked (XOR of all).
- **Tags**: bits, xor, single-number. LeetCode 136.

### ALG-BITS-011 LoneElementAmongTriples
`func LoneElementAmongTriples[T integer](values []T) T`
- **Source**: Knuth §7.1.3 (bitwise operations act on all positions in parallel) plus the proof below; no
  textbook names this exact problem. FLAGGED for Devin as a cited-step composition (like L15).
- **Proof**: fix a bit position. If every value but one occurs 3 times, the number of values with that bit
  set is 3m + b, where b is the lone value's bit. Keep a two-bit counter (twos, ones) of that number modulo
  3. The updates ones ← (ones ⊕ v) ∧ ¬twos, then twos ← (twos ⊕ v) ∧ ¬ones map, for a set bit of v:
  (0,0) → (0,1) → (1,0) → (0,0), and leave the state unchanged for a clear bit (checked case by case: with
  v = 0, ones ← ones ∧ ¬twos = ones because the state (1,1) never occurs, and likewise twos). So after all
  values the counter is (0, b) and ones = b at every position, i.e. ones is the lone value. Two's
  complement makes this hold for signed types too.
- **Best known**: O(n) time, O(1) space. **Lost**: per-bit column count modulo 3 (O(w·n), w = 64);
  hash-map count (O(n) space).
- **Worked example**: [4 9 4 1 9 4 9] → 1. **Errors**: none; empty → 0; precondition unchecked.
- **Tags**: bits, single-number, modular-counting. LeetCode 137.

### ALG-BITS-012 BitwiseAndRange
`func BitwiseAndRange(lo, hi uint64) (uint64, error)`
- **Source**: Knuth §7.1.3 (AND; the leftmost bit of x ⊕ y is the first differing position) plus the proof
  below. FLAGGED for Devin as a cited-step composition (like L15).
- **Proof**: if lo = hi the AND is lo. Otherwise let d be the highest bit where lo and hi differ; then lo
  has 0 there and hi has 1, and they agree above d. Every v in [lo, hi] agrees with them above d (they are
  the bounds of an interval inside one aligned block of size 2^(d+1)). The value m = (hi with bits < d
  cleared) satisfies lo < m ≤ hi, and m − 1 ≥ lo; m has bits 0…d−1 clear and m − 1 has bit d clear, so
  every bit ≤ d is 0 in the AND. Result: lo with bits 0…d cleared = lo >> (d+1) << (d+1).
- **Best known**: O(1). **Lost**: shift both right until equal (O(w)); Kernighan clearing hi & (hi−1)
  while hi > lo (O(w)).
- **Worked example**: [5, 7] → 4; [12, 15] → 12; [8, 16] → 0. **Errors**: ErrRange (lo > hi).
- **Tags**: bits, and, range. LeetCode 201.

### ALG-BITS-013 HammingDistance
`func HammingDistance(x, y uint64) int`
- **Source**: Hamming 1950 (distance between code points); Knuth §7.1.3 ν(x ⊕ y). **Best known**: O(1).
- **Worked example**: 0b1011101 vs 0b1001001 → 2. **Errors**: none. **Tags**: bits, hamming-distance. LeetCode 461, 477.
- Hamming weight dropped (= math/bits.OnesCount64).

### ALG-BITS-014 Submasks
`func Submasks(mask uint64) iter.Seq[uint64]`
- **Source**: Knuth §7.1.3 (running through subsets with x ← (x − χ) & χ); Laaksonen ch. 10 (subset iteration).
- **Best known**: O(2^k) for k set bits (the output size), O(1) per submask. **Lost**: iterate 0…mask and
  test s & ^mask (O(mask), exponential in the mask's top bit, not its popcount).
- **Worked example**: 0b1101 → 1101 1100 1001 1000 0101 0100 0001 0000. **Errors**: none; early break supported.
- **Tags**: bits, submask, subset-enumeration, iterator, bitmask-dp.

### ALG-BITS-015 MortonEncode2D · ALG-BITS-016 MortonDecode2D
`func MortonEncode2D(x, y uint32) uint64` · `func MortonDecode2D(code uint64) (x, y uint32)`
- **Source**: Morton 1966 (the interleaved file sequence); Warren §7-2 "Shuffling Bits" (outer perfect shuffle
  by masks and shifts in log2 w steps); Knuth §7.1.3 (bit interleaving / zipper). Knuth names it in 7.1.3;
  equation numbers not confirmed (section-level citation).
- **Best known**: O(log w) with mask-and-shift (5 steps per coordinate), O(1) with BMI2 PDEP/PEXT (not
  reachable from portable Go). **Lost**: bit-by-bit loop (O(w)); 256-entry lookup tables (globals).
- **Worked example**: (0,0) (1,0) (0,1) (1,1) → 0 1 2 3; (2,3) → 14; round trip for every input.
- **Errors**: none. **Tags**: bits, morton, z-order, spatial-index.

---

# Catalog: algo/intmath (SHORT = INTMTH), new package

Builder W3-B10, 2026-11-12. Integer algorithms on machine words and digit strings. Inputs never modified.
Every result that can exceed its type returns ErrOverflow. Signed results follow Go's truncated division.

## Not built here (one function, one home)
- GCD/LCM, ExtendedGCD, ModPow, ModInverse, SieveOfEratosthenes, IsPrimeDeterministic64 → num/integer (L4, D-005;
  IsPrime is NUM-INTEG-001).
- **SqrtNewton DROPPED**: math.Sqrt is correctly rounded (IEEE 754-2008 §5.4.1 requires squareRoot to be
  correctly rounded; Go's math.Sqrt documents IEEE behaviour and compiles to the hardware instruction). A
  Newton iteration on float64 can at best equal it and is slower. Integer ⌊√n⌋ is built (INTMTH-002).
- **NumberOfDigitOne HELD (no source)**: the per-position counting formula for "how many digit 1s appear in
  1…n" was found only on problem sites and blogs (search 2026-11-12); no textbook or paper located.

Sources used (confirmed 2026-11-12 by search / table of contents unless marked):
- Knuth, *TAOCP* Vol. 2, 3rd ed., 1997: §3.1 ex. 6 (Floyd's cycle finding, confirmed); §4.3.1 Algorithms A, D, M;
  §4.4 radix conversion Methods (1a), (1b), (2a); §4.6.3 Algorithm A (right-to-left binary exponentiation).
- Cohen, *A Course in Computational Algebraic Number Theory*, Springer GTM 138, 1993, Algorithm 1.7.1 (integer
  square root; the book is confirmed, the algorithm number is from the standard reference list, not seen).
- Crandall, Pomerance, *Prime Numbers: A Computational Perspective*, 2nd ed., 2005, Algorithm 9.2.11 (not seen; second source).
- Dijkstra, *A Discipline of Programming*, Prentice-Hall, 1976, ch. 17 "An exercise attributed to R. W. Hamming" (confirmed).
- Graham, Knuth, Patashnik, *Concrete Mathematics*, 2nd ed., 1994, §4.4 "Factorial Factors" (ε_p(n!); the
  equation number (4.24) is from memory of the text, section confirmed).
- Gauss, *Disquisitiones Arithmeticae*, 1801, art. 12 (digit sum ≡ number mod 9, confirmed) and art. 291–292
  (three-square theorem, confirmed).
- Hardy, Wright, *An Introduction to the Theory of Numbers*, 6th ed., OUP, 2008: Theorem 2 (fundamental theorem
  of arithmetic); §9.2 "Terminating and recurring decimals" (confirmed); §20.5 Theorem 369 (four squares) and
  §20.10 (g(2) = 4 via 8m + 7) (confirmed in the table of contents).
- Guy, *Unsolved Problems in Number Theory*, 3rd ed., Springer, 2004, E34 "Happy numbers".
- Cormen et al., *Introduction to Algorithms*, 3rd ed., 2009, §31.6 (repeated squaring); Warren, *Hacker's
  Delight*, 2nd ed., 2013, ch. 9 "Integer Division".

### ALG-INTMTH-001 PowerBySquaring
`func PowerBySquaring(base, exp int64) (int64, error)`
- **Source**: Knuth Vol. 2 §4.6.3 Algorithm A; CLRS §31.6. **Best known**: O(log e) multiplications
  (addition chains can save a constant, not the order; Knuth §4.6.3).
- **Approach**: right-to-left binary on |base| in uint64 with bits.Mul64; the limit is 2^63 − 1, or 2^63
  when the result is negative. The square of the running base is checked only when a set bit of e remains,
  because it is then a factor of the result. **Lost**: left-to-right (same count; needs the top bit first);
  float math.Pow (inexact above 2^53).
- **Worked example**: 2^23 = 8388608 (Knuth's x^23). (−2)^63 = MinInt64; 0^0 = 1.
- **Errors**: ErrNegative (exp < 0); ErrOverflow. **Tags**: integer, power, binary-exponentiation. LeetCode 50 (integer form).

### ALG-INTMTH-002 ISqrt
`func ISqrt(n uint64) uint64`
- **Source**: Cohen Algorithm 1.7.1; Crandall–Pomerance Algorithm 9.2.11.
- **Termination proof**: let s = ⌊√n⌋. (i) Never below s: since x is an integer, ⌊(x + ⌊n/x⌋)/2⌋ =
  ⌊(x + n/x)/2⌋ (nested floors with an integer term, Concrete Mathematics §3.2), and (x + n/x)/2 ≥ √n by AM–GM, so y ≥ s.
  (ii) Strict decrease above s: if x ≥ s + 1 > √n then n/x < x, so (x + n/x)/2 < x and y < x. The iterates
  are integers that fall strictly while > s and never drop below s, so they reach s; there (i) gives y ≥ s = x
  and the loop returns s.
  Start x₀ = 2^⌈b/2⌉ > √n (b = bit length). Overflow: x ≤ 2^32 and ⌊n/x⌋ ≤ 2^32 + 2.
- **Best known**: O(1) word operations for 64 bits; Newton from within a factor 2 converges quadratically
  (≤ 7 steps). **Lost**: float math.Sqrt then correct by ±1 (correct but depends on float rounding near
  2^64; the integer proof is simpler to audit); bit-by-bit digit method (32 steps).
- **Worked example**: 99 → 9, 100 → 10, 2^64 − 1 → 2^32 − 1. **Errors**: none. **Tags**: integer, isqrt. LeetCode 69, 367.

### ALG-INTMTH-003 LongDivision
`func LongDivision[T signedInteger](dividend, divisor T) (quotient, remainder T, err error)`
- **Source**: Knuth §4.3.1 Algorithm D in radix 2 (each quotient digit is one comparison); Warren ch. 9.
- **Best known**: O(w) bit steps without the hardware divide (w = bit length of |dividend|).
- **Approach**: restoring shift-and-subtract on magnitudes in uint64, then Go's signs (quotient toward zero,
  remainder with the dividend's sign). **Lost**: doubling the divisor and subtracting from the top (same
  O(w) but needs overflow care on the doubled divisor); repeated subtraction (O(quotient)).
- **Worked example**: Go spec: 5/3 = 1 r 2, −5/3 = −1 r −2, 5/−3 = −1 r 2, −5/−3 = 1 r −2.
- **Errors**: ErrDivideByZero; ErrOverflow for MinT / −1 (documented rule: refused, not wrapped).
- **Tags**: integer, division, shift-subtract. LeetCode 29.

### ALG-INTMTH-004 ReverseDigits
`func ReverseDigits[T signedInteger](x T) (T, error)`
- **Source**: Knuth §4.4 Methods (1a) and (1b). **Best known**: O(d), d ≤ 19 digits.
- **Approach**: reverse |x| in uint64 (cannot overflow: < 10^19 < 2^64), then range-check against T.
  **Lost**: reverse in T with a pre-multiply bound check per digit (needs T's max, which a generic body
  cannot name without reflection).
- **Worked example**: 1234 → 4321, −560 → −65. **Errors**: ErrOverflow. **Tags**: integer, digits, reverse. LeetCode 7.

### ALG-INTMTH-005 IsPalindromeNumber
`func IsPalindromeNumber[T integer](x T) bool`
- **Source**: Knuth §4.4 Methods (1a), (1b). **Best known**: O(d). **Approach**: reverse only the low half
  and compare (no overflow, d/2 divisions). **Lost**: format to a string (allocates); full reversal (can overflow).
- **Worked example**: 121, 1221, 0 → true; −121, 10 → false. **Errors**: none. **Tags**: integer, palindrome. LeetCode 9.

### ALG-INTMTH-006 IsHappyNumber
`func IsHappyNumber(n uint64) bool`
- **Source**: Knuth Vol. 2 §3.1 ex. 6 (Floyd); Guy E34 (definition, the unhappy cycle 4 → … → 4).
- **Best known**: O(log n) (read the digits once; after one step every value ≤ 1620). **Lost**: a visited set
  (O(cycle) memory, a map); hard-coding the cycle {4, …} (correct, but rests on a fact the code would not check).
- **Worked example** (Guy E34): happy numbers ≤ 50 are 1 7 10 13 19 23 28 31 32 44 49. **Errors**: none; 0 → false.
- **Tags**: integer, happy-number, cycle-detection. LeetCode 202.

### ALG-INTMTH-007 DigitalRoot
`func DigitalRoot(n uint64) uint64`
- **Source**: Gauss art. 12 (n ≡ digit sum mod 9). Since the root is a single digit 1…9 for n ≥ 1 and is
  ≡ n mod 9, it is 1 + (n − 1) mod 9. **Best known**: O(1). **Lost**: repeated digit sums (O(log n)).
- **Worked example**: 38 → 2. **Errors**: none; 0 → 0. **Tags**: integer, digital-root. LeetCode 258.

### ALG-INTMTH-008 FactorialTrailingZeros
`func FactorialTrailingZeros(n uint64) uint64`
- **Source**: Concrete Mathematics §4.4 (Legendre: ε_p(n!) = Σ ⌊n/p^k⌋); zeros = min(ε_2, ε_5) = ε_5.
- **Best known**: O(log n). **Lost**: computing n! (impossible beyond 20!).
- **Worked example**: 100! ends in 24 zeros. **Errors**: none. **Tags**: integer, factorial, legendre-formula. LeetCode 172.

### ALG-INTMTH-009 IsHammingNumber · ALG-INTMTH-010 NthHammingNumber
`func IsHammingNumber(n uint64) bool` · `func NthHammingNumber(n int) (uint64, error)`
- **Source**: Dijkstra 1976 ch. 17 (set definition; generation in order by three indexes into the sequence).
- **Best known**: membership O(log n). n-th: O(n) time and space by Dijkstra's merge; O(n^(2/3))-ish counting
  methods exist but need floating logarithms and lose exactness; n is at most 13282 for a uint64 result.
- **Lost** for 010: a min-heap of candidates (O(n log n)); enumerate triples and sort (O(n log n)).
- **Overflow**: 13282 Hamming numbers are ≤ 2^64 − 1 (counted exactly with math/big; the test recounts);
  the largest is 18_432_000_000_000_000_000 = 2^26·3^2·5^15. n > 13282 → ErrOverflow before allocating.
- **Worked example** (Dijkstra): 1 2 3 4 5 6 8 9 10 12 15 16 18 20 24 25 27 30.
- **Errors**: 010: ErrOutOfRange (n < 1), ErrOverflow. **Tags**: integer, hamming-number, 5-smooth, ugly-number. LeetCode 263, 264.

### ALG-INTMTH-011 MinSquaresToSum
`func MinSquaresToSum(n uint64) int`
- **Source**: Hardy–Wright §20.5 Theorem 369 (four squares); Gauss art. 291–292 (n is a sum of three
  squares iff n ≠ 4^a(8b+7)); H–W §20.10 (necessity). Two squares: 4m = x² + y² forces x, y even (squares are
  0 or 1 mod 4), so stripping factors 4 preserves the answer; m ≡ 3 (mod 4) is never a sum of two squares.
- **Best known**: O(√n) with this method (a factorisation-based two-square test, Fermat/H–W Theorem 366,
  would be faster only with a fast factoriser; num/integer has Pollard rho but algo may not import num).
  **Lost**: the DP over all m ≤ n (O(n√n) time, O(n) space) and BFS (same).
- **Worked example**: 7, 15, 28 → 4; 12 → 3; 13 → 2; 25 → 1. **Errors**: none; 0 → 0.
- **Tags**: integer, sum-of-squares, lagrange, legendre. LeetCode 279.

### ALG-INTMTH-012 IsPowerOfThree
`func IsPowerOfThree[T integer](x T) bool`
- **Source**: Hardy–Wright Theorem 2 (unique factorisation ⇒ divisors of 3^40 are exactly 3^j). 3^40 =
  12157665459056928801 is the largest power of 3 ≤ 2^64 − 1 (checked with math/big). **Best known**: O(1).
  **Lost**: divide by 3 in a loop (O(log x)); float log (inexact).
- **Worked example**: 27 → true, 45 → false. **Errors**: none. **Tags**: integer, power-of-three. LeetCode 326.

### ALG-INTMTH-013 FractionToDecimal
`func FractionToDecimal(numerator, denominator int64, maxFractionDigits int) (string, error)`
- **Source**: Hardy–Wright §9.2 (a rational has a terminating or recurring decimal); Knuth §4.4 Method (2a).
- **Approach**: long division; remainders r ← 10r mod d in 128-bit (bits.Mul64/Div64, so |d| up to 2^63 is
  exact); a map from remainder to the position of its digit finds the first repeat; "(" is inserted there.
  **maxFractionDigits** is a caller cap because the period can be |d| − 1 digits (the output alone could be
  9·10^18 bytes); ErrTooLong past it. **Best known**: O(L) for L output digits. **Lost**: compute the
  pre-period max(v₂, v₅) and the period ord_q(10) first (needs gcd and order, i.e. num/integer, and still
  emits L digits); Floyd/Brent cycle detection (O(1) memory but ~3× the divisions, and the result is O(L) anyway).
- **Worked example** (H–W §9.2): 1/3 = 0.(3), 3/7 = 0.(428571); 1/6 = 0.1(6); −1/6 = "-0.1(6)"; 22/7 = 3.(142857).
- **Errors**: ErrDivideByZero; ErrNegative (cap < 0); ErrTooLong. **Tags**: fraction, repeating-decimal. LeetCode 166.

### ALG-INTMTH-014 AddBinaryStrings
`func AddBinaryStrings(a, b string) (string, error)`
- **Source**: Knuth §4.3.1 Algorithm A, radix 2. **Best known**: O(m + n). **Lost**: math/big parse and add
  (three allocations and conversions; this is one pass).
- **Worked example**: "1011" + "11" = "1110". Leading zeros accepted, none in the output.
- **Errors**: ErrEmpty; ErrInvalidDigit. **Tags**: binary, addition, digit-string. LeetCode 67.

### ALG-INTMTH-015 MultiplyDecimalStrings
`func MultiplyDecimalStrings(a, b string) (string, error)`
- **Source**: Knuth §4.3.1 Algorithm M, radix 10. **Best known** for schoolbook: O(m·n); asymptotically
  faster methods exist (Karatsuba O(n^1.585), Knuth §4.3.3; Schönhage–Strassen / Harvey–van der Hoeven
  O(n log n)) and are what math/big uses for large operands. Schoolbook kept: the request named it, it is
  exact and simple, and a caller with long operands should use math/big. **Lost**: Karatsuba (faster above
  a few hundred digits; not requested, and math/big already ships it).
- **Worked example**: "123" × "45" = "5535". **Errors**: ErrEmpty; ErrInvalidDigit. **Tags**: multiplication, digit-string. LeetCode 43.

### ALG-INTMTH-016 IncrementDigits
`func IncrementDigits(dst, digits []uint8) ([]uint8, error)`
- **Source**: Knuth §4.3.1 Algorithm A with v = 1, radix 10. **Best known**: O(n) worst case (all nines);
  O(1) amortised over a run of increments. **Approach**: find the last non-9, copy, bump, zero the tail.
- **Worked example**: [1 2 9] → [1 3 0]; [9 9] → [1 0 0]; [0 9] → [1 0] (leading zeros kept).
- **Errors**: ErrEmpty; ErrInvalidDigit (> 9). **Input modified**: no; appends to dst. **Tags**: digits, plus-one. LeetCode 66.

---

# Catalog: algo/text (SHORT = TEXT), new package

Builder W3-B11, 2026-11-12. Oracles are capped (≤ 80 bytes, ≤ 2000 trials; the exponential wildcard
oracle ≤ 9 runes) per the go-testing.md 60 s rule. New package, numbers from 001. `algo/text` is in doc/taxonomy.md
("string matching and distance").

**Conventions (stated once in doc.go).** Byte semantics by default: the matching, hashing and
palindrome functions are generic over `S ~string | ~[]byte`, index bytes, and report byte offsets
(one implementation serves both string and []byte callers without conversion). Rune semantics only
where the answer changes meaning on UTF-8 text: IsAnagram, WildcardMatch ('?' is one rune) and
RegexMatch ('.' is one rune); invalid UTF-8 bytes read as U+FFFD there (the Go conversion rule, as
ALG-HASH-010). ReverseWords and StringToInt use ASCII white space, which is UTF-8 safe (no multi-byte
rune contains an ASCII byte). Searches allocate only their pattern tables (KMP: the prefix function,
m ints); Rabin–Karp allocates nothing beyond dst. Large outputs append to dst (D-003 style).
Lost alternative: separate string and []byte functions (doubles the API) / rune-level everywhere
(the classical algorithms are defined on symbol arrays; bytes are the symbols of a Go string, and
rune decoding would make every index ambiguous).

**Confirmations (2026-11-12)**: Crossref confirmed KMP 1977 (10.1137/0206024), Karp–Rabin 1987
(10.1147/rd.312.0249), Manacher 1975 (10.1145/321892.321896), Apostolico–Breslauer–Galil 1995
(10.1016/0304-3975(94)00083-U), Foster 1947 (10.2307/3029479). Read first hand: SemVer 2.0.0 (§2, §9–§11,
BNF), RFC 8259 §6 ABNF, ISO/IEC 9899:2011 draft N1570 §7.22.1.4 (strtol), TIFF 6.0 Section 9 (p. 42).
Thompson 1968 (10.1145/363347.363387): Crossref returned 429; citation from memory, reviewer to check.
CLRS, Gusfield, Bentley, Knuth, Aho–Lam–Sethi–Ullman, Laaksonen: textbook sections cited from memory.

### ALG-TEXT-001 PrefixFunction
`func PrefixFunction[S ~string | ~[]byte](dst []int, s S) []int`
- **Source**: Knuth, Morris, Pratt, "Fast Pattern Matching in Strings", SIAM J. Comput. 6(2):323–350, 1977,
  doi:10.1137/0206024; CLRS 3rd ed. §32.4, COMPUTE-PREFIX-FUNCTION.
- **Best known**: O(m) time (amortised: k rises at most m times), Ω(m) to read s. O(1) extra.
- **Approach**: CLRS's loop, π[q] = length of the longest proper prefix of s[:q+1] that is also its
  suffix (0-based). **Lost**: KMP's optimised `next` table (same bound; the plain π is what callers reuse
  for periods and borders).
- **Worked example**: CLRS Fig. 32.10: ababaca → 0 0 1 2 3 0 1. **Errors**: none. **Tags**: LeetCode 1392, 459.

### ALG-TEXT-002 KMPSearch
`func KMPSearch[S ~string | ~[]byte](dst []int, text, pattern S) ([]int, error)`
- **Source**: as 001 (KMP 1977; CLRS §32.4 KMP-MATCHER).
- **Best known**: O(n + m) worst case, O(m) extra (the π table); Ω(n + m).
- **Approach**: KMP matcher appending every start (overlapping). **Lost**: Boyer–Moore (sublinear on
  average, but O(nm) for all occurrences in its basic form); naive O(nm).
- **Worked example**: (hand) text abababacaba, pattern ababaca (CLRS Fig. 32.10's pattern) → [2].
- **Errors**: ErrEmptyPattern. **Tags**: LeetCode 28 (StrStr: DROPPED — first match only is exactly
  strings.Index / bytes.Index, which the standard library already does in O(n + m)).

### ALG-TEXT-003 ZFunction
`func ZFunction[S ~string | ~[]byte](dst []int, s S) []int`
- **Source**: Gusfield, *Algorithms on Strings, Trees, and Sequences*, Cambridge UP, 1997, §1.3 (Z values)
  and §1.4 (the Z algorithm, linear time).
- **Best known**: O(n), O(1) extra. **Approach**: the Z-box [l, r) algorithm. **Lost**: naive O(n²).
- **Convention**: z[0] = n (Gusfield leaves Z₁ undefined). **Worked example**: Gusfield §1.3, aabcaabxaaz:
  Z₅ = 3, Z₆ = 1, Z₇ = Z₈ = 0, Z₉ = 2 (1-based) → z = 11 1 0 0 3 1 0 0 2 1 0. **Errors**: none.

### ALG-TEXT-004 RollingHash (NewRollingHash; methods Append, Roll, Sum, Len, Reset — one ID, L23)
`func NewRollingHash(base, modulus uint64, windowLen int) (RollingHash, error)`
- **Source**: Karp, Rabin, "Efficient randomized pattern-matching algorithms", IBM J. Res. Dev.
  31(2):249–260, 1987, doi:10.1147/rd.312.0249; CLRS §32.2 (Horner's rule and the rolling update mod q).
- **Cost**: O(1) per Append/Roll (one 128-bit product and remainder, math/bits). **Lost**: Mersenne-only
  2⁶¹−1 reduction (faster, but cannot reproduce CLRS's q = 13 example); a fixed modulus 2⁶⁴ (no prime,
  easy to break with Thue–Morse strings).
- **Worked example**: CLRS Fig. 32.5 (base 10, q = 13, window 5): windows 31415 and 67399 hash alike.
- **Errors**: ErrHashParams (modulus < 2, base < 1 or ≥ modulus), ErrWindowLen (< 1).

### ALG-TEXT-005 RabinKarpSearch
`func RabinKarpSearch[S ~string | ~[]byte](dst []int, text, pattern S, base, modulus uint64) ([]int, error)`
- **Source**: as 004. Every fingerprint hit is verified byte by byte, so there are no false positives.
- **Best known**: O(n + m) expected for a random base; worst O(nm) (every window a hit). O(1) extra.
- **Lost**: KMP (002) — deterministic, but Rabin–Karp extends to many patterns and 2-D. Both kept.
- **Worked example**: CLRS Fig. 32.5: 2359023141526739921, 31415, q = 13 → [6]; spurious hit at 12 rejected.
- **Errors**: ErrEmptyPattern, ErrHashParams.

### ALG-TEXT-006 RepeatedSubstrings
`func RepeatedSubstrings[S ~string | ~[]byte](dst []int, text S, length int) ([]int, error)`
- **Source**: Karp–Rabin 1987 (fingerprints of all windows); every hash match verified.
- **Cost**: O(n + L·h) expected, h = windows whose fingerprint equals a stored one (O(nL) worst, e.g. aⁿ);
  O(n) extra for the table. **Lost**: suffix array + LCP (deterministic O(n log n)); much more code.
- **Output**: the start of the first occurrence of each distinct length-L substring that occurs ≥ 2 times
  (overlaps count), ascending. **Worked example** (hand): xyzxyzxy, L = 3 → [0, 1, 2].
- **Errors**: ErrLength (L < 1). **Tags**: LeetCode 187, 1044 (fixed-L step).

### ALG-TEXT-007 LongestCommonPrefix
`func LongestCommonPrefix[S ~string | ~[]byte](strs []S) int`
- **Source**: Knuth, TAOCP Vol. 3, 2nd ed., §6.3 (digital searching: the common prefix of a key set is the
  trie path above the first branching node). Walking that path column by column needs no trie.
- **Best known**: O(k·p + k) where p = result, k = len(strs); Ω the same (every string must be read to p).
- **Lost**: build a trie (ALG-TRIE) O(Σ|s|) and memory; sort then compare first/last O(Σ|s| log k).
- **Worked example** (hand): carpet carpool cartoon → 3. **Errors**: none; no strings → 0.

### ALG-TEXT-008 PalindromeRadii (Manacher — the one home)
`func PalindromeRadii[S ~string | ~[]byte](dst []int, s S) []int`
- **Source**: Manacher, JACM 22(3):346–351, 1975, doi:10.1145/321892.321896; Apostolico, Breslauer, Galil,
  "Parallel detection of all palindromes in a string", TCS 141:163–173, 1995 (Manacher's method yields all
  maximal palindromes, odd and even centres).
- **Best known**: O(n), Ω(n). **Approach**: Manacher over the implicit separator-interleaved string
  (no copy). **Lost**: expand around centre O(n²) (that one is ALG-DP's).
- **Output**: 2n−1 lengths; entry c is the longest palindrome centred at byte c/2 (c even) or between
  bytes (c−1)/2 and (c+1)/2 (c odd). **Worked example** (hand): abaaba → 1 0 3 0 1 6 1 0 3 0 1.

### ALG-TEXT-009 LongestPalindromicSubstring
`func LongestPalindromicSubstring[S ~string | ~[]byte](s S) (start, length int)` — via 008; leftmost on ties.
- **Worked example** (hand): xyracecarz → (2, 7). **Tags**: LeetCode 5, 647.

### ALG-TEXT-010 IsAnagram
`func IsAnagram(a, b string) bool` — rune multisets equal.
- **Source**: Bentley, *Programming Pearls*, 2nd ed., Column 2 (anagram signatures); CLRS §8.2 (counting).
- **Best known**: O(n) expected; ASCII fast path O(1) extra (a 128-entry table), else O(d) map.
- **Lost**: sort the runes (Bentley's signature; O(n log n)). **Tags**: LeetCode 242.

### ALG-TEXT-011 RunLengthEncode / ALG-TEXT-012 RunLengthDecode (PackBits)
`func RunLengthEncode(dst, src []byte) []byte` · `func RunLengthDecode(dst, src []byte, maxLen int) ([]byte, error)`
- **Format**: Apple PackBits as specified in TIFF Revision 6.0, Adobe, 3 June 1992, Section 9 "PackBits
  Compression", p. 42: header n (int8) 0…127 → n+1 literal bytes follow; −127…−1 → the next byte
  repeated 1−n times; −128 → no-op. Binary safe and unambiguous (a text "3a2b" form is ambiguous for
  digit symbols and has no standard).
- **Encoder policy** (TIFF's advice): runs of ≥ 3 → replicate; a run of 2 joins the literal in progress
  when it is followed by a single byte, otherwise replicate; runs longer than 128 split. Output ≤
  n + ⌊n/3⌋ + 1 bytes (each non-final, non-full literal header is followed by a replicate run, so it
  covers ≥ 3 input bytes; TIFF's "1 extra byte per 128" is the toolbox encoder's bound, not this policy's). **Decode**: ErrTruncated (header without its data), ErrTooLong (output would exceed
  maxLen), ErrMaxLen (maxLen < 0). Counts are at most 128 per header, so no count can overflow; maxLen is
  the cap against decompression bombs (expansion ≤ 64×).
- **Worked example** (hand, from the rules): AA AA AA 80 00 2A AA AA AA AA 80 00 2A 22 AA×10 →
  FE AA 02 80 00 2A FD AA 03 80 00 2A 22 F7 AA. **Time** O(n + output), O(1) extra.

### ALG-TEXT-013 ReverseWords
`func ReverseWords(s string) string`
- **Source**: Bentley, *Programming Pearls*, 2nd ed., Column 2.3 (rotation by reversals: reverse the
  whole, then each piece). Words = maximal runs of non-ASCII-white-space bytes; output joined by one space.
- **Best known**: O(n), one allocation. **Lost**: strings.Fields + slices.Reverse + strings.Join (same
  output for ASCII white space, but k+2 allocations; Fields also splits on Unicode spaces).

### ALG-TEXT-014 MinimumWindowSubstring
`func MinimumWindowSubstring(text, target string) (start, length int, found bool)` — byte multiset.
- **Source**: Laaksonen, *Competitive Programmer's Handbook*, 2018, §8.1 "Two pointers method" (L18):
  right pointer extends until the window covers target; left pointer advances while it still covers.
  Each pointer moves n times → O(n + m), 256-entry table O(1) extra. **Lost**: binary search on length
  O(n log n). Ties: leftmost. Empty target → (0, 0, true). **Tags**: LeetCode 76.

### ALG-TEXT-015 StringToInt
`func StringToInt(s string) (value int64, n int, err error)`
- **Source**: ISO/IEC 9899:2011 §7.22.1.4 paragraphs 2–5 (strtol, base 10), read in draft N1570;
  white space per §7.4.1.10 ("C" locale: space, \f \n \r \t \v).
- **Rule**: skip white space, optional sign, longest digit run; n = bytes consumed (0 when there are no
  digits: ErrNoDigits). Overflow: ErrOverflow with value 0 (not strtol's clamp). **Not strconv.ParseInt**:
  ParseInt rejects leading white space and trailing text and does not report how much it consumed.
- **Tags**: atoi, strtol. LeetCode 8 (its clamping rule is not ours).

### ALG-TEXT-016 IsJSONNumber
`func IsJSONNumber(s string) bool`
- **Source**: RFC 8259 (STD 90), 2017, §6 ABNF `number = [ minus ] int [ frac ] [ exp ]`.
- **Approach**: a 9-state DFA, O(n), no allocation. **Differs from strconv.ParseFloat**, which also
  accepts +1, .5, 5., 01, Inf, NaN, hex floats and underscores-free forms, and rejects out-of-range values
  that JSON's grammar accepts (1e999). json.Valid also accepts strings, objects and surrounding white space.

### ALG-TEXT-017 ColumnLabel / ALG-TEXT-018 ColumnNumber
`func ColumnLabel(number int) (string, error)` · `func ColumnNumber(label string) (int, error)`
- **Source**: Foster, "A Number System without a Zero-Symbol", Mathematics Magazine 21(1):39–41, 1947,
  doi:10.2307/3029479 (bijective base-k numeration; here k = 26, digits A…Z = 1…26).
- **Worked example** (hand): 1 A · 26 Z · 27 AA · 52 AZ · 53 BA · 702 ZZ · 703 AAA · 16384 XFD.
- **Errors**: ErrColumnNumber (< 1), ErrColumnLabel (empty or not A–Z), ErrOverflow (> MaxInt).
  Names avoid the product name (spreadsheet columns; tag excel). **Tags**: LeetCode 168, 171.

### ALG-TEXT-019 CompareSemVer
`func CompareSemVer(a, b string) (int, error)`
- **Source**: Semantic Versioning 2.0.0, https://semver.org/spec/v2.0.0.html (read 2026-11-12): §2 (X.Y.Z,
  no leading zeros), §9 (pre-release), §10 (build metadata ignored), §11 (precedence), the BNF.
- **Approach**: validate and compare in one scan without allocation; numeric parts compared as digit
  strings (no size limit, no overflow). **Lost**: plain dotted integers ("1.01" = "1.1") — no citable
  rule; HELD as CompareVersions. **Worked example**: §11's chain 1.0.0-alpha < … < 1.0.0.
- **Errors**: ErrSemVer.

### ALG-TEXT-020 WildcardMatch
`func WildcardMatch(text, pattern string) bool` — '?' one rune, '*' any run, whole text.
- **Source**: semantics: IEEE Std 1003.1-2017 (POSIX), XCU §2.13.1 (? and *, without brackets);
  method: NFA simulation, Thompson 1968 CACM 11(6):419–422, doi:10.1145/363347.363387; Aho, Lam, Sethi,
  Ullman, *Compilers*, 2nd ed., 2006, §3.7.2 Algorithm 3.22. Recurrence D[i][j] (text i runes, pattern j):
  '*' → D[i][j−1] ∨ D[i−1][j]; '?' or equal → D[i−1][j−1].
- **Cost**: O(n·m), O(m) extra (two rows). **Lost**: greedy with backtracking (O(1) space, no acceptable
  source found). **Not path.Match**: no brackets, no escapes, '*' crosses '/'. **Tags**: glob. LeetCode 44.

### ALG-TEXT-021 RegexMatch
`func RegexMatch(text, pattern string) (bool, error)` — '.' any rune, 'x*' zero or more x, whole text.
- **Source**: as 020 (Thompson 1968; ALSU §3.3 regular expressions, §3.7.2 simulation). Recurrence: x* →
  D[i][j−1] ∨ (x matches tᵢ ∧ D[i−1][j]); x → x matches tᵢ ∧ D[i−1][j−1].
- **Cost**: O(n·m), O(m) extra. **Errors**: ErrPattern ('*' first or after '*'). **Tags**: LeetCode 10.

### HELD / DROPPED
- **IntegerToRoman / RomanToInteger — HELD (no source)**: no primary reference for the strict subtractive
  grammar could be read (encyclopediaofmath.org returned 502; only secondary pages). Candidates for Devin:
  Ifrah, *The Universal History of Numbers*, Wiley 2000, ch. 16; Menninger, *Number Words and Number Symbols*.
- **StrStr — DROPPED**: identical to strings.Index / bytes.Index.
- **CompareVersions (plain dotted integers) — HELD**: no citable rule; CompareSemVer built instead.

---

# Catalog: algo/enumerate (SHORT = ENUM), new package

Builder W3-B12, 2026-11-12. Proposed by the combinatorics subject (doc/combinatorics_2026-11-12/PENDING_SHARED.md)
and planned in doc/taxonomy.md line 41; counts (n!, C(n,k), Catalan numbers) stay in num/combinatorics.
Every large output is an iter.Seq; breaking the range loop stops the generator at once. No recursion: every
generator is a loop over an explicit state array.

**Decision E1 (buffer reuse)**: the iterators yield ONE buffer, rewritten in place between steps. A caller
that keeps an element copies it (slices.Clone, string(b)). **Lost**: a fresh slice per step: one allocation
per visit makes enumeration allocation-bound (n! or 2ⁿ allocations); Knuth's algorithms are all "visit the
current array" algorithms, and Go's bufio.Scanner.Bytes uses the same convention. The yielded buffer must
not be modified by the caller (it is the generator's state, or mirrors it).

**Decision E2 (invalid sizes)**: functions whose only arguments are slices never fail (k < 0 or k > n is
the empty enumeration, C(n,k) = 0). Functions taking a bare count (CombinationSum, GenerateParentheses,
GrayCodeSequence) return (iter.Seq, error) so a negative or too-large count is an error, not silence.

### ALG-ENUM-001 NextPermutationInPlace
`func NextPermutationInPlace[T cmp.Ordered](values []T) bool`
- **Source**: Knuth, TAOCP Vol. 4A, Combinatorial Algorithms Part 1, Addison-Wesley, 2011, §7.2.1.2,
  Algorithm L (lexicographic permutation generation; steps L2 find j, L3 increase a_j, L4 reverse a_{j+1..n}).
  Works on multisets (Knuth states Algorithm L for a_1 ≤ … ≤ a_n with repeats).
- **Best known**: O(n) worst case per step, O(1) amortised over a full run of n! steps (Knuth §7.2.1.2,
  the analysis after Algorithm L: the comparisons and interchanges per step average a constant); Ω(1) per step.
- **Approach**: Algorithm L exactly. On the last permutation (non-increasing) it reverses to the first and
  returns false (so a do-while loop visits every permutation once). **Lost**: Heap's method (not lexicographic).
- **Worked example**: [1 2 3] → [1 3 2] → [2 1 3] → [2 3 1] → [3 1 2] → [3 2 1] → false, [1 2 3].
- **Errors**: none; empty and one-element input return false. **Modifies input**: yes (name says so).
- **Tags**: permutation, lexicographic, next-permutation, combinatorics. LeetCode 31.

### ALG-ENUM-002 Permutations
`func Permutations[T any](values []T) iter.Seq[[]T]`
- **Source**: Knuth §7.2.1.2 Algorithm L run on the index array 0..n−1 (all n! orderings of the positions in
  lexicographic order); every swap and reversal is mirrored on a copy of the values, so T needs no ordering.
- **Best known**: O(1) amortised per permutation (Algorithm L), O(n! ) total; Ω(n!) output.
- **Approach**: as above. **Lost**: Heap's algorithm (Sedgewick, "Permutation Generation Methods", ACM
  Computing Surveys 9(2):137–164, 1977, doi:10.1145/356689.356692): one swap per step but not lexicographic,
  and ordered output is what callers compare against. **Worked example**: [a b c] → abc acb bac bca cab cba.
- **Errors**: none; empty input yields one empty permutation. Input copied when iteration starts; never modified.
- **Tags**: permutation, backtracking, enumeration, iterator. LeetCode 46.

### ALG-ENUM-003 PermutationsUnique
`func PermutationsUnique[T cmp.Ordered](values []T) iter.Seq[[]T]`
- **Source**: Knuth §7.2.1.2 Algorithm L on the sorted multiset (repeats produce each distinct
  arrangement once). Calls ALG-ENUM-001.
- **Best known**: O(1) amortised per distinct permutation; O(n log n) to sort first.
- **Lost**: generate all n! and deduplicate with a set (n!/∏mᵢ! of them may be far fewer than n!).
- **Worked example**: [1 1 2] → 112 121 211. **Errors**: none. **Tags**: permutation, multiset, duplicates. LeetCode 47.

### ALG-ENUM-004 Combinations
`func Combinations[T any](values []T, k int) iter.Seq[[]T]`
- **Source**: Kreher, Stinson, *Combinatorial Algorithms: Generation, Enumeration, and Search*, CRC Press,
  1999, §2.3.1 (lexicographic ordering of k-subsets; the successor algorithm kSubsetLexSuccessor — algorithm
  number not confirmed first-hand); Knuth §7.2.1.3 (generating all combinations; Knuth's Algorithm L
  visits them in colex order, the order lost here).
- **Best known**: O(1) amortised per combination for the index successor, O(k) to refresh the yielded
  values in the worst step; Ω(C(n,k)) output.
- **Approach**: index array c₀ < … < c_{k−1}; find the rightmost i with cᵢ < n−k+i, increment it, reset the
  tail to consecutive values. **Lost**: Knuth Algorithm L/T (colex order; faster by a constant, but the
  lexicographic order is the one callers expect).
- **Worked example**: [1 2 3 4], k = 2 → 12 13 14 23 24 34. **Errors**: none; k < 0 or k > n yields nothing;
  k = 0 yields one empty combination. **Tags**: combination, k-subset, lexicographic. LeetCode 77.

### ALG-ENUM-005 Subsets
`func Subsets[T any](values []T) iter.Seq[[]T]`
- **Source**: Knuth §7.2.1.1, Algorithm M (mixed-radix generation) with every radix 2: binary counting
  over membership flags, values[0] the least significant digit.
- **Best known**: O(n) per subset to rebuild the yielded slice (output size), O(1) amortised counter step.
- **Lost**: Gray order (Knuth §7.2.1.1 Algorithm G): one flag changes per step, but an inserted or removed
  element in the middle of the yielded slice still costs O(n), so nothing is saved; a uint64 mask (limits n to 63).
- **Worked example**: [a b c] → {} a b ab c ac bc abc. **Errors**: none. **Tags**: subset, power-set, binary-counting. LeetCode 78.

### ALG-ENUM-006 SubsetsUnique
`func SubsetsUnique[T cmp.Ordered](values []T) iter.Seq[[]T]`
- **Source**: Knuth §7.2.1.1 Algorithm M: a sub-multiset of a multiset with distinct values d₁<…<d_r of
  multiplicities m₁…m_r is a mixed-radix tuple (a₁…a_r), 0 ≤ aⱼ ≤ mⱼ.
- **Best known**: O(n) per sub-multiset (output), ∏(mⱼ+1) of them; O(n log n) sort first.
- **Lost**: all 2ⁿ subsets plus a set of seen keys (exponential memory).
- **Worked example**: [1 2 2] → {} 1 2 12 22 122. **Errors**: none. **Tags**: subset, multiset, duplicates. LeetCode 90.

### ALG-ENUM-007 CombinationSum
`func CombinationSum(candidates []int, target int) (iter.Seq[[]int], error)`
- **Source**: backtracking: Knuth, TAOCP Vol. 4B, Combinatorial Algorithms Part 2, Addison-Wesley, 2022,
  §7.2.2, Algorithm B (basic backtrack); Skiena, *The Algorithm Design Manual*, 2nd ed., Springer, 2008,
  §7.1 "Backtracking". The search tree: the multisets of candidates in non-decreasing order, pruned when the
  next candidate exceeds the remainder (candidates sorted ascending).
- **Best known**: output-sensitive; O(depth) per tree node, depth ≤ target / min(candidates). No
  polynomial bound in general (the count is a partition count).
- **Approach**: iterative backtracking on an explicit stack (no recursion). **Lost**: recursion (depth
  target/min could exceed the goroutine stack for large targets); DP counting tables (count, not list).
- **Worked example**: candidates [2 3 5], target 8 → [2 2 2 2] [2 3 3] [3 5].
- **Errors**: ErrNonPositive (a candidate ≤ 0: infinitely many combinations), ErrNegative (target < 0).
  Target 0 yields one empty combination. Duplicate candidates are merged. **Tags**: backtracking,
  combination-sum, partition. LeetCode 39.

### ALG-ENUM-008 GenerateParentheses
`func GenerateParentheses(pairs int) (iter.Seq[[]byte], error)`
- **Source**: Knuth, TAOCP Vol. 4A, 2011, §7.2.1.6, Algorithm P (nested parentheses in lexicographic order).
  Order: from ()()…() to ((…)), i.e. increasing when ')' is taken to precede '(' (hand-simulated for n = 2, 3).
- **Best known**: O(1) amortised per string (Knuth §7.2.1.6); C_n = Catalan(n) strings, Ω(C_n) output.
- **Lost**: backtracking on (open, close) counts (O(n) per string, recursion depth 2n).
- **Worked example**: 3 → ()()() ()(()) (())() (()()) ((())).
- **Errors**: ErrNegative (pairs < 0), ErrTooLarge (2·pairs+1 overflows int). 0 yields one empty string.
- **Tags**: parentheses, catalan, dyck-word, enumeration. LeetCode 22.

### ALG-ENUM-009 GrayCodeSequence
`func GrayCodeSequence(bits int) (iter.Seq[uint64], error)`
- **Source**: Knuth §7.2.1.1, the binary reflected Gray code g(k) = k ⊕ ⌊k/2⌋; Gray, "Pulse code
  communication", US Patent 2,632,058, 1953.
- **Best known**: O(1) per code word; 2^bits words.
- **Worked example**: 3 → 0 1 3 2 6 7 5 4. **Errors**: ErrBitsRange (bits < 0 or > 64). bits = 0 yields 0 once.
- **Tags**: gray-code, bits, enumeration. LeetCode 89.

---

# Catalog: algo/sample (SHORT = SAMPLE), new package

Builder W3-B12, 2026-11-12. NOT in doc/taxonomy.md: created on the lead's brief (W3-B12); the taxonomy row
`ALG | algo/sample | SAMPLE | random shuffles and samples from a caller-supplied source` is for the lead/Devin.
Every function takes a *rand.Rand (math/rand/v2) from the caller; there is no global source, no clock,
no goroutine. Nil rng is ErrNilRand. Uniformity is exact given a uniform source: math/rand/v2's IntN and
Uint64N are unbiased (Lemire's multiply-and-reject, per the standard library documentation).

### ALG-SAMPLE-001 ShuffleInPlace
`func ShuffleInPlace[T any](rng *rand.Rand, values []T) error`
- **Source**: Knuth, TAOCP Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §3.4.2,
  Algorithm P (shuffling); Durstenfeld, "Algorithm 235: Random permutation", Communications of the ACM
  7(7):420, 1964, doi:10.1145/364520.364540 (confirmed on Crossref 2026-11-12). Fisher & Yates (1938) is
  the pencil-and-paper original.
- **Best known**: O(n) time, n−1 random draws, O(1) extra; Ω(n) (every element may move).
- **Approach**: for j = n−1 down to 1, swap values[j] with values[rng.IntN(j+1)]. **Lost**: sort by random
  keys (O(n log n), ties bias it); rand.Rand.Shuffle (same algorithm, but takes a swap callback; ours is
  generic, documents the source and returns ErrNilRand instead of panicking).
- **Caveat**: a source with s bits of state reaches at most 2^s of the n! orders (PCG: 2^128 < 35!).
- **Worked example**: chi-square-free check: over 60 000 shuffles of [a b c] each of the 6 orders appears
  10 000 ± 5σ times. **Errors**: ErrNilRand. **Modifies input**: yes (name says so).
- **Tags**: shuffle, fisher-yates, knuth-shuffle, random-permutation. LeetCode 384.

### ALG-SAMPLE-002 ReservoirSample
`func ReservoirSample[T any](dst []T, rng *rand.Rand, seq iter.Seq[T], k int) ([]T, error)`
- **Source**: Vitter, "Random sampling with a reservoir", ACM Transactions on Mathematical Software
  11(1):37–57, 1985, doi:10.1145/3147.3165, Algorithm R (credited there to Waterman); Knuth Vol. 2, 3rd ed.,
  §3.4.2, Algorithm R (reservoir sampling).
- **Best known**: Algorithm R: O(N) time, N−k draws, O(k) space; Vitter's Algorithm Z: O(k(1+log(N/k)))
  draws expected, same output distribution.
- **Approach**: Algorithm R: keep the first k; the t-th item (t > k, 1-based) replaces a uniform slot with
  probability k/t. **Lost**: Algorithm Z (fewer random draws, but a sequence must still be read item by
  item through iter.Seq, so the O(N) read dominates and Z's skip arithmetic gains little; Algorithm L, Li
  1994, likewise).
- **Order**: the reservoir's slot order, not random; shuffle if order matters.
- **Worked example**: sample 2 of 4 items 120 000 times: each of the 6 pairs appears 20 000 ± 5σ times.
- **Errors**: ErrNilRand, ErrNilSeq, ErrNegative (k < 0). k = 0 returns dst unchanged without reading seq.
  Fewer than k items: all of them. **Tags**: reservoir-sampling, streaming, random-sample. LeetCode 382, 398.

### ALG-SAMPLE-003 NewWeightedPicker · ALG-SAMPLE-004 WeightedPicker.Pick
`func NewWeightedPicker[W integer](weights []W) (*WeightedPicker[W], error)` · `func (p *WeightedPicker[W]) Pick(rng *rand.Rand) (int, error)`
- **Source**: the inversion method on a discrete distribution: Devroye, *Non-Uniform Random Variate
  Generation*, Springer, 1986, ch. III, §2 (inversion by sequential or binary search of the cumulative
  table; subsection not confirmed first-hand); Knuth Vol. 2, 3rd ed., §3.4.1 (discrete distributions). Prefix sums: ALG-PREFIX-001; search:
  ALG-SEARCH-002 UpperBound — both called, not copied.
- **Best known**: O(n) build, O(log n) per pick by binary search; O(1) per pick by Walker's alias method
  (Walker, "An efficient method for generating discrete random variables with general distributions", ACM
  TOMS 3(3):253–256, 1977, doi:10.1145/355744.355749). **Lost**: the alias method — O(1) picks, but its
  table is built in floating point (or rationals) and a pick needs two draws; with integer weights the
  inversion is exact and O(log n) is fast in practice. Kept as a later variant.
- **Approach**: prefix[i] = w₀+…+w_{i−1}; draw r uniform in [0, total); answer UpperBound(prefix, r) − 1.
  Zero weights are never picked.
- **Worked example**: weights [1 3] → index 1 three times as often as 0 (r ∈ {0} → 0, r ∈ {1,2,3} → 1).
- **Errors**: ErrNegativeWeight, ErrZeroTotal (empty or all zero), ErrOverflow (total does not fit W);
  Pick: ErrNilRand, ErrNilPicker. Weights are copied (prefix sums) and never modified.
- **Tags**: weighted-random, inversion, prefix-sum, binary-search, discrete-distribution. LeetCode 528.

### ALG-SAMPLE-005 RejectionSample
`func RejectionSample[T any](rng *rand.Rand, propose func(*rand.Rand) T, accept func(candidate T, u float64) bool, maxTries int) (T, error)`
- **Source**: von Neumann, "Various techniques used in connection with random digits", in *Monte Carlo
  Method*, National Bureau of Standards Applied Mathematics Series 12, 1951, pp. 36–38; Devroye 1986,
  ch. II, §3 (the rejection method; the number of tries is geometric with mean c).
- **Best known**: expected c tries, c = sup f/(c·g) bound of the proposal; the cap turns an unbounded
  loop into an error.
- **Approach**: repeat: x ← propose(rng), u ← rng.Float64() in [0, 1); return x if accept(x, u). The caller
  writes accept as u < f(x)/(c·g(x)), or as a plain membership test ignoring u.
- **Worked example**: points uniform in the unit disc from the square [−1,1]²: acceptance rate π/4 (checked
  to 4σ over 100 000 tries). **Errors**: ErrNilRand, ErrNilFunc, ErrMaxTries (maxTries < 1 is
  ErrNonPositive). **Tags**: rejection-sampling, accept-reject, monte-carlo. LeetCode 470.

### ALG-SAMPLE-006 RandomElement
`func RandomElement[T any](rng *rand.Rand, seq iter.Seq[T]) (T, error)`
- **Name**: the brief's "LinkedListRandomNode". It takes any iter.Seq (a list's values, a channel drained by
  the caller, a file's lines), so "LinkedList" in the name would mislead; the problem is "a uniform element
  of a sequence of unknown length". Recorded for the lead.
- **Source**: Vitter 1985, Algorithm R with k = 1; calls ReservoirSample (ALG-SAMPLE-002).
- **Best known**: O(N) time (one pass), O(1) space. **Errors**: ErrNilRand, ErrNilSeq, ErrEmpty.
- **Worked example**: 5 items, 50 000 draws, each 10 000 ± 5σ. **Tags**: reservoir-sampling, streaming, random-element. LeetCode 382.

### ALG-SAMPLE-007 ShuffleDeterministic
`func ShuffleDeterministic(dst []int, n int, seed uint64) ([]int, error)`
- **Source**: ShuffleInPlace (ALG-SAMPLE-001) on 0…n−1 with a PCG source (O'Neill, "PCG: A Family of
  Simple Fast Space-Efficient Statistically Good Algorithms for Random Number Generation", Harvey Mudd
  College HMC-CS-2014-0905, 2014; math/rand/v2 NewPCG) seeded (seed, fixed stream constant).
- **Not cryptographic**: anyone who knows the seed knows the permutation; never use it for secrets.
  Reproducible within one Go toolchain; math/rand/v2 does not promise identical IntN output across Go
  releases, so tests check same-seed equality and the permutation property, not golden values.
- **Errors**: ErrNegative (n < 0), ErrTooLarge (len(dst)+n ints exceed a slice). Appends to dst. **Tags**: shuffle, seeded, deterministic, testing.

---

# Catalog: algo/matrix (SHORT = MATRIX), new package

Builder W3-B12, 2026-11-12. NOT in doc/taxonomy.md: created on the lead's brief; the taxonomy row
`ALG | algo/matrix | MATRIX | grid algorithms on [][]T: rotation, transposition, sparse products, cellular automata`
is for the lead/Devin. Distinction from num/linalg: num/linalg owns float64 linear algebra (solves,
decompositions, condition numbers); algo/matrix owns algorithms on rectangular [][]T grids that are
type-generic or combinatorial. SearchSortedMatrix is ALG-SEARCH-009/010 (L9), not here.

Grid conventions: grid[r][c], r the row; every row the same length (else ErrRagged); a grid with no rows
is empty. Functions producing a grid reuse dst when it already has the result's shape, else allocate
once (one backing array + one row-header slice); dst must not share memory with the input.

### ALG-MATRIX-001 Rotate90InPlace
`func Rotate90InPlace[T any](grid [][]T) error`
- **Source**: Knuth, TAOCP Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §1.3.3
  "Applications to permutations" (a permutation applied in place cycle by cycle, with one temporary per
  cycle; ex. 12 there is in-place transposition). The clockwise rotation (r, c) → (c, n−1−r) is a
  permutation of the n² cells whose cycles are {(r,c), (c,n−1−r), (n−1−r,n−1−c), (n−1−c,r)}: length 4,
  or 1 at the centre of odd n (proof: applying the map four times is the identity, and it fixes a cell
  only when r = c = (n−1)/2). One representative per cycle: r < ⌊n/2⌋, r ≤ c < n−1−r.
- **Best known**: O(n²), every cell moves; O(1) extra. **Lost**: transpose then reverse each row (same
  bound, 2 passes and ~1.5× the writes; shown in the test as the oracle).
- **Worked example**: [[1 2 3] [4 5 6] [7 8 9]] → [[7 4 1] [8 5 2] [9 6 3]].
- **Errors**: ErrNotSquare (a row length ≠ number of rows). Empty grid: no-op. **Modifies input**: yes.
- **Tags**: matrix, rotate, in-place, permutation-cycles, grid. LeetCode 48.

### ALG-MATRIX-002 Transpose
`func Transpose[T any](dst, grid [][]T) ([][]T, error)`
- **Source**: Cormen et al., *Introduction to Algorithms*, 3rd ed., MIT Press, 2009, Appendix D.1 (the
  transpose Aᵀ, (Aᵀ)ᵢⱼ = aⱼᵢ). **Best known**: O(mn), Ω(mn). Result n×m.
- **Lost**: Knuth §1.3.3 ex. 12 in-place transposition of a rectangular array stored flat: [][]T rows
  cannot change length in place, and a cycle-following transpose is slower in cache than a copy.
- **Worked example**: [[1 2 3] [4 5 6]] → [[1 4] [2 5] [3 6]]. **Errors**: ErrRagged.
- **Tags**: matrix, transpose, grid. LeetCode 867.

### ALG-MATRIX-003 NewCSR · ALG-MATRIX-004 SparseMultiply
`func NewCSR[T number](dense [][]T) (*CSR[T], error)` · `func SparseMultiply[T number](a, b *CSR[T]) (*CSR[T], error)`
- **Source**: Gustavson, "Two Fast Algorithms for Sparse Matrices: Multiplication and Permuted
  Transposition", ACM Transactions on Mathematical Software 4(3):250–269, 1978,
  doi:10.1145/355791.355796 (confirmed on Crossref 2026-11-12): row-by-row product with a dense
  accumulator and a marker array (the "sparse accumulator").
- **Best known**: O(m + flops + Σᵢ nnz(Cᵢ) log nnz(Cᵢ)) here: flops = Σ over nonzeros a_ik of nnz(row k of
  B); the log term sorts each output row's columns so the result is canonical (Gustavson's output is
  unordered; his permuted transposition sorts in O(nnz) — lost for simplicity, a transpose twice).
- **Accumulator**: dense array of b.Cols when b.Cols ≤ 4·nnz(B)+64, else a map (exact either way); so the
  extra memory is O(nnz(B) + b.Cols) bounded by O(nnz(B)) and an exported huge Cols field cannot make
  make() panic. Explicit zeros of the product are dropped (−0.0 too). Duplicate column entries in an input
  row are summed; unsorted rows are fine.
- **Numbers**: signed integers (products and sums overflow-checked → ErrOverflow) and floats (IEEE: Inf and
  NaN propagate as arithmetic does; documented, not an error).
- **Lost**: map-of-rows representation (hash per entry; CSR is the format Gustavson and every BLAS uses);
  dense O(mnp) multiply.
- **Worked example**: A = [[1 0 2] [0 3 0]], B = [[0 4] [5 0] [6 0]] → C = [[12 4] [15 0]].
- **Errors**: ErrRagged (NewCSR), ErrInvalidCSR (a malformed CSR: lengths, RowStart not non-decreasing,
  column out of range, negative dimension), ErrShape (a.Cols ≠ b.Rows), ErrOverflow.
- **Tags**: sparse-matrix, csr, matrix-multiply, gustavson. LeetCode 311.

### ALG-MATRIX-005 GameOfLifeStep
`func GameOfLifeStep(dst, grid [][]bool) ([][]bool, error)`
- **Source**: Gardner, "Mathematical Games: The fantastic combinations of John Conway's new solitaire game
  'life'", Scientific American 223(4):120–123, 1970, doi:10.1038/scientificamerican1070-120 (confirmed on
  Crossref 2026-11-12): survival with 2 or 3 neighbours, birth with exactly 3, death otherwise; the 8
  neighbours of the Moore neighbourhood.
- **Boundary (ASSUMPTION, for the lead)**: Gardner's board is infinite; on a finite grid every cell outside
  is dead (no wrap-around). A torus is a later variant.
- **Best known**: O(mn) per generation (each cell reads 8 neighbours). **Approach**: double buffer (dst).
  **Lost**: in place with a 2-bit state encoding (needs an integer cell type; bool cannot hold "was alive,
  now dead", and the double buffer is what a caller animating generations keeps anyway).
- **Worked example**: Gardner's blinker: a row of three live cells becomes a column of three, and back.
- **Errors**: ErrRagged. **Tags**: game-of-life, cellular-automaton, grid, simulation. LeetCode 289.

### ALG-MATRIX-006 MaxRectangleSumNoLargerThanK
`func MaxRectangleSumNoLargerThanK[T signed](grid [][]T, k T) (T, error)`
- **Source**: the 2-D → 1-D reduction over pairs of rows (columns summed between them): Bentley,
  *Programming Pearls*, 2nd ed., Addison-Wesley, 2000, Column 8, Problems (the two-dimensional maximum-sum
  subarray; problem number not confirmed first-hand). The 1-D step "max of P[j] − P[i] ≤ k over i < j"
  is the cross-pair count of merge sort: Kleinberg, Tardos, *Algorithm Design*, Addison-Wesley, 2006, §5.3
  "Counting Inversions" (Merge-and-Count: pairs i < j split across halves are examined with two pointers
  while merging sorted halves); merge sort: Knuth Vol. 3, 2nd ed., §5.2.4. Prefix sums via ALG-PREFIX-001.
  Proof of the cross step: with both halves sorted ascending, for a right value r the best left value is the
  smallest l ≥ r − k; as r grows r − k grows, so the pointer only moves right (as in Merge-and-Count).
- **Best known**: O(s² · L log L), s = min(m, n), L = max(m, n) (iterate pairs along the short side).
- **Approach**: bottom-up merge sort (no recursion). **Lost**: a balanced ordered set with a ceiling query
  per prefix (same bound; algo/tree's BST is unbalanced, O(L) per query on sorted input, and has no ceiling
  query); Kadane (does not handle the ≤ k bound).
- **Worked example** (own): [[1 0 1] [0 −2 3]], k = 2 → 2 (the rectangle [[0 1] [−2 3]]); k = 3 → 3.
- **Errors**: ErrRagged, ErrEmpty (no cells), ErrNotFound (every rectangle sums to more than k), ErrOverflow
  (a partial column or prefix sum leaves int64, or the answer does not fit T — the L28 policy).
- **Tags**: matrix, max-sum-rectangle, prefix-sum, merge-sort, divide-and-conquer. LeetCode 363.

## HELD (no source)
- **SpiralOrder**: the only descriptions found (2026-11-12) are problem sites and interview-question books
  (Aziz–Lee–Prakash, *Elements of Programming Interviews*; McDowell, *Cracking the Coding Interview*);
  no textbook or paper. HELD until the lead or Devin accepts an interview book as a source class.
- **SetZeroesInPlace**: the O(1)-extra method (first row and column as markers) is found only in
  McDowell, *Cracking the Coding Interview*, 6th ed., 2015, Problem 1.8, and problem sites. HELD as above.

---

# algo/geometry2d (SHORT = GEO2D)

Builder W4-B13. New package, assigned by the lead's brief; not yet in `doc/taxonomy.md` (lead to add:
`ALG | algo/geometry2d | GEO2D | computational geometry over planar point sets: orientation, hulls, closest
pair, point location, segment intersection`). `num/geometry` keeps the measurement formulas (PolygonArea,
PointDistance, RectangleArea = w·h); this package does not import `num/*`.

## Package decision: exact int64 coordinates only (builder choice, for the lead to log)

`Point{X, Y int64}`. Every predicate is exact on the whole int64 range: a coordinate difference is carried as
a 64-bit magnitude and a sign (65 bits), a product of two differences as an unsigned 128-bit magnitude
(`math/bits.Mul64`) and a sign, and a sum of two squares in 192 bits. Nothing overflows, so no input is
rejected for size. **Lost**: (a) `float64` points with Shewchuk's adaptive predicates (Shewchuk, "Adaptive
precision floating-point arithmetic and fast robust geometric predicates", Discrete & Computational Geometry
18:305–363, 1997, doi:10.1007/PL00009321): several hundred lines of expansion arithmetic for a predicate the
integer type gives exactly; callers with real-valued data scale to a fixed grid (for example micrometres in
int64) and lose nothing a float would keep. (b) `float64` with an epsilon: not robust (hull and
point-location answers become inconsistent near degeneracies). (c) plain int64 products with a coordinate cap (|coordinate| < 2³⁰ for the cross product to fit):
an invented cap that rejects valid input, for a saving of a few nanoseconds. The
consequence: no float segment-intersection POINT is offered (it would need the float type); the predicate
SegmentsIntersect is exact.

Units: coordinates are in any one length unit u the caller chooses; areas are in u².

| ID | Name | Signature |
|---|---|---|
| ALG-GEO2D-001 | Orientation | `func Orientation(a, b, c Point) int` |
| ALG-GEO2D-002 | ConvexHull | `func ConvexHull(dst, points []Point) []Point` |
| ALG-GEO2D-003 | ClosestPair | `func ClosestPair(points []Point) (i, j int, err error)` |
| ALG-GEO2D-004 | PointInPolygon | `func PointInPolygon(polygon []Point, q Point) (Location, error)` |
| ALG-GEO2D-005 | SegmentsIntersect | `func SegmentsIntersect(s, t Segment) bool` |
| ALG-GEO2D-006 | MaxPointsOnLine | `func MaxPointsOnLine(points []Point) int` |
| ALG-GEO2D-007 | RectanglesUnionArea | `func RectanglesUnionArea(r1, r2 Rect) (uint64, error)` |
| ALG-GEO2D-008 | HasVerticalSymmetryAxis | `func HasVerticalSymmetryAxis(points []Point) bool` |
| — | PerfectRectangle | HELD (no source) |

### ALG-GEO2D-001 Orientation
- **Source**: CLRS (Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009),
  §33.1 "Line-segment properties", the cross product (p1 − p0) × (p2 − p0) and DIRECTION. O'Rourke,
  Computational Geometry in C, 2nd ed., Cambridge University Press, 1998, §1.3 (Area2 / Left, the sign of
  twice the signed triangle area).
- **Definition**: the sign of (b − a) × (c − a): +1 when a, b, c turn counterclockwise (c left of a→b),
  −1 clockwise, 0 collinear (including any repeated point).
- **Best known**: O(1), exact. **Approach**: 65-bit differences, 128-bit products via `bits.Mul64`, signed
  compare. **Lost**: int64 products (overflow once |coordinate| ≥ 2³⁰); `math/big` (allocates).
- **Worked example**: hand: (0,0),(4,0),(0,3) → +1; swapped → −1; (0,0),(1,1),(5,5) → 0. Extremes:
  (MinInt64,MinInt64),(MaxInt64,MinInt64),(MinInt64,MaxInt64) → +1 (cross ≈ 2¹²⁸, exact).
- **Errors**: none. No allocation. **Tags**: orientation, ccw, cross-product, predicate.

### ALG-GEO2D-002 ConvexHull
- **Source**: Andrew, "Another efficient algorithm for convex hulls in two dimensions", Information
  Processing Letters 9(5):216–219, 1979, doi:10.1016/0020-0190(79)90072-3 (verified on Crossref
  2026-11-12) — the monotone chain. CLRS 3rd ed. §33.3 (Graham scan; the Ω(n log n) lower bound is stated
  there, with reference to Yao 1981).
- **Best known**: O(n log n), Ω(n log n) in the algebraic decision-tree model; output-sensitive O(n log h)
  (Kirkpatrick–Seidel 1986, Chan 1996) is better only when h ≪ n. Space O(n).
- **Approach**: clone, sort by (X, Y), drop duplicates, build lower then upper chain popping on
  Orientation ≤ 0. **Collinear policy**: only strict convex vertices; points in the interior of a hull edge
  and duplicates are dropped. Output counterclockwise starting at the least (X, Y) point; all collinear → the
  two extremes; one distinct point → that point; empty → nothing. **Lost**: Graham scan (angular sort is
  awkward exactly; same bound); Jarvis march O(nh); Chan O(n log h) (much more code, wins only for tiny h).
- **Output**: appended to dst (D-003). **Input modified**: no. **Errors**: none.
- **Worked example**: hand: square (0,0),(4,0),(4,4),(0,4) with (2,2) inside and (2,0) on an edge → (0,0),
  (4,0),(4,4),(0,4).
- **Tags**: convex-hull, monotone-chain. LeetCode: 587 (that problem keeps edge points; ours does not).

### ALG-GEO2D-003 ClosestPair
- **Source**: Shamos, Hoey, "Closest-point problems", 16th Annual Symposium on Foundations of Computer Science
  (SFCS 1975), IEEE, pp. 151–162, doi:10.1109/SFCS.1975.8 (verified on Crossref 2026-11-12). CLRS 3rd ed.
  §33.4 "Finding the closest pair of points" (divide and conquer with presorting; at most 7 strip
  successors).
- **Best known**: Θ(n log n) in the algebraic decision-tree model (element uniqueness reduces to it,
  Shamos–Hoey 1975); O(n) expected with randomisation and floor (Rabin 1976) is outside the model and needs
  hashing of real cells. Space O(n).
- **Approach**: sort indices by X, recurse, merge by Y on the way back (merge sort, so O(n log n) not
  O(n log² n)); strip test with exact squared distances in 192 bits. Recursion depth ⌈log₂ n⌉.
  **Lost**: O(n²) all pairs; sweep line with a balanced tree (same bound, needs an ordered set).
- **Result**: indices i < j of one closest pair (Euclidean). Duplicates are a pair at distance 0.
- **Errors**: `ErrTooFewPoints` (n < 2). **Input modified**: no.
- **Worked example**: hand: (0,0),(10,10),(3,4),(11,12),(20,0) → the pair (10,10),(11,12), indices 1, 3.
- **Tags**: closest-pair, divide-and-conquer, nearest. LeetCode: none exact.

### ALG-GEO2D-004 PointInPolygon
- **Source**: Haines, "Point in Polygon Strategies", in Heckbert (ed.), Graphics Gems IV, Academic Press, 1994,
  pp. 24–46 (crossings test, +X ray, half-open vertex rule). O'Rourke, Computational Geometry in C, 2nd ed.,
  1998, §7.4 "Point in Polygon" (§7.4.2 ray crossings, p. 239 ff.; confirmed from the book's table of
  contents 2026-11-12).
- **Best known**: Θ(n) for one query on an unpreprocessed polygon (every edge may matter). O(1) extra.
- **Approach**: one pass; each edge first tested for q on it (exact, → OnBoundary); otherwise an edge whose
  endpoints straddle y = q.Y (one strictly above, the other not) is a crossing when q is strictly on the side
  of it facing −X, decided by Orientation (no division). **Boundary policy**: OnBoundary is reported
  exactly, before parity. Even–odd rule for self-intersecting input. Either vertex order; closed implicitly.
  **Lost**: winding number (§7.4.1; same cost, differs only for self-intersecting polygons; a WindingNumber
  variant may follow with its own ID); float ray test (inexact at vertices).
- **Errors**: `ErrTooFewPoints` (< 3 vertices). No allocation. **Input modified**: no.
- **Worked example**: hand: U shape (0,0),(6,0),(6,6),(4,6),(4,2),(2,2),(2,6),(0,6): (1,4) Inside, (3,4)
  Outside (in the notch), (3,2) OnBoundary, (5,6) OnBoundary, (1,2) Inside (its ray runs along the notch's
  floor through two vertices: the half-open rule counts 3 crossings).
- **Tags**: point-in-polygon, crossing-number, ray-casting, even-odd.

### ALG-GEO2D-005 SegmentsIntersect
- **Source**: CLRS 3rd ed. §33.1, SEGMENTS-INTERSECT and ON-SEGMENT (closed segments; collinear overlap and
  touching endpoints count).
- **Best known**: O(1), exact with Orientation. No allocation. **Lost**: solving the 2×2 system in floats.
- **Degenerate**: a zero-length segment is a point; it intersects when it lies on the other.
- **Worked example**: hand: (0,0)–(4,4) × (0,4)–(4,0) true; (0,0)–(2,0) × (3,0)–(5,0) false; (0,0)–(3,0) ×
  (2,0)–(5,0) true (overlap); (0,0)–(2,2) × (2,2)–(4,0) true (touch).
- **Errors**: none. **Tags**: segment-intersection, predicate.

### ALG-GEO2D-006 MaxPointsOnLine
- **Source**: the problem (maximum number of collinear points) and its O(n²) bound: Gajentaan, Overmars, "On a
  class of O(n²) problems in computational geometry", Computational Geometry 5(3):165–185, 1995,
  doi:10.1016/0925-7721(95)00022-2 (verified on Crossref 2026-11-12): deciding whether any three are collinear
  is 3SUM-hard, so o(n²) is not expected. The per-anchor step is proved here (L15 pattern): a line through
  ≥ 2 of the points has a first point p (least index); every other point q on it gives the same direction
  q − p up to sign, and the direction reduced by gcd of |dx|, |dy| with a sign flag is a canonical key for
  that line through p; so the largest bucket for anchor p, plus copies of p, plus p, is the largest line
  whose first point is p. Gcd by Stein's binary algorithm (Stein, "Computational problems associated with Racah algebra", J. Comput. Phys. 1(3):397–405, 1967,
  doi:10.1016/0021-9991(67)90047-2), private (L4: the exported gcd belongs to num/integer).
- **Best known**: O(n²) (3SUM-hard). This: O(n²) expected (hash map), O(n) extra.
  **Lost**: sort directions per anchor (O(n² log n), deterministic); slope as float (wrong for steep/near
  slopes).
- **Duplicates**: counted separately (a multiset): every copy is on every line through that point.
- **Worked example**: hand: (0,0),(1,1),(2,2),(3,3),(0,1),(1,0),(2,0) → 4. Empty → 0, one → 1.
- **Errors**: none. **Tags**: collinear, lines, gcd. LeetCode: 149.

### ALG-GEO2D-007 RectanglesUnionArea
- Named per L10 (not RectangleArea: NUM-GEOM owns w·h).
- **Source**: inclusion–exclusion |A ∪ B| = |A| + |B| − |A ∩ B| (Stanley, Enumerative Combinatorics Vol. 1,
  2nd ed., Cambridge, 2012, §2.1, for finite sets; for areas it is additivity of Lebesgue measure). The
  intersection of two axis-aligned rectangles is the product of the intersections of their x and y
  intervals (proof: a point is in both iff its x is in both x ranges and its y in both y ranges). Klee's
  measure problem (Klee, Amer. Math. Monthly 84(4):284–285, 1977, doi:10.1080/00029890.1977.11994336) is
  the n-rectangle generalisation (sweep, O(n log n)); not built here.
- **Best known**: O(1). Exact: widths as 64-bit magnitudes, areas in 128 bits, sum in 192.
- **Result**: uint64 area in u²; `ErrOverflow` when the union exceeds 2⁶⁴ − 1. Degenerate (zero-width)
  rectangles are allowed (area 0). **Errors**: `ErrInvertedRect` (Min > Max on an axis), `ErrOverflow`.
- **Worked example**: hand: [−3,3]×[0,4] ∪ [0,9]×[−1,2] = 24 + 27 − 6 = 45.
- **Tags**: rectangle, union-area, inclusion-exclusion. LeetCode: 223.

### ALG-GEO2D-008 HasVerticalSymmetryAxis
- **Source**: Atallah, "On symmetry detection", IEEE Transactions on Computers C-34(7):663–666, 1985,
  doi:10.1109/TC.1985.1676605 (all symmetry axes of a planar point set in O(n log n)). The vertical-axis
  specialisation is proved here (L15 pattern): a reflection in x = c maps the set onto itself and so maps
  its least x to its greatest, hence 2c = minX + maxX; the set is symmetric iff every point's mirror
  (minX + maxX − x, y) is in the set.
- **Best known**: O(n) expected with hashing; Ω(n log n) deterministic in the algebraic decision-tree model
  (it decides set equality). This: O(n) expected time, O(n) extra (one map).
  **Lost**: sort + binary search (O(n log n), deterministic); Atallah's general method (all axes; more work).
- **Points as a set**: duplicates are ignored. Empty → true. Mirror computed without overflow.
- **Worked example**: hand: (1,1),(−1,1) → true (axis x = 0); (1,1),(−1,−1) → false; (0,0),(3,0),(1,5),(2,5)
  → true (axis x = 1.5).
- **Tags**: symmetry, reflection, line-reflection. LeetCode: 356 (the problem's name: LineReflection).

### PerfectRectangle — HELD (no source)
The area-plus-corner-parity test (an exact tiling iff total area equals the bounding box's and every corner
point except the four outer ones occurs an even number of times) was found only on problem-site and
tutorial pages (searched 2026-11-12); no textbook or paper stating and proving it. L22/L31: not built.
Candidate for Devin: a sweep over x with an interval set (a Klee-style measure sweep that also detects
overlap) is sourceable from Klee 1977 / Bentley's sweep if wanted.

---

# algo/ratelimit (SHORT = RATELM)

Builder W4-B14. Not yet in `doc/taxonomy.md` (lead to add: `ALG | algo/ratelimit | RATELM | rate meters with
the clock passed in: token bucket, GCRA leaky bucket, fixed (jumping) window, sliding (moving) window log`).

| ID | Name | Signature |
|---|---|---|
| ALG-RATELM-001 | TokenBucket | `func NewTokenBucket(tokensPerInterval int64, interval time.Duration, burst int64) (*TokenBucket, error)` with `Allow(now time.Time, n int64) (ok bool, retryAfter time.Duration)`, `Available(now time.Time) int64` |
| ALG-RATELM-002 | LeakyBucket | `func NewLeakyBucket(emissionInterval, tolerance time.Duration) (*LeakyBucket, error)` with `Allow(now time.Time) (ok bool, retryAfter time.Duration)` |
| ALG-RATELM-003 | FixedWindow | `func NewFixedWindow(limit int64, window time.Duration) (*FixedWindow, error)` with `Allow(now time.Time, n int64) (ok bool, retryAfter time.Duration)` |
| ALG-RATELM-004 | SlidingWindowLog | `func NewSlidingWindowLog(limit int64, window time.Duration) (*SlidingWindowLog, error)` with `Allow(now time.Time, n int64) (ok bool, retryAfter time.Duration)` |
| — | SlidingWindowCounter | **HELD (no source)**: the "weighted previous window" estimate is described only by vendor blogs and library docs (searched 2026-11-12); no paper, RFC or standard found. No ID given. |

## Conventions shared by the package (builder decisions, for the lead's table)

- **Clock as a parameter.** Every call takes `now time.Time`; no type reads the clock, starts a goroutine or
  holds a lock. Not safe for concurrent use: the caller's mutex guards a meter (one meter per key, the
  caller's map). Times should come from one clock (`time.Now` readings carry the monotonic clock, which
  `Sub`/`Before` use when both operands have it).
- **Clock going backwards** (decision R1): a `now` earlier than the latest `now` a meter has seen is treated
  as that latest time: no elapsed time, no refill, no window change; the call is still decided. **Lost**:
  returning an error (a wall clock steps back under NTP; the caller can do nothing better than "no time
  passed", and an error on the hot path would be ignored or turned into an outage).
- **Time arithmetic**: only `Time.Sub` (saturates at ±292 years) and `Time.Add` of a non-negative duration
  never larger than an already-measured gap, or of the emission interval with a check that the result moved
  forward. Scaled integer products are checked with `math/bits.Mul64` at construction.
- **Requests that can never conform** (`n` greater than the burst or the limit): `(false, Never)` where
  `Never = time.Duration(math.MaxInt64)`. `n <= 0` is `(true, 0)` and changes nothing. A zero-value meter
  refuses everything with `Never` (its burst/limit is 0). **Lost**: a third error result (the signature the
  lead asked for is `(ok, retryAfter)`; the two documented values cover both cases without one).
- **Size cap** (L13/L26): SlidingWindowLog keeps at most `limit` entries of 32 bytes, allocated once; limit is
  capped at 1<<22 (128 MiB) and larger returns `ErrTooLarge`.

### ALG-RATELM-001 TokenBucket
- **Source**: Heinanen and Guérin, "A Single Rate Three Color Marker", RFC 2697, IETF, September 1999, §3
  "Metering" (read 2026-11-12, https://www.rfc-editor.org/rfc/rfc2697.txt): the committed bucket C starts
  full, Tc(0) = CBS; Tc is incremented CIR times per second up to CBS; a packet of size B conforms (green) when
  Tc − B ≥ 0 and then Tc is decremented by B, otherwise nothing is decremented. Only the C bucket is used
  (no excess bucket E: EBS = 0 makes the yellow branch empty). Turner, "New directions in communications (or
  which way to the information age?)", IEEE Communications Magazine 24(10):8–15, 1986,
  doi:10.1109/MCOM.1986.1092946 is the historical origin of the term (Crossref metadata confirmed; full text not
  read; not used for any rule).
- **Best known**: O(1) time and O(1) space per call; this is the lower bound (each call must read the state).
- **Approach**: exact integer fluid model. One token = `interval` nanoseconds of "credit units"; each elapsed
  nanosecond adds `tokensPerInterval` units, so the level is exact at every nanosecond and never drifts
  (capacity = burst·interval units, checked with `bits.Mul64`). Refill is clamped to the capacity before
  multiplying, so `elapsed·tokensPerInterval` never overflows. `retryAfter` = ⌈deficit / tokensPerInterval⌉ ns.
  This credits whole tokens at the same average rate as RFC 2697's "CIR times per second", in finer steps.
  **Lost**: float64 tokens (as in golang.org/x/time/rate; rounding accumulates and equality tests become
  tolerances); RFC 2697's discrete one-token ticks (needs a tick phase the caller would have to choose).
- **Equivalence used as an oracle**: with `tokensPerInterval = 1`, `interval = T`, `burst = b`, TokenBucket
  makes the same decisions and returns the same retryAfter as LeakyBucket(T, (b−1)·T), by I.371 Annex A.2:
  tokens = (b·T − max(X′, 0))/T where X′ is the continuous-state bucket content. Proof in the test comment.
- **Worked example (hand-derived from RFC 2697 §3)**: 2 tokens per second, burst 3, at t=0: Allow(n=3) ok;
  Allow(1) refused, retryAfter 500 ms; at t=500 ms Allow(1) ok; at t=10 s Available = 3 (capped at burst).
- **Errors**: `ErrRate` (tokensPerInterval < 1), `ErrInterval` (interval ≤ 0), `ErrBurst` (burst < 1),
  `ErrTooLarge` (burst·interval > MaxInt64). Input never modified.
- **Consumer evidence**: red_lite_cache `internal/adminauth/limiter.go:52-89` (a map of
  golang.org/x/time/rate limiters with a sweep and a cap) and the consumer backend
  `api/middleware/rate_limit.go:273-288` (`checkIPLimit`: `rate.NewLimiter(RequestsPerSecond, BurstSize)` in a
  `sync.Map`, `Allow()`). **What TokenBucket replaces**: the per-key bucket arithmetic only, i.e. the
  `*rate.Limiter` value and its `Allow`/`ReserveN` call (limiter.go `bucket.l`; rate_limit.go:279/287).
  Their call becomes `ok, retry := b.Allow(now, 1)` and `retry` feeds Retry-After directly (limiter.go
  computes it from a Reservation). **What stays theirs**: the key → bucket map, the idle sweep and
  MaxEntries cap, the clock (`now func() time.Time` / `time.Now()`), and the locking. Note for the consumer backend:
  `rate.Limiter` locks internally and is used from a `sync.Map` without an outer lock; TokenBucket does not
  lock, so the entry needs its own `sync.Mutex` around `Allow` (or a striped lock). Removes the
  golang.org/x/time/rate dependency from both. Consumers are listed, never edited.
- **Tags**: rate-limiting, token-bucket, traffic-policing, throttling, retry-after. Pattern pointers only:
  code-heim go_82 (not a source).

### ALG-RATELM-002 LeakyBucket
- **Source**: ITU-T Recommendation I.371 (03/2004), "Traffic control and congestion control in B-ISDN",
  Annex A "Generic Cell Rate Algorithm GCRA(T,τ)", §A.1 Virtual Scheduling Algorithm and §A.2
  Continuous-state leaky bucket algorithm, pp. 86–87 (read 2026-11-12 from the ITU free PDF): TAT is
  initialised to ta(1); a cell with ta(k) < TAT − τ is non-conforming and TAT is unchanged; if
  TAT − τ ≤ ta(k) < TAT, TAT += T; if ta(k) > TAT, TAT = ta(k) + T. The NOTE states A.1 and A.2 decide the
  same cells. Second source: Berger, Eckberg, Hou, Lucantoni, "Performance characterizations of traffic
  monitoring, and associated control, mechanisms for broadband 'packet' networks", Proc. IEEE GLOBECOM '90,
  pp. 350–354, 1990, doi:10.1109/GLOCOM.1990.116535, §4 p. 351 (the leaky-bucket monitor: a counter
  incremented per arrival and drained at a constant rate, excess over a threshold non-conforming; read).
- **Best known**: O(1) time, O(1) space per call (lower bound).
- **Approach**: the virtual scheduling form (A.1): one `time.Time` (TAT) of state, comparisons by
  `TAT.Sub(now)` (saturating). A cell conforms iff TAT − now ≤ τ; retryAfter = (TAT − now) − τ. The next TAT
  is checked to have moved forward (Time.Add saturates at the end of time; then `(false, Never)`).
  **Lost**: the continuous-state form (A.2) with X and LCT (two fields for the same decisions; kept as the
  test oracle); a queue-based "leaky bucket as a queue" (that shapes/delays traffic, needs a scheduler or
  goroutine; this package meters).
- **Weight**: one cell per call, as the source defines it. A caller wanting n units per request uses
  TokenBucket.
- **Worked example (hand-derived from I.371 A.1)**: T = 100 ms, τ = 200 ms, all at t=0: three calls conform
  (TAT 100, 200, 300 ms), the fourth is refused with retryAfter 100 ms; at t=100 ms one conforms (TAT 400 ms).
- **Errors**: `ErrInterval` (emissionInterval ≤ 0), `ErrTolerance` (tolerance < 0).
- **Tags**: rate-limiting, leaky-bucket, gcra, virtual-scheduling, traffic-policing, atm.

### ALG-RATELM-003 FixedWindow
- **Source**: Berger, Eckberg, Hou, Lucantoni, GLOBECOM '90 (above), §4 p. 351, the jumping window scheme:
  "the counter starts at zero, at the beginning of a fixed time interval (window). The counter is incremented
  upon cell arrivals within the window, and resets to zero at the end of the window. The next window starts
  immediately after the preceding window" (read). Rathgeb, "Modeling and performance comparison of policing
  mechanisms for ATM networks", IEEE J. Sel. Areas Commun. 9(3):325–334, 1991, doi:10.1109/49.76630
  compares the same jumping and moving windows (Crossref and abstract confirmed; full text not read).
- **Best known**: O(1) time, O(1) space per call.
- **Approach**: the first call opens window 0 at its `now`; windows are [start + jW, start + (j+1)W). On a
  call past the window, start jumps by ⌊elapsed/W⌋·W (never past now, so no overflow) and the count resets.
  A request conforms when count + n ≤ limit (checked as n ≤ limit − count). retryAfter is the time to the
  window's end. **Lost**: windows aligned to the Unix epoch (`now.Truncate(W)`; useful across machines but
  the source anchors at a window beginning, and the caller can align by passing aligned times); the triggered
  jumping window (Rathgeb's TJW: window starts at the first arrival after the last expired) — not requested.
- **Known property (tested)**: up to 2·limit can pass in a span shorter than W across a window edge; the
  sliding log (004) does not allow this.
- **Worked example (hand-derived)**: limit 3, W = 1 s: at 0, 0.2, 0.9 s ok; at 0.95 s refused, retryAfter
  50 ms; at 1.0 s ok (new window).
- **Errors**: `ErrLimit` (limit < 1), `ErrWindow` (window ≤ 0).
- **Tags**: rate-limiting, fixed-window, jumping-window, counter, traffic-policing.

### ALG-RATELM-004 SlidingWindowLog
- **Source**: Berger, Eckberg, Hou, Lucantoni, GLOBECOM '90, §4 p. 351, the moving window scheme: "the maximum
  number of cell arrivals within a window of fixed time interval (W time slots) is limited as in the jumping
  window scheme. The difference is that adjacent windows overlap by W−1 slots … it requires that the arrival
  time of a cell be stored for W slots" (read). Rathgeb 1991 (above) as the second, metadata-confirmed source.
- **Best known**: O(1) amortized time per accepted call (each entry is pushed and evicted once), O(limit)
  space: the log must remember every accepted arrival still inside the window, and up to `limit` of them
  can be (lower bound for an exact moving window). A refused call costs O(min(n, limit)) to compute
  retryAfter (O(1) for n = 1).
- **Approach**: a ring buffer of (time, n) entries, `queue.CircularBuffer` (ALG-QUEUE-002, reused by call per
  L11), capacity `limit` (each entry holds ≥ 1 unit, so it never fills), and a running total. Entry t counts
  for now ∈ [t, t+W). retryAfter = the time until enough of the oldest entries expire for n to fit.
  **Lost**: a slice of timestamps re-sliced from the front (leaks the head until growth copies);
  per-unit entries (n entries per request; O(n) per call).
- **Worked example (hand-derived from the moving-window rule)**: limit 3, W = 1 s: at 0, 0.2, 0.9 s ok; at
  0.95 s refused, retryAfter 50 ms (the 0 s entry expires at 1.0 s); at 1.0 s ok (the 0 s entry has left);
  at 1.1 s refused, retryAfter 100 ms (the 0.2 s entry expires at 1.2 s). Same arrivals as FixedWindow's
  example, so the two tests show where the schemes differ.
- **Errors**: `ErrLimit`, `ErrWindow`, `ErrTooLarge` (limit > 1<<22).
- **Tags**: rate-limiting, sliding-window, moving-window, sliding-log, traffic-policing. Pattern pointers
  only: code-heim go_69/go_77 (not sources).

### SlidingWindowCounter — HELD (no source)
The estimate count ≈ previous·(1 − elapsed/W) + current is described by vendor engineering blogs and library
documentation only (search 2026-11-12: limits.readthedocs.io, Kong, Vercel, Appwrite, others); no paper, RFC,
ITU recommendation or textbook found. Berger et al.'s "exponential smoothing scheme" (§4) is a different,
cited mechanism and could be offered instead if Devin wants a fixed-memory sliding estimate.

---

# algo/bloom (SHORT = BLOOM)

Builder W4-B14. Not yet in `doc/taxonomy.md` (lead to add: `ALG | algo/bloom | BLOOM | Bloom filters:
approximate set membership with false positives, parameter formulas`).

| ID | Name | Signature |
|---|---|---|
| ALG-BLOOM-001 | BloomFilter | `func NewBloomFilter(bitCount uint64, hashCount int) (*BloomFilter, error)` with `Add(item []byte)`, `MayContain(item []byte) bool`, `EstimatedFPR() float64`, `BitCount() uint64`, `HashCount() int`, `SetBits() uint64` |
| ALG-BLOOM-002 | FalsePositiveRate | `func FalsePositiveRate(bitCount, items uint64, hashCount int) (float64, error)` |
| ALG-BLOOM-003 | OptimalParams | `func OptimalParams(expectedItems uint64, falsePositiveRate float64) (bitCount uint64, hashCount int, err error)` |
| ALG-BLOOM-004 | BloomFilter.UnionWith | `func (f *BloomFilter) UnionWith(other *BloomFilter) error` |
| — | Intersection | **Not built**: Broder–Mitzenmacher §2.4 uses two filters only to *estimate* the size of an intersection; a bitwise-AND "intersection filter" is not in the sources read. Optional in the brief. |

## Decisions (for the lead's table)

- **B1 Hash function: SHA-256 (`crypto/sha256.Sum256`, FIPS 180-4)**, h₁ = digest bytes 0–7, h₂ = bytes 8–15
  (big-endian). Deterministic across processes, machines and Go versions, so a filter can be persisted,
  exchanged and unioned between processes and tests pin exact bit positions; its output is close to the
  "fully random" hashes the analysis assumes; `Sum256` keeps no state, so MayContain only reads and
  concurrent MayContain calls are safe (Add and UnionWith still need the caller's lock); no allocation.
  **Lost (measured)**: FNV-1a-128 (`hash/fnv`), the first choice: its two 64-bit halves are correlated
  for short keys (the prime is 2⁸⁸ + 0x13B), and a Python model of this exact construction gave a measured
  FPR of 0.00093 and a fill-based estimate of 0.0055 against the predicted 0.0100 (m = 9586, k = 7,
  1000 keys "member-i", 100 000 probes "other-i"); with SHA-256 the same run measured 0.01004 (scratch
  script, 2026-11-12). **Lost**: `hash/maphash` with a caller's `maphash.Seed` (faster and keyed, but a Seed
  cannot be serialised, so a filter means nothing outside its process, and tests could only check
  statistics). Cost: SHA-256 is slower per byte than FNV or maphash (hardware-accelerated on amd64/arm64);
  the benchmark records it. The hash is public, so an adversary choosing items can search offline for
  ones that hit set bits; a keyed variant would be a new function.
- **B2 Index arithmetic**: a = h₁ mod m, b = h₂ mod m, gᵢ = (a + i·b) mod m for i = 0…k−1 computed by
  repeated addition with one conditional subtraction (exact, no overflow since m ≤ 2³⁰). b = 0 (probability
  1/m) makes all k indices equal; the source's analysis includes that case (h₂ uniform on [m]).
- **B3 Size cap** (L13/L26): bitCount ≤ `MaxBits` = 1<<30 bits (128 MiB), exported, else `ErrTooLarge`; OptimalParams
  returns `ErrTooLarge` when the m it computes exceeds the cap.
- **B4 Items are `[]byte`**. A string caller converts (`[]byte(s)`). **Lost**: a generic `T comparable` API
  (needs maphash.Comparable, see B1).

### ALG-BLOOM-001 BloomFilter
- **Source**: Bloom, "Space/time trade-offs in hash coding with allowable errors", Communications of the ACM
  13(7):422–426, 1970, doi:10.1145/362686.362692 (Crossref metadata confirmed 2026-11-12; full text not read: the
  origin, credited by both papers below). The statement read and implemented: Kirsch and Mitzenmacher, "Less hashing, same performance: building a better Bloom filter",
  Random Structures & Algorithms 33(2):187–218, 2008, doi:10.1002/rsa.20208, §2 "Standard Bloom filters"
  (read: m bits initially 0, k hash functions, set bits hᵢ(x), test whether all hᵢ(y) are 1; given the
  fraction ρ of 0 bits, the false-positive probability is (1 − ρ)ᵏ) and §5.2, Theorem 5.2 (read: double
  hashing, location h₁(u) + i·h₂(u) mod m for i ∈ [k] with h₁, h₂ independent uniform on [m], has the same
  asymptotic false-positive probability (1 − e^{−k/c})ᵏ as k independent hashes, c = m/n).
- **Best known**: Add and MayContain O(L + k) time (L = item length in bytes, hashed once; then k index steps),
  O(1) extra space; O(m) bits total. Space lower bound for any set representation with false-positive rate ε:
  n·log₂(1/ε) bits; a Bloom filter with optimal k uses log₂e ≈ 1.44 times that (Broder–Mitzenmacher §2.2).
- **Approach**: bits packed in `[]uint64`; one SHA-256 digest per call, k indices by double hashing (B2).
  A running count of set bits (updated when Add flips a 0) gives EstimatedFPR = (setBits/m)ᵏ in O(1)
  (KM §2: conditioned on ρ, FPR = (1−ρ)ᵏ), which, unlike a count of Adds, ignores duplicates and stays right
  after UnionWith. **Lost**: k independent hash functions (k digests per call; KM Theorem 5.2 shows no loss
  in the asymptotic FPR); the partitioned variant (KM §5.1; same asymptotics, harder to union with an
  unpartitioned filter); a counter of Adds for the FPR (overcounts duplicates).
- **Worked example (computed independently with Python hashlib, 2026-11-12; the test also recomputes gᵢ with
  math/big)**: m = 64, k = 3: "apple" → bits 41, 48, 55; "banana" → bits 13, 35, 57. After both Adds,
  SetBits = 6, both MayContain, EstimatedFPR = (6/64)³.
- **Errors**: `ErrBitCount` (bitCount < 1), `ErrHashCount` (hashCount < 1), `ErrTooLarge`
  (bitCount > 1<<30). Items are read, never modified or retained.
- **Tags**: bloom-filter, probabilistic, set-membership, approximate, double-hashing, sketch.

### ALG-BLOOM-002 FalsePositiveRate
- **Source**: Broder and Mitzenmacher, "Network applications of Bloom filters: a survey", Internet
  Mathematics 1(4):485–509, 2004, doi:10.1080/15427951.2004.10129096, §2.1 (read 2026-11-12): after n
  items, a given bit is still 0 with probability p′ = (1 − 1/m)^{kn} ≈ e^{−kn/m}; FPR f = (1 − e^{−kn/m})ᵏ.
  KM 2008 §2 states the same.
- **Best known**: O(1) time and space (one exp, one pow).
- **Approach**: `-math.Expm1(-k·n/m)` for 1 − e^{−kn/m} (accurate when kn/m is small), then `math.Pow(·, k)`.
  **Lost**: the exact (1 − (1 − 1/m)^{kn})ᵏ (differs by O(1/m); the brief and both sources use the
  exponential form).
- **Worked example (from survey §2.1)**: k = ln 2 · m/n gives f = (1/2)ᵏ; m/n = 10, k = 7 gives
  (1 − e^{−0.7})⁷ ≈ 0.008194 (hand-computed, checked against the closed form within 1e-12 relative).
- **Errors**: `ErrBitCount` (bitCount < 1), `ErrHashCount` (hashCount < 1). items = 0 gives 0.
- **Tags**: bloom-filter, false-positive-rate, probability.

### ALG-BLOOM-003 OptimalParams
- **Source**: Broder–Mitzenmacher 2004, §2.1 (k = ln 2 · (m/n) minimises f) and §2.2 (f ≤ ε requires
  m ≥ n·log₂(1/ε)/ln 2 = −n·ln ε/(ln 2)²) (read).
- **Best known**: O(1).
- **Approach**: m = ⌈−n·ln p/(ln 2)²⌉, k = max(1, round(m/n · ln 2)). Rounding k to the nearest integer is the
  survey's "in practice k must be an integer" (§2.1); the test checks the returned (m, k) meet p by
  FalsePositiveRate within a 1% relative margin (rounding k moves f by less than that for m ≥ 64).
  **Lost**: floor(k) (always cheaper but can exceed p by more).
- **Worked example**: n = 1000, p = 0.01 → m = ⌈9585.06⌉ = 9586, k = round(6.6446) = 7 (hand-computed).
- **Errors**: `ErrItems` (expectedItems < 1), `ErrProbability` (p not in (0, 1), or NaN), `ErrTooLarge`
  (m > 1<<30).
- **Tags**: bloom-filter, sizing, optimal-parameters.

### ALG-BLOOM-004 BloomFilter.UnionWith
- **Source**: Broder–Mitzenmacher 2004, §2.4 "Standard Bloom filter tricks" (read): two filters with the same
  number of bits and the same hash functions represent the union of their sets by the OR of their bit vectors.
- **Best known**: O(m/64) time, O(1) extra space (every word must be read).
- **Approach**: word-wise OR into the receiver, then recount set bits with `math/bits.OnesCount64`.
- **Errors**: `ErrIncompatible` when other is nil or has a different bitCount or hashCount (the hash
  function is fixed package-wide, so equal (m, k) is compatible). The receiver is modified (method on the
  filter it updates); other is read only.
- **Tags**: bloom-filter, union, set-operations, merge.

---

# algo/seq (SHORT = SEQ)

Builder W4-B15. New package (lead decision L3: `algo/seq`, not `algo/iter`). Not yet in `doc/taxonomy.md`
(lead to add: `ALG | algo/seq | SEQ | range-over-func helpers over iter.Seq / iter.Seq2`).

**Primary sources for every function** (confirmed 2026-11-12):
- The Go Programming Language Specification (Go 1.23+), "For statements with range clause", range over a
  function `func(yield func(K) bool)`: the loop body runs once per yield call; when the body terminates
  (break, return, panic) yield returns false and must not be called again. https://go.dev/ref/spec#For_range.
  (The fetched page was truncated before that section; the heading and rule are cited from the spec as
  published since Go 1.23 — the reviewer should open the anchor once.)
- Go standard library, package `iter` documentation, https://pkg.go.dev/iter: "Single-Use Iterators", "Pulling
  Values" — `Pull` returns (next, stop); stop "must be called when the caller is no longer interested in next
  values and next has not yet signaled that the sequence is over"; calling stop more than once, or after the
  end, is allowed (read 2026-11-12).

**Common shape.** Every helper is a pure function returning an iterator; nothing is read ahead; the helper's own
`for v := range upstream` returns as soon as the consumer's yield returns false, so the upstream's yield returns
false at once and it is never asked for another element (tested per helper with a probe upstream that records
produced / stopped / finished). A nil iterator is the empty sequence. A bare count or size returns
`(iter.Seq, error)` (run convention, L32). Lower bound for all: Ω(m) to read m elements; every helper is O(1)
(amortized) per element.

| ID | Name | Signature |
|---|---|---|
| ALG-SEQ-001 | Map | `func Map[T, U any](seq iter.Seq[T], fn func(T) U) iter.Seq[U]` |
| ALG-SEQ-002 | Filter | `func Filter[T any](seq iter.Seq[T], keep func(T) bool) iter.Seq[T]` |
| ALG-SEQ-003 | Take | `func Take[T any](seq iter.Seq[T], n int) (iter.Seq[T], error)` |
| ALG-SEQ-004 | Chunk | `func Chunk[T any](seq iter.Seq[T], size int) (iter.Seq[[]T], error)` |
| ALG-SEQ-005 | Zip | `func Zip[A, B any](first iter.Seq[A], second iter.Seq[B]) iter.Seq2[A, B]` |
| ALG-SEQ-006 | Flatten | `func Flatten[T any](seqs iter.Seq[iter.Seq[T]]) iter.Seq[T]` |
| ALG-SEQ-007 | Batch | `func Batch[T any](seq iter.Seq[T], maxSize int, endsBatch func(T) bool) (iter.Seq[[]T], error)` |
| ALG-SEQ-008 | Window | `func Window[T any](seq iter.Seq[T], width int) (iter.Seq[[]T], error)` |

### ALG-SEQ-001 Map / ALG-SEQ-002 Filter
- **Best known**: O(1) per element + the cost of fn/keep; O(1) extra. Met.
- **Approach**: range over seq, yield fn(v) / v when keep(v). **Lost**: a version taking iter.Seq2 too (one
  concept per function; a Seq2 variant can be added with its own ID if a consumer asks).
- **Errors**: none; nil seq or nil fn/keep is the empty sequence (documented zero value, never a panic).
- **Worked example**: Map(1 2 3, square) → 1 4 9; Filter(1…6, even) → 2 4 6. Hand-derived.
- **Allocation**: one closure per call, one per range loop (the loop body closure escapes into seq); none per
  element (tested: allocations over 10 and 10 000 elements are equal).

### ALG-SEQ-003 Take
- **Best known**: O(min(n, m)). Met: after the n-th element Take returns from its range loop, so the upstream
  is stopped at its n-th yield and never produces element n+1 (tested on an infinite probe).
- **Errors**: ErrCount (n < 0). n = 0 never starts the upstream.
- **Worked example**: Take(0 1 2 …, 3) → 0 1 2, upstream produced exactly 3. Hand-derived.

### ALG-SEQ-004 Chunk
- **Best known**: O(1) amortized per element, O(size) extra. Met (O(min(size, m)): the buffer grows by append,
  so a huge size never allocates up front or panics in make).
- **Decision (buffer)**: ONE reused slice per range loop, rewritten between groups (valid until the body
  returns; slices.Clone to keep). **Lost**: a fresh slice per group (one allocation per group; callers that
  only stream the groups pay for nothing they use). Reason: run convention L32 (enumeration iterators yield
  one reused buffer). Documented in doc.go and the doc comment; a test checks the backing array is reused.
- **Errors**: ErrSize (size < 1). Last group holds the remaining 1…size elements.
- **Worked example**: Chunk(1…7, 3) → [1 2 3] [4 5 6] [7]. Hand-derived. Oracle: re-slicing the input.
- **Not slices.Chunk**: the standard library's slices.Chunk splits a slice; this groups an iterator.

### ALG-SEQ-005 Zip
- **Best known**: O(1) per pair. Met.
- **Approach**: range over first; pull second with `iter.Pull`, created lazily at first's first element
  (second is not started when first is empty); `stop` is deferred and therefore ALWAYS called: on consumer
  break, when first ends, when second ends, and when the loop body panics (each tested). first is stopped by
  returning from its range loop. iter.Pull runs second as a coroutine (no goroutine of ours, no channel).
  **Lost**: pulling both with two iter.Pull (two coroutine switches per pair instead of one); collecting
  either side into a slice (O(m) memory, breaks infinite iterators).
- **Errors**: none; nil on either side is the empty sequence.
- **Worked example**: Zip(1 2 3, "x" "y") → (1,x) (2,y). Hand-derived.

### ALG-SEQ-006 Flatten
- **Best known**: O(1) per element and per inner iterator. Met.
- **Approach**: nested range loops; on stop, both the current inner and the outer iterator are stopped and
  no later inner is started (tested). nil inner iterators are skipped.
- **Errors**: none. **Worked example**: [1 2] [] [3] → 1 2 3. Hand-derived.

### ALG-SEQ-007 Batch
- **Decision**: KEPT, not dropped: Batch is NOT identical to Chunk. A batch closes at maxSize elements OR when
  `endsBatch(v)` is true for the element just added (that element is the batch's last). With a nil
  endsBatch it groups exactly as Chunk (tested). Use: records up to a terminator, messages up to a commit
  marker, with a size cap. **Lost**: dropping it (the predicate flush has no other home); a time-based flush
  (needs a clock or a timer goroutine: not a pure iterator).
- **Best known**: O(1) amortized per element; O(maxSize) extra. Met. Same reused-buffer rule as Chunk.
- **Errors**: ErrSize (maxSize < 1).
- **Worked example**: 5 0 1 2 3 4 0 6, terminator 0, maxSize 4 → [5 0] [1 2 3 4] [0] [6]. Hand-derived.

### ALG-SEQ-008 Window
- **Source (ring)**: Knuth, TAOCP Vol. 1, 3rd ed., 1997, §2.2.2 (circular array, indices modulo M).
- **Best known**: O(1) per element, O(width) extra, and no copy of the window per step. Met.
- **Approach**: a ring of 2·width slots; element i is written to slot i mod width AND slot i mod width +
  width, so the window starting at slot s is the contiguous slice [s, s+width): two writes per step, never a
  copy. The ring is allocated once the first window is full (before that, the elements are collected by
  append), so a huge width on a short seq costs nothing and never panics in make.
  **Lost**: a ring of width slots copied into a separate output slice per window (O(width) per step);
  `queue.CircularBuffer` (its elements are not contiguous, so each window would need a copy).
- **Errors**: ErrSize (width < 1). m < width gives no windows.
- **Worked example**: Window(1…5, 3) → [1 2 3] [2 3 4] [3 4 5]. Hand-derived. Oracle: in[i:i+width].

## Channel patterns: NOT library material (they own goroutines)
calcbench functions start no goroutine the caller did not ask for (CLAUDE.md, "What every function is"). The
code-heim channel exercises are listed here only to record why they are not built:
- go_61 errgroup (run tasks, cancel on first error): owns goroutines and a context; that is
  golang.org/x/sync/errgroup's job, not a pure function.
- go_62 or-done (wrap a channel so a done signal ends the read): a goroutine forwarding between channels.
- go_62 tee (copy one channel to two): a goroutine per tee with unbounded blocking coupling.
- go_56 / go_60 generators (a goroutine feeding a channel): replaced by iter.Seq, which needs no goroutine.
- fan-out / fan-in (several workers, merged results): goroutines and synchronisation by definition.

## Tags (problem sites, never in identifiers)
Take/Chunk/Window/Zip/Flatten: iterator-pattern exercises (code-heim go_56/go_60 generators reworked as
iter.Seq); Window: sliding-window pattern (LeetCode 239/643 use the same windows).

---

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
