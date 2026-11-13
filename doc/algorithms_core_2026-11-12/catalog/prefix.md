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
