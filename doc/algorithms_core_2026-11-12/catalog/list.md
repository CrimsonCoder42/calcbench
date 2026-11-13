# algo/list (SHORT = LIST) — S2 additions (builder W2-B5)

New package; `algo/list` is in `doc/taxonomy.md` ("linked lists, LRU"). Every ID below is new (001…017).

## Decisions for the lead (each with the alternative that lost)

| # | Decision | Alternative | Why it lost |
|---|---|---|---|
| B5-1 | Functions that relink a caller's nodes carry `InPlace` (`ReverseInPlace`, `ReverseKGroupInPlace`, `MergeTwoSortedInPlace`, `RemoveNthFromEndInPlace`, `ReorderInPlace`, `PartitionInPlace`), as L21 | the bare names of the brief (`Reverse`, …) | BRIEF_COMMON: "inputs never modified unless the name says so (…InPlace)"; a copying `Reverse` would allocate n nodes to do what the pattern exists to avoid |
| B5-2 | Every singly linked walk is guarded by Floyd's two-pointer test (`walk`, unexported): a cyclic input returns `ErrCycle` and is left unchanged; nil is the empty list | trust the caller (no check) / a visited-set map | an unguarded walk hangs on a cycle (the brief forbids it); a map costs O(n) space and allocation; Floyd costs ~1.5 extra pointer reads per node and O(1) space |
| B5-3 | `DoublyLinkedList[T]` is built although `container/list` exists | use `container/list` | container/list stores `any`: every Value is boxed (an extra allocation for non-pointer T) and type-asserted at each read; no `iter.Seq`; ours is generic, unboxed, has `All`/`Backward` iterators, and the same O(1) ownership check (an Element from another list is refused, never corrupts). LRU and LFU are built on it |
| B5-4 | Package name `list` (taxonomy-approved) | `linkedlist` | a caller importing `container/list` beside ours must rename one; `container/list` is not in CLAUDE.md's forbidden list and our `DoublyLinkedList` replaces its use. Flag for Devin |
| B5-5 | Types (Node, DoublyLinkedList, LRUCache, LFUCache) take one ID each, methods listed in the type's doc comment (ALG-QUEUE-002 shape) | an ID per method (ALG-HEAP-004…007 shape) | fewer rows for one concept; both shapes exist in S2 — lead to pick one for the index |
| B5-6 | A key not equal to itself (a float NaN, or a struct holding one) is refused by the caches' Put with `ErrKeyNotComparable` | store it | Go map lookups and deletes never find a NaN key, so it could never be read or evicted: the map would grow past capacity (L17 treats NaN as never equal) |
| B5-7 | LFU ties (equal use count) broken by least recent use; Put on an existing key counts as a use | arbitrary within a frequency (the paper's SET) | deterministic eviction; the order costs nothing with a linked node list |

## Consumer evidence

`/Users/devinanderson/Desktop/red_lite_cache/internal/cache` (read 2026-11-12): **no hand-rolled LRU or
eviction list found.** `store.go:15-16` "the store itself does not evict" (bounded by MaxKeys and TTL);
`groups_bench_test.go:15` and `:163` show the eviction it benchmarks is a uniform random draw from a dense slot
array (its D-043), not recency. `ordered.go:30` is a B-tree leaf chain (`next *bnode`), not a cache list.
So LRUCache/LFUCache answer a stated need ("LRU needs") but replace no existing code there.

| ID | Name | Signature |
|---|---|---|
| ALG-LIST-001 | Node | `type Node[T any] struct { Value T; Next *Node[T] }` |
| ALG-LIST-002 | FromSlice | `func FromSlice[T any](values []T) *Node[T]` |
| ALG-LIST-003 | AppendValues | `func AppendValues[T any](dst []T, head *Node[T]) ([]T, error)` |
| ALG-LIST-004 | Length | `func Length[T any](head *Node[T]) (int, error)` |
| ALG-LIST-005 | DoublyLinkedList | `type DoublyLinkedList[T any]`; `type Element[T any]`; methods PushFront, PushBack, InsertBefore, InsertAfter, Remove, MoveToFront, MoveToBack, Front, Back, Len, All, Backward; Element.Next, Element.Prev |
| ALG-LIST-006 | ReverseInPlace | `func ReverseInPlace[T any](head *Node[T]) (*Node[T], error)` |
| ALG-LIST-007 | ReverseKGroupInPlace | `func ReverseKGroupInPlace[T any](head *Node[T], k int) (*Node[T], error)` |
| ALG-LIST-008 | DetectCycle | `func DetectCycle[T any](head *Node[T]) bool` |
| ALG-LIST-009 | CycleStart | `func CycleStart[T any](head *Node[T]) (start *Node[T], cycleLen int)` |
| ALG-LIST-010 | Middle | `func Middle[T any](head *Node[T]) (*Node[T], error)` |
| ALG-LIST-011 | MergeTwoSortedInPlace | `func MergeTwoSortedInPlace[T cmp.Ordered](a, b *Node[T]) (*Node[T], error)` |
| ALG-LIST-012 | RemoveNthFromEndInPlace | `func RemoveNthFromEndInPlace[T any](head *Node[T], fromEnd int) (*Node[T], error)` |
| ALG-LIST-013 | ReorderInPlace | `func ReorderInPlace[T any](head *Node[T]) (*Node[T], error)` |
| ALG-LIST-014 | PartitionInPlace | `func PartitionInPlace[T cmp.Ordered](head *Node[T], pivot T) (*Node[T], error)` |
| ALG-LIST-015 | LRUCache | `func NewLRUCache[K comparable, V any](capacity int) (*LRUCache[K, V], error)`; Get, Peek, Put, Delete, Len, Cap |
| ALG-LIST-016 | LFUCache | `func NewLFUCache[K comparable, V any](capacity int) (*LFUCache[K, V], error)`; Get, Put, Delete, Len, Cap |
| ALG-LIST-017 | AddTwoNumbers | `func AddTwoNumbers(a, b *Node[int], radix int) (*Node[int], error)` |

Errors (package `list`): `ErrCycle`, `ErrSharedNodes`, `ErrGroupSize`, `ErrPosition`, `ErrCapacity`,
`ErrKeyNotComparable`, `ErrRadix`, `ErrDigit`. nil is the empty list everywhere; no function panics on a
nil list, a cyclic list, a nil receiver or a foreign Element.

### ALG-LIST-001 Node · 002 FromSlice · 003 AppendValues · 004 Length
- **Source**: Knuth, *The Art of Computer Programming*, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.3 (linked
  allocation: INFO and LINK fields, Λ the null link). Cormen et al., *Introduction to Algorithms*, 3rd ed.,
  MIT Press, 2009, §10.2 (singly linked lists).
- **Best known**: building or reading n nodes Θ(n). Length of a list that may be cyclic: Θ(n) time, O(1) space
  with Floyd's test (ALG-LIST-008's source).
