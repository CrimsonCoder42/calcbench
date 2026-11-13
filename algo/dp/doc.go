// Package dp holds dynamic-programming patterns: subsequences, edit distance, knapsacks, change
// making, counting paths, segmentation, palindromes, interval DP and one-pass scans.
//
// ID short form: DP (IDs are ALG-DP-NNN).
//
// Conventions. No function modifies its input. Sequences are generic: ordered comparisons use
// cmp.Compare (a NaN below every other float); matches use == (a NaN matches nothing); for text pass
// []rune(s) to work in characters or []byte(s) to work in bytes. Values, prices and costs are int64,
// counts uint64, sizes and weights int; every integer accumulation is checked and returns ErrOverflow
// instead of wrapping. A count or path sum that the recurrence holds for a smaller sub-problem is part
// of its state, so ErrOverflow is reported when such a partial value overflows even if the final answer
// would fit (lead decision L28); MatrixChainOrder, whose costs are all non-negative, reports only a true
// overflow of the answer. Tables beyond O(input) are capped at 128 MiB (ErrTooLarge, before any
// allocation; lead decisions L13, L26). Invalid input returns a sentinel error wrapped with the
// offending parameter; nothing panics on input. Recursion appears only in LongestCommonSubsequence
// (depth O(log n)) and MatrixChainOrder's printing (depth ≤ n ≤ about 3300).
//
// Related homes: the longest palindromic substring (Manacher's linear-time algorithm) is algo/text
// LongestPalindromicSubstring, ALG-TEXT-009; maximum subarray (Kadane) is algo/prefix MaxSubarray;
// MaximalRectangle is in algo/stack. The "climbing stairs" count with steps of 1 and 2 is
// Fibonacci(n+1).
//
// Functions:
//
//	ALG-DP-001  LongestIncreasingSubsequence  indices of one longest strictly increasing subsequence (patience sorting)
//	ALG-DP-002  LongestCommonSubsequence      one longest common subsequence in linear space (Hirschberg)
//	ALG-DP-003  LCSLength                     length of a longest common subsequence, one row of memory
//	ALG-DP-004  EditDistance                  Levenshtein distance, one row of memory
//	ALG-DP-005  OneEditDistance               whether two sequences are exactly one edit apart
//	ALG-DP-006  Knapsack01                    best value and items of the 0-1 knapsack
//	ALG-DP-007  UnboundedKnapsack             best value and copies per item, items reusable
//	ALG-DP-008  CoinChangeMin                 fewest coins that make an amount
//	ALG-DP-009  CoinChangeCount               number of coin combinations that make an amount
//	ALG-DP-010  Fibonacci                     F(n) for n ≤ 93; Fibonacci(n+1) is the climbing-stairs count
//	ALG-DP-011  MaxNonAdjacentSum             largest sum of pairwise non-adjacent elements (path)
//	ALG-DP-012  MaxNonAdjacentSumCircular     the same on a ring
//	ALG-DP-013  GridPathCount                 right/down paths across a grid with blocked cells
//	ALG-DP-014  MinPathSumGrid                least-cost right/down path across a grid
//	ALG-DP-015  TextSegmentation              split text into dictionary words
//	ALG-DP-016  (unused: LongestPalindromicSubstring has its one home in algo/text, ALG-TEXT-009; lead decision L33)
//	ALG-DP-017  PalindromePartitionMinCuts    fewest cuts into palindromic pieces
//	ALG-DP-018  RodCutting                    best revenue and pieces for a rod (CLRS §15.1)
//	ALG-DP-019  MatrixChainOrder              fewest multiplications and parenthesisation (CLRS §15.2)
//	ALG-DP-020  CanJumpToEnd                  whether forward jumps reach the last index
//	ALG-DP-021  MinJumpsToEnd                 fewest forward jumps to the last index
//	ALG-DP-022  MaxProfitOneTrade             best single buy-then-sell
//	ALG-DP-023  MaxProfitUnlimitedTrades      best total over any number of non-overlapping trades
//	ALG-DP-024  MinPathSumTriangle            least-sum apex-to-base path in a number triangle
//	ALG-DP-025  SubsequenceOccurrences        times a pattern occurs in a text as a subsequence
package dp
