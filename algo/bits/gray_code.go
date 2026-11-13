package bits

// GrayEncode returns the reflected binary Gray code of x: x XOR (x >> 1). Consecutive integers have
// codes that differ in exactly one bit, and GrayDecode(GrayEncode(x)) = x for every x.
//
// ID:      ALG-BITS-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Gray, F., "Pulse code communication", U.S. Patent 2,632,058, filed 13 Nov 1947, issued 17 Mar 1953 (the reflected binary code)
// Source:  Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.2.1.1 "Generating all n-tuples", Gray binary code g(k) = k ⊕ ⌊k/2⌋
// Tags:    bits, gray-code, reflected-binary, encoding
func GrayEncode(x uint64) uint64 {
	return x ^ x>>1
}

// GrayDecode returns the number whose reflected binary Gray code is g: bit i of the result is the XOR
// of bits i, i+1, …, 63 of g, computed as a prefix XOR in six shift-and-XOR steps.
//
// ID:      ALG-BITS-006
// Time:    O(log w)   w = 64 bits; six shifts
// Space:   O(1), no allocation
// Source:  Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.2.1.1 "Generating all n-tuples" (inverse Gray code: k = g ⊕ ⌊g/2⌋ ⊕ ⌊g/4⌋ ⊕ …)
// Source:  Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, ch. 13 "Gray Code", §13-1 (conversion from Gray code by a prefix XOR in log2 w steps)
// Tags:    bits, gray-code, reflected-binary, decoding, prefix-xor
func GrayDecode(g uint64) uint64 {
	g ^= g >> 1
	g ^= g >> 2
	g ^= g >> 4
	g ^= g >> 8
	g ^= g >> 16
	g ^= g >> 32
	return g
}