- **Approach**: FromSlice allocates the n nodes as one block (one allocation) and links them. AppendValues and
  Length walk with the Floyd guard (B5-2). **Lost**: one allocation per node (n allocations; slower to build,
  nothing gained for a library type the caller may relink anyway; documented that the block stays alive
  while any of its nodes is reachable).
- **Input modified**: never. **Errors**: `ErrCycle` (AppendValues returns dst unchanged).
- **Worked example**: FromSlice([1 2 3]) → 1→2→3→nil; Length 3; AppendValues → [1 2 3]. Hand-derived.
- **Tags**: linked-list, singly-linked.

### ALG-LIST-005 DoublyLinkedList
- **Source**: Knuth Vol. 1, 3rd ed., §2.2.5 (doubly linked lists with a list head: LLINK/RLINK, deletion and
  insertion in O(1) given the node). CLRS 3rd ed., §10.2 (LIST-INSERT, LIST-DELETE, sentinels, Figure 10.4).
- **Best known**: O(1) worst case per insert/remove/move given an element; O(n) to iterate; O(n) space.
- **Approach**: circular list around a sentinel element embedded in the list value (lazy init: the zero
  value is an empty list). Each Element records its list; operations given an Element of another list (or a
  removed one, or nil) do nothing and report false/nil. `All`/`Backward` are `iter.Seq[T]`; the next element
  is read before yield, so removing the current element during iteration is safe. **Lost**: `container/list`
  (B5-3); a non-circular list with nil ends (two more branches per operation).
