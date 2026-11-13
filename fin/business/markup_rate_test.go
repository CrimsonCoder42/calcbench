package business

import (
	"errors"
	"testing"
)

// Contemporary Mathematics (OpenStax 2023) §6.2, Example 6.16: cost $62.00, retail $71.30 → 15 % markup.
func TestMarkupRate_sourceExample(t *testing.T) {
	got, err := MarkupRate(6200, 7130)
	if err != nil || got != 930.0/6200.0 {
		t.Fatalf("got %v, %v; want %v", got, err, 930.0/6200.0)
	}
}

func TestMarkupRate_table(t *testing.T) {
	cases := []struct {
		name        string
		cost, price int64
		want        float64
	}{
		{"cost 120, price 160 → 1/3", 12000, 16000, 4000.0 / 12000.0},
		{"no markup", 500, 500, 0},
		{"sold below cost", 1000, 800, -0.2},
		{"free item", 1000, 0, -1},
		{"largest values", maxExactCents, maxExactCents, 0},
	}
	for _, c := range cases {
		got, err := MarkupRate(c.cost, c.price)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestMarkupRate_errors(t *testing.T) {
	cases := []struct {
		name        string
		cost, price int64
		want        error
	}{
		{"zero cost", 0, 100, ErrZeroDenominator},
		{"negative cost", -1, 100, ErrNegative},
		{"negative price", 100, -1, ErrNegative},
		{"cost beyond 2^53", maxExactCents + 1, 100, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := MarkupRate(c.cost, c.price); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestMarkupRate_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _ = MarkupRate(6200, 7130) }); a != 0 {
		t.Errorf("allocations = %v, want 0", a)
	}
}

func BenchmarkMarkupRate(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = MarkupRate(6200, 7130)
	}
}
