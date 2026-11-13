package tree

import "errors"

// Sentinel errors shared by the package's types and functions.
var (
	// ErrRange: an index or bound is outside the structure.
	ErrRange = errors.New("tree: index out of range")
	// ErrOverflow: an integer sum does not fit its type, or a float sum became infinite.
	ErrOverflow = errors.New("tree: sum overflows")
	// ErrNotFinite: a floating-point value is NaN or ±Inf.
	ErrNotFinite = errors.New("tree: value is NaN or infinite")
	// ErrNilFunc: a required function argument is nil.
	ErrNilFunc = errors.New("tree: function is nil")
	// ErrEmptyRange: a query over an empty range on a structure with no identity element.
	ErrEmptyRange = errors.New("tree: empty range has no value")
)

// number is the set of types whose sums the ready-made sum and range-add structures keep.
type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// signed is the set of types the Fenwick structures accept: they subtract, so unsigned types are out.
type signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// finite reports whether x is not NaN and not ±Inf (always true for integers: x−x is 0).
func finite[T number](x T) bool {
	d := x - x
	return d == d && d == 0
}

// isFloat reports whether T is a floating-point type.
func isFloat[T number]() bool {
	var one, two T = 1, 2
	return one/two != 0
}

// addChecked returns a+b and false when the integer sum wraps or the float sum is not finite.
func addChecked[T number](a, b T) (T, bool) {
	s := a + b
	if !finite(s) || (b > 0 && s < a) || (b < 0 && s > a) {
		return s, false
	}
	return s, true
}

// subChecked returns a−b and false when the integer difference wraps or the float one is not finite.
func subChecked[T number](a, b T) (T, bool) {
	s := a - b
	if !finite(s) || (b > 0 && s > a) || (b < 0 && s < a) {
		return s, false
	}
	return s, true
}

// mulCount returns u·count (count >= 1) and false when it does not fit T or is not finite.
func mulCount[T number](u T, count int) (T, bool) {
	if u == 0 {
		return 0, true
	}
	c := T(count)
	if isFloat[T]() {
		p := u * c
		return p, finite(p)
	}
	if int(c) != count || c <= 0 {
		return 0, false // count does not fit T, and |u| >= 1
	}
	p := u * c
	if p/c != u {
		return p, false
	}
	return p, true
}
