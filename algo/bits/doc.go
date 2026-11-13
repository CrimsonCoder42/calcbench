// Package bits holds bit-manipulation algorithms on machine words: power-of-two tests and rounding,
// isolating the lowest and highest set bit, bit reversal of a field, Gray code, a popcount table,
// XOR and modular bit counting to find the lone element, Hamming distance, submask enumeration and
// 2-D Morton (Z-order) codes.
//
// ID short form: BITS (IDs are ALG-BITS-NNN).
//
// Name: the standard library also has a package math/bits; a caller that needs both imports one under
// another name (import calcbits "github.com/CrimsonCoder42/calcbench/algo/bits"). The name is the one
// approved in doc/taxonomy.md (the same reasoning as L24 for algo/list).
//
// What is not here: a function that would be exactly one math/bits call is not repeated. The population
// count is math/bits.OnesCount64, the index of the lowest set bit is math/bits.TrailingZeros64, the index
// of the highest set bit is math/bits.Len64(x)-1, and a full 64-bit reversal is math/bits.Reverse64.
//
// Words: functions take uint64 unless the algorithm is the same for every integer type, in which case
// they are generic. No function allocates except PopCountTable when dst lacks capacity; no input slice
// is modified.
//
// Functions:
//
//	ALG-BITS-001  IsPowerOfTwo             whether x is 2^k for some k ≥ 0
//	ALG-BITS-002  IsPowerOfFour            whether x is 4^k for some k ≥ 0
//	ALG-BITS-003  NextPowerOfTwo           least power of two ≥ x (ErrOverflow above 2^63)
//	ALG-BITS-004  ReverseBits              the low width bits of x in reverse order
//	ALG-BITS-005  GrayEncode               the reflected binary Gray code of x
//	ALG-BITS-006  GrayDecode               the number whose Gray code is g
//	ALG-BITS-007  LowestSetBit             x with every bit but its lowest set bit cleared
//	ALG-BITS-008  HighestSetBit            x with every bit but its highest set bit cleared
//	ALG-BITS-009  PopCountTable            the number of set bits of 0, 1, …, n, in O(n)
//	ALG-BITS-010  LoneElementAmongPairs    the one value present an odd number of times when the rest are paired
//	ALG-BITS-011  LoneElementAmongTriples  the one value present once when every other is present three times
//	ALG-BITS-012  BitwiseAndRange          lo & (lo+1) & … & hi, in O(1)
//	ALG-BITS-013  HammingDistance          the number of bit positions where x and y differ
//	ALG-BITS-014  Submasks                 every submask of a mask, descending, as an iterator
//	ALG-BITS-015  MortonEncode2D           interleave the bits of (x, y) into a Z-order code
//	ALG-BITS-016  MortonDecode2D           split a Z-order code back into (x, y)
package bits
