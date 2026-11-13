package tree

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
)

// ErrMalformed: the bytes given to Deserialize are not in the format Serialize writes.
var ErrMalformed = errors.New("tree: malformed serialized tree")

// emptyMark is the byte that stands for an empty subtree in the serialized form.
const emptyMark = '#'

// Serialize appends to dst the tree at root in this package's text format and returns the extended
// slice. The format is the preorder sequence of the tree in which every empty subtree is written as '#'
// and every node as a netstring of its value's bytes: the byte count in decimal (no leading zeros), ':',
// the bytes, ','. The empty tree is "#"; a single node whose value encodes as "7" is "1:7,##"; the tree
// 1 → (2, 3) is "1:1,1:2,##1:3,##". The value bytes are whatever appendValue appends for it, so any byte
// may occur in them. dst == nil allocates. The tree is read, never modified.
// Errors: ErrNilFunc when appendValue is nil (dst is returned unchanged).
//
// ID:      ALG-TREE-019
// Time:    O(n + B)   n = nodes, B = total bytes of the encoded values
// Space:   O(h) extra (an explicit stack; no recursion)   h = height; O(n + B) appended to dst
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §14.4 "Binary trees", p. 139 (preorder), https://cses.fi/book/book.pdf (read 2026-11-12); empty-subtree markers make the preorder unique (catalog/tree.md)
// Source:  Bernstein, "Netstrings", 1997, https://cr.yp.to/proto/netstrings.txt (read 2026-11-12) (the value framing)
// Tags:    tree, binary-tree, serialize, encode, preorder, netstring
func Serialize[T any](dst []byte, root *Node[T], appendValue func(dst []byte, v T) []byte) ([]byte, error) {
	if appendValue == nil {
		return dst, fmt.Errorf("%w: appendValue", ErrNilFunc)
	}
	stack := []*Node[T]{root} // nil entries are empty subtrees still to write
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if n == nil {
			dst = append(dst, emptyMark)
			continue
		}
		// The length is known only after the value is written: write value, then "len:", and rotate the
		// prefix to the front in place (three reversals, no allocation).
		start := len(dst)
		dst = appendValue(dst, n.Value)
		length := len(dst) - start
		dst = strconv.AppendInt(dst, int64(length), 10)
		dst = append(dst, ':')
		slices.Reverse(dst[start : start+length])
		slices.Reverse(dst[start+length:])
		slices.Reverse(dst[start:])
		dst = append(dst, ',')
		stack = append(stack, n.Right, n.Left)
	}
	return dst, nil
}

// Deserialize returns the tree that Serialize wrote into data, decoding each value with parseValue. The
// bytes passed to parseValue are a sub-slice of data: parseValue must copy them if it keeps them.
// The empty tree ("#") is (nil, nil). data is read, never modified.
// Errors: ErrNilFunc when parseValue is nil; ErrMalformed (with the byte offset) for a bad or
// leading-zero length, a missing ':' or ',', a length past the end, data that ends before the tree is
// complete, or bytes left after it; an error from parseValue, wrapped with the offset.
//
// ID:      ALG-TREE-020
// Time:    O(L + Σ parse)   L = len(data); plus parseValue's own cost per node
// Space:   O(h) extra (an explicit stack of open slots; no recursion)   h = height; O(n) for the result
// Source:  Laaksonen, Competitive Programmer's Handbook, draft of 3 July 2018, §14.4 "Binary trees", p. 139 (preorder), https://cses.fi/book/book.pdf (read 2026-11-12)
// Source:  Bernstein, "Netstrings", 1997, https://cr.yp.to/proto/netstrings.txt (read 2026-11-12) (the value framing)
// Tags:    tree, binary-tree, deserialize, decode, parse, preorder, netstring
func Deserialize[T any](data []byte, parseValue func(b []byte) (T, error)) (*Node[T], error) {
	if parseValue == nil {
		return nil, fmt.Errorf("%w: parseValue", ErrNilFunc)
	}
	var root *Node[T]
	slots := []**Node[T]{&root} // where the next subtrees go, the next one on top
	pos := 0
	for len(slots) > 0 {
		slot := slots[len(slots)-1]
		slots = slots[:len(slots)-1]
		if pos >= len(data) {
			return nil, fmt.Errorf("%w: data ends at byte %d with %d subtrees missing", ErrMalformed, pos, len(slots)+1)
		}
		if data[pos] == emptyMark {
			pos++
			continue // *slot stays nil
		}
		start := pos
		length := 0
		for pos < len(data) && data[pos] >= '0' && data[pos] <= '9' {
			if length > (len(data)-pos)/10 {
				return nil, fmt.Errorf("%w: length at byte %d exceeds the data", ErrMalformed, start)
			}
			length = length*10 + int(data[pos]-'0')
			pos++
		}
		switch {
		case pos == start:
			return nil, fmt.Errorf("%w: byte %d is %q, want a digit or '#'", ErrMalformed, pos, data[pos])
		case data[start] == '0' && pos-start > 1:
			return nil, fmt.Errorf("%w: length at byte %d has a leading zero", ErrMalformed, start)
		case pos >= len(data) || data[pos] != ':':
			return nil, fmt.Errorf("%w: want ':' at byte %d", ErrMalformed, pos)
		}
		pos++
		if length > len(data)-pos-1 {
			return nil, fmt.Errorf("%w: value of %d bytes at byte %d runs past the end", ErrMalformed, length, pos)
		}
		if data[pos+length] != ',' {
			return nil, fmt.Errorf("%w: want ',' at byte %d", ErrMalformed, pos+length)
		}
		v, err := parseValue(data[pos : pos+length])
		if err != nil {
			return nil, fmt.Errorf("tree: value at byte %d: %w", pos, err)
		}
		pos += length + 1
		nd := &Node[T]{Value: v}
		*slot = nd
		slots = append(slots, &nd.Right, &nd.Left)
	}
	if pos != len(data) {
		return nil, fmt.Errorf("%w: %d bytes left after the tree at byte %d", ErrMalformed, len(data)-pos, pos)
	}
	return root, nil
}
