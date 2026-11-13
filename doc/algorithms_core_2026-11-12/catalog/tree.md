# algo/tree (SHORT = TREE) — S2 wave 2 (builder W2-B6)

New package (planned in doc/taxonomy.md: `algo/tree`, TREE, "balanced trees, traversals"). No balanced
(red-black/AVL) type is built; `SortedToBalancedBST` gives a perfectly balanced BST from sorted keys.

Sources confirmed 2026-11-12:
- **CLRS**: Cormen, Leiserson, Rivest, Stein, *Introduction to Algorithms*, 3rd ed., MIT Press, 2009, ch. 12
  (§12.1 BST property and INORDER-TREE-WALK, Theorem 12.1; §12.2 TREE-SEARCH, TREE-MINIMUM; §12.3
  TREE-INSERT, TRANSPLANT, TREE-DELETE), §22.2 (breadth-first search), Appendix B.5.3 (height of a tree).
- **Hibbard** 1962: T. N. Hibbard, "Some combinatorial properties of certain trees with applications to
  searching and sorting", J. ACM 9(1):13–28, doi:10.1145/321105.321108 (Crossref checked).
- **Knuth** TAOCP Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.3.1: Algorithm T (inorder with a stack),
  exercise 12 (preorder analogue), exercise 13 (postorder: a second state per stack entry), exercise 7
  (preorder + inorder determine the tree). Exercise numbers checked against an index of §2.3.1 exercises.
- **Laaksonen**, *Competitive Programmer's Handbook*, draft of 3 July 2018, https://cses.fi/book/book.pdf
  (read 2026-11-12; textbook per lead decision L18): §9.2 Binary indexed tree pp. 86–88; §9.3 Segment tree
  pp. 89–93; §9.4 "Range updates" pp. 93–94; §14.2 Diameter, Algorithm 1 pp. 135–136; §14.4 Binary trees
  p. 139; §18.3 Lowest common ancestor p. 167; §28.1 Lazy propagation pp. 258–260.
- **Fenwick** 1994: P. M. Fenwick, "A new data structure for cumulative frequency tables", Software:
  Practice and Experience 24(3):327–336, doi:10.1002/spe.4380240306 (Crossref checked).
- **Wirth**, *Algorithms and Data Structures*, Oberon version, 2004 (author's PDF AD2012), §4.4.1,
  pp. 151–152 (perfectly balanced tree: root, n DIV 2 nodes left, n − n DIV 2 − 1 right).
- **Aho, Hopcroft, Ullman**, "On finding lowest common ancestors in trees", SIAM J. Comput. 5(1):115–132,
  1976, doi:10.1137/0205011 (Crossref checked).
- **Bernstein**, "Netstrings", 1997, https://cr.yp.to/proto/netstrings.txt (also IETF
  draft-bernstein-netstrings-02): `len ":" bytes ","`, no leading zeros.
- **Go spec**, "For statements with range clause" (go1.27 doc/go_spec.html, #For_range): "If the loop body
  terminates (such as by a break statement), yield returns false and must not be called again."

Shared worked tree (Laaksonen §14.4 p. 139): root 1; 1 → (2, 3); 2 → (4, 5); 5 → (6, nil); 3 → (nil, 7).
Preorder 1 2 4 5 6 3 7 · inorder 4 2 6 5 1 3 7 · postorder 4 6 5 2 7 3 1.

Conventions: every traversal is iterative (explicit stack or queue), so a degenerate tree of 1e5 nodes is
safe; segment-tree recursion is over the index range, depth ⌈log2 n⌉ + 1 ≤ 64. The input to every tree
function must be a tree (no node reachable twice); a cyclic structure is a precondition breach (documented,
not detected). Keys compare with cmp.Compare (NaN first, D-002). Integer sums are overflow-checked
(ErrOverflow); floats reject NaN/±Inf (ErrNotFinite). Every mutating method is atomic: on error the
structure is unchanged.

## Types
`Node[T any]{Value T; Left, Right *Node[T]}` — the general binary tree. `BST[T cmp.Ordered]` — zero value is
an empty set (no duplicates). `FenwickTree[T]`, `FenwickRangeUpdate[T]` (T signed integer or float),
`SegmentTree[T any]`, `LazySegmentTree[T, U any]` with `LazyOps[T, U]`.

## BST — ALG-TREE-001 Insert · 002 Search · 003 Delete · 004 Len · 005 Root
- **Source**: CLRS §12.2–12.3; Hibbard 1962 (deletion by successor). Iterative.
- **Best known**: O(h) per operation for an unbalanced BST (h = height; n worst, O(log n) expected for random
  insertion order, CLRS §12.4). Lower bound Ω(log n) comparisons in the worst case for any comparison
  search tree. **Lost**: a self-balancing tree (red-black, CLRS ch. 13) — not requested; recorded for later.
- **Approach**: CLRS TREE-DELETE relinks the successor node (no value copy, so outstanding *Node pointers
  stay valid). **Lost**: Hibbard's original copy-successor-key form (invalidates node identity).
- **Worked example**: CLRS Figure 12.2 shape built by inserting 15 6 18 3 7 17 20 2 4 13 9 (hand-checked
  inorder 2 3 4 6 7 9 13 15 17 18 20), then deleting 6 (two children: successor 7), 13 (one child), 2 (leaf).
- **Errors**: none; Insert false for a duplicate, Delete false when absent; nil *BST safe.

## Traversals — ALG-TREE-006 Inorder · 007 Preorder · 008 Postorder · 009 LevelOrder (iter.Seq[T])
- **Source**: Knuth §2.3.1 Algorithm T, ex. 12, ex. 13; CLRS §22.2 (BFS by levels); Go spec #For_range.
- **Best known**: Θ(n) time; O(h) extra for depth-first with a stack, O(w) for level order (w = max level
  width). Morris threading gives O(1) extra but writes into the tree temporarily — **lost**: the input must
  not be modified and an early break would leave threads behind.
- **Early break**: each iterator returns immediately when yield returns false (test counts calls).
- **Worked example**: the shared tree's three orders; level order 1 2 3 4 5 7 6.

## ALG-TREE-010 Height · 011 IsBalanced
- Height = edges on the longest root-to-leaf path (CLRS B.5.3): empty −1, one node 0 (MaxDepth = Height+1).
- IsBalanced: every node's subtree heights differ by at most 1 (Adel'son-Vel'skii–Landis condition as
  stated in Knuth TAOCP Vol. 3, 2nd ed., 1998, §6.2.3). Computed by one iterative postorder fold, O(n) time,
  O(h) extra; stops at the first unbalanced node. **Lost**: recursive height per node (O(n²) on a chain).