- **Errors**: none; Remove/Move report ok=false on a foreign element. Not safe for concurrent use; must not be
  copied after first use (the sentinel's address is the list's identity).
- **Worked example**: PushBack 1,2,3; PushFront 0; MoveToBack(front) → 1 2 3 0; Remove(back) → 1 2 3. Hand-derived.
- **Tags**: linked-list, doubly-linked, deque, data-structure. LeetCode 707 (pattern only).

### ALG-LIST-006 ReverseInPlace
- **Source**: Knuth Vol. 1, 3rd ed., §2.2.3, Exercise 7 ("invert" a linked linear list by changing its links;
  confirmed 2026-11-12). CLRS 3rd ed., Exercise 10.2-7 (Θ(n) nonrecursive reversal using constant storage).
- **Best known**: Θ(n) time (every link changes), O(1) extra space.
- **Approach**: one pass with three pointers after the Floyd guard. **Lost**: recursion (O(n) stack depth,
  unbounded on a long list); copying into a slice and relinking (O(n) space).
- **Input modified**: yes (links). **Errors**: `ErrCycle` (list untouched). nil → nil.
- **Worked example**: 1→2→3→4 → 4→3→2→1. Hand-derived from Knuth's description.
- **Tags**: linked-list, reverse, in-place. LeetCode 206.

### ALG-LIST-007 ReverseKGroupInPlace
- **Source**: same as 006 (inversion applied to each consecutive block of k nodes); Knuth Vol. 1 §2.2.3 Ex. 7.
- **Best known**: Θ(n), O(1) extra.
- **Approach**: length first (guarded), then ⌊n/k⌋ block reversals splicing each block's tail to the next
  block; a final block of fewer than k nodes keeps its order. **Lost**: recursion per block (O(n/k) stack);
  a stack of k nodes (O(k) space).
- **Errors**: `ErrGroupSize` (k < 1), `ErrCycle`. k = 1 or k > n: unchanged.
- **Worked example**: 1→2→3→4→5, k = 2 → 2→1→4→3→5; k = 3 → 3→2→1→4→5. Hand-derived.
- **Tags**: linked-list, reverse, in-place. LeetCode 25, 24 (k = 2).

### ALG-LIST-008 DetectCycle · ALG-LIST-009 CycleStart
- **Source**: Knuth, *The Art of Computer Programming*, Vol. 2, 3rd ed., Addison-Wesley, 1997, §3.1,
  Exercise 6 (find μ and λ in O(μ + λ) steps with a bounded number of memory locations; Knuth ascribes the
  X_n = X_2n method to R. W. Floyd) and its answer. Confirmed 2026-11-12 (exercise statement). Brent's
  power-of-two variant is Exercise 7 (R. P. Brent, 1977).
- **Best known**: Θ(μ + λ) time, O(1) space; μ = nodes before the cycle, λ = cycle length. Lower bound:
  every node must be read (Ω(μ + λ)).
- **Approach**: Floyd's tortoise and hare; CycleStart then restarts one pointer at head and steps both by one
  until they meet (the entry, μ steps), then walks the cycle once (λ). **Lost**: Brent's method (fewer
  evaluations of f, which matters for an expensive f; here f is a pointer read and Floyd's meeting point is
  needed for the entry anyway); a visited set (O(n) space).
- **Errors**: none. DetectCycle(nil) false; CycleStart on an acyclic list → (nil, 0).
- **Worked example**: 1→2→3→4→5→(back to 3): μ = 2, λ = 3, start = node 3. Self loop: λ = 1. Hand-derived.
- **Tags**: linked-list, cycle-detection, floyd, two-pointer, tortoise-hare. LeetCode 141, 142, 287 (pattern).

### ALG-LIST-010 Middle
- **Source**: Floyd's two-speed walk (Knuth Vol. 2 §3.1 Ex. 6, as above): when the fast pointer runs off the
  end after k double steps, the slow pointer is at index k = ⌊n/2⌋. Laaksonen, *Guide to Competitive
  Programming*, Springer (L18) — not needed; Knuth suffices.
