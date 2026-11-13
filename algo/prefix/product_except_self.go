package prefix

import (
	"fmt"
	"math/bits"
	"slices"
)

// ProductExceptSelf appends to dst, for each i, the product of every element of values except
// values[i], and returns the extended slice (n results for n values). No division is used: result i is
// the exclusive prefix product of values[:i] times the exclusive suffix product of values[i+1:]. The
// empty product is 1, so a single value gives [1].
// Zeros are counted first: two or more make every result 0; exactly one makes every result 0 except at
// the zero's index. Products are carried as a 64-bit magnitude and a sign, so overflow is reported
// exactly: a result is in error only when its true value does not fit T.
// values is read, never modified. dst must not overlap values.
// Errors: ErrOverflow when a result does not fit T (its index is in the message). On error dst is
// returned with its length unchanged.
//
// ID:      ALG-PREFIX-010
// Time:    O(n)   n = len(values); two passes, at most 2n 64×64-bit multiplications
// Space:   O(1) extra (the prefix products are stored in the result); O(n) for the result, one allocation only if dst lacks capacity
// Source:  Blelloch, "Prefix Sums and Their Applications", in Reif (ed.), Synthesis of Parallel Algorithms, Morgan Kaufmann, 1993, ch. 1, §1.1 (the prescan, an exclusive scan under any associative operator with identity; here × with identity 1, run forwards and backwards); also CMU-CS-90-190, 1990
// Tags:    prefix-product, suffix-product, scan, no-division, integer, overflow-checked
func ProductExceptSelf[T integer](dst, values []T) ([]T, error) {
	n := len(values)
	zeros, zeroAt, negatives := 0, -1, 0
	for i, v := range values {
		if v == 0 {
			zeros++
			zeroAt = i
		} else if _, neg := magnitude(v); neg {
			negatives++
		}
	}
	base := len(dst)
	out := slices.Grow(dst, n)[:base+n]
	res := out[base:]
	switch {
	case zeros >= 2:
		clear(res)
		return out, nil
	case zeros == 1:
		// Only the zero's own result can be non-zero: the product of all the others.
		mag, neg, ok := uint64(1), negatives%2 == 1, true
		for i, v := range values {
			if i == zeroAt {
				continue
			}
			m, _ := magnitude(v)
			var hi uint64
			hi, mag = bits.Mul64(mag, m)
			if hi != 0 {
				ok = false
				break
			}
		}
		var p T
		if ok {
			p, ok = fromMagnitude[T](mag, neg)
		}
		if !ok {
			return dst, fmt.Errorf("%w: product except values[%d]", ErrOverflow, zeroAt)
		}
		clear(res)
		res[zeroAt] = p
		return out, nil
	}
	// No zeros: every |values[i]| >= 1, so a partial product's magnitude never shrinks and a result is
	// at least as large as either of its two partial products.
	// Pass 1: res[i] = |values[0] × … × values[i-1]|, stored as its negation for a signed T, whose
	// negative range is one larger (so |MinInt| is storable). If it cannot be stored, result i and every
	// later result exceed the type.
	prefix, prefixBig := uint64(1), false
	for i, v := range values {
		var stored T
		ok := !prefixBig
		if ok {
			stored, ok = fromMagnitude[T](prefix, isSigned[T]())
		}
		if !ok {
			return dst, fmt.Errorf("%w: product except values[%d]", ErrOverflow, i)
		}
		res[i] = stored
		m, _ := magnitude(v)
		var hi uint64
		hi, prefix = bits.Mul64(prefix, m)
		prefixBig = prefixBig || hi != 0
	}
	// Pass 2, right to left: multiply in the suffix magnitude and give each result its sign, the parity
	// of the negative values other than values[i].
	suffix, suffixBig := uint64(1), false
	for i := n - 1; i >= 0; i-- {
		v := values[i]
		m, vNeg := magnitude(v)
		pm, _ := magnitude(res[i])
		hi, mag := bits.Mul64(pm, suffix)
		ok := !suffixBig && hi == 0
		var r T
		if ok {
			others := negatives
			if vNeg {
				others--
			}
			r, ok = fromMagnitude[T](mag, others%2 == 1)
		}
		if !ok {
			return dst, fmt.Errorf("%w: product except values[%d]", ErrOverflow, i)
		}
		res[i] = r
		hi, suffix = bits.Mul64(suffix, m)
		suffixBig = suffixBig || hi != 0
	}
	return out, nil
}
