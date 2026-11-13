package business

import (
	"errors"
	"math"
	"testing"
)

// Contemporary Mathematics (OpenStax 2023) §6.2, Example 6.16: $62.00 at 15 % markup → $71.30.
func TestRetailPriceFromMarkup_sourceExample(t *testing.T) {
	got, err := RetailPriceFromMarkup(6200, 0.15)
	if err != nil || got != 7130 {
		t.Fatalf("got %d, %v; want 7130", got, err)
	}
}

func TestRetailPriceFromMarkup_table(t *testing.T) {
	cases := []struct {
		name string
		cost int64
		rate float64
		want int64
	}{
		{"no markup", 12345, 0, 12345},
		{"doubling", 12345, 1, 24690},
		{"zero cost", 0, 0.5, 0},
		{"half cent rounds to even (down)", 1, 0.5, 2},  // 1.5 → 2
		{"half cent rounds to even (up)", 5, 0.5, 8},    // 7.5 → 8
		{"half cent rounds to even (keep)", 9, 0.5, 14}, // 13.5 → 14
		{"cost 120 at 1/3", 12000, 1.0 / 3, 16000},
	}
	for _, c := range cases {
		got, err := RetailPriceFromMarkup(c.cost, c.rate)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestRetailPriceFromMarkup_errors(t *testing.T) {
	cases := []struct {
		name string
		cost int64
		rate float64
		want error
	}{
		{"negative cost", -1, 0.1, ErrNegative},
		{"cost beyond 2^53", maxExactCents + 1, 0.1, ErrOutOfRange},
		{"NaN rate", 100, math.NaN(), ErrNotFinite},
		{"Inf rate", 100, math.Inf(1), ErrNotFinite},
		{"negative rate", 100, -0.1, ErrRate},
		{"result beyond 2^53", maxExactCents, 1, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := RetailPriceFromMarkup(c.cost, c.rate); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestRetailPriceFromMarkup_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _ = RetailPriceFromMarkup(6200, 0.15) }); a != 0 {
		t.Errorf("allocations = %v, want 0", a)
	}
}

func BenchmarkRetailPriceFromMarkup(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = RetailPriceFromMarkup(6200, 0.15)
	}
}
