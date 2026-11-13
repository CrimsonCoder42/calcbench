# algo/hashing (SHORT = HASH)

Builder W1-B3. Package name `hashing`, not `hash`: callers import the standard `hash` package beside ours.
Not yet in `doc/taxonomy.md` (lead to add the row: `ALG | algo/hashing | HASH | counting, grouping and set
operations with hash tables`).

Common to the package:
- **Hash-table costs**: CLRS = Cormen, Leiserson, Rivest, Stein, *Introduction to Algorithms*, 3rd ed., MIT
  Press, 2009, ch. 11 (Theorems 11.1 and 11.2: under simple uniform hashing, a search costs Θ(1 + α) expected).
  Go's map is a hash table whose operations are expected O(1); every "Time" below is expected (average) time,
  and the worst case of each map operation is O(m) (all keys colliding). Second source: Sedgewick & Wayne,
  *Algorithms*, 4th ed., Addison-Wesley, 2011, §3.4 (hash tables, Proposition K: constant expected cost under
  the uniform hashing assumption). Knuth, TAOCP Vol. 3, 2nd ed., 1998, §6.4 (hashing).
- **Order**: the Go language specification, "For statements with range clause": "The iteration order over maps
  is not specified and is not guaranteed to be the same from one iteration to the next." No output of this
  package depends on map iteration order: groups and set results follow the order of first appearance in
  the input; ranked results have a total order on (count, value).
- **Equality** is Go's `==`. A float NaN is not equal to itself, so each NaN is its own key: it never matches
  and is never merged with another NaN. -0 and +0 are one key.
- No function modifies its input.

| ID | Name | Signature |
|---|---|---|
| ALG-HASH-001 | FrequencyCount | `func FrequencyCount[T comparable](values []T) map[T]int` |
| ALG-HASH-002 | TopKFrequent | `func TopKFrequent[T cmp.Ordered](values []T, k int) ([]Count[T], error)` |
| ALG-HASH-003 | TopKByCount | `func TopKByCount[T cmp.Ordered](counts map[T]int, k int) ([]Count[T], error)` |
| ALG-HASH-004 | GroupBy | `func GroupBy[T any, K comparable](values []T, key func(T) K) ([]Group[K, T], error)` |
| ALG-HASH-005 | FirstUnique | `func FirstUnique[T comparable](values []T) int` |
| ALG-HASH-006 | SetIntersection | `func SetIntersection[T comparable](a, b []T) []T` |
| ALG-HASH-007 | SetUnion | `func SetUnion[T comparable](a, b []T) []T` |
| ALG-HASH-008 | SetDifference | `func SetDifference[T comparable](a, b []T) []T` |
| ALG-HASH-009 | LongestConsecutive | `func LongestConsecutive[T Integer](values []T) (start T, length int)` |
| ALG-HASH-010 | GroupAnagrams | `func GroupAnagrams(words []string) [][]string` |
| ALG-HASH-011 | GroupAnagramsByCount | `func GroupAnagramsByCount(words []string, alphabet string) ([][]string, error)` |
| ALG-HASH-012 | TwoSum | `func TwoSum[T Integer](values []T, target T) (i, j int)` |
| ALG-HASH-013 | ContainsNearbyDuplicate | `func ContainsNearbyDuplicate[T comparable](values []T, maxGap int) (bool, error)` |

`Count[T]` is `{Value T; Count int}`. `Group[K, T]` is `{Key K; Items []T}`. `Integer` is the set of Go's
built-in integer types (`~int | ~int8 | … | ~uintptr`).

### ALG-HASH-001 FrequencyCount
- **Source**: CLRS ch. 11 (direct-address and hash tables as a dictionary); Sedgewick & Wayne §3.1, client
  FrequencyCounter (a symbol table counting occurrences), §3.4.
- **Best known**: O(n) expected time, Ω(n) to read the input; O(m) space for the m distinct values (the
  result). A sort-and-scan is O(n log n) and needs an ordered type.
- **Approach**: one map, not pre-sized (pre-sizing to n wastes memory when m ≪ n); it grows as keys arrive. **Lost**:
  sort a copy then count runs (O(n log n), O(n) extra, needs `cmp.Ordered`).
- **Input modified**: no. **Errors**: none. Empty input → empty, non-nil map.
- **Worked example**: Sedgewick & Wayne §3.1's FrequencyCounter idea on our own words: "to be or not to be"
  → to 2, be 2, or 1, not 1. Hand-derived.

