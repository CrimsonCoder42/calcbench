# Work: S2 algorithms — the core set, 2026-11-12

Request: catalog and build the core algorithm set in `algo/` (lead + builders in four waves + one reviewer).
Owned: `algo/**` except `algo/enumerate`/`algo/checksum` if another subject created them (neither exists:
combinatorics only proposed `algo/enumerate` in its PENDING_SHARED, so this run builds it; `algo/checksum`
is planned by D-016 for another subject and is not touched), this folder, `doc/catalog/algorithms.md`, and the
shared files for this run.

| File | What |
|---|---|
| BRIEF_COMMON.md | the brief every builder reads first |
| catalog/<package>.md | each builder's catalog section; the lead assembles doc/catalog/algorithms.md |
| PROGRESS.md | one line per function as it lands |
| REPORT.md | the final report |
| out/ | raw go test and benchmark output (gitignored) |

## Lead decisions (each with the alternative that lost)

| # | Decision | Alternative | Why it lost |
|---|---|---|---|
| L1 | `algo/sorting` (SORT) | `algo/sort` as requested | CLAUDE.md forbids package names a caller imports beside ours (`sort`) |
| L2 | `algo/hashing` (HASH) | `algo/hash` | `hash` is a standard package imported beside ours |
| L3 | `algo/seq` (SEQ) for range-over-func helpers | `algo/iter` | `iter` is the standard package every caller of these helpers imports |
| L4 | GCD/LCM (binary gcd), ExtendedGCD, ModPow, ModInverse, SieveOfEratosthenes, IsPrimeDeterministic64 routed to `num/integer` (request in PENDING_SHARED) | build them in `algo/intmath` | D-005 makes `num/integer` the home of primes, gcd, lcm and modular arithmetic; NUM-INTEG-001 IsPrime already exists; one function, one home |
| L5 | IntervalSchedulingMaxNonOverlap in `algo/interval` | `algo/heap` (as listed) | it is a greedy over intervals sorted by finish (CLRS §16.1), no heap |
| L6 | the monotonic queue's one home is `algo/window` (exported MonotonicDeque) | a second copy in `algo/queue` | one home; window already had the deque (unexported) |
| L7 | MovingMedian built in wave 2 (by the trie builder, owning only `algo/window/moving_median*.go`) | wave 1 with the window builder | it needs the exported heap built in the same wave by another builder |
| L8 | catalog sections in this folder per builder, assembled by the lead | builders editing doc/catalog/algorithms.md | two builders never own the same file |
| L9 | SearchSortedMatrix lives in `algo/search` | `algo/matrix` | it is a search; the request allowed either |
| L10 | geometry2d's two-rectangle area named RectanglesUnionArea | RectangleArea | NUM-GEOM RectangleArea (w·h) already exists; a second "RectangleArea" would confuse search |
| L11 | intra-algo imports allowed (acyclic) and reuse is by call | copying S1 code | the request: "every S1 reuse is a call, not a copy" |
| L12 | ThreeWayPartition exported once in `algo/sorting`; S1's unexported `partitionThreeWay` in selection/introselect.go stays | make selection call sorting | sorting imports selection (PartialSort → introselect), so the call would be a cycle; moving PartialSort into selection to break it is LEFT FOR DEVIN |
| L13 | CountingSort caps maxKey at 1<<24 (ASSUMPTION: a 128 MiB count table is the largest a library call should allocate silently) | no cap | an unbounded key would let one call allocate gigabytes |
| L14 | MedianOfTwoSorted HELD | build from secondary summaries | only a secondary summary of Frederickson–Johnson 1982 was readable; no-source rule |
| L15 | MaximalRectangle kept: the row→histogram reduction is proved in catalog/stack.md and the histogram step is cited (Berkman–Schieber–Vishkin 1993); flagged for Devin | HOLD it | it composes a cited method by a short proof; the question goes to Devin in REPORT |
| L16 | SimplifyPath DROPPED | build it | identical to the standard library's path.Clean on a rooted path |
| L17 | NaN is its own hash key in algo/hashing (never equal); TwoSum returns the pair with the smallest second index | other conventions | documented in the functions; deterministic outputs |
| L18 | Laaksonen, *Guide to Competitive Programming* (Springer; the Competitive Programmer's Handbook is its free text) accepted as a textbook source | reject it | it is a published textbook teaching the patterns, not a problem site; cited with section numbers |
| L19 | ExponentialMovingAverage in `algo/window` | `num/series` | it is a streaming pass over a sequence beside the other window functions; NIST §6.4.3.1 is the source |
| L20 | int128 / integer-constraint helpers stay unexported per package | one exported shared package | they are implementation details; exporting would create a public API nobody asked for |
| L21 | MoveZeroesInPlace (name says it modifies the input) | MoveZeroes | naming rule: a function that modifies its input says so |
| L22 | MaxProductSubarray, ContainerMostWater, TrappingRainWater HELD | build from problem-site knowledge | no acceptable primary source found |
| L23 | ID shape: a type gets one ID when its methods are plain accessors/mutators of one structure (queue, list caches); a method that is itself an algorithm (heap Push/Pop, Fenwick update/query) may carry its own ID; both shapes are recorded in INDEX | force one shape | wave-1 builders already shipped both; renumbering unreleased IDs buys nothing; the index lists what is callable either way |
| L24 | package `algo/list` keeps its taxonomy name despite container/list | rename | the name is approved in doc/taxonomy.md; callers alias on import if they need both |
| L25 | list functions that relink the caller's nodes are named …InPlace | bare names | naming rule (as L21) |
| L26 | constructors taking a size (NewDirected, NewUndirected, NewDisjointSet, NewHeap capacityHint, …) must return an error, not let make() panic, when n·elemSize exceeds what a slice length can hold; real out-of-memory is the runtime's | no guard | "a library function never panics on input"; the guard is the language's own limit, not an invented cap — to be fixed in the review pass |
| L27 | graph weights int64 with overflow checks; grid edits …InPlace; Cell defined in unionfind, aliased in graph; bidirectional WordLadder not built | float64 weights; dst copies | builder's catalog/graph.md gives each losing alternative |
| L28 | integer range sums report ErrOverflow when a partial sum overflows even if the total fits; Fenwick types accept signed and float types; min/max trees return ErrEmptyRange on an empty range | exact-total int128 arithmetic everywhere | documented; the partial sums are the data structure's state |
| L29 | practitioner names over problem-site names: PopCountTable, LoneElementAmongPairs/Triples, Submasks, PowerBySquaring, DigitalRoot, IsHammingNumber/NthHammingNumber, MinSquaresToSum, MultiplyDecimalStrings, IncrementDigits, LongDivision | the request's LeetCode-flavoured names (CountingBits, SingleNumber, UglyNumber, PlusOne…) | naming rule: no problem-site names; the request's names are listed under Tags |
| L30 | functions that are exactly one standard-library call are DROPPED (OnesCount64, TrailingZeros64, Len64, Reverse64, math.Sqrt) | thin wrappers | one function, one home (the standard library) |
| L31 | interview books (CtCI, EPI) are not accepted as primary sources: SpiralOrder and SetZeroesInPlace HELD | accept them | they collect problems; they are not the origin or a recognised textbook of the method (rules: problem sites teach the pattern only) |
| L32 | enumeration iterators yield one reused buffer (callers copy); functions taking a bare count return (iter.Seq, error); RandomElement replaces "LinkedListRandomNode" (any iter.Seq) | fresh slice per step; list-only API | zero allocation per step; the general iterator covers the list case |
| L33 | LongestPalindromicSubstring's one home is algo/text (O(n) on Manacher's PalindromeRadii); dp does not build it | dp's expand-around-centre version too | one name, one home; the O(n) method dominates |
| L34 | text is byte-level by default (generic ~string \| ~[]byte), rune-level where documented; RLE uses PackBits (TIFF 6.0 §9); StringToInt follows C11 strtol with an error on overflow; ColumnLabel/ColumnNumber; CompareSemVer (SemVer 2.0.0 §11) | ad-hoc formats; product names | each has a citable grammar or standard |
| L35 | DP tables capped at 128 MiB (ErrTooLarge), as L13; ClimbingStairs folded into Fibonacci (documented as Fibonacci(n+1)); CoinChangeCount uint64 with ErrOverflow | no cap; separate ClimbingStairs; big.Int | one bound across the run; one function per concept |
| L36 | ALG-DP-016 left unused (never released, not reused) | renumber 017…025 down | an ID once given is never reused or renumbered |
| L37 | geometry2d points are exact int64 only (128/192-bit internal products: every predicate exact over the whole range) | float64 with Shewchuk adaptive predicates | exactness without a robustness layer; no float intersection point is offered |
| L38 | ratelimit: a clock going backwards counts as the latest time seen (no error); impossible requests return (false, Never); n ≤ 0 → (true, 0) | an error on a backwards clock | wall clocks step backwards in practice; an error would deny service on NTP adjustment |
| L39 | Bloom indices from SHA-256 (crypto/sha256, standard library) split into h₁, h₂ for Kirsch–Mitzenmacher double hashing | hash/fnv (requested option); hash/maphash | the builder's measurement: 128-bit FNV-1a halves are correlated (FPR 0.00093 vs 0.0100 predicted); maphash's seed cannot be persisted, so a filter could not be saved |
| L40 | SlidingWindowCounter HELD; Bloom intersection not built | build from vendor blogs | no acceptable source |
