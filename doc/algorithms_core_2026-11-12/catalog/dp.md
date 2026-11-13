# Catalog: algo/dp (SHORT = DP), new package

Builder W3-B9, 2026-11-12. Package in doc/taxonomy.md (`algo/dp`, DP). Conventions (doc.go): inputs never
modified; ordered comparisons by cmp.Compare, matches by == (NaN matches nothing); values int64, counts
uint64, sizes/weights int; every integer accumulation checked (ErrOverflow); partial-state overflow policy
as L28; tables beyond O(input) capped at 128 MiB (ErrTooLarge before allocating; L13, L26). Sentinels live
in checked_arith.go. Citations were checked 2026-11-12 against Crossref and author PDFs (a research agent);
locators marked (print) could not be opened online and should be checked in print by the reviewer.

## Decisions (each with the alternative that lost)

| # | Decision | Alternative | Why it lost |
|---|---|---|---|
| D1 | Names are practitioner terms: MaxNonAdjacentSum(Circular), GridPathCount, TextSegmentation, CanJumpToEnd, MinJumpsToEnd, MaxProfitOneTrade, MaxProfitUnlimitedTrades, SubsequenceOccurrences | the brief's problem-site titles (HouseRobber, UniquePaths, WordBreak, JumpGame, BestTimeToBuySell, DistinctSubsequences) | BRIEF_COMMON: never copy identifiers; "WordBreak" also clashes with Unicode word-boundary breaking |
| D2 | ClimbingStairs not a separate function: it is Fibonacci(n+1), said in the doc comment, doc.go and tags | a second exported function | one function, one home |
| D3 | Fibonacci overflows at n ≥ 94 (uint64), not 93 as the brief says | error at 93 | F(93) = 12200160415121876738 < 2⁶⁴; 93 is the int64 limit's n+1 |
| D4 | CoinChangeCount returns uint64 + ErrOverflow | a big.Int variant | one allocation per cell; no consumer; a Big variant can be added later with its own ID |
| D5 | CoinChangeCount rejects repeated denominations (ErrDuplicate); CoinChangeMin allows them | dedupe silently | a repeat would double-count combinations; for the minimum it is harmless |
| D6 | One 128 MiB table cap for every DP table (1-D and 2-D) | per-function caps | one rule, same as L13; also stops make() panicking (L26) |
| D7 | Knapsack01 returns the items (n·W choice bits) | value only in O(W) | callers want the set; bits are 1/64 of an int table; capacity is first capped at Σ weights |
| D8 | ALG-DP-016 unused: LongestPalindromicSubstring removed (lead L33; home algo/text ALG-TEXT-009); dp/palindrome.go and its test moved to _delete/algo/dp/ | keep an O(1)-space variant | one home |

## Entries

### ALG-DP-001 LongestIncreasingSubsequence
`func LongestIncreasingSubsequence[T cmp.Ordered](values []T) []int`
- **Source**: Fredman, Discrete Mathematics 11(1):29–35, 1975, doi:10.1016/0012-365X(75)90103-X (confirmed);
  Aldous, Diaconis, Bull. AMS 36(4):413–432, 1999, doi:10.1090/S0273-0979-99-00796-X, §1.1 Lemma 1 (piles =
  LIS length; preprint numbering, published numbering (print)).
- **Best known**: O(n log n); Fredman: n log n − n log log n + O(n) comparisons needed. Reading input Ω(n).
- **Approach**: patience sorting with search.LowerBound (ALG-SEARCH-001) on pile tops, back-pointers for one
  subsequence. **Lost**: O(n²) recurrence (slower).
- **Worked example**: deck 7 2 8 1 3 4 10 6 9 5 → 5 piles; indices [3 4 5 7 8] (1 3 4 6 9).
- **Errors**: none. **Tags**: dp, lis, patience-sorting. LeetCode 300, 673.

### ALG-DP-002 LongestCommonSubsequence · ALG-DP-003 LCSLength
`func LongestCommonSubsequence[T comparable](first, second []T) []T` · `func LCSLength[T comparable](first, second []T) int`
- **Source**: Hirschberg, CACM 18(6):341–343, 1975, doi:10.1145/360825.360861 (confirmed); CLRS 3rd ed. §15.4
  eq. (15.9), Exercise 15.4-4 (print).
- **Best known**: O(mn) time (no strongly subquadratic algorithm under SETH); O(min(m, n)) space for the
  length, O(m + n) for the sequence (Hirschberg). **Approach**: one row (length); Hirschberg divide and
  conquer, halving the longer input, depth ⌈log2⌉. **Lost**: full table + traceback (O(mn) memory);
  Hunt–Szymanski (better only when matches are sparse; worse worst case).
