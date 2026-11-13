package enumerate

import (
	"fmt"
	"iter"
)

// maxAllocBytes is the largest buffer this package allocates from a size argument: half the runtime's
// maxAlloc (runtime/malloc.go), above which makeslice panics. 2⁴⁷ on 64-bit, 2³¹ on 32-bit (rule L41).
const maxAllocBytes = 1 << (31 + 16*(^uint(0)>>63))

// GenerateParentheses returns an iterator over the C(pairs) = Catalan(pairs) balanced strings of pairs
// pairs of parentheses, in the order of Knuth's Algorithm P: from ()()…() to ((…)), increasing when ')'
// is taken to come before '('. For 3: ()()() ()(()) (())() (()()) ((())). pairs = 0 yields one empty
// string. The yielded []byte is reused between steps: string(b) or slices.Clone to keep it, never modify it.
// Errors: ErrNegative when pairs < 0; ErrTooLarge when the 2·pairs+1-byte buffer exceeds maxAllocBytes
// (2⁴⁷ bytes on 64-bit, 2³¹ on 32-bit), above which make would panic.
//
// ID:      ALG-ENUM-008
// Time:    O(1) amortised per string, O(n) worst step; Catalan(n) strings   n = pairs
// Space:   O(n) extra: one buffer of 2n+1 bytes, allocated once per iteration
// Source:  Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.6, Algorithm P (nested parentheses in lexicographic order)
// Tags:    parentheses, catalan, dyck-word, enumeration, iterator, balanced
func GenerateParentheses(pairs int) (iter.Seq[[]byte], error) {
	if pairs < 0 {
		return nil, fmt.Errorf("%w: pairs = %d", ErrNegative, pairs)
	}
	if pairs > (maxAllocBytes-1)/2 {
		return nil, fmt.Errorf("%w: pairs = %d", ErrTooLarge, pairs)
	}
	return func(yield func([]byte) bool) {
		n := pairs
		switch n {
		case 0:
			yield([]byte{})
			return
		case 1:
			yield([]byte("()"))
			return
		}
		// a[1..2n] is the string; a[0] = ')' is Knuth's sentinel.
		a := make([]byte, 2*n+1)
		a[0] = ')'
		for k := 1; k <= n; k++ {
			a[2*k-1] = '('
			a[2*k] = ')'
		}
		m := 2*n - 1 // a[m] is the last '('; a[m+1..2n] are all ')'
		for {
			if !yield(a[1:]) { // P2: visit
				return
			}
			a[m] = ')' // P3: easy case
			if a[m-1] == ')' {
				a[m-1] = '('
				m--
				continue
			}
			j, k := m-1, 2*n-1 // P4: find j
			for a[j] == '(' {
				a[j] = ')'
				a[k] = '('
				j--
				k -= 2
			}
			if j == 0 { // P5: increase a[j]
				return
			}
			a[j] = '('
			m = 2*n - 1
		}
	}, nil
}
