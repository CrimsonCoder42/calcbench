// Package combinatorics counts arrangements and selections exactly: factorials, binomial coefficients,
// permutations, multisets, circular arrangements, derangements, Catalan, Stirling and Bell numbers, and
// the log forms that probability needs.
//
// ID short form: COMBI (IDs are NUM-COMBI-NNN).
//
// Results: exact counts are uint64. A count above 2⁶⁴ − 1 is ErrOverflow, never a wrapped value; the
// functions whose results grow past that quickly have a *big.Int twin (the Big suffix). Factorial,
// Binomial and Permutations need no twin: their big forms are one standard-library call each,
// new(big.Int).MulRange(1, n), new(big.Int).Binomial(n, k) and new(big.Int).MulRange(n-k+1, n).
// Log forms are float64 natural logs. Every function is pure; the uint64 and float64 functions allocate
// nothing when they succeed (an error allocates its wrapped message). No function modifies its input.
// Errors are wrapped with the offending parameter; test with errors.Is.
//
// Functions:
//
//	NUM-COMBI-001  Factorial                       n! exactly, n ≤ 20
//	NUM-COMBI-002  LogFactorial                    natural log of n!
//	NUM-COMBI-003  Binomial                        C(n,k) exactly, overflow detected
//	NUM-COMBI-004  LogBinomial                     natural log of C(n,k)
//	NUM-COMBI-005  Permutations                    nPr = n!/(n−k)! exactly
//	NUM-COMBI-006  PermutationsWithRepetition      nᵏ sequences of length k over n symbols
//	NUM-COMBI-007  CombinationsWithRepetition      multisets of size k from n kinds, C(n+k−1,k)
//	NUM-COMBI-008  Multinomial                     (Σcᵢ)!/Πcᵢ!, permutations of a multiset
//	NUM-COMBI-009  CircularPermutations            arrangements of k of n around a circle, rotations equal
//	NUM-COMBI-010  CircularPermutationsUndirected  circular arrangements, rotations and reflections equal
//	NUM-COMBI-011  Derangements                    permutations of n with no fixed point, n ≤ 20
//	NUM-COMBI-012  Catalan                         n-th Catalan number, n ≤ 36
//	NUM-COMBI-013  StirlingFirstUnsigned           permutations of n with exactly k cycles
//	NUM-COMBI-014  StirlingSecond                  partitions of an n-set into k blocks
//	NUM-COMBI-015  Bell                            number of partitions of an n-set, n ≤ 25
//	NUM-COMBI-016  StirlingApproximation           √(2πn)(n/e)ⁿ, with Robbins's error bound
//	NUM-COMBI-017  MultinomialBig                  Multinomial as a *big.Int
//	NUM-COMBI-018  DerangementsBig                 Derangements as a *big.Int
//	NUM-COMBI-019  CatalanBig                      Catalan as a *big.Int
//	NUM-COMBI-020  StirlingFirstUnsignedBig        StirlingFirstUnsigned as a *big.Int
//	NUM-COMBI-021  StirlingSecondBig               StirlingSecond as a *big.Int
//	NUM-COMBI-022  BellBig                         Bell as a *big.Int
package combinatorics
