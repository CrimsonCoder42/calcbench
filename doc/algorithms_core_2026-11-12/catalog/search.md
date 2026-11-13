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
  Science of Programming*, Springer, 1981 (saddleback search, p. 215 according to a search summary
  2026-11-12; no section number found, and neither page nor section checked against a copy); Bird, "Improving saddleback search: a lesson in algorithm design",
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
