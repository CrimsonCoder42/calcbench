package text

// IsJSONNumber reports whether all of s is a number in the grammar of RFC 8259 §6:
// number = [ "-" ] int [ frac ] [ exp ]; int = "0" / ( digit1-9 *DIGIT ); frac = "." 1*DIGIT;
// exp = ( "e" / "E" ) [ "-" / "+" ] 1*DIGIT. No white space, no leading '+', no leading zeros, no bare
// '.' at either end, no Infinity or NaN. Only the syntax is checked: 1e999 is a valid JSON number even
// though it overflows float64.
// It differs from strconv.ParseFloat, which also accepts "+1", ".5", "5.", "01", "Inf", "NaN" and hex
// floats, and rejects out-of-range values; and from json.Valid, which accepts any JSON value and
// surrounding white space.
//
// ID:      ALG-TEXT-016
// Time:    O(n)   n = len(s); one transition of a 9-state automaton per byte
// Space:   O(1), no allocation
// Source:  Bray (ed.), "The JavaScript Object Notation (JSON) Data Interchange Format", RFC 8259 (STD 90), IETF, December 2017, §6 "Numbers" (the ABNF above)
// Tags:    string, parse, validate, json, number, state-machine, automaton
func IsJSONNumber(s string) bool {
	const (
		start      = iota // nothing read
		afterMinus        // "-"
		zero              // int is "0"
		intDigits         // int is digit1-9 *DIGIT
		afterPoint        // "." without a digit yet
		fracDigits        // frac has a digit
		afterE            // "e" or "E"
		afterSign         // exponent sign
		expDigits         // exponent has a digit
	)
	state := start
	for i := range len(s) {
		c := s[i]
		digit := c >= '0' && c <= '9'
		switch state {
		case start, afterMinus:
			switch {
			case c == '-' && state == start:
				state = afterMinus
			case c == '0':
				state = zero
			case digit:
				state = intDigits
			default:
				return false
			}
		case zero, intDigits, fracDigits:
			switch {
			case digit && state != zero:
				// stay
			case c == '.' && state != fracDigits:
				state = afterPoint
			case c == 'e' || c == 'E':
				state = afterE
			default:
				return false
			}
		case afterPoint:
			if !digit {
				return false
			}
			state = fracDigits
		case afterE:
			switch {
			case c == '+' || c == '-':
				state = afterSign
			case digit:
				state = expDigits
			default:
				return false
			}
		case afterSign, expDigits:
			if !digit {
				return false
			}
			state = expDigits
		}
	}
	return state == zero || state == intDigits || state == fracDigits || state == expDigits
}
