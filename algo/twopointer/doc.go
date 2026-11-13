// Package twopointer holds algorithms that walk a slice with two indexes, each moving in one direction
// only, so a whole pass costs O(n): pair and triple sums on sorted data, in-place deduplication,
// merging, partitioning and palindrome checks over slices.
//
// ID short form: TWOPTR (IDs are ALG-TWOPTR-NNN).
//
// Modification: a function whose name ends in InPlace modifies its slice and says how in its comment;
// every other function only reads its input. Variable-length sliding windows (shortest or longest run
// with a property) live in algo/window; string forms (palindromic substrings and the like) live in
// algo/text.
//
// Numbers: TwoSumSorted and ThreeSum take any built-in integer type and compare sums exactly in 128-bit
// arithmetic, so overflow never hides or invents an answer. MergeSortedInPlace orders by cmp.Compare
// (NaN first, doc/DECISIONS.md D-002); the others use ==.
//
// Functions:
//
//	ALG-TWOPTR-001  TwoSumSorted             indexes of two elements of a sorted slice with a given sum
//	ALG-TWOPTR-002  ThreeSum                 every distinct triple with a given sum
//	ALG-TWOPTR-003  RemoveDuplicatesInPlace  drop adjacent duplicates, return the new length
//	ALG-TWOPTR-004  MergeSortedInPlace       merge a sorted slice into a buffer's sorted prefix, back to front
//	ALG-TWOPTR-005  IsPalindrome             whether a slice equals its reversal
//	ALG-TWOPTR-006  MoveZeroesInPlace        move zero values to the end, keeping the others' order
package twopointer
