package interval

import "fmt"

// SummaryRanges returns ascending integer values as the closed ranges of consecutive integers they
// form: 0 1 2 4 5 7 → [0,2] [4,5] [7,7]. values must be ascending; duplicates are allowed and join their
// range. No values give an empty, non-nil slice. values is read, never modified. No overflow at the
// extremes of T.
// Errors: ErrNotSorted, with the index, when values is not ascending.
//
// ID:      ALG-INTVL-006
// Time:    O(n)   n = len(values)
// Space:   O(1) extra; O(r) for the result, r = number of ranges ≤ n (amortized appends)
// Source:  Klee, "Can the measure of ∪[aᵢ,bᵢ] be computed in less than O(n log n) steps?", American Mathematical Monthly, 1977, 84(4):284–285, doi:10.1080/00029890.1977.11994336 (union of intervals by one sweep in sorted order; here the unit intervals [v, v+1), which join when consecutive)
// Tags:    interval, ranges, runs, consecutive, integers, compress, sorted
func SummaryRanges[T integer](values []T) ([]Interval[T], error) {
	out := []Interval[T]{}
	for i, v := range values {
		if i == 0 {
			out = append(out, Interval[T]{v, v})
			continue
		}
		last := &out[len(out)-1]
		switch {
		case v < last.End:
			return nil, fmt.Errorf("%w: values[%d] = %v < values[%d] = %v", ErrNotSorted, i, v, i-1, values[i-1])
		case v == last.End:
			// duplicate: already in the range
		case v-last.End == 1: // v > last.End, so the difference cannot be a wrapped 1
			last.End = v
		default:
			out = append(out, Interval[T]{v, v})
		}
	}
	return out, nil
}
