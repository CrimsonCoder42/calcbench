package dp

import "slices"

// EditDistance returns the Levenshtein distance between from and to: the fewest single-element
// insertions, deletions and substitutions, each costing 1, that turn from into to. Elements match by
// ==, so a NaN matches nothing. For text, pass []rune(s) to count characters or []byte(s) to count
// bytes. The distance is symmetric. Inputs are read, never modified.
//
// ID:      ALG-DP-004
// Time:    O(m·n)   m = len(from), n = len(to)
// Space:   O(min(m, n)) extra (one row, one allocation; none when either input is empty)
// Source:  Wagner, Fischer, "The string-to-string correction problem", Journal of the ACM, 1974, 21(1):168–173, doi:10.1145/321796.321811 (the cost recurrence over prefixes and its table algorithm, here with unit costs)
// Source:  Gusfield, Algorithms on Strings, Trees, and Sequences, Cambridge University Press, 1997, §11.3 (the edit-distance recurrence and its table)
// Tags:    dp, edit-distance, levenshtein, string, diff, fuzzy-match, spell-check
func EditDistance[T comparable](from, to []T) int {
	if len(to) > len(from) {
		from, to = to, from // unit costs: the distance is symmetric, so the row spans the shorter input
	}
	if len(to) == 0 {
		return len(from)
	}
	row := make([]int, len(to)+1) // row[j] = distance(from[:i], to[:j])
	for j := range row {
		row[j] = j
	}
	for i, x := range from {
		diag := row[0] // distance(from[:i], to[:0])
		row[0] = i + 1
		for j, y := range to {
			up := row[j+1]
			best := diag // substitution (free when x == y)
			if x != y {
				best = min(diag, up, row[j]) + 1 // substitute, delete x, insert y
			}
			row[j+1] = best
			diag = up
		}
	}
	return row[len(to)]
}

// OneEditDistance reports whether from and to are exactly one edit apart: Levenshtein distance 1
// (one insertion, deletion or substitution). Equal inputs are 0 apart and give false. Elements match
// by ==. Inputs are read, never modified.
//
// ID:      ALG-DP-005
// Time:    O(min(m, n))   m = len(from), n = len(to); O(1) when the lengths differ by more than 1
// Space:   O(1), no allocation
// Source:  Ukkonen, "Algorithms for approximate string matching", Information and Control, 1985, 64(1–3):100–118, doi:10.1016/S0019-9958(85)80046-2 (deciding distance ≤ k needs only the diagonals |i − j| ≤ k of the table: O(k·min(m, n)); here k = 1)
// Source:  Wagner, Fischer, "The string-to-string correction problem", Journal of the ACM, 1974, 21(1):168–173, doi:10.1145/321796.321811 (the edit operations)
// Tags:    dp, edit-distance, levenshtein, string, one-edit, typo, fuzzy-match
func OneEditDistance[T comparable](from, to []T) bool {
	if len(to) > len(from) {
		from, to = to, from
	}
	m, n := len(from), len(to)
	if m-n > 1 {
		return false
	}
	i := 0
	for i < n && from[i] == to[i] {
		i++
	}
	if m == n {
		// One substitution at the first mismatch, and the rest equal; no mismatch means distance 0.
		return i < n && slices.Equal(from[i+1:], to[i+1:])
	}
	// m == n+1: delete from[i] (i may be n, the last element of from).
	return slices.Equal(from[i+1:], to[i:])
}
