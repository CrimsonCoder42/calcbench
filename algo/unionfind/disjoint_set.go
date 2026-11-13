package unionfind

import (
	"errors"
	"fmt"
	"strconv"
)

// ErrNegativeSize: an element count or a grid dimension is less than 0.
var ErrNegativeSize = errors.New("unionfind: size must not be negative")

// ErrOutOfRange: an element is outside 0..Len()-1, or a cell is outside the grid.
var ErrOutOfRange = errors.New("unionfind: element out of range")

// ErrTooLarge: n is so large that the structure's arrays would exceed maxAllocBytes.
var ErrTooLarge = errors.New("unionfind: size too large to allocate")

// maxAllocBytes is the largest single allocation this package makes from a size argument: 2⁴⁷ bytes on
// 64-bit, 2³¹ on 32-bit, half the Go runtime's maxAlloc, above which makeslice panics (runtime/malloc.go;
// lead decision L41).
const maxAllocBytes = 1 << (31 + 16*(^uint(0)>>63))

// DisjointSet is a collection of disjoint sets over the elements 0..Len()-1, kept as a forest of parent
// pointers with union by rank and path compression. Each set is named by its root (its representative).
//
// The zero value is an empty, usable DisjointSet (MakeSet adds elements). On a nil *DisjointSet, Count
// and Len return 0, MakeSet returns -1, and Find, Union and Connected return ErrOutOfRange; no method
// panics. Find compresses paths, so even reads modify the structure: a DisjointSet is not safe for
// concurrent use.
type DisjointSet struct {
	parent []int
	rank   []uint8 // rank[r] <= floor(log2 n) <= 63 (CLRS Corollary 21.5), so uint8 never overflows
	count  int
}

// NewDisjointSet returns n singleton sets {0}, {1}, ..., {n-1}.
// Errors: ErrNegativeSize when n < 0; ErrTooLarge when the parent array (one int per element) would
// exceed 2⁴⁷ bytes (2³¹ on 32-bit).
//
// ID:      ALG-UNION-001
// Time:    O(n)   n = the element count
// Space:   O(n) for the structure (two allocations when n > 0)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.3 (MAKE-SET for each element)
// Source:  Galler, Fischer, "An improved equivalence algorithm", Communications of the ACM, 1964, 7(5):301–303, doi:10.1145/364099.364331
// Tags:    union-find, disjoint-set, dsu, make-set
func NewDisjointSet(n int) (*DisjointSet, error) {
	if n < 0 {
		return nil, fmt.Errorf("%w: n = %d", ErrNegativeSize, n)
	}
	if n > maxAllocBytes/(strconv.IntSize/8) { // the int parent array is the larger of the two
		return nil, fmt.Errorf("%w: n = %d", ErrTooLarge, n)
	}
	s := &DisjointSet{parent: make([]int, n), rank: make([]uint8, n), count: n}
	for i := range s.parent {
		s.parent[i] = i
	}
	return s, nil
}

// MakeSet adds a new element in a set of its own and returns it (the old Len()). On a nil *DisjointSet it
// adds nothing and returns -1.
//
// ID:      ALG-UNION-002
// Time:    O(1) amortized (the buffers double when full)
// Space:   O(1) extra amortized
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.3 (MAKE-SET)
// Tags:    union-find, disjoint-set, make-set
func (s *DisjointSet) MakeSet() int {
	if s == nil {
		return -1
	}
	x := len(s.parent)
	s.parent = append(s.parent, x)
	s.rank = append(s.rank, 0)
	s.count++
	return x
}

// Find returns the representative (root) of the set containing x, and points every element on the
// path from x to the root directly at the root (path compression).
// Errors: ErrOutOfRange when x is outside 0..Len()-1.
//
// ID:      ALG-UNION-003
// Time:    O(α(n)) amortized   n = Len(), α = the inverse Ackermann function (at most 4 for any real n); O(log n) worst case for one call
// Space:   O(1) extra (iterative, two passes over the path)
// Source:  Tarjan, "Efficiency of a good but not linear set union algorithm", Journal of the ACM, 1975, 22(2):215–225, doi:10.1145/321879.321884
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.3 (FIND-SET with path compression), §21.4 (the α bound)
// Tags:    union-find, disjoint-set, find, path-compression
func (s *DisjointSet) Find(x int) (int, error) {
	if s == nil || x < 0 || x >= len(s.parent) {
		return 0, fmt.Errorf("%w: x = %d, Len = %d", ErrOutOfRange, x, s.Len())
	}
	return s.find(x), nil
}

