# S1: Reference algorithms

Status: SHIPPED (2026-11-12)
Work folder: `doc/reference_algorithms_2026-11-12/`
Decisions: D-001 (introselect), D-002 (NaN: total order; reject in sums), D-003 (window dst), D-004 (Median shape).
Assumptions: A-001…A-011.

All functions are generic over `cmp.Ordered` except RollingSum/RollingMean, which do arithmetic and take
`float64`. No function modifies its input. Ordering is `cmp.Compare`/`cmp.Less` (D-002). Standard library only.
Values are unitless to the algorithm: a result is in the unit of its input values.

## Functions

| ID | Name | Signature |
|---|---|---|
| ALG-SEARCH-001 | LowerBound | `func LowerBound[T cmp.Ordered](sorted []T, target T) int` |
| ALG-SEARCH-002 | UpperBound | `func UpperBound[T cmp.Ordered](sorted []T, target T) int` |
| ALG-SEARCH-003 | BinarySearchFirst | `func BinarySearchFirst[T cmp.Ordered](sorted []T, target T) (index int, found bool)` |
| ALG-SEARCH-004 | BinarySearchLast | `func BinarySearchLast[T cmp.Ordered](sorted []T, target T) (index int, found bool)` |
| ALG-SELECT-001 | KthSmallest | `func KthSmallest[T cmp.Ordered](values []T, k int) (T, error)` |
| ALG-SELECT-002 | Median | `func Median[T cmp.Ordered](values []T) (lower, upper T, err error)` |
| ALG-WINDOW-001 | RollingSum | `func RollingSum(dst, values []float64, windowLen int) ([]float64, error)` |
| ALG-WINDOW-002 | RollingMean | `func RollingMean(dst, values []float64, windowLen int) ([]float64, error)` |
| ALG-WINDOW-003 | RollingMin | `func RollingMin[T cmp.Ordered](dst, values []T, windowLen int) ([]T, error)` |
| ALG-WINDOW-004 | RollingMax | `func RollingMax[T cmp.Ordered](dst, values []T, windowLen int) ([]T, error)` |
| ALG-HEAP-001 | TopK | `func TopK[T cmp.Ordered](values []T, k int) ([]T, error)` |
| ALG-MERGE-001 | MergeKSorted | `func MergeKSorted[T cmp.Ordered](lists [][]T) []T` |
| ALG-INTVL-001 | MergeOverlapping | `func MergeOverlapping[T cmp.Ordered](intervals []Interval[T]) ([]Interval[T], error)` |

`Interval[T]` is a closed range `{Start, End T}` (A-005), declared in `algo/interval`.

"KNUTH16" below means the 16 keys Knuth uses as a running example in TAOCP Vol. 3 (2nd ed.), ch. 5 and §6.2.1:
503 087 512 061 908 170 897 275 653 426 154 509 612 677 765 703. Sorted: 061 087 154 170 275 426 503 509
512 612 653 677 703 765 897 908.

## algo/search (SHORT = SEARCH)

### ALG-SEARCH-001 LowerBound · -002 UpperBound · -003 BinarySearchFirst · -004 BinarySearchLast
- **Source**: Knuth, *The Art of Computer Programming*, Vol. 3, 2nd ed., §6.2.1, Algorithm B (binary search).
  The leftmost-match form whose invariant brackets the first occurrence: Bentley, *Programming Pearls*, 2nd
  ed., Column 9.3.
- **Best known**: O(log n) time, O(1) space. Lower bound: ⌈log2(n+1)⌉ comparisons in the worst case for
  any comparison-based search of n sorted keys (Knuth §6.2.1, the comparison-tree argument). The half-open
  bisection does at most ⌊log2 n⌋+1 = ⌈log2(n+1)⌉ comparisons for n ≥ 1, so it meets the bound.
