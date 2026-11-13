package business

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// eoqTolerance is the relative tolerance against the 256-bit oracle: S ÷ H, √(2·S/H), √D and the product
// each round once (≤ ½ ulp; a square root halves its input's relative error), so the total stays below
// 4 ulp of the result.
const eoqTolerance = 4 * 0x1p-52

// eoqOracle returns √(2·D·S ÷ H) in 256-bit math/big, rounded to float64.
func eoqOracle(d float64, s, h int64) float64 {
	const prec = 256
	x := new(big.Float).SetPrec(prec).SetFloat64(d)
	x.Mul(x, new(big.Float).SetPrec(prec).SetInt64(2*s))
	x.Quo(x, new(big.Float).SetPrec(prec).SetInt64(h))
	f, _ := x.Sqrt(x).Float64()
	return f
}

// Harris 1915, "What Quantity to Make at Once", p. 51: movement M = 1,000 units a month, set-up cost $2,
// unit cost $0.10, 10 % a year for interest and depreciation → "the most economical size of lot is 2,190
// units". With D = 12·M = 12,000 a year, S = 200¢ and H = 10 % of 10¢ = 1¢ a unit a year:
// √(2·12,000·200 ÷ 1) = √4,800,000 = 2,190.890… (Harris truncates). Tolerance as eoqTolerance.
func TestEconomicOrderQuantity_sourceExample(t *testing.T) {
	got, err := EconomicOrderQuantity(12000, 200, 1)
	if err != nil {
		t.Fatal(err)
	}
	want := math.Sqrt(4800000)
	if math.Abs(got-want) > want*eoqTolerance || math.Floor(got) != 2190 {
		t.Fatalf("got %v, want %v (Harris: 2,190)", got, want)
	}
}

// Perfect squares give exact roots; compared with eoqTolerance (the product of two roots can be off by
// an ulp even when the exact value is an integer).
func TestEconomicOrderQuantity_table(t *testing.T) {
	cases := []struct {
		name string
		d    float64
		s, h int64
		want float64
	}{
		{"zero demand", 0, 200, 1, 0},
		{"zero order cost", 1000, 0, 5, 0},
		{"perfect square", 1000, 500, 1, 1000},
		{"holding above order cost", 100, 1, 2, 10},
		{"fractional demand", 0.5, 4, 1, 2},
		{"max order cost", 2, maxExactCents, 1, math.Sqrt(4) * math.Sqrt(maxExactCents)},
		{"max holding cost", maxExactCents, 1, maxExactCents, math.Sqrt(2)},
		{"huge demand stays finite", math.MaxFloat64, maxExactCents, 1, math.Sqrt(2*maxExactCents) * math.Sqrt(math.MaxFloat64)},
		{"tiny demand", math.SmallestNonzeroFloat64, 1, 2, math.Sqrt(math.SmallestNonzeroFloat64)},
	}
	for _, c := range cases {
		got, err := EconomicOrderQuantity(c.d, c.s, c.h)
		if err != nil || math.IsInf(got, 0) || math.Abs(got-c.want) > c.want*eoqTolerance {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestEconomicOrderQuantity_errors(t *testing.T) {
	cases := []struct {
		name string
		d    float64
		s, h int64
		want error
	}{
		{"NaN demand", math.NaN(), 1, 1, ErrNotFinite},
		{"+Inf demand", math.Inf(1), 1, 1, ErrNotFinite},
		{"-Inf demand", math.Inf(-1), 1, 1, ErrNotFinite},
		{"negative demand", -1, 1, 1, ErrNegative},
		{"negative order cost", 1, -1, 1, ErrNegative},
		{"negative holding cost", 1, 1, -1, ErrNegative},
		{"order cost above 2^53", 1, maxExactCents + 1, 1, ErrOutOfRange},
		{"holding cost above 2^53", 1, 1, maxExactCents + 1, ErrOutOfRange},
		{"holding cost max int64", 1, 1, math.MaxInt64, ErrOutOfRange},
		{"zero holding cost", 1, 1, 0, ErrZeroDenominator},
	}
	for _, c := range cases {
		if _, err := EconomicOrderQuantity(c.d, c.s, c.h); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Against the 256-bit oracle on 500 fixed-seed inputs, demand spread over 1e-6 … 1e12 units a year.
func TestEconomicOrderQuantity_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		d := math.Pow(10, -6+18*rng.Float64())
		s := int64(rng.Uint64N(maxExactCents + 1))
		h := 1 + int64(rng.Uint64N(maxExactCents))
		got, err := EconomicOrderQuantity(d, s, h)
		if err != nil {
			t.Fatalf("(%v,%d,%d): %v", d, s, h, err)
		}
		want := eoqOracle(d, s, h)
		if math.Abs(got-want) > want*eoqTolerance {
			t.Fatalf("(%v,%d,%d) = %v, want %v", d, s, h, got, want)
		}
	}
}

func TestEconomicOrderQuantity_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = EconomicOrderQuantity(12000, 200, 1) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkEconomicOrderQuantity(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = EconomicOrderQuantity(12000, 200, 1)
	}
}