- Worked example (hand-checked): shared tree height 3 (path 1 2 5 6); balanced (node 1: 2 vs 1, node 2:
  0 vs 1, node 5: 0 vs −1, node 3: −1 vs 0); the chain 1→2→3 is not (node 1: −1 vs 1).

## ALG-TREE-012 LowestCommonAncestorBST · 013 LowestCommonAncestor
- **Source**: definition AHU 1976; Laaksonen §18.3. General tree: one postorder fold counting p and q;
  the first node whose subtree holds both is the LCA. O(n), O(h) extra. **Lost**: Euler tour + RMQ or
  binary lifting (O(n log n) preprocessing, worth it only for many queries; a future type).
- BST: the LCA of keys p ≤ q is the first node on the search path with p ≤ key ≤ q (proof: above it both
  keys lie on one side, so the node is a common ancestor's descendant on the shared path; at it they split
  or one equals it, so no deeper node is an ancestor of both). Search via CLRS §12.2; both keys then
  confirmed present. O(h), O(1) extra.
- Worked (hand-checked on the insertion tree 15 6 18 3 7 17 20 2 4 13 9): LCA(2, 4) = 3, LCA(9, 17) = 15,
  LCA(13, 9) = 13, LCA(4, 9) = 6; the shared tree (general): LCA(4, 6) = 2, LCA(6, 7) = 1.

## ALG-TREE-014 FromPreIn · 015 FromInPost
- **Source**: Knuth §2.3.1 ex. 7; Laaksonen §14.4 p. 139 (pre+in and post+in determine the tree; pre+post
  does not). T comparable, values distinct.
- **Best known**: Θ(n). **Approach**: one stack pass (a node becomes a right child when the inorder cursor
  matches the stack top), nodes in one allocation, then both traversals of the result are compared with the
  inputs (O(n)), so any inconsistent input is rejected. **Lost**: recursive split with an index map
  (O(n) but recursion depth n on a chain).
- **Errors**: ErrLength (lengths differ), ErrDuplicate, ErrInconsistent (not the traversals of one tree;
  a NaN never matches itself, so float NaN input is ErrInconsistent).
- Worked example: shared tree from (pre, in) and from (in, post).

## ALG-TREE-016 SortedToBalancedBST
- **Source**: Wirth §4.4.1 pp. 151–152 (perfectly balanced: left gets n DIV 2 nodes). Θ(n) time, one
  allocation of n nodes, O(log n) extra stack (explicit). Errors: ErrNotSorted (not strictly ascending).
- Worked: 1..7 → root 4, (2: 1, 3), (6: 5, 7).

