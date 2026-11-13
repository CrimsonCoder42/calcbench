package text

import "unicode/utf8"

// WildcardMatch reports whether pattern matches all of text, where '?' matches exactly one rune, '*'
// matches any run of runes (including none), and every other rune matches itself. There are no bracket
// expressions and no escapes, and '*' crosses '/' (unlike path.Match). Runes, not bytes: bytes that are
// not valid UTF-8 each read as U+FFFD. text and pattern are read, never modified.
// Method: the row-by-row table D[i][j] = "the first j pattern runes match the first i text runes", with
// D[i][j] = D[i][j−1] or D[i−1][j] for '*', and D[i−1][j−1] for '?' or an equal rune; only two rows are
// kept. This is the state-set simulation of the pattern's nondeterministic automaton.
//
// ID:      ALG-TEXT-020
// Time:    O(n·m)   n = runes in text, m = len(pattern) in bytes
// Space:   O(m) extra: at most one allocation of 2(m+1) booleans (the compiler may place small tables on the stack)
// Source:  IEEE Std 1003.1-2017 (POSIX.1), Shell and Utilities, §2.13.1 "Patterns Matching a Single Character" and §2.13.2 "Patterns Matching Multiple Characters" (the meaning of ? and *)
// Source:  Thompson, "Programming Techniques: Regular expression search algorithm", Communications of the ACM, 1968, 11(6):419–422, doi:10.1145/363347.363387 (simulating the automaton one text symbol at a time)
// Source:  Aho, Lam, Sethi, Ullman, Compilers: Principles, Techniques, and Tools, 2nd ed., Addison-Wesley, 2006, §3.7.2, Algorithm 3.22 (simulating an NFA)
// Tags:    string, wildcard, glob, pattern-matching, dynamic-programming, nfa
func WildcardMatch(text, pattern string) bool {
	m := len(pattern)
	rows := make([]bool, 2*(m+1))
	prev, cur := rows[:m+1], rows[m+1:]
	// Row 0: only a prefix of '*'s matches the empty text.
	prev[0] = true
	for p := 0; p < m; {
		r, w := utf8.DecodeRuneInString(pattern[p:])
		prev[p+w] = prev[p] && r == '*'
		p += w
	}
	for _, t := range text {
		cur[0] = false
		for p := 0; p < m; {
			r, w := utf8.DecodeRuneInString(pattern[p:])
			q := p + w
			switch {
			case r == '*':
				cur[q] = cur[p] || prev[q]
			case r == '?' || r == t:
				cur[q] = prev[p]
			default:
				cur[q] = false
			}
			p = q
		}
		prev, cur = cur, prev
	}
	return prev[m]
}
