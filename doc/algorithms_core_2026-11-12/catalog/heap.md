# algo/heap (SHORT = HEAP) — S2 additions (builder W1-B4)

ALG-HEAP-001 TopK is S1 (shipped); its sift is unexported and stays as is. `Heap[T]` below is the
reusable exported type that wave 2 imports (graph Dijkstra/Prim, list LFU, window MovingMedian).

## Heap API (for wave 2)

```go
var ErrNilLess = errors.New("heap: less function is nil")
var ErrTooLarge = errors.New("heap: requested size is too large") // NewHeap, NewMedianStream (rule L41)
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
- **Errors**: `ErrNilLess`; `ErrTooLarge` (NewHeap: capacityHint·sizeof(T) > 2⁴⁷ bytes, 2³¹ on 32-bit; rule L41). Empty Pop/Peek → (zero, false). Zero/nil Heap: Push false.
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
  enumeration and sort (O(nm log nm)). **Lost (fix M8)**: Frederickson, "An optimal algorithm for selection
  in a min-heap", Information and Computation 104(2):197–214, 1993, doi:10.1006/inco.1993.1030: the pair
  matrix is an implicit heap-ordered tree, and his method finds the k smallest in O(k) time, beating
  O(k log m). It lost because it returns the k pairs unordered (sorting them costs O(k log k) again), uses
  recursive grouping structures with large constants, and no consumer needs k large enough to repay that.
  It is the candidate for an unordered variant under a new ID.
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
`type MedianStream[T cmp.Ordered]` (zero value ready) ·
`NewMedianStream[T](capacityHint int) (*MedianStream[T], error)` (signature changed by fix H3) · `Add(v T) bool` · `Median() (lower, upper T, err error)` · `Len() int`.
- **Source**: Sedgewick, Wayne, *Algorithms*, 4th ed., Addison-Wesley, 2011, §2.4, Exercise 2.4.30
  "Dynamic median-finding": insert O(log n), find the median O(1), using a max-heap for keys below the median
  and a min-heap for keys above (booksite algs4.cs.princeton.edu/24pq, read 2026-11-12, lists it as Creative
  Problem 4 with the two-heap hint; the book number 2.4.30 confirmed by a search listing, fix M5). Lower/upper median: CLRS ch. 9 introduction.
- **Best known**: O(log n) per insert, O(1) median, O(n) space (the textbook bound above). No lower bound is
  claimed here (none cited).
- **Approach**: two heaps, lower half max-heap holds the extra value. **Lost**: an order-statistic tree
  (O(log n) both, more memory and code); re-selecting each time (O(n) per query); a sorted slice (O(n) insert).
- **Median shape**: (lower, upper, err) as ALG-SELECT-002, so callers average or not as their domain says.
- **Worked example**: KNUTH16 in order, medians after 1…5 keys: (503,503) (87,503) (503,503) (87,503)
  (503,503); after 16: (509,512), matching ALG-SELECT-002. Hand-derived.
- **Errors**: `ErrEmpty`; `ErrTooLarge` from NewMedianStream (rule L41). nil *MedianStream: Add false.
- **Tags**: heap, median, running-median, streaming, two-heaps. LeetCode: 295 (480 is window MovingMedian, wave 2).

## Not here
- IntervalSchedulingMaxNonOverlap was listed under heap in the request; it needs no heap (sort + greedy) and
  its one home is algo/interval (ALG-INTVL-008).
