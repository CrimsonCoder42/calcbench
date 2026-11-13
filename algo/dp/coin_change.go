package dp

import (
	"fmt"
	"math"
	"slices"
)

// CoinChangeMin returns the fewest coins, each of a denomination in coins and each denomination usable
// any number of times, whose values sum to exactly amount; ok is false when no combination does.
// amount 0 needs 0 coins. Repeated denominations are allowed and change nothing. The dynamic program
// is exact for every coin system, including those where the greedy "largest coin first" rule fails.
// coins is read, never modified.
// Errors: ErrNegative when amount < 0; ErrNotPositive when a denomination is ≤ 0; ErrTooLarge when the
// table of amount+1 entries exceeds 128 MiB.
//
// ID:      ALG-DP-008
// Time:    O(k·A)   k = len(coins), A = amount (pseudo-polynomial)
// Space:   O(A) extra (one table of A+1 ints, one allocation)
// Source:  Wright, "The change-making problem", Journal of the ACM, 1975, 22(1):125–128, doi:10.1145/321864.321874 (the fewest coins for an amount, by dynamic programming over amounts)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Problem 16-1(d) (an O(nk) dynamic program for coin changing with k denominations; CLRS assumes a 1-cent coin, which the recurrence here does not need)
// Tags:    dp, coin-change, change-making, unbounded-knapsack, minimum-coins
func CoinChangeMin(coins []int, amount int) (count int, ok bool, err error) {
	if err := checkCoins(coins, amount); err != nil {
		return 0, false, err
	}
	fewest := make([]int, amount+1) // fewest[a] = fewest coins for a, or unreachable
	const unreachable = math.MaxInt
	for a := 1; a <= amount; a++ {
		best := unreachable
		for _, c := range coins {
			if c <= a && fewest[a-c] != unreachable && fewest[a-c]+1 < best {
				best = fewest[a-c] + 1
			}
		}
		fewest[a] = best
	}
	if fewest[amount] == unreachable {
		return 0, false, nil
	}
	return fewest[amount], true, nil
}

// CoinChangeCount returns the number of ways to make amount from the denominations in coins, each
// usable any number of times, counting combinations (multisets of coins), not orders: 1+2 and 2+1 are
// one way. amount 0 has one way (no coins). coins is read, never modified.
// Errors: ErrNegative when amount < 0; ErrNotPositive when a denomination is ≤ 0; ErrDuplicate when a
// denomination repeats (it would count the same coin twice); ErrTooLarge when the table of amount+1
// entries exceeds 128 MiB; ErrOverflow when the count, or the count for a smaller amount that the
// recurrence must hold on the way (lead decision L28), exceeds 2⁶⁴−1.
//
// ID:      ALG-DP-009
// Time:    O(k·A)   k = len(coins), A = amount (pseudo-polynomial); plus O(k log k) for the duplicate check
// Space:   O(A + k) extra (one table of A+1 uint64s and a sorted copy of coins, two allocations)
// Source:  Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §7.1 (ways to change an amount with given coins: the generating function ∏ 1/(1 − z^c), expanded one denomination at a time)
// Source:  Pólya, "On picture-writing", American Mathematical Monthly, 1956, 63(10):689–697, doi:10.2307/2309555 (counting the ways to make change by generating functions)
// Tags:    dp, coin-change, counting, partitions, generating-function, combinations
func CoinChangeCount(coins []int, amount int) (uint64, error) {
	if err := checkCoins(coins, amount); err != nil {
		return 0, err
	}
	sorted := slices.Clone(coins)
	slices.Sort(sorted)
	for i := 1; i < len(sorted); i++ {
		if sorted[i] == sorted[i-1] {
			return 0, fmt.Errorf("%w: coin %d", ErrDuplicate, sorted[i])
		}
	}
	// ways[a], after denominations c1..cj, = ways to make a from c1..cj: multiplying the series by
	// 1/(1 − z^c) is the running sum ways[a] += ways[a−c] in increasing a.
	ways := make([]uint64, amount+1)
	ways[0] = 1
	for _, c := range sorted {
		for a := c; a <= amount; a++ {
			s, ok := addUint64(ways[a], ways[a-c])
			if !ok {
				return 0, fmt.Errorf("%w: ways to make %d", ErrOverflow, a)
			}
			ways[a] = s
		}
	}
	return ways[amount], nil
}

// checkCoins validates the shared coin-change inputs and the table size.
func checkCoins(coins []int, amount int) error {
	if amount < 0 {
		return fmt.Errorf("%w: amount = %d", ErrNegative, amount)
	}
	for i, c := range coins {
		if c <= 0 {
			return fmt.Errorf("%w: coins[%d] = %d", ErrNotPositive, i, c)
		}
	}
	if amount == math.MaxInt || !tableFits(amount+1, 1, 8) {
		return fmt.Errorf("%w: amount = %d", ErrTooLarge, amount)
	}
	return nil
}
