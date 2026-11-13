// Package hashing counts, groups and compares collections with hash tables.
//
// ID short form: HASH (IDs are ALG-HASH-NNN).
//
// The package is named hashing, not hash, so that callers can import the standard hash package beside it.
//
// Cost: Go maps are hash tables. Every Time below is expected (average) time under the uniform hashing
// assumption (Cormen et al., Introduction to Algorithms, 3rd ed., ch. 11); one map operation is O(m) in
// the worst case, when all m keys collide.
//
// Order: Go map iteration order is unspecified (the Go language specification, "For statements with range
// clause"), and it never reaches an output here. Groups and set results follow the order of first
// appearance in the input; ranked results are ordered by count descending, then value ascending.
//
// Equality is Go's ==. A floating-point NaN is not equal to itself, so each NaN is its own key: it never
// matches another value and is never merged with another NaN. -0 and +0 are one key.
// No function modifies its input.
//
// Functions:
//
//	ALG-HASH-001  FrequencyCount           how many times each distinct value occurs
//	ALG-HASH-002  TopKFrequent             the k most frequent values, with their counts
//	ALG-HASH-003  TopKByCount              the k largest entries of a count map
//	ALG-HASH-004  GroupBy                  values grouped by a key function, in first-appearance order
//	ALG-HASH-005  FirstUnique              index of the first value that occurs exactly once
//	ALG-HASH-006  SetIntersection          distinct values of a that occur in b
//	ALG-HASH-007  SetUnion                 distinct values of a, then of b
//	ALG-HASH-008  SetDifference            distinct values of a that do not occur in b
//	ALG-HASH-009  LongestConsecutive       longest run of consecutive integers present
//	ALG-HASH-010  GroupAnagrams            words grouped by sorted-rune signature
//	ALG-HASH-011  GroupAnagramsByCount     words grouped by rune counts over an alphabet
//	ALG-HASH-012  TwoSum                   indices of two values that sum to a target
//	ALG-HASH-013  ContainsNearbyDuplicate  whether equal values occur within a given index gap
//	ALG-HASH-014  BinCount                 dense count of each integer 0…bins-1 (direct-address table)
//	ALG-HASH-015  BinCountWeighted         compensated sum of weights per integer 0…bins-1
//	ALG-HASH-016  IsIn                     for each value, whether it occurs in a test set
//	ALG-HASH-017  SetSymmetricDifference   distinct values in exactly one of a and b
package hashing
