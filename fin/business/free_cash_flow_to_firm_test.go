package business

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// The CFA reading's summary states FCFF = EBIT(1 − t) + Dep − FCInv − WCInv without numbers.
// Hand-checked: EBIT $100,000, t = 0.25, Dep $10,000, FCInv $20,000, WCInv $5,000 →
// 75,000 + 10,000 − 20,000 − 5,000 = $60,000. Exact cents.
func TestFreeCashFlowToFirm_sourceExample(t *testing.T) {
	got, err := FreeCashFlowToFirm(10000000, 0.25, 1000000, 2000000, 500000)
	if err != nil || got != 6000000 {
		t.Fatalf("got %v, %v; want 6000000", got, err)
	}
}

func TestFreeCashFlowToFirm_table(t *testing.T) {
	cases := []struct {
		name           string
		ebit           int64
		tax            float64
		dep, capex, wc int64
		want           int64
	}{
		{"zero tax", 100, 0, 0, 0, 0, 100},
		{"full tax", 100, 1, 5, 0, 0, 5},
		{"loss", -100, 0.25, 0, 0, 0, -75},
		{"net disposals", 0, 0.25, 0, -50, 0, 50},
		{"working capital released", 0, 0.25, 0, 0, -20, 20},
		{"tie to even down", 1, 0.5, 0, 0, 0, 0},
		{"tie to even up", 3, 0.5, 0, 0, 0, 2},
		{"tie with odd offset", 1, 0.5, 1, 0, 0, 2},
		{"negative tie", -1, 0.5, 0, 0, 0, 0},
		{"negative tie up", -3, 0.5, 0, 0, 0, -2},
		{"max", maxExactCents, 0, 0, 0, 0, maxExactCents},
	}
	for _, c := range cases {
		got, err := FreeCashFlowToFirm(c.ebit, c.tax, c.dep, c.capex, c.wc)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestFreeCashFlowToFirm_errors(t *testing.T) {
	cases := []struct {
		name           string
		ebit           int64
		tax            float64
		dep, capex, wc int64
		want           error
	}{
		{"ebit above 2^53", maxExactCents + 1, 0, 0, 0, 0, ErrOutOfRange},
		{"dep above 2^53", 0, 0, maxExactCents + 1, 0, 0, ErrOutOfRange},
		{"capex below -2^53", 0, 0, 0, -maxExactCents - 1, 0, ErrOutOfRange},
		{"wc above 2^53", 0, 0, 0, 0, maxExactCents + 1, ErrOutOfRange},
		{"negative dep", 0, 0, -1, 0, 0, ErrNegative},
		{"NaN tax", 0, math.NaN(), 0, 0, 0, ErrNotFinite},
		{"Inf tax", 0, math.Inf(1), 0, 0, 0, ErrNotFinite},
		{"negative tax", 0, -0.1, 0, 0, 0, ErrRate},
		{"tax above 1", 0, 1.1, 0, 0, 0, ErrRate},
		{"result above 2^53", maxExactCents, 0, 1, 0, 0, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := FreeCashFlowToFirm(c.ebit, c.tax, c.dep, c.capex, c.wc); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Oracle at the full ±2⁵³ range: the result must equal the exact rational EBIT − EBIT·t + Dep − FCInv −
// WCInv rounded half-even (math/big), with no slack.
func TestFreeCashFlowToFirm_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	draw := func() int64 { return int64(rng.Uint64N(2*maxExactCents+1)) - maxExactCents }
	for i := 0; i < 500; i++ {
		ebit, capex, wc := draw(), draw(), draw()
		dep := int64(rng.Uint64N(maxExactCents + 1))
		tax := rng.Float64()
		if i%4 == 0 {
			ebit, dep, capex, wc, tax = int64(rng.IntN(2001))-1000, int64(rng.IntN(10)), 0, 0, 0.5
		}
		want := exactRound(ebit+dep-capex-wc, ebit, tax, -1)
		got, err := FreeCashFlowToFirm(ebit, tax, dep, capex, wc)
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

func TestFreeCashFlowToFirm_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = FreeCashFlowToFirm(10000000, 0.25, 1000000, 2000000, 500000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkFreeCashFlowToFirm(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = FreeCashFlowToFirm(10000000, 0.25, 1000000, 2000000, 500000)
	}
}