- **Worked example**: CLRS Fig. 15.8: ABCBDAB, BDCABA → length 4 (BCBA is one).
- **Errors**: none. **Tags**: dp, lcs, diff. LeetCode 1143, 1092.

### ALG-DP-004 EditDistance · ALG-DP-005 OneEditDistance
`func EditDistance[T comparable](from, to []T) int` · `func OneEditDistance[T comparable](from, to []T) bool`
- **Source**: Wagner, Fischer, JACM 21(1):168–173, 1974, doi:10.1145/321796.321811 (confirmed); Gusfield
  1997 §11.3 (print); Ukkonen, Information and Control 64(1–3):100–118, 1985, doi:10.1016/S0019-9958(85)80046-2
  (testing distance ≤ t in O(t·min(m, n)): confirmed from the abstract).
- **Best known**: O(mn) / O(min) space; for "distance = 1" O(min(m, n)) (Ukkonen band, t = 1), O(1) space.
- **Approach**: one row, shorter input spans it. OneEditDistance: first mismatch, then compare the rest
  (slices.Equal). **Lost**: full table; for OneEditDistance, running EditDistance (O(mn)).
- **Worked examples**: Jurafsky, Martin, *Speech and Language Processing*, 3rd ed. draft (Dec 2021), §2.5
  "Minimum Edit Distance", p. 22, Fig. 2.14 (read 2026-11-12, web.stanford.edu/~jurafsky/slp3): intention →
  execution = 5 with Levenshtein unit costs, the model EditDistance uses (J&M eq. (2.8), substitution cost 2,
  gives 8, Fig. 2.18; not this function). Our own: harbor → arbour = 2 (hand, checked offline).
- **Errors**: none. **Tags**: edit-distance, levenshtein. LeetCode 72, 161.

### ALG-DP-006 Knapsack01 · ALG-DP-007 UnboundedKnapsack
`func Knapsack01(weights []int, values []int64, capacity int) (int64, []int, error)` ·
`func UnboundedKnapsack(weights []int, values []int64, capacity int) (int64, []int, error)`
- **Source**: Kellerer, Pferschy, Pisinger, Knapsack Problems, Springer 2004, doi:10.1007/978-3-540-24777-7,
  §2.3 Dynamic Programming and ch. 8 (§8.3 DP for the unbounded problem) (TOC confirmed); CLRS Exercise 16.2-2.
- **Best known**: NP-hard; O(nW) pseudo-polynomial. **Approach**: Bellman recursion over capacities, one row;
  0-1 records choices in n·W bits; unbounded keeps the last item per capacity. **Lost**: value-indexed DP
  (O(n·V), better only when values are small); branch and bound (exponential worst case).
- **Worked examples**: CLRS Fig. 16.2: (10,60),(20,100),(30,120), W 50 → 220, items [1 2]. Unbounded with the
  CLRS rod prices as items, W 4 → 10 (two of length 2).
- **Errors**: ErrLengthMismatch, ErrNegative, ErrOverflow, ErrTooLarge; ErrUnbounded (zero weight, positive
  value). **Tags**: knapsack. LeetCode 416, 474.

### ALG-DP-008 CoinChangeMin · ALG-DP-009 CoinChangeCount
- **Source**: Wright, "The change-making problem", JACM 22(1):125–128, 1975, doi:10.1145/321864.321874
  (confirmed); CLRS Problem 16-1(d) (assumes a penny; ours does not need one). Count: Graham, Knuth, Patashnik,
  Concrete Mathematics 2nd ed. §7.1 "Domino Theory and Change" (∏ 1/(1 − z^c); 292 ways for a dollar:
  confirmed via secondary sources and recomputed offline); Pólya 1956 doi:10.2307/2309555.
- **Best known**: O(kA) for arbitrary denominations (change-making is weakly NP-hard). **Lost**: greedy
  (wrong for {1,6,10}); BFS over amounts (same bound, more memory).
- **Worked examples**: {1,6,10}, 12 → 2 (greedy pays 3); {1,5,10,25,50}, 100 → 292 ways.
- **Errors**: ErrNegative, ErrNotPositive, ErrTooLarge; count also ErrDuplicate, ErrOverflow (partial, L28).
- **Tags**: coin-change. LeetCode 322, 518.

