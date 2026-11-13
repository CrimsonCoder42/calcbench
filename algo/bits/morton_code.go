package bits

// spreadBits32 moves bit i of v to bit 2i, leaving the odd positions zero, in five mask-and-shift steps.
func spreadBits32(v uint32) uint64 {
	x := uint64(v)
	x = (x | x<<16) & 0x0000_FFFF_0000_FFFF
	x = (x | x<<8) & 0x00FF_00FF_00FF_00FF
	x = (x | x<<4) & 0x0F0F_0F0F_0F0F_0F0F
	x = (x | x<<2) & 0x3333_3333_3333_3333
	x = (x | x<<1) & 0x5555_5555_5555_5555
	return x
}

// compactBits32 is the inverse of spreadBits32: it gathers the even bits 2i of x into bit i.
func compactBits32(x uint64) uint32 {
	x &= 0x5555_5555_5555_5555
	x = (x | x>>1) & 0x3333_3333_3333_3333
	x = (x | x>>2) & 0x0F0F_0F0F_0F0F_0F0F
	x = (x | x>>4) & 0x00FF_00FF_00FF_00FF
	x = (x | x>>8) & 0x0000_FFFF_0000_FFFF
	x = (x | x>>16) & 0x0000_0000_FFFF_FFFF
	return uint32(x)
}

// MortonEncode2D returns the Morton (Z-order) code of the point (x, y): bit i of x goes to bit 2i and
// bit i of y to bit 2i+1. Sorting points by their codes visits them along the Z-order curve, which keeps
// nearby points mostly near each other (spatial indexing, quadtrees, cache-friendly grids).
//
// ID:      ALG-BITS-015
// Time:    O(log w)   w = 32 bits per coordinate; five mask-and-shift steps each
// Space:   O(1), no allocation
// Source:  Morton, G. M., "A computer oriented geodetic data base and a new technique in file sequencing", IBM Ltd., Ottawa, 1966 (technical report; the interleaved-bit file sequence)
// Source:  Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §7-2 "Shuffling Bits" (outer perfect shuffle by masks and shifts in log2 w steps)
// Source:  Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques" (bit interleaving, the "zipper" of two words, by δ-shifts)
// Tags:    bits, morton, z-order, interleave, spatial-index, quadtree
func MortonEncode2D(x, y uint32) uint64 {
	return spreadBits32(x) | spreadBits32(y)<<1
}

// MortonDecode2D returns the point (x, y) whose Morton code is code: the even bits of code form x and the
// odd bits form y. MortonDecode2D(MortonEncode2D(x, y)) = (x, y) for every x, y.
//
// ID:      ALG-BITS-016
// Time:    O(log w)   w = 32 bits per coordinate; five mask-and-shift steps each
// Space:   O(1), no allocation
// Source:  Morton, G. M., "A computer oriented geodetic data base and a new technique in file sequencing", IBM Ltd., Ottawa, 1966
// Source:  Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §7-2 "Shuffling Bits" (inverse perfect shuffle)
// Tags:    bits, morton, z-order, deinterleave, spatial-index
func MortonDecode2D(code uint64) (x, y uint32) {
	return compactBits32(code), compactBits32(code >> 1)
}
