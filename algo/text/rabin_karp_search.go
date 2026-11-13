package text

import "fmt"

// RabinKarpSearch appends to dst the start of every occurrence of pattern in text, ascending and
// overlapping occurrences included, and returns the extended slice. Windows are compared by their
// Karp–Rabin fingerprint (RollingHash with base and modulus), and every fingerprint hit is confirmed
// byte by byte, so a hash collision never produces a false match. DefaultRollingBase and
// DefaultRollingModulus are good defaults. text and pattern are read, never modified.
// Errors: ErrEmptyPattern when pattern is empty; ErrHashParams as NewRollingHash. dst is then returned
// unchanged.
//
// ID:      ALG-TEXT-005
// Time:    O(n + m·(k + f))   n = len(text), m = len(pattern), k = occurrences, f = false fingerprint hits (each verified, then rejected); O(n·m) worst case
// Space:   O(1) extra, no allocation beyond the result; O(k) for the result
// Source:  Karp, Rabin, "Efficient randomized pattern-matching algorithms", IBM Journal of Research and Development, 1987, 31(2):249–260, doi:10.1147/rd.312.0249
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §32.2, RABIN-KARP-MATCHER and Figure 32.5
// Tags:    string, search, pattern-matching, rabin-karp, rolling-hash, fingerprint, find-all
func RabinKarpSearch[S ~string | ~[]byte](dst []int, text, pattern S, base, modulus uint64) ([]int, error) {
	m := len(pattern)
	if m == 0 {
		return dst, fmt.Errorf("pattern: %w", ErrEmptyPattern)
	}
	ph, err := NewRollingHash(base, modulus, m)
	if err != nil {
		return dst, err
	}
	n := len(text)
	if m > n {
		return dst, nil
	}
	th := ph
	for i := range m {
		ph.Append(pattern[i])
		th.Append(text[i])
	}
	want := ph.Sum()
	for s := 0; ; s++ {
		if th.Sum() == want && equalAt(text, s, pattern, 0, m) {
			dst = append(dst, s)
		}
		if s+m >= n {
			break
		}
		th.Roll(text[s], text[s+m])
	}
	return dst, nil
}

// equalAt reports whether a[i:i+n] and b[j:j+n] hold the same bytes; the caller keeps the ranges in bounds.
func equalAt[S ~string | ~[]byte](a S, i int, b S, j, n int) bool {
	for k := range n {
		if a[i+k] != b[j+k] {
			return false
		}
	}
	return true
}
