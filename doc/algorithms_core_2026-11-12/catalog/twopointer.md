# algo/twopointer (SHORT = TWOPTR) — S2, new package (builder W1-B2)

**Not in doc/taxonomy.md**: the lead's brief created it ("NEW, SHORT TWOPTR"). Taxonomy needs the row
`| ALG | algo/twopointer | TWOPTR | two-index scans: pair and triple sums, in-place dedup, merge, partition |`
(lead's file). Variable-length windows live in algo/window, string forms in algo/text.

Error: `ErrLength` (merge_sorted_in_place.go:10).

## Sources (verification state)
- **Laaksonen**, *Competitive Programmer's Handbook*, 2018, §8.1 "Two pointers method" pp. 77–79 (read
  first-hand): pointers that move one way only, O(n); "Subarray sum" (positive values) and "2SUM problem"
  (1 4 5 6 7 9 9 10, x = 12 → 5 + 7).
- **Sedgewick & Wayne**, *Algorithms*, 4th ed., 2011, Exercise 1.4.15 "Faster 3-sum" (linear two-sum on a
  sorted array, then quadratic 3-sum). Exercise number and content from search results; the book's
  exercise text was not read first-hand (the booksite lists only some exercises). §2.5 "Duplicates" read
  first-hand on algs4.cs.princeton.edu/25applications.
- **Gajentaan & Overmars** 1995, doi:10.1016/0925-7721(95)00022-2 — Crossref-confirmed: Computational
  Geometry 5(3):165–185.
- **Knuth** TAOCP Vol. 3, 2nd ed., §5.2.4 "Sorting by Merging", Algorithm M, p. 158 (section and page from a
  bibliographic search; not re-read). **CLRS** 3rd ed. §2.3.1 MERGE (Figure 2.3 example), §7.1 PARTITION
  (Lomuto) — standard; not re-read.
- **Manacher** 1975, doi:10.1145/321892.321896 — Crossref-confirmed: JACM 22(3):346–351.
- **Go standard library**, slices.Compact, https://pkg.go.dev/slices#Compact.

### ALG-TWOPTR-001 TwoSumSorted (two_sum_sorted.go:23)
`func TwoSumSorted[T integer](sorted []T, target T) (i, j int, ok bool)`
- **Source**: Laaksonen §8.1 "2SUM" pp. 78–79; S&W Ex. 1.4.15.
- **Best known**: Θ(n) on sorted input (must read the data in the worst case), O(1). **Lost**: binary search
  per element (O(n log n), S&W TwoSumFast); hash set (O(n) expected, O(n) extra; needed only when unsorted).
- **Approach**: pointers from both ends; sums compared in int128, so wrap-around never matches (uint8
  100 + 200 ≠ 44). Result: smallest i, then largest j (proved in the oracle against an O(n²) scan).
- **Worked example**: Laaksonen → (2, 4). **Errors**: none (absent → (−1, −1, false)). Unsorted input:
  unspecified result, no panic. **Input modified**: no.
- **Tags**: two-pointers, two-sum. LeetCode 167, 1.

### ALG-TWOPTR-002 ThreeSum (three_sum.go:19)
`func ThreeSum[T integer](dst [][3]T, values []T, target T) [][3]T`
- **Source**: S&W Ex. 1.4.15; Gajentaan & Overmars 1995; Laaksonen §8.1.
- **Best known**: O(n²) is the practical bound; subquadratic word-RAM algorithms exist (Grønlund & Pettie,
  FOCS 2014; not verified further here, impractical); output can be Θ(n²) triples. **Lost**: hash per pair
  (O(n²) expected, O(n) extra, harder dedup); O(n³) enumeration.
- **Approach**: sorted copy (input not modified, one O(n) allocation), fixed first element, two pointers,
  duplicate skipping on all three positions; exact int128 sums. Output: each triple ascending, triples
  lexicographic, appended to dst.
- **Worked example**: none numeric in the sources; hand-derived 30 −40 −20 −10 40 0 10 5 → 4 triples.
- **Errors**: none. **Tags**: two-pointers, 3sum. LeetCode 15, 16, 18.

