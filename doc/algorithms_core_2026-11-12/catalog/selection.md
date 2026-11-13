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