- **Best known**: Θ(n) (the end must be found), O(1) space; one pass.
- **Approach**: slow/fast; returns the node at index ⌊n/2⌋ (0-based) — for even n the second of the two
  middles. A meeting of the pointers proves a cycle (`ErrCycle`), so no separate guard pass. **Lost**: count
  then walk n/2 (1.5 passes as well, but needs its own cycle guard).
- **Errors**: `ErrCycle`. nil → nil.
- **Worked example**: 1→2→3→4→5 → 3; 1→2→3→4 → 3. Hand-derived.
- **Tags**: linked-list, middle, two-pointer. LeetCode 876.

### ALG-LIST-011 MergeTwoSortedInPlace
- **Source**: Knuth, *TAOCP* Vol. 3, 2nd ed., Addison-Wesley, 1998, §5.2.4, Algorithm L (list merge sort; its
  merge step relinks two sorted linked lists). CLRS 3rd ed., §2.3.1 (MERGE).
- **Best known**: Θ(m + n) time, ≤ m + n − 1 comparisons, O(1) extra space; lower bound for merging equal
  sizes is 2n − 1 comparisons (Knuth Vol. 3 §5.3.2), met.
- **Approach**: a dummy head, take from a while !(b < a): stable (equal values: a's before b's).
  Values compared with `cmp.Less` (NaN first, D-002). **Lost**: recursion (O(m + n) stack).
- **Input modified**: both lists are relinked into the result. Unsorted input gives a merged permutation, no
  error. **Errors**: `ErrCycle`, `ErrSharedNodes` (the two lists share a node: their last nodes coincide —
  merging would build a cycle).
- **Worked example**: 1→3→5 and 1→2→4 → 1(a)→1(b)→2→3→4→5. Hand-derived.
- **Tags**: linked-list, merge, sorted, stable, in-place. LeetCode 21.

### ALG-LIST-012 RemoveNthFromEndInPlace
- **Source**: Knuth Vol. 1 §2.2.3 (deleting a node by relinking its predecessor) + the length of the list.
- **Best known**: Θ(L) time, O(1) space. One pass with a gap of n between two pointers is possible, but the
  cycle guard reads the list once anyway, so length-then-walk is the same order with simpler code. **Lost**:
  the gap pointer (it alone cannot detect a cycle).
- **Input modified**: yes; the removed node's Next is set to nil. `fromEnd` is 1-based (1 = last node).
- **Errors**: `ErrPosition` (fromEnd < 1 or > length), `ErrCycle`.
- **Worked example**: 1→2→3→4→5, fromEnd 2 → 1→2→3→5. Hand-derived.
- **Tags**: linked-list, two-pointer, in-place. LeetCode 19.

