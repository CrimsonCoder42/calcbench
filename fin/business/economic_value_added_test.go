package business

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Principles of Accounting Vol. 2 §12.3, new machinery project: after-tax income $250,000 × (1 − 0.40) =
// $150,000, invested capital $1,500,000, WACC 9 % → EVA $15,000. The same section's residual income uses
// the same computation: $250,000 − $1,500,000 × 0.18 = −$20,000. Exact cents.
func TestEconomicValueAdded_sourceExample(t *testing.T) {
	got, err := EconomicValueAdded(15000000, 150000000, 0.09)
	if err != nil || got != 1500000 {
		t.Fatalf("EVA: got %v, %v; want 1500000", got, err)
	}
	got, err = EconomicValueAdded(25000000, 150000000, 0.18)
	if err != nil || got != -2000000 {
		t.Fatalf("residual income: got %v, %v; want -2000000", got, err)
	}
}

func TestEconomicValueAdded_table(t *testing.T) {
	cases := []struct {
		name           string
		nopat, capital int64
		rate           float64
		want           int64
	}{
		{"zero rate", 100, 1000, 0, 100},
		{"zero capital", 100, 0, 0.5, 100},
		{"tie to even down", 0, 1, 0.5, 0},
		{"tie to even up", 0, 3, 0.5, -2},
		{"tie with odd offset", 1, 1, 0.5, 0},
		{"negative rate", 0, 100, -0.5, 50},
		{"max nopat", maxExactCents, 0, 0.1, maxExactCents},
		{"min nopat", -maxExactCents, 0, 0.1, -maxExactCents},
	}
	for _, c := range cases {
		got, err := EconomicValueAdded(c.nopat, c.capital, c.rate)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestEconomicValueAdded_errors(t *testing.T) {
	cases := []struct {
		name           string
		nopat, capital int64
		rate           float64
		want           error
	}{
		{"nopat above 2^53", maxExactCents + 1, 0, 0.1, ErrOutOfRange},
		{"capital above 2^53", 0, maxExactCents + 1, 0.1, ErrOutOfRange},
		{"negative capital", 0, -1, 0.1, ErrNegative},
		{"NaN rate", 0, 1, math.NaN(), ErrNotFinite},
		{"+Inf rate", 0, 1, math.Inf(1), ErrNotFinite},
		{"-Inf rate", 0, 1, math.Inf(-1), ErrNotFinite},
		{"rate -1", 0, 1, -1, ErrRate},
		{"result below -2^53", -maxExactCents, maxExactCents, 1, ErrOverflow},
		{"huge rate", 0, maxExactCents, 1e300, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := EconomicValueAdded(c.nopat, c.capital, c.rate); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// exactRound is the math/big oracle for the exact money results: k + sign·x·r, with r taken as the exact
// rational value of its float64, rounded half-even.
func exactRound(k, x int64, r float64, sign int64) int64 {
	prod := new(big.Rat).SetFloat64(r)
	prod.Mul(prod, new(big.Rat).SetInt64(x))
	if sign < 0 {
		prod.Neg(prod)
	}
	return halfEvenRat(prod.Add(prod, new(big.Rat).SetInt64(k)))
}

// Hidden ties: inputs where fl(x·r) ends in exactly ½ but the exact product does not (or the reverse).
// Rounding the float product alone would pick the wrong cent; the FMA-corrected rounding must match the
// exact oracle. The test fails if the search finds no such input, so it keeps covering the case.
func TestEconomicValueAdded_hiddenTies(t *testing.T) {
	const rate = 0.3
	found := 0
	for x := int64(1 << 52); x < 1<<52+4000; x++ {
		p := float64(float64(x) * rate) // explicit conversion forces the rounded product; without it the compiler may fuse p − Trunc(p) into an FMA
		if math.Abs(p-math.Trunc(p)) != 0.5 || math.FMA(float64(x), rate, -p) == 0 {
			continue
		}
		found++
		got, err := EconomicValueAdded(0, x, rate)
		if want := exactRound(0, x, rate, -1); err != nil || got != want {
			t.Fatalf("EVA(0, %d, %v) = %v, %v; want %v", x, rate, got, err, want)
		}
		got, err = FreeCashFlowToFirm(x, rate, 0, 0, 0)
		if want := exactRound(x, x, rate, -1); err != nil || got != want {
			t.Fatalf("FCFF(%d, %v) = %v, %v; want %v", x, rate, got, err, want)
		}
		got, err = FreeCashFlowToEquity(0, x, rate, 0)
		if want := exactRound(-x, x, rate, 1); err != nil || got != want {
			t.Fatalf("FCFE(0, %d, %v) = %v, %v; want %v", x, rate, got, err, want)
		}
	}
	if found == 0 {
		t.Fatal("no hidden tie in the search range")
	}
}

// Oracle at the full ±2⁵³ range: the result must equal the exact rational nopat − capital·rate rounded
// half-even (math/big), with no slack.
func TestEconomicValueAdded_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		nopat := int64(rng.Uint64N(2*maxExactCents+1)) - maxExactCents
		capital := int64(rng.Uint64N(maxExactCents + 1))
		rate := rng.Float64()*0.3 - 0.05
		if i%4 == 0 {
			capital = int64(rng.IntN(1000))
			rate = 0.5 // ties
		}
		want := exactRound(nopat, capital, rate, -1)
		got, err := EconomicValueAdded(nopat, capital, rate)
		if want > maxExactCents || want < -maxExactCents {
			if !errors.Is(err, ErrOverflow) {
				t.Fatalf("(%d,%d,%v): err = %v, want ErrOverflow", nopat, capital, rate, err)
			}
			continue
		}
		if err != nil || got != want {
			t.Fatalf("(%d,%d,%v) = %v, %v; want %v", nopat, capital, rate, got, err, want)
		}
	}
}

func TestEconomicValueAdded_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = EconomicValueAdded(15000000, 150000000, 0.09) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkEconomicValueAdded(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = EconomicValueAdded(15000000, 150000000, 0.09)
	}
}
