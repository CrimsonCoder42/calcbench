package text

// MinimumWindowSubstring returns the shortest window text[start : start+length] that contains every byte
// of target at least as many times as target does (a multiset of bytes; order does not matter). Among
// shortest windows it returns the leftmost. found is false, with (0, 0), when no window qualifies; an
// empty target gives (0, 0, true). Bytes are compared. text and target are read, never modified.
// Two pointers move right only: the right one extends the window until it covers target, then the left
// one advances while the window still covers it, and each covering window is a candidate.
//
// ID:      ALG-TEXT-014
// Time:    O(n + m)   n = len(text), m = len(target); each pointer moves at most n times
// Space:   O(1) extra, no allocation (a 256-entry count table)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §8.1 "Two pointers method", p. 77 (both pointers move one way only), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    string, window, sliding-window, two-pointers, minimum-window, cover, multiset
func MinimumWindowSubstring(text, target string) (start, length int, found bool) {
	if len(target) == 0 {
		return 0, 0, true
	}
	var need [256]int // need[c] > 0: the window still lacks that many c
	for i := range len(target) {
		need[target[i]]++
	}
	missing := len(target)
	lo := 0
	for hi := range len(text) {
		c := text[hi]
		if need[c] > 0 {
			missing--
		}
		need[c]--
		for missing == 0 {
			if w := hi + 1 - lo; !found || w < length {
				start, length, found = lo, w, true
			}
			d := text[lo]
			need[d]++
			if need[d] > 0 {
				missing++
			}
			lo++
		}
	}
	return start, length, found
}