- **Approach**: half-open `[lo, hi)` bisection on `cmp.Less` with midpoint `lo + (hi-lo)/2`, and no equality
  test inside the loop. First and Last are LowerBound and UpperBound plus one equality check.
  **Lost**: Knuth's three-way Algorithm B, which stops at the first equal key it meets. That is not
  necessarily the first or last occurrence, and it costs two comparisons per step. Also lost: a linear scan
  (O(n)), and Eytzinger / branchless layouts (Khuong & Morin, ACM JEA 22, 2017). Those need a re-laid-out
  input, not a plain sorted slice.
- **Input modified**: no.
- **Worked example**: sorted KNUTH16, target 653 → `BinarySearchFirst` = (10, true), the 0-based form of
  Knuth's position 11. `LowerBound(400)` = 5 (426 is the first key ≥ 400). `UpperBound(653)` = 11. Hand-derived.
- **Errors**: none. Unsorted input is a precondition, not detected. The result is still in [0, n] and
  nothing panics. When the target is not found, First and Last return (insertion point, false).

## algo/selection (SHORT = SELECT)

### ALG-SELECT-001 KthSmallest
- **Source**: Hoare, "Algorithm 65: Find", CACM 4(7):321–322, 1961 (quickselect). Blum, Floyd, Pratt, Rivest,
  Tarjan, "Time bounds for selection", JCSS 7(4):448–461, 1973 (median of medians). Musser, "Introspective
  sorting and selection algorithms", Software: Practice and Experience 27(8):983–993, 1997 (the introspective
  switch). Textbook: Cormen et al., *Introduction to Algorithms*, 3rd ed., §9.2–9.3.
- **Best known**: O(n) worst-case time (BFPRT). The lower bound is Ω(n), since every element must be read.
  Space: O(1) extra when reordering in place. This function copies (A-002), so O(n) extra, one allocation,
  plus O(log n) stack for the median-of-medians recursion (depth ≤ log5 n + 1).
- **Approach** (D-001): introselect. Median-of-3 pivot, then a Dijkstra three-way partition, so equal keys
  end in one pass. A range of ≤ 16 is finished by insertion sort. Once the median-of-3 rounds have
  partitioned more than 6n elements in total (A-010), the pivot becomes the median of medians of groups
  of 5 for the rest of the call. That bounds the total at O(n) worst case: at most 7n before the switch,
  then a linear median-of-medians tail. (The first rule, "halve the range every two rounds", switched
  on 59–80% of random inputs in verify and was replaced.)
  **Lost**: pure median-of-medians, randomized quickselect (D-001), and sort-then-index (O(n log n)).
- **Input modified**: no (a copy is reordered).
- **Worked example**: KNUTH16 in Knuth's original order. k = 1 → 61, k = 5 → 275, k = 8 → 509, k = 16 → 908.
  Read off the sorted list by hand.
- **Errors**: `ErrEmpty` (len(values) = 0); `ErrRank` (k < 1 or k > n). k is 1-based (A-001).

### ALG-SELECT-002 Median
- **Source**: as ALG-SELECT-001. The median as the ⌊(n+1)/2⌋-th and ⌈(n+1)/2⌉-th order statistics:
  Cormen et al., 3rd ed., ch. 9 intro (lower and upper median).
- **Best known**: O(n) time; Ω(n) lower bound. Space as ALG-SELECT-001.
- **Approach**: one selection for the lower median at index (n−1)/2. For even n, the upper median is the
  minimum of the part to its right, which the partition guarantees is ≥ the lower median (one O(n) scan).
  **Lost**: two independent selections (twice the work); a full sort (O(n log n)).
- **Input modified**: no.
- **Worked example**: KNUTH16 → (509, 512). The first 15 keys of KNUTH16 (703 dropped) → (509, 509).
  Hand-derived.
- **Errors**: `ErrEmpty`.

## algo/window (SHORT = WINDOW)

For all four: windows are `values[i-w+1 : i+1]` for i = w−1 … n−1, so the output has n−w+1 elements and
each one includes the current value. Results are appended to dst (D-003). w > n → dst unchanged (A-004).

