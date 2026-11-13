package list

import (
	"errors"
	"fmt"
)

// Errors of AddTwoNumbers, wrapped with the offending value; test with errors.Is.
var (
	// ErrRadix: the radix is less than 2.
	ErrRadix = errors.New("list: radix must be at least 2")
	// ErrDigit: a digit is negative or not less than the radix.
	ErrDigit = errors.New("list: digit outside [0, radix)")
)

// AddTwoNumbers returns the sum of two nonnegative integers written as digit lists in the given radix,
// least significant digit first (342 in radix 10 is 2→4→3), as a new list in the same form. nil is zero;
// a shorter operand is read as if padded with zeros. The result has max(m, n) digits, plus one when the
// last carry is 1; nil + nil is nil. Leading (most significant) zeros of the inputs are kept. Inputs are
// read, never modified. Any radix from 2 to math.MaxInt works: digit sums are formed in uint, which holds
// 2·(radix−1)+1 without overflow.
// Errors: ErrRadix when radix < 2; ErrDigit when a digit is outside [0, radix); ErrCycle when either
// list has a cycle. The result is then nil.
//
// ID:      ALG-LIST-017
// Time:    O(max(m, n))   m, n = digits in a and b
// Space:   O(1) extra; O(max(m, n)) for the result, one allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 2, 3rd ed., Addison-Wesley, 1997, §4.3.1, Algorithm A (addition of nonnegative integers in radix b, carry k ∈ {0, 1})
// Tags:    linked-list, big-integer, addition, carry, radix
func AddTwoNumbers(a, b *Node[int], radix int) (*Node[int], error) {
	if radix < 2 {
		return nil, fmt.Errorf("%w: radix = %d", ErrRadix, radix)
	}
	m, _, err := walk(a)
	if err != nil {
		return nil, fmt.Errorf("%w (a)", err)
	}
	n, _, err := walk(b)
	if err != nil {
		return nil, fmt.Errorf("%w (b)", err)
	}
	for i, p := 0, a; p != nil; i, p = i+1, p.Next {
		if p.Value < 0 || p.Value >= radix {
			return nil, fmt.Errorf("%w: a digit %d = %d, radix %d", ErrDigit, i, p.Value, radix)
		}
	}
	for i, p := 0, b; p != nil; i, p = i+1, p.Next {
		if p.Value < 0 || p.Value >= radix {
			return nil, fmt.Errorf("%w: b digit %d = %d, radix %d", ErrDigit, i, p.Value, radix)
		}
	}
	size := max(m, n)
	if size == 0 {
		return nil, nil
	}
	nodes := make([]Node[int], size+1) // one spare node for a final carry
	r := uint(radix)
	var carry uint
	for i := range size {
		var s uint = carry
		if a != nil {
			s += uint(a.Value)
			a = a.Next
		}
		if b != nil {
			s += uint(b.Value)
			b = b.Next
		}
		carry = 0
		if s >= r {
			s -= r
			carry = 1
		}
		nodes[i].Value = int(s)
		if i > 0 {
			nodes[i-1].Next = &nodes[i]
		}
	}
	if carry == 1 {
		nodes[size].Value = 1
		nodes[size-1].Next = &nodes[size]
	}
	return &nodes[0], nil
}
