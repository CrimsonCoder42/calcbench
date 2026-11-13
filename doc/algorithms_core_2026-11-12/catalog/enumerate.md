# Catalog: algo/enumerate (SHORT = ENUM), new package

Builder W3-B12, 2026-11-12. Proposed by the combinatorics subject (doc/combinatorics_2026-11-12/PENDING_SHARED.md)
and planned in doc/taxonomy.md line 41; counts (n!, C(n,k), Catalan numbers) stay in num/combinatorics.
Every large output is an iter.Seq; breaking the range loop stops the generator at once. No recursion: every
generator is a loop over an explicit state array.

**Decision E1 (buffer reuse)**: the iterators yield ONE buffer, rewritten in place between steps. A caller
that keeps an element copies it (slices.Clone, string(b)). **Lost**: a fresh slice per step: one allocation
per visit makes enumeration allocation-bound (n! or 2ⁿ allocations); Knuth's algorithms are all "visit the
current array" algorithms, and Go's bufio.Scanner.Bytes uses the same convention. The yielded buffer must
not be modified by the caller (it is the generator's state, or mirrors it).

**Decision E2 (invalid sizes)**: functions whose only arguments are slices never fail (k < 0 or k > n is
the empty enumeration, C(n,k) = 0). Functions taking a bare count (CombinationSum, GenerateParentheses,
GrayCodeSequence) return (iter.Seq, error) so a negative or too-large count is an error, not silence.

### ALG-ENUM-001 NextPermutationInPlace
`func NextPermutationInPlace[T cmp.Ordered](values []T) bool`
- **Source**: Knuth, TAOCP Vol. 4A, Combinatorial Algorithms Part 1, Addison-Wesley, 2011, §7.2.1.2,
  Algorithm L (lexicographic permutation generation; steps L2 find j, L3 increase a_j, L4 reverse a_{j+1..n}).
  Works on multisets (Knuth states Algorithm L for a_1 ≤ … ≤ a_n with repeats).
- **Best known**: O(n) worst case per step, O(1) amortised over a full run of n! steps (Knuth §7.2.1.2,
  the analysis after Algorithm L: the comparisons and interchanges per step average a constant); Ω(1) per step.
- **Approach**: Algorithm L exactly. On the last permutation (non-increasing) it reverses to the first and
  returns false (so a do-while loop visits every permutation once). **Lost**: Heap's method (not lexicographic).
- **Worked example**: [1 2 3] → [1 3 2] → [2 1 3] → [2 3 1] → [3 1 2] → [3 2 1] → false, [1 2 3].
- **Errors**: none; empty and one-element input return false. **Modifies input**: yes (name says so).
- **Tags**: permutation, lexicographic, next-permutation, combinatorics. LeetCode 31.

### ALG-ENUM-002 Permutations
`func Permutations[T any](values []T) iter.Seq[[]T]`
- **Source**: Knuth §7.2.1.2 Algorithm L run on the index array 0..n−1 (all n! orderings of the positions in
  lexicographic order); every swap and reversal is mirrored on a copy of the values, so T needs no ordering.
- **Best known**: O(1) amortised per permutation (Algorithm L), O(n! ) total; Ω(n!) output.
- **Approach**: as above. **Lost**: Heap's algorithm (Sedgewick, "Permutation Generation Methods", ACM
  Computing Surveys 9(2):137–164, 1977, doi:10.1145/356689.356692): one swap per step but not lexicographic,
  and ordered output is what callers compare against. **Worked example**: [a b c] → abc acb bac bca cab cba.
- **Errors**: none; empty input yields one empty permutation. Input copied when iteration starts; never modified.
- **Tags**: permutation, backtracking, enumeration, iterator. LeetCode 46.

### ALG-ENUM-003 PermutationsUnique
`func PermutationsUnique[T cmp.Ordered](values []T) iter.Seq[[]T]`
- **Source**: Knuth §7.2.1.2 Algorithm L on the sorted multiset (repeats produce each distinct
  arrangement once). Calls ALG-ENUM-001.
- **Best known**: O(1) amortised per distinct permutation; O(n log n) to sort first.
- **Lost**: generate all n! and deduplicate with a set (n!/∏mᵢ! of them may be far fewer than n!).
- **Worked example**: [1 1 2] → 112 121 211. **Errors**: none. **Tags**: permutation, multiset, duplicates. LeetCode 47.

### ALG-ENUM-004 Combinations
`func Combinations[T any](values []T, k int) iter.Seq[[]T]`
- **Source**: Kreher, Stinson, *Combinatorial Algorithms: Generation, Enumeration, and Search*, CRC Press,
  1999, §2.3.1 (lexicographic ordering of k-subsets; the successor algorithm kSubsetLexSuccessor — algorithm
  number not confirmed first-hand); Knuth §7.2.1.3 (generating all combinations; Knuth's Algorithm L
  visits them in colex order, the order lost here).
- **Best known**: O(1) amortised per combination for the index successor, O(k) to refresh the yielded
  values in the worst step; Ω(C(n,k)) output.
