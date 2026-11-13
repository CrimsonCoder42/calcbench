// Package tree holds binary trees and the tree-shaped range-query structures: a generic binary tree node,
// an unbalanced binary search tree, traversal iterators, tree queries (height, balance, lowest common
// ancestor, path sums, views), reconstruction from traversals, a text serialization, Fenwick trees and
// segment trees (generic monoid, with lazy propagation for range updates).
//
// ID short form: TREE (IDs are ALG-TREE-NNN).
//
// Conventions: every walk over a *Node is iterative (an explicit stack or queue), so a degenerate tree
// of any depth is safe; segment-tree recursion is over the index range, depth ≤ ⌈log2 n⌉ + 1 ≤ 64.
// A *Node argument must be a tree (no node reachable twice); a cycle is not detected. Keys compare with
// cmp.Compare (a NaN sorts first, -0 equals +0; doc/DECISIONS.md D-002). Integer sums are exact or
// ErrOverflow, never wrapped; float inputs to sums must be finite (ErrNotFinite). No function modifies
// its input; every method that changes a structure is atomic (on error nothing changes). Ranges are
// half-open [lo, hi) and 0-based. Iterators stop at once when the range loop breaks.
//
// Functions:
//
//	ALG-TREE-001  BST.Insert  add a key, false if present
//	ALG-TREE-002  BST.Search  whether a key is present
//	ALG-TREE-003  BST.Delete  remove a key (Hibbard deletion by successor), false if absent
//	ALG-TREE-004  BST.Len  the number of keys
//	ALG-TREE-005  BST.Root  the root node, for the traversals and queries
//	ALG-TREE-006  Inorder  iterator over values in inorder
//	ALG-TREE-007  Preorder  iterator over values in preorder
//	ALG-TREE-008  Postorder  iterator over values in postorder
//	ALG-TREE-009  LevelOrder  iterator over values level by level
//	ALG-TREE-010  Height  edges on the longest root-to-leaf path (−1 when empty)
//	ALG-TREE-011  IsBalanced  whether subtree heights differ by at most 1 everywhere
//	ALG-TREE-012  LowestCommonAncestorBST  the LCA of two keys in a BST
//	ALG-TREE-013  LowestCommonAncestor  the LCA of two nodes in a binary tree
//	ALG-TREE-014  FromPreIn  the tree with a given preorder and inorder
//	ALG-TREE-015  FromInPost  the tree with a given inorder and postorder
//	ALG-TREE-016  SortedToBalancedBST  a perfectly balanced BST from strictly ascending keys
//	ALG-TREE-017  KthSmallestBST  the k-th smallest key of a BST
//	ALG-TREE-018  ValidateBST  nil when a tree is a BST, ErrNotBST otherwise
//	ALG-TREE-019  Serialize  append a tree in preorder with '#' for empty subtrees and netstring values
//	ALG-TREE-020  Deserialize  the tree Serialize wrote
//	ALG-TREE-021  MaxPathSum  the largest sum along any path
//	ALG-TREE-022  RightSideView  the last value of each level
//	ALG-TREE-023  NewFenwickTree  a Fenwick tree over a copy of values, built in O(n)
//	ALG-TREE-024  FenwickTree.Add  add to one element
//	ALG-TREE-025  FenwickTree.PrefixSum  the sum of the first k elements
//	ALG-TREE-026  FenwickTree.RangeSum  the sum of [lo, hi)
//	ALG-TREE-027  FenwickTree.Len  the number of elements
//	ALG-TREE-028  NewFenwickRangeUpdate  range-add, point-query structure over a copy of values
//	ALG-TREE-029  FenwickRangeUpdate.RangeAdd  add a constant to [lo, hi)
//	ALG-TREE-030  FenwickRangeUpdate.Get  one element
//	ALG-TREE-031  FenwickRangeUpdate.Len  the number of elements
//	ALG-TREE-032  NewSegmentTree  a segment tree over any monoid (combine + identity)
//	ALG-TREE-033  NewSumSegmentTree  a segment tree of sums
//	ALG-TREE-034  NewMinSegmentTree  a segment tree of minimums
//	ALG-TREE-035  NewMaxSegmentTree  a segment tree of maximums
//	ALG-TREE-036  SegmentTree.Set  replace one element
//	ALG-TREE-037  SegmentTree.Query  the combined value of [lo, hi)
//	ALG-TREE-038  SegmentTree.Len  the number of elements
//	ALG-TREE-039  NewLazySegmentTree  a lazy segment tree defined by LazyOps
//	ALG-TREE-040  NewRangeAddSumTree  range add, range sum
//	ALG-TREE-041  NewRangeAddMinTree  range add, range minimum
//	ALG-TREE-042  NewRangeAddMaxTree  range add, range maximum
//	ALG-TREE-043  LazySegmentTree.RangeUpdate  apply an update to [lo, hi)
//	ALG-TREE-044  LazySegmentTree.Query  the combined value of [lo, hi)
//	ALG-TREE-045  LazySegmentTree.Len  the number of elements
package tree
