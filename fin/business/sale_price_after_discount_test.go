package business

import (
	"errors"
	"math"
	"testing"
)

// Contemporary Mathematics (OpenStax 2023) §6.2, Example 6.10: $75.80 at 25 % off → $56.85.
func TestSalePriceAfterDiscount_sourceExample(t *testing.T) {
	got, err := SalePriceAfterDiscount(7580, 0.25)
	if err != nil || got != 5685 {
		t.Fatalf("got %d, %v; want 5685", got, err)
	}
}

func TestSalePriceAfterDiscount_table(t *testing.T) {
	cases := []struct {
		name  string
		price int64
		rate  float64
		want  int64
	}{
		{"no discount", 12345, 0, 12345},
		{"free", 12345, 1, 0},
		{"half off, tie to even (down)", 5, 0.5, 2}, // 2.5 → 2
		{"half off, tie to even (up)", 7, 0.5, 4},   // 3.5 → 4
		{"largest price", maxExactCents, 0.5, maxExactCents / 2},
		{"zero price", 0, 0.3, 0},
	}
	for _, c := range cases {
		got, err := SalePriceAfterDiscount(c.price, c.rate)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestSalePriceAfterDiscount_errors(t *testing.T) {
	cases := []struct {
		name  string
		price int64
		rate  float64
		want  error
	}{
		{"negative price", -1, 0.1, ErrNegative},
		{"price beyond 2^53", maxExactCents + 1, 0.1, ErrOutOfRange},
		{"NaN rate", 100, math.NaN(), ErrNotFinite},
		{"Inf rate", 100, math.Inf(-1), ErrNotFinite},
		{"negative rate", 100, -0.1, ErrRate},
		{"rate above 1", 100, 1.01, ErrRate},
	}
	for _, c := range cases {
		if _, err := SalePriceAfterDiscount(c.price, c.rate); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestSalePriceAfterDiscount_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _ = SalePriceAfterDiscount(7580, 0.25) }); a != 0 {
		t.Errorf("allocations = %v, want 0", a)
	}
}

func BenchmarkSalePriceAfterDiscount(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = SalePriceAfterDiscount(7580, 0.25)
	}
}