- **Approach**: index array c₀ < … < c_{k−1}; find the rightmost i with cᵢ < n−k+i, increment it, reset the
  tail to consecutive values. **Lost**: Knuth Algorithm L/T (colex order; faster by a constant, but the
  lexicographic order is the one callers expect).
- **Worked example**: [1 2 3 4], k = 2 → 12 13 14 23 24 34. **Errors**: none; k < 0 or k > n yields nothing;
  k = 0 yields one empty combination. **Tags**: combination, k-subset, lexicographic. LeetCode 77.

### ALG-ENUM-005 Subsets
`func Subsets[T any](values []T) iter.Seq[[]T]`
- **Source**: Knuth §7.2.1.1, Algorithm M (mixed-radix generation) with every radix 2: binary counting
  over membership flags, values[0] the least significant digit.
- **Best known**: O(n) per subset to rebuild the yielded slice (output size), O(1) amortised counter step.
- **Lost**: Gray order (Knuth §7.2.1.1 Algorithm G): one flag changes per step, but an inserted or removed
  element in the middle of the yielded slice still costs O(n), so nothing is saved; a uint64 mask (limits n to 63).
- **Worked example**: [a b c] → {} a b ab c ac bc abc. **Errors**: none. **Tags**: subset, power-set, binary-counting. LeetCode 78.

### ALG-ENUM-006 SubsetsUnique
`func SubsetsUnique[T cmp.Ordered](values []T) iter.Seq[[]T]`
- **Source**: Knuth §7.2.1.1 Algorithm M: a sub-multiset of a multiset with distinct values d₁<…<d_r of
  multiplicities m₁…m_r is a mixed-radix tuple (a₁…a_r), 0 ≤ aⱼ ≤ mⱼ.
- **Best known**: O(n) per sub-multiset (output), ∏(mⱼ+1) of them; O(n log n) sort first.
- **Lost**: all 2ⁿ subsets plus a set of seen keys (exponential memory).
- **Worked example**: [1 2 2] → {} 1 2 12 22 122. **Errors**: none. **Tags**: subset, multiset, duplicates. LeetCode 90.

### ALG-ENUM-007 CombinationSum
`func CombinationSum(candidates []int, target int) (iter.Seq[[]int], error)`
- **Source**: backtracking: Knuth, TAOCP Vol. 4B, Combinatorial Algorithms Part 2, Addison-Wesley, 2022,
  §7.2.2, Algorithm B (basic backtrack); Skiena, *The Algorithm Design Manual*, 2nd ed., Springer, 2008,
  §7.1 "Backtracking". The search tree: the multisets of candidates in non-decreasing order, pruned when the
  next candidate exceeds the remainder (candidates sorted ascending).
- **Best known**: output-sensitive; O(depth) per tree node, depth ≤ target / min(candidates). No
  polynomial bound in general (the count is a partition count).
- **Approach**: iterative backtracking on an explicit stack (no recursion). **Lost**: recursion (depth
  target/min could exceed the goroutine stack for large targets); DP counting tables (count, not list).
- **Worked example**: candidates [2 3 5], target 8 → [2 2 2 2] [2 3 3] [3 5].
- **Errors**: ErrNonPositive (a candidate ≤ 0: infinitely many combinations), ErrNegative (target < 0).
  Target 0 yields one empty combination. Duplicate candidates are merged. **Tags**: backtracking,
  combination-sum, partition. LeetCode 39.

### ALG-ENUM-008 GenerateParentheses
`func GenerateParentheses(pairs int) (iter.Seq[[]byte], error)`
- **Source**: Knuth, TAOCP Vol. 4A, 2011, §7.2.1.6, Algorithm P (nested parentheses in lexicographic order).
  Order: from ()()…() to ((…)), i.e. increasing when ')' is taken to precede '(' (hand-simulated for n = 2, 3).
- **Best known**: O(1) amortised per string (Knuth §7.2.1.6); C_n = Catalan(n) strings, Ω(C_n) output.
- **Lost**: backtracking on (open, close) counts (O(n) per string, recursion depth 2n).
- **Worked example**: 3 → ()()() ()(()) (())() (()()) ((())).
- **Errors**: ErrNegative (pairs < 0), ErrTooLarge (2·pairs+1 overflows int). 0 yields one empty string.
- **Tags**: parentheses, catalan, dyck-word, enumeration. LeetCode 22.

### ALG-ENUM-009 GrayCodeSequence
`func GrayCodeSequence(bits int) (iter.Seq[uint64], error)`
- **Source**: Knuth §7.2.1.1, the binary reflected Gray code g(k) = k ⊕ ⌊k/2⌋; Gray, "Pulse code
  communication", US Patent 2,632,058, 1953.
- **Best known**: O(1) per code word; 2^bits words.
- **Worked example**: 3 → 0 1 3 2 6 7 5 4. **Errors**: ErrBitsRange (bits < 0 or > 64). bits = 0 yields 0 once.
- **Tags**: gray-code, bits, enumeration. LeetCode 89.
