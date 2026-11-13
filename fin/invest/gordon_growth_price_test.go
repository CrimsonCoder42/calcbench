package invest

import (
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// OpenStax Principles of Finance §11.2: current dividend $5.00, growth 4 %, required return 8 % →
// D₁ = $5.20 and value $130.00.
func TestGordonGrowthPrice_openStax(t *testing.T) {
	got, err := GordonGrowthPrice(520, 0.08, 0.04)
	if err != nil || got != 13000 {
		t.Errorf("GordonGrowthPrice(520, 0.08, 0.04) = %d, %v; want 13000", got, err)
	}
}

func TestGordonGrowthPrice_cases(t *testing.T) {
	for _, c := range []struct {
		name string
		d    int64
		r, g float64
		want int64
	}{
		{"zero growth equals D/r", 200, 0.06, 0, 3333},
		{"negative growth", 1000, 0.10, -0.10, 5000},
		{"zero dividend", 0, 0.1, 0.05, 0},
		{"half-even tie down", 1, 1.5, -0.5, 0}, // 1/2 = 0.5 → 0
		{"half-even tie up", 3, 1.5, -0.5, 2},   // 3/2 = 1.5 → 2
	} {
		got, err := GordonGrowthPrice(c.d, c.r, c.g)
		if err != nil || got != c.want {
			t.Errorf("%s: %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

// Oracle: D/(r − g) in 256-bit big.Float from the exact float64 inputs, rounded half-even (F-5
// contract), on 300 random inputs.
func TestGordonGrowthPrice_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		d := r.Int64N(1_000_000)
		req := float64(r.IntN(3000)+1) / 10000
		g := req - float64(r.IntN(2000)+1)/10000
		got, err := GordonGrowthPrice(d, req, g)
		noErr(t, "GordonGrowthPrice", err)
		den := new(big.Float).SetPrec(256).Sub(big.NewFloat(req), big.NewFloat(g))
		q := new(big.Float).SetPrec(256).Quo(new(big.Float).SetPrec(256).SetInt64(d), den)
		exact, _ := q.Float64()
		checkRounded(t, "GordonGrowthPrice", got, exact)
	}
}

func TestGordonGrowthPrice_errors(t *testing.T) {
	_, err := GordonGrowthPrice(100, 0.05, 0.05)
	wantErr(t, "r = g", err, ErrDomain)
	_, err = GordonGrowthPrice(100, 0.04, 0.05)
	wantErr(t, "r < g", err, ErrDomain)
	_, err = GordonGrowthPrice(-1, 0.1, 0.05)
	wantErr(t, "dividend < 0", err, ErrNegative)
	_, err = GordonGrowthPrice(maxCents+1, 0.1, 0.05)
	wantErr(t, "dividend > 2^53", err, ErrOutOfRange)
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1), -1} {
		_, err = GordonGrowthPrice(100, bad, 0)
		wantErr(t, "bad r", err, ErrRate)
		_, err = GordonGrowthPrice(100, 0.1, bad)
		wantErr(t, "bad g", err, ErrRate)
	}
	_, err = GordonGrowthPrice(maxCents, 0.1, 0.1-1e-12)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestGordonGrowthPrice_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkCents, sinkErr = GordonGrowthPrice(520, 0.08, 0.04) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkGordonGrowthPrice(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, sinkErr = GordonGrowthPrice(520, 0.08, 0.04)
	}
}
