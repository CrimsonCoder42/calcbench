package window

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
)

// ErrNegativeWindow: the time window is negative.
var ErrNegativeWindow = errors.New("window: time window must not be negative")

// LabeledEvent is one event: an instant At in the caller's integer time unit (Unix seconds, Unix
// nanoseconds, minutes since an epoch) and a label.
type LabeledEvent[L cmp.Ordered] struct {
	At    int64
	Label L
}

// CoOccurrence is the number of events labelled A that have at least one event labelled B within the
// window (A ≠ B). Count(A, B) and Count(B, A) differ in general: each counts events of its own A label.
type CoOccurrence[L cmp.Ordered] struct {
	A, B  L
	Count int
}

// CoOccurrenceWithinWindow returns, for every ordered pair of distinct labels (A, B), the number of events
// labelled A with at least one event labelled B no more than window apart in time (|Δt| ≤ window, both
// sides), sorted by A then B; pairs with count 0 are left out. Events need not be sorted; equal instants
// count as within any window. The events are sorted by time (a copy), and two pointers sweep the window
// [t − window, t + window] around each event while a set of the labels present in it is kept with
// O(1) insert and delete; each event then visits only the labels in its window. events is read, never
// modified.
// Errors: ErrNegativeWindow when window < 0.
//
// ID:      ALG-WINDOW-016
// Time:    O(n log n + P)   n = len(events), P = Σ over events of the number of distinct labels within the window of that event (P ≤ n·k, k = number of distinct labels); O(n log n + P + R log R) when k² > 4n + 1024 and counts go to a hash map (R = result pairs ≤ P)
// Space:   O(n + k) extra (the sorted copy, label ids, window counts; plus k² counters when k² ≤ 4n + 1024, else a map of R entries); O(R) for the result
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §8.1 "Two pointers method", p. 77 (both pointers move one way only, so O(n) moves in total), https://cses.fi/book/book.pdf (read 2026-11-12)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 2 and §6.4 (sorting), ch. 11 "Hash Tables" (O(1) expected insert and search) and ch. 17 (amortized analysis of the sweep)
// Tags:    window, sliding-window, two-pointers, sweep, co-occurrence, events, pairs, time-window
// Example: events a@0, a@30, b@10, c@100, window 24 → (a, b) 2, (b, a) 1 (hand-derived; the test reproduces it)
func CoOccurrenceWithinWindow[L cmp.Ordered](events []LabeledEvent[L], window int64) ([]CoOccurrence[L], error) {
	if window < 0 {
		return nil, fmt.Errorf("%w: window = %d", ErrNegativeWindow, window)
	}
	n := len(events)
	if n < 2 {
		return []CoOccurrence[L]{}, nil
	}
	sorted := slices.Clone(events)
	slices.SortStableFunc(sorted, func(a, b LabeledEvent[L]) int { return cmp.Compare(a.At, b.At) })

	// Dense label ids 0…k−1 in label order.
	labels := make([]L, n)
	for i, e := range sorted {
		labels[i] = e.Label
	}
	slices.Sort(labels)
	labels = slices.Compact(labels)
	k := len(labels)
	if k < 2 {
		return []CoOccurrence[L]{}, nil
	}
	ids := make([]int32, n)
	for i, e := range sorted {
		j, _ := slices.BinarySearch(labels, e.Label)
		ids[i] = int32(j)
	}

	// within reports whether later − earlier ≤ window for earlier ≤ later, without int64 overflow.
	within := func(earlier, later int64) bool { return uint64(later)-uint64(earlier) <= uint64(window) }

	inWindow := make([]int, k) // events of each label in the current window
	pos := make([]int32, k)    // index of a label in present, valid while inWindow > 0
	present := make([]int32, 0, k)
	add := func(id int32) {
		if inWindow[id] == 0 {
			pos[id] = int32(len(present))
			present = append(present, id)
		}
		inWindow[id]++
	}
	remove := func(id int32) {
		inWindow[id]--
		if inWindow[id] == 0 {
			last := present[len(present)-1]
			present[pos[id]] = last
			pos[last] = pos[id]
			present = present[:len(present)-1]
		}
	}

	dense := k*k <= 4*n+1024
	var grid []int
	var sparse map[int64]int
	if dense {
		grid = make([]int, k*k)
	} else {
		sparse = make(map[int64]int)
	}

	lo, hi := 0, 0
	for i := range n {
		t := sorted[i].At
		for hi < n && within(t, sorted[hi].At) {
			add(ids[hi])
			hi++
		}
		for !within(sorted[lo].At, t) {
			remove(ids[lo])
			lo++
		}
		a := ids[i]
		for _, b := range present {
			if b == a {
				continue
			}
			if dense {
				grid[int(a)*k+int(b)]++
			} else {
				sparse[int64(a)*int64(k)+int64(b)]++
			}
		}
	}

	if dense {
		r := 0
		for _, c := range grid {
			if c > 0 {
				r++
			}
		}
		out := make([]CoOccurrence[L], 0, r)
		for a := range k {
			for b := range k {
				if c := grid[a*k+b]; c > 0 {
					out = append(out, CoOccurrence[L]{A: labels[a], B: labels[b], Count: c})
				}
			}
		}
		return out, nil
	}
	keys := make([]int64, 0, len(sparse))
	for key := range sparse {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	out := make([]CoOccurrence[L], len(keys))
	for i, key := range keys {
		out[i] = CoOccurrence[L]{A: labels[key/int64(k)], B: labels[key%int64(k)], Count: sparse[key]}
	}
	return out, nil
}
