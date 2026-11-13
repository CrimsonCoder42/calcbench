package prefix

import (
	"fmt"
	"slices"
)

// DifferenceArray holds an array of n integers as the differences between neighbours, so adding a
// constant to a whole range costs O(1); Materialize turns it back into the array (its prefix sums).
// Differences are kept in 128-bit arithmetic, so intermediate values never overflow and Materialize
// reports exactly the elements that do not fit T. The zero value is an empty array.
// A DifferenceArray is not safe for concurrent use.
type DifferenceArray[T integer] struct {
	diff []int128 // len n+1: diff[0] = a[0], diff[i] = a[i] − a[i−1]; diff[n] absorbs updates that reach the end
}

// NewDifferenceArray returns a DifferenceArray holding a copy of initial. Pass make([]T, n) for n
// zeros. initial is read, never modified.
//
// ID:      ALG-PREFIX-005
// Time:    O(n)   n = len(initial)
// Space:   O(n) for the result (one allocation of n+1 128-bit entries)
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.4 "Additional techniques", "Range updates", p. 93 (the difference array, whose prefix sums are the array), https://cses.fi/book/book.pdf (read 2026-11-12)
// Source:  Blelloch, "Prefix Sums and Their Applications", in Reif (ed.), Synthesis of Parallel Algorithms, Morgan Kaufmann, 1993, ch. 1, §1.1 (the scan, which inverts the difference)
// Tags:    difference-array, prefix-sum, range-update, integer
func NewDifferenceArray[T integer](initial []T) *DifferenceArray[T] {
	diff := make([]int128, len(initial)+1)
	prev := int128{}
	for i, v := range initial {
		w := widen(v)
		diff[i] = w.sub(prev)
		prev = w
	}
	return &DifferenceArray[T]{diff: diff}
}

// Len returns the number of elements n. A nil or zero DifferenceArray has length 0.
//
// ID:      ALG-PREFIX-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.4 "Additional techniques", "Range updates", p. 93, https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    difference-array, length
func (d *DifferenceArray[T]) Len() int {
	if d == nil || len(d.diff) == 0 {
		return 0
	}
	return len(d.diff) - 1
}

// RangeAdd adds delta to every element of the half-open range [lo, hi). lo == hi is a no-op.
// Exact for fewer than 2⁶³ calls (the 128-bit differences cannot overflow before then).
// Errors: ErrRange unless 0 <= lo <= hi <= Len(); the array is then unchanged.
//
// ID:      ALG-PREFIX-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.4 "Additional techniques", "Range updates", p. 93 (increase the difference at the start of the range, decrease it after the end), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    difference-array, range-update, range-add, integer
func (d *DifferenceArray[T]) RangeAdd(lo, hi int, delta T) error {
	n := d.Len()
	if lo < 0 || hi < lo || hi > n {
		return fmt.Errorf("%w: [lo, hi) = [%d, %d) with Len() = %d", ErrRange, lo, hi, n)
	}
	if lo == hi {
		return nil
	}
	w := widen(delta)
	d.diff[lo] = d.diff[lo].add(w)
	d.diff[hi] = d.diff[hi].sub(w)
	return nil
}

// Materialize appends the n current elements to dst and returns the extended slice: element i is
// diff[0] + … + diff[i], summed in 128-bit arithmetic. The DifferenceArray is not changed.
// Errors: ErrOverflow when an element does not fit T (its index is in the message). On error dst is
// returned with its length unchanged.
//
// ID:      ALG-PREFIX-007
// Time:    O(n)   n = Len()
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §9.4 "Additional techniques", "Range updates", p. 93 ("the original array is the prefix sum array of the difference array"), https://cses.fi/book/book.pdf (read 2026-11-12)
// Tags:    difference-array, prefix-sum, range-update, integer, overflow-checked
func (d *DifferenceArray[T]) Materialize(dst []T) ([]T, error) {
	n := d.Len()
	out := slices.Grow(dst, n)
	var running int128
	for i := range n {
		running = running.add(d.diff[i])
		v, ok := narrow[T](running)
		if !ok {
			return dst, fmt.Errorf("%w: element %d", ErrOverflow, i)
		}
		out = append(out, v)
	}
	return out, nil
}