// find is Find without the range check. Pass 1 walks to the root; pass 2 repoints the path at it.
func (s *DisjointSet) find(x int) int {
	root := x
	for s.parent[root] != root {
		root = s.parent[root]
	}
	for s.parent[x] != root {
		next := s.parent[x]
		s.parent[x] = root
		x = next
	}
	return root
}

// Union merges the sets containing a and b and reports whether they were separate (merged true) or
// already one set (merged false, nothing changes). The root of lower rank is linked under the other
// (union by rank).
// Errors: ErrOutOfRange when a or b is outside 0..Len()-1.
//
// ID:      ALG-UNION-004
// Time:    O(α(n)) amortized   n = Len(), α = the inverse Ackermann function
// Space:   O(1) extra
// Source:  Tarjan, "Efficiency of a good but not linear set union algorithm", Journal of the ACM, 1975, 22(2):215–225, doi:10.1145/321879.321884
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.3 (UNION, LINK with union by rank)
// Tags:    union-find, disjoint-set, union, union-by-rank
func (s *DisjointSet) Union(a, b int) (merged bool, err error) {
	n := s.Len()
	if a < 0 || a >= n {
		return false, fmt.Errorf("%w: a = %d, Len = %d", ErrOutOfRange, a, n)
	}
	if b < 0 || b >= n {
		return false, fmt.Errorf("%w: b = %d, Len = %d", ErrOutOfRange, b, n)
	}
	return s.union(a, b), nil
}

// union is Union without the range checks.
func (s *DisjointSet) union(a, b int) bool {
	ra, rb := s.find(a), s.find(b)
	if ra == rb {
		return false
	}
	switch {
	case s.rank[ra] < s.rank[rb]:
		s.parent[ra] = rb
	case s.rank[ra] > s.rank[rb]:
		s.parent[rb] = ra
	default:
		s.parent[rb] = ra
		s.rank[ra]++
	}
	s.count--
	return true
}

// Connected reports whether a and b are in the same set.
// Errors: ErrOutOfRange when a or b is outside 0..Len()-1.
//
// ID:      ALG-UNION-005
// Time:    O(α(n)) amortized   n = Len(), α = the inverse Ackermann function
// Space:   O(1) extra
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.1 (SAME-COMPONENT)
// Source:  Tarjan, "Efficiency of a good but not linear set union algorithm", Journal of the ACM, 1975, 22(2):215–225, doi:10.1145/321879.321884
// Tags:    union-find, disjoint-set, connected, same-set
func (s *DisjointSet) Connected(a, b int) (bool, error) {
	n := s.Len()
	if a < 0 || a >= n {
		return false, fmt.Errorf("%w: a = %d, Len = %d", ErrOutOfRange, a, n)
	}
	if b < 0 || b >= n {
		return false, fmt.Errorf("%w: b = %d, Len = %d", ErrOutOfRange, b, n)
	}
	return s.find(a) == s.find(b), nil
}

// Count returns the number of disjoint sets (0 for a nil *DisjointSet).
//
// ID:      ALG-UNION-006
// Time:    O(1)
// Space:   O(1)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.1 (each UNION reduces the number of sets by one)
// Tags:    union-find, disjoint-set, count, components
func (s *DisjointSet) Count() int {
	if s == nil {
		return 0
	}
	return s.count
}

// Len returns the number of elements (0 for a nil *DisjointSet).
//
// ID:      ALG-UNION-007
// Time:    O(1)
// Space:   O(1)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.1 (n = the number of MAKE-SET operations)
// Tags:    union-find, disjoint-set, size
func (s *DisjointSet) Len() int {
	if s == nil {
		return 0
	}
	return len(s.parent)
}