### ALG-HASH-002 TopKFrequent · ALG-HASH-003 TopKByCount
- **Source**: counting by hash table (CLRS ch. 11) followed by a bounded min-heap of the k best (Sedgewick &
  Wayne §2.4, client TopM; Williams, "Algorithm 232: Heapsort", CACM 7(6):347–348, 1964,
  doi:10.1145/512274.3734138; Floyd, "Algorithm 245: Treesort", CACM 7(12):701, 1964,
  doi:10.1145/355588.365103). Same pattern as ALG-HEAP-001, keyed on (count, value).
- **Order (deterministic)**: count descending; equal counts by value ascending (`cmp.Compare`). Keys of a map
  are distinct, so (count, value) is a total order and map iteration order never reaches the output.
- **Best known**: O(n + m) expected (counting plus selection, CLRS §9.3) then O(k log k) to order the k;
  the lower bound is Ω(n + k log k). Bounded heap: O(n + m log k) expected time, O(m) for the counts (TopKFrequent
  only) and O(k) for the result.
- **Approach**: the bounded heap. **Lost**: bucket sort by count (CLRS §8.4): O(n + m) time, but n+1 buckets
  of O(n) extra space, and the deterministic tie-break needs every bucket sorted by value, which is
  O(m log m) when all counts are equal: no better than the heap in the worst case and more memory.
  Also lost: sorting all m (count, value) pairs (O(m log m)).
- **Input modified**: no. **Errors**: `ErrNegativeK` (k < 0). k > m returns all m, ranked. k = 0 → empty,
  non-nil slice.
- **Worked example**: counts of "to be or not to be" with k = 2 → (be 2), (to 2): equal counts, "be" < "to".
  Hand-derived.
- **Consumer evidence**: `the consumer backend/knowledgebase/series/rank.go:19-37` hand-rolls a `countHeap` over
  `container/heap` and a `TopK(counts map[K]int, k)` with exactly this order (count desc, key asc); its caller
  is `knowledgebase/analysis/symptoms/symptoms.go:661`. Consumer change: replace `series.TopK(counts, k)` with
  `hashing.TopKByCount(counts, k)` (field `Key` becomes `Value`; k < 0 is an error instead of nil; k = 0 gives
  an empty slice instead of nil).
- **Tags**: LeetCode 347, 692.

### ALG-HASH-004 GroupBy
- **Source**: CLRS ch. 11 (dictionary operations INSERT and SEARCH); Sedgewick & Wayne §3.4.
- **Best known**: O(n) expected time (n calls of key, n map operations), Ω(n). O(g) extra for the map of g
  distinct keys; O(n) for the result.
- **Approach**: a map from key to group index into a result slice; groups in order of their key's first
  appearance, items in input order (stable). **Lost**: sort by key then split runs (O(n log n), needs an
  ordered key, loses input order unless stable).
- **Input modified**: no (key is called once per value). **Errors**: `ErrNilKeyFunc` (key == nil).
- **Worked example**: 1…7 grouped by parity → [odd: 1 3 5 7], [even: 2 4 6]. Hand-derived.

### ALG-HASH-005 FirstUnique
- **Source**: counting with a hash table (CLRS ch. 11), then a second scan.
- **Best known**: O(n) expected, Ω(n) (an unread element might be the unique one). O(m) extra.
- **Approach**: count, then the first index whose count is 1. **Lost**: a linked hash map of candidates
  (one pass, but more bookkeeping and no gain on a slice that can be re-read); O(n²) pairwise scan.
- **Returns**: the index, or −1 when every value repeats or the slice is empty. **Errors**: none.
- **Worked example**: the bytes of "statistics" (s 3, t 3, a 1, i 2, c 1) → index 2 ('a'). Hand-derived.
- **Tags**: LeetCode 387.

### ALG-HASH-006 SetIntersection · ALG-HASH-007 SetUnion · ALG-HASH-008 SetDifference
- **Source**: set operations on hash tables, CLRS ch. 11; Knuth, TAOCP Vol. 3, 2nd ed., §6.4.
- **Order (deterministic)**: the result holds distinct values. Intersection: the values of a that occur in b,
  in order of first appearance in a. Union: the distinct values of a in order of first appearance, then the
  values of b not in a, in order of first appearance in b. Difference: the values of a not in b, in order of
  first appearance in a.
- **Best known**: O(|a| + |b|) expected time, Ω(|a| + |b|). O(|a| + |b|) extra (a set of b and a seen-set).
- **Approach**: a set of b and a seen-set for de-duplication. **Lost**: sort both and merge (O(N log N), needs
  `cmp.Ordered`, and the output order would be sorted rather than input order).
