package hashing

// Integer is the set of Go's built-in integer types and the types defined on them.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// LongestConsecutive returns the start and length of the longest run of consecutive integers
// start, start+1, …, start+length-1 that all occur in values, in any order. Duplicates count once.
// Equal-length runs are broken by the smallest start. Empty input gives (0, 0). Runs never wrap past
// the type's maximum or minimum. values is read, never modified.
//
// ID:      ALG-HASH-009
// Time:    O(n) expected   n = len(values); each distinct value is looked up O(1) times
// Space:   O(n) extra for the set of values
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (expected O(1) SEARCH in a hash table)
// Tags:    hash, set, consecutive, run, sequence, streak
func LongestConsecutive[T Integer](values []T) (start T, length int) {
	present := make(map[T]struct{}, len(values))
	for _, v := range values {
		present[v] = struct{}{}
	}
	// Map order is unspecified, but the result is not: ties go to the smallest start.
	for v := range present {
		if prev := v - 1; prev < v { // prev < v is false only when v is the type's minimum
			if _, ok := present[prev]; ok {
				continue // v is inside a run; the run is counted from its first value
			}
		}
		runLen := 1
		for cur := v; ; runLen++ {
			next := cur + 1
			if next < cur { // wrapped past the type's maximum
				break
			}
			if _, ok := present[next]; !ok {
				break
			}
			cur = next
		}
		if runLen > length || (runLen == length && v < start) {
			start, length = v, runLen
		}
	}
	return start, length
}
