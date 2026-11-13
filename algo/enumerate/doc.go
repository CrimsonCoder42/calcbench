// Package enumerate generates combinatorial patterns in order: permutations, combinations, subsets,
// multiset variants, combinations with a fixed sum, nested parentheses and Gray codes. It lists them;
// counting them (n!, C(n,k), Catalan numbers) is num/combinatorics.
//
// ID short form: ENUM (IDs are ALG-ENUM-NNN).
//
// Iterators: every generator returns an iter.Seq. Breaking the range loop stops it at once. No generator
// recurses: each is a loop over an explicit state array, so depth is never a limit.
//
// Buffer reuse: an iterator yields ONE slice, rewritten in place between steps (no allocation per
// visit). Copy an element to keep it (slices.Clone, or string(b) for parentheses) and never modify the
// yielded slice: it is, or mirrors, the generator's state. The input slice is copied when iteration
// starts (Permutations, PermutationsUnique, SubsetsUnique) or read during iteration (Combinations,
// Subsets); it is never modified, except by NextPermutationInPlace, whose name says so.
//
// Sizes: a function whose only size is a slice never fails (k < 0 or k > n is the empty enumeration). A
// function taking a bare count returns (iter.Seq, error) so an invalid count is an error, not silence.
// Values compare with cmp.Compare order (a NaN sorts first, -0 equals +0).
//
// Functions:
//
//	ALG-ENUM-001  NextPermutationInPlace  the next permutation in lexicographic order, in place
//	ALG-ENUM-002  Permutations            all n! orderings of the positions, lexicographic
//	ALG-ENUM-003  PermutationsUnique      the distinct permutations of a multiset, lexicographic
//	ALG-ENUM-004  Combinations            the k-element combinations of the positions, lexicographic
//	ALG-ENUM-005  Subsets                 all 2ⁿ subsets, binary counting order
//	ALG-ENUM-006  SubsetsUnique           the distinct sub-multisets of a multiset
//	ALG-ENUM-007  CombinationSum          the multisets of candidates (reusable) with a given sum
//	ALG-ENUM-008  GenerateParentheses     the balanced strings of n pairs of parentheses
//	ALG-ENUM-009  GrayCodeSequence        the binary reflected Gray code on a number of bits
package enumerate
