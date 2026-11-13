package business

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Principles of Finance §6.6, Clear Lake Sporting Goods: ($35,000 − $5,000) ÷ $95,000, printed 0.32.
// Exact against the correctly rounded 3,000,000 ÷ 9,500,000 and to 0.005 against the printed value.
func TestReturnOnCommonEquity_sourceExample(t *testing.T) {
	got, err := ReturnOnCommonEquity(3500000, 500000, 9500000)
	if err != nil {
		t.Fatal(err)
	}
	if want := 3000000.0 / 9500000.0; got != want || math.Abs(got-0.32) > 0.005 {
		t.Fatalf("got %v, want %v (printed 0.32)", got, want)
	}
}

func TestReturnOnCommonEquity_table(t *testing.T) {
	cases := []struct {
		name         string
		ni, pref, eq int64
		want         float64
	}{
		{"no preferred", 10, 0, 100, 0.1},
		{"loss", -10, 0, 100, -0.1},
		{"negative equity", 10, 0, -100, -0.1},
		{"preferred above income", 10, 30, 100, -0.2},
		{"extremes", -maxExactCents, maxExactCents, 2, -maxExactCents},
	}
	for _, c := range cases {
		got, err := ReturnOnCommonEquity(c.ni, c.pref, c.eq)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestReturnOnCommonEquity_errors(t *testing.T) {
	cases := []struct {
		name         string
		ni, pref, eq int64
		want         error
	}{
		{"zero equity", 1, 0, 0, ErrZeroDenominator},
		{"income above 2^53", maxExactCents + 1, 0, 1, ErrOutOfRange},
		{"income below -2^53", -maxExactCents - 1, 0, 1, ErrOutOfRange},
		{"preferred above 2^53", 1, maxExactCents + 1, 1, ErrOutOfRange},
		{"equity above 2^53", 1, 0, maxExactCents + 1, ErrOutOfRange},
		{"negative preferred", 1, -1, 1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := ReturnOnCommonEquity(c.ni, c.pref, c.eq); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Oracle against math/big: exact when |ni − pref| ≤ 2⁵³ (one rounding), else within 2⁻⁵² relative.
func TestReturnOnCommonEquity_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		ni := int64(rng.Uint64N(2*maxExactCents+1)) - maxExactCents
		pref := int64(rng.Uint64N(maxExactCents + 1))
		eq := int64(rng.Uint64N(2*maxExactCents+1)) - maxExactCents
		if eq == 0 {
			continue
		}
		got, err := ReturnOnCommonEquity(ni, pref, eq)
		if err != nil {
			t.Fatal(err)
		}
		want := ratOracle(ni-pref, eq)
		d := ni - pref
		if d >= -maxExactCents && d <= maxExactCents {
			if got != want {
				t.Fatalf("(%d,%d,%d) = %v, want %v", ni, pref, eq, got, want)
			}
		} else if math.Abs(got-want) > math.Abs(want)*0x1p-52 {
			t.Fatalf("(%d,%d,%d) = %v, want %v within 2^-52 relative", ni, pref, eq, got, want)
		}
	}
}

func TestReturnOnCommonEquity_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = ReturnOnCommonEquity(3500000, 500000, 9500000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkReturnOnCommonEquity(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ReturnOnCommonEquity(3500000, 500000, 9500000)
	}
}
