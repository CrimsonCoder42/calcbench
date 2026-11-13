package business

import (
	"errors"
	"math/rand/v2"
	"testing"
)

// Principles of Accounting Vol. 2 §3.2: Hicks Manufacturing $18,000 ÷ ($100 − $20) = 225 units; Marshall &
// Hirito $14,000 ÷ ($400 − $150) = 56 returns; College Creations $100,000 ÷ ($500 − $300) = 500 units.
// Exact (integer quotients of exact float64 values).
func TestBreakEvenUnits_sourceExample(t *testing.T) {
	cases := []struct {
		fixed, price, variable int64
		want                   float64
	}{
		{1800000, 10000, 2000, 225},
		{1400000, 40000, 15000, 56},
		{10000000, 50000, 30000, 500},
	}
	for _, c := range cases {
		got, err := BreakEvenUnits(c.fixed, c.price, c.variable)
		if err != nil || got != c.want {
			t.Errorf("%v: got %v, %v; want %v", c, got, err, c.want)
		}
	}
}

func TestBreakEvenUnits_table(t *testing.T) {
	cases := []struct {
		name                   string
		fixed, price, variable int64
		want                   float64
	}{
		{"no fixed costs", 0, 10, 5, 0},
		{"fractional units", 10, 4, 0, 2.5},
		{"one cent margin", maxExactCents, maxExactCents, maxExactCents - 1, maxExactCents},
	}
	for _, c := range cases {
		got, err := BreakEvenUnits(c.fixed, c.price, c.variable)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestBreakEvenUnits_errors(t *testing.T) {
	cases := []struct {
		name                   string
		fixed, price, variable int64
		want                   error
	}{
		{"zero margin", 1, 5, 5, ErrDomain},
		{"negative margin", 1, 5, 6, ErrDomain},
		{"fixed above 2^53", maxExactCents + 1, 2, 1, ErrOutOfRange},
		{"price above 2^53", 1, maxExactCents + 1, 1, ErrOutOfRange},
		{"variable above 2^53", 1, 2, maxExactCents + 1, ErrOutOfRange},
		{"negative fixed", -1, 2, 1, ErrNegative},
		{"negative price", 1, -2, 1, ErrNegative},
		{"negative variable", 1, 2, -1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := BreakEvenUnits(c.fixed, c.price, c.variable); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Oracle: fixed/(price − variable) against math/big; exact (all values exact in float64).
func TestBreakEvenUnits_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		fixed := int64(rng.Uint64N(maxExactCents + 1))
		variable := int64(rng.Uint64N(maxExactCents))
		price := variable + 1 + int64(rng.Uint64N(uint64(maxExactCents-variable)))
		got, err := BreakEvenUnits(fixed, price, variable)
		if err != nil {
			t.Fatal(err)
		}
		if want := ratOracle(fixed, price-variable); got != want {
			t.Fatalf("(%d,%d,%d) = %v, want %v", fixed, price, variable, got, want)
		}
	}
}

func TestBreakEvenUnits_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = BreakEvenUnits(1800000, 10000, 2000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkBreakEvenUnits(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = BreakEvenUnits(1800000, 10000, 2000)
	}
}
