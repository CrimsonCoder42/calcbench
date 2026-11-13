package text

import "slices"

// RunLengthEncode appends to dst the PackBits encoding of src and returns the extended slice.
// Format (TIFF 6.0, Section 9): a header byte n read as int8, then data. n in 0…127: the next n+1 bytes
// are literal. n in −127…−1: the next byte is repeated 1−n times. n = −128 (0x80) is never written.
// Encoder policy, from the same section's advice: a run of 3 or more equal bytes is a replicate run; a
// run of 2 is a replicate run unless a literal run is in progress and the next run is a single byte, in
// which case it joins the literal; runs and literals longer than 128 bytes are split. The output is at
// most n + ⌊n/3⌋ + 1 bytes for n = len(src): each header costs one byte, and every literal that is not
// the last or a full 128 bytes is followed by a replicate run, so each such header covers at least 3
// input bytes. src is read, never modified. RunLengthDecode inverts it.
//
// ID:      ALG-TEXT-011
// Time:    O(n)   n = len(src)
// Space:   O(1) extra; at most n + ⌊n/3⌋ + 1 bytes for the result, one allocation only if dst lacks capacity
// Source:  Adobe Developers Association, TIFF Revision 6.0, Final, 3 June 1992, Section 9 "PackBits Compression", p. 42 (Apple Macintosh PackBits)
// Tags:    compression, run-length, rle, packbits, tiff, bytes
func RunLengthEncode(dst, src []byte) []byte {
	n := len(src)
	dst = slices.Grow(dst, n+n/3+1)
	litStart := -1 // start of the literal run in progress, or −1
	for i := 0; i < n; {
		j := i + 1
		for j < n && src[j] == src[i] {
			j++
		}
		run := j - i
		joinLiteral := run == 1 ||
			(run == 2 && litStart >= 0 && (j+1 == n || (j+1 < n && src[j+1] != src[j])))
		if joinLiteral {
			if litStart < 0 {
				litStart = i
			}
			for j-litStart >= 128 { // a full literal: write it, keep the rest open
				dst = appendLiteral(dst, src[litStart:litStart+128])
				litStart += 128
			}
			if litStart == j {
				litStart = -1
			}
			i = j
			continue
		}
		if litStart >= 0 {
			dst = appendLiteral(dst, src[litStart:i])
			litStart = -1
		}
		for run >= 2 {
			k := min(run, 128)
			dst = append(dst, byte(int8(1-k)), src[i])
			run -= k
		}
		if run == 1 { // a lone byte left after splitting a long run starts a literal
			litStart = j - 1
		}
		i = j
	}
	if litStart >= 0 {
		dst = appendLiteral(dst, src[litStart:n])
	}
	return dst
}

// appendLiteral appends one literal run of 1…128 bytes: the header len−1, then the bytes.
func appendLiteral(dst, lit []byte) []byte {
	dst = append(dst, byte(len(lit)-1))
	return append(dst, lit...)
}