### ALG-DP-010 Fibonacci
- **Source**: Knuth TAOCP Vol. 1 3rd ed. §1.2.8 (definition; addition identity is eq. (6), confirmed);
  Erickson, Algorithms (2019) §3.1 "Mātrāvṛtta" (F(n+1) meters of total length n: the climbing-stairs count;
  confirmed). **Best known**: O(log n) by doubling, but n ≤ 93, so the loop is O(1) in practice (lost: doubling,
  no gain at n ≤ 93; closed form, float error). **Errors**: ErrNegative, ErrOverflow (n ≥ 94).
- **Worked example**: F(0…12) table; F(93) = 12200160415121876738. **Tags**: fibonacci, climbing-stairs. LeetCode 509, 70.

### ALG-DP-011 MaxNonAdjacentSum · ALG-DP-012 MaxNonAdjacentSumCircular
- **Source**: Roughgarden, Algorithms Illuminated Part 3 (2019) ch. 16 §16.1–16.2 (WIS in path graphs) (print);
  Kleinberg–Tardos §6.1 (same recurrence). **Circular, proof**: a cycle independent set cannot contain both
  values[0] and values[n−1], so it is independent in values[1:] or values[:n−1]; conversely each path set is
  a cycle set. So the optimum is the max of two path optima. (Composition of the cited recurrence, as L15.)
- **Best known**: O(n), O(1). **Worked example**: path 1 4 5 4 → 8. **Errors**: ErrOverflow (exact: an
  overflowing candidate is a feasible sum > MaxInt64). **Tags**: independent-set. LeetCode 198, 213.

### ALG-DP-013 GridPathCount · ALG-DP-014 MinPathSumGrid · ALG-DP-024 MinPathSumTriangle
- **Source**: CLRS Exercise 22.4-2 (number of paths in a DAG, linear time) (print); CLRS §24.2 (DAG shortest
  paths); Concrete Mathematics §5.1 (addition formula: the open grid gives C(r+c−2, r−1)). Krattenthaler,
  "Lattice Path Enumeration" (Handbook of Enumerative Combinatorics, 2015, ch. 10, arXiv:1503.05930) §10.7,
  Lemma 10.7.2 counts paths avoiding points by a determinant — the alternative that lost (O(k³) for k
  obstacles, needs big binomials; the DP is exact and O(rc)). The grid with right/down moves is a DAG whose
  row-major order is topological; the triangle likewise bottom-up.
- **Best known**: O(cells) time, O(row) space; reading the grid Ω(cells).
- **Worked examples**: 3×7 open → 28; [[2 9 1][4 1 8][7 3 2]] → 12; triangle 3/6 2/1 9 4/8 5 7 2 → 11.
- **Errors**: ErrNotRectangular, ErrOverflow (partial, L28); ErrEmpty (sums); ErrNotTriangle.
- **Tags**: grid, lattice-path, dag. LeetCode 62, 63, 64, 120.

### ALG-DP-015 TextSegmentation
- **Source**: Erickson, Algorithms (2019) §2.5 "Text Segmentation (Interpunctio Verborum)" and §3.3
  "Interpunctio Verborum Redux" (confirmed); Kleinberg–Tardos ch. 6 Exercise 5 (print).
- **Best known**: O(D + n·L) with a dictionary trie. **Approach**: trie in one hash map (edge → child), forward
  reachability with back-pointers. **Lost**: substring hashing per length (O(n·L²)); algo/trie's Trie (no
  exported single-step walk). **Worked example**: Erickson's BOTHEARTHANDSATURNSPIN splits.
- **Errors**: none; (nil, false) when no split. **Tags**: word-break. LeetCode 139, 140.

### ALG-DP-017 PalindromePartitionMinCuts
- **Source**: Fici, Gagie, Kärkkäinen, Kempa, J. Discrete Algorithms 28:41–48, 2014,
  doi:10.1016/j.jda.2014.08.001, §2 "A simple quadratic algorithm" (confirmed).
- **Best known**: O(n log n) time, O(n) space (same paper; also I et al. 2014). **Built**: O(n²) worst, O(n)
  space, expanding centres in increasing order so cuts[i] is final when read. **Why not the best**: the
  O(n log n) method needs a palindromic tree with series links; LEFT FOR DEVIN (go-complexity rule 3: needs a
  benchmark-backed reason or the faster build). **Worked examples**: racecarxyz → 3; noonabbad → 2 (hand,
  checked offline). **Tags**: palindrome, factorization. LeetCode 132.

### ALG-DP-018 RodCutting
- **Source**: CLRS §15.1 eq. (15.2), BOTTOM-UP-CUT-ROD, EXTENDED-BOTTOM-UP-CUT-ROD, Fig. 15.1 (print).
- **Best known**: O(n·min(n, m)). **Worked example**: r1…r10 = 1 5 8 10 13 17 18 22 25 30; 4 → 2+2; 7 → 1+6.
- Lengths beyond the price table are waste (earn 0). **Errors**: ErrNegative, ErrTooLarge, ErrOverflow.

