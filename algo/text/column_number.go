package text

import (
	"errors"
	"fmt"
	"math"
)

// ErrColumnLabel: a spreadsheet column label is empty or holds a byte other than 'A'…'Z'.
var ErrColumnLabel = errors.New("text: column label must be one or more letters A-Z")

// ColumnNumber returns the number of the spreadsheet column label: "A" → 1, "Z" → 26, "AA" → 27,
// "XFD" → 16384. Labels are bijective base-26 numerals with digits A…Z worth 1…26; only upper-case A–Z
// is accepted (convert case first if needed). It is the inverse of ColumnLabel.
// Errors: ErrColumnLabel for an empty label or any other byte; ErrOverflow when the number exceeds
// math.MaxInt.
//
// ID:      ALG-TEXT-018
// Time:    O(n)   n = len(label)
// Space:   O(1), no allocation
// Source:  Foster, "A Number System without a Zero-Symbol", Mathematics Magazine, 1947, 21(1):39–41, doi:10.2307/3029479 (bijective base-k numeration)
// Tags:    string, spreadsheet, excel, column, bijective-numeration, base-26, parse
func ColumnNumber(label string) (int, error) {
	if len(label) == 0 {
		return 0, fmt.Errorf("label: %w", ErrColumnLabel)
	}
	number := 0
	for i := range len(label) {
		c := label[i]
		if c < 'A' || c > 'Z' {
			return 0, fmt.Errorf("label byte %d: %w", i, ErrColumnLabel)
		}
		d := int(c-'A') + 1
		if number > (math.MaxInt-d)/26 {
			return 0, fmt.Errorf("label: %w", ErrOverflow)
		}
		number = number*26 + d
	}
	return number, nil
}
