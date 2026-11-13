# Review fixes (lead's list; builders apply only the items in their own packages)

Common rule L41 (size guard): every function or constructor that allocates from a size argument returns
ErrTooLarge (a package sentinel, created if missing) when n·elemSize > maxAllocBytes, with a private
`const maxAllocBytes = 1 << (31 + 16*(^uint(0)>>63))` (2⁴⁷ bytes on 64-bit, 2³¹ on 32-bit: half the runtime's
maxAlloc, above which makeslice panics; runtime/malloc.go). Guards of the form n ≤ MaxInt/elemSize are not
enough. Packages that already apply the 128 MiB policy (L13/L35) keep it. Each new ErrTooLarge gets a test
(e.g. NewX(1<<60) returns ErrTooLarge, no panic) and the doc comment's Errors line lists it.

| # | Package (owner) | Fix |
|---|---|---|
| H1a | selection (B1) | kth_largest_test.go:14 expected 509 → 512 (8th largest of Knuth's 16 keys) |
| H1b | geometry2d (B13) | rectangles_union_area_test.go:34 expected 14 → 13 |
| H1c | matrix (B12) | sparse_multiply_test.go:26 denseMultiply must take p explicitly (m×0 · 0×p = m×p zeros) |
| H1d | prefix (B2) | summed_area_table_test.go:141 loop to table.Cols(), not the requested cols, when rows = 0 |
| H1e | enumerate (B12) | gray_code_test.go:96 assert the allocation count does not depend on bits (3 per pass measured), not ≤1 |
| H2 | intmath (B10) | min_squares_to_sum_test.go:52 cap a, b below 2²⁰ (the package took 245.7 s) |
| H3 | unionfind, graph (B8); heap (B4); queue (B3); bits (B10); sample, enumerate (B12) | L41 guards: NewDisjointSet; NewDirected/NewUndirected; FloydWarshall (V² table); NewHeap capacityHint; NewMedianStream; NewCircularBuffer; PopCountTable; ShuffleDeterministic; GenerateParentheses |
| H5 | selection, search, sorting (B1); stack (B3); dp (B9) | tests reusing problem-site example data → our own data with recomputed answers: kth_largest_test.go:25,28; search_sorted_matrix_test.go:28; search_rotated_dup_test.go:37; peak_element_test.go:28; dedupe_sorted_test.go:19; daily_temperatures_test.go:12; largest_rectangle_histogram_test.go:25; maximal_rectangle_test.go:49 and catalog/stack.md:72; dp intention→execution: cite Jurafsky–Martin, Speech and Language Processing, §2.5 (confirm) in the test comment and catalog/dp.md |
| M1 | graph (B8) | Dijkstra and BellmanFord: a relaxation that overflows positively is "no improvement" (skip); only report what the doc defines; test: NewDirected(2, {0→1 w1, 1→0 MaxInt64}) from 0 → [0, 1] |
| M2 | text (B11) | zero-value RollingHash: guard modulus 0 (error or no-op; never divide by zero); test it |
| M3 | tree (B6) | RangeAddMin/RangeAddMax: document that only the node extreme is overflow-checked (elements not aggregated can overflow and surface on a later point query), or keep both extremes per node; decide, document, test the documented behaviour |
| M5 | search (B1); heap (B4); stack (B3) | citations: Gries 1981 saddleback section number; Sedgewick–Wayne Exercise 2.4.30 for dynamic median; Naamad–Lee–Hsu DOI |
| M6 | text (B11) | drop the "(w−1)/q for a random base" collision claim from rolling_hash.go (not in Karp–Rabin or CLRS §32.2), or cite a polynomial-identity source that states it |
| M7 | bloom (B14) | step = h₂ mod m can be 0 or share a factor with m: use step = 1 + h₂ mod (m−1) (Dillinger–Manolios 2004 enhanced double hashing, doi:10.1007/978-3-540-30494-4_26 — confirm) or force an odd step with m a power of two; document; test step ≠ 0 |
| M8 | heap (B4); sample (B12); intmath (B10) | catalog "Lost" lines: KSmallestPairs vs Frederickson 1993 O(k); WeightedPicker vs Walker/Vose alias O(1); ReservoirSample vs Li 1994 Algorithm L; MinSquaresToSum worst case (~4.3e9 steps) stated in Time |
| L1 | bits (B10); sample (B12) | define n in "Time: O(n)" (pop_count_table.go:16, shuffle_deterministic.go:20) |
| L2 | text (B11) | WildcardMatch/RegexMatch Space: measured 0 allocations — correct the claim to what the code does |
| L3 | sorting (B1) | bucket_sort.go: document that dst must not partly overlap values |
| L4 | bloom (B14) | cap hashCount (k ≤ m, and a documented upper bound) |
| B1 | graph (B8) | benchmarks GridBFS/GridDFS/FloodFillInPlace: all-open grid or start inside the large component, so the work grows with n |
| B2 | matrix (B12) | MaxRectangleSumNoLargerThanK benchmark: even-valued grid with an odd k so it never exits early |
