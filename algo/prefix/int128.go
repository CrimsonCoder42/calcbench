package prefix

import "math/bits"

// integer is the set of Go's built-in integer types and the types defined on them. Every one fits in
// 64 bits, so a sum or difference of a few of them fits exactly in int128.
type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// int128 is a two's-complement 128-bit integer: hi holds the upper 64 bits (signed), lo the lower 64.
// It carries exact sums of int64 and uint64 values; additions never overflow it before 2⁶³ terms.
type int128 struct {
	hi int64
	lo uint64
}

// isSigned reports whether T is a signed integer type: ^0 is -1 for a signed type and the maximum for
// an unsigned one.
func isSigned[T integer]() bool {
	var zero T
	return ^zero < zero
}

// widen returns v as an int128, exactly.
func widen[T integer](v T) int128 {
	if isSigned[T]() {
		x := int64(v)
		return int128{hi: x >> 63, lo: uint64(x)}
	}
	return int128{lo: uint64(v)}
}

func (a int128) add(b int128) int128 {
	lo, carry := bits.Add64(a.lo, b.lo, 0)
	return int128{hi: a.hi + b.hi + int64(carry), lo: lo}
}

func (a int128) sub(b int128) int128 {
	lo, borrow := bits.Sub64(a.lo, b.lo, 0)
	return int128{hi: a.hi - b.hi - int64(borrow), lo: lo}
}

// narrow returns a as a T and true when a is in T's range, or (0, false).
func narrow[T integer](a int128) (T, bool) {
	if isSigned[T]() {
		x := int64(a.lo)
		if a.hi != x>>63 {
			return 0, false
		}
		t := T(x)
		if int64(t) != x {
			return 0, false
		}
		return t, true
	}
	if a.hi != 0 {
		return 0, false
	}
	t := T(a.lo)
	if uint64(t) != a.lo {
		return 0, false
	}
	return t, true
}

// checkedAdd returns a + b and true, or the wrapped sum and false when the true sum is outside T.
func checkedAdd[T integer](a, b T) (T, bool) {
	s := a + b
	if (b > 0 && s < a) || (b < 0 && s > a) {
		return s, false
	}
	return s, true
}

// checkedSub returns a - b and true, or the wrapped difference and false when the true difference is
// outside T.
func checkedSub[T integer](a, b T) (T, bool) {
	d := a - b
	if (b > 0 && d > a) || (b < 0 && d < a) {
		return d, false
	}
	return d, true
}

// magnitude returns |v| and whether v is negative. |math.MinInt64| = 2⁶³ fits in uint64.
func magnitude[T integer](v T) (uint64, bool) {
	if isSigned[T]() {
		x := int64(v)
		if x < 0 {
			return uint64(-x), true // -MinInt64 wraps to MinInt64, whose uint64 is 2⁶³: correct
		}
		return uint64(x), false
	}
	return uint64(v), false
}

// fromMagnitude returns the T whose magnitude is mag and whose sign is negative when neg, and true;
// or (0, false) when that value is outside T.
func fromMagnitude[T integer](mag uint64, neg bool) (T, bool) {
	w := int128{lo: mag}
	if neg {
		w = int128{}.sub(w)
	}
	return narrow[T](w)
}