### ALG-WINDOW-001 RollingSum · ALG-WINDOW-002 RollingMean
- **Source**: the single moving average M_t = (X_t + … + X_{t−N+1}) / N: NIST/SEMATECH *e-Handbook of
  Statistical Methods*, §6.4.2.1 "Single Moving Average" (verified 2026-11-12). Compensated summation:
  Neumaier, ZAMM 54(1):39–51, 1974; Higham, *Accuracy and Stability of Numerical Algorithms*, 2nd ed., §4.3.
- **Best known**: O(n) time (Ω(n) to read the input). O(1) extra space. The result takes O(n−w+1) and is
  allocated once, only when dst lacks capacity.
- **Approach**: a Neumaier-compensated sliding sum (subtract the leaving value, add the entering one),
  restarted from scratch every w windows. A compensated sum that is never restarted still leaks: rounding
  left in the compensation term by large values that have left the window erased later tiny windows. The
  reviewer found this in verify (an error of 9e-4 at n = 1e6, and a wrong sign on a 4-window example).
  With the restart, a window's error depends only on the at most 2w−1 values since the restart. That
  costs at most 3n additions, still O(n). Values are checked for NaN and ±Inf in a first pass.
  Mean = sum / w.
  **Lost**: an uncompensated running sum (error grows with n); a compensated sum without restarts (above);
  recomputing each window (O(n·w)); prefix sums (O(n) extra space, and the cancellation error of
  subtracting two large prefixes).
- **Input modified**: no.
- **Worked example**: NIST §6.4.2.1 supplier data 9, 8, 9, 12, 9, 12, 11, 7, 13, 9, 11, 10 (1000-dollar
  units), N = 3. NIST prints the means 8.667, 9.667, 10.000, 11.000, 10.667, 10.000, 10.333, 9.667, 11.000,
  10.000. The exact sums 26, 29, 30, 33, 32, 30, 31, 29, 33, 30 are hand-derived; the test checks the means as
  exact fractions and also against NIST's printed values to ±0.0005.