### ALG-TWOPTR-003 RemoveDuplicatesInPlace (remove_duplicates_in_place.go:17)
`func RemoveDuplicatesInPlace[T comparable](values []T) int`
- **Source**: S&W §2.5 "Duplicates"; Go slices.Compact.
- **Best known**: Θ(n), O(1). **Lost**: calling slices.Compact (returns a slice, and clears the tail: an extra
  pass of writes). The tail values[n:] is left as is (documented). Equality ==: NaN kept, −0/+0 merged.
- **Note for the lead**: this overlaps slices.Compact closely; kept because the brief asks for the
  length-returning in-place form.
- **Worked example**: hand-made 1 1 2 3 3 3 7 → 4. **Input modified**: yes (name says so).
- **Tags**: dedup, in-place. LeetCode 26, 80.

### ALG-TWOPTR-004 MergeSortedInPlace (merge_sorted_in_place.go:26)
`func MergeSortedInPlace[T cmp.Ordered](buf []T, m int, other []T) error`
- **Source**: Knuth §5.2.4 Algorithm M; CLRS §2.3.1 MERGE.
- **Best known**: Θ(m + k), O(1) extra (the buffer already has room). **Lost**: front-to-back merge (needs
  O(m) scratch or overwrites unread elements); append then sort (O((m+k) log(m+k))).
- **Approach**: back to front, take from buf only when strictly greater, so equal elements of buf stay first
  (stable; visible with −0/+0 in the test). cmp.Compare order (NaN first).
- **Worked example**: CLRS Figure 2.3, 2 4 5 7 + 1 2 3 6 → 1 2 2 3 4 5 6 7. **Errors**: ErrLength (m < 0 or no
  room), buf unchanged. other must not overlap buf (documented, not detectable). **Input modified**: buf yes.
- **Tags**: merge, in-place, stable. LeetCode 88.

### ALG-TWOPTR-005 IsPalindrome (is_palindrome.go:14)
`func IsPalindrome[T comparable](values []T) bool`
- **Source**: Manacher 1975 (definition); Laaksonen §8.1 (pointer pattern). Slices only; strings → algo/text.
- **Best known**: Θ(n) worst case, O(1). **Worked example**: hand-made 3 1 4 1 3. NaN off-centre → false.
- **Tags**: palindrome. LeetCode 125, 234 (on slices).

### ALG-TWOPTR-006 MoveZeroesInPlace (move_zeroes_in_place.go:14)
`func MoveZeroesInPlace[T comparable](values []T) int`
- **Source**: CLRS 3rd ed. §7.1 PARTITION (Lomuto: one boundary index, swap kept elements to it).
- **Best known**: Θ(n), O(1), ≤ n swaps. **Lost**: copy non-zeros forward then fill zeros (fewer writes when
  zeros are few, but rewrites −0 as +0, so the multiset changes); stable partition by rotations (O(n log n)).
- **Decision (lead to log)**: named **MoveZeroesInPlace**, not MoveZeroes: the project rule says a function
  that modifies its input says so in its name (…InPlace). Returns the non-zero count. Zero = T's zero value
  (== compare): −0 counts as zero (sign bit preserved by swapping), NaN does not.
- **Worked example**: hand-made 4 0 0 −2 7 0 → 4 −2 7 0 0 0. **Tags**: partition, in-place. LeetCode 283.

## HELD (no source)
- **ContainerMostWater**: the "move the shorter side inward" two-pointer and its correctness argument appear
  only on problem and interview sites; no textbook or paper found (searched 2026-11-12). LeetCode 11.
- **TrappingRainWater**: the formula water(i) = min(max left, max right) − h(i) and its two-pointer form have
  no textbook or paper source found (searched 2026-11-12); only problem/interview sites. LeetCode 42.
  (If Devin accepts the formula as a definition, it is a max-prescan (Blelloch §1.1) each way plus a sum.)

## Consumer evidence
the consumer backend/knowledgebase/series/ has no hand-rolled pair-sum, merge or partition; rank.go's
MergeIntervals belongs to algo/interval. None to cite.
