# calcbench functions (from INDEX.md, 2026-11-13)

Every live exported function: 1,297 rows (the 28 REMOVED rows are left out). Columns: ID | Function | Package | Does.
Grep by meaning; open a package only to confirm a signature (`go doc ./<pkg> <Func>`).

| ID | Function | Package | Does |
|---|---|---|---|
| ALG-BITS-001 | IsPowerOfTwo | algo/bits | whether x is 2^k for some k ≥ 0 |
| ALG-BITS-002 | IsPowerOfFour | algo/bits | whether x is 4^k for some k ≥ 0 |
| ALG-BITS-003 | NextPowerOfTwo | algo/bits | least power of two ≥ x (ErrOverflow above 2^63) |
| ALG-BITS-004 | ReverseBits | algo/bits | the low width bits of x in reverse order |
| ALG-BITS-005 | GrayEncode | algo/bits | the reflected binary Gray code of x |
| ALG-BITS-006 | GrayDecode | algo/bits | the number whose Gray code is g |
| ALG-BITS-007 | LowestSetBit | algo/bits | x with every bit but its lowest set bit cleared |
| ALG-BITS-008 | HighestSetBit | algo/bits | x with every bit but its highest set bit cleared |
| ALG-BITS-009 | PopCountTable | algo/bits | the number of set bits of 0, 1, …, n, in O(n) |
| ALG-BITS-010 | LoneElementAmongPairs | algo/bits | the one value present an odd number of times when the rest are paired |
| ALG-BITS-011 | LoneElementAmongTriples | algo/bits | the one value present once when every other is present three times |
| ALG-BITS-012 | BitwiseAndRange | algo/bits | lo & (lo+1) & … & hi, in O(1) |
| ALG-BITS-013 | HammingDistance | algo/bits | the number of bit positions where x and y differ |
| ALG-BITS-014 | Submasks | algo/bits | every submask of a mask, descending, as an iterator |
| ALG-BITS-015 | MortonEncode2D | algo/bits | interleave the bits of (x, y) into a Z-order code |
| ALG-BITS-016 | MortonDecode2D | algo/bits | split a Z-order code back into (x, y) |
| ALG-BLOOM-001 | BloomFilter | algo/bloom | NewBloomFilter, Add, MayContain, EstimatedFPR, BitCount, HashCount, SetBits |
| ALG-BLOOM-002 | FalsePositiveRate | algo/bloom | (1 − e^(−kn/m))^k for m bits, n items, k hashes |
| ALG-BLOOM-003 | OptimalParams | algo/bloom | m and k for n items at a target false-positive rate |
| ALG-BLOOM-004 | BloomFilter.UnionWith | algo/bloom | OR another filter's bits into this one |
| ALG-DP-001 | LongestIncreasingSubsequence | algo/dp | indices of one longest strictly increasing subsequence (patience sorting) |
| ALG-DP-002 | LongestCommonSubsequence | algo/dp | one longest common subsequence in linear space (Hirschberg) |
| ALG-DP-003 | LCSLength | algo/dp | length of a longest common subsequence, one row of memory |
| ALG-DP-004 | EditDistance | algo/dp | Levenshtein distance, one row of memory |
| ALG-DP-005 | OneEditDistance | algo/dp | whether two sequences are exactly one edit apart |
| ALG-DP-006 | Knapsack01 | algo/dp | best value and items of the 0-1 knapsack |
| ALG-DP-007 | UnboundedKnapsack | algo/dp | best value and copies per item, items reusable |
| ALG-DP-008 | CoinChangeMin | algo/dp | fewest coins that make an amount |
| ALG-DP-009 | CoinChangeCount | algo/dp | number of coin combinations that make an amount |
| ALG-DP-010 | Fibonacci | algo/dp | F(n) for n ≤ 93; Fibonacci(n+1) is the climbing-stairs count |
| ALG-DP-011 | MaxNonAdjacentSum | algo/dp | largest sum of pairwise non-adjacent elements (path) |
| ALG-DP-012 | MaxNonAdjacentSumCircular | algo/dp | the same on a ring |
| ALG-DP-013 | GridPathCount | algo/dp | right/down paths across a grid with blocked cells |
| ALG-DP-014 | MinPathSumGrid | algo/dp | least-cost right/down path across a grid |
| ALG-DP-015 | TextSegmentation | algo/dp | split text into dictionary words |
| ALG-DP-017 | PalindromePartitionMinCuts | algo/dp | fewest cuts into palindromic pieces |
| ALG-DP-018 | RodCutting | algo/dp | best revenue and pieces for a rod (CLRS §15.1) |
| ALG-DP-019 | MatrixChainOrder | algo/dp | fewest multiplications and parenthesisation (CLRS §15.2) |
| ALG-DP-020 | CanJumpToEnd | algo/dp | whether forward jumps reach the last index |
| ALG-DP-021 | MinJumpsToEnd | algo/dp | fewest forward jumps to the last index |
| ALG-DP-022 | MaxProfitOneTrade | algo/dp | best single buy-then-sell |
| ALG-DP-023 | MaxProfitUnlimitedTrades | algo/dp | best total over any number of non-overlapping trades |
| ALG-DP-024 | MinPathSumTriangle | algo/dp | least-sum apex-to-base path in a number triangle |
| ALG-DP-025 | SubsequenceOccurrences | algo/dp | times a pattern occurs in a text as a subsequence |
| ALG-ENUM-001 | NextPermutationInPlace | algo/enumerate | the next permutation in lexicographic order, in place |
| ALG-ENUM-002 | Permutations | algo/enumerate | all n! orderings of the positions, lexicographic |
| ALG-ENUM-003 | PermutationsUnique | algo/enumerate | the distinct permutations of a multiset, lexicographic |
| ALG-ENUM-004 | Combinations | algo/enumerate | the k-element combinations of the positions, lexicographic |
| ALG-ENUM-005 | Subsets | algo/enumerate | all 2ⁿ subsets, binary counting order |
| ALG-ENUM-006 | SubsetsUnique | algo/enumerate | the distinct sub-multisets of a multiset |
| ALG-ENUM-007 | CombinationSum | algo/enumerate | the multisets of candidates (reusable) with a given sum |
| ALG-ENUM-008 | GenerateParentheses | algo/enumerate | the balanced strings of n pairs of parentheses |
| ALG-ENUM-009 | GrayCodeSequence | algo/enumerate | the binary reflected Gray code on a number of bits |
| ALG-GEO2D-001 | Orientation | algo/geometry2d | the turn of three points: counterclockwise, clockwise or collinear |
| ALG-GEO2D-002 | ConvexHull | algo/geometry2d | the strict convex hull, counterclockwise (Andrew's monotone chain) |
| ALG-GEO2D-003 | ClosestPair | algo/geometry2d | the indices of two points at the least distance, O(n log n) |
| ALG-GEO2D-004 | PointInPolygon | algo/geometry2d | Inside, Outside or OnBoundary of a polygon (crossing number) |
| ALG-GEO2D-005 | SegmentsIntersect | algo/geometry2d | whether two closed segments share a point |
| ALG-GEO2D-006 | MaxPointsOnLine | algo/geometry2d | the most points on one straight line |
| ALG-GEO2D-007 | RectanglesUnionArea | algo/geometry2d | the area covered by two axis-aligned rectangles |
| ALG-GEO2D-008 | HasVerticalSymmetryAxis | algo/geometry2d | whether a vertical line reflects the point set onto itself |
| ALG-GRAPH-001 | NewDirected | algo/graph | a directed graph from an edge list |
| ALG-GRAPH-002 | NewUndirected | algo/graph | an undirected graph from an edge list |
| ALG-GRAPH-003 | Graph.AddEdge | algo/graph | add one edge |
| ALG-GRAPH-004 | Graph.VertexCount | algo/graph | the number of vertices |
| ALG-GRAPH-005 | Graph.EdgeCount | algo/graph | the number of edges |
| ALG-GRAPH-006 | Graph.Directed | algo/graph | whether the graph is directed |
| ALG-GRAPH-007 | Graph.Neighbors | algo/graph | a vertex's (neighbour, weight) pairs as an iter.Seq2 |
| ALG-GRAPH-008 | Graph.Edges | algo/graph | every edge once, as a slice |
| ALG-GRAPH-009 | BFS | algo/graph | vertices reachable from a source, breadth-first, as an iter.Seq |
| ALG-GRAPH-010 | DFS | algo/graph | vertices reachable from a source, depth-first preorder, as an iter.Seq |
| ALG-GRAPH-011 | ShortestPathUnweighted | algo/graph | a fewest-edges path between two vertices (BFS) |
| ALG-GRAPH-012 | Dijkstra | algo/graph | single-source shortest paths, non-negative weights |
| ALG-GRAPH-013 | BellmanFord | algo/graph | single-source shortest paths, any weights, negative-cycle check |
| ALG-GRAPH-014 | FloydWarshall | algo/graph | all-pairs shortest-path distances |
| ALG-GRAPH-015 | TopologicalSort | algo/graph | a dependency order of a DAG (Kahn), cycle detection |
| ALG-GRAPH-016 | ConnectedComponents | algo/graph | a component label per vertex (union-find) |
| ALG-GRAPH-017 | IsBipartite | algo/graph | a two-side split with every edge across, if one exists |
| ALG-GRAPH-018 | Kruskal | algo/graph | a minimum spanning forest (sorted edges, union-find) |
| ALG-GRAPH-019 | Prim | algo/graph | a minimum spanning forest (grown with a heap) |
| ALG-GRAPH-020 | CloneGraph | algo/graph | a deep copy of a graph |
| ALG-GRAPH-021 | WordLadder | algo/graph | a shortest one-letter-at-a-time word ladder |
| ALG-GRAPH-022 | GridBFS | algo/graph | grid cells reachable from a start, breadth-first, as an iter.Seq |
| ALG-GRAPH-023 | GridDFS | algo/graph | grid cells reachable from a start, depth-first preorder, as an iter.Seq |
| ALG-GRAPH-024 | NumberOfIslands | algo/graph | the number of connected land groups on a grid |
| ALG-GRAPH-025 | FloodFillInPlace | algo/graph | repaint the region of one value around a cell (modifies the grid) |
| ALG-GRAPH-026 | SurroundedRegionsInPlace | algo/graph | fill regions of one value that cannot reach the border (modifies the grid) |
| ALG-HASH-001 | FrequencyCount | algo/hashing | how many times each distinct value occurs |
| ALG-HASH-002 | TopKFrequent | algo/hashing | the k most frequent values, with their counts |
| ALG-HASH-003 | TopKByCount | algo/hashing | the k largest entries of a count map |
| ALG-HASH-004 | GroupBy | algo/hashing | values grouped by a key function, in first-appearance order |
| ALG-HASH-005 | FirstUnique | algo/hashing | index of the first value that occurs exactly once |
| ALG-HASH-006 | SetIntersection | algo/hashing | distinct values of a that occur in b |
| ALG-HASH-007 | SetUnion | algo/hashing | distinct values of a, then of b |
| ALG-HASH-008 | SetDifference | algo/hashing | distinct values of a that do not occur in b |
| ALG-HASH-009 | LongestConsecutive | algo/hashing | longest run of consecutive integers present |
| ALG-HASH-010 | GroupAnagrams | algo/hashing | words grouped by sorted-rune signature |
| ALG-HASH-011 | GroupAnagramsByCount | algo/hashing | words grouped by rune counts over an alphabet |
| ALG-HASH-012 | TwoSum | algo/hashing | indices of two values that sum to a target |
| ALG-HASH-013 | ContainsNearbyDuplicate | algo/hashing | whether equal values occur within a given index gap |
| ALG-HEAP-001 | TopK | algo/heap | the k largest elements, in descending order |
| ALG-HEAP-002 | NewHeap | algo/heap | an empty Heap ordered by a less function |
| ALG-HEAP-003 | Heapify | algo/heap | a Heap built from a copy of a slice in O(n) (Floyd) |
| ALG-HEAP-004 | Heap.Push | algo/heap | add an element, O(log n) |
| ALG-HEAP-005 | Heap.Pop | algo/heap | remove and return the first element, ok false when empty |
| ALG-HEAP-006 | Heap.Peek | algo/heap | the first element without removing it |
| ALG-HEAP-007 | Heap.Len | algo/heap | the number of elements |
| ALG-HEAP-008 | KSmallestPairs | algo/heap | the k index pairs with the smallest sums from two ascending slices |
| ALG-HEAP-009 | KthSmallestInSortedMatrix | algo/heap | the k-th smallest element of a row- and column-ascending matrix |
| ALG-HEAP-010 | NewMedianStream | algo/heap | a presized running-median stream (the zero MedianStream is ready too) |
| ALG-HEAP-011 | MedianStream.Add | algo/heap | add a value, O(log n) |
| ALG-HEAP-012 | MedianStream.Median | algo/heap | the lower and upper medians so far, O(1) |
| ALG-HEAP-013 | MedianStream.Len | algo/heap | the number of values added |
| ALG-INTMTH-001 | PowerBySquaring | algo/intmath | base^exp for int64, overflow checked |
| ALG-INTMTH-002 | ISqrt | algo/intmath | ⌊√n⌋ for every uint64, exact |
| ALG-INTMTH-003 | LongDivision | algo/intmath | quotient and remainder by shift-and-subtract, no division operator |
| ALG-INTMTH-004 | ReverseDigits | algo/intmath | the decimal digits of x in reverse order, overflow checked |
| ALG-INTMTH-005 | IsPalindromeNumber | algo/intmath | whether x reads the same in decimal both ways |
| ALG-INTMTH-006 | IsHappyNumber | algo/intmath | whether iterating the sum of squared digits reaches 1 |
| ALG-INTMTH-007 | DigitalRoot | algo/intmath | the repeated digit sum of n, in O(1) |
| ALG-INTMTH-008 | FactorialTrailingZeros | algo/intmath | the number of trailing zeros of n! in decimal |
| ALG-INTMTH-009 | IsHammingNumber | algo/intmath | whether n = 2^i·3^j·5^k (a regular, 5-smooth number) |
| ALG-INTMTH-010 | NthHammingNumber | algo/intmath | the n-th Hamming number, ascending from 1 |
| ALG-INTMTH-011 | MinSquaresToSum | algo/intmath | the least number of perfect squares that sum to n |
| ALG-INTMTH-012 | IsPowerOfThree | algo/intmath | whether x is 3^k for some k ≥ 0 |
| ALG-INTMTH-013 | FractionToDecimal | algo/intmath | the decimal expansion of a fraction, the repetend in parentheses |
| ALG-INTMTH-014 | AddBinaryStrings | algo/intmath | the sum of two binary digit strings |
| ALG-INTMTH-015 | MultiplyDecimalStrings | algo/intmath | the product of two decimal digit strings |
| ALG-INTMTH-016 | IncrementDigits | algo/intmath | a decimal digit slice plus one |
| ALG-INTVL-001 | MergeOverlapping | algo/interval | the union of closed intervals as disjoint sorted intervals |
| ALG-INTVL-002 | InsertInterval | algo/interval | a sorted, disjoint list with one interval merged in |
| ALG-INTVL-003 | IntersectLists | algo/interval | the intersection of two sorted, disjoint lists |
| ALG-INTVL-004 | MinMeetingRooms | algo/interval | the fewest rooms that hold every interval (the depth) |
| ALG-INTVL-005 | MissingRanges | algo/interval | the integers of [lower, upper] absent from an ascending slice, as ranges |
| ALG-INTVL-006 | SummaryRanges | algo/interval | ascending integers as ranges of consecutive values |
| ALG-INTVL-007 | Covers | algo/interval | whether a union of intervals contains a target interval |
| ALG-INTVL-008 | IntervalSchedulingMaxNonOverlap | algo/interval | a largest set of pairwise disjoint intervals (greedy by End) |
| ALG-INTVL-009 | Skyline | algo/interval | the outline of rectangular buildings as key points |
| ALG-LIST-001 | Node | algo/list | singly linked list node (Value, Next) |
| ALG-LIST-002 | FromSlice | algo/list | a list holding a slice's values, one allocation |
| ALG-LIST-003 | AppendValues | algo/list | append a list's values to a slice |
| ALG-LIST-004 | Length | algo/list | the number of nodes, ErrCycle on a cycle |
| ALG-LIST-005 | DoublyLinkedList | algo/list | generic doubly linked list (Element, iterators), O(1) per operation |
| ALG-LIST-006 | ReverseInPlace | algo/list | reverse a list by relinking |
| ALG-LIST-007 | ReverseKGroupInPlace | algo/list | reverse each block of k nodes |
| ALG-LIST-008 | DetectCycle | algo/list | whether a list loops (Floyd) |
| ALG-LIST-009 | CycleStart | algo/list | the first node of the cycle and its length |
| ALG-LIST-010 | Middle | algo/list | the node at index ⌊n/2⌋ |
| ALG-LIST-011 | MergeTwoSortedInPlace | algo/list | stable merge of two ascending lists |
| ALG-LIST-012 | RemoveNthFromEndInPlace | algo/list | unlink the k-th node from the end |
| ALG-LIST-013 | ReorderInPlace | algo/list | relink to L0, Ln−1, L1, Ln−2, … |
| ALG-LIST-014 | PartitionInPlace | algo/list | stable partition around a pivot |
| ALG-LIST-015 | LRUCache | algo/list | least-recently-used cache, O(1) per operation |
| ALG-LIST-016 | LFUCache | algo/list | least-frequently-used cache (LRU among ties), O(1) per operation |
| ALG-LIST-017 | AddTwoNumbers | algo/list | sum of two digit lists, least significant digit first |
| ALG-MATRIX-001 | Rotate90InPlace | algo/matrix | rotate a square grid a quarter turn clockwise, in place |
| ALG-MATRIX-002 | Transpose | algo/matrix | the transpose of a grid |
| ALG-MATRIX-003 | NewCSR | algo/matrix | compressed sparse row form of a dense grid |
| ALG-MATRIX-004 | SparseMultiply | algo/matrix | the product of two CSR matrices (Gustavson) |
| ALG-MATRIX-005 | GameOfLifeStep | algo/matrix | one generation of Conway's Game of Life |
| ALG-MATRIX-006 | MaxRectangleSumNoLargerThanK | algo/matrix | the largest sum of a sub-rectangle that is ≤ k |
| ALG-MERGE-001 | MergeKSorted | algo/merge | one ascending slice from k ascending slices, stable |
| ALG-PREFIX-001 | PrefixSums | algo/prefix | the n+1 prefix sums of a slice |
| ALG-PREFIX-002 | RangeSum | algo/prefix | the sum of values[lo:hi] from its prefix sums, O(1) |
| ALG-PREFIX-003 | PrefixSums2D | algo/prefix | the summed-area table of a grid |
| ALG-PREFIX-004 | SummedAreaTable.RectangleSum | algo/prefix | the sum of a rectangle of the grid, O(1) |
| ALG-PREFIX-005 | NewDifferenceArray | algo/prefix | a difference array holding a copy of a slice |
| ALG-PREFIX-006 | DifferenceArray.RangeAdd | algo/prefix | add a constant to a range, O(1) |
| ALG-PREFIX-007 | DifferenceArray.Materialize | algo/prefix | the current elements, O(n) |
| ALG-PREFIX-008 | CountSubarraysSumEquals | algo/prefix | the number of subarrays with a given sum |
| ALG-PREFIX-009 | MaxSubarray | algo/prefix | the largest subarray sum and its range (Kadane) |
| ALG-PREFIX-010 | ProductExceptSelf | algo/prefix | for each i, the product of all values but values[i], no division |
| ALG-PREFIX-011 | EquilibriumIndex | algo/prefix | the first index whose left and right sums are equal |
| ALG-PREFIX-012 | DifferenceArray.Len | algo/prefix | the number of elements |
| ALG-PREFIX-013 | SummedAreaTable.Rows | algo/prefix | the number of grid rows |
| ALG-PREFIX-014 | SummedAreaTable.Cols | algo/prefix | the number of grid columns |
| ALG-QUEUE-001 | QueueFromTwoStacks | algo/queue | FIFO queue from two stacks, O(1) amortized per operation |
| ALG-QUEUE-002 | CircularBuffer | algo/queue | fixed-capacity ring buffer that rejects or overwrites when full |
| ALG-RATELM-001 | TokenBucket | algo/ratelimit | token bucket, exact integer refill, n tokens per request (RFC 2697 §3) |
| ALG-RATELM-002 | LeakyBucket | algo/ratelimit | GCRA virtual scheduling, one cell per request (ITU-T I.371 Annex A) |
| ALG-RATELM-003 | FixedWindow | algo/ratelimit | jumping window: at most limit per consecutive window of length W |
| ALG-RATELM-004 | SlidingWindowLog | algo/ratelimit | moving window: at most limit in any span of length W, exact log |
| ALG-SAMPLE-001 | ShuffleInPlace | algo/sample | uniform random permutation in place (Fisher–Yates) |
| ALG-SAMPLE-002 | ReservoirSample | algo/sample | k uniform items from a sequence of unknown length (Algorithm R) |
| ALG-SAMPLE-003 | NewWeightedPicker | algo/sample | prepares integer weights for weighted picks (prefix sums) |
| ALG-SAMPLE-004 | WeightedPicker.Pick | algo/sample | an index chosen with probability weight/total, O(log n) |
| ALG-SAMPLE-005 | RejectionSample | algo/sample | a draw from a proposal accepted by the caller's test |
| ALG-SAMPLE-006 | RandomElement | algo/sample | one uniform item from a sequence of unknown length |
| ALG-SAMPLE-007 | ShuffleDeterministic | algo/sample | a permutation of 0…n−1 fixed by a seed, for tests |
| ALG-SEARCH-001 | LowerBound | algo/search | first index whose element is not less than the target |
| ALG-SEARCH-002 | UpperBound | algo/search | first index whose element is greater than the target |
| ALG-SEARCH-003 | BinarySearchFirst | algo/search | index of the first element equal to the target |
| ALG-SEARCH-004 | BinarySearchLast | algo/search | index of the last element equal to the target |
| ALG-SEARCH-005 | SearchRotated | algo/search | index of the target in a rotated sorted slice of distinct values |
| ALG-SEARCH-006 | SearchRotatedDup | algo/search | the same with duplicates allowed (O(n) worst case) |
| ALG-SEARCH-007 | PeakElement | algo/search | index of an element not less than its neighbours |
| ALG-SEARCH-008 | SearchPredicate | algo/search | smallest x in [lo, hi) where a monotone predicate is true |
| ALG-SEARCH-009 | SearchSortedMatrix | algo/search | position of the target in a row-major sorted matrix |
| ALG-SEARCH-010 | SearchSortedMatrixSaddleback | algo/search | position of the target in a row- and column-sorted matrix |
| ALG-SELECT-001 | KthSmallest | algo/selection | the k-th smallest element, k = 1 being the minimum |
| ALG-SELECT-002 | Median | algo/selection | the lower and upper middle elements |
| ALG-SELECT-003 | KthLargest | algo/selection | the k-th largest element, k = 1 being the maximum |
| ALG-SELECT-004 | KthSmallestInto | algo/selection | the k-th smallest, selected in a caller-supplied buffer (no allocation) |
| ALG-SEQ-001 | Map | algo/seq | each element transformed by a function |
| ALG-SEQ-002 | Filter | algo/seq | the elements a predicate keeps |
| ALG-SEQ-003 | Take | algo/seq | the first n elements, then stop the upstream |
| ALG-SEQ-004 | Chunk | algo/seq | consecutive groups of exactly size elements (the last may be shorter) |
| ALG-SEQ-005 | Zip | algo/seq | pairs from two iterators, stopping at the shorter (iter.Pull, stop always called) |
| ALG-SEQ-006 | Flatten | algo/seq | the elements of a sequence of sequences, in order |
| ALG-SEQ-007 | Batch | algo/seq | groups closed at a size limit OR when a predicate marks an element as the last |
| ALG-SEQ-008 | Window | algo/seq | every run of width consecutive elements (sliding window), reused ring buffer |
| ALG-SORT-001 | MergeSortStable | algo/sorting | stable sort by a compare func (bottom-up merge sort) |
| ALG-SORT-002 | InsertionSort | algo/sorting | sort for small or nearly sorted input |
| ALG-SORT-003 | CountingSort | algo/sorting | sort of ints in [0, maxKey] in O(n + maxKey) |
| ALG-SORT-004 | RadixSortUint64 | algo/sorting | LSD radix sort of uint64 keys, 8-bit digits |
| ALG-SORT-005 | BucketSort | algo/sorting | sort of floats in [0, 1), expected linear time |
| ALG-SORT-006 | IsSorted | algo/sorting | whether a slice is ascending, and where it first is not |
| ALG-SORT-007 | DedupeSorted | algo/sorting | the distinct values of a sorted slice |
| ALG-SORT-008 | PartialSort | algo/sorting | the k smallest values, ascending |
| ALG-SORT-009 | SortIndices | algo/sorting | the stable sorting permutation (argsort) |
| ALG-SORT-010 | ThreeWayPartition | algo/sorting | less / equal / greater blocks around a pivot (Dutch national flag) |
| ALG-STACK-001 | NextGreater | algo/stack | index of the next strictly greater element |
| ALG-STACK-002 | NextSmaller | algo/stack | index of the next strictly smaller element |
| ALG-STACK-003 | PreviousGreater | algo/stack | index of the previous strictly greater element |
| ALG-STACK-004 | PreviousSmaller | algo/stack | index of the previous strictly smaller element |
| ALG-STACK-005 | LargestRectangleHistogram | algo/stack | area of the largest rectangle under a histogram |
| ALG-STACK-006 | MaximalRectangle | algo/stack | area of the largest all-true rectangle in a grid |
| ALG-STACK-007 | DailyTemperatures | algo/stack | positions to wait for a strictly greater reading |
| ALG-STACK-008 | MinStack | algo/stack | stack with Push, Pop, Top and Min in O(1) |
| ALG-STACK-009 | EvaluateRPN | algo/stack | value of an int64 reverse Polish expression |
| ALG-STACK-010 | BalancedBrackets | algo/stack | whether brackets of the given pairs are balanced |
| ALG-STREAM-001 | MovingAverageStream | algo/stream | mean of the last windowLen values, O(1) amortized per value (queue.CircularBuffer) |
| ALG-STREAM-002 | HitCounter | algo/stream | hits in the last window of time, in per-resolution buckets of a ring |
| ALG-STREAM-003 | DisjointIntervals | algo/stream | integers added so far, as sorted disjoint closed ranges |
| ALG-STREAM-004 | RunningStats | algo/stream | count, mean, population and sample variance, min, max; Merge combines two |
| ALG-TEXT-001 | PrefixFunction | algo/text | longest proper border of every prefix (KMP failure function) |
| ALG-TEXT-002 | KMPSearch | algo/text | every start of a pattern in a text, Knuth–Morris–Pratt |
| ALG-TEXT-003 | ZFunction | algo/text | longest common prefix of the string and each suffix |
| ALG-TEXT-004 | NewRollingHash | algo/text | polynomial hash of a sliding window (RollingHash: Append, Roll, Sum, Len, Reset) |
| ALG-TEXT-005 | RabinKarpSearch | algo/text | every start of a pattern, by fingerprints, each hit verified |
| ALG-TEXT-006 | RepeatedSubstrings | algo/text | first starts of the length-L substrings that occur twice or more |
| ALG-TEXT-007 | LongestCommonPrefix | algo/text | length of the prefix shared by every string |
| ALG-TEXT-008 | PalindromeRadii | algo/text | longest palindrome at every centre, Manacher |
| ALG-TEXT-009 | LongestPalindromicSubstring | algo/text | the longest palindromic substring |
| ALG-TEXT-010 | IsAnagram | algo/text | whether two strings hold the same runes |
| ALG-TEXT-011 | RunLengthEncode | algo/text | PackBits run-length encoding |
| ALG-TEXT-012 | RunLengthDecode | algo/text | PackBits run-length decoding, with an output cap |
| ALG-TEXT-013 | ReverseWords | algo/text | the words in reverse order, single-spaced |
| ALG-TEXT-014 | MinimumWindowSubstring | algo/text | shortest window holding every byte of a target |
| ALG-TEXT-015 | StringToInt | algo/text | strtol-style base-10 prefix parse, overflow is an error |
| ALG-TEXT-016 | IsJSONNumber | algo/text | whether a string is a JSON number (RFC 8259 §6) |
| ALG-TEXT-017 | ColumnLabel | algo/text | spreadsheet column label of a number (bijective base 26) |
| ALG-TEXT-018 | ColumnNumber | algo/text | number of a spreadsheet column label |
| ALG-TEXT-019 | CompareSemVer | algo/text | Semantic Versioning 2.0.0 precedence of two versions |
| ALG-TEXT-020 | WildcardMatch | algo/text | whole-text match with '?' and '*' |
| ALG-TEXT-021 | RegexMatch | algo/text | whole-text match with '.' and '*' |
| ALG-TREE-001 | BST.Insert | algo/tree | add a key, false if present |
| ALG-TREE-002 | BST.Search | algo/tree | whether a key is present |
| ALG-TREE-003 | BST.Delete | algo/tree | remove a key (Hibbard deletion by successor), false if absent |
| ALG-TREE-004 | BST.Len | algo/tree | the number of keys |
| ALG-TREE-005 | BST.Root | algo/tree | the root node, for the traversals and queries |
| ALG-TREE-006 | Inorder | algo/tree | iterator over values in inorder |
| ALG-TREE-007 | Preorder | algo/tree | iterator over values in preorder |
| ALG-TREE-008 | Postorder | algo/tree | iterator over values in postorder |
| ALG-TREE-009 | LevelOrder | algo/tree | iterator over values level by level |
| ALG-TREE-010 | Height | algo/tree | edges on the longest root-to-leaf path (−1 when empty) |
| ALG-TREE-011 | IsBalanced | algo/tree | whether subtree heights differ by at most 1 everywhere |
| ALG-TREE-012 | LowestCommonAncestorBST | algo/tree | the LCA of two keys in a BST |
| ALG-TREE-013 | LowestCommonAncestor | algo/tree | the LCA of two nodes in a binary tree |
| ALG-TREE-014 | FromPreIn | algo/tree | the tree with a given preorder and inorder |
| ALG-TREE-015 | FromInPost | algo/tree | the tree with a given inorder and postorder |
| ALG-TREE-016 | SortedToBalancedBST | algo/tree | a perfectly balanced BST from strictly ascending keys |
| ALG-TREE-017 | KthSmallestBST | algo/tree | the k-th smallest key of a BST |
| ALG-TREE-018 | ValidateBST | algo/tree | nil when a tree is a BST, ErrNotBST otherwise |
| ALG-TREE-019 | Serialize | algo/tree | append a tree in preorder with '#' for empty subtrees and netstring values |
| ALG-TREE-020 | Deserialize | algo/tree | the tree Serialize wrote |
| ALG-TREE-021 | MaxPathSum | algo/tree | the largest sum along any path |
| ALG-TREE-022 | RightSideView | algo/tree | the last value of each level |
| ALG-TREE-023 | NewFenwickTree | algo/tree | a Fenwick tree over a copy of values, built in O(n) |
| ALG-TREE-024 | FenwickTree.Add | algo/tree | add to one element |
| ALG-TREE-025 | FenwickTree.PrefixSum | algo/tree | the sum of the first k elements |
| ALG-TREE-026 | FenwickTree.RangeSum | algo/tree | the sum of [lo, hi) |
| ALG-TREE-027 | FenwickTree.Len | algo/tree | the number of elements |
| ALG-TREE-028 | NewFenwickRangeUpdate | algo/tree | range-add, point-query structure over a copy of values |
| ALG-TREE-029 | FenwickRangeUpdate.RangeAdd | algo/tree | add a constant to [lo, hi) |
| ALG-TREE-030 | FenwickRangeUpdate.Get | algo/tree | one element |
| ALG-TREE-031 | FenwickRangeUpdate.Len | algo/tree | the number of elements |
| ALG-TREE-032 | NewSegmentTree | algo/tree | a segment tree over any monoid (combine + identity) |
| ALG-TREE-033 | NewSumSegmentTree | algo/tree | a segment tree of sums |
| ALG-TREE-034 | NewMinSegmentTree | algo/tree | a segment tree of minimums |
| ALG-TREE-035 | NewMaxSegmentTree | algo/tree | a segment tree of maximums |
| ALG-TREE-036 | SegmentTree.Set | algo/tree | replace one element |
| ALG-TREE-037 | SegmentTree.Query | algo/tree | the combined value of [lo, hi) |
| ALG-TREE-038 | SegmentTree.Len | algo/tree | the number of elements |
| ALG-TREE-039 | NewLazySegmentTree | algo/tree | a lazy segment tree defined by LazyOps |
| ALG-TREE-040 | NewRangeAddSumTree | algo/tree | range add, range sum |
| ALG-TREE-041 | NewRangeAddMinTree | algo/tree | range add, range minimum |
| ALG-TREE-042 | NewRangeAddMaxTree | algo/tree | range add, range maximum |
| ALG-TREE-043 | LazySegmentTree.RangeUpdate | algo/tree | apply an update to [lo, hi) |
| ALG-TREE-044 | LazySegmentTree.Query | algo/tree | the combined value of [lo, hi) |
| ALG-TREE-045 | LazySegmentTree.Len | algo/tree | the number of elements |
| ALG-TRIE-001 | Trie.Insert | algo/trie | store a value under a key; reports a replacement |
| ALG-TRIE-002 | Trie.Search | algo/trie | the value under a key |
| ALG-TRIE-003 | Trie.HasPrefix | algo/trie | whether any key starts with a prefix |
| ALG-TRIE-004 | Trie.Delete | algo/trie | remove a key, pruning branches that lead to no key |
| ALG-TRIE-005 | Trie.Len | algo/trie | the number of keys |
| ALG-TRIE-006 | Trie.WithPrefix | algo/trie | iterator over the keys with a prefix, in byte order |
| ALG-TRIE-007 | Trie.AutocompleteTopK | algo/trie | the k best keys with a prefix by a caller's rank, ties by byte order |
| ALG-TRIE-008 | Trie.WildcardSearch | algo/trie | iterator over the keys matching a pattern where '.' is any one byte |
| ALG-TRIE-009 | Trie.ShortestPrefixOf | algo/trie | the shortest key that is a prefix of a word |
| ALG-TRIE-010 | ReplaceWordsByRoot | algo/trie | each word replaced by its shortest root prefix |
| ALG-TRIE-011 | NewPrefixSuffixIndex | algo/trie | a permuterm index of words for prefix-and-suffix queries |
| ALG-TRIE-012 | PrefixSuffixIndex.Lookup | algo/trie | the indexes of the words with a prefix and a suffix, ascending |
| ALG-TRIE-013 | PrefixSuffixIndex.LookupLast | algo/trie | the largest such index |
| ALG-TWOPTR-001 | TwoSumSorted | algo/twopointer | indexes of two elements of a sorted slice with a given sum |
| ALG-TWOPTR-002 | ThreeSum | algo/twopointer | every distinct triple with a given sum |
| ALG-TWOPTR-003 | RemoveDuplicatesInPlace | algo/twopointer | drop adjacent duplicates, return the new length |
| ALG-TWOPTR-004 | MergeSortedInPlace | algo/twopointer | merge a sorted slice into a buffer's sorted prefix, back to front |
| ALG-TWOPTR-005 | IsPalindrome | algo/twopointer | whether a slice equals its reversal |
| ALG-TWOPTR-006 | MoveZeroesInPlace | algo/twopointer | move zero values to the end, keeping the others' order |
| ALG-UNION-001 | NewDisjointSet | algo/unionfind | n singleton sets {0} … {n-1} |
| ALG-UNION-002 | DisjointSet.MakeSet | algo/unionfind | add one new singleton element |
| ALG-UNION-003 | DisjointSet.Find | algo/unionfind | the representative of an element's set (with path compression) |
| ALG-UNION-004 | DisjointSet.Union | algo/unionfind | merge two sets (union by rank), reporting whether they were separate |
| ALG-UNION-005 | DisjointSet.Connected | algo/unionfind | whether two elements share a set |
| ALG-UNION-006 | DisjointSet.Count | algo/unionfind | the number of sets |
| ALG-UNION-007 | DisjointSet.Len | algo/unionfind | the number of elements |
| ALG-UNION-008 | IslandsOnline | algo/unionfind | the island count after each land addition on a grid |
| ALG-WINDOW-001 | RollingSum | algo/window | compensated sum of each window |
| ALG-WINDOW-002 | RollingMean | algo/window | mean of each window (the single moving average) |
| ALG-WINDOW-003 | RollingMin | algo/window | least element of each window |
| ALG-WINDOW-004 | RollingMax | algo/window | greatest element of each window |
| ALG-WINDOW-005 | NewMonotonicDeque | algo/window | an empty min- or max-deque |
| ALG-WINDOW-006 | MonotonicDeque.PushBack | algo/window | evict beaten elements from the back, then append, O(1) amortized |
| ALG-WINDOW-007 | MonotonicDeque.Front | algo/window | the least (or greatest) element held |
| ALG-WINDOW-008 | MonotonicDeque.PopFront | algo/window | remove the front element |
| ALG-WINDOW-009 | MonotonicDeque.PopFrontIfEqual | algo/window | remove the front element if it is the value leaving the window |
| ALG-WINDOW-010 | MonotonicDeque.Len | algo/window | the number of elements held |
| ALG-WINDOW-011 | MinLengthSubarraySumAtLeast | algo/window | shortest run of positive values with sum at least a target |
| ALG-WINDOW-012 | LongestWithAtMostKDistinct | algo/window | longest run with at most k distinct values |
| ALG-WINDOW-013 | LongestWithoutRepeat | algo/window | longest run with no repeated value |
| ALG-WINDOW-014 | ExponentialMovingAverage | algo/window | single exponential smoothing (NIST §6.4.3.1) |
| ALG-WINDOW-015 | MovingMedian | algo/window | median of each window (mean of the two middles when even), two heaps |
| ALG-WINDOW-016 | CoOccurrenceWithinWindow | algo/window | per ordered label pair, events of A with an event of B within the window |
| FIN-BIZ-001 | CurrentRatio | fin/business | current assets ÷ current liabilities |
| FIN-BIZ-002 | QuickRatio | fin/business | (cash + marketable securities + receivables) ÷ current liabilities |
| FIN-BIZ-003 | DebtRatio | fin/business | total liabilities ÷ total assets |
| FIN-BIZ-004 | DebtToEquityRatio | fin/business | total liabilities ÷ total equity |
| FIN-BIZ-005 | EquityMultiplier | fin/business | total assets ÷ total equity |
| FIN-BIZ-007 | TimesInterestEarned | fin/business | EBIT ÷ interest expense |
| FIN-BIZ-008 | FixedChargeCoverageRatio | fin/business | (EBIT + lease payments) ÷ (interest + lease payments) |
| FIN-BIZ-012 | AssetTurnover | fin/business | net sales ÷ average total assets |
| FIN-BIZ-014 | ReceivablesTurnover | fin/business | net credit sales ÷ average receivables |
| FIN-BIZ-015 | AverageCollectionPeriod | fin/business | days × average receivables ÷ net credit sales |
| FIN-BIZ-016 | InventoryTurnover | fin/business | cost of goods sold ÷ average inventory |
| FIN-BIZ-017 | DaysInInventory | fin/business | days × inventory ÷ cost of goods sold |
| FIN-BIZ-018 | GrossProfitMargin | fin/business | (net sales − cost of goods sold) ÷ net sales |
| FIN-BIZ-019 | OperatingProfitMargin | fin/business | operating income ÷ net sales |
| FIN-BIZ-020 | NetProfitMargin | fin/business | net income ÷ net sales |
| FIN-BIZ-021 | EarningsPerShare | fin/business | (net income − preferred dividends) ÷ weighted average shares, cents |
| FIN-BIZ-022 | DilutedEarningsPerShare | fin/business | (net income − preferred + convertible interest) ÷ diluted shares, cents |
| FIN-BIZ-023 | PriceToEarningsRatio | fin/business | price per share ÷ earnings per share |
| FIN-BIZ-024 | PriceToBookRatio | fin/business | price per share ÷ book value per share |
| FIN-BIZ-025 | PriceToSalesRatio | fin/business | market capitalization ÷ revenue |
| FIN-BIZ-026 | BookValuePerShare | fin/business | (total equity − preferred equity) ÷ shares, cents |
| FIN-BIZ-027 | DividendsPerShare | fin/business | total dividends ÷ shares outstanding, cents |
| FIN-BIZ-028 | DividendPayoutRatio | fin/business | dividends ÷ net income |
| FIN-BIZ-029 | RetentionRatio | fin/business | (net income − dividends) ÷ net income |
| FIN-BIZ-030 | ReturnOnAssets | fin/business | net income ÷ average total assets |
| FIN-BIZ-031 | ReturnOnEquity | fin/business | net income ÷ average equity |
| FIN-BIZ-032 | ReturnOnCommonEquity | fin/business | (net income − preferred dividends) ÷ average common equity |
| FIN-BIZ-034 | EconomicValueAdded | fin/business | after-tax income − invested capital × cost of capital, cents |
| FIN-BIZ-035 | FreeCashFlowToFirm | fin/business | EBIT(1 − t) + depreciation − capex − Δ working capital, cents |
| FIN-BIZ-036 | FreeCashFlowToEquity | fin/business | FCFF − interest(1 − t) + net borrowing, cents |
| FIN-BIZ-037 | NetWorkingCapital | fin/business | current assets − current liabilities, cents |
| FIN-BIZ-038 | DegreeOfOperatingLeverage | fin/business | contribution margin ÷ operating income |
| FIN-BIZ-040 | ContributionMargin | fin/business | sales − variable costs, cents |
| FIN-BIZ-041 | ContributionMarginRatio | fin/business | (sales − variable costs) ÷ sales |
| FIN-BIZ-042 | BreakEvenUnits | fin/business | fixed costs ÷ (price − variable cost) per unit |
| FIN-BIZ-043 | EconomicOrderQuantity | fin/business | √(2 × annual demand × order cost ÷ holding cost), units |
| FIN-BIZ-046 | DepreciationStraightLine | fin/business | straight-line schedule, remainder in the last period |
| FIN-BIZ-047 | DepreciationDecliningBalance | fin/business | declining-balance schedule, never below salvage |
| FIN-BIZ-048 | GoodwillFromAcquisition | fin/business | purchase price − fair value of net assets, cents |
| FIN-BIZ-049 | BankEfficiencyRatio | fin/business | noninterest expense ÷ (net interest income + noninterest income) |
| FIN-BIZ-050 | LoanToDepositRatio | fin/business | net loans and leases ÷ total deposits |
| FIN-ECON-001 | ConsumerPriceIndex | fin/economics | fixed-basket price index, base period = 100 |
| FIN-ECON-002 | InflationRate | fin/economics | decimal change between two index levels |
| FIN-ECON-003 | RealInterestRate | fin/economics | exact Fisher relation (1+i)/(1+π) − 1 |
| FIN-ECON-004 | HerfindahlHirschmanIndex | fin/economics | sum of squared percent market shares, 0–10 000 |
| FIN-INT-001 | FutureValue | fin/interest | PV·(1+i)ⁿ, cents |
| FIN-INT-002 | PresentValue | fin/interest | FV·(1+i)⁻ⁿ, cents |
| FIN-INT-003 | FutureValueContinuous | fin/interest | PV·e^(rt), cents |
| FIN-INT-004 | PresentValueContinuous | fin/interest | FV·e^(−rt), cents |
| FIN-INT-005 | SinglePaymentCompoundAmount | fin/interest | (F/P,i,n) |
| FIN-INT-006 | SinglePaymentPresentWorth | fin/interest | (P/F,i,n) |
| FIN-INT-007 | UniformSeriesCompoundAmount | fin/interest | (F/A,i,n) |
| FIN-INT-008 | UniformSeriesSinkingFund | fin/interest | (A/F,i,n) |
| FIN-INT-009 | UniformSeriesPresentWorth | fin/interest | (P/A,i,n) |
| FIN-INT-010 | CapitalRecovery | fin/interest | (A/P,i,n) |
| FIN-INT-011 | UniformGradientPresentWorth | fin/interest | (P/G,i,n) |
| FIN-INT-012 | UniformGradientFutureWorth | fin/interest | (F/G,i,n) |
| FIN-INT-013 | UniformGradientUniformSeries | fin/interest | (A/G,i,n) |
| FIN-INT-014 | FutureValueAnnuity | fin/interest | FV of an ordinary annuity, cents |
| FIN-INT-015 | FutureValueAnnuityDue | fin/interest | FV of an annuity due, cents |
| FIN-INT-016 | PresentValueAnnuity | fin/interest | PV of an ordinary annuity, cents |
| FIN-INT-017 | PresentValueAnnuityDue | fin/interest | PV of an annuity due, cents |
| FIN-INT-018 | AnnuityPaymentFromPresentValue | fin/interest | level end-of-period payment for a PV (loan payment) |
| FIN-INT-019 | AnnuityPaymentFromFutureValue | fin/interest | level end-of-period payment for an FV (sinking fund) |
| FIN-INT-020 | AnnuityDuePaymentFromPresentValue | fin/interest | level start-of-period payment for a PV |
| FIN-INT-021 | AnnuityDuePaymentFromFutureValue | fin/interest | level start-of-period payment for an FV |
| FIN-INT-022 | PresentValueGrowingAnnuity | fin/interest | PV of a growing ordinary annuity, cents |
| FIN-INT-023 | FutureValueGrowingAnnuity | fin/interest | FV of a growing ordinary annuity, cents |
| FIN-INT-024 | GrowingAnnuityPaymentFromPresentValue | fin/interest | first payment of a growing annuity for a PV |
| FIN-INT-025 | GrowingAnnuityPaymentFromFutureValue | fin/interest | first payment of a growing annuity for an FV |
| FIN-INT-026 | PresentValuePerpetuity | fin/interest | C/r, cents |
| FIN-INT-027 | PresentValueGrowingPerpetuity | fin/interest | C/(r−g), cents; also terminal value |
| FIN-INT-029 | EffectiveAnnualRate | fin/interest | (1 + r/m)^m − 1 |
| FIN-INT-030 | EffectiveAnnualRateContinuous | fin/interest | e^r − 1 |
| FIN-INT-031 | AnnualPercentageYield | fin/interest | Regulation DD Appendix A APY, decimal |
| FIN-INT-032 | AnnualPercentageRate | fin/interest | Regulation Z Appendix J APR, equal payments, decimal |
| FIN-INT-033 | RuleOf72 | fin/interest | approximate periods to double, 72/(100·i) |
| FIN-INT-034 | DoublingTime | fin/interest | exact periods to double, ln 2/ln(1+i) |
| FIN-INT-035 | SimpleInterest | fin/interest | P·r·t, cents |
| FIN-INT-036 | YearFraction30360 | fin/interest | MSRB Rule G-33 30/360 year fraction |
| FIN-INT-037 | YearFractionActual365 | fin/interest | actual days/365 year fraction |
| FIN-INT-038 | NetPresentValue | fin/interest | Σ cₜ/(1+i)ᵗ, cents |
| FIN-INT-039 | InternalRateOfReturn | fin/interest | rate where NPV = 0 |
| FIN-INT-040 | GrowthRatePerPeriod | fin/interest | (FV/PV)^(1/n) − 1 (CAGR) |
| FIN-INT-042 | AnnuityRate | fin/interest | rate solving PV = PMT·(P/A) + FV·(P/F) (bond yield) |
| FIN-INVEST-001 | CurrentYield | fin/invest | annual coupon / price |
| FIN-INVEST-002 | BondPrice | fin/invest | price of a level-coupon bond on a coupon date |
| FIN-INVEST-003 | YieldToMaturity | fin/invest | nominal annual yield to maturity from price |
| FIN-INVEST-004 | MacaulayDuration | fin/invest | PV-weighted average time of a bond's flows, years |
| FIN-INVEST-005 | ModifiedDuration | fin/invest | −(1/P)·dP/dy, years |
| FIN-INVEST-006 | Convexity | fin/invest | (1/P)·d²P/dy², years² |
| FIN-INVEST-007 | BankDiscountYield | fin/invest | Treasury bill discount rate, 360-day basis |
| FIN-INVEST-009 | BondEquivalentYield | fin/invest | Treasury bill investment rate (≤ half a year) |
| FIN-INVEST-010 | HoldingPeriodYield | fin/invest | total return over a holding period |
| FIN-INVEST-011 | AnnualizedHoldingPeriodYield | fin/invest | effective annual rate of a holding-period return |
| FIN-INVEST-012 | CapitalGainsYield | fin/invest | price change / beginning price |
| FIN-INVEST-013 | DividendYield | fin/invest | annual dividend / price |
| FIN-INVEST-014 | GordonGrowthPrice | fin/invest | D₁/(r − g) |
| FIN-INVEST-015 | ZeroGrowthStockPrice | fin/invest | D/r (also preferred stock) |
| FIN-INVEST-016 | CAPMExpectedReturn | fin/invest | Rf + β(Rm − Rf) |
| FIN-INVEST-017 | RiskPremium | fin/invest | expected return − risk-free rate |
| FIN-INVEST-018 | BlackScholesCall | fin/invest | European call value |
| FIN-INVEST-019 | BlackScholesPut | fin/invest | European put value |
| FIN-INVEST-020 | ModifiedInternalRateOfReturn | fin/invest | MIRR with finance and reinvestment rates |
| FIN-INVEST-021 | ProfitabilityIndex | fin/invest | PV(inflows)/PV(outflows) |
| FIN-INVEST-022 | PaybackPeriod | fin/invest | periods to recover the investment |
| FIN-INVEST-023 | DiscountedPaybackPeriod | fin/invest | payback on discounted flows |
| FIN-INVEST-025 | SharpeRatio | fin/invest | (Rp − Rf)/σp |
| FIN-INVEST-026 | WeightedAverageCostOfCapital | fin/invest | WACC with debt and equity |
| FIN-INVEST-027 | TaxEquivalentYield | fin/invest | tax-free yield/(1 − T) |
| FIN-INVEST-028 | NetAssetValuePerShare | fin/invest | (assets − liabilities)/shares |
| FIN-INVEST-030 | BidAskSpread | fin/invest | relative quoted spread |
| FIN-LOAN-001 | AmortizationSchedule | fin/loan | per-period payment, interest, principal and balance |
| FIN-LOAN-002 | RemainingBalance | fin/loan | balance after k payments, equal to schedule row k |
| FIN-LOAN-003 | BalloonBalance | fin/loan | balance after k payments of a given amount |
| FIN-LOAN-004 | BalloonLoanPayment | fin/loan | level payment that leaves a stated balloon |
| FIN-LOAN-005 | ExtraPaymentPayoff | fin/loan | periods to payoff and interest saved with an extra payment |
| FIN-LOAN-006 | LoanToValueRatio | fin/loan | loan / property value |
| FIN-LOAN-007 | PITIPayment | fin/loan | principal and interest + taxes + insurance (+ mortgage insurance) |
| FIN-LOAN-008 | DebtServiceCoverageRatio | fin/loan | net operating income / debt service |
| FIN-LOAN-009 | DebtToIncomeRatio | fin/loan | monthly debt payments / gross monthly income |
| FIN-LOAN-010 | AffordableLoanAmount | fin/loan | largest loan whose housing payment fits a housing ratio |
| LIF-FIT-001 | LoadProgressionACSM | life/fitness | ACSM 2009 2–10% load-increase range and whether its rep condition is met |
| LIF-FIT-002 | IPFGLCoefficient | life/fitness | IPF GL coefficient (2020) for body weight, sex and competition type |
| LIF-FIT-003 | IPFGLPoints | life/fitness | IPF GL points (2020) for a result, rounded per the IPF procedure |
| LIF-FIT-006 | RepsInReserve | life/fitness | repetitions in reserve from an RPE of 7–10 (10 − RPE) |
| LIF-FIT-007 | LoadAtPercentOneRepMax | life/fitness | training load at a percent of 1 RM |
| LIF-FIT-008 | VolumeLoad | life/fitness | one set's volume load, repetitions × external load |
| LIF-FIT-009 | Pace | life/fitness | pace, in seconds per unit of distance |
| LIF-FIT-010 | Splits | life/fitness | split times of a cumulative distance stream for whole units of unitKmOrMi |
| LIF-FIT-011 | TotalVolumeLoad | life/fitness | volume load of a session or a week |
| LIF-FIT-012 | OneRepMaxBrzycki | life/fitness | 1 RM from reps to fatigue, Brzycki: load / (1.0278 − 0.0278 × reps) |
| LIF-FIT-014 | OneRepMaxWathen | life/fitness | 1 RM from reps to fatigue, Wathen: 100 × load / (48.8 + 53.8 e^(−0.075 reps)) |
| LIF-NUTR-001 | RestingEnergyExpenditureMifflinStJeor | life/nutrition | REE (kcal/day), Mifflin 1990 equation with 166 × sex |
| LIF-NUTR-002 | RestingEnergyExpenditureMifflinStJeorSimplified | life/nutrition | REE (kcal/day), Mifflin 1990 sex-specific simplified forms |
| LIF-NUTR-003 | AMDRGrams | life/nutrition | protein, carbohydrate and fat g/day at the AMDR low and high ends |
| LIF-NUTR-004 | TotalWaterAI | life/nutrition | total-water and beverage AI (L/day), IOM 2005 Table S-2 |
| LIF-NUTR-005 | ProteinCarbohydrateDRIFor | life/nutrition | protein and carbohydrate RDA or AI (g/day), NASEM DRI summary table |
| LIF-NUTR-006 | PercentDailyValue | life/nutrition | percent Daily Value of a DRV nutrient, 21 CFR 101.9(d)(7)(ii) |
| LIF-NUTR-007 | ProteinRangeISSN | life/nutrition | ISSN 2017 protein range 1.4–2.0 g/kg/day (g/day) |
| MED-ANTHRO-001 | BodyMassIndex | med/anthropometry | body mass index (kg/m²) from kg and m |
| MED-ANTHRO-002 | BodyMassIndexImperial | med/anthropometry | body mass index (kg/m²) from lb and in, the CDC's × 703 form |
| MED-ANTHRO-003 | BMICategoryCDCAdult | med/anthropometry | CDC adult BMI category, the CDC's own label |
| MED-ANTHRO-004 | BodySurfaceAreaHaycock | med/anthropometry | body surface area (m²), Haycock 1978 |
| MED-ANTHRO-005 | AdjustedBodyWeight | med/anthropometry | ideal weight plus 40% of the excess (kg), Schwartz 1978 |
| MED-ANTHRO-006 | IdealBodyMassTraubKichen | med/anthropometry | ideal body mass of a child aged 1–17 (kg), Traub & Kichen 1983 |
| MED-ANTHRO-007 | EstimatedWeightInfantBestGuess | med/anthropometry | estimated weight of an infant under 12 months (kg), Best Guess |
| MED-ANTHRO-008 | BodyFatPercentDeurenbergAdult | med/anthropometry | body-fat percentage from BMI, age and sex, Deurenberg adult formula |
| MED-ANTHRO-009 | BodyFatPercentDeurenbergChild | med/anthropometry | body-fat percentage from BMI, age and sex, Deurenberg child formula |
| MED-ANTHRO-010 | BasalMetabolicRateHarrisBenedict | med/anthropometry | predicted basal heat production |
| MED-ANTHRO-011 | RelativeFatMass | med/anthropometry | relative fat mass (RFM), an estimate of whole-body fat percentage |
| MED-ANTHRO-012 | MUSTScore | med/anthropometry | 'Malnutrition Universal Screening Tool' ('MUST') total, 0–6 |
| MED-ANTHRO-013 | MUSTRiskCategory | med/anthropometry | 'MUST' step 4 overall risk of malnutrition for a total |
| MED-ANTHRO-014 | BodyRoundnessIndex | med/anthropometry | body roundness index from waist and height (cm), Thomas 2013 |
| MED-ANTHRO-015 | BodyRoundnessIndexImperial | med/anthropometry | body roundness index from waist and height (in) |
| MED-CARDIO-001 | MaxHeartRateTanaka | med/cardio | age-predicted maximal heart rate (Tanaka 2001) |
| MED-CARDIO-002 | QTcFridericia | med/cardio | QT interval corrected to RR = 1 s by the cube root (Fridericia 1920) |
| MED-CARDIO-003 | PulseWaveVelocityMoensKorteweg | med/cardio | pressure-wave speed in an elastic tube (Korteweg 1878, eq. 10) |
| MED-CARDIO-004 | TIMIScoreUANSTEMI | med/cardio | TIMI risk score for UA/NSTEMI, 0–7 (Antman 2000) |
| MED-CARDIO-005 | TIMIUANSTEMIEventRatePercent | med/cardio | 14-day event rate per TIMI UA/NSTEMI score (Antman 2000) |
| MED-CARDIO-006 | TIMIRiskIndex | med/cardio | HR × (age/10)² / SBP for STEMI (Morrow 2001) |
| MED-CARDIO-007 | MarburgHeartScore | med/cardio | chest-pain CAD score in primary care, 0–5 (Bösner 2010) |
| MED-CARDIO-008 | INTERCHESTScore | med/cardio | chest-pain CAD score in primary care, −1 to 5 (Aerts 2017) |
| MED-CARDIO-010 | OESILScore | med/cardio | syncope score, 0–4 (Colivicchi 2003) |
| MED-CARDIO-011 | OESILMortalityPercent | med/cardio | 12-month mortality per OESIL score (Colivicchi 2003) |
| MED-CARDIO-012 | RevisedCardiacRiskIndex | med/cardio | perioperative cardiac risk factors, 0–6 (Lee 1999) |
| MED-CARDIO-013 | SubtleAnteriorSTEMI4Variable | med/cardio | 4-variable LAD-occlusion formula (Driver 2017) |
| MED-CARDIO-014 | SubtleAnteriorSTEMI3Variable | med/cardio | 3-variable LAD-occlusion formula (Smith 2012) |
| MED-CARDIO-015 | ShanghaiBrugadaScore | med/cardio | Shanghai score for Brugada syndrome (Antzelevitch 2016) |
| MED-CARDIO-016 | TisdaleQTRiskScore | med/cardio | QTc-prolongation risk score, 0–21 (Tisdale 2013) |
| MED-CARDIO-017 | TisdaleQTRiskGroup | med/cardio | Tisdale Low / Moderate / High |
| MED-CARDIO-018 | LACEIndex | med/cardio | 30-day death or readmission index, 0–19 (van Walraven 2010) |
| MED-CARDIO-019 | CHADS2Score | med/cardio | CHADS₂, 0–6 (Gage 2001) |
| MED-CARDIO-020 | CHADS2StrokeRatePer100PersonYears | med/cardio | stroke rate per CHADS₂ score (Gage 2001) |
| MED-CARDIO-021 | ATRIABleedingScore | med/cardio | warfarin hemorrhage score, 0–10 (Fang 2011) |
| MED-CARDIO-022 | ATRIABleedingGroup | med/cardio | ATRIA bleeding low / intermediate / high risk |
| MED-CARDIO-023 | ATRIAStrokeScore | med/cardio | AF stroke score, 0–15 (Singer 2013) |
| MED-CARDIO-024 | ATRIAStrokeGroup | med/cardio | ATRIA stroke low / moderate / high |
| MED-CARDIO-025 | ORBITBleedingScore | med/cardio | AF bleeding score on OAC, 0–7 (O'Brien 2015) |
| MED-CARDIO-026 | ORBITBleedingGroup | med/cardio | ORBIT low / medium / high |
| MED-CARDIO-027 | DAPTScore | med/cardio | DAPT score, −2 to 10 (Yeh 2016) |
| MED-CARDIO-028 | HeartRateRecovery1Min | med/cardio | peak minus 1-minute heart rate (Cole 1999) |
| MED-CARDIO-029 | MaxHeartRateGulati | med/cardio | 206 − 0.88 × age, women (Gulati 2010) |
| MED-CARDIO-030 | MaxHeartRateNes | med/cardio | 211 − 0.64 × age (Nes 2013, HUNT) |
| MED-CARDIO-031 | CardiacPowerOutputW | med/cardio | MAP × CO / 451, watts (Fincke 2004) |
| MED-CARDIO-032 | DiastolicShockIndex | med/cardio | HR / diastolic pressure (Ospina-Tascón 2020) |
| MED-CARDIO-033 | RespiratoryAdjustedShockIndex | med/cardio | HR / SBP × RR / 10 (Caputo 2018) |
| MED-CARDIO-034 | ACEFII | med/cardio | ACEF II cardiac-surgery score (Ranucci 2018) |
| MED-CARDIO-035 | IEMortalityScore | med/cardio | 6-month mortality score in infective endocarditis (Park 2016) |
| MED-CARDIO-036 | GillmoreATTRStage | med/cardio | ATTR-CM stage I–III (Gillmore 2018) |
| MED-CARDIO-037 | SIPACutoff | med/cardio | age-adjusted pediatric shock-index cut-off (Acker 2015) |
| MED-CARDIO-038 | DutchLipidClinicFHScore | med/cardio | Dutch Lipid Clinic Network FH score, 0–26 (Nordestgaard 2013) |
| MED-CARDIO-039 | DutchLipidClinicFHGroup | med/cardio | definite / probable / possible / unlikely FH |
| MED-CARDIO-040 | WellsPEScore | med/cardio | Wells PE score, 0–12.5 (Wells 2000) |
| MED-CARDIO-041 | GenevaScoreRevised | med/cardio | revised Geneva PE score, 0–22 (Le Gal 2006) |
| MED-CARDIO-042 | DDimerAgeAdjustedCutoffUgL | med/cardio | age × 10 µg/L above 50 years (Douma 2010) |
| MED-CARDIO-043 | HERDOO2Count | med/cardio | HERDOO2 characteristics, 0–4 (Rodger 2008) |
| MED-CARDIO-044 | IMPEDEVTEScore | med/cardio | IMPEDE VTE score in myeloma (Sanfilippo 2019) |
| MED-CARDIO-045 | TRiPCastScore | med/cardio | TRiP(cast) VTE score after leg immobilisation (Nemeth 2020) |
| MED-CARDIO-046 | PulsePressureVariationPercent | med/cardio | respiratory variation of the arterial pulse pressure |
| MED-CARDIO-047 | CARTScore | med/cardio | Cardiac Arrest Risk Triage score, 0–57 (Churpek 2012) |
| MED-CARDIO-048 | CRUSADEBleedingScore | med/cardio | CRUSADE bleeding score, 1–96 (Subherwal 2009) |
| MED-CARDIO-049 | CRUSADEBleedingGroupOf | med/cardio | CRUSADE Very Low / Low / Moderate / High / Very High |
| MED-CARDIO-050 | READMITSScore | med/cardio | AMI READMITS 30-day readmission score, 0–33 (Nguyen 2018) |
| MED-CARDIO-051 | READMITSRiskGroupOf | med/cardio | READMITS extremely low … extremely high |
| MED-CARDIO-052 | RESCUEIHCAScore | med/cardio | RESCUE-IHCA score for ECPR in-hospital arrest (Tonna 2022) |
| MED-CARDIO-053 | ARCHBRHighBleedingRisk | med/cardio | ARC-HBR: ≥1 major or ≥2 minor criteria (Urban 2019) |
| MED-CARDIO-054 | DukeISCVIDCategoryOf | med/cardio | 2023 Duke-ISCVID Definite / Possible / Rejected IE (Fowler 2023) |
| MED-CARDIO-055 | BostonSyncopeRuleCount | med/cardio | Boston Syncope Rule groups present, 0–8 (Grossman 2007) |
| MED-CARDIO-056 | ModifiedShockIndex | med/cardio | heart rate / mean arterial pressure (Liu 2012) |
| MED-CARDIO-057 | HEARTScore | med/cardio | HEART score, 0–10, HEART Pathway bands (Mahler 2015) |
| MED-CARDIO-058 | HEARTPathwayLowRisk | med/cardio | HEART 0–3 and negative 0 h and 3 h troponin (Mahler 2015) |
| MED-CARDIO-059 | HEARScore | med/cardio | HEART without troponin, 0–8 (Smith 2020) |
| MED-CARDIO-060 | PCPHF10YearRiskPercent | med/cardio | 10-year heart-failure risk, PCP-HF equations (Khan 2019) |
| MED-CARDIO-062 | YEARSPulmonaryEmbolismExcluded | med/cardio | YEARS items + D-dimer 1000/500 ng/mL rule (van der Hulle 2017) |
| MED-CARDIO-063 | SanFranciscoSyncopeRuleCount | med/cardio | San Francisco Syncope Rule items present, 0–5 (Quinn 2004) |
| MED-CARDIO-064 | FourPEPSScore | med/cardio | 4PEPS score (−4 to 22), Roy 2021 Table 3 |
| MED-CARDIO-065 | FourPEPSLevel | med/cardio | 4PEPS clinical pretest probability level, the source's label |
| MED-CHEM-001 | CorrectedCalciumPayne | med/chemistry | albumin-adjusted calcium, mg/dL (Payne 1973) |
| MED-CHEM-002 | CorrectedCalciumPayneSI | med/chemistry | albumin-adjusted calcium, mmol/L from mmol/L and g/L |
| MED-CHEM-003 | SerumOsmolalityDorwart | med/chemistry | calculated serum osmolality, mOsm/kg, glucose and BUN in mg/dL |
| MED-CHEM-004 | SerumOsmolalityDorwartSI | med/chemistry | calculated serum osmolality, mOsm/kg, glucose and urea in mmol/L |
| MED-CHEM-005 | EstimatedAverageGlucose | med/chemistry | estimated average glucose from HbA1c, mg/dL (ADAG) |
| MED-CHEM-006 | EstimatedAverageGlucoseMmolL | med/chemistry | estimated average glucose from HbA1c, mmol/L (ADAG) |
| MED-CHEM-007 | AdjustedAnionGapFigge | med/chemistry | albumin-adjusted anion gap, albumin in g/L (Figge 1998) |
| MED-CHEM-008 | AdjustedAnionGapFiggeGDL | med/chemistry | albumin-adjusted anion gap, albumin in g/dL (Figge 1998) |
| MED-CHEM-009 | LDLCholesterolMartinHopkins | med/chemistry | LDL cholesterol, mg/dL, Martin–Hopkins 180-cell factor |
| MED-CHEM-010 | LDLCholesterolMartinHopkinsSI | med/chemistry | LDL cholesterol, mmol/L, Martin–Hopkins 180-cell factor |
| MED-CHEM-011 | AnionGap | med/chemistry | serum anion gap Na − (Cl + HCO3), mEq/L, Kitabchi 2009 |
| MED-CHEM-012 | EffectiveSerumOsmolality | med/chemistry | effective osmolality 2Na + glucose/18, mOsm/kg, Kitabchi 2009 |
| MED-CHEM-013 | CorrectedSodiumHyperglycemia | med/chemistry | Na + 1.6 per 100 mg/dL glucose above 100, mEq/L, Kitabchi 2009 |
| MED-CRIT-001 | NEWS2 | med/critical | NEWS2 parameter scores and aggregate |
| MED-CRIT-002 | NEWS2ClinicalRisk | med/critical | NEWS2 clinical-risk band for an aggregate score and the red-score flag |
| MED-CRIT-003 | SCORTEN | med/critical | SCORTEN (0–7), the count of its seven risk factors present |
| MED-CRIT-004 | ATLASScore | med/critical | five-component ATLAS score for Clostridium difficile infection, 0–10 |
| MED-CRIT-005 | MASCCRiskIndex | med/critical | MASCC risk-index score for febrile neutropenia |
| MED-CRIT-006 | NeutrophilLymphocyteRatio | med/critical | ratio of the neutrophil count to the lymphocyte count |
| MED-CRIT-007 | PECARNFebrileInfantLowRisk | med/critical | whether a febrile infant 60 days and younger meets the PECARN rule's low-risk |
| MED-CRIT-008 | QSOFAScore | med/critical | quick SOFA (qSOFA) score, 0–3 |
| MED-CRIT-009 | RochesterLowRisk | med/critical | whether a febrile infant 60 days of age or younger meets every Rochester |
| MED-CRIT-010 | ShapiroBloodCultureRule | med/critical | ShapiroBloodCultureRule applies the Shapiro decision rule for blood cultures in |
| MED-CRIT-011 | SICScore | med/critical | sepsis-induced coagulopathy (SIC) score of Iba 2017 Table 3 |
| MED-CRIT-012 | CPOTScore | med/critical | Critical-Care Pain Observation Tool score (0–8) |
| MED-CRIT-013 | SIRSCriteria | med/critical | the four SIRS criteria, their count and "two or more" |
| MED-CRIT-014 | SOFAScore | med/critical | SOFA organ scores and total, 0–24 |
| MED-CRIT-015 | MSOFAScore | med/critical | modified SOFA (Grissom 2010) components and total, 0–19 |
| MED-CRIT-016 | SOFA2Score | med/critical | SOFA-2 (Ranzani 2025) system scores and total, 0–24 |
| MED-CRIT-017 | PhoenixSepsisScore | med/critical | Phoenix Sepsis Score of a child and the consensus criteria |
| MED-CRIT-018 | REMSScore | med/critical | Rapid Emergency Medicine Score, 0–26 |
| MED-CRIT-019 | LRINECScore | med/critical | LRINEC score for necrotizing fasciitis, 0–13 |
| MED-CRIT-020 | SAPS3Score | med/critical | SAPS 3 admission score and hospital mortality, bilirubin and creatinine in mg/dL |
| MED-CRIT-021 | SAPS3ScoreSI | med/critical | SAPS 3 admission score and hospital mortality, bilirubin and creatinine in µmol/L |
| MED-EMERG-001 | TWISTScore | med/emergency | TWIST score (0–7) for testicular torsion in children |
| MED-EMERG-002 | NEDOCSScore | med/emergency | NEDOCS emergency-department crowding score, Weiss 2004 |
| MED-EMERG-003 | NEDOCSLevel | med/emergency | NEDOCS level, the source's label |
| MED-ENDO-001 | GlucoseManagementIndicatorPercent | med/endocrine | GMI (%) from mean CGM glucose in mg/dL, Bergenstal 2018 |
| MED-ENDO-002 | GlucoseManagementIndicatorMmolMol | med/endocrine | GMI (mmol/mol) from mean CGM glucose in mmol/L, Bergenstal 2018 |
| MED-ENDO-003 | BeAMValue | med/endocrine | bedtime minus pre-breakfast glucose (mg/dL), Zisman 2016 |
| MED-ENDO-004 | CPeptideGlucoseRatio | med/endocrine | C-peptide (ng/mL) / glucose (mg/dL) × 100, Saisho 2016 |
| MED-ENDO-005 | GDMThresholdIADPSGMgDL | med/endocrine | whether a 75-g OGTT value meets an IADPSG threshold (mg/dL) |
| MED-ENDO-006 | GDMThresholdIADPSGMmolL | med/endocrine | whether a 75-g OGTT value meets an IADPSG threshold (mmol/L) |
| MED-ENDO-007 | DKAMortalityScore | med/endocrine | DKA mortality prediction model points (0–25), Efstathiou 2002 |
| MED-ENDO-008 | AUSDRISKScore | med/endocrine | AUSDRISK points (0–35), Chen 2010 Box 3 |
| MED-ENDO-009 | DiabetesDistressScore17 | med/endocrine | DDS17 mean-item score (1–6) from the 17 ratings |
| MED-ENDO-010 | DiabetesDistressCategoryFisher | med/endocrine | Fisher 2012 DDS17 category, the source's own label |
| MED-ENDO-011 | TimeInRangesMgDL | med/endocrine | percent of glucose readings (mg/dL) in each consensus range |
| MED-ENDO-012 | TimeInRangesMmolL | med/endocrine | percent of glucose readings (mmol/L) in each consensus range |
| MED-ENDO-013 | A1CCategoryADAPercent | med/endocrine | ADA category of a hemoglobin A1C in percent (NGSP) |
| MED-ENDO-014 | A1CCategoryADAMmolMol | med/endocrine | ADA category of a hemoglobin A1C in mmol/mol |
| MED-ENDO-015 | FastingGlucoseCategoryADAMgDL | med/endocrine | ADA category of a fasting plasma glucose (FPG) in mg/dL |
| MED-ENDO-016 | FastingGlucoseCategoryADAMmolL | med/endocrine | ADA category of a fasting plasma glucose in mmol/L |
| MED-ENDO-017 | OGTTTwoHourCategoryADAMgDL | med/endocrine | ADA category of the 2-h plasma glucose |
| MED-ENDO-018 | OGTTTwoHourCategoryADAMmolL | med/endocrine | ADA category of the 2-h plasma glucose during a 75-g OGTT in mmol/L |
| MED-ENDO-019 | RandomGlucoseCriterionADAMgDL | med/endocrine | whether a random plasma glucose meets the fourth ADA Table 2.1 criterion |
| MED-ENDO-020 | RandomGlucoseCriterionADAMmolL | med/endocrine | whether a random plasma glucose in mmol/L meets the fourth ADA Table 2.1 |
| MED-ENDO-021 | HypoglycemiaLevelADAMgDL | med/endocrine | ADA hypoglycemia level of a glucose in mg/dL |
| MED-ENDO-022 | HypoglycemiaLevelADAMmolL | med/endocrine | ADA hypoglycemia level of a glucose in mmol/L |
| MED-ENDO-023 | CGMTargetsMetBattelino | med/endocrine | which of the 2019 International Consensus on Time in Range targets for people |
| MED-ENDO-024 | HOMAIR | med/endocrine | HOMA-IR, insulin (µU/mL) × glucose (mmol/L) / 22.5 |
| MED-ENDO-025 | HOMABeta | med/endocrine | HOMA-%B, 20 × insulin / (glucose mmol/L − 3.5), percent |
| MED-ENDO-026 | QUICKI | med/endocrine | 1/(log10 insulin µU/mL + log10 glucose mg/dL) |
| MED-ENDO-027 | METSIR | med/endocrine | METS-IR, ln(2G + TG) × BMI / ln(HDL), mg/dL |
| MED-ENDO-028 | MyxedemaComaScore | med/endocrine | myxedema coma diagnostic score (0–230), Popoveniuc 2014 |
| MED-ENDO-029 | OSTAIndex | med/endocrine | OSTA index, trunc(0.2 × (kg − years)), Koh 2001 |
| MED-ENDO-030 | OSTAIndexImperial | med/endocrine | OSTA index from pounds |
| MED-ENDO-031 | OSTARiskCategory | med/endocrine | OSTA risk category, the source's label |
| MED-GASTRO-001 | BISAPScore | med/gastro | BISAP (0–5), Wu 2008 |
| MED-GASTRO-002 | RockallScorePreEndoscopy | med/gastro | Rockall score before endoscopy (0–7), Rockall 1996 Table III |
| MED-GASTRO-003 | RockallScoreComplete | med/gastro | complete Rockall score (0–11), Rockall 1996 Table III |
| MED-GASTRO-004 | PediatricAppendicitisScore | med/gastro | Pediatric Appendicitis Score (0–10), Samuel 2002 |
| MED-GASTRO-005 | RIPASAScore | med/gastro | RIPASA score (1.5–16), Chong 2010 |
| MED-GASTRO-006 | ISEEScore | med/gastro | I-SEE total for eosinophilic esophagitis, Dellon 2022 Table 1 |
| MED-GASTRO-007 | ISEESeverity | med/gastro | I-SEE severity category, the source's own label |
| MED-GASTRO-008 | AcutePancreatitisDiagnosisScore | med/gastro | Jin acute pancreatitis model (−3 to 14), Jin 2024 Table 1 |
| MED-GASTRO-009 | RansonScore | med/gastro | Ranson's 11 early signs, non-gallstone (0–11) |
| MED-GASTRO-010 | GlasgowBlatchfordScore | med/gastro | Glasgow-Blatchford score (0–23), Blatchford 2000 via Banister 2018 Table 1 |
| MED-GASTRO-011 | AlvaradoScore | med/gastro | Alvarado score (0–10), Alvarado 1986 via Alvarado 2016 |
| MED-GASTRO-013 | OaklandScore | med/gastro | Oakland score for lower GI bleeding (0–33), Oakland 2017 Table 2 |
| MED-GASTRO-014 | PUCAIScore | med/gastro | Pediatric Ulcerative Colitis Activity Index (0–85), Turner 2007 |
| MED-HEMAT-001 | NCCNIPIScore | med/hematology | NCCN-IPI score (0–8) for diffuse large B-cell lymphoma |
| MED-HEMAT-002 | NCCNIPIRiskGroup | med/hematology | NCCN-IPI risk group for a score |
| MED-HEMAT-003 | GIPSSScore | med/hematology | GIPSS genetic score for primary myelofibrosis |
| MED-HEMAT-004 | GIPSSRiskCategory | med/hematology | GIPSS risk category for a score |
| MED-HEMAT-005 | MAPSScore | med/hematology | MAPS clinical-molecular score (0–6) for systemic mastocytosis |
| MED-HEMAT-006 | MALTIPIScore | med/hematology | MALT-IPI, the number of its three risk factors present |
| MED-HEMAT-007 | MALTIPIRiskGroup | med/hematology | MALT-IPI risk group for a score |
| MED-HEMAT-008 | MYSECPMScore | med/hematology | MYSEC-PM score for secondary myelofibrosis (points, with age × 0.15) |
| MED-HEMAT-009 | RIETEBleedingScore | med/hematology | RIETE major-bleeding score (0–8) during anticoagulation for VTE |
| MED-HEMAT-010 | RIETERiskGroup | med/hematology | RIETE risk group for a score |
| MED-HEMAT-011 | DuvalCIBMTRScore | med/hematology | Duval 2010 pre-transplant score (0–6) for AML in relapse or induction failure |
| MED-HEMAT-012 | ICCPrimaryMyelofibrosis | med/hematology | ICC 2022 PMF stage whose criteria are met: pre-PMF or overt PMF |
| MED-HEMAT-013 | ICCSystemicMastocytosis | med/hematology | whether the ICC 2022 systemic mastocytosis criteria are met |
| MED-HEMAT-016 | IMPROVEBleedingScore | med/hematology | IMPROVE bleeding risk score (0–30.5) for medical inpatients |
| MED-HEMAT-017 | CRPUpperReferenceLimitWenerMgDL | med/hematology | age- and sex-adjusted CRP upper reference limit, mg/dL (Wener 2000) |
| MED-HEMAT-019 | PLASMICScore | med/hematology | PLASMIC score (0–7), Bendapudi 2017 |
| MED-HEMAT-020 | PLASMICRiskGroup | med/hematology | PLASMIC risk group, the source's label |
| MED-HEPAT-001 | FIB4 | med/hepatic | FIB-4 index from age, AST, ALT and platelets |
| MED-HEPAT-002 | MELD | med/hepatic | MELD score (not rounded) |
| MED-HEPAT-003 | MELDNa | med/hepatic | MELD-Na score (not rounded) |
| MED-HEPAT-004 | MELD3 | med/hepatic | MELD 3.0 score, rounded to the nearest integer |
| MED-HEPAT-005 | BARDScore | med/hepatic | BARD score, 0–4, for nonalcoholic fatty liver disease |
| MED-HEPAT-006 | GlasgowAlcoholicHepatitisScore | med/hepatic | GAHS (5–12), Forrest 2005 Table 4, SI units |
| MED-HEPAT-007 | MaddreyDiscriminantFunction | med/hepatic | 4.6 × (PT − control PT) + bilirubin mg/dL |
| MED-HEPAT-008 | NAFLDActivityScore | med/hepatic | NAS (0–8), steatosis + lobular inflammation + ballooning |
| MED-INFECT-001 | CentorScore | med/infectious | Centor 1981 count (0–4) of its four clinical variables |
| MED-INFECT-002 | MenzaScore | med/infectious | Menza 2009 simple-model score (0–19) for HIV acquisition |
| MED-INFECT-003 | FACEDROPSScore | med/infectious | FACE DROPS score (−1 to 27) for Lyme facial palsy |
| MED-INFECT-004 | HIVRASPTotalRisk | med/infectious | RASP total risk of HIV transmission for one exposure (probability) |
| MED-NEURO-001 | ABC2HemorrhageVolumeML | med/neuro | intracerebral hemorrhage volume by ABC/2 (Kothari 1996) |
| MED-NEURO-002 | ABCD2Score | med/neuro | ABCD² score (0–7) after TIA (Johnston 2007) |
| MED-NEURO-003 | ABCD2RiskGroupOf | med/neuro | ABCD² stratum: low, moderate or high risk |
| MED-NEURO-004 | CPSSSScore | med/neuro | Cincinnati Prehospital Stroke Severity Scale (0–4) from NIHSS items |
| MED-NEURO-005 | DRAGONScore | med/neuro | DRAGON score (0–10) after IV thrombolysis |
| MED-NEURO-006 | ICHScore | med/neuro | ICH Score (0–6) (Hemphill 2001) |
| MED-NEURO-007 | SEDANScore | med/neuro | SEDAN score (0–6) for post-thrombolysis hemorrhage |
| MED-NEURO-008 | SICHScore | med/neuro | secondary intracerebral hemorrhage score (0–6) |
| MED-NEURO-009 | FASTEDScore | med/neuro | FAST-ED score (0–9) |
| MED-NEURO-010 | RACEScore | med/neuro | RACE scale (0–9) |
| MED-NEURO-011 | PREScore | med/neuro | Pittsburgh Response to Endovascular therapy score |
| MED-NEURO-012 | TwoHELPS2BScore | med/neuro | 2HELPS2B EEG seizure-risk score (0–7) |
| MED-NEURO-013 | GCSPupilsScore | med/neuro | GCS-Pupils score (1–15) |
| MED-NEURO-014 | TokuhashiRevisedGroup | med/neuro | revised Tokuhashi group for a total (0–15) |
| MED-NEURO-015 | FourATCategoryOf | med/neuro | 4AT band for a total (0–12) |
| MED-NEURO-016 | CanadianCTHeadRule | med/neuro | Canadian CT Head Rule: high- and medium-risk factor groups |
| MED-NEURO-017 | CATCHRule | med/neuro | CATCH rule: high- and medium-risk factor groups (children) |
| MED-NEURO-018 | NewOrleansCriteriaCount | med/neuro | New Orleans criteria present (0–7) |
| MED-NEURO-019 | NEXUSHeadCTLowRisk | med/neuro | NEXUS Head CT low-risk profile: under 65, none of seven findings |
| MED-NEURO-020 | OttawaSAHRuleCount | med/neuro | Ottawa SAH Rule findings present (0–6) |
| MED-NEURO-021 | UCDavisHeadTraumaLowRisk | med/neuro | UC Davis rule: none of five predictors (children) |
| MED-NEURO-022 | PECARNHeadInjuryVeryLowRisk | med/neuro | PECARN head-trauma very-low-risk rule for the child's age |
| MED-NEURO-023 | HINTSAnyFinding | med/neuro | any of the three HINTS findings present |
| MED-NEURO-024 | ASPECTSScore | med/neuro | ASPECTS (0–10): 10 minus regions with early ischemic change |
| MED-NEURO-025 | LAMSScore | med/neuro | Los Angeles Motor Scale (0–5) |
| MED-NEURO-026 | FOURScore | med/neuro | FOUR score (0–16) |
| MED-NEURO-027 | CAMICUPositive | med/neuro | CAM-ICU algorithm: features 1 and 2 plus 3 or 4 |
| MED-NEURO-028 | ModifiedFisherGrade | med/neuro | modified Fisher grade (0–4) from SAH thickness and IVH |
| MED-NEURO-030 | CanadianTIAScore | med/neuro | Canadian TIA Score (−3 to 23) |
| MED-NEURO-032 | AcuteStrokeDysphagiaScreenFailed | med/neuro | Acute Stroke Dysphagia Screen: any item present |
| MED-NEURO-033 | MTOQ4Score | med/neuro | mTOQ-4 sum (0–8) |
| MED-NEURO-034 | MTOQ4CategoryOf | med/neuro | mTOQ-4 treatment-efficacy category for a sum (0–8) |
| MED-NEURO-035 | MoCATotal | med/neuro | MoCA total (0–30) with the education point |
| MED-OB-001 | TermCategoryACOG | med/obstetric | ACOG CO 579 category: "early term", "full term", "late term" or "postterm" |
| MED-OB-002 | IOTASimpleRules | med/obstetric | IOTA simple rules: "malignant", "benign" or "inconclusive" from M and B feature counts |
| MED-ONCO-001 | ASTCTCRSGrade | med/oncology | ASTCT consensus grade (1–4) of cytokine release syndrome, 0 when the fever condition is unmet |
| MED-PEDS-001 | KocherCriteriaCount | med/pediatrics | number of Kocher's four predictors present (0–4) |
| MED-PEDS-002 | RecommendedSleepHoursAASM | med/pediatrics | AASM 2016 recommended sleep per 24 h for an age, as [min, max] hours |
| MED-PEDS-003 | DHAKADehydrationScore | med/pediatrics | DHAKA dehydration score (0–12) from four clinical signs |
| MED-PEDS-004 | DHAKADehydrationCategory | med/pediatrics | DHAKA dehydration category for a score |
| MED-PEDS-005 | WAT1Score | med/pediatrics | Withdrawal Assessment Tool-1 (0–12), Franck 2008 |
| MED-PEDS-006 | PEWSScore | med/pediatrics | Brighton Paediatric Early Warning Score (0–13), Monaghan 2005 |
| MED-PEDS-007 | BacterialMeningitisScore | med/pediatrics | Bacterial Meningitis Score, 0–6 (Nigrovic 2002; validated 2007) |
| MED-PEDS-008 | WHOWeightIncrement1MonthPercentiles | med/pediatrics | WHO 1-month weight-increment percentiles, birth–12 months, g (WHO 2009 data files) |
| MED-PSYCH-001 | COMMScore | med/psychiatry | Current Opioid Misuse Measure total (0–68) from its 17 item ratings |
| MED-PSYCH-002 | OpioidRiskToolCategory | med/psychiatry | Opioid Risk Tool band for a total: low, moderate or high risk |
| MED-PSYCH-003 | GAD7Score | med/psychiatry | GAD-7 anxiety severity score, 0–21 |
| MED-PSYCH-004 | GAD7SeverityOf | med/psychiatry | GAD-7 severity band for a total (0–21) |
| MED-PSYCH-005 | PCPTSD5Score | med/psychiatry | Primary Care PTSD Screen for DSM-5 total, 0–5 |
| MED-PSYCH-006 | ORTOUDScore | med/psychiatry | revised Opioid Risk Tool for opioid use disorder (ORT-OUD) total, 0–9 |
| MED-PSYCH-007 | CRAFFTScore | med/psychiatry | CRAFFT total, 0–6: the count of "yes" answers to its six items |
| MED-PSYCH-008 | AUDITCScore | med/psychiatry | AUDIT-C total, 0–12 |
| MED-PSYCH-009 | CAGEScore | med/psychiatry | CAGE total, 0–4: the count of "yes" answers |
| MED-PSYCH-010 | CoronavirusAnxietyScaleScore | med/psychiatry | Coronavirus Anxiety Scale total, 0–20 |
| MED-PSYCH-011 | DAST10Score | med/psychiatry | DAST-10 total, 0–10 (question 3 scores for "no") |
| MED-PSYCH-012 | DAST10LevelOf | med/psychiatry | DAST-10 level for a total (0–10) |
| MED-PSYCH-013 | ASRSScreenerCount | med/psychiatry | ASRS v1.1 screener items at threshold, 0–6 |
| MED-PSYCH-014 | MDQScreen | med/psychiatry | Mood Disorder Questionnaire symptom count and screening rule |
| MED-PSYCH-015 | GDS15Score | med/psychiatry | 15-item Geriatric Depression Scale total, 0–15 |
| MED-PSYCH-016 | COWSScore | med/psychiatry | Clinical Opiate Withdrawal Scale total, 0–48 |
| MED-PSYCH-017 | COWSSeverityOf | med/psychiatry | COWS withdrawal band for a total (0–48) |
| MED-PSYCH-018 | PSS3Positive | med/psychiatry | Patient Safety Screener (ED-SAFE) positive-screen rule |
| MED-RENAL-001 | EGFRCKDEPI2021 | med/renal | eGFR, 2021 CKD-EPI creatinine, creatinine in mg/dL |
| MED-RENAL-002 | EGFRCKDEPI2021SI | med/renal | eGFR, 2021 CKD-EPI creatinine, creatinine in µmol/L |
| MED-RENAL-003 | EGFRCKDEPI2021CystatinC | med/renal | eGFR, 2021 CKD-EPI creatinine–cystatin C, creatinine in mg/dL |
| MED-RENAL-004 | EGFRCKDEPI2021CystatinCSI | med/renal | eGFR, 2021 CKD-EPI creatinine–cystatin C, creatinine in µmol/L |
| MED-RENAL-005 | CreatinineClearanceCockcroftGault | med/renal | creatinine clearance (mL/min), Cockcroft–Gault, creatinine in mg/dL |
| MED-RENAL-006 | CreatinineClearanceCockcroftGaultSI | med/renal | creatinine clearance (mL/min), Cockcroft–Gault, creatinine in µmol/L |
| MED-RENAL-007 | EGFRSchwartzBedside | med/renal | eGFR of a child, bedside Schwartz, creatinine in mg/dL |
| MED-RENAL-008 | EGFRSchwartzBedsideSI | med/renal | eGFR of a child, bedside Schwartz, creatinine in µmol/L |
| MED-RENAL-009 | FractionalExcretionSodium | med/renal | fractional excretion of sodium (FENa), percent |
| MED-RENAL-010 | FractionalExcretionUrea | med/renal | fractional excretion of urea (FEUrea), percent |
| MED-RENAL-011 | GFRCategoryKDIGO | med/renal | KDIGO GFR category G1–G5 with the source's term |
| MED-RENAL-012 | AKINStage | med/renal | Acute Kidney Injury Network stage, 0–3 (0 |
| MED-RENAL-013 | STONEScore | med/renal | STONE score, 0–13, from Moore 2014 Table 3 |
| MED-RENAL-014 | STONEProbabilityGroup | med/renal | STONE-score group of Moore 2014 |
| MED-RENAL-015 | UrineProteinCreatinineRatio | med/renal | protein/creatinine ratio of a single voided urine sample |
| MED-RENAL-016 | UrineAnionGap | med/renal | urine anion gap, Na⁺ + K⁺ − Cl⁻, in mmol/L, |
| MED-RENAL-017 | EGFRCKiDU25Creatinine | med/renal | eGFR of a child or young adult, CKiD U25 height/creatinine, creatinine in mg/dL |
| MED-RENAL-018 | EGFRCKiDU25CreatinineSI | med/renal | eGFR, CKiD U25 height/creatinine, creatinine in µmol/L |
| MED-RENAL-019 | EGFRCKiDU25CystatinC | med/renal | eGFR of a child or young adult, CKiD U25 cystatin C |
| MED-RENAL-020 | McMahonScore | med/renal | McMahon rhabdomyolysis risk score, 0–19 |
| MED-RENAL-021 | HIVCKDRiskScore | med/renal | Scherzer 2014 5-year CKD risk score in HIV, 0–15 |
| MED-RENAL-022 | EGFRMDRD | med/renal | eGFR, 4-variable MDRD (IDMS, 175), creatinine in mg/dL |
| MED-RENAL-023 | EGFRMDRDSI | med/renal | eGFR, 4-variable MDRD (IDMS, 175), creatinine in µmol/L |
| MED-RENAL-024 | KtV | med/renal | dialysis dose Kt/V from clearance, time and urea volume |
| MED-RENAL-025 | RIFLEClassify | med/renal | ADQI RIFLE class from creatinine and urine output, creatinine in mg/dL |
| MED-RENAL-026 | RIFLEClassifySI | med/renal | ADQI RIFLE class, creatinine in µmol/L |
| MED-RENAL-027 | ThakarScore | med/renal | Thakar (Cleveland Clinic) acute renal failure after cardiac surgery score (0–17) |
| MED-RESP-001 | DiffusingCapacityRoughtonForster | med/respiratory | lung diffusing capacity from DM, θ and VC (Roughton–Forster) |
| MED-RESP-002 | SixMinuteWalkDistanceEnright | med/respiratory | reference six-minute walk distance |
| MED-RESP-003 | SixMinuteWalkDistanceEnrightImperial | med/respiratory | reference six-minute walk distance |
| MED-RESP-004 | AIRQScore | med/respiratory | Asthma Impairment and Risk Questionnaire total (0–10) |
| MED-RESP-005 | AIRQControlLevel | med/respiratory | AIRQ control level for a total score |
| MED-RESP-006 | DeltaPScore | med/respiratory | Dutch-English LEMS Tumor Association Prediction (DELTA-P) score (0–6) |
| MED-RESP-007 | PackYears | med/respiratory | cumulative smoking exposure in pack-years |
| MED-RESP-008 | PackYearsFromCigarettes | med/respiratory | pack-years from cigarettes per day |
| MED-RESP-009 | NEWS2SpO2Scale1Points | med/respiratory | NEWS2 points (0–3) for an oxygen saturation on SpO₂ Scale 1 |
| MED-RESP-010 | NEWS2SpO2Scale2Points | med/respiratory | NEWS2 points (0–3) for an oxygen saturation on SpO₂ Scale 2 |
| MED-RESP-011 | CURB65Score | med/respiratory | CURB-65 score (0–5) for community-acquired pneumonia |
| MED-RESP-012 | ROXIndex | med/respiratory | ROX index: (SpO₂ / FiO₂) / respiratory rate |
| MED-RESP-013 | PaO2FiO2Ratio | med/respiratory | PaO₂/FiO₂ ratio, in mmHg: arterial oxygen tension divided |
| MED-RESP-014 | BerlinARDSOxygenationCategory | med/respiratory | Berlin Definition's oxygenation category for a PaO₂/FiO₂ ratio |
| MED-RESP-015 | SpO2FiO2Ratio | med/respiratory | SpO₂/FiO₂ (S/F) ratio: oxygen saturation by pulse oximetry |
| MED-RESP-016 | SpO2FiO2FromPaO2FiO2Rice | med/respiratory | S/F ratio that the Rice 2007 regression S/F = 64 + 0.84 × |
| MED-RESP-017 | RapidShallowBreathingIndex | med/respiratory | rapid shallow breathing index f/VT, in breaths/min/L |
| MED-RESP-018 | OxygenationIndex | med/respiratory | oxygenation index OI = mean airway pressure × FiO₂ / PaO₂ |
| MED-RESP-019 | OxygenSaturationIndex | med/respiratory | oxygen saturation index OSI = mean airway pressure × FiO₂ / SpO₂ |
| MED-RESP-020 | MuLBSTAScore | med/respiratory | MuLBSTA score (0–20) for viral pneumonia |
| MED-RESP-021 | PESARDScore | med/respiratory | PE-SARD bleeding score (0–5 points, in steps of 0.5) |
| MED-RESP-022 | PESARDRiskGroup | med/respiratory | PE-SARD bleeding-risk group for a score: 0 low, 1–2.5 intermediate, > 2.5 high |
| MED-RESP-023 | PIROScore | med/respiratory | PIRO score for community-acquired pneumonia (0–8) |
| MED-RESP-024 | PIRORiskLevel | med/respiratory | PIRO risk level for a score: 0–2 Low, 3 Mild, 4 high, 5–8 Very high |
| MED-RESP-025 | SCARFScore | med/respiratory | Sequential Clinical Assessment of Respiratory Function score |
| MED-RESP-026 | PeakFlowZoneNHLBI | med/respiratory | NHLBI Asthma Action Plan zone for a peak expiratory flow measured against the |
| MED-RESP-027 | CHOSENScore | med/respiratory | CHOSEN COVID-19 risk score (0–55) from age, SpO₂ and albumin |
| MED-RESP-028 | PESIScore | med/respiratory | Pulmonary Embolism Severity Index total (age + item points) |
| MED-RESP-029 | PESIRiskClass | med/respiratory | PESI risk class I–V for a total |
| MED-RESP-030 | OttawaCOPDRiskScale | med/respiratory | Ottawa COPD Risk Scale total (0–16) |
| MED-RESP-031 | OttawaCOPDRiskCategory | med/respiratory | Ottawa COPD risk category: 0 Low, 1–2 Medium, 3–4 High, ≥ 5 Very high |
| MED-RESP-032 | DRIPScore | med/respiratory | Drug Resistance in Pneumonia score (0–14) |
| MED-RESP-033 | FeverPAINScore | med/respiratory | FeverPAIN score (0–5) |
| MED-RESP-034 | PRESETScore | med/respiratory | PRESET ECMO score (0–15) |
| MED-RESP-035 | PRESETRiskClass | med/respiratory | PRESET risk class: 0–5 I, 6–9 II, 10–15 III |
| MED-RESP-036 | QCSIScore | med/respiratory | quick COVID-19 Severity Index (0–12) |
| MED-RESP-037 | POMPECMortalityPercent | med/respiratory | POMPE-C 30-day mortality, percent (weight in lb) |
| MED-RESP-038 | POMPECMortalityPercentKg | med/respiratory | POMPE-C 30-day mortality, percent (weight in kg) |
| MED-RESP-039 | RESECT90MortalityProbability | med/respiratory | RESECT-90 probability of death within 90 days of lung resection |
| MED-RESP-040 | STUMBLScore | med/respiratory | blunt chest wall trauma (STUMBL) risk score |
| MED-RESP-041 | UtahCOVID19Score | med/respiratory | Utah COVID-19 risk score (0.5–24.5) |
| MED-RESP-042 | PARSScore | med/respiratory | Pediatric Asthma Risk Score (0–14) |
| MED-RESP-043 | LIPSScore | med/respiratory | Lung Injury Prediction Score |
| MED-RESP-045 | CRB65Score | med/respiratory | CRB-65 score (0–4) for community-acquired pneumonia |
| MED-RESP-046 | BODEIndex | med/respiratory | BODE index (0–10) for COPD |
| MED-RESP-047 | LightsCriteria | med/respiratory | Light's criteria met (0–3) and whether a pleural effusion is an exudate |
| MED-RESP-048 | COPDAssessmentTestScore | med/respiratory | COPD Assessment Test total (0–40) |
| MED-RESP-049 | SilvermanAndersenScore | med/respiratory | Silverman–Andersen respiratory severity score (0–10) |
| MED-RESP-050 | PEFAmplitudePercentMean | med/respiratory | one day's PEF amplitude percent mean (highest − lowest) / mean × 100, GINA |
| MED-RESP-051 | PEFDiurnalVariabilityPercent | med/respiratory | average daily diurnal PEF variability over the days given, GINA |
| MED-RESP-052 | PEFDiurnalVariabilityExcessiveGINA | med/respiratory | whether diurnal PEF variability exceeds GINA Box 1-2 (> 10 % adults, > 13 % children 6–11) |
| MED-RESP-053 | PEFChangeConsistentWithAsthmaGINA | med/respiratory | whether a PEF change is at least 20 % (GINA, adults) |
| MED-RESP-054 | AcuteAsthmaPEFBandSIGN | med/respiratory | SIGN 158 Table 15 PEF cell (< 33 / 33–50 / > 50–75 % best or predicted), adults |
| MED-RHEUM-001 | SjogrenACREULAR2016Score | med/rheumatology | 2016 ACR/EULAR primary Sjögren's syndrome weighted score (0–9) |
| MED-RHEUM-002 | PMRScoreWithoutUltrasound | med/rheumatology | 2012 EULAR/ACR PMR algorithm score without ultrasound (0–6) |
| MED-RHEUM-003 | PMRScoreWithUltrasound | med/rheumatology | 2012 EULAR/ACR PMR algorithm score with ultrasound (0–8) |
| MED-RHEUM-004 | GoutACREULAR2015Score | med/rheumatology | 2015 ACR/EULAR gout criteria domain score (−6 to 23) |
| MED-RHEUM-005 | SLICC2012Classified | med/rheumatology | whether the SLICC 2012 SLE classification rule is met |
| MED-RHEUM-006 | AcuteGoutDiagnosticRuleScore | med/rheumatology | Janssens 2010 acute gout rule score (0–13) |
| MED-RHEUM-009 | PASScore | med/rheumatology | PAS / PAS-II (0–10), Wolfe 2005 |
| MED-RHEUM-010 | PASActivityCategory | med/rheumatology | PAS / PAS-II activity category, the source's label |
| MED-SCALE-001 | EpworthSleepinessScore | med/scales | Epworth Sleepiness Scale total (0–24) from the 8 item ratings |
| MED-SCALE-002 | EpworthCategory | med/scales | Epworth band for a total, the official site's own label |
| MED-SCALE-003 | GlasgowComaScore | med/scales | Glasgow Coma Scale total (3–15) from eye, verbal and motor ratings |
| MED-SCALE-004 | GlasgowComaSeverityRimel | med/scales | Rimel 1982 head-injury group for a GCS total: severe, moderate, minor |
| MED-SCALE-005 | ApgarScore | med/scales | Apgar score (0–10) from the five signs, each rated 0–2 |
| MED-SCALE-006 | ApgarCondition1953 | med/scales | Apgar 1953 group for a score: good, fair or poor condition |
| MED-SCALE-007 | HARKScore | med/scales | HARK score (0–4): one point per yes to the four HARK questions |
| MED-SCALE-008 | DangerAssessment5Count | med/scales | count (0–5) of yes answers to the five DA-5 questions |
| MED-SCALE-011 | NeckDisabilityIndex | med/scales | Neck Disability Index (0–50), Vernon & Mior 1991 |
| MED-SCALE-012 | NDIDisability | med/scales | Neck Disability Index band, the source's label |
| MED-SCALE-013 | WASTScore | med/scales | Woman Abuse Screening Tool sum (8–24), Brown 1996 |
| MED-SCALE-014 | FLACCScore | med/scales | FLACC behavioural pain score (0–10), arithmetic only |
| MED-SCALE-015 | PAINADScore | med/scales | PAINAD pain score (0–10), arithmetic only |
| MED-THERMO-001 | HypothermiaStageSwiss | med/thermoregulation | accidental-hypothermia stage I–III for a core temperature in °C |
| MED-THERMO-002 | HypothermiaStageSwissF | med/thermoregulation | same stage for a core temperature in °F |
| MED-TRAUMA-001 | ABCScore | med/trauma | ABC score (0–4) for massive transfusion |
| MED-TRAUMA-002 | CanadianCSpineRule | med/trauma | Canadian C-Spine Rule's three questions and outcome |
| MED-TRAUMA-003 | NEXUSChest | med/trauma | count of the seven NEXUS Chest criteria present |
| MED-TRAUMA-004 | NEXUSChestCT | med/trauma | NEXUS Chest CT-All and Chest CT-Major: any criterion present |
| MED-TRAUMA-005 | NEXUSCSpineLowProbability | med/trauma | NEXUS cervical spine: all five low-probability criteria met |
| MED-TRAUMA-006 | NewTraumaScore | med/trauma | New Trauma Score, triage sum and probability of survival |
| MED-TRAUMA-007 | OttawaAnkleRule1992 | med/trauma | 1992 Ottawa ankle malleolar and midfoot finding sets |
| MED-TRAUMA-008 | OttawaKneeRule | med/trauma | count of the five Ottawa knee rule variables present |
| MED-TRAUMA-009 | PECARNIntraAbdominalInjuryVeryLowRisk | med/trauma | PECARN blunt abdominal rule: every item met |
| MED-TRAUMA-010 | PedSRCVeryLowRisk | med/trauma | PedSRC rule: count of its five variables and none present |
| MED-TRAUMA-011 | PittsburghKneeRule | med/trauma | Pittsburgh knee rule combination |
| MED-TRAUMA-012 | RibScore | med/trauma | RibScore (0–6) from six radiographic variables |
| MED-TRAUMA-013 | TASHScore | med/trauma | TASH score (0–28) |
| MED-TRAUMA-014 | WisconsinCriteria | med/trauma | count of the five Wisconsin maxillofacial CT criteria present |
| MED-TRAUMA-015 | TLICS | med/trauma | Thoracolumbar Injury Classification and Severity Score, 1–10 |
| MED-TRAUMA-016 | InjurySeverityScore | med/trauma | Injury Severity Score from six AIS region maxima, 0–75 |
| MED-TRAUMA-017 | SINSScore | med/trauma | Spinal Instability Neoplastic Score, 0–18 |
| MED-TRAUMA-018 | PediatricTraumaScore | med/trauma | Pediatric Trauma Score, −6 to 12 |
| MED-TRAUMA-019 | MESSScore | med/trauma | Mangled Extremity Severity Score |
| MED-TRAUMA-020 | RevisedTraumaScore | med/trauma | Revised Trauma Score: GCS/SBP/RR codes, T-RTS 0–12 and weighted RTS |
| MED-TRAUMA-021 | NottinghamHipFractureScore | med/trauma | Nottingham Hip Fracture Score (0–10), haemoglobin in g/dL |
| MED-TRAUMA-022 | NottinghamHipFractureScoreSI | med/trauma | Nottingham Hip Fracture Score (0–10), haemoglobin in g/L |
| MED-TRAUMA-023 | PECARNCSpineRule | med/trauma | PECARN cervical spine rule group: high-risk, CART-identified or neither factor |
| NUM-ALGEB-001 | PolynomialEvaluate | num/algebra | value of a polynomial at x by Horner's rule |
| NUM-ALGEB-002 | PolynomialValueAndDerivative | num/algebra | value and first derivative of a polynomial at x in one pass |
| NUM-ALGEB-003 | PolynomialMultiply | num/algebra | coefficients of the product of two polynomials |
| NUM-ALGEB-004 | LinearEvaluate | num/algebra | y = m·x + b with one rounding |
| NUM-ALGEB-005 | LinearSolve | num/algebra | x from y = m·x + b |
| NUM-ALGEB-006 | QuadraticDiscriminant | num/algebra | b² − 4ac to within 2 ulps |
| NUM-ALGEB-007 | QuadraticRealRoots | num/algebra | real roots of a·x² + b·x + c without cancellation |
| NUM-ALGEB-008 | QuadraticRoots | num/algebra | both complex roots of a·x² + b·x + c |
| NUM-ALGEB-009 | CubicRoots | num/algebra | the three complex roots of a cubic, Newton-polished |
| NUM-ALGEB-010 | PolynomialRoots | num/algebra | all complex roots of a real polynomial (Aberth–Ehrlich) |
| NUM-COMBI-001 | Factorial | num/combinatorics | n! exactly, n ≤ 20 |
| NUM-COMBI-002 | LogFactorial | num/combinatorics | natural log of n! |
| NUM-COMBI-003 | Binomial | num/combinatorics | C(n,k) exactly, overflow detected |
| NUM-COMBI-004 | LogBinomial | num/combinatorics | natural log of C(n,k) |
| NUM-COMBI-005 | Permutations | num/combinatorics | nPr = n!/(n−k)! exactly |
| NUM-COMBI-006 | PermutationsWithRepetition | num/combinatorics | nᵏ sequences of length k over n symbols |
| NUM-COMBI-007 | CombinationsWithRepetition | num/combinatorics | multisets of size k from n kinds, C(n+k−1,k) |
| NUM-COMBI-008 | Multinomial | num/combinatorics | (Σcᵢ)!/Πcᵢ!, permutations of a multiset |
| NUM-COMBI-009 | CircularPermutations | num/combinatorics | arrangements of k of n around a circle, rotations equal |
| NUM-COMBI-010 | CircularPermutationsUndirected | num/combinatorics | circular arrangements, rotations and reflections equal |
| NUM-COMBI-011 | Derangements | num/combinatorics | permutations of n with no fixed point, n ≤ 20 |
| NUM-COMBI-012 | Catalan | num/combinatorics | n-th Catalan number, n ≤ 36 |
| NUM-COMBI-013 | StirlingFirstUnsigned | num/combinatorics | permutations of n with exactly k cycles |
| NUM-COMBI-014 | StirlingSecond | num/combinatorics | partitions of an n-set into k blocks |
| NUM-COMBI-015 | Bell | num/combinatorics | number of partitions of an n-set, n ≤ 25 |
| NUM-COMBI-016 | StirlingApproximation | num/combinatorics | √(2πn)(n/e)ⁿ, with Robbins's error bound |
| NUM-COMBI-017 | MultinomialBig | num/combinatorics | Multinomial as a *big.Int |
| NUM-COMBI-018 | DerangementsBig | num/combinatorics | Derangements as a *big.Int |
| NUM-COMBI-019 | CatalanBig | num/combinatorics | Catalan as a *big.Int |
| NUM-COMBI-020 | StirlingFirstUnsignedBig | num/combinatorics | StirlingFirstUnsigned as a *big.Int |
| NUM-COMBI-021 | StirlingSecondBig | num/combinatorics | StirlingSecond as a *big.Int |
| NUM-COMBI-022 | BellBig | num/combinatorics | Bell as a *big.Int |
| NUM-GEOM-001 | PointDistance | num/geometry | distance between two points |
| NUM-GEOM-002 | Midpoint | num/geometry | midpoint of a segment |
| NUM-GEOM-003 | Slope | num/geometry | rise over run of the line through two points |
| NUM-GEOM-004 | SlopeAngleRad | num/geometry | angle of a direction from rise and run |
| NUM-GEOM-005 | PointLineDistance | num/geometry | perpendicular distance from a point to a line |
| NUM-GEOM-006 | AngleBetweenLinesRad | num/geometry | acute angle between two lines |
| NUM-GEOM-007 | PolarToCartesian | num/geometry | (r, θ) to (x, y) |
| NUM-GEOM-008 | CartesianToPolar | num/geometry | (x, y) to (r, θ) |
| NUM-GEOM-009 | ComplementaryAngleRad | num/geometry | π/2 − θ |
| NUM-GEOM-010 | SupplementaryAngleRad | num/geometry | π − θ |
| NUM-GEOM-011 | PolygonInteriorAngleSumRad | num/geometry | sum of the interior angles of an n-gon |
| NUM-GEOM-012 | RegularPolygonInteriorAngleRad | num/geometry | one interior angle of a regular n-gon |
| NUM-GEOM-013 | RegularPolygonExteriorAngleRad | num/geometry | one exterior angle of a regular n-gon |
| NUM-GEOM-014 | TrianglePerimeter | num/geometry | perimeter of a triangle from three sides |
| NUM-GEOM-015 | TriangleAreaBaseHeight | num/geometry | area of a triangle from base and height |
| NUM-GEOM-016 | TriangleAreaHeron | num/geometry | area of a triangle from three sides (Kahan's Heron) |
| NUM-GEOM-017 | TriangleAreaSAS | num/geometry | area of a triangle from two sides and the included angle |
| NUM-GEOM-020 | RightTriangleHypotenuse | num/geometry | hypotenuse of a right triangle from two legs |
| NUM-GEOM-021 | RightTriangleLeg | num/geometry | other leg of a right triangle from hypotenuse and leg |
| NUM-GEOM-022 | RectangleArea | num/geometry | area of a rectangle (square; parallelogram by base × height) |
| NUM-GEOM-023 | RectanglePerimeter | num/geometry | perimeter of a quadrilateral with two pairs of equal sides |
| NUM-GEOM-024 | ParallelogramAreaSAS | num/geometry | area of a parallelogram from two sides and the angle |
| NUM-GEOM-025 | TrapezoidArea | num/geometry | area of a trapezoid |
| NUM-GEOM-026 | KiteArea | num/geometry | area of a quadrilateral with perpendicular diagonals |
| NUM-GEOM-027 | RegularPolygonAreaFromSide | num/geometry | area of a regular n-gon from its side |
| NUM-GEOM-028 | RegularPolygonAreaFromApothem | num/geometry | area of a regular n-gon from its apothem |
| NUM-GEOM-029 | RegularPolygonAreaFromCircumradius | num/geometry | area of a regular n-gon from its circumradius |
| NUM-GEOM-030 | RegularPolygonSideInscribed | num/geometry | side of a regular n-gon inscribed in a circle |
| NUM-GEOM-031 | RegularPolygonSideCircumscribed | num/geometry | side of a regular n-gon circumscribed about a circle |
| NUM-GEOM-032 | PolygonArea | num/geometry | area of a simple polygon from its vertices (shoelace) |
| NUM-GEOM-033 | PolygonPerimeter | num/geometry | perimeter of a polygon from its vertices |
| NUM-GEOM-034 | CircleCircumference | num/geometry | circumference of a circle |
| NUM-GEOM-035 | CircleArea | num/geometry | area of a circle |
| NUM-GEOM-036 | AnnulusArea | num/geometry | area of an annulus (ring) |
| NUM-GEOM-037 | ArcLength | num/geometry | length of a circular arc |
| NUM-GEOM-038 | ChordLength | num/geometry | length of a circular chord |
| NUM-GEOM-039 | SectorArea | num/geometry | area of a circular sector |
| NUM-GEOM-040 | SegmentArea | num/geometry | area of a circular segment |
| NUM-GEOM-041 | ArcRadiusFromChordSagitta | num/geometry | radius of an arc from its chord and sagitta |
| NUM-GEOM-042 | CircleTangentLength | num/geometry | tangent length from an external point to a circle |
| NUM-GEOM-043 | EllipseArea | num/geometry | area of an ellipse |
| NUM-GEOM-044 | EllipseChordFromCenter | num/geometry | chord of an ellipse at an offset from the centre |
| NUM-GEOM-045 | EllipseChordFromEdge | num/geometry | chord of an ellipse at a depth from the vertex |
| NUM-GEOM-046 | EllipsePerimeter | num/geometry | perimeter of an ellipse, exact (AGM) |
| NUM-GEOM-047 | EllipsePerimeterRamanujan | num/geometry | perimeter of an ellipse, Ramanujan's approximation |
| NUM-GEOM-048 | PrismVolume | num/geometry | volume of a prism from base area and height |
| NUM-GEOM-049 | CylinderVolume | num/geometry | volume of a cylinder |
| NUM-GEOM-050 | CylinderLateralArea | num/geometry | lateral area of a cylinder |
| NUM-GEOM-051 | CylinderSurfaceArea | num/geometry | total surface area of a cylinder |
| NUM-GEOM-052 | HollowCylinderVolume | num/geometry | volume of a hollow cylinder (tube) |
| NUM-GEOM-053 | HollowCylinderSurfaceArea | num/geometry | total surface area of a hollow cylinder (tube) |
| NUM-GEOM-054 | PyramidVolume | num/geometry | volume of a pyramid, any base |
| NUM-GEOM-055 | RegularPyramidSurfaceArea | num/geometry | total surface area of a right pyramid with a regular base |
| NUM-GEOM-056 | PyramidFrustumVolume | num/geometry | volume of a frustum of a pyramid |
| NUM-GEOM-057 | PrismoidVolume | num/geometry | volume by the prismoidal formula |
| NUM-GEOM-058 | ConeVolume | num/geometry | volume of a cone |
| NUM-GEOM-059 | ConeLateralArea | num/geometry | lateral area of a right circular cone |
| NUM-GEOM-060 | ConeSurfaceArea | num/geometry | total surface area of a right circular cone |
| NUM-GEOM-061 | ConeFrustumVolume | num/geometry | volume of a frustum of a cone |
| NUM-GEOM-062 | ConeFrustumSurfaceArea | num/geometry | total surface area of a frustum of a cone |
| NUM-GEOM-063 | SphereVolume | num/geometry | volume of a sphere |
| NUM-GEOM-064 | SphereSurfaceArea | num/geometry | surface area of a sphere |
| NUM-GEOM-065 | SphericalCapVolume | num/geometry | volume of a spherical cap |
| NUM-GEOM-066 | SphericalCapArea | num/geometry | curved area of a spherical cap |
| NUM-GEOM-067 | TorusVolume | num/geometry | volume of a ring torus |
| NUM-GEOM-068 | TorusSurfaceArea | num/geometry | surface area of a ring torus |
| NUM-GEOM-069 | EllipsoidVolume | num/geometry | volume of an ellipsoid |
| NUM-GEOM-070 | EllipsoidSurfaceArea | num/geometry | surface area of an ellipsoid, exact (Carlson) |
| NUM-GEOM-071 | TetrahedronVolume | num/geometry | volume of a regular tetrahedron |
| NUM-GEOM-072 | TetrahedronSurfaceArea | num/geometry | surface area of a regular tetrahedron |
| NUM-GEOM-073 | CubeVolume | num/geometry | volume of a cube |
| NUM-GEOM-074 | CubeSurfaceArea | num/geometry | surface area of a cube |
| NUM-GEOM-075 | OctahedronVolume | num/geometry | volume of a regular octahedron |
| NUM-GEOM-076 | OctahedronSurfaceArea | num/geometry | surface area of a regular octahedron |
| NUM-GEOM-077 | DodecahedronVolume | num/geometry | volume of a regular dodecahedron |
| NUM-GEOM-078 | DodecahedronSurfaceArea | num/geometry | surface area of a regular dodecahedron |
| NUM-GEOM-079 | IcosahedronVolume | num/geometry | volume of a regular icosahedron |
| NUM-GEOM-080 | IcosahedronSurfaceArea | num/geometry | surface area of a regular icosahedron |
| NUM-GEOM-081 | HelixArcLength | num/geometry | arc length of a circular helix |
| NUM-GEOM-082 | HelixHeight | num/geometry | height of a circular helix from its arc length |
| NUM-INTEG-001 | IsPrime | num/integer | whether n is prime, exact for every uint64 |
| NUM-INTEG-002 | PrimeFactors | num/integer | prime factors of n with multiplicity, ascending |
| NUM-INTEG-003 | SmallestPrimeFactor | num/integer | least prime factor of n |
| NUM-LINALG-001 | Vec2Add | num/linalg | a + b in the plane |
| NUM-LINALG-002 | Vec2Sub | num/linalg | a − b in the plane |
| NUM-LINALG-003 | Vec2Scale | num/linalg | k·v in the plane |
| NUM-LINALG-004 | Vec2Div | num/linalg | v/k in the plane, error at k = 0 |
| NUM-LINALG-005 | Vec2Dot | num/linalg | dot product of two 2D vectors, within 2 ulps |
| NUM-LINALG-006 | Vec2Cross | num/linalg | 2D scalar cross product (determinant, right-hand-rule sign) |
| NUM-LINALG-007 | Vec2Norm | num/linalg | length of a 2D vector without overflow |
| NUM-LINALG-008 | Vec2Normalize | num/linalg | unit vector of a 2D vector, error at zero |
| NUM-LINALG-009 | Vec2AngleBetween | num/linalg | angle between two 2D vectors (Kahan's formula) |
| NUM-LINALG-010 | Vec2Rotate | num/linalg | 2D vector turned counterclockwise by an angle |
| NUM-LINALG-011 | Vec3Add | num/linalg | a + b in space |
| NUM-LINALG-012 | Vec3Sub | num/linalg | a − b in space |
| NUM-LINALG-013 | Vec3Scale | num/linalg | k·v in space |
| NUM-LINALG-014 | Vec3Div | num/linalg | v/k in space, error at k = 0 |
| NUM-LINALG-015 | Vec3Dot | num/linalg | dot product of two 3D vectors |
| NUM-LINALG-016 | Vec3Cross | num/linalg | cross product, each component within 2 ulps |
| NUM-LINALG-017 | Vec3Norm | num/linalg | length of a 3D vector without overflow |
| NUM-LINALG-018 | Vec3Normalize | num/linalg | unit vector of a 3D vector, error at zero |
| NUM-LINALG-019 | Vec3Distance | num/linalg | distance between two points in space |
| NUM-LINALG-020 | Vec3AngleBetween | num/linalg | angle between two 3D vectors (Kahan's formula) |
| NUM-LINALG-021 | Vec3DirectionAngles | num/linalg | angles between a vector and the x, y, z axes |
| NUM-LINALG-022 | Vec3ComponentAlong | num/linalg | scalar component of v along an axis |
| NUM-LINALG-023 | Vec3Project | num/linalg | projection of v onto a line |
| NUM-LINALG-024 | Vec3Reject | num/linalg | component of v perpendicular to a line |
| NUM-LINALG-025 | Vec3Reflect | num/linalg | mirror image of v in a plane through the origin |
| NUM-LINALG-026 | Vec3ParallelogramArea | num/linalg | area of the parallelogram spanned by two vectors |
| NUM-LINALG-027 | Vec3TriangleArea | num/linalg | area of the triangle spanned by two vectors |
| NUM-LINALG-028 | Vec3ScalarTripleProduct | num/linalg | a·(b × c), the signed volume (mixed product) |
| NUM-LINALG-029 | Vec3VectorTripleProduct | num/linalg | a × (b × c) |
| NUM-LINALG-030 | Vec3PlaneNormal | num/linalg | unit normal of the plane through three points |
| NUM-LINALG-031 | Vec3Rotate | num/linalg | v turned about an axis by an angle (Rodrigues) |
| NUM-LINALG-032 | Vec3Lerp | num/linalg | linear interpolation, exact at both ends |
| NUM-LINALG-033 | Vec3Min | num/linalg | componentwise minimum |
| NUM-LINALG-034 | Vec3Max | num/linalg | componentwise maximum |
| NUM-LINALG-035 | VecAdd | num/linalg | a + b in n dimensions |
| NUM-LINALG-036 | VecSub | num/linalg | a − b in n dimensions |
| NUM-LINALG-037 | VecScale | num/linalg | k·v in n dimensions |
| NUM-LINALG-038 | VecDiv | num/linalg | v/k in n dimensions, error at k = 0 |
| NUM-LINALG-039 | VecDot | num/linalg | dot product in n dimensions |
| NUM-LINALG-040 | VecDotCompensated | num/linalg | dot product as if in twice the working precision (Dot2) |
| NUM-LINALG-041 | VecNorm | num/linalg | Euclidean norm, scaled and compensated |
| NUM-LINALG-042 | VecNorm1 | num/linalg | 1-norm, compensated sum |
| NUM-LINALG-043 | VecNormInf | num/linalg | infinity norm (largest \|element\|) |
| NUM-LINALG-044 | VecNormP | num/linalg | p-norm for p ≥ 1, scaled |
| NUM-LINALG-045 | VecNormalize | num/linalg | unit vector in n dimensions, error at zero |
| NUM-LINALG-046 | VecDistance | num/linalg | Euclidean distance between two points in n dimensions |
| NUM-LINALG-047 | VecAngleBetween | num/linalg | angle between two n-D vectors (Kahan's formula) |
| NUM-LINALG-048 | VecProject | num/linalg | projection of v onto a line in n dimensions |
| NUM-LINALG-049 | VecLerp | num/linalg | linear interpolation in n dimensions, exact at both ends |
| NUM-LINALG-050 | VecMin | num/linalg | element-by-element minimum |
| NUM-LINALG-051 | VecMax | num/linalg | element-by-element maximum |
| NUM-LINALG-052 | CartesianToCylindrical | num/linalg | (x, y, z) to cylindrical (R, φ, z) |
| NUM-LINALG-053 | CylindricalToCartesian | num/linalg | cylindrical (R, φ, z) to (x, y, z) |
| NUM-LINALG-054 | CartesianToSpherical | num/linalg | (x, y, z) to spherical (R, θ polar, φ azimuth) |
| NUM-LINALG-055 | SphericalToCartesian | num/linalg | spherical (R, θ, φ) to (x, y, z) |
| NUM-LINALG-056 | Mat2Identity | num/linalg | 2×2 identity |
| NUM-LINALG-057 | Mat2Mul | num/linalg | 2×2 matrix product, elements within 2 ulps |
| NUM-LINALG-058 | Mat2MulVec | num/linalg | 2×2 matrix times a 2D vector |
| NUM-LINALG-059 | Mat2Transpose | num/linalg | 2×2 transpose |
| NUM-LINALG-060 | Mat2Det | num/linalg | 2×2 determinant within 2 ulps |
| NUM-LINALG-061 | Mat2Inverse | num/linalg | 2×2 inverse, error when singular |
| NUM-LINALG-062 | Mat2Trace | num/linalg | 2×2 trace |
| NUM-LINALG-063 | Mat2Rotation | num/linalg | 2×2 rotation matrix |
| NUM-LINALG-064 | Mat2SymmetricEigenvalues | num/linalg | eigenvalues of a symmetric 2×2 matrix without cancellation |
| NUM-LINALG-065 | Mat3Identity | num/linalg | 3×3 identity |
| NUM-LINALG-066 | Mat3Mul | num/linalg | 3×3 matrix product |
| NUM-LINALG-067 | Mat3MulVec | num/linalg | 3×3 matrix times a 3D vector |
| NUM-LINALG-068 | Mat3Transpose | num/linalg | 3×3 transpose |
| NUM-LINALG-069 | Mat3Det | num/linalg | 3×3 determinant as a triple product |
| NUM-LINALG-070 | Mat3Inverse | num/linalg | 3×3 inverse by cofactors, error when singular |
| NUM-LINALG-071 | Mat3Trace | num/linalg | 3×3 trace |
| NUM-LINALG-072 | Mat3RotationX | num/linalg | rotation matrix about the x axis |
| NUM-LINALG-073 | Mat3RotationY | num/linalg | rotation matrix about the y axis |
| NUM-LINALG-074 | Mat3RotationZ | num/linalg | rotation matrix about the z axis |
| NUM-LINALG-075 | Mat3RotationAxisAngle | num/linalg | rotation matrix about any axis (Rodrigues) |
| NUM-LINALG-076 | Mat3SymmetricEigenvalues | num/linalg | eigenvalues of a symmetric 3×3 matrix (Smith's closed form) |
| NUM-LINALG-077 | Mat4Identity | num/linalg | 4×4 identity |
| NUM-LINALG-078 | Mat4Mul | num/linalg | 4×4 matrix product |
| NUM-LINALG-079 | Mat4MulVec | num/linalg | 4×4 matrix times a 4-vector |
| NUM-LINALG-080 | Mat4Transpose | num/linalg | 4×4 transpose |
| NUM-LINALG-081 | Mat4Det | num/linalg | 4×4 determinant by Laplace expansion |
| NUM-LINALG-082 | Mat4Inverse | num/linalg | 4×4 inverse by cofactors, error when singular |
| NUM-LINALG-083 | Mat4Trace | num/linalg | 4×4 trace |
| NUM-LINALG-084 | Mat4Translation | num/linalg | homogeneous translation matrix |
| NUM-LINALG-085 | Mat4Scale | num/linalg | homogeneous scaling matrix |
| NUM-LINALG-086 | Mat4FromMat3 | num/linalg | 3×3 linear map embedded in a 4×4 homogeneous matrix |
| NUM-LINALG-087 | Mat4TransformPoint | num/linalg | point transformed by a homogeneous matrix (with the w divide) |
| NUM-LINALG-088 | Mat4TransformDirection | num/linalg | direction transformed by a homogeneous matrix (no translation) |
| NUM-LINALG-089 | MatIdentity | num/linalg | n×n identity |
| NUM-LINALG-090 | MatMul | num/linalg | product of two general matrices |
| NUM-LINALG-091 | MatVecMul | num/linalg | general matrix times a vector |
| NUM-LINALG-092 | MatTranspose | num/linalg | transpose of a general matrix |
| NUM-LINALG-093 | MatTrace | num/linalg | trace of an n×n matrix, compensated |
| NUM-LINALG-094 | LUDecompose | num/linalg | PA = LU with partial pivoting |
| NUM-LINALG-095 | LUSolve | num/linalg | solve Ax = b from an LU factorisation |
| NUM-LINALG-096 | MatSolve | num/linalg | solve Ax = b by LU with partial pivoting |
| NUM-LINALG-097 | MatDet | num/linalg | n×n determinant from LU, without intermediate overflow |
| NUM-LINALG-098 | MatInverse | num/linalg | n×n inverse from LU, error when singular |
| NUM-LINALG-099 | MatRank | num/linalg | numerical rank by complete pivoting with a tolerance |
| NUM-LINALG-100 | MatConditionEstimate1 | num/linalg | estimate of the 1-norm condition number (Hager–Higham) |
| NUM-LINALG-101 | QuatAdd | num/linalg | p + q for quaternions |
| NUM-LINALG-102 | QuatSub | num/linalg | p − q for quaternions |
| NUM-LINALG-103 | QuatMul | num/linalg | Hamilton product p·q |
| NUM-LINALG-104 | QuatConjugate | num/linalg | quaternion conjugate |
| NUM-LINALG-105 | QuatNorm | num/linalg | quaternion magnitude without overflow |
| NUM-LINALG-106 | QuatNormalize | num/linalg | unit quaternion (versor), error at zero |
| NUM-LINALG-107 | QuatInverse | num/linalg | quaternion inverse, error at zero |
| NUM-LINALG-108 | QuatFromAxisAngle | num/linalg | unit quaternion of a rotation about an axis |
| NUM-LINALG-109 | QuatToAxisAngle | num/linalg | axis and angle of a quaternion's rotation |
| NUM-LINALG-110 | QuatFromMat3 | num/linalg | unit quaternion of a rotation matrix (Shepperd) |
| NUM-LINALG-111 | QuatToMat3 | num/linalg | rotation matrix of a quaternion |
| NUM-LINALG-112 | QuatRotateVec3 | num/linalg | 3D vector rotated by a quaternion |
| NUM-LINALG-113 | QuatSlerp | num/linalg | spherical linear interpolation between rotations |
| NUM-LINALG-114 | QuatFromEulerZYX | num/linalg | quaternion from yaw, pitch, roll (aerospace Z-Y-X) |
| NUM-LINALG-115 | QuatToEulerZYX | num/linalg | yaw, pitch, roll from a quaternion, gimbal lock handled |
| NUM-LINALG-116 | MatAdd | num/linalg | A + B entry by entry |
| NUM-LINALG-117 | MatSub | num/linalg | A − B entry by entry |
| NUM-LINALG-118 | MatScale | num/linalg | k·A |
| NUM-LINALG-119 | MatHadamard | num/linalg | Hadamard (entrywise) product A ∘ B |
| NUM-LINALG-120 | VecOuter | num/linalg | outer product u·vᵀ |
| NUM-LINALG-121 | MatKronecker | num/linalg | Kronecker product A ⊗ B |
| NUM-LINALG-122 | MatPower | num/linalg | Aᵏ by repeated squaring |
| NUM-LINALG-123 | Mat3Minor | num/linalg | 3×3 minor: 2×2 determinant with row i, column j deleted |
| NUM-LINALG-124 | Mat2Adjugate | num/linalg | 2×2 adjugate [d −b; −c a] |
| NUM-LINALG-125 | Mat3Adjugate | num/linalg | 3×3 adjugate (transposed cofactors) |
| NUM-LINALG-126 | Mat4Adjugate | num/linalg | 4×4 adjugate (transposed cofactors) |
| NUM-LINALG-127 | MatMinor | num/linalg | n×n minor by LU |
| NUM-LINALG-128 | MatCofactorMatrix | num/linalg | n×n cofactor matrix (det·A⁻ᵀ, or minors when singular) |
| NUM-LINALG-129 | MatAdjugate | num/linalg | n×n adjugate, A·adj A = det(A)·I |
| NUM-LINALG-130 | MatRowEchelon | num/linalg | row echelon form with partial pivoting and tolerance |
| NUM-LINALG-131 | MatRREF | num/linalg | reduced row echelon form (Gauss–Jordan) |
| NUM-LINALG-132 | MatNullSpace | num/linalg | null-space basis from the RREF |
| NUM-LINALG-133 | MatColumnSpace | num/linalg | column-space basis: the pivot columns of A |
| NUM-LINALG-134 | Mat2Characteristics | num/linalg | det, trace, rank, inverse, eigenvalues of a 2×2 |
| NUM-LINALG-135 | Mat3Characteristics | num/linalg | det, trace, rank, inverse, char. polynomial, eigenvalues of a 3×3 |
| NUM-LINALG-136 | MatIsSymmetric | num/linalg | symmetric to a tolerance |
| NUM-LINALG-137 | MatIsDiagonal | num/linalg | diagonal to a tolerance |
| NUM-LINALG-138 | MatIsUpperTriangular | num/linalg | upper triangular to a tolerance |
| NUM-LINALG-139 | MatIsLowerTriangular | num/linalg | lower triangular to a tolerance |
| NUM-LINALG-140 | MatIsOrthogonal | num/linalg | QᵀQ = I to a tolerance |
| NUM-LINALG-141 | MatIsPositiveDefinite | num/linalg | positive definite (Cholesky succeeds) |
| NUM-LINALG-142 | MatIsPositiveSemidefinite | num/linalg | positive semidefinite (smallest eigenvalue ≥ −tol) |
| NUM-LINALG-143 | MatIsSingular | num/linalg | numerical rank below n |
| NUM-LINALG-144 | Mat2SolveCramer | num/linalg | 2×2 system by Cramer's rule |
| NUM-LINALG-145 | Mat3SolveCramer | num/linalg | 3×3 system by Cramer's rule |
| NUM-LINALG-146 | Mat2CharPoly | num/linalg | 2×2 characteristic polynomial |
| NUM-LINALG-147 | Mat3CharPoly | num/linalg | 3×3 characteristic polynomial |
| NUM-LINALG-148 | MatCharPolyFaddeevLeVerrier | num/linalg | n×n characteristic polynomial (Faddeev–LeVerrier) |
| NUM-LINALG-149 | Mat2Eigen | num/linalg | eigenvalues and eigenvectors of a general 2×2 (real or complex) |
| NUM-LINALG-150 | MatSymmetricEigenJacobi | num/linalg | symmetric eigenvalues and vectors by cyclic Jacobi |
| NUM-LINALG-151 | MatPowerIteration | num/linalg | dominant eigenpair by the power method |
| NUM-LINALG-152 | Cholesky | num/linalg | A = L·Lᵀ for symmetric positive definite A |
| NUM-LINALG-153 | CholeskySolve | num/linalg | solve Ax = b from the Cholesky factor |
| NUM-LINALG-154 | QRHouseholder | num/linalg | thin QR by Householder reflections |
| NUM-LINALG-155 | GramSchmidtModified | num/linalg | thin QR by modified Gram–Schmidt |
| NUM-LINALG-156 | LeastSquaresQR | num/linalg | min ‖Ax − b‖₂ by Householder QR |
| NUM-LINALG-157 | ForwardSubstitution | num/linalg | solve Lx = b, L lower triangular |
| NUM-LINALG-158 | BackSubstitution | num/linalg | solve Ux = b, U upper triangular |
| NUM-LINALG-159 | TridiagonalSolve | num/linalg | tridiagonal system by the Thomas algorithm |
| NUM-LINALG-160 | SVDJacobi | num/linalg | thin SVD by one-sided Jacobi (Hestenes) |
| NUM-LINALG-161 | MatPseudoInverse | num/linalg | Moore–Penrose pseudo-inverse from the SVD |
| NUM-LINALG-162 | MatCondition2 | num/linalg | 2-norm condition number σ_max/σ_min |
| NUM-LINALG-163 | MatNormFrobenius | num/linalg | Frobenius norm |
| NUM-LINALG-164 | MatNorm1 | num/linalg | 1-norm: largest absolute column sum |
| NUM-LINALG-165 | MatNormInf | num/linalg | ∞-norm: largest absolute row sum |
| NUM-LINALG-166 | MatNorm2 | num/linalg | spectral norm σ_max |
| NUM-LINALG-167 | CovarianceToCorrelation | num/linalg | correlation matrix and standard deviations from a covariance matrix |
| NUM-LINALG-168 | CorrelationToCovariance | num/linalg | covariance matrix from correlations and standard deviations |
| NUM-LOGS-001 | GoldenRatio | num/logs | the golden ratio φ = (1+√5)/2 |
| NUM-PROB-001 | ProbabilityFromCounts | num/prob | favorable/total for equally likely outcomes |
| NUM-PROB-002 | ComplementProbability | num/prob | 1 − P(A) |
| NUM-PROB-003 | IndependentJointProbability | num/prob | P(A)·P(B) for independent events |
| NUM-PROB-004 | JointProbability | num/prob | P(A given B)·P(B), the multiplication rule |
| NUM-PROB-005 | UnionProbability | num/prob | P(A) + P(B) − P(A ∩ B) |
| NUM-PROB-006 | ConditionalProbability | num/prob | P(A ∩ B)/P(B) |
| NUM-PROB-007 | TotalProbability | num/prob | Σ P(A given Bᵢ)·P(Bᵢ) over a partition |
| NUM-PROB-008 | BayesTheorem | num/prob | P(B given A)·P(A)/P(B) |
| NUM-PROB-009 | BayesPosterior | num/prob | posteriors over a partition, appended to dst |
| NUM-PROB-010 | PositivePredictiveValue | num/prob | P(condition given a positive test) from Se, Sp, prevalence |
| NUM-PROB-011 | NegativePredictiveValue | num/prob | P(no condition given a negative test) from Se, Sp, prevalence |
| NUM-PROB-012 | OddsFromProbability | num/prob | odds in favour p/(1 − p) |
| NUM-PROB-013 | ProbabilityFromOdds | num/prob | probability o/(1 + o) from odds in favour |
| NUM-PROB-014 | DecimalOddsFromProbability | num/prob | fair decimal odds 1/p |
| NUM-PROB-015 | ImpliedProbability | num/prob | probability 1/d implied by decimal odds |
| NUM-PROB-016 | Overround | num/prob | booksum Σ 1/dᵢ minus 1 of one market |
| NUM-PROB-017 | NormalizedImpliedProbabilities | num/prob | implied probabilities with the overround removed proportionally |
| NUM-PROB-018 | ExpectedValue | num/prob | E[X] of a discrete distribution |
| NUM-PROB-019 | DiscreteVariance | num/prob | Var X of a discrete distribution, two-pass |
| NUM-PROB-020 | ChebyshevBound | num/prob | Chebyshev's bound 1/k² on a deviation of k or more standard deviations |
| NUM-PROB-021 | ChebyshevBoundEpsilon | num/prob | Chebyshev's bound σ²/ε² on a deviation of ε or more |
| NUM-PROB-022 | MarkovBound | num/prob | Markov's bound E[X]/a on P(X ≥ a) |
| NUM-PROB-023 | BernoulliPMF | num/prob | Bernoulli P(X = k) |
| NUM-PROB-024 | BernoulliCDF | num/prob | Bernoulli P(X ≤ k) |
| NUM-PROB-025 | BernoulliMeanVariance | num/prob | Bernoulli mean and variance |
| NUM-PROB-026 | BinomialPMF | num/prob | binomial P(X = k), saddle-point form |
| NUM-PROB-027 | BinomialCDF | num/prob | binomial P(X ≤ k), compensated tail sum |
| NUM-PROB-028 | BinomialMeanVariance | num/prob | binomial mean and variance |
| NUM-PROB-029 | GeometricTrialsPMF | num/prob | geometric P(X = k), X = trials to the first success |
| NUM-PROB-030 | GeometricTrialsCDF | num/prob | geometric P(X ≤ k), X = trials to the first success |
| NUM-PROB-031 | GeometricTrialsMeanVariance | num/prob | geometric mean and variance, trials convention |
| NUM-PROB-032 | NegativeBinomialTrialsPMF | num/prob | negative binomial P(X = x), X = trials to the r-th success |
| NUM-PROB-033 | NegativeBinomialTrialsCDF | num/prob | negative binomial P(X ≤ x), X = trials to the r-th success |
| NUM-PROB-034 | NegativeBinomialTrialsMeanVariance | num/prob | negative binomial mean and variance, trials convention |
| NUM-PROB-035 | HypergeometricPMF | num/prob | hypergeometric P(X = k), saddle-point form |
| NUM-PROB-036 | HypergeometricCDF | num/prob | hypergeometric P(X ≤ k), compensated tail sum |
| NUM-PROB-037 | HypergeometricMeanVariance | num/prob | hypergeometric mean and variance |
| NUM-PROB-038 | DiscreteUniformPMF | num/prob | discrete uniform P(X = k) on a…b |
| NUM-PROB-039 | DiscreteUniformCDF | num/prob | discrete uniform P(X ≤ k) on a…b |
| NUM-PROB-040 | DiscreteUniformMeanVariance | num/prob | discrete uniform mean and variance on a…b |
| NUM-PROB-041 | PoissonPMF | num/prob | Poisson P(X = k), saddle-point form |
| NUM-PROB-042 | PoissonCDF | num/prob | Poisson P(X ≤ k), compensated tail sum |
| NUM-PROB-043 | PoissonMeanVariance | num/prob | Poisson mean and variance |
| NUM-PROB-044 | UniformPDF | num/prob | continuous uniform density on [a, b] |
| NUM-PROB-045 | UniformCDF | num/prob | continuous uniform P(X ≤ x) on [a, b] |
| NUM-PROB-046 | UniformMeanVariance | num/prob | continuous uniform mean and variance on [a, b] |
| NUM-PROB-047 | NormalPDF | num/prob | normal density |
| NUM-PROB-048 | NormalCDF | num/prob | normal P(X ≤ x), accurate in the lower tail |
| NUM-PROB-049 | NormalSurvival | num/prob | normal P(X > x), accurate in the upper tail |
| NUM-PROB-050 | NormalIntervalProbability | num/prob | normal P(lo ≤ X ≤ hi), tail-aware |
| NUM-PROB-051 | LognormalPDF | num/prob | lognormal density |
| NUM-PROB-052 | LognormalCDF | num/prob | lognormal P(X ≤ x) |
| NUM-PROB-053 | LognormalMeanVariance | num/prob | lognormal mean and variance |
| NUM-PROB-054 | ExponentialPDF | num/prob | exponential density, rate form |
| NUM-PROB-055 | ExponentialCDF | num/prob | exponential P(X ≤ x), rate form |
| NUM-PROB-056 | ExponentialMeanVariance | num/prob | exponential mean and variance, rate form |
| NUM-PROB-057 | GeometricTrialsQuantile | num/prob | least k with P(X ≤ k) ≥ u, trials convention |
| NUM-PROB-058 | DiscreteUniformQuantile | num/prob | least k in a…b with P(X ≤ k) ≥ u |
| NUM-PROB-059 | UniformQuantile | num/prob | continuous uniform quantile a + u(b − a) |
| NUM-PROB-060 | ExponentialQuantile | num/prob | exponential quantile −ln(1 − u)/λ |
| NUM-PROB-061 | NormalQuantile | num/prob | normal quantile by Wichura's PPND16 |
| NUM-PROB-062 | LognormalQuantile | num/prob | lognormal quantile exp(μ + σΦ⁻¹(u)) |
| NUM-PROB-063 | ChiSquareCDF | num/prob | chi-square P(X ≤ x), regularized incomplete gamma P(ν/2, x/2) |
| NUM-PROB-064 | ChiSquareSurvival | num/prob | chi-square P(X > x), the p-value, accurate in the upper tail |
| NUM-PROB-065 | ChiSquareQuantile | num/prob | chi-square quantile (critical value), safeguarded Newton |
| NUM-PROB-066 | LikelihoodRatioPositive | num/prob | LR+ = sensitivity/(1 − specificity) |
| NUM-PROB-067 | LikelihoodRatioNegative | num/prob | LR− = (1 − sensitivity)/specificity |
| NUM-PROB-068 | PostTestProbability | num/prob | probability after a result with likelihood ratio LR (odds × LR) |
| NUM-SERIES-001 | TheilSenSlope | num/series | median of the pairwise slopes, randomized slope selection |
| NUM-SERIES-002 | TheilSenSlopePerDay | num/series | Theil–Sen slope of point values per day of elapsed time |
| NUM-SERIES-003 | TrailingTimeWindowSum | num/series | sum and count of the values in (t − window, t] at each point |
| NUM-SERIES-004 | LocalDayWindow | num/series | [from, to) covering the last N local calendar days |
| NUM-SERIES-005 | LocalDayBuckets | num/series | index ranges of the points of each local calendar date |
| NUM-SERIES-006 | FixedStrideBuckets | num/series | min, mean, max and count of the points in each fixed stride |
| NUM-SERIES-007 | WeekStart | num/series | local midnight of the first day of an instant's week |
| NUM-SERIES-008 | WeeklySums | num/series | sum and count of the values in each local week of a window |
| NUM-SERIES-009 | AverageRealVariability | num/series | mean absolute change between consecutive readings (ARV) |
| NUM-SERIES-010 | AverageRateOfChangePerDay | num/series | (last − first)/elapsed days between two points, the secant slope |
| NUM-SERIES-011 | LocalHourWindowMean | num/series | mean and count of the points in a window of local clock hours |
| NUM-STATS-001 | ArithmeticMean | num/stats | sum of the values divided by their count |
| NUM-STATS-002 | RelativeFrequency | num/stats | each class count divided by the total count |
| NUM-STATS-003 | CoefficientOfVariation | num/stats | sample standard deviation divided by the mean |
| NUM-STATS-004 | CohenD | num/stats | difference of two means in standard deviations (Cohen's d) |
| NUM-STATS-005 | GroupedMean | num/stats | estimated mean of grouped data from class marks |
| NUM-STATS-006 | GroupedMedian | num/stats | estimated median of grouped data by interpolation |
| NUM-STATS-008 | FrequencyDistribution | num/stats | count of values in each class |
| NUM-STATS-009 | LinearRegression | num/stats | intercept and slope of the least-squares line |
| NUM-STATS-010 | MatthewsCorrelation | num/stats | Matthews correlation coefficient of a 2×2 table |
| NUM-STATS-011 | MeanAbsoluteDeviation | num/stats | average absolute deviation about the mean |
| NUM-STATS-012 | Median | num/stats | middle value, or mean of the two middle values |
| NUM-STATS-013 | Midrange | num/stats | mean of the smallest and largest values |
| NUM-STATS-014 | PearsonCorrelation | num/stats | Pearson correlation coefficient of paired samples |
| NUM-STATS-015 | Percentile | num/stats | p-quantile by a chosen Hyndman–Fan method |
| NUM-STATS-016 | PooledStandardDeviation | num/stats | pooled standard deviation of two samples |
| NUM-STATS-017 | PopulationVariance | num/stats | variance with divisor N |
| NUM-STATS-018 | Range | num/stats | largest value minus smallest value |
| NUM-STATS-019 | SampleStandardDeviation | num/stats | standard deviation with divisor n − 1 |
| NUM-STATS-020 | SummaryStatistics | num/stats | count, mean, sample SD, min, max, median and range |
| NUM-STATS-021 | SampleVariance | num/stats | variance with divisor n − 1 |
| NUM-STATS-022 | PopulationStandardDeviation | num/stats | standard deviation with divisor N |
| NUM-STATS-023 | StandardErrorOfMean | num/stats | sample standard deviation divided by √n |
| NUM-STATS-024 | Mode | num/stats | most frequent values, all ties, ascending |
| NUM-STATS-025 | SumOfSquares | num/stats | sum of squared deviations about the mean |
| NUM-STATS-026 | WeightedMean | num/stats | weighted mean Σw·x / Σw |
| NUM-STATS-028 | ZScore | num/stats | standard score (x − mean)/sd |
| NUM-STATS-029 | RawScoreFromZScore | num/stats | value with a given standard score, mean + z·sd |
| NUM-STATS-030 | PercentRelativeStandardDeviation | num/stats | 100·s/abs(mean), the percent relative standard deviation |
| NUM-STATS-031 | MedianAbsoluteDeviation | num/stats | median absolute deviation about the median |
| NUM-STATS-032 | TukeyFences | num/stats | quartiles, IQ and inner/outer fences Q1 − 1.5·IQ … Q3 + 3·IQ |
| NUM-STATS-033 | ChiSquareGoodnessOfFit | num/stats | Pearson χ² = Σ (O − E)²/E against expected counts |
| NUM-STATS-034 | ChiSquareGoodnessOfFitUniform | num/stats | Pearson χ² against equal expected counts N/k |
| NUM-STATS-035 | FrequencyWeightedPopulationStandardDeviation | num/stats | population SD of a frequency table, √(Σf(x − μ)²/N) |
| NUM-STATS-036 | FrequencyWeightedSampleStandardDeviation | num/stats | sample SD of a frequency table, √(Σf(x − x̄)²/(n − 1)) |
| NUM-STATS-037 | FrequencyWeightedMedian | num/stats | median of a frequency table without expanding it |
| NUM-TRIG-001 | DegToRad | num/trig | degrees to radians |
| NUM-TRIG-002 | RadToDeg | num/trig | radians to degrees |
| NUM-TRIG-003 | NormalizeAngle | num/trig | angle reduced to [0, 2π) |
| NUM-TRIG-004 | NormalizeAngleSigned | num/trig | angle reduced to (−π, π] |
| NUM-TRIG-005 | AngleDifference | num/trig | shortest signed turn from one angle to another, in (−π, π] |
| NUM-TRIG-006 | Sec | num/trig | secant |
| NUM-TRIG-007 | Csc | num/trig | cosecant; error at 0 |
| NUM-TRIG-008 | Cot | num/trig | cotangent; error at 0 |
| NUM-TRIG-009 | Asec | num/trig | inverse secant, principal value in [0, π] |
| NUM-TRIG-010 | Acsc | num/trig | inverse cosecant, principal value in [−π/2, π/2] |
| NUM-TRIG-011 | Acot | num/trig | inverse cotangent, principal value in (−π/2, π/2] |
| NUM-TRIG-012 | Sech | num/trig | hyperbolic secant |
| NUM-TRIG-013 | Csch | num/trig | hyperbolic cosecant; error at 0 |
| NUM-TRIG-014 | Coth | num/trig | hyperbolic cotangent; error at 0 |
| NUM-TRIG-015 | Asech | num/trig | inverse hyperbolic secant on (0, 1] |
| NUM-TRIG-016 | Acsch | num/trig | inverse hyperbolic cosecant; error at 0 |
| NUM-TRIG-017 | Acoth | num/trig | inverse hyperbolic cotangent on \|x\| > 1 |
| NUM-TRIG-018 | LawOfSinesSide | num/trig | side from a side and two angles (law of sines) |
| NUM-TRIG-019 | LawOfSinesAngle | num/trig | the 0, 1 or 2 angles from two sides and a non-included angle (SSA) |
| NUM-TRIG-020 | LawOfCosinesSide | num/trig | third side from two sides and the included angle |
| NUM-TRIG-021 | LawOfCosinesAngle | num/trig | angle opposite a side from three sides |
| NUM-TRIG-022 | SolveTriangleSSS | num/trig | all parts of a triangle from three sides |
| NUM-TRIG-023 | SolveTriangleSAS | num/trig | all parts from two sides and the included angle |
| NUM-TRIG-024 | SolveTriangleASA | num/trig | all parts from two angles and the included side |
| NUM-TRIG-025 | SolveTriangleAAS | num/trig | all parts from two angles and a non-included side |
| NUM-TRIG-026 | SolveTriangleSSA | num/trig | the 0, 1 or 2 triangles from two sides and a non-included angle |
| NUM-TRIG-027 | SolveRightTriangleLegs | num/trig | right triangle from its two legs |
| NUM-TRIG-028 | SolveRightTriangleHypotenuseLeg | num/trig | right triangle from hypotenuse and one leg |
| NUM-TRIG-029 | SolveRightTriangleHypotenuseAngle | num/trig | right triangle from hypotenuse and an acute angle |
| NUM-TRIG-030 | SolveRightTriangleLegAngle | num/trig | right triangle from a leg and its opposite acute angle |
| NUM-TRIG-031 | CentralAngle | num/trig | central angle between two latitude/longitude points (haversine form) |
| UNT-LEN-001 | CentimetersToInches | unit/length | l/in = (l/cm)/2.54 |
| UNT-LEN-002 | InchesToCentimeters | unit/length | l/cm = 2.54 · l/in |
| UNT-MASS-001 | KilogramsToPounds | unit/mass | m/lb = (m/kg)/0.453 592 37 |
| UNT-MASS-002 | PoundsToKilograms | unit/mass | m/kg = 0.453 592 37 · m/lb |
| UNT-TEMP-001 | FahrenheitToCelsius | unit/temperature | t/°C = (t/°F − 32)/1.8 |
| UNT-TEMP-002 | CelsiusToFahrenheit | unit/temperature | t/°F = 1.8 · t/°C + 32 |
| UNT-TEMP-003 | CelsiusToKelvin | unit/temperature | T/K = t/°C + 273.15 |
| UNT-TEMP-004 | KelvinToCelsius | unit/temperature | t/°C = T/K − 273.15 |
