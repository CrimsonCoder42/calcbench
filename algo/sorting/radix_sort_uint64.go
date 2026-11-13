package sorting

// RadixSortUint64 returns values sorted ascending by least-significant-digit radix sort with 8-bit
// digits: eight stable counting passes, lowest byte first, never comparing two keys. A byte position on
// which every key agrees is skipped. The result is written into dst (reused when cap(dst) >= len(values),
// else one new slice); dst may be values itself. values is never modified unless it is dst.
//
// ID:      ALG-SORT-004
// Time:    O(n)   n = len(values); exactly: one counting pass plus at most 8 scatter passes, O(8·(n + 256))
// Space:   O(n) extra (one allocation, the scatter buffer; none when n < 2 or all keys are equal) plus 16 KiB of counts; O(n) for the result
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §8.3 (RADIX-SORT, Lemma 8.4: b-bit keys in r-bit digits, Θ((b/r)(n + 2^r)))
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2.5 (radix sorting)
// Tags:    sorting, radix-sort, lsd, integer-sort, linear-time, non-comparison
func RadixSortUint64(dst, values []uint64) []uint64 {
	n := len(values)
	out := resize(dst, n)
	copy(out, values)
	if n < 2 {
		return out
	}
	// counts[d][b] = how many keys have byte b at position d (0 = least significant).
	var counts [8][256]int
	for _, v := range out {
		for d := range 8 {
			counts[d][byte(v>>(8*d))]++
		}
	}
	var scratch []uint64
	src := out
	for d := range 8 {
		c := &counts[d]
		if c[byte(src[0]>>(8*d))] == n {
			continue // every key has the same byte here; this pass would not move anything.
		}
		if scratch == nil {
			scratch = make([]uint64, n)
		}
		to := scratch
		if &src[0] == &scratch[0] {
			to = out
		}
		// Turn counts into starting offsets, then scatter stably.
		offset := 0
		for b := range 256 {
			c[b], offset = offset, offset+c[b]
		}
		for _, v := range src {
			b := byte(v >> (8 * d))
			to[c[b]] = v
			c[b]++
		}
		src = to
	}
	if &src[0] != &out[0] {
		copy(out, src)
	}
	return out
}
