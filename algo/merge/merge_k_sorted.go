package merge

import "cmp"

// MergeKSorted merges k ascending slices into one ascending slice.
// Equal elements keep the order of their input slices (stable): ties go to the lower list index.
// lists and its slices are read, never modified. No lists, or only empty ones, give an empty, non-nil
// slice.
//
// ID:      ALG-MERGE-001
// Time:    O(n log k + m)   n = total elements, k = number of non-empty slices, m = len(lists)
// Space:   O(k) extra (one allocation when k >= 3), O(n) for the result (one allocation)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.5 and Exercise 6.5-9 (k-way merge with a min-heap)
// Source:  Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.4.1 (multiway merging)
// Source:  Floyd, "Algorithm 245: Treesort", Communications of the ACM, 1964, 7(12):701, doi:10.1145/355588.365103 (bottom-up heap construction)
// Tags:    merge, heap, sorted, k-way, stable
func MergeKSorted[T cmp.Ordered](lists [][]T) []T {
	total, nonEmpty := 0, 0
	for _, list := range lists {
		total += len(list)
		if len(list) > 0 {
			nonEmpty++
		}
	}
	out := make([]T, 0, total)
	switch nonEmpty {
	case 0:
		return out
	case 1:
		for _, list := range lists {
			out = append(out, list...)
		}
		return out
	case 2:
		var pair [2][]T
		i := 0
		for _, list := range lists {
			if len(list) > 0 {
				pair[i] = list
				i++
			}
		}
		return mergeTwo(out, pair[0], pair[1])
	}

	h := make([]cursor[T], 0, nonEmpty)
	for i, list := range lists {
		if len(list) > 0 {
			h = append(h, cursor[T]{head: list[0], list: i, next: 1})
		}
	}
	for i := len(h)/2 - 1; i >= 0; i-- {
		siftDown(h, i)
	}
	for len(h) > 0 {
		top := &h[0]
		out = append(out, top.head)
		if src := lists[top.list]; top.next < len(src) {
			top.head = src[top.next]
			top.next++
		} else {
			last := len(h) - 1
			h[0] = h[last]
			h = h[:last]
		}
		siftDown(h, 0)
	}
	return out
}

// cursor is one input slice's position in the merge: its current head value, which slice it is, and
// the index of the element after the head.
type cursor[T cmp.Ordered] struct {
	head T
	list int
	next int
}

// before reports whether cursor x must be emitted before y: smaller head first, then lower list index.
func before[T cmp.Ordered](x, y *cursor[T]) bool {
	if c := cmp.Compare(x.head, y.head); c != 0 {
		return c < 0
	}
	return x.list < y.list
}

// siftDown restores the min-heap property of h below index i, assuming both subtrees of i are heaps.
// O(log len(h)), iterative.
func siftDown[T cmp.Ordered](h []cursor[T], i int) {
	n := len(h)
	for {
		least := i
		left, right := 2*i+1, 2*i+2
		if left < n && before(&h[left], &h[least]) {
			least = left
		}
		if right < n && before(&h[right], &h[least]) {
			least = right
		}
		if least == i {
			return
		}
		h[i], h[least] = h[least], h[i]
		i = least
	}
}

// mergeTwo appends the stable merge of first and second to out (ties take first's element first):
// one comparison per element, no heap.
func mergeTwo[T cmp.Ordered](out, first, second []T) []T {
	i, j := 0, 0
	for i < len(first) && j < len(second) {
		if cmp.Less(second[j], first[i]) {
			out = append(out, second[j])
			j++
		} else {
			out = append(out, first[i])
			i++
		}
	}
	out = append(out, first[i:]...)
	return append(out, second[j:]...)
}
