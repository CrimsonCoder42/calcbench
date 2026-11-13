package selection

import "cmp"

// insertionMax is the range length at or below which insertion sort finishes the selection.
const insertionMax = 16

// workBudget is how many times len(a) the median-of-3 rounds may partition in total before every later
// pivot is the median of medians. Median-of-3 quickselect averages about 2.75n for the median
// (Kirschenhofer, Prodinger, Martínez, "Analysis of Hoare's FIND algorithm with median-of-three
// partition", Random Structures & Algorithms, 1997, 10(1-2):143–156,
// doi:10.1002/(SICI)1098-2418(199701/03)10:1/2<143::AID-RSA7>3.0.CO;2-V; measured 2.75n here). The probe in
// doc/reference_algorithms_2026-11-12/REPORT.md saw no switch at 6n in 600 random inputs. A "halve the
// range every two rounds" rule switched on 59-80% of them.
const workBudget = 6

// selectInPlace reorders a so that a[k] is the element a sorted copy would hold at index k, every
// element of a[:k] is <= a[k] and every element of a[k+1:] is >= a[k] (cmp.Compare order).
// The caller guarantees 0 <= k < len(a).
//
// Introselect (Musser 1997): quickselect with a median-of-3 pivot until the lengths it has partitioned
// exceed workBudget·n, so at most (workBudget+1)·n before the switch. After it, every pivot is the
// median of medians (Blum et al. 1973), which leaves at most 7/10 of the range, so the rest is O(n) too.
// Worst case O(n).
func selectInPlace[T cmp.Ordered](a []T, k int) {
	selectRange(a, k, false)
}

// selectRange is selectInPlace with the median-of-medians pivot forced from the start when forceMoM is
// true. Tests use it to exercise the fallback path on ordinary inputs.
func selectRange[T cmp.Ordered](a []T, k int, forceMoM bool) {
	lo, hi := 0, len(a) // a[k] lies in a[lo:hi]; everything left of lo is <= it, right of hi is >= it.
	useMoM := forceMoM
	// uint64 so the budget cannot overflow on a 32-bit int.
	var work uint64
	budget := workBudget * uint64(len(a))
	for hi-lo > insertionMax {
		var pivot T
		if useMoM {
			pivot = medianOfMedians(a[lo:hi])
		} else {
			pivot = medianOfThree(a[lo], a[lo+(hi-lo)/2], a[hi-1])
			if work += uint64(hi - lo); work > budget {
				useMoM = true // this round still runs; every later pivot is the median of medians.
			}
		}
		less, notGreater := partitionThreeWay(a[lo:hi], pivot)
		switch {
		case k < lo+less:
			hi = lo + less
		case k >= lo+notGreater:
			lo += notGreater
		default:
			return // a[k] is in the block equal to the pivot.
		}
	}
	insertionSort(a[lo:hi])
}

// partitionThreeWay reorders s into three blocks, s[:less] < pivot, s[less:notGreater] == pivot and
// s[notGreater:] > pivot, and returns the two boundaries (Dijkstra's three-way partition).
// Runs of equal keys therefore cost one pass, not a quadratic number of rounds.
func partitionThreeWay[T cmp.Ordered](s []T, pivot T) (less, notGreater int) {
	lt, i, gt := 0, 0, len(s)
	for i < gt {
		switch c := cmp.Compare(s[i], pivot); {
		case c < 0:
			s[lt], s[i] = s[i], s[lt]
			lt++
			i++
		case c > 0:
			gt--
			s[i], s[gt] = s[gt], s[i]
		default:
			i++
		}
	}
	return lt, gt
}

// medianOfThree returns the middle value of x, y, z.
func medianOfThree[T cmp.Ordered](x, y, z T) T {
	if cmp.Less(y, x) {
		x, y = y, x
	}
	if cmp.Less(z, y) {
		y = z
		if cmp.Less(y, x) {
			y = x
		}
	}
	return y
}

// medianOfMedians returns the median of the medians of the groups of 5 of s (the last group may be
// smaller). At least 3/10 of s is <= it and at least 3/10 is >= it. It reorders s. The recursion
// through selectRange works on len(s)/5 elements, so its depth is at most log5(len(s)) + 1.
func medianOfMedians[T cmp.Ordered](s []T) T {
	medians := 0
	for start := 0; start < len(s); start += 5 {
		end := min(start+5, len(s))
		insertionSort(s[start:end])
		mid := start + (end-start-1)/2
		// medians <= start, so this only overwrites groups already read.
		s[medians], s[mid] = s[mid], s[medians]
		medians++
	}
	k := (medians - 1) / 2
	selectRange(s[:medians], k, false)
	return s[k]
}

// insertionSort sorts s ascending in cmp.Compare order. Used only on ranges of at most insertionMax.
func insertionSort[T cmp.Ordered](s []T) {
	for i := 1; i < len(s); i++ {
		v := s[i]
		j := i
		for j > 0 && cmp.Less(v, s[j-1]) {
			s[j] = s[j-1]
			j--
		}
		s[j] = v
	}
}
