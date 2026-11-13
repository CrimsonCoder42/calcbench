package business

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Principles of Accounting Vol. 1 §6.6, California Business Solutions: ($293,500 − $180,000) ÷ $293,500,
// printed 0.39. Exact against the correctly rounded quotient 11,350,000 ÷ 29,350,000, and to 0.005 (half
// the printed last digit) against 0.39.
func TestGrossProfitMargin_sourceExample(t *testing.T) {
	got, err := GrossProfitMargin(29350000, 18000000)
	if err != nil {
		t.Fatal(err)
	}
	if want := 11350000.0 / 29350000.0; got != want || math.Abs(got-0.39) > 0.005 {
		t.Fatalf("got %v, want %v (printed 0.39)", got, want)
	}
}

func TestGrossProfitMargin_table(t *testing.T) {
	cases := []struct {
		name        string
		sales, cogs int64
		want        float64
	}{
		{"no cost", 100, 0, 1},
		{"break even", 100, 100, 0},
		{"loss", 100, 150, -0.5},
		{"max", maxExactCents, maxExactCents / 2, 0.5},
		{"cost at max, sales one", 1, maxExactCents, 1 - maxExactCents},
	}
	for _, c := range cases {
		got, err := GrossProfitMargin(c.sales, c.cogs)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestGrossProfitMargin_errors(t *testing.T) {
	cases := []struct {
		name        string
		sales, cogs int64
		want        error
	}{
		{"zero sales", 0, 1, ErrZeroDenominator},
		{"sales above 2^53", maxExactCents + 1, 1, ErrOutOfRange},
		{"cogs above 2^53", 1, maxExactCents + 1, ErrOutOfRange},
		{"negative sales", -1, 1, ErrNegative},
		{"negative cogs", 1, -1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := GrossProfitMargin(c.sales, c.cogs); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Oracle: (sales − cogs)/sales against big.Rat on 500 fixed-seed inputs. The difference is exact in
// float64 (|·| ≤ 2⁵³), so the quotient is correctly rounded and the comparison is exact.
func TestGrossProfitMargin_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		sales := 1 + int64(rng.Uint64N(maxExactCents))
		cogs := int64(rng.Uint64N(maxExactCents + 1))
		got, err := GrossProfitMargin(sales, cogs)
		if err != nil {
			t.Fatal(err)
		}
		want, _ := new(big.Rat).SetFrac(big.NewInt(sales-cogs), big.NewInt(sales)).Float64()
		if got != want {
			t.Fatalf("(%d, %d) = %v, want %v", sales, cogs, got, want)
		}
	}
}

func TestGrossProfitMargin_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = GrossProfitMargin(29350000, 18000000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkGrossProfitMargin(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = GrossProfitMargin(29350000, 18000000)
	}
}
