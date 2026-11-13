package text

import (
	"errors"
	"fmt"
)

// Errors returned by RunLengthDecode.
var (
	// ErrTruncated: a PackBits header is not followed by all the data it announces.
	ErrTruncated = errors.New("text: run-length data ends inside a run")
	// ErrTooLong: the decoded output would exceed the caller's cap.
	ErrTooLong = errors.New("text: decoded output exceeds the length cap")
	// ErrMaxLen: the length cap is negative.
	ErrMaxLen = errors.New("text: length cap must not be negative")
)

// RunLengthDecode appends to dst the bytes that the PackBits data src decodes to and returns the extended
// slice. Format (TIFF 6.0, Section 9), header byte n read as int8: 0…127 copies the next n+1 bytes;
// −127…−1 repeats the next byte 1−n times; −128 is a no-op. Every PackBits stream is accepted, not only
// RunLengthEncode's. maxLen caps the number of bytes this call appends (a run can expand 64-fold, so
// untrusted input needs a cap); no single count can overflow, since a header announces at most 128 bytes.
// src is read, never modified.
// Errors: ErrMaxLen when maxLen < 0; ErrTruncated when src ends inside a run; ErrTooLong when the output
// would exceed maxLen. On error dst is returned with its length unchanged.
//
// ID:      ALG-TEXT-012
// Time:    O(n + k)   n = len(src), k = bytes appended (at most maxLen)
// Space:   O(1) extra; O(k) for the result
// Source:  Adobe Developers Association, TIFF Revision 6.0, Final, 3 June 1992, Section 9 "PackBits Compression", p. 42 (the unpacking loop)
// Tags:    compression, run-length, rle, packbits, tiff, bytes, decode
func RunLengthDecode(dst, src []byte, maxLen int) ([]byte, error) {
	if maxLen < 0 {
		return dst, fmt.Errorf("maxLen %d: %w", maxLen, ErrMaxLen)
	}
	base := len(dst)
	for i := 0; i < len(src); {
		n := int(int8(src[i]))
		i++
		switch {
		case n >= 0:
			k := n + 1
			if len(src)-i < k {
				return dst[:base], fmt.Errorf("src offset %d: %w", i-1, ErrTruncated)
			}
			if len(dst)-base > maxLen-k {
				return dst[:base], fmt.Errorf("maxLen %d: %w", maxLen, ErrTooLong)
			}
			dst = append(dst, src[i:i+k]...)
			i += k
		case n == -128:
			// no-op
		default:
			k := 1 - n
			if i >= len(src) {
				return dst[:base], fmt.Errorf("src offset %d: %w", i-1, ErrTruncated)
			}
			if len(dst)-base > maxLen-k {
				return dst[:base], fmt.Errorf("maxLen %d: %w", maxLen, ErrTooLong)
			}
			c := src[i]
			i++
			for range k {
				dst = append(dst, c)
			}
		}
	}
	return dst, nil
}