### ALG-DP-019 MatrixChainOrder
- **Source**: CLRS §15.2 eq. (15.7), MATRIX-CHAIN-ORDER, PRINT-OPTIMAL-PARENS, Fig. 15.5 (print).
- **Best known**: O(n log n), Hu–Shing, SIAM J. Comput. 11(2):362–373, 1982, doi:10.1137/0211028 (confirmed) and
  Part II 13(2):228–251, 1984. **Built**: O(n³)/O(n²) CLRS. **Why not the best**: Hu–Shing is a long,
  intricate polygon-partition algorithm; LEFT FOR DEVIN. Overflowing splits are +∞ (costs ≥ 0), so only a true
  overflow of the answer errors. **Worked example**: 30 35 15 5 10 20 25 → 15125, ((A1(A2A3))((A4A5)A6)).
- **Errors**: ErrEmpty, ErrNotPositive, ErrTooLarge (n ≳ 3340), ErrOverflow. LeetCode none (pattern: interval DP).

### ALG-DP-020 CanJumpToEnd · ALG-DP-021 MinJumpsToEnd
- **Source**: CLRS §22.2 (BFS gives fewest-edge distances). **Proof (for the one-pass form)**: edges go
  i → i+1 … i+a_i. Claim: each BFS layer L_d is a contiguous run (e_{d−1}, e_d]. L_0 = {0}. If every index
  ≤ e_d is reached by layer d, the indices reachable in one more move are exactly (e_d, max_{i ≤ e_d}(i+a_i)],
  a run, since every j ≤ the max is covered by the edge from the index achieving it. So layer boundaries are
  the running maxima, and reachability is "farthest ≥ n−1". (Composition of a cited method, as L15; flagged.)
- **Best known**: O(n), O(1). **Lost**: explicit BFS (O(n + Σa_i)). **Worked example**: [2 0 1 0] → 2 moves.
- **Errors**: ErrNegative. **Tags**: bfs, greedy. LeetCode 55, 45.

### ALG-DP-022 MaxProfitOneTrade · ALG-DP-023 MaxProfitUnlimitedTrades
- **Source**: Bentley, Programming Pearls 2nd ed., Column 8 §8.4 "A Scanning Algorithm" (print): p[s] − p[b]
  is the sum of the changes b+1 … s, so one trade is the maximum subvector of the changes; the scan state is
  the lowest price so far. **Unlimited, proof**: a trade's profit is a sum of changes ≤ the sum of its positive
  changes; trades are disjoint, so any set earns ≤ Σ positive changes; trading every rise attains it.
  (Flagged as a short proof on a cited identity.)
- **Best known**: O(n), O(1). **Worked example**: [9 4 6 2 8 5] → one trade 6 (buy 3, sell 4); unlimited 8.
- **Errors**: ErrNegative (prices ≥ 0, so differences fit), ErrOverflow (unlimited). LeetCode 121, 122.

### ALG-DP-025 SubsequenceOccurrences
- **Source**: Lothaire, Combinatorics on Words, Addison-Wesley 1983 (CUP 1997), ch. 6 "Subwords" (Sakarovitch,
  Simon) §6.3: (ua choose vb) = (u choose vb) + δ_{a,b}(u choose v), (u choose ε) = 1 (section confirmed via
  Rigo et al., orbi.uliege.be; the book page (print)). Elzinga–Rahmann–Wang 2008 was checked and does NOT
  cover this count; not cited. **Best known**: O(nm) time, O(m) space. **Worked example**: "see" in
  "subsequence" → 6. **Errors**: ErrOverflow (partial, L28). LeetCode 115.

## HELD (no acceptable source)
- DecodeWays (digit-string decodings 1–26): no textbook or paper found.
- MaxProfit with at most k trades: Csűrös, IEEE/ACM TCBB 1(4):139–150, 2004, doi:10.1109/TCBB.2004.43 gives
  O(n log n) for k maximum-scoring segments, but whether it states the O(nk) state DP is unverified, and its
  O(n log n) method was not read in full. HELD until the paper is read.
- MaxProfit with cooldown (state-machine DP): no source found.
- MaximalSquare (largest all-true square): no citable source for the min-of-three recurrence found
  (Rosenfeld–Pfaltz 1966 distance transforms are related but not this recurrence without a proof of ours).

## Consumer evidence
None given in the brief.
