package reduce

// TwoSum returns s = fl(a + b) and the rounding error e, so that s + e = a + b exactly in real
// arithmetic (round to nearest, no overflow). It needs no comparison of |a| and |b|. When a + b overflows,
// s is ±Inf and e is NaN; inputs are not checked.
//
// ID:      ALG-RED-003
// Time:    O(1)   (six floating-point operations)
// Space:   O(1) (no allocation)
// Source:  Knuth, The Art of Computer Programming, Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1998, §4.2.2, Theorem B
// Source:  Higham, Accuracy and Stability of Numerical Algorithms, 2nd ed., SIAM, 2002, §4.3 (the error-free transformation TwoSum)
// Tags:    sum, error-free-transformation, two-sum, compensated, exact, no-alloc
func TwoSum(a, b float64) (s, e float64) {
	s = a + b
	bVirtual := s - a
	aVirtual := s - bVirtual
	e = (a - aVirtual) + (b - bVirtual)
	return s, e
}
