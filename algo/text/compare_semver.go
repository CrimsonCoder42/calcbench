package text

import (
	"errors"
	"fmt"
	"strings"
)

// ErrSemVer: a string is not a valid Semantic Versioning 2.0.0 version.
var ErrSemVer = errors.New("text: not a valid semantic version")

// CompareSemVer returns −1, 0 or +1 as version a has lower, equal or higher precedence than version b
// under Semantic Versioning 2.0.0. Both must match the specification's grammar: MAJOR.MINOR.PATCH,
// non-negative integers with no leading zeros (§2); an optional pre-release "-" followed by non-empty
// dot-separated identifiers of [0-9A-Za-z-], numeric ones without leading zeros (§9); optional build
// metadata "+" followed by non-empty dot-separated identifiers of [0-9A-Za-z-] (§10). No "v" prefix.
// Precedence (§11): major, minor and patch compare numerically (of any size; no overflow); a version
// with a pre-release is lower than the same version without; pre-release identifiers compare left to
// right, numeric ones numerically, others in ASCII order, numeric below alphanumeric, and a longer list
// is higher when all earlier identifiers are equal. Build metadata is ignored, so 0 does not mean the
// strings are equal.
// Errors: ErrSemVer, naming a or b, when either does not match the grammar.
//
// ID:      ALG-TEXT-019
// Time:    O(n)   n = len(a) + len(b)
// Space:   O(1), no allocation on valid input
// Source:  Preston-Werner, Semantic Versioning 2.0.0, https://semver.org/spec/v2.0.0.html (read 2026-11-12), items 2, 9, 10 and 11 and the Backus–Naur form grammar
// Tags:    string, version, semver, compare, precedence, parse
func CompareSemVer(a, b string) (int, error) {
	va, ok := parseSemVer(a)
	if !ok {
		return 0, fmt.Errorf("a: %w", ErrSemVer)
	}
	vb, ok := parseSemVer(b)
	if !ok {
		return 0, fmt.Errorf("b: %w", ErrSemVer)
	}
	for k := range 3 {
		if c := compareNumeric(va.core[k], vb.core[k]); c != 0 {
			return c, nil
		}
	}
	switch {
	case !va.hasPre && !vb.hasPre:
		return 0, nil
	case !va.hasPre:
		return 1, nil
	case !vb.hasPre:
		return -1, nil
	}
	pa, pb := va.pre, vb.pre
	for {
		if pa == "" || pb == "" {
			switch {
			case pa == "" && pb == "":
				return 0, nil
			case pa == "":
				return -1, nil
			default:
				return 1, nil
			}
		}
		var ia, ib string
		ia, pa = nextIdentifier(pa)
		ib, pb = nextIdentifier(pb)
		na, nb := isDigits(ia), isDigits(ib)
		var c int
		switch {
		case na && nb:
			c = compareNumeric(ia, ib)
		case na:
			c = -1
		case nb:
			c = 1
		default:
			c = strings.Compare(ia, ib)
		}
		if c != 0 {
			return c, nil
		}
	}
}

// semVer holds the parts of a parsed version as substrings of the input.
type semVer struct {
	core   [3]string // major, minor, patch digits
	pre    string    // pre-release identifiers joined by '.', without the '-'
	hasPre bool
}

// parseSemVer splits and validates v against the Semantic Versioning 2.0.0 grammar.
func parseSemVer(v string) (semVer, bool) {
	var out semVer
	if i := strings.IndexByte(v, '+'); i >= 0 {
		if !validIdentifiers(v[i+1:], false) {
			return out, false
		}
		v = v[:i]
	}
	if i := strings.IndexByte(v, '-'); i >= 0 {
		out.pre, out.hasPre = v[i+1:], true
		if !validIdentifiers(out.pre, true) {
			return out, false
		}
		v = v[:i]
	}
	for k := range 3 {
		var part string
		if k < 2 {
			i := strings.IndexByte(v, '.')
			if i < 0 {
				return out, false
			}
			part, v = v[:i], v[i+1:]
		} else {
			part, v = v, ""
		}
		if !isDigits(part) || (len(part) > 1 && part[0] == '0') {
			return out, false
		}
		out.core[k] = part
	}
	return out, true
}

// validIdentifiers reports whether s is one or more non-empty '.'-separated identifiers of [0-9A-Za-z-];
// with noLeadingZero, an all-digit identifier must not start with '0' unless it is "0".
func validIdentifiers(s string, noLeadingZero bool) bool {
	if s == "" {
		return false
	}
	for s != "" {
		var id string
		id, s = nextIdentifier(s)
		if id == "" {
			return false
		}
		for i := range len(id) {
			c := id[i]
			if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '-') {
				return false
			}
		}
		if noLeadingZero && isDigits(id) && len(id) > 1 && id[0] == '0' {
			return false
		}
	}
	return true
}

// nextIdentifier splits s at its first '.': the identifier before it, and the rest after it. A trailing
// '.' leaves rest = "." so that the empty identifier after it is still seen by the caller.
func nextIdentifier(s string) (id, rest string) {
	i := strings.IndexByte(s, '.')
	if i < 0 {
		return s, ""
	}
	if i == len(s)-1 {
		return s[:i], "."
	}
	return s[:i], s[i+1:]
}

// isDigits reports whether s is non-empty and all ASCII digits.
func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := range len(s) {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// compareNumeric compares two digit strings without leading zeros by value: longer is larger, then ASCII.
func compareNumeric(x, y string) int {
	if len(x) != len(y) {
		if len(x) < len(y) {
			return -1
		}
		return 1
	}
	return strings.Compare(x, y)
}
