package business

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Hand-checked (the CFA ratio list and Regulation S-K give no numbers): EBIT $50,000, lease payments
// $10,000, interest $20,000 → (50,000 + 10,000) ÷ (20,000 + 10,000) = 60,000 ÷ 30,000 = 2. Exact: the
// quotient is an integer.
func TestFixedChargeCoverageRatio_sourceExample(t *testing.T) {
	got, err := FixedChargeCoverageRatio(5000000, 1000000, 2000000)
	if err != nil || got != 2 {
		t.Fatalf("got %v, %v; want 2", got, err)
	}
}

// Exact comparisons: every sum is within ±2⁵³, so the quotient is correctly rounded and equals the
// float64 literal of the same rational.
func TestFixedChargeCoverageRatio_table(t *testing.T) {
	cases := []struct {
		name                string
		ebit, fixed, intExp int64
		want                float64
	}{
		{"no lease equals times interest earned", 4300000, 0, 200000, 21.5},
		{"no interest", 1000, 1000, 0, 2},
		{"zero ebit", 0, 10, 30, 0.25},
		{"negative ebit", -30, 10, 10, -1},
		{"ebit offsets lease", -10, 10, 10, 0},
		{"max ebit", maxExactCents - 1, 1, 0, maxExactCents},
		{"min ebit", -maxExactCents, 0, 1, -maxExactCents},
	}
	for _, c := range cases {
		got, err := FixedChargeCoverageRatio(c.ebit, c.fixed, c.intExp)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestFixedChargeCoverageRatio_errors(t *testing.T) {
	cases := []struct {
		name                string
		ebit, fixed, intExp int64
		want                error
	}{
		{"zero denominator", 1, 0, 0, ErrZeroDenominator},
		{"ebit above 2^53", maxExactCents + 1, 1, 1, ErrOutOfRange},
		{"ebit below -2^53", -maxExactCents - 1, 1, 1, ErrOutOfRange},
		{"lease above 2^53", 1, maxExactCents + 1, 1, ErrOutOfRange},
		{"interest max int64", 1, 1, math.MaxInt64, ErrOutOfRange},
		{"negative lease", 1, -1, 1, ErrNegative},
		{"negative interest", 1, 1, -1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := FixedChargeCoverageRatio(c.ebit, c.fixed, c.intExp); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Against math/big on 500 fixed-seed inputs. Tolerance 2 ulp (relative 2⁻⁵¹): a sum above 2⁵³ in
// magnitude is rounded once when converted to float64, before the correctly rounded division.
func TestFixedChargeCoverageRatio_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		ebit := int64(rng.Uint64N(2*maxExactCents+1)) - maxExactCents
		fixed := int64(rng.Uint64N(maxExactCents + 1))
		intExp := int64(rng.Uint64N(maxExactCents + 1))
		if fixed+intExp == 0 {
			continue
		}
		got, err := FixedChargeCoverageRatio(ebit, fixed, intExp)
		if err != nil {
			t.Fatalf("(%d,%d,%d): %v", ebit, fixed, intExp, err)
		}
		want := ratOracle(ebit+fixed, intExp+fixed)
		if math.Abs(got-want) > math.Abs(want)*0x1p-51 {
			t.Fatalf("(%d,%d,%d) = %v, want %v", ebit, fixed, intExp, got, want)
		}
	}
}

func TestFixedChargeCoverageRatio_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = FixedChargeCoverageRatio(5000000, 1000000, 2000000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkFixedChargeCoverageRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = FixedChargeCoverageRatio(5000000, 1000000, 2000000)
	}
}
