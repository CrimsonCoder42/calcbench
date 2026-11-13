package invest

import (
	"fmt"
	"math"
	"testing"
)

func TestGarmanKlassVolatility_handWorked(t *testing.T) {
	got, err := GarmanKlassVolatility([]float64{100}, []float64{110}, []float64{100}, []float64{105})
	hl, co := math.Log(1.1), math.Log(1.05)
	want := math.Sqrt(0.5*hl*hl - (2*math.Ln2-1)*co*co)
	if err != nil || !near(got, want, 1e-16) || !near(got, 0.060187, 1e-5) {
		t.Errorf("GarmanKlassVolatility = %v, %v; want %v", got, err, want)
	}
}

// With open = close the estimator is √½·ln H/L per bar.
func TestGarmanKlassVolatility_openEqualsClose(t *testing.T) {
	got, err := GarmanKlassVolatility([]float64{105}, []float64{110}, []float64{100}, []float64{105})
	if want := math.Sqrt(0.5) * math.Log(1.1); err != nil || !near(got, want, 1e-16) {
		t.Errorf("O = C: %v, %v; want %v", got, err, want)
	}
}

// Every term is ≥ 0 on valid bars; the result is finite and non-negative on generated bars.
func TestGarmanKlassVolatility_generated(t *testing.T) {
	o, h, l, c := genBars(500)
	got, err := GarmanKlassVolatility(o, h, l, c)
	if err != nil || got <= 0 || got > 0.1 {
		t.Errorf("generated: %v, %v; want in (0, 0.1]", got, err)
	}
}

func TestGarmanKlassVolatility_errors(t *testing.T) {
	_, err := GarmanKlassVolatility(nil, nil, nil, nil)
	wantErr(t, "empty", err, ErrTooFew)
	_, err = GarmanKlassVolatility([]float64{100}, []float64{110}, []float64{100}, []float64{105, 106})
	wantErr(t, "lengths", err, ErrLengthMismatch)
	_, err = GarmanKlassVolatility([]float64{math.NaN()}, []float64{110}, []float64{100}, []float64{105})
	wantErr(t, "NaN", err, ErrNotFinite)
	_, err = GarmanKlassVolatility([]float64{120}, []float64{110}, []float64{100}, []float64{105})
	wantErr(t, "open > high", err, ErrDomain)
	_, err = GarmanKlassVolatility([]float64{100}, []float64{110}, []float64{100}, []float64{99})
	wantErr(t, "close < low", err, ErrDomain)
	_, err = GarmanKlassVolatility([]float64{100}, []float64{100}, []float64{110}, []float64{105})
	wantErr(t, "low > high", err, ErrDomain)
}

func TestGarmanKlassVolatility_noAlloc(t *testing.T) {
	o, h, l, c := []float64{100}, []float64{110}, []float64{100}, []float64{105}
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = GarmanKlassVolatility(o, h, l, c) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkGarmanKlassVolatility(b *testing.B) {
	for _, n := range benchSizes {
		o, h, l, c := genBars(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = GarmanKlassVolatility(o, h, l, c)
			}
		})
	}
}
