package hashing

import (
	"encoding/binary"
	"errors"
	"fmt"
	"unicode/utf8"
)

// Errors of GroupAnagramsByCount. They are wrapped with the offending input; test with errors.Is.
var (
	// ErrAlphabet: the alphabet repeats a rune, holds U+FFFD, or is not valid UTF-8.
	ErrAlphabet = errors.New("hashing: alphabet must be valid UTF-8 with no repeated rune and no U+FFFD")
	// ErrNotInAlphabet: a word holds a rune that is not in the alphabet.
	ErrNotInAlphabet = errors.New("hashing: word holds a rune outside the alphabet")
)

// GroupAnagramsByCount returns words grouped into anagram classes, like GroupAnagrams, using as each
// word's signature its count of every rune of alphabet instead of its sorted runes. The groups are the
// same as GroupAnagrams' whenever every word is spelled from alphabet. Groups are in the order of their
// first word; words keep input order, duplicates included. Bytes that are not valid UTF-8 read as
// U+FFFD, which no alphabet may hold, so such a word is an error. words is read, never modified.
// Errors: ErrAlphabet when alphabet repeats a rune, holds U+FFFD or is not valid UTF-8;
// ErrNotInAlphabet when a word holds a rune outside alphabet (or a byte that is not valid UTF-8).
//
// ID:      ALG-HASH-011
// Time:    O(Σ Lᵢ + n·σ) expected   Lᵢ = runes in words[i], n = len(words), σ = runes in alphabet
// Space:   O(n·σ) extra for the signatures (a uvarint per alphabet rune); O(n) for the result
// Source:  Bentley, Programming Pearls, 2nd ed., Addison-Wesley, 2000, Column 2, problem C (anagram signatures, §2.8, pp. 18–20)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §8.2 (counting) and ch. 11 (hash tables)
// Tags:    hash, anagram, signature, counting, group, string, deterministic
func GroupAnagramsByCount(words []string, alphabet string) ([][]string, error) {
	if !utf8.ValidString(alphabet) {
		return nil, fmt.Errorf("%w: alphabet is not valid UTF-8", ErrAlphabet)
	}
	// slot maps a rune to its position in alphabet: a table for ASCII, a map only for other runes.
	var ascii [utf8.RuneSelf]int32
	for i := range ascii {
		ascii[i] = -1
	}
	var other map[rune]int32
	sigma := int32(0)
	for _, r := range alphabet {
		if r == utf8.RuneError {
			return nil, fmt.Errorf("%w: alphabet holds U+FFFD", ErrAlphabet)
		}
		if r < utf8.RuneSelf {
			if ascii[r] >= 0 {
				return nil, fmt.Errorf("%w: %q repeats", ErrAlphabet, r)
			}
			ascii[r] = sigma
		} else {
			if other == nil {
				other = make(map[rune]int32)
			}
			if _, dup := other[r]; dup {
				return nil, fmt.Errorf("%w: %q repeats", ErrAlphabet, r)
			}
			other[r] = sigma
		}
		sigma++
	}

	groups := [][]string{}
	index := make(map[string]int)
	counts := make([]uint64, sigma)
	var key []byte
	for wi, w := range words {
		clear(counts)
		for _, r := range w {
			s := int32(-1)
			if r < utf8.RuneSelf {
				s = ascii[r]
			} else if slot, ok := other[r]; ok {
				s = slot
			}
			if s < 0 {
				return nil, fmt.Errorf("%w: words[%d] holds %q", ErrNotInAlphabet, wi, r)
			}
			counts[s]++
		}
		key = key[:0]
		for _, c := range counts {
			key = binary.AppendUvarint(key, c)
		}
		g, ok := index[string(key)] // a lookup with string(key) does not allocate
		if !ok {
			g = len(groups)
			index[string(key)] = g
			groups = append(groups, nil)
		}
		groups[g] = append(groups[g], w)
	}
	return groups, nil
}