## ALG-TREE-017 KthSmallestBST · 018 ValidateBST
- **Source**: CLRS §12.1 Theorem 12.1 (inorder walk lists a BST's keys in sorted order). Kth: O(h + k),
  O(h) extra, stops early. Errors ErrRange (k < 1 or k > n). **Lost**: size-augmented nodes (O(h) per
  query; needs a different node type). Validate: strictly ascending inorder, O(n), O(h); ErrNotBST names
  the first pair out of order.

## ALG-TREE-019 Serialize · 020 Deserialize
- **Format (ours)**: the preorder sequence with every empty subtree written as `#` and every node's value
  as a netstring (Bernstein 1997): `len:bytes,`. Empty tree `#`; a leaf 7 encoded as "7" is `1:7,##`.
  Uniqueness: preorder with explicit empty markers is a prefix code of extended binary trees (each node
  token is followed by exactly two subtree encodings), so parsing is unambiguous; netstring framing makes
  any byte value safe. Values are encoded by the caller's appendValue / parseValue.
- Θ(n + total value bytes); O(h) extra stack, iterative both ways. Errors: ErrNilFunc, ErrMalformed
  (bad length, leading zero, missing comma, truncated, trailing bytes), parse errors wrapped.

## ALG-TREE-021 MaxPathSum
- **Source**: Laaksonen §14.2 Algorithm 1 (every path has a highest node; per node the best downward
  path and the best path through it), weighted, a downward gain clipped at 0. Θ(n), O(h) extra.
- Overflow: only a positive overflow is possible (a node plus non-negative gains), which means the true
  maximum does not fit T → ErrOverflow. ErrEmptyTree, ErrNotFinite.
- Worked: −10 → (9, 20 → (15, 7)) gives 42 (15+20+7), hand-checked.

## ALG-TREE-022 RightSideView
- Last value of each level in level order (CLRS §22.2). Θ(n), O(w) extra; appends to dst.
- Worked: shared tree → 1 3 7 6.

## Fenwick — ALG-TREE-023 NewFenwickTree · 024 Add · 025 PrefixSum · 026 RangeSum · 027 Len
## ALG-TREE-028 NewFenwickRangeUpdate · 029 RangeAdd · 030 Get · 031 Len
- **Source**: Fenwick 1994; Laaksonen §9.2 (tree[k] = sum of (k − p(k), k], p(k) = k & −k) and §9.4
  "Range updates" (the tree over the difference array: range add = two point adds, point value = prefix).
- **Best known**: O(log n) update and prefix query; Ω(log n / log log n) cell-probe lower bound for
  dynamic prefix sums (Fredman–Saks 1989, not re-read; recorded as a pointer). Build O(n): each node adds
  itself to its parent k + p(k) once (follows from the definition).
- **Decision**: range-update + point-query variant built (sourced, §9.4). Range-update + range-query by two
  Fenwick trees **not built** (no primary source read for it); that need is met by the lazy segment tree
  (ALG-TREE-039…045), sourced in §28.1.
- RangeSum walks the two index chains until they meet (equivalent to PrefixSum(hi) − PrefixSum(lo), fewer
  terms, fewer spurious overflows).
- T = signed integers and floats (the range variant subtracts). ErrRange, ErrNotFinite, ErrOverflow
  (a partial sum the tree forms does not fit T; can happen for a range whose total fits — use a wider T).
- Worked: Laaksonen p. 87: [1 3 4 8 6 1 4 2] → tree 1 4 4 16 6 7 4 29; sum of the first 7 = 27.
  §9.4 p. 93–94: [3 3 1 1 1 5 2 2], add 5 to positions 1..4 → 3 8 6 6 6 5 2 2.

## Segment trees — ALG-TREE-032 NewSegmentTree · 033 NewSumSegmentTree · 034 NewMinSegmentTree
## 035 NewMaxSegmentTree · 036 Set · 037 Query · 038 Len
## Lazy — ALG-TREE-039 NewLazySegmentTree · 040 NewRangeAddSumTree · 041 NewRangeAddMinTree
## 042 NewRangeAddMaxTree · 043 RangeUpdate · 044 Query · 045 Len
- **Source**: Laaksonen §9.3 (segment tree; any associative operation), §28.1 (lazy propagation: node
  value plus a pending update, combined updates compose). Generic: a monoid (associative combine +
  identity; combine need not commute) and, for lazy, an update type U with Apply(u, aggregate, length),
  Compose(newer, older) and NoUpdate.
- **Best known**: O(log n) point update and range query, O(log n) range update with lazy propagation;
  O(n) build. Layout: 2n − 1 nodes (left child i+1, right child i + 2·(mid − lo)), not 4n.
  **Lost**: the bottom-up iterative tree over 2n slots (needs commutativity or padding to a power of two
  plus an identity for the padding; min/max ready-mades have no natural identity for every T) and the
  4n heap layout (twice the memory).
- Query on the lazy tree is read-only (pending updates are composed on the way down, not pushed), so
  concurrent readers do not race with each other. Set and RangeUpdate journal the nodes they write and
  roll back on overflow (atomic).
- Empty range [lo, lo): the identity for NewSegmentTree/NewLazySegmentTree/sum trees; ErrEmptyRange for
  min/max ready-mades (no identity for every T).
- ErrRange, ErrNilFunc, ErrNotFinite (number ready-mades), ErrOverflow, ErrEmptyRange.
- Fix M3 (review): the range-add min/max trees keep both extremes per node (min and max), so an add
  that would push any element out of T is rejected atomically, not only one that moves the node's aggregate.
- Worked: Laaksonen §28.1 pp. 258–259: [5 8 6 3 2 7 2 6 7 1 7 5 6 2 3 2] sums to 72; adding 2 to positions
  5..13 makes the total 90 (root 90/0 in the book's figure) and the left half 45, right half 45.

Tags (LeetCode, pattern only): 94, 144, 145, 102, 104, 110, 235, 236, 105, 106, 108, 230, 98, 297, 124, 199,
307, 315, 450, 701.
