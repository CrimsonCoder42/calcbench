package prefix

// CountSubarraysSumEquals returns how many non-empty contiguous subarrays values[lo:hi] have a sum
// equal to target. Values may be negative. A subarray sums to target exactly when two prefix sums
// differ by target, so one pass counts, for each prefix sum P, how many earlier prefix sums equal
// P − target. Prefix sums are kept in 128-bit arithmetic, so no input overflows.
// values is read, never modified.
//
// ID:      ALG-PREFIX-008
// Time:    O(n) expected   n = len(values); one hash-map lookup and one update per value
// Space:   O(n) extra (a map of at most n+1 distinct prefix sums)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.1 "Static array queries", "Sum queries", p. 84 (a range sum is a difference of two prefix sums), https://cses.fi/book/book.pdf (read 2026-11-12)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 "Hash Tables" (O(1) expected search and insert)
// Tags:    prefix-sum, hash-map, subarray, count, target-sum, integer
func CountSubarraysSumEquals[T integer](values []T, target T) int {
	seen := make(map[int128]int, len(values)+1)
	var prefix int128
	seen[prefix] = 1
	t := widen(target)
	count := 0
	for _, v := range values {
		prefix = prefix.add(widen(v))
		count += seen[prefix.sub(t)]
		seen[prefix]++
	}
	return count
}
