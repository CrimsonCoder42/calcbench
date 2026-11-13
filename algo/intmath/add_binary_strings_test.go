package intmath

import (
	"errors"
	"math/big"
	"math/rand/v2"
	"strings"
	"testing"
)

// Knuth Vol. 2 §4.3.1, Algorithm A in radix 2: 1011 + 11 = 1110 (11 + 3 = 14).
func TestAddBinaryStrings_sourceExample(t *testing.T) {
	if got, err := AddBinaryStrings("1011", "11"); err != nil || got != "1110" {
		t.Errorf("AddBinaryStrings(1011, 11) = (%q, %v), want 1110", got, err)
	}
}

func TestAddBinaryStrings_cases(t *testing.T) {
	tests := []struct{ a, b, want string }{
		{"0", "0", "0"}, {"1", "0", "1"}, {"1", "1", "10"}, {"111", "1", "1000"}, {"0001", "0001", "10"},
		{"000", "0", "0"}, {"10", "1011", "1101"},
		{strings.Repeat("1", 200), "1", "1" + strings.Repeat("0", 200)},
	}
	for _, tt := range tests {
		if got, err := AddBinaryStrings(tt.a, tt.b); err != nil || got != tt.want {
			t.Errorf("AddBinaryStrings(%q, %q) = (%q, %v), want %q", tt.a, tt.b, got, err, tt.want)
		}
	}
}

func TestAddBinaryStrings_errors(t *testing.T) {
	tests := []struct {
		a, b    string
		wantErr error
	}{{"", "1", ErrEmpty}, {"1", "", ErrEmpty}, {"102", "1", ErrInvalidDigit}, {"1", "1a", ErrInvalidDigit}, {"-1", "1", ErrInvalidDigit}, {" 1", "1", ErrInvalidDigit}}
	for _, tt := range tests {
		if _, err := AddBinaryStrings(tt.a, tt.b); !errors.Is(err, tt.wantErr) {
			t.Errorf("AddBinaryStrings(%q, %q) error = %v, want %v", tt.a, tt.b, err, tt.wantErr)
		}
	}
}

// randomDigits returns a string of n random digits below radix, possibly with leading zeros.
func randomDigits(rng *rand.Rand, n int, radix int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('0' + rng.IntN(radix))
	}
	return string(b)
}

// Oracle: math/big addition.
func TestAddBinaryStrings_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 400 {
		a, b := randomDigits(rng, 1+rng.IntN(150), 2), randomDigits(rng, 1+rng.IntN(150), 2)
		x, _ := new(big.Int).SetString(a, 2)
		y, _ := new(big.Int).SetString(b, 2)
		want := new(big.Int).Add(x, y).Text(2)
		if got, err := AddBinaryStrings(a, b); err != nil || got != want {
			t.Fatalf("AddBinaryStrings(%q, %q) = (%q, %v), want %q", a, b, got, err, want)
		}
	}
}

func BenchmarkAddBinaryStrings(b *testing.B) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range benchSizes {
		x, y := "1"+randomDigits(rng, size.n-1, 2), "1"+randomDigits(rng, size.n-1, 2)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				s, _ := AddBinaryStrings(x, y)
				benchSink = uint64(len(s))
			}
		})
	}
}
