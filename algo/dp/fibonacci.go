package dp

import "fmt"

// maxFibonacciIndex is the largest n with F(n) < 2⁶⁴: F(93) = 12200160415121876738; F(94) does not fit.
const maxFibonacciIndex = 93

// Fibonacci returns the Fibonacci number F(n): F(0) = 0, F(1) = 1, F(n) = F(n−1) + F(n−2).
// F(n+1) is also the number of ways to climb n stairs taking 1 or 2 steps at a time (the sequences of
// 1s and 2s that sum to n), so the "climbing stairs" count is Fibonacci(n+1).
// Errors: ErrNegative when n < 0; ErrOverflow when n > 93 (F(94) exceeds 2⁶⁴−1).
//
// ID:      ALG-DP-010
// Time:    O(n)   n ≤ 93, so at most 92 additions: O(1) in practice
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 1: Fundamental Algorithms, 3rd ed., Addison-Wesley, 1997, §1.2.8, eq. (1) and the table of F(n) (definition)
// Source:  Erickson, Algorithms, 1st ed., self-published, 2019, ISBN 978-1-792-64483-2, §3.1 "Mātrāvṛtta" (meters of short and long syllables of total length n number F(n+1), computed bottom-up)
// Tags:    dp, fibonacci, recurrence, climbing-stairs, compositions, counting
func Fibonacci(n int) (uint64, error) {
	if n < 0 {
		return 0, fmt.Errorf("%w: n = %d", ErrNegative, n)
	}
	if n > maxFibonacciIndex {
		return 0, fmt.Errorf("%w: n = %d (F(n) fits uint64 only for n ≤ %d)", ErrOverflow, n, maxFibonacciIndex)
	}
	var prev, cur uint64 = 1, 0 // F(−1) = 1, F(0) = 0
	for range n {
		prev, cur = cur, prev+cur
	}
	return cur, nil
}
