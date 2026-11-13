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
