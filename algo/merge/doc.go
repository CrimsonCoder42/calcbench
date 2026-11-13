// Package merge combines sorted sequences into one sorted sequence.
//
// ID short form: MERGE (IDs are ALG-MERGE-NNN).
//
// Ordering: values are compared with cmp.Compare, a total order in which a floating-point NaN sorts
// before every other value and -0 equals +0 (doc/DECISIONS.md D-002). Inputs must be ascending in that
// order (a slice sorted by slices.Sort is). Sortedness is a precondition and is not checked: on unsorted
// input every element still appears exactly once in the result, and nothing panics. No function
// modifies its input.
//
// Functions:
//
//	ALG-MERGE-001  MergeKSorted  one ascending slice from k ascending slices, stable
package merge
