package text

import (
	"errors"
	"math"
	"math/rand/v2"
	"strconv"
	"testing"
)

// ISO/IEC 9899:2011 §7.22.1.4 ¶2–4 (N1570), applied by hand: "   -42xyz" splits into white space "   ",
// the subject sequence "-42" and the final string "xyz", so the value is −42 and 6 bytes are consumed.
func TestStringToInt_strtolDecomposition(t *testing.T) {
	v, n, err := StringToInt("   -42xyz")
	if v != -42 || n != 6 || err != nil {
		t.Errorf("StringToInt(\"   -42xyz\") = (%d, %d, %v), want (-42, 6, nil)", v, n, err)
	}
}

// Our own wording of the usual exercise cases, under the strtol rule (overflow is an error here).
func TestStringToInt_cases(t *testing.T) {
	tests := []struct {
		name  string
		s     string
		value int64
		n     int
	}{
		{"plain", "42", 42, 2},
		{"plus sign", "+7", 7, 2},
		{"zero", "0", 0, 1},
		{"negative zero", "-0", 0, 2},
		{"leading zeros", "000123", 123, 6},
		{"all six white-space bytes", " \t\n\v\f\r9", 9, 7},
		{"stops at the first non-digit", "12 34", 12, 2},
		{"no hex prefix in base 10", "0x1F", 0, 1},
		{"no decimal point", "3.99", 3, 1},
		{"maximum", "9223372036854775807", math.MaxInt64, 19},
		{"minimum", "-9223372036854775808", math.MinInt64, 20},
		{"long run of leading zeros does not overflow", "-00000000000000000000009", -9, 24},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, n, err := StringToInt(tt.s)
			if v != tt.value || n != tt.n || err != nil {
				t.Errorf("StringToInt(%q) = (%d, %d, %v), want (%d, %d, nil)", tt.s, v, n, err, tt.value, tt.n)
			}
		})
	}
}

func TestStringToInt_errors(t *testing.T) {
	tests := []struct {
		name string
		s    string
		n    int
		want error
	}{
		{"empty", "", 0, ErrNoDigits},
		{"only white space", "   ", 0, ErrNoDigits},
		{"sign alone", "  -", 0, ErrNoDigits},
		{"space after sign", "- 5", 0, ErrNoDigits},
		{"two signs", "+-5", 0, ErrNoDigits},
		{"letters first", "abc12", 0, ErrNoDigits},
		{"non-ASCII space is not skipped", " 5", 0, ErrNoDigits},
		{"one past the maximum", "9223372036854775808", 19, ErrOverflow},
		{"one past the minimum", "-9223372036854775809x", 20, ErrOverflow},
		{"far too long, n covers every digit", "123456789012345678901234567890!", 30, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, n, err := StringToInt(tt.s)
			if v != 0 || n != tt.n || !errors.Is(err, tt.want) {
				t.Errorf("StringToInt(%q) = (%d, %d, %v), want (0, %d, %v)", tt.s, v, n, err, tt.n, tt.want)
			}
		})
	}
}

func TestStringToInt_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { v, _, _ := StringToInt("  -123456"); benchSink = int(v) }); allocs != 0 {
		t.Errorf("StringToInt allocates %v times on valid input, want 0", allocs)
	}
}

// On strings that are just an optional sign and digits, StringToInt and strconv.ParseInt agree on the
// value and on overflow, and StringToInt consumes the whole string.
func TestStringToInt_oracleParseInt(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 600 {
		s := randomText(rng, 1+rng.IntN(21), "0123456789")
		if trial%3 == 1 {
			s = "-" + s
		} else if trial%3 == 2 {
			s = "+" + s
		}
		v, n, err := StringToInt(s)
		want, werr := strconv.ParseInt(s, 10, 64)
		switch {
		case werr != nil:
			if !errors.Is(err, ErrOverflow) || n != len(s) {
				t.Fatalf("trial %d: StringToInt(%q) = (%d, %d, %v); ParseInt says %v", trial, s, v, n, err, werr)
			}
		case v != want || n != len(s) || err != nil:
			t.Fatalf("trial %d: StringToInt(%q) = (%d, %d, %v), want (%d, %d, nil)", trial, s, v, n, err, want, len(s))
		}
	}
}

func BenchmarkStringToInt(b *testing.B) {
	b.Run("19digits", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			v, _, _ := StringToInt("  -9223372036854775807 rest")
			benchSink = int(v & 1)
		}
	})
}
