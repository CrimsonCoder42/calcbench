package text

import (
	"errors"
	"fmt"
	"slices"
)

// ErrLength: a substring length is less than 1.
var ErrLength = errors.New("text: substring length must be at least 1")

// RepeatedSubstrings appends to dst, ascending, the start of the first occurrence of every distinct
// substring of the given length that occurs at least twice in text (overlapping occurrences count), and
// returns the extended slice. Windows are grouped by Karp–Rabin fingerprint (DefaultRollingBase,
// DefaultRollingModulus) and every grouping is confirmed byte by byte, so a collision never reports a
// substring that does not repeat and never hides one that does. text is read, never modified.
// Errors: ErrLength when length < 1; dst is then returned unchanged. length > len(text) gives none.
//
// ID:      ALG-TEXT-006
// Time:    O(n + L·h) expected   n = len(text), L = length, h = windows whose fingerprint equals an earlier window's (every later occurrence of a repeat counts); O(n·L) worst case, e.g. a text of one repeated byte
// Space:   O(w) extra for the fingerprint table, w = n−L+1 windows; O(r) for the result, r = substrings reported
// Source:  Karp, Rabin, "Efficient randomized pattern-matching algorithms", IBM Journal of Research and Development, 1987, 31(2):249–260, doi:10.1147/rd.312.0249
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §32.2 (rolling fingerprints) and ch. 11 (hash tables)
// Tags:    string, rolling-hash, rabin-karp, repeated, duplicate-substring, fingerprint, hash-map
func RepeatedSubstrings[S ~string | ~[]byte](dst []int, text S, length int) ([]int, error) {
	if length < 1 {
		return dst, fmt.Errorf("length %d: %w", length, ErrLength)
	}
	n := len(text)
	if length > n {
		return dst, nil
	}
	h, err := NewRollingHash(DefaultRollingBase, DefaultRollingModulus, length)
	if err != nil {
		return dst, err
	}
	type representative struct {
		start    int
		reported bool
	}
	seen := make(map[uint64][]representative)
	base := len(dst)
	for i := range length {
		h.Append(text[i])
	}
	for s := 0; ; s++ {
		reps := seen[h.Sum()]
		matched := false
		for k := range reps {
			if equalAt(text, reps[k].start, text, s, length) {
				if !reps[k].reported {
					reps[k].reported = true
					dst = append(dst, reps[k].start)
				}
				matched = true
				break
			}
		}
		if !matched {
			seen[h.Sum()] = append(reps, representative{start: s})
		}
		if s+length >= n {
			break
		}
		h.Roll(text[s], text[s+length])
	}
	slices.Sort(dst[base:])
	return dst, nil
}
