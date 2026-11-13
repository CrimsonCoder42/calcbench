package text

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

// ErrPattern: a '*' has nothing to repeat (it is first, or follows another '*').
var ErrPattern = errors.New("text: '*' must follow a rune or '.'")

// RegexMatch reports whether pattern matches all of text, where '.' matches any one rune, x* matches zero
// or more of x (x a rune or '.'), and every other rune matches itself. There are no other operators and
// no escapes; for full regular expressions use the standard regexp package. Runes, not bytes: bytes that
// are not valid UTF-8 each read as U+FFFD. text and pattern are read, never modified.
// Method: the row-by-row table D[i][j] = "the first j pattern elements match the first i text runes",
// with D[i][j] = D[i][j−1] or (x matches rune i and D[i−1][j]) for x*, and (x matches rune i and
// D[i−1][j−1]) for x; only two rows are kept. This is the state-set simulation of the pattern's
// nondeterministic automaton, so the time is polynomial for every pattern (no backtracking blow-up).
// Errors: ErrPattern when a '*' is first or follows another '*'.
//
// ID:      ALG-TEXT-021
// Time:    O(n·m)   n = runes in text, m = len(pattern) in bytes
// Space:   O(m) extra: at most one allocation of 2(m+1) booleans (the compiler may place small tables on the stack)
// Source:  Thompson, "Programming Techniques: Regular expression search algorithm", Communications of the ACM, 1968, 11(6):419–422, doi:10.1145/363347.363387
// Source:  Aho, Lam, Sethi, Ullman, Compilers: Principles, Techniques, and Tools, 2nd ed., Addison-Wesley, 2006, §3.3 (regular expressions: concatenation and Kleene closure) and §3.7.2, Algorithm 3.22 (simulating an NFA)
// Tags:    string, regex, regular-expression, kleene-star, pattern-matching, dynamic-programming, nfa
func RegexMatch(text, pattern string) (bool, error) {
	m := len(pattern)
	for p := 0; p < m; {
		if pattern[p] == '*' {
			return false, fmt.Errorf("pattern offset %d: %w", p, ErrPattern)
		}
		_, w := utf8.DecodeRuneInString(pattern[p:])
		p += w
		if p < m && pattern[p] == '*' {
			p++
		}
	}
	rows := make([]bool, 2*(m+1))
	prev, cur := rows[:m+1], rows[m+1:]
	// Row 0: the empty text matches a prefix of starred elements only.
	prev[0] = true
	for p := 0; p < m; {
		_, q, starred := regexElement(pattern, p)
		prev[q] = prev[p] && starred
		p = q
	}
	for _, t := range text {
		cur[0] = false
		for p := 0; p < m; {
			r, q, starred := regexElement(pattern, p)
			matches := r == '.' || r == t
			if starred {
				cur[q] = cur[p] || (matches && prev[q])
			} else {
				cur[q] = matches && prev[p]
			}
			p = q
		}
		prev, cur = cur, prev
	}
	return prev[m], nil
}

// regexElement decodes the element at byte offset p of a validated pattern: its rune, the offset just
// after it (past its '*' when it has one), and whether it is starred.
func regexElement(pattern string, p int) (r rune, next int, starred bool) {
	r, w := utf8.DecodeRuneInString(pattern[p:])
	next = p + w
	if next < len(pattern) && pattern[next] == '*' {
		return r, next + 1, true
	}
	return r, next, false
}
