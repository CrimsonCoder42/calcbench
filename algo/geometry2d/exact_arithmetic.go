package geometry2d

import "math/bits"

// The predicates here are exact on the whole int64 range. A difference of two int64 coordinates needs 65
// bits; it is carried as a uint64 magnitude and a sign. A product of two such magnitudes needs 128 bits
// (bits.Mul64), and a sum of two such products 129, held in 192.

// coordDiff returns a − b exactly, as its magnitude and whether it is negative. The true difference has
// magnitude below 2⁶⁴, so the wrapped uint64 subtraction is exact.
func coordDiff(a, b int64) (mag uint64, neg bool) {
	if a >= b {
		return uint64(a) - uint64(b), false
	}
	return uint64(b) - uint64(a), true
}

// uint128 is an unsigned 128-bit integer.
type uint128 struct{ hi, lo uint64 }

func mul64(x, y uint64) uint128 {
	hi, lo := bits.Mul64(x, y)
	return uint128{hi, lo}
}

func (a uint128) cmp(b uint128) int {
	switch {
	case a.hi != b.hi:
		if a.hi < b.hi {
			return -1
		}
		return 1
	case a.lo != b.lo:
		if a.lo < b.lo {
			return -1
		}
		return 1
	}
	return 0
}

func (a uint128) isZero() bool { return a.hi == 0 && a.lo == 0 }

// int129 is a signed product: a 128-bit magnitude and a sign; zero is never negative.
type int129 struct {
	mag uint128
	neg bool
}

// signedProduct returns (±m1)·(±m2) exactly.
func signedProduct(m1 uint64, neg1 bool, m2 uint64, neg2 bool) int129 {
	p := mul64(m1, m2)
	return int129{mag: p, neg: neg1 != neg2 && !p.isZero()}
}

// compareSigned returns the sign of p − q.
func compareSigned(p, q int129) int {
	switch {
	case p.neg && !q.neg:
		return -1
	case !p.neg && q.neg:
		return 1
	case p.neg:
		return -p.mag.cmp(q.mag)
	}
	return p.mag.cmp(q.mag)
}

// uint192 is an unsigned 192-bit integer; top holds the bits above 128.
type uint192 struct{ top, hi, lo uint64 }

func widen128(a uint128) uint192 { return uint192{hi: a.hi, lo: a.lo} }

func (a uint192) add(b uint192) uint192 {
	lo, c := bits.Add64(a.lo, b.lo, 0)
	hi, c := bits.Add64(a.hi, b.hi, c)
	top, _ := bits.Add64(a.top, b.top, c)
	return uint192{top, hi, lo}
}

// sub returns a − b; the callers guarantee b ≤ a.
func (a uint192) sub(b uint192) uint192 {
	lo, c := bits.Sub64(a.lo, b.lo, 0)
	hi, c := bits.Sub64(a.hi, b.hi, c)
	top, _ := bits.Sub64(a.top, b.top, c)
	return uint192{top, hi, lo}
}

func (a uint192) less(b uint192) bool {
	if a.top != b.top {
		return a.top < b.top
	}
	if a.hi != b.hi {
		return a.hi < b.hi
	}
	return a.lo < b.lo
}

// squaredDistance returns |p − q|² exactly (at most 2·(2⁶⁴ − 1)², which needs 129 bits).
func squaredDistance(p, q Point) uint192 {
	dx, _ := coordDiff(p.X, q.X)
	dy, _ := coordDiff(p.Y, q.Y)
	return widen128(mul64(dx, dx)).add(widen128(mul64(dy, dy)))
}
