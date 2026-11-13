package dp

import "fmt"

// MaxProfitOneTrade returns the largest prices[sell] − prices[buy] with buy < sell (one purchase
// followed by one sale), and its buy and sell indices (the smallest sell index, then the smallest buy
// index, among the best pairs). When no sale gains (fewer than two prices, or prices never rise) it
// returns profit 0 with buy = sell = −1: no trade.
// Prices are in any one integer unit (cents, ticks) and must not be negative, so a difference always
// fits in int64. prices is read, never modified.
// Errors: ErrNegative when a price is negative.
//
// ID:      ALG-DP-022
// Time:    O(n)   n = len(prices); one pass keeping the lowest price so far
// Space:   O(1), no allocation
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 8, §8.4 (the scanning algorithm for the maximum-sum subvector; prices[s] − prices[b] is the sum of the day-to-day changes b+1 … s, so the best trade is the maximum subvector of the changes, and the scan's running state reduces to the lowest price so far)
// Tags:    dp, stock, max-profit, max-difference, kadane, scanning, trading
func MaxProfitOneTrade(prices []int64) (profit int64, buy, sell int, err error) {
	if err := checkPrices(prices); err != nil {
		return 0, -1, -1, err
	}
	buy, sell = -1, -1
	low := 0 // index of the lowest price so far (the earliest, on ties)
	for i := 1; i < len(prices); i++ {
		if gain := prices[i] - prices[low]; gain > profit {
			profit, buy, sell = gain, low, i
		}
		if prices[i] < prices[low] {
			low = i
		}
	}
	return profit, buy, sell, nil
}

// MaxProfitUnlimitedTrades returns the largest total profit from any number of non-overlapping trades
// (buy, then sell at a later index, then buy again no earlier than that sale) over prices: the sum of
// every positive day-to-day change, holding at most one unit at a time. Prices must not be negative.
// prices is read, never modified.
// Why the sum is the optimum: a trade bought at b and sold at s earns the sum of the changes b+1 … s,
// which is at most the sum of the positive ones among them; trades do not overlap, so any set of trades
// earns at most the sum of all positive changes; buying before each rise and selling after it earns
// exactly that.
// Errors: ErrNegative when a price is negative; ErrOverflow when the total does not fit in int64.
//
// ID:      ALG-DP-023
// Time:    O(n)   n = len(prices)
// Space:   O(1), no allocation
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 8, §8.4 (a trade's profit as a subvector sum of day-to-day changes; the unlimited case is the sum of the positive changes, proved above and in catalog/dp.md)
// Tags:    dp, greedy, stock, max-profit, trading
func MaxProfitUnlimitedTrades(prices []int64) (int64, error) {
	if err := checkPrices(prices); err != nil {
		return 0, err
	}
	var total int64
	for i := 1; i < len(prices); i++ {
		if rise := prices[i] - prices[i-1]; rise > 0 {
			s, ok := addInt64(total, rise)
			if !ok {
				return 0, fmt.Errorf("%w: total profit through day %d", ErrOverflow, i)
			}
			total = s
		}
	}
	return total, nil
}

// checkPrices reports a negative price.
func checkPrices(prices []int64) error {
	for i, p := range prices {
		if p < 0 {
			return fmt.Errorf("%w: prices[%d] = %d", ErrNegative, i, p)
		}
	}
	return nil
}
