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
- **Errors**: `ErrNegativeSize` (n < 0), `ErrTooLarge` (L41: the int parent array would exceed
  2⁴⁷ bytes, 2³¹ on 32-bit; test NewDisjointSet(math.MaxInt)), `ErrOutOfRange` (element outside 0..Len()−1). nil *DisjointSet:
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
