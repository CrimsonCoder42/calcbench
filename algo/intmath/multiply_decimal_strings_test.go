package intmath

import (
	"errors"
	"math/big"
	"math/rand/v2"
	"strings"
	"testing"
)

// Knuth Vol. 2 §4.3.1, Algorithm M in radix 10: 123 × 45 = 5535.
func TestMultiplyDecimalStrings_sourceExample(t *testing.T) {
	if got, err := MultiplyDecimalStrings("123", "45"); err != nil || got != "5535" {
		t.Errorf("MultiplyDecimalStrings(123, 45) = (%q, %v), want 5535", got, err)
	}
}

func TestMultiplyDecimalStrings_cases(t *testing.T) {
	tests := []struct{ a, b, want string }{
		{"0", "0", "0"}, {"0", "987", "0"}, {"987", "0", "0"}, {"1", "1", "1"}, {"9", "9", "81"},
		{"007", "006", "42"}, {"000", "5", "0"}, {"100", "100", "10000"}, {"99", "99", "9801"},
		{"12", "304", "3648"}, // a zero digit in the multiplier skips its row
		{strings.Repeat("9", 50), strings.Repeat("9", 50), strings.Repeat("9", 49) + "8" + strings.Repeat("0", 49) + "1"},
	}
	for _, tt := range tests {
		if got, err := MultiplyDecimalStrings(tt.a, tt.b); err != nil || got != tt.want {
			t.Errorf("MultiplyDecimalStrings(%q, %q) = (%q, %v), want %q", tt.a, tt.b, got, err, tt.want)
		}
	}
}

func TestMultiplyDecimalStrings_errors(t *testing.T) {
	tests := []struct {
		a, b    string
		wantErr error
	}{{"", "1", ErrEmpty}, {"1", "", ErrEmpty}, {"12a", "1", ErrInvalidDigit}, {"1", "-3", ErrInvalidDigit}, {"1.5", "2", ErrInvalidDigit}}
	for _, tt := range tests {
		if _, err := MultiplyDecimalStrings(tt.a, tt.b); !errors.Is(err, tt.wantErr) {
			t.Errorf("MultiplyDecimalStrings(%q, %q) error = %v, want %v", tt.a, tt.b, err, tt.wantErr)
		}
	}
}

// Oracle: math/big multiplication.
func TestMultiplyDecimalStrings_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 400 {
		a, b := randomDigits(rng, 1+rng.IntN(80), 10), randomDigits(rng, 1+rng.IntN(80), 10)
		x, _ := new(big.Int).SetString(a, 10)
		y, _ := new(big.Int).SetString(b, 10)
		want := new(big.Int).Mul(x, y).String()
		if got, err := MultiplyDecimalStrings(a, b); err != nil || got != want {
			t.Fatalf("MultiplyDecimalStrings(%q, %q) = (%q, %v), want %q", a, b, got, err, want)
		}
	}
}

func BenchmarkMultiplyDecimalStrings(b *testing.B) {
	// O(m·n): digit lengths 1e2, 1e3, 1e4 for both factors (1e4 × 1e4 = 1e8 steps).
	rng := rand.New(rand.NewPCG(1, 2))
	for _, size := range []struct {
		name string
		n    int
	}{{"1e2", 100}, {"1e3", 1_000}, {"1e4", 10_000}} {
		x, y := "1"+randomDigits(rng, size.n-1, 10), "1"+randomDigits(rng, size.n-1, 10)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				s, _ := MultiplyDecimalStrings(x, y)
				benchSink = uint64(len(s))
			}
		})
	}
}