### ALG-LIST-013 ReorderInPlace
- **Source**: composition of cited steps: the middle (ALG-LIST-010's walk), inversion of the second half
  (Knuth Vol. 1 §2.2.3 Ex. 7), and alternating merge of two lists by relinking (Knuth Vol. 1 §2.2.3).
  Proof of the order: first half L0…L⌈n/2⌉−1, reversed second half Ln−1…L⌈n/2⌉; alternating them gives
  L0, Ln−1, L1, Ln−2, …, which is the target order.
- **Best known**: Θ(n), O(1) extra. **Lost**: copying nodes into a slice and relinking by index (O(n) space).
- **Errors**: `ErrCycle`. nil, one or two nodes: unchanged.
- **Worked example**: 1→2→3→4 → 1→4→2→3; 1→2→3→4→5 → 1→5→2→4→3. Hand-derived.
- **Tags**: linked-list, reorder, two-pointer, in-place. LeetCode 143.

### ALG-LIST-014 PartitionInPlace
- **Source**: stable partition of a linked list into two sublists by relinking (Knuth Vol. 1 §2.2.3, linked
  allocation; CLRS 3rd ed. §7.1 PARTITION for the predicate; a linked list allows the stable form in O(1)
  space, unlike the array form).
- **Best known**: Θ(n), O(1) extra.
- **Approach**: two dummy-headed tails (values < pivot, values ≥ pivot by `cmp.Less`), each keeping input
  order, then concatenated. NaN pivot: nothing is less than NaN under cmp.Less, so all go to the second part.
- **Errors**: `ErrCycle`.
- **Worked example**: 1→4→3→2→5→2, pivot 3 → 1→2→2→4→3→5. Hand-derived.
- **Tags**: linked-list, partition, stable, in-place. LeetCode 86.

### ALG-LIST-015 LRUCache
- **Source**: Mattson, Gecsei, Slutz, Traiger, "Evaluation techniques for storage hierarchies", IBM Systems
  Journal 9(2):78–117, 1970, doi:10.1147/sj.92.0078 (LRU replacement as a stack algorithm; verified on
  Crossref 2026-11-12). Implementation: hash table (CLRS 3rd ed., ch. 11) + doubly linked list in recency
  order (CLRS §10.2; this package's ALG-LIST-005).
- **Best known**: O(1) per Get/Put/Delete (expected, hash map), O(capacity) space. **Lost**: a timestamp per
  entry with an O(n) scan or a heap (O(log n)) for the victim; `container/list` (B5-3).
- **Approach**: `map[K]*Element[entry]` + DoublyLinkedList, front = most recent. Put on a full cache reuses the
  evicted element (no allocation in steady state). Get marks the key most recent; Peek does not.
  `Put(key, value) (evicted Entry[K, V], didEvict bool, err error)`.
- **Errors**: `ErrCapacity` (capacity < 1), `ErrKeyNotComparable` (B5-6). Zero value or nil: capacity 0,
  Get misses, Put stores nothing and reports the given pair as evicted. No clock, no goroutines; not safe for
  concurrent use.
- **Worked example**: capacity 2: Put(1,a) Put(2,b) Get(1)=a Put(3,c) evicts 2; Get(2) miss; Put(4,d)
  evicts 1. Hand-derived (Mattson's LRU stack order).
- **Tags**: cache, lru, eviction, hash-map, linked-list. LeetCode 146.

### ALG-LIST-016 LFUCache
- **Source**: Shah, Mitra, Matani, "An O(1) algorithm for implementing the LFU cache eviction scheme",
  technical report, 16 August 2010, arXiv:2110.11602 (posted 2021; read 2026-11-12): §5, a doubly linked
  *frequency list* of nodes, each holding the items with that access count, a hash table by key, and every
  item pointing to its frequency node; insert, access and evict O(1).
- **Best known**: O(1) per operation (the paper), O(capacity) space. **Lost**: min-heap on (count, recency)
  with a map (O(log n) per op, the paper's §4 baseline; algo/heap has no decrease-key).
- **Approach**: the paper's structure, with the per-frequency SET replaced by a recency-ordered linked list
  (B5-7): the victim is the least recently used item of the lowest count. A bucket whose only item moves up
  to an absent count+1 is renumbered in place (no allocation). The count saturates at 2^64−1.
- **Errors**: `ErrCapacity`, `ErrKeyNotComparable`. Zero value/nil as LRUCache. Put on an existing key replaces
  the value and counts as a use.
- **Worked example**: capacity 2: Put(1) Put(2) Get(1) Put(3) evicts 2 (count 1 vs 2); Get(3) (count 2);
  Put(4) evicts 1 (counts tie at 2, 1 used less recently). Hand-derived.
- **Tags**: cache, lfu, eviction, hash-map, linked-list. LeetCode 460.

### ALG-LIST-017 AddTwoNumbers
- **Source**: Knuth, *TAOCP* Vol. 2, 3rd ed., 1997, §4.3.1, Algorithm A (addition of nonnegative integers in
  radix b, least significant digit first, carry k ∈ {0, 1}).
- **Best known**: Θ(max(m, n)) time; output max(m, n) or +1 digits. O(1) extra besides the result.
- **Approach**: lengths first (guard), result allocated as one block of max(m, n)+1 nodes; a shorter operand is
  read as zeros (Knuth assumes equal length n; padding with zeros is the same number). Digit sums in uint
  (no overflow for any radix ≤ MaxInt). nil is zero; nil + nil = nil.
- **Input modified**: never. **Errors**: `ErrRadix` (radix < 2), `ErrDigit` (digit < 0 or ≥ radix), `ErrCycle`.
- **Worked example**: 342 + 465 in radix 10, as 2→4→3 + 5→6→4 → 7→0→8 (807); 99 + 1 → 0→0→1. Hand-derived.
- **Tags**: linked-list, big-integer, addition, carry. LeetCode 2, 445 (pattern).