- **Input modified**: no. **Errors**: none. Empty results are non-nil.
- **Worked example**: a = 1 2 2 3 4, b = 4 3 3 5 → ∩ = 3 4, ∪ = 1 2 3 4 5, a∖b = 1 2. Hand-derived.
- **Tags**: LeetCode 349.

### ALG-HASH-009 LongestConsecutive
- **Source**: membership in a hash set, CLRS ch. 11 (expected O(1) SEARCH); the run is counted only from a
  value whose predecessor is absent, so each value is visited O(1) times.
- **Best known**: O(n) expected time, O(n) extra; the comparison lower bound for the sorting approach is
  Ω(n log n) (element distinctness reduces to it in the comparison model), which hashing sidesteps.
- **Approach**: set of values; for each value v with v−1 absent, walk v+1, v+2, … while present. Overflow
  safe: the walk stops at the type's maximum and the predecessor test skips the type's minimum.
  Ties broken by the smallest start. **Lost**: sort and scan (O(n log n), O(n) copy).
- **Returns**: (start, length); empty input → (0, 0). Duplicates count once. **Errors**: none.
- **Worked example**: 10 4 20 1 3 2 11 → start 1, length 4 (1 2 3 4). Hand-derived.
- **Tags**: LeetCode 128.

### ALG-HASH-010 GroupAnagrams · ALG-HASH-011 GroupAnagramsByCount
- **Source**: the anagram signature: Bentley, *Programming Pearls*, 2nd ed., Addison-Wesley, 2000, Column 2,
  problem C (posed in §2.1; signatures in §2.8, pp. 18–20): sort the letters of each word to get its signature; words with equal
  signatures are anagrams. Grouping by signature: a hash table (CLRS ch. 11).
- **Definition**: two words are anagrams when they hold the same multiset of runes (Unicode code points).
  Groups are in order of first appearance; words keep input order; duplicates are kept.
- **Best known**: GroupAnagrams: O(Σ Lᵢ log Lᵢ) time (sorting each word), O(Σ Lᵢ) extra for the keys.
  GroupAnagramsByCount: O(Σ Lᵢ + n·σ) time with σ = number of runes in the alphabet, O(n·σ) extra for the
  keys; it is faster when words are long and σ is small (the classic 26-letter count key).
- **Approach**: sorted-rune key (Bentley's signature); count key over a caller-given alphabet, encoded as
  uvarint counts. **Lost**: a product-of-primes signature (overflows past a few letters); pairwise anagram
  tests (O(n²)).
- **Input modified**: no. **Errors**: ByCount: `ErrAlphabet` (a rune repeated in alphabet, or invalid
  UTF-8 in it), `ErrNotInAlphabet` (a word holds a rune outside alphabet).
- **Worked example**: Bentley's own (Column 2): "deposit", "dopiest", "posited", "topside" share the signature
  "deiopst". Test: those four plus "pots", "stop", "tops" → 2 groups.
- **Tags**: LeetCode 49.

### ALG-HASH-012 TwoSum
- **Source**: complement lookup in a hash table (CLRS ch. 11); the problem is CLRS 3rd ed., Exercise 2.3-7
  (a pair summing to x, solved there in Θ(n lg n) by sorting), here in O(n) expected with hashing.
- **Best known**: O(n) expected time, O(n) extra. Sorting plus two pointers is O(n log n) and loses the
  original indices without an O(n) index array.
- **Approach**: scan j = 0…n−1; look up target − values[j] in a map of value → first index. Returns the pair
  with the smallest j, and for it the smallest i < j. The subtraction is overflow-checked: if target − v is
  out of the type's range, no partner exists for v. **Lost**: sort + two pointers (above).
- **Returns**: (i, j) with i < j, or (−1, −1). **Errors**: none.
- **Worked example**: 3 8 11 −2 7, target 9 → (2, 3), since 11 + (−2) = 9 and no earlier pair sums to 9.
  Hand-derived.
- **Tags**: LeetCode 1.

### ALG-HASH-013 ContainsNearbyDuplicate
- **Source**: a sliding window held in a hash set (CLRS ch. 11 dictionary; the window as in ALG-WINDOW).
- **Best known**: O(n) expected time; O(min(n, maxGap)) extra, the window's set.
- **Approach**: a set of the last maxGap values; insert, and delete the one leaving. **Lost**: a map of last
  index per value (O(m) extra, m distinct, which is larger than the window when maxGap ≪ m).
- **Returns**: true when some i < j has values[i] == values[j] and j − i ≤ maxGap. maxGap = 0 → false.
- **Errors**: `ErrNegativeGap` (maxGap < 0).
- **Worked example**: 5 1 7 5 with maxGap 3 → true; with maxGap 2 → false. Hand-derived.
- **Tags**: LeetCode 219.
