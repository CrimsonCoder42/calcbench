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
