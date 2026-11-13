package text

import (
	"errors"
	"fmt"
	"math"
)

// Errors returned by StringToInt and ColumnNumber.
var (
	// ErrNoDigits: after white space and an optional sign there is no decimal digit.
	ErrNoDigits = errors.New("text: no digits to convert")
	// ErrOverflow: the value does not fit the result type.
	ErrOverflow = errors.New("text: value out of range")
)

// StringToInt converts the leading decimal integer of s, the way C's strtol does with base 10, and returns
// the value and n, the number of bytes of s it consumed: s[n:] is the text after the number.
// Rule: skip ASCII white space (space, \t, \n, \v, \f, \r), then an optional '+' or '-', then the longest
// run of digits 0–9; anything after the digits is left unread. Overflow is an error, not a clamp.
// Unlike strconv.ParseInt, it accepts leading white space and trailing text and reports what it consumed.
// Errors: ErrNoDigits, with value 0 and n = 0 (as strtol's endptr = nptr), when there are no digits;
// ErrOverflow, with value 0 and n past the whole digit run, when the value is outside the int64 range.
//
// ID:      ALG-TEXT-015
// Time:    O(n)   n = len(s)
// Space:   O(1), no allocation
// Source:  ISO/IEC 9899:2011, Programming languages — C, §7.22.1.4 "The strtol, strtoll, strtoul, and strtoull functions", paragraphs 2–5 (read in committee draft N1570, 12 April 2011); white space per §7.4.1.10
// Tags:    string, parse, integer, atoi, strtol, overflow
func StringToInt(s string) (value int64, n int, err error) {
	i := 0
	for i < len(s) && isASCIISpace(s[i]) {
		i++
	}
	negative := false
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		negative = s[i] == '-'
		i++
	}
	start := i
	limit := uint64(math.MaxInt64) // the largest magnitude allowed
	if negative {
		limit++ // |MinInt64|
	}
	var magnitude uint64
	overflow := false
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		d := uint64(s[i] - '0')
		if !overflow {
			if magnitude > (limit-d)/10 {
				overflow = true
			} else {
				magnitude = magnitude*10 + d
			}
		}
		i++
	}
	if i == start {
		return 0, 0, fmt.Errorf("s: %w", ErrNoDigits)
	}
	if overflow {
		return 0, i, fmt.Errorf("s: %w", ErrOverflow)
	}
	if negative {
		return int64(-magnitude), i, nil // two's complement: −(1<<63) wraps to MinInt64
	}
	return int64(magnitude), i, nil
}
