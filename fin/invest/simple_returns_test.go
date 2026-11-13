package invest

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
)

// Shared fixtures for the return-series tests (FIN-INVEST-032 onward).

var (
	sinkInt     int
	sinkFloats  []float64
	benchSizes  = []int{1_000, 10_000, 100_000}
	tailReturns = []float64{0.02, -0.05, 0.04, -0.01, 0.00, 0.03, -0.03, 0.01, 0.05, 0.02}
)

// genReturns returns n fixed-seed returns, uniform in [−0.05, 0.05).
func genReturns(n int) []float64 {
	r := rand.New(rand.NewPCG(1, 2))
	out := make([]float64, n)
	for i := range out {
		out[i] = 0.1*r.Float64() - 0.05
	}
	return out
}

// genPrices returns n fixed-seed positive prices, a random walk from 100.
func genPrices(n int) []float64 {
	out := make([]float64, n)
	p := 100.0
	for i, r := range genReturns(n) {
		p *= 1 + r
		out[i] = p
	}
	return out
}

// genBars returns n fixed-seed valid OHLC bars.
func genBars(n int) (open, high, low, close []float64) {
	r := rand.New(rand.NewPCG(3, 4))
	open, high, low, close = make([]float64, n), make([]float64, n), make([]float64, n), make([]float64, n)
	p := 100.0
	for i := range n {
		o := p
		c := o * (1 + 0.04*r.Float64() - 0.02)
		h := math.Max(o, c) * (1 + 0.01*r.Float64())
		l := math.Min(o, c) * (1 - 0.01*r.Float64())
		open[i], high[i], low[i], close[i] = o, h, l, c
		p = c
	}
	return open, high, low, close
}

func TestSimpleReturns_handWorked(t *testing.T) {
	got, err := SimpleReturns(nil, []float64{100, 110, 99})
	if err != nil || len(got) != 2 || !near(got[0], 0.1, 1e-15) || !near(got[1], -0.1, 1e-15) {
		t.Errorf("SimpleReturns = %v, %v; want [0.1 -0.1]", got, err)
	}
}

func TestSimpleReturns_reusesDst(t *testing.T) {
	dst := make([]float64, 0, 4)
	got, err := SimpleReturns(dst, []float64{1, 2, 4})
	if err != nil || &got[0] != &dst[:1][0] || got[0] != 1 || got[1] != 1 {
		t.Errorf("SimpleReturns = %v, %v; want [1 1] in dst", got, err)
	}
}

func TestSimpleReturns_errors(t *testing.T) {
	_, err := SimpleReturns(nil, []float64{100})
	wantErr(t, "one price", err, ErrTooFew)
	_, err = SimpleReturns(nil, nil)
	wantErr(t, "empty", err, ErrTooFew)
	_, err = SimpleReturns(nil, []float64{100, math.NaN()})
	wantErr(t, "NaN", err, ErrNotFinite)
	_, err = SimpleReturns(nil, []float64{100, math.Inf(1)})
	wantErr(t, "Inf", err, ErrNotFinite)
	_, err = SimpleReturns(nil, []float64{100, 0})
	wantErr(t, "zero", err, ErrDomain)
	_, err = SimpleReturns(nil, []float64{-1, 100})
	wantErr(t, "negative", err, ErrDomain)
	_, err = SimpleReturns(nil, []float64{1e-300, 1e300})
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestSimpleReturns_noAlloc(t *testing.T) {
	prices, dst := []float64{100, 110, 99}, make([]float64, 2)
	if a := testing.AllocsPerRun(100, func() { sinkFloats, sinkErr = SimpleReturns(dst, prices) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkSimpleReturns(b *testing.B) {
	for _, n := range benchSizes {
		prices, dst := genPrices(n), make([]float64, n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloats, sinkErr = SimpleReturns(dst, prices)
			}
		})
	}
}
