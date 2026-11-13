package stack

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

// ErrPairs: the bracket pairs are not a list of distinct opener/closer rune pairs.
var ErrPairs = errors.New("stack: pairs must be opener/closer runes, each rune used once, valid UTF-8, no U+FFFD")

// BalancedBrackets reports whether the brackets in text are balanced: every closer matches the most
// recent unmatched opener of its pair, and no opener is left unmatched. pairs lists the brackets as
// opener/closer rune pairs, for example "()[]{}" or "<>«»". Runes of text that are in no pair are
// ignored; bytes of text that are not valid UTF-8 are ignored too (they read as U+FFFD, which no pair
// may hold). text and pairs are read, never modified.
// Errors: ErrPairs when pairs has an odd number of runes, an opener equal to its closer, a rune in more
// than one place, U+FFFD, or invalid UTF-8.
//
// ID:      ALG-STACK-010
// Time:    O(n + p)   n = runes of text, p = runes of pairs
// Space:   O(n) extra worst case for the stack of open brackets (none when text has no opener)
// Source:  Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §1.3, Exercise 1.3.4 (Parentheses)
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.1 (stacks)
// Tags:    stack, brackets, parentheses, balanced, matching, delimiters, dyck
func BalancedBrackets(text, pairs string) (bool, error) {
	if !utf8.ValidString(pairs) {
		return false, fmt.Errorf("%w: pairs is not valid UTF-8", ErrPairs)
	}
	// role of a rune: 0 none, +k opener of pair k, -k closer of pair k (k from 1). A table for ASCII,
	// a map only when pairs holds other runes.
	var ascii [utf8.RuneSelf]int32
	var other map[rune]int32
	setRole := func(r rune, role int32) error {
		if r == utf8.RuneError {
			return fmt.Errorf("%w: pairs holds U+FFFD", ErrPairs)
		}
		if r < utf8.RuneSelf {
			if ascii[r] != 0 {
				return fmt.Errorf("%w: %q is used twice", ErrPairs, r)
			}
			ascii[r] = role
			return nil
		}
		if other == nil {
			other = make(map[rune]int32)
		}
		if _, dup := other[r]; dup {
			return fmt.Errorf("%w: %q is used twice", ErrPairs, r)
		}
		other[r] = role
		return nil
	}
	pairCount := int32(0)
	for rest := pairs; rest != ""; {
		opener, w1 := utf8.DecodeRuneInString(rest)
		rest = rest[w1:]
		if rest == "" {
			return false, fmt.Errorf("%w: %q has no closer", ErrPairs, opener)
		}
		closer, w2 := utf8.DecodeRuneInString(rest)
		rest = rest[w2:]
		if opener == closer {
			return false, fmt.Errorf("%w: %q opens and closes", ErrPairs, opener)
		}
		pairCount++
		if err := setRole(opener, pairCount); err != nil {
			return false, err
		}
		if err := setRole(closer, -pairCount); err != nil {
			return false, err
		}
	}

	var open []int32 // pair numbers of the unmatched openers
	for _, r := range text {
		role := int32(0)
		if r < utf8.RuneSelf {
			role = ascii[r]
		} else if other != nil {
			role = other[r]
		}
		switch {
		case role > 0:
			open = append(open, role)
		case role < 0:
			if len(open) == 0 || open[len(open)-1] != -role {
				return false, nil
			}
			open = open[:len(open)-1]
		}
	}
	return len(open) == 0, nil
}
