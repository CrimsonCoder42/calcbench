# S10 algorithms core — report (requested as "S2"; S2 was taken by num/logs)

**Per package** (new IDs · `go test` final, out/go_test_final.txt · benchmarks out/review_bench.txt + bench_fix_1.txt):
sorting 10 ok 0.28s · search +6 ok 0.26s · selection +2 ok 0.22s · window +11 ok 2.36s · prefix 14 ok 0.29s ·
twopointer 6 ok 0.22s · hashing 13 ok 0.25s · stack 10 ok 0.27s · queue 2 ok 0.29s · heap +12 ok 0.38s ·
interval +8 ok 0.23s · list 17 ok 0.22s · tree 45 ok 0.54s · trie 13 ok 0.23s · graph 26 ok 0.33s · unionfind 8 ok 0.39s ·
dp 24 ok 0.27s · bits 16 ok 0.23s · intmath 16 ok 0.39s · text 21 ok 0.29s · enumerate 9 ok 0.24s · sample 7 ok 0.36s ·
matrix 6 ok 0.19s · geometry2d 8 ok 0.46s · ratelimit 4 ok 0.21s · bloom 4 ok 0.44s · seq 8 ok 0.29s · stream 4 ok 0.24s;
merge unchanged ok. 29 packages, all `ok`. Benchmarks: growth matches every documented Time (fixed-size ops 0 allocs;
n-D ×10 per decade; n log n and n² where stated; grid and MaxRectangleSum benchmarks re-run after fixture fixes).
Known heavy: stream DisjointIntervals O(k) per Add (documented); trie NewPrefixSuffixIndex 835 MB at 1e5 words.

**ALG functions in INDEX.md**: 13 before → 343 after (+330), plus the unused-ID row ALG-DP-016. Counts: code 343 |
INDEX 343 (+1 unused) | doc.go 343 (+1 unused); every row matches its doc comment.

**Reviewer** (one pass): 5 failing tests, all wrong expectations (KthLargest 512, RectanglesUnionArea 13,
denseMultiply p, Cols() on 0 rows, Gray allocation) — fixed; intmath test > 60 s (fixed: two oracles capped, 245 s → 0.4 s);
L26 not applied (size guards added everywhere, L41: n·elemSize ≤ 2⁴⁷ bytes → ErrTooLarge); problem-site example data in
9 tests (replaced with our own, hand-checked); Dijkstra/BellmanFord false ErrOverflow (fixed); zero-value RollingHash panic
(fixed); lazy min/max tree overflow (fixed: both extremes per node); Bloom zero/degenerate step (fixed: Dillinger–Manolios
enhanced double hashing) and unbounded k (capped); unsupported rolling-hash bound (removed); missing locators (added);
"Lost" lines for faster methods (added). All 11 L15 proofs checked correct by the reviewer.

**Held (no acceptable source), 17**: MedianOfTwoSorted · MaxProductSubarray · ContainerMostWater · TrappingRainWater ·
NumberOfDigitOne · IntegerToRoman · RomanToInteger · plain CompareVersions · DecodeWays · MaxProfit (k trades) ·
MaxProfit (cooldown) · MaximalSquare · SpiralOrder · SetZeroesInPlace · PerfectRectangle · SlidingWindowCounter ·
Bloom intersection. Routed to num/integer (D-005): GCD/LCM, ExtendedGCD, ModPow, ModInverse, Sieve, IsPrimeDeterministic64.
Dropped (= standard library): SimplifyPath, StrStr, SqrtNewton, four math/bits wrappers.

**Consumer changes**: 8 (moved to the S14 port work folder as CONSUMER_CHANGES_algorithms_core.md: consumer backend rank.go:19-37, rank.go:70, feed.go:725, rate_limit.go:273,
moments.go:119; red_lite_cache limiter.go:52, ordered.go:45, ordered.go:64/118/201). Nothing in the consumers was edited.

**FOUND, NOT FIXED**: graph's arc-count guard (graph.go:125) has no test (needs trillions of edges); Gries 1981 saddleback
cited at p. 215, not checked in print; not first-hand: Fredkin 1960, Härdle–Steiger 1995, Cohen Alg. 1.7.1,
Crandall–Pomerance 9.2.11, Concrete Mathematics (4.24), Knuth vol. 2 §4.3.1 Alg. A, Thompson 1968, Kreher–Stinson,
Devroye subsections, Vose 1991 / Li 1994 DOIs; sorting.ThreeWayPartition duplicates selection's unexported partition (L12);
PalindromePartitionMinCuts O(n²) and MatrixChainOrder O(n³) where faster methods exist.

**LEFT FOR DEVIN**: the L15 proofs (A-036); interview books as a source class (L31: 17 held); L41's cap (2⁴⁷ vs 128 MiB);
algo/matrix Transpose vs num/linalg MatTranspose (two shapes, one name); Laaksonen as a textbook (A-037); moving
PartialSort into selection to end L12; rule sync: add `iter` and `hash` to CLAUDE.md's forbidden package names
(CLAUDE.md "Naming" lists sort, slices, maps, … but not these two); go-functions.md says nothing about IDs for methods
(L23 used both shapes). Scratch: _delete/algo/dp/palindrome*.go (L33), _delete/algo/heap/probe_batch_select_test.go and
_delete/algo/selection/probe_fallback_test.go (S1 probes, already there).
