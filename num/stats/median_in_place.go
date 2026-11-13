package stats

// medianInPlace returns the median of finite, non-empty a (the mean of the two middle values when len(a)
// is even), reordering a. It exists so that MedianAbsoluteDeviation can select inside a buffer it
// already owns; algo/selection exports only selectors that copy their input (ALG-SELECT-001/002).
// Time O(n) worst case, O(1) extra space plus O(log n) stack, n = len(a).
func medianInPlace(a []float64) float64 {
	mid := (len(a) - 1) / 2
	selectFloatInPlace(a, mid)
	if len(a)%2 == 1 {
		return a[mid]
	}
	// Everything right of mid is ≥ a[mid], so the upper median is the least of it.
	upper := a[mid+1]
	for _, v := range a[mid+2:] {
		upper = min(upper, v)
	}
	return midpoint(a[mid], upper)
}

// selectInsertionMax is the range length at or below which insertion sort finishes the selection.
const selectInsertionMax = 16

// selectWorkBudget is how many times len(a) the median-of-3 rounds may partition in total before every
// later pivot is the median of medians (the value algo/selection uses; see its introselect.go).
const selectWorkBudget = 6

// selectFloatInPlace reorders finite a so that a[k] is the element a sorted copy would hold at index k,
// with a[:k] ≤ a[k] ≤ a[k+1:]. 0 ≤ k < len(a). Introselect: quickselect with a median-of-3 pivot until
// the partitioned lengths exceed selectWorkBudget·n, then median-of-medians pivots, each leaving at most
// 7/10 of the range, so O(n) worst case.
// Source: Musser, "Introspective sorting and selection algorithms", Software: Practice and Experience,
// 1997, 27(8):983–993; Blum, Floyd, Pratt, Rivest, Tarjan, "Time bounds for selection", JCSS, 1973,
// 7(4):448–461 (the same algorithm as ALG-SELECT-001, specialised to float64 and in place).
func selectFloatInPlace(a []float64, k int) {
	selectFloatRange(a, k, false)
}

// selectFloatRange is selectFloatInPlace with median-of-medians pivots from the start when forceMoM is
// true; tests use it to exercise the fallback on ordinary inputs.
func selectFloatRange(a []float64, k int, forceMoM bool) {
	lo, hi := 0, len(a)
	useMoM := forceMoM
	var work uint64
	budget := selectWorkBudget * uint64(len(a))
	for hi-lo > selectInsertionMax {
		var pivot float64
		if useMoM {
			pivot = floatMedianOfMedians(a[lo:hi])
		} else {
			pivot = floatMedianOfThree(a[lo], a[lo+(hi-lo)/2], a[hi-1])
			if work += uint64(hi - lo); work > budget {
				useMoM = true
			}
		}
		less, notGreater := floatPartitionThreeWay(a[lo:hi], pivot)
		switch {
		case k < lo+less:
			hi = lo + less
		case k >= lo+notGreater:
			lo += notGreater
		default:
			return // a[k] is in the block equal to the pivot
		}
	}
	floatInsertionSort(a[lo:hi])
}

// floatPartitionThreeWay reorders s into s[:less] < pivot, s[less:notGreater] == pivot and
// s[notGreater:] > pivot (Dijkstra's three-way partition) and returns the two boundaries.
func floatPartitionThreeWay(s []float64, pivot float64) (less, notGreater int) {
	lt, i, gt := 0, 0, len(s)
	for i < gt {
		switch {
		case s[i] < pivot:
			s[lt], s[i] = s[i], s[lt]
			lt++
			i++
		case s[i] > pivot:
			gt--
			s[i], s[gt] = s[gt], s[i]
		default:
			i++
		}
	}
	return lt, gt
}

// floatMedianOfThree returns the middle value of x, y, z.
func floatMedianOfThree(x, y, z float64) float64 {
	return max(min(x, y), min(max(x, y), z))
}

// floatMedianOfMedians returns the median of the medians of the groups of 5 of s, reordering s. At
// least 3/10 of s is ≤ it and at least 3/10 is ≥ it. Its recursion through selectFloatInPlace works on
// len(s)/5 elements, so the depth is at most log5(len(s)) + 1.
func floatMedianOfMedians(s []float64) float64 {
	medians := 0
	for start := 0; start < len(s); start += 5 {
		end := min(start+5, len(s))
		floatInsertionSort(s[start:end])
		mid := start + (end-start-1)/2
		s[medians], s[mid] = s[mid], s[medians] // medians ≤ start: only groups already read are overwritten
		medians++
	}
	k := (medians - 1) / 2
	selectFloatInPlace(s[:medians], k)
	return s[k]
}

// floatInsertionSort sorts s ascending; used only on ranges of at most selectInsertionMax.
func floatInsertionSort(s []float64) {
	for i := 1; i < len(s); i++ {
		v := s[i]
		j := i
		for j > 0 && v < s[j-1] {
			s[j] = s[j-1]
			j--
		}
		s[j] = v
	}
}