- **Errors**: `ErrWindowLen` (w < 1); `ErrNotFinite` (a NaN or ±Inf value, with its index); `ErrOverflow` (a
  window's sum exceeds the float64 range, A-007). On error dst is returned with its length unchanged.

### ALG-WINDOW-003 RollingMin · ALG-WINDOW-004 RollingMax
- **Source**: Lemire, "Streaming maximum-minimum filter using no more than three comparisons per element",
  Nordic Journal of Computing 13(4):328–339, 2006 (the monotonic-deque filter).
- **Best known**: O(n) time, at most 2 comparisons per element amortized for a single filter (Lemire 2006);
  Ω(n) lower bound. O(w) extra space for the deque. The deque lives on the stack (no allocation) when
  w ≤ 64, and is one allocation otherwise.
- **Approach**: a ring-buffer deque of indices whose values are strictly monotone. Each index is pushed and
  popped at most once. **Lost**: van Herk / Gil–Werman blocks (O(n) and ≤ 3 comparisons per element, but
  two O(n) or O(w) prefix/suffix arrays, and no streaming); a heap or tree of the window (O(n log w));
  rescanning each window (O(n·w)).
- **Input modified**: no.
- **Worked example**: NIST §6.4.2.1 supplier data, w = 3. Min 8, 8, 9, 9, 9, 7, 7, 7, 9, 9; Max 9, 12, 12, 12,
  12, 12, 13, 13, 13, 11. Hand-derived from the source's data.
- **Errors**: `ErrWindowLen` (w < 1). NaN follows the total order (D-002): it is the minimum.

## algo/heap (SHORT = HEAP)

### ALG-HEAP-001 TopK
- **Source**: a bounded min-heap of the k largest seen so far: Sedgewick & Wayne, *Algorithms*, 4th ed.,
  §2.4, client TopM (verified 2026-11-12). Heaps: Williams, "Algorithm 232: Heapsort", CACM 7(6):347–348,
  1964; bottom-up heap construction: Floyd, "Algorithm 245: Treesort", CACM 7(12):701, 1964.
- **Best known**: O(n + k log k) time, which matches the lower bound (read n, emit k sorted). The method that
  reaches it is batched selection: a 2k buffer pruned by linear-time selection (CLRS §9.3), also in O(k) space.
  The bounded heap is O(n log k) worst case, O(k) space.
- **Approach**: the bounded heap (k slots, Floyd heapify of the first k values, then one comparison with the
  root per further value) and a final in-place heapsort that yields descending order. Space is O(1) beyond
  the result. **Lost**: batched selection. It has the better worst-case bound, but every element is copied
  into the buffer and partitioned. The heap rejects an element with a single comparison, which is the usual
  case once the heap holds large values. The heap is also the pattern this package is named for. It
  is kept on benchmark evidence (A-011; raw output in the work folder's `out/probe_topk_batch.txt`). Also lost: a full sort (O(n log n),
  O(n) space).
- **Input modified**: no.
- **Worked example**: KNUTH16, k = 3 → 908, 897, 765. Hand-derived.
- **Errors**: `ErrNegativeK` (k < 0). k > n returns all n values, descending (A-003).

## algo/merge (SHORT = MERGE)

### ALG-MERGE-001 MergeKSorted
- **Source**: k-way merge with a priority queue: Cormen et al., 3rd ed., §6.5 and Exercise 6.5-9; Knuth,
  TAOCP Vol. 3, 2nd ed., §5.4.1 (multiway merging).
- **Best known**: O(n log k) time. Lower bound Ω(n log k) comparisons: there are n!/∏(nᵢ!) interleavings,
  whose log is n log k − O(n) for equal lists. O(k) extra space; O(n) for the result.
- **Approach**: a binary min-heap of (head value, list index, position). On each step the root is replaced
  and sifted down; ties go to the lower list index, so the merge is stable. k = 1 is a copy.
  **Lost**: a tree of losers (Knuth §5.4.1). It is the same O(n log k) with about half the comparisons per
  element, but more state, and it is an optimisation to take only on benchmark evidence. Also lost:
  pairwise divide-and-conquer merging (same time, O(n) extra buffer), and concatenating then sorting
  (O(n log n)).
- **Input modified**: no.
- **Worked example**: KNUTH16 cut in Knuth's order into four runs of four, each sorted:
  [61 87 503 512], [170 275 897 908], [154 426 509 653], [612 677 703 765] → sorted KNUTH16.
- **Errors**: none. Unsorted lists are a precondition (A-009): every element still appears exactly once.

## algo/interval (SHORT = INTVL)

### ALG-INTVL-001 MergeOverlapping
- **Source**: sort by left endpoint, then one sweep: Klee, "Can the measure of ∪[aᵢ,bᵢ] be computed in less
  than O(n log n) steps?", American Mathematical Monthly 84(4):284–285, 1977. Lower bound: Fredman & Weide,
  "On the complexity of computing the measure of ∪[aᵢ,bᵢ]", CACM 21(7):540–544, 1978.
- **Best known**: O(n log n) time. Ω(n log n) in the algebraic decision-tree model (Fredman & Weide). O(n)
  for the result. O(log n) extra for the sort (pdqsort in `slices.SortFunc`).
- **Approach**: copy (that copy is the result buffer), sort by Start, merge in place in one sweep.
  **Lost**: an interval tree (CLRS ch. 14.3), which has the same bound plus O(n) nodes; a bitmap over an
  integer domain (O(n + range), not general).
- **Input modified**: no.
- **Worked example**: [9,12] [1,4] [3,5] [12,13] [20,25] [6,6] → [1,5] [6,6] [9,13] [20,25]. Written here;
  the result is hand-derived.
- **Errors**: `ErrInverted` (End < Start, with the index).

## Out of scope

In-place selection variants, a reusable heap or priority-queue type, interval intersection and coverage,
prefix sums, the arithmetic median (`num/stats`), and integer rolling sums (they need an overflow policy).
Each needs a new ID when it is requested.
