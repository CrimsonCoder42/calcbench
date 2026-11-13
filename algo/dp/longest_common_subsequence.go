package dp

// LCSLength returns the length of a longest common subsequence of first and second: the largest k such
// that some k elements of first, in order, equal some k elements of second, in order. Elements match by
// ==, so a floating-point NaN matches nothing. Inputs are read, never modified.
//
// ID:      ALG-DP-003
// Time:    O(m·n)   m = len(first), n = len(second)
// Space:   O(min(m, n)) extra (one row of the table, one allocation; none when either input is empty)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §15.4, eq. (15.9) (the LCS recurrence) and Exercise 15.4-4 (one row of min(m, n) entries)
// Source:  Hirschberg, "A linear space algorithm for computing maximal common subsequences", Communications of the ACM, 1975, 18(6):341–343, doi:10.1145/360825.360861 (Algorithm A: the length in linear space)
// Tags:    dp, subsequence, lcs, diff, string, sequence-alignment
func LCSLength[T comparable](first, second []T) int {
	if len(second) > len(first) {
		first, second = second, first
	}
	if len(second) == 0 {
		return 0
	}
	row := make([]int, len(second)+1)
	lcsRowForward(first, second, row)
	return row[len(second)]
}

// LongestCommonSubsequence returns one longest common subsequence of first and second, as a new slice
// of their elements (empty and non-nil when there is none). Elements match by ==, so a NaN matches
// nothing. Which of several longest ones is returned is fixed by the input (Hirschberg's split takes the
// first best split point), not by chance. Inputs are read, never modified.
//
// ID:      ALG-DP-002
// Time:    O(m·n)   m = len(first), n = len(second); about twice the cell updates of LCSLength
// Space:   O(min(m, n)) extra (two rows) plus O(log max(m, n)) recursion depth; O(min(m, n)) for the result
// Source:  Hirschberg, "A linear space algorithm for computing maximal common subsequences", Communications of the ACM, 1975, 18(6):341–343, doi:10.1145/360825.360861 (Algorithms B and C: split the longer sequence in half, find the best split of the other by a forward and a backward row)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §15.4, eq. (15.9) (the LCS recurrence)
// Tags:    dp, subsequence, lcs, diff, string, sequence-alignment, divide-and-conquer, linear-space
func LongestCommonSubsequence[T comparable](first, second []T) []T {
	// The longer sequence is the one halved (recursion depth log of it); rows span the shorter one.
	if len(second) > len(first) {
		first, second = second, first
	}
	out := make([]T, 0, len(second))
	if len(second) == 0 {
		return out
	}
	rows := make([]int, 2*(len(second)+1))
	fwd, bwd := rows[:len(second)+1], rows[len(second)+1:]
	return hirschberg(first, second, fwd, bwd, out)
}

// hirschberg appends an LCS of a and b to out. fwd and bwd are scratch rows of at least len(b)+1.
// Recursion halves a, so its depth is at most ⌈log2 len(a)⌉ + 1.
func hirschberg[T comparable](a, b []T, fwd, bwd []int, out []T) []T {
	if len(a) == 0 || len(b) == 0 {
		return out
	}
	if len(a) == 1 {
		for _, y := range b {
			if y == a[0] {
				return append(out, a[0])
			}
		}
		return out
	}
	mid := len(a) / 2
	f, g := fwd[:len(b)+1], bwd[:len(b)+1]
	lcsRowForward(a[:mid], b, f)  // f[j] = LCS(a[:mid], b[:j])
	lcsRowBackward(a[mid:], b, g) // g[j] = LCS(a[mid:], b[j:])
	split, best := 0, -1
	for j := range f {
		if s := f[j] + g[j]; s > best {
			split, best = j, s
		}
	}
	out = hirschberg(a[:mid], b[:split], fwd, bwd, out)
	return hirschberg(a[mid:], b[split:], fwd, bwd, out)
}

// lcsRowForward sets row[j] = LCS(a, b[:j]) for j = 0..len(b); row must have len(b)+1 entries.
func lcsRowForward[T comparable](a, b []T, row []int) {
	clear(row)
	for _, x := range a {
		diag := 0 // the old row[j]
		for j, y := range b {
			up := row[j+1]
			if x == y {
				row[j+1] = diag + 1
			} else if row[j] > up {
				row[j+1] = row[j]
			}
			diag = up
		}
	}
}

// lcsRowBackward sets row[j] = LCS(a, b[j:]) for j = 0..len(b); row must have len(b)+1 entries.
func lcsRowBackward[T comparable](a, b []T, row []int) {
	clear(row)
	for i := len(a) - 1; i >= 0; i-- {
		x := a[i]
		diag := 0 // the old row[j+1]
		for j := len(b) - 1; j >= 0; j-- {
			up := row[j]
			if x == b[j] {
				row[j] = diag + 1
			} else if row[j+1] > up {
				row[j] = row[j+1]
			}
			diag = up
		}
	}
}
