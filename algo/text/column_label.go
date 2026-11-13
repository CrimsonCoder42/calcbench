package text

import (
	"errors"
	"fmt"
)

// ErrColumnNumber: a spreadsheet column number is less than 1.
var ErrColumnNumber = errors.New("text: column number must be at least 1")

// ColumnLabel returns the spreadsheet column label of number: 1 → "A", 26 → "Z", 27 → "AA", 702 → "ZZ",
// 703 → "AAA". This is bijective base-26 numeration with digits A…Z worth 1…26 (there is no zero digit),
// so every positive integer has exactly one label. ColumnNumber is its inverse.
// Errors: ErrColumnNumber when number < 1.
//
// ID:      ALG-TEXT-017
// Time:    O(log₂₆ number)   at most 14 digits for an int64
// Space:   O(1) extra; one allocation for the result string
// Source:  Foster, "A Number System without a Zero-Symbol", Mathematics Magazine, 1947, 21(1):39–41, doi:10.2307/3029479 (bijective base-k numeration)
// Tags:    string, spreadsheet, excel, column, bijective-numeration, base-26, conversion
func ColumnLabel(number int) (string, error) {
	if number < 1 {
		return "", fmt.Errorf("number %d: %w", number, ErrColumnNumber)
	}
	var buf [16]byte
	i := len(buf)
	for v := uint64(number); v > 0; {
		v-- // shift digits 1…26 to 0…25
		i--
		buf[i] = byte('A' + v%26)
		v /= 26
	}
	return string(buf[i:]), nil
}
