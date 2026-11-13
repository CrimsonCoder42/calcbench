package business

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// The CFA reading's summary states FCFE = FCFF − Int(1 − t) + Net borrowing without numbers.
// Hand-checked: FCFF $60,000, interest $4,000, t = 0.25, net borrowing $2,000 → 60,000 − 3,000 + 2,000 =
// $59,000. Exact cents.
func TestFreeCashFlowToEquity_sourceExample(t *testing.T) {
	got, err := FreeCashFlowToEquity(6000000, 400000, 0.25, 200000)
	if err != nil || got != 5900000 {
		t.Fatalf("got %v, %v; want 5900000", got, err)
	}
}

func TestFreeCashFlowToEquity_table(t *testing.T) {
	cases := []struct {
		name           string
		fcff, interest int64
		tax            float64
		borrow         int64
		want           int64
	}{
		{"no interest", 100, 0, 0.3, 0, 100},
		{"full tax shield", 100, 50, 1, 0, 100},
		{"net repayment", 100, 0, 0, -40, 60},
		{"negative fcff", -100, 10, 0, 0, -110},
		{"tie to even", 0, 1, 0.5, 0, 0},
		{"tie to even down", 0, 3, 0.5, 0, -2},
		{"tie with odd offset", 1, 1, 0.5, 0, 0},
		{"extremes cancel", maxExactCents, 0, 0, -maxExactCents, 0},
	}
	for _, c := range cases {
		got, err := FreeCashFlowToEquity(c.fcff, c.interest, c.tax, c.borrow)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestFreeCashFlowToEquity_errors(t *testing.T) {
	cases := []struct {
		name           string
		fcff, interest int64
		tax            float64
		borrow         int64
		want           error
	}{
		{"fcff above 2^53", maxExactCents + 1, 0, 0, 0, ErrOutOfRange},
		{"interest above 2^53", 0, maxExactCents + 1, 0, 0, ErrOutOfRange},
		{"borrowing below -2^53", 0, 0, 0, -maxExactCents - 1, ErrOutOfRange},
		{"negative interest", 0, -1, 0, 0, ErrNegative},
		{"NaN tax", 0, 0, math.NaN(), 0, ErrNotFinite},
		{"-Inf tax", 0, 0, math.Inf(-1), 0, ErrNotFinite},
		{"tax above 1", 0, 0, 2, 0, ErrRate},
		{"result above 2^53", maxExactCents, 0, 0, 1, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := FreeCashFlowToEquity(c.fcff, c.interest, c.tax, c.borrow); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Oracle at the full ±2⁵³ range: exact rational FCFF − Int + Int·t + net borrowing rounded half-even.
func TestFreeCashFlowToEquity_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	draw := func() int64 { return int64(rng.Uint64N(2*maxExactCents+1)) - maxExactCents }
	for i := 0; i < 500; i++ {
		fcff, borrow := draw(), draw()
		interest := int64(rng.Uint64N(maxExactCents + 1))
		tax := rng.Float64()
		want := exactRound(fcff-interest+borrow, interest, tax, 1)
		got, err := FreeCashFlowToEquity(fcff, interest, tax, borrow)
		if want > maxExactCents || want < -maxExactCents {
			if !errors.Is(err, ErrOverflow) {
				t.Fatalf("i=%d: err = %v, want ErrOverflow", i, err)
			}
			continue
		}
		if err != nil || got != want {
			t.Fatalf("i=%d: got %v, %v; want %v", i, got, err, want)
		}
	}
}

func TestFreeCashFlowToEquity_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = FreeCashFlowToEquity(6000000, 400000, 0.25, 200000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkFreeCashFlowToEquity(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = FreeCashFlowToEquity(6000000, 400000, 0.25, 200000)
	}
}
